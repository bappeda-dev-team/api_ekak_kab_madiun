package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"ekak_kabupaten_madiun/model/web/tujuanopd"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type TujuanOpdServiceImpl struct {
	TujuanOpdRepository    repository.TujuanOpdRepository
	OpdRepository          repository.OpdRepository
	PeriodeRepository      repository.PeriodeRepository
	BidangUrusanRepository repository.BidangUrusanRepository
	LockDataRepository     repository.LockDataRepository
	DB                     *sql.DB
}

func NewTujuanOpdServiceImpl(tujuanOpdRepository repository.TujuanOpdRepository, opdRepository repository.OpdRepository, periodeRepository repository.PeriodeRepository, bidangUrusanRepository repository.BidangUrusanRepository, lockDataRepository repository.LockDataRepository, DB *sql.DB) *TujuanOpdServiceImpl {
	return &TujuanOpdServiceImpl{
		TujuanOpdRepository:    tujuanOpdRepository,
		OpdRepository:          opdRepository,
		PeriodeRepository:      periodeRepository,
		BidangUrusanRepository: bidangUrusanRepository,
		LockDataRepository:     lockDataRepository,
		DB:                     DB,
	}
}

func (service *TujuanOpdServiceImpl) Create(ctx context.Context, request tujuanopd.TujuanOpdCreateRequest) (tujuanopd.TujuanOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	periode, err := service.PeriodeRepository.FindById(ctx, tx, request.PeriodeId)
	if err != nil {
		if err == sql.ErrNoRows {
			return tujuanopd.TujuanOpdResponse{}, fmt.Errorf("periode dengan id %d tidak ditemukan", request.PeriodeId)
		}
		return tujuanopd.TujuanOpdResponse{}, err
	}
	tahunAwal, err := strconv.Atoi(periode.TahunAwal)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, fmt.Errorf("format tahun awal periode tidak valid: %s", periode.TahunAwal)
	}
	tahunAkhir, err := strconv.Atoi(periode.TahunAkhir)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, fmt.Errorf("format tahun akhir periode tidak valid: %s", periode.TahunAkhir)
	}
	_, err = service.BidangUrusanRepository.FindByKodeBidangUrusan(ctx, tx, request.KodeBidangUrusan)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}
	tujuanOpdDomain := domain.TujuanOpd{
		KodeOpd:          request.KodeOpd,
		KodeBidangUrusan: request.KodeBidangUrusan,
		Tujuan:           request.Tujuan,
		PeriodeId: domain.Periode{
			Id:           request.PeriodeId,
			TahunAwal:    periode.TahunAwal,
			TahunAkhir:   periode.TahunAkhir,
			JenisPeriode: periode.JenisPeriode,
		},
		TahunAwal:    periode.TahunAwal,
		TahunAkhir:   periode.TahunAkhir,
		JenisPeriode: periode.JenisPeriode,
	}
	for _, indikatorReq := range request.Indikator {
		uuidInd := uuid.New().String()[:5]
		kodeIndikator := fmt.Sprintf("IND-TJN-%s", uuidInd)
		indikatorDomain := domain.Indikator{
			KodeIndikator:       kodeIndikator,
			Jenis:               indikatorReq.Jenis,                                                    // FIX: mapping Jenis
			DefinisiOperasional: sql.NullString{String: indikatorReq.DefinisiOperasional, Valid: true}, // FIX
			Indikator:           indikatorReq.Indikator,
			RumusPerhitungan:    sql.NullString{String: indikatorReq.RumusPerhitungan, Valid: true},
			SumberData:          sql.NullString{String: indikatorReq.SumberData, Valid: true},
		}
		tahunMap := make(map[string]bool)
		if len(indikatorReq.Target) == 0 {
			return tujuanopd.TujuanOpdResponse{}, fmt.Errorf(
				"indikator harus memiliki minimal 1 target dalam rentang periode %d-%d",
				tahunAwal, tahunAkhir,
			)
		}
		for _, targetReq := range indikatorReq.Target {
			tahunTarget, err := strconv.Atoi(targetReq.Tahun)
			if err != nil {
				return tujuanopd.TujuanOpdResponse{}, fmt.Errorf("format tahun target tidak valid: %s", targetReq.Tahun)
			}
			if tahunTarget < tahunAwal || tahunTarget > tahunAkhir {
				return tujuanopd.TujuanOpdResponse{}, fmt.Errorf(
					"tahun target %d harus berada dalam rentang periode %d-%d",
					tahunTarget, tahunAwal, tahunAkhir,
				)
			}
			if tahunMap[targetReq.Tahun] {
				return tujuanopd.TujuanOpdResponse{}, fmt.Errorf("tahun target %s duplikat", targetReq.Tahun)
			}
			tahunMap[targetReq.Tahun] = true
			if targetReq.Target == "" {
				return tujuanopd.TujuanOpdResponse{}, fmt.Errorf("target untuk tahun %s tidak boleh kosong", targetReq.Tahun)
			}
			if targetReq.Satuan == "" {
				return tujuanopd.TujuanOpdResponse{}, fmt.Errorf("satuan untuk tahun %s tidak boleh kosong", targetReq.Tahun)
			}
			uuidTrg := uuid.New().String()[:5]
			targetDomain := domain.Target{
				Id:          fmt.Sprintf("TRG-TJN-%s", uuidTrg),
				IndikatorId: kodeIndikator, // FIX: pakai kodeIndikator, bukan Id yang kosong
				Target:      targetReq.Target,
				Satuan:      targetReq.Satuan,
				Tahun:       targetReq.Tahun,
			}
			indikatorDomain.Target = append(indikatorDomain.Target, targetDomain)
		}
		tujuanOpdDomain.Indikator = append(tujuanOpdDomain.Indikator, indikatorDomain)
	}
	tujuanOpdResult, err := service.TujuanOpdRepository.Create(ctx, tx, tujuanOpdDomain)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}
	return helper.ToTujuanOpdResponse(tujuanOpdResult), nil
}

func (service *TujuanOpdServiceImpl) Update(ctx context.Context, request tujuanopd.TujuanOpdUpdateRequest) (tujuanopd.TujuanOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	// Ambil existing untuk preserve tahun & kode_indikator
	existing, err := service.TujuanOpdRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}
	// Tahun/periode TIDAK diubah dari request — pakai existing
	// periode, err := service.PeriodeRepository.FindById(ctx, tx, request.PeriodeId) ← HAPUS/COMMENT
	tahunAwal, _ := strconv.Atoi(existing.TahunAwal)
	tahunAkhir, _ := strconv.Atoi(existing.TahunAkhir)
	// Build map indikator existing by Id
	existingIndMap := make(map[string]domain.Indikator)
	for _, ind := range existing.Indikator {
		existingIndMap[ind.Id] = ind
	}
	_, err = service.BidangUrusanRepository.FindByKodeBidangUrusan(ctx, tx, request.KodeBidangUrusan)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}
	tujuanOpd := domain.TujuanOpd{
		Id:               request.Id,
		KodeOpd:          request.KodeOpd,
		KodeBidangUrusan: request.KodeBidangUrusan,
		Tujuan:           request.Tujuan,
		// Periode tidak diubah — dari existing
		TahunAwal:    existing.TahunAwal,
		TahunAkhir:   existing.TahunAkhir,
		JenisPeriode: existing.JenisPeriode,
	}
	for _, indikatorReq := range request.Indikator {
		var kodeIndikator string
		if indikatorReq.Id != "" {
			// UPDATE existing — pertahankan kode_indikator lama
			ex, ok := existingIndMap[indikatorReq.Id]
			if !ok {
				return tujuanopd.TujuanOpdResponse{}, fmt.Errorf("indikator id %s tidak ditemukan", indikatorReq.Id)
			}
			kodeIndikator = ex.KodeIndikator
		} else {
			// INSERT baru — generate kode baru
			if indikatorReq.KodeIndikator != "" {
				kodeIndikator = indikatorReq.KodeIndikator
			} else {
				uuidInd := uuid.New().String()[:5]
				kodeIndikator = fmt.Sprintf("IND-TJN-%s", uuidInd)
			}
		}
		indikatorDomain := domain.Indikator{
			Id:                  indikatorReq.Id,
			KodeIndikator:       kodeIndikator,
			Jenis:               "renstra", // selalu renstra di endpoint ini
			Indikator:           indikatorReq.Indikator,
			RumusPerhitungan:    sql.NullString{String: indikatorReq.RumusPerhitungan, Valid: true},
			SumberData:          sql.NullString{String: indikatorReq.SumberData, Valid: true},
			DefinisiOperasional: sql.NullString{String: indikatorReq.DefinisiOperasional, Valid: true},
		}
		// Validasi & build target (sama seperti sebelumnya)
		tahunMap := make(map[string]bool)
		for _, targetReq := range indikatorReq.Target {
			tahunTarget, _ := strconv.Atoi(targetReq.Tahun)
			if tahunTarget < tahunAwal || tahunTarget > tahunAkhir {
				return tujuanopd.TujuanOpdResponse{}, fmt.Errorf(
					"tahun target %s di luar range periode %d-%d", targetReq.Tahun, tahunAwal, tahunAkhir)
			}
			if tahunMap[targetReq.Tahun] {
				return tujuanopd.TujuanOpdResponse{}, fmt.Errorf("tahun target %s duplikat", targetReq.Tahun)
			}
			tahunMap[targetReq.Tahun] = true
			indikatorDomain.Target = append(indikatorDomain.Target, domain.Target{
				Id:          targetReq.Id, // kosong = INSERT baru, isi = UPDATE
				IndikatorId: kodeIndikator,
				Target:      targetReq.Target,
				Satuan:      targetReq.Satuan,
				Tahun:       targetReq.Tahun,
				Jenis:       "renstra",
			})
		}
		tujuanOpd.Indikator = append(tujuanOpd.Indikator, indikatorDomain)
	}
	err = service.TujuanOpdRepository.Update(ctx, tx, tujuanOpd)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}
	// Ambil ulang dari DB untuk response yang akurat
	updated, err := service.TujuanOpdRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}
	return helper.ToTujuanOpdResponse(updated), nil
}

func (service *TujuanOpdServiceImpl) Delete(ctx context.Context, tujuanOpdId int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	_, err = service.TujuanOpdRepository.FindById(ctx, tx, tujuanOpdId)
	if err != nil {
		return err
	}

	return service.TujuanOpdRepository.Delete(ctx, tx, tujuanOpdId)
}

func (service *TujuanOpdServiceImpl) FindById(ctx context.Context, tujuanOpdId int) (tujuanopd.TujuanOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	tujuanOpd, err := service.TujuanOpdRepository.FindById(ctx, tx, tujuanOpdId)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}

	// Ambil data OPD
	opd, err := service.OpdRepository.FindByKodeOpd(ctx, tx, tujuanOpd.KodeOpd)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}

	// Ambil data bidang urusan
	bidangUrusan, err := service.BidangUrusanRepository.FindByKodeBidangUrusan(ctx, tx, tujuanOpd.KodeBidangUrusan)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}

	response := tujuanopd.TujuanOpdResponse{
		Id:               tujuanOpd.Id,
		KodeBidangUrusan: tujuanOpd.KodeBidangUrusan,
		NamaBidangUrusan: bidangUrusan.NamaBidangUrusan,
		KodeOpd:          tujuanOpd.KodeOpd,
		NamaOpd:          opd.NamaOpd,
		Tujuan:           tujuanOpd.Tujuan,
		TahunAwal:        tujuanOpd.TahunAwal,
		TahunAkhir:       tujuanOpd.TahunAkhir,
		JenisPeriode:     tujuanOpd.JenisPeriode,
		Indikator:        make([]tujuanopd.IndikatorResponse, 0),
	}

	for _, indikator := range tujuanOpd.Indikator {
		indikatorResponse := tujuanopd.IndikatorResponse{
			Id:                  indikator.Id,
			IdTujuanOpd:         tujuanOpd.Id,
			NamaIndikator:       indikator.Indikator,
			RumusPerhitungan:    indikator.RumusPerhitungan.String,
			DefinisiOperasional: indikator.DefinisiOperasional.String,
			SumberData:          indikator.SumberData.String,
			Target:              make([]tujuanopd.TargetResponse, 0),
		}

		tahunAwalInt, _ := strconv.Atoi(tujuanOpd.TahunAwal)
		tahunAkhirInt, _ := strconv.Atoi(tujuanOpd.TahunAkhir)

		// Buat map untuk target yang ada
		targetMap := make(map[string]domain.Target)
		for _, t := range indikator.Target {
			if t.Id != "" {
				targetMap[t.Tahun] = t
			}
		}

		// Generate target untuk setiap tahun dalam range
		for year := tahunAwalInt; year <= tahunAkhirInt; year++ {
			tahunStr := strconv.Itoa(year)
			if target, exists := targetMap[tahunStr]; exists {
				targetResponse := tujuanopd.TargetResponse{
					Id:              target.Id,
					IndikatorId:     indikator.KodeIndikator,
					Tahun:           tahunStr,
					TargetIndikator: target.Target,
					SatuanIndikator: target.Satuan,
				}
				indikatorResponse.Target = append(indikatorResponse.Target, targetResponse)
			} else {
				targetResponse := tujuanopd.TargetResponse{
					Id:              "",
					IndikatorId:     indikator.KodeIndikator,
					Tahun:           tahunStr,
					TargetIndikator: "",
					SatuanIndikator: "",
				}
				indikatorResponse.Target = append(indikatorResponse.Target, targetResponse)
			}
		}

		response.Indikator = append(response.Indikator, indikatorResponse)
	}

	return response, nil
}

func (service *TujuanOpdServiceImpl) FindAll(
	ctx context.Context,
	kodeOpd, tahunAwal, tahunAkhir, jenisPeriode string,
) ([]tujuanopd.TujuanOpdwithBidangUrusanResponse, error) {
	if len(tahunAwal) != 4 || len(tahunAkhir) != 4 {
		return nil, fmt.Errorf("format tahun tidak valid")
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	opd, err := service.OpdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return nil, err
	}
	// Gunakan FindAllByPeriod yang sudah pakai tb_indikator_matrix
	tujuanOpds, err := service.TujuanOpdRepository.FindAllByPeriod(
		ctx, tx, kodeOpd, tahunAwal, tahunAkhir, jenisPeriode, "renstra",
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]tujuanopd.TujuanOpdwithBidangUrusanResponse, 0), nil
		}
		return nil, err
	}
	bidangUrusanMap, err := service.fetchBidangUrusanMap(ctx, tx, tujuanOpds)
	if err != nil {
		return nil, err
	}
	return BuildTujuanOpdBidangResponse(tujuanOpds, opd, bidangUrusanMap, nil), nil
}

func (service *TujuanOpdServiceImpl) FindTujuanOpdOnlyName(ctx context.Context, kodeOpd string, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]tujuanopd.TujuanOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi tahun
	if len(tahunAwal) != 4 || len(tahunAkhir) != 4 {
		return nil, fmt.Errorf("format tahun tidak valid")
	}
	if _, err := strconv.Atoi(tahunAwal); err != nil {
		return nil, fmt.Errorf("tahun awal harus berupa angka")
	}
	if _, err := strconv.Atoi(tahunAkhir); err != nil {
		return nil, fmt.Errorf("tahun akhir harus berupa angka")
	}

	// Ambil data OPD
	opd, err := service.OpdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return nil, err
	}

	// Ambil semua tujuan OPD
	tujuanOpds, err := service.TujuanOpdRepository.FindAllByPeriod(ctx, tx, kodeOpd, tahunAwal, tahunAkhir, jenisPeriode, "renstra")
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]tujuanopd.TujuanOpdResponse, 0), nil
		}
		return nil, err
	}

	var responses []tujuanopd.TujuanOpdResponse
	for _, tujuan := range tujuanOpds {
		// Ambil data bidang urusan
		bidangUrusan, err := service.BidangUrusanRepository.FindByKodeBidangUrusan(ctx, tx, tujuan.KodeBidangUrusan)
		if err != nil {
			log.Printf("Warning: Gagal mendapatkan data bidang urusan: %v", err)
			continue
		}

		// var indikatorResponses []tujuanopd.IndikatorResponse
		// for _, indikator := range tujuan.Indikator {
		// 	var targetResponses []tujuanopd.TargetResponse
		// 	for _, target := range indikator.Target {
		// 		if target.Id != "" { // Hanya tambahkan target yang valid
		// 			targetResponses = append(targetResponses, tujuanopd.TargetResponse{
		// 				Id:              target.Id,
		// 				IndikatorId:     target.IndikatorId,
		// 				TargetIndikator: target.Target,
		// 				SatuanIndikator: target.Satuan,
		// 			})
		// 		}
		// 	}

		// 	indikatorResponses = append(indikatorResponses, tujuanopd.IndikatorResponse{
		// 		Id:            indikator.Id,
		// 		NamaIndikator: indikator.Indikator,
		// 		Target:        targetResponses,
		// 	})
		// }

		tujuanResponse := tujuanopd.TujuanOpdResponse{
			Id:               tujuan.Id,
			KodeBidangUrusan: tujuan.KodeBidangUrusan,
			NamaBidangUrusan: bidangUrusan.NamaBidangUrusan,
			KodeOpd:          tujuan.KodeOpd,
			NamaOpd:          opd.NamaOpd,
			Tujuan:           tujuan.Tujuan,
			TahunAwal:        tujuan.TahunAwal,
			TahunAkhir:       tujuan.TahunAkhir,
			JenisPeriode:     tujuan.JenisPeriode,
			// Indikator:        indikatorResponses,
		}

		responses = append(responses, tujuanResponse)
	}

	// Jika tidak ada data, kembalikan slice kosong
	if len(responses) == 0 {
	}

	// Urutkan berdasarkan ID
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].Id < responses[j].Id
	})

	return responses, nil
}

func (service *TujuanOpdServiceImpl) FindTujuanOpdByTahun(
	ctx context.Context,
	kodeOpd, tahun, jenisPeriode string,
) ([]tujuanopd.TujuanOpdwithBidangUrusanResponse, error) {
	if len(tahun) != 4 {
		return nil, fmt.Errorf("format tahun tidak valid")
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	opd, err := service.OpdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return nil, err
	}
	// Gunakan guard chain penetapan (paling lengkap: renstra→ranwal→rankhir→penetapan)
	tujuanOpds, err := service.loadLayerPenetapan(ctx, tx, kodeOpd, tahun, jenisPeriode)
	if err != nil {
		return nil, err
	}
	if len(tujuanOpds) == 0 {
		return make([]tujuanopd.TujuanOpdwithBidangUrusanResponse, 0), nil
	}
	bidangUrusanMap, err := service.fetchBidangUrusanMap(ctx, tx, tujuanOpds)
	if err != nil {
		return nil, err
	}
	return service.buildTujuanOpdResponse(tujuanOpds, opd, bidangUrusanMap), nil
}

// renstra renja
// ─────────────────────────────────────────────────────────────────
// HELPER: kumpulkan kode_bidang_urusan unik → batch fetch 1 query
// ─────────────────────────────────────────────────────────────────
func (service *TujuanOpdServiceImpl) fetchBidangUrusanMap(
	ctx context.Context,
	tx *sql.Tx,
	tujuanOpds []domain.TujuanOpd,
) (map[string]domainmaster.BidangUrusan, error) {
	uniqueKodes := make(map[string]struct{})
	for _, t := range tujuanOpds {
		if t.KodeBidangUrusan != "" {
			uniqueKodes[t.KodeBidangUrusan] = struct{}{}
		}
	}
	kodeList := make([]string, 0, len(uniqueKodes))
	for k := range uniqueKodes {
		kodeList = append(kodeList, k)
	}
	return service.TujuanOpdRepository.FindBidangUrusanBatch(ctx, tx, kodeList)
}

// ═════════════════════════════════════════════════════════════════
// GUARD CHAIN: renstra → ranwal → rankhir → penetapan
//
// Setiap layer memanggil layer sebelumnya sebagai base (fallback),
// lalu menimpa target per-indikator jika layer saat ini punya data.
// ═════════════════════════════════════════════════════════════════

// applyTargetOverride menerapkan target dari overrideTujuans ke atas baseTujuans.
// Untuk setiap indikator dalam setiap tujuan:
//   - Jika override punya target terisi (Id != "") → pakai target override.
//   - Jika tidak → pertahankan target dari base (fallback ke layer sebelumnya).
//
// Fungsi ini melakukan deep-copy agar tidak memodifikasi slice asli.
func applyTargetOverride(baseTujuans, overrideTujuans []domain.TujuanOpd) []domain.TujuanOpd {
	// Build lookup: tujuanId → kodeIndikator → target terisi
	type overrideKey struct {
		tujuanId      int
		kodeIndikator string
	}
	overrideMap := make(map[overrideKey]domain.Target, len(overrideTujuans)*4)
	for _, t := range overrideTujuans {
		for _, ind := range t.Indikator {
			for _, tg := range ind.Target {
				if tg.Id != "" {
					overrideMap[overrideKey{t.Id, ind.KodeIndikator}] = tg
					break
				}
			}
		}
	}

	result := make([]domain.TujuanOpd, len(baseTujuans))
	for i, tujuan := range baseTujuans {
		t := tujuan
		t.Indikator = make([]domain.Indikator, len(tujuan.Indikator))
		for j, ind := range tujuan.Indikator {
			newInd := ind
			key := overrideKey{tujuan.Id, ind.KodeIndikator}
			if tg, ok := overrideMap[key]; ok {
				newInd.Target = []domain.Target{tg}
			} else {
				// Salin slice target dari base agar tidak berbagi referensi
				newInd.Target = make([]domain.Target, len(ind.Target))
				copy(newInd.Target, ind.Target)
			}
			t.Indikator[j] = newInd
		}
		result[i] = t
	}
	return result
}

// loadLayerRenstra — layer dasar.
// Mengembalikan semua tujuan OPD dengan indikator renstra + target renstra
// (termasuk target dengan jenis=” untuk backward compat data lama).
func (service *TujuanOpdServiceImpl) loadLayerRenstra(
	ctx context.Context, tx *sql.Tx,
	kodeOpd, tahun, jenisPeriode string,
) ([]domain.TujuanOpd, error) {
	opds, err := service.TujuanOpdRepository.FindAllByTahun(ctx, tx, kodeOpd, tahun, jenisPeriode, "renstra")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return opds, nil
}

// loadLayerRanwal — guard level 2.
// Base = renstra. Target ranwal menimpa target renstra jika tersedia.
func (service *TujuanOpdServiceImpl) loadLayerRanwal(
	ctx context.Context, tx *sql.Tx,
	kodeOpd, tahun, jenisPeriode string,
) ([]domain.TujuanOpd, error) {
	base, err := service.loadLayerRenstra(ctx, tx, kodeOpd, tahun, jenisPeriode)
	if err != nil {
		return nil, err
	}
	ranwalOpds, err := service.TujuanOpdRepository.FindAllByTahun(ctx, tx, kodeOpd, tahun, jenisPeriode, "ranwal")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return applyTargetOverride(base, ranwalOpds), nil
}

// loadLayerRankhir — guard level 3.
// Base = ranwal (→ renstra). Target rankhir menimpa jika tersedia.
func (service *TujuanOpdServiceImpl) loadLayerRankhir(
	ctx context.Context, tx *sql.Tx,
	kodeOpd, tahun, jenisPeriode string,
) ([]domain.TujuanOpd, error) {
	base, err := service.loadLayerRanwal(ctx, tx, kodeOpd, tahun, jenisPeriode)
	if err != nil {
		return nil, err
	}
	rankhirOpds, err := service.TujuanOpdRepository.FindAllByTahun(ctx, tx, kodeOpd, tahun, jenisPeriode, "rankhir")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return applyTargetOverride(base, rankhirOpds), nil
}

// loadLayerPenetapan — guard level 4 (paling atas).
// Base = rankhir (→ ranwal → renstra). Target penetapan menimpa jika tersedia.
func (service *TujuanOpdServiceImpl) loadLayerPenetapan(
	ctx context.Context, tx *sql.Tx,
	kodeOpd, tahun, jenisPeriode string,
) ([]domain.TujuanOpd, error) {
	base, err := service.loadLayerRankhir(ctx, tx, kodeOpd, tahun, jenisPeriode)
	if err != nil {
		return nil, err
	}
	penetapanOpds, err := service.TujuanOpdRepository.FindAllByTahun(ctx, tx, kodeOpd, tahun, jenisPeriode, "penetapan")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return applyTargetOverride(base, penetapanOpds), nil
}

// ═════════════════════════════════════════════════════════════════
// DUAL-SLOT HELPERS (mirip Tujuan Pemda loadLayerRankhirDual / loadLayerPenetapanDual)
//
// Setiap indikator mendapat tepat 2 slot target untuk tahun yang diminta:
//   rankhir  → [ranwal, rankhir]
//   penetapan→ [rankhir, penetapan]
// ═════════════════════════════════════════════════════════════════

type opdIndikatorKey struct {
	tujuanId      int
	kodeIndikator string
}

// fillOpdTargetSlot mengisi slot jenis tertentu di base dari data layerData DB.
func fillOpdTargetSlot(base *[]domain.TujuanOpd, layerData []domain.TujuanOpd, jenis string) {
	lookup := make(map[opdIndikatorKey]domain.Target)
	for _, tp := range layerData {
		for _, ind := range tp.Indikator {
			for _, tg := range ind.Target {
				raw := strings.TrimSpace(tg.Target)
				if raw != "" && raw != "-" {
					lookup[opdIndikatorKey{tp.Id, ind.KodeIndikator}] = tg
					break
				}
			}
		}
	}
	if len(lookup) == 0 {
		return
	}
	for i := range *base {
		tp := &(*base)[i]
		for j := range tp.Indikator {
			ind := &tp.Indikator[j]
			k := opdIndikatorKey{tp.Id, ind.KodeIndikator}
			tg, ok := lookup[k]
			if !ok {
				continue
			}
			for t := range ind.Target {
				if ind.Target[t].Jenis == jenis {
					ind.Target[t] = tg
					ind.Target[t].Jenis = jenis
					break
				}
			}
		}
	}
}

// loadLayerRankhirDualOpd — skeleton dari renstra, lalu 2 slot: [ranwal, rankhir].
// Slot ranwal diisi dari data ranwal DB; jika tidak ada, fallback ke target renstra.
// Slot rankhir diisi dari data rankhir DB; jika tidak ada, placeholder "-".
func (service *TujuanOpdServiceImpl) loadLayerRankhirDualOpd(
	ctx context.Context, tx *sql.Tx,
	kodeOpd, tahun, jenisPeriode string,
) ([]domain.TujuanOpd, error) {
	renstraOpds, err := service.TujuanOpdRepository.FindAllByTahun(ctx, tx, kodeOpd, tahun, jenisPeriode, "renstra")
	if err != nil {
		return nil, err
	}
	for i := range renstraOpds {
		for j := range renstraOpds[i].Indikator {
			ind := &renstraOpds[i].Indikator[j]
			kode := ind.KodeIndikator
			// Slot 1: ranwal — pakai renstra target sebagai awal (akan di-override jika ranwal ada)
			var ranwalSlot domain.Target
			if len(ind.Target) > 0 && strings.TrimSpace(ind.Target[0].Target) != "" && ind.Target[0].Target != "-" {
				tg := ind.Target[0]
				ranwalSlot = domain.Target{
					Id: tg.Id, IndikatorId: kode,
					Target: tg.Target, Satuan: tg.Satuan, Tahun: tahun, Jenis: "ranwal",
				}
			} else {
				ranwalSlot = domain.Target{IndikatorId: kode, Target: "-", Satuan: "-", Tahun: tahun, Jenis: "ranwal"}
			}
			// Slot 2: rankhir placeholder
			rankhirSlot := domain.Target{IndikatorId: kode, Target: "-", Satuan: "-", Tahun: tahun, Jenis: "rankhir"}
			ind.Target = []domain.Target{ranwalSlot, rankhirSlot}
		}
	}
	// Override slot ranwal dengan data aktual ranwal dari DB (jika ada)
	ranwalOpds, err := service.TujuanOpdRepository.FindAllByTahun(ctx, tx, kodeOpd, tahun, jenisPeriode, "ranwal")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	fillOpdTargetSlot(&renstraOpds, ranwalOpds, "ranwal")
	// Isi slot rankhir dari DB
	rankhirOpds, err := service.TujuanOpdRepository.FindAllByTahun(ctx, tx, kodeOpd, tahun, jenisPeriode, "rankhir")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	fillOpdTargetSlot(&renstraOpds, rankhirOpds, "rankhir")
	return renstraOpds, nil
}

// loadLayerPenetapanDualOpd — skeleton dari renstra, lalu 2 slot: [rankhir, penetapan].
func (service *TujuanOpdServiceImpl) loadLayerPenetapanDualOpd(
	ctx context.Context, tx *sql.Tx,
	kodeOpd, tahun, jenisPeriode string,
) ([]domain.TujuanOpd, error) {
	renstraOpds, err := service.TujuanOpdRepository.FindAllByTahun(ctx, tx, kodeOpd, tahun, jenisPeriode, "renstra")
	if err != nil {
		return nil, err
	}
	for i := range renstraOpds {
		for j := range renstraOpds[i].Indikator {
			ind := &renstraOpds[i].Indikator[j]
			kode := ind.KodeIndikator
			ind.Target = []domain.Target{
				{IndikatorId: kode, Target: "-", Satuan: "-", Tahun: tahun, Jenis: "rankhir"},
				{IndikatorId: kode, Target: "-", Satuan: "-", Tahun: tahun, Jenis: "penetapan"},
			}
		}
	}
	rankhirOpds, err := service.TujuanOpdRepository.FindAllByTahun(ctx, tx, kodeOpd, tahun, jenisPeriode, "rankhir")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	fillOpdTargetSlot(&renstraOpds, rankhirOpds, "rankhir")
	penetapanOpds, err := service.TujuanOpdRepository.FindAllByTahun(ctx, tx, kodeOpd, tahun, jenisPeriode, "penetapan")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	fillOpdTargetSlot(&renstraOpds, penetapanOpds, "penetapan")
	return renstraOpds, nil
}

// ─────────────────────────────────────────────────────────────────
// HELPER: bangun response TujuanOpdwithBidangUrusanResponse
//
//	dari domain, opd, dan bidangUrusanMap (sudah di-batch)
func (service *TujuanOpdServiceImpl) buildTujuanOpdResponse(
	tujuanOpds []domain.TujuanOpd,
	opd domainmaster.Opd,
	bidangUrusanMap map[string]domainmaster.BidangUrusan,
) []tujuanopd.TujuanOpdwithBidangUrusanResponse {
	return BuildTujuanOpdBidangResponse(tujuanOpds, opd, bidangUrusanMap, nil)
}

// ─────────────────────────────────────────────────────────────────
// GET /tujuan_opd/renstra/:kode_opd/:tahun_awal/:tahun_akhir
// jenis indikator hardcode = "renstra"
// target: slot setiap tahun dalam range
// ─────────────────────────────────────────────────────────────────
func (service *TujuanOpdServiceImpl) FindTujuanRenstra(
	ctx context.Context,
	kodeOpd, tahunAwal, tahunAkhir, jenisPeriode string,
) ([]tujuanopd.TujuanOpdwithBidangUrusanResponse, error) {
	if len(tahunAwal) != 4 || len(tahunAkhir) != 4 {
		return nil, fmt.Errorf("format tahun tidak valid")
	}
	if _, err := strconv.Atoi(tahunAwal); err != nil {
		return nil, fmt.Errorf("tahun awal harus berupa angka")
	}
	if _, err := strconv.Atoi(tahunAkhir); err != nil {
		return nil, fmt.Errorf("tahun akhir harus berupa angka")
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	opd, err := service.OpdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return nil, err
	}
	tujuanOpds, err := service.TujuanOpdRepository.FindAllByPeriod(
		ctx, tx, kodeOpd, tahunAwal, tahunAkhir, jenisPeriode, "renstra",
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]tujuanopd.TujuanOpdwithBidangUrusanResponse, 0), nil
		}
		return nil, err
	}
	bidangUrusanMap, err := service.fetchBidangUrusanMap(ctx, tx, tujuanOpds)
	if err != nil {
		return nil, err
	}
	return BuildTujuanOpdBidangResponse(tujuanOpds, opd, bidangUrusanMap, &TujuanOpdBidangResponseOpts{ForceIndicatorJenisRankhir: true}), nil
}

// ─────────────────────────────────────────────────────────────────
// GET /tujuan_opd/renja_ranwal/:kode_opd/:tahun
// jenis indikator hardcode = "ranwal"
// target: 1 slot untuk tahun yang diminta
// ─────────────────────────────────────────────────────────────────
// FindTujuanRanwal menampilkan tujuan OPD dengan guard chain renstra → ranwal.
// Target ranwal digunakan jika tersedia; jika tidak, fallback ke target renstra.
func (service *TujuanOpdServiceImpl) FindTujuanRanwal(
	ctx context.Context,
	kodeOpd, tahun, jenisPeriode string,
) ([]tujuanopd.TujuanOpdwithBidangUrusanResponse, error) {
	if len(tahun) != 4 {
		return nil, fmt.Errorf("format tahun tidak valid")
	}
	if _, err := strconv.Atoi(tahun); err != nil {
		return nil, fmt.Errorf("tahun harus berupa angka")
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	opd, err := service.OpdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return nil, err
	}
	// Guard: renstra (default) → ranwal (override)
	tujuanOpds, err := service.loadLayerRanwal(ctx, tx, kodeOpd, tahun, jenisPeriode)
	if err != nil {
		return nil, err
	}
	if len(tujuanOpds) == 0 {
		return make([]tujuanopd.TujuanOpdwithBidangUrusanResponse, 0), nil
	}
	bidangUrusanMap, err := service.fetchBidangUrusanMap(ctx, tx, tujuanOpds)
	if err != nil {
		return nil, err
	}
	return service.buildTujuanOpdResponse(tujuanOpds, opd, bidangUrusanMap), nil
}

// ─────────────────────────────────────────────────────────────────
// GET /tujuan_opd/rankhir/:kode_opd/:tahun
// jenis indikator hardcode = "rankhir"
// target: 1 slot untuk tahun yang diminta
// ─────────────────────────────────────────────────────────────────

// FindTujuanRankhir menampilkan tujuan OPD dengan 2 slot target per indikator:
// [ranwal, rankhir]. Indikator dari renstra; target masing-masing diisi dari DB.
func (service *TujuanOpdServiceImpl) FindTujuanRankhir(
	ctx context.Context,
	kodeOpd, tahun, jenisPeriode string,
) ([]tujuanopd.TujuanOpdwithBidangUrusanResponse, error) {
	if len(tahun) != 4 {
		return nil, fmt.Errorf("format tahun tidak valid")
	}
	if _, err := strconv.Atoi(tahun); err != nil {
		return nil, fmt.Errorf("tahun harus berupa angka")
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	opd, err := service.OpdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return nil, err
	}
	tujuanOpds, err := service.loadLayerRankhirDualOpd(ctx, tx, kodeOpd, tahun, jenisPeriode)
	if err != nil {
		return nil, err
	}
	if len(tujuanOpds) == 0 {
		return make([]tujuanopd.TujuanOpdwithBidangUrusanResponse, 0), nil
	}
	bidangUrusanMap, err := service.fetchBidangUrusanMap(ctx, tx, tujuanOpds)
	if err != nil {
		return nil, err
	}
	return BuildTujuanOpdRankhirDualSlotResponse(tujuanOpds, opd, bidangUrusanMap), nil
}

func (service *TujuanOpdServiceImpl) CreateTujuanRenjaIndikator(
	ctx context.Context,
	tujuanOpdId int,
	jenis string,
	requests []tujuanopd.IndikatorCreateRequest,
) ([]tujuanopd.IndikatorResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	_, err = service.TujuanOpdRepository.FindById(ctx, tx, tujuanOpdId)
	if err != nil {
		return nil, fmt.Errorf("tujuan opd id %d tidak ditemukan", tujuanOpdId)
	}
	var indikatorDomains []domain.Indikator
	var responses []tujuanopd.IndikatorResponse
	for _, req := range requests {
		if req.Indikator == "" {
			return nil, fmt.Errorf("nama indikator tidak boleh kosong")
		}
		if len(req.Target) != 1 {
			return nil, fmt.Errorf("setiap indikator harus memiliki tepat 1 target")
		}
		if req.Target[0].Target == "" {
			return nil, fmt.Errorf("nilai target tidak boleh kosong")
		}
		if req.Target[0].Satuan == "" {
			return nil, fmt.Errorf("satuan tidak boleh kosong")
		}
		if req.Target[0].Tahun == "" {
			return nil, fmt.Errorf("tahun target tidak boleh kosong")
		}
		kodeIndikator := fmt.Sprintf("IND-TJN-%s", uuid.New().String()[:5])
		targetId := fmt.Sprintf("TRG-TJN-%s", uuid.New().String()[:5])
		ind := domain.Indikator{
			KodeIndikator:       kodeIndikator,
			Jenis:               jenis,
			DefinisiOperasional: sql.NullString{String: req.DefinisiOperasional, Valid: true},
			Indikator:           req.Indikator,
			RumusPerhitungan:    sql.NullString{String: req.RumusPerhitungan, Valid: true},
			SumberData:          sql.NullString{String: req.SumberData, Valid: true},
			Target: []domain.Target{{
				Id:          targetId,
				IndikatorId: kodeIndikator,
				Target:      req.Target[0].Target,
				Satuan:      req.Target[0].Satuan,
				Tahun:       req.Target[0].Tahun,
			}},
		}
		indikatorDomains = append(indikatorDomains, ind)
		responses = append(responses, tujuanopd.IndikatorResponse{
			Id:                  kodeIndikator,
			KodeIndikator:       kodeIndikator,
			IdTujuanOpd:         tujuanOpdId,
			NamaIndikator:       req.Indikator,
			RumusPerhitungan:    req.RumusPerhitungan,
			SumberData:          req.SumberData,
			DefinisiOperasional: req.DefinisiOperasional,
			Jenis:               jenis,
			Target: []tujuanopd.TargetResponse{{
				Id: targetId, IndikatorId: kodeIndikator,
				Tahun: req.Target[0].Tahun, TargetIndikator: req.Target[0].Target,
				SatuanIndikator: req.Target[0].Satuan,
			}},
		})
	}
	if err := service.TujuanOpdRepository.CreateRenjaIndikator(ctx, tx, tujuanOpdId, indikatorDomains); err != nil {
		return nil, err
	}
	return responses, nil
}

func (service *TujuanOpdServiceImpl) UpdateTujuanRenjaIndikator(
	ctx context.Context,
	kodeIndikator string, // ← dari URL param
	jenis string,
	request tujuanopd.IndikatorUpdateRequest, // ← single object
) (tujuanopd.IndikatorResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanopd.IndikatorResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	// Validasi: pastikan kode_indikator ada di DB
	_, err = service.TujuanOpdRepository.FindIndikatorByKodeIndikator(ctx, tx, kodeIndikator)
	if err != nil {
		return tujuanopd.IndikatorResponse{}, fmt.Errorf("indikator dengan kode %s tidak ditemukan", kodeIndikator)
	}
	// Validasi field wajib
	if request.Indikator == "" {
		return tujuanopd.IndikatorResponse{}, fmt.Errorf("nama indikator tidak boleh kosong")
	}
	if len(request.Target) != 1 {
		return tujuanopd.IndikatorResponse{}, fmt.Errorf("harus memiliki tepat 1 target")
	}
	if request.Target[0].Target == "" {
		return tujuanopd.IndikatorResponse{}, fmt.Errorf("nilai target tidak boleh kosong")
	}
	if request.Target[0].Tahun == "" {
		return tujuanopd.IndikatorResponse{}, fmt.Errorf("tahun target tidak boleh kosong")
	}
	targetId := request.Target[0].Id
	if targetId == "" {
		targetId = fmt.Sprintf("TRG-TJN-%s", uuid.New().String()[:5])
	}
	ind := domain.Indikator{
		KodeIndikator:       kodeIndikator, // pakai dari URL, bukan dari body
		Jenis:               jenis,
		DefinisiOperasional: sql.NullString{String: request.DefinisiOperasional, Valid: true},
		Indikator:           request.Indikator,
		RumusPerhitungan:    sql.NullString{String: request.RumusPerhitungan, Valid: true},
		SumberData:          sql.NullString{String: request.SumberData, Valid: true},
		Target: []domain.Target{{
			Id: targetId, IndikatorId: kodeIndikator,
			Target: request.Target[0].Target,
			Satuan: request.Target[0].Satuan,
			Tahun:  request.Target[0].Tahun,
		}},
	}
	if err := service.TujuanOpdRepository.UpdateRenjaIndikator(ctx, tx, []domain.Indikator{ind}); err != nil {
		return tujuanopd.IndikatorResponse{}, err
	}
	return tujuanopd.IndikatorResponse{
		Id:                  kodeIndikator,
		KodeIndikator:       kodeIndikator,
		NamaIndikator:       request.Indikator,
		RumusPerhitungan:    request.RumusPerhitungan,
		SumberData:          request.SumberData,
		DefinisiOperasional: request.DefinisiOperasional,
		Jenis:               jenis,
		Target: []tujuanopd.TargetResponse{{
			Id: targetId, IndikatorId: kodeIndikator,
			Tahun: request.Target[0].Tahun, TargetIndikator: request.Target[0].Target,
			SatuanIndikator: request.Target[0].Satuan,
		}},
	}, nil
}
func (service *TujuanOpdServiceImpl) DeleteTujuanRenjaIndikator(ctx context.Context, kodeIndikator string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)
	_, err = service.TujuanOpdRepository.FindIndikatorByKodeIndikator(ctx, tx, kodeIndikator)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("kode indikator %s tidak ditemukan", kodeIndikator)
		}
		return err // ← tampilkan error asli (bukan dibungkus)
	}
	return service.TujuanOpdRepository.DeleteIndikatorTargetRenja(ctx, tx, kodeIndikator)
}

func (service *TujuanOpdServiceImpl) FindTujuanPenetapan(
	ctx context.Context,
	kodeOpd, tahun, jenisPeriode string,
) ([]tujuanopd.TujuanOpdwithBidangUrusanResponse, error) {
	if len(tahun) != 4 {
		return nil, fmt.Errorf("format tahun tidak valid")
	}
	if _, err := strconv.Atoi(tahun); err != nil {
		return nil, fmt.Errorf("tahun harus berupa angka")
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	opd, err := service.OpdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return nil, err
	}
	tujuanOpds, err := service.TujuanOpdRepository.FindAllByTahun(
		ctx, tx, kodeOpd, tahun, jenisPeriode, "penetapan",
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]tujuanopd.TujuanOpdwithBidangUrusanResponse, 0), nil
		}
		return nil, err
	}
	bidangUrusanMap, err := service.fetchBidangUrusanMap(ctx, tx, tujuanOpds)
	if err != nil {
		return nil, err
	}
	return service.buildTujuanOpdResponse(tujuanOpds, opd, bidangUrusanMap), nil
}

const lockJenisTujuanOpd = "tujuan_opd"

// TujuanOpdPenetapan menampilkan tujuan OPD dengan 2 slot target per indikator:
// [rankhir, penetapan]. Indikator dari renstra; target diisi dari DB masing-masing.
// Field is_lock mencerminkan status lock dari tb_lock_data.
func (service *TujuanOpdServiceImpl) TujuanOpdPenetapan(
	ctx context.Context, kodeOpd, tahun, jenisPeriode string,
) ([]tujuanopd.TujuanOpdPenetapanResponse, error) {
	if len(tahun) != 4 {
		return nil, fmt.Errorf("format tahun tidak valid")
	}
	if _, err := strconv.Atoi(tahun); err != nil {
		return nil, fmt.Errorf("tahun harus berupa angka")
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	opd, err := service.OpdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return nil, err
	}
	isLocked, err := service.LockDataRepository.IsLocked(ctx, tx, lockJenisTujuanOpd, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	tujuanOpds, err := service.loadLayerPenetapanDualOpd(ctx, tx, kodeOpd, tahun, jenisPeriode)
	if err != nil {
		return nil, err
	}
	if len(tujuanOpds) == 0 {
		return []tujuanopd.TujuanOpdPenetapanResponse{}, nil
	}
	bidangUrusanMap, err := service.fetchBidangUrusanMap(ctx, tx, tujuanOpds)
	if err != nil {
		return nil, err
	}
	return buildTujuanOpdPenetapanDualSlotResponse(tujuanOpds, opd, bidangUrusanMap, isLocked), nil
}

func BuildTujuanOpdBidangResponseDualJenis(
	tujuanOpds []domain.TujuanOpd,
	opd domainmaster.Opd,
	bidangUrusanMap map[string]domainmaster.BidangUrusan,
	jenisOverride string,
) []tujuanopd.TujuanOpdwithBidangUrusanResponse {
	responseMap := make(map[string]*tujuanopd.TujuanOpdwithBidangUrusanResponse)
	for _, tujuan := range tujuanOpds {
		indikatorResponses := make([]tujuanopd.IndikatorResponse, 0, len(tujuan.Indikator))
		for _, ind := range tujuan.Indikator {
			// Target override (rankhir/penetapan)
			targets := make([]tujuanopd.TargetResponse, 0, len(ind.Target))
			for _, t := range ind.Target {
				targets = append(targets, tujuanopd.TargetResponse{
					Id: t.Id, IndikatorId: ind.KodeIndikator,
					Tahun: t.Tahun, TargetIndikator: t.Target,
					SatuanIndikator: t.Satuan, Jenis: jenisOverride,
				})
			}
			// Target renstra (full period sebagai perbandingan)
			targetRenstra := make([]tujuanopd.TargetResponse, 0, len(ind.TargetRenstra))
			for _, t := range ind.TargetRenstra {
				targetRenstra = append(targetRenstra, tujuanopd.TargetResponse{
					Id: t.Id, IndikatorId: ind.KodeIndikator,
					Tahun: t.Tahun, TargetIndikator: t.Target,
					SatuanIndikator: t.Satuan, Jenis: "renstra",
				})
			}
			indikatorResponses = append(indikatorResponses, tujuanopd.IndikatorResponse{
				Id:                  ind.Id,
				KodeIndikator:       ind.KodeIndikator,
				IdTujuanOpd:         tujuan.Id,
				NamaIndikator:       ind.Indikator,
				RumusPerhitungan:    ind.RumusPerhitungan.String,
				SumberData:          ind.SumberData.String,
				DefinisiOperasional: ind.DefinisiOperasional.String,
				Jenis:               jenisOverride,
				Target:              targets,
				TargetRenstra:       targetRenstra, // ← 2 jenis
			})
		}
		tujuanResp := tujuanopd.TujuanOpdResponse{
			Id: tujuan.Id, Tujuan: tujuan.Tujuan,
			TahunAwal: tujuan.TahunAwal, TahunAkhir: tujuan.TahunAkhir,
			JenisPeriode: tujuan.JenisPeriode, Indikator: indikatorResponses,
		}
		mapKey := tujuan.KodeBidangUrusan
		if mapKey == "" {
			mapKey = "000"
		}
		if existing, ok := responseMap[mapKey]; ok {
			existing.TujuanOpd = append(existing.TujuanOpd, tujuanResp)
		} else {
			bu := bidangUrusanMap[tujuan.KodeBidangUrusan]
			kodeUrusan := ""
			if len(bu.KodeBidangUrusan) > 0 {
				kodeUrusan = bu.KodeBidangUrusan[:1]
			}
			responseMap[mapKey] = &tujuanopd.TujuanOpdwithBidangUrusanResponse{
				Urusan: bu.NamaUrusan, KodeUrusan: kodeUrusan,
				KodeBidangUrusan: bu.KodeBidangUrusan, NamaBidangUrusan: bu.NamaBidangUrusan,
				KodeOpd: tujuan.KodeOpd, NamaOpd: opd.NamaOpd,
				TujuanOpd: []tujuanopd.TujuanOpdResponse{tujuanResp},
			}
		}
	}
	result := make([]tujuanopd.TujuanOpdwithBidangUrusanResponse, 0, len(responseMap))
	for _, r := range responseMap {
		result = append(result, *r)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].KodeBidangUrusan < result[j].KodeBidangUrusan
	})
	return result
}

// buildTujuanOpdPenetapanDualJenis — khusus penetapan dengan is_lock + 2 jenis target
func buildTujuanOpdPenetapanDualJenis(
	tujuanOpds []domain.TujuanOpd,
	opd domainmaster.Opd,
	bidangUrusanMap map[string]domainmaster.BidangUrusan,
	isLock bool,
	jenisOverride string,
) []tujuanopd.TujuanOpdPenetapanResponse {
	responseMap := make(map[string]*tujuanopd.TujuanOpdPenetapanResponse)
	for _, tujuan := range tujuanOpds {
		indikatorResponses := make([]tujuanopd.IndikatorResponse, 0, len(tujuan.Indikator))
		for _, ind := range tujuan.Indikator {
			targets := make([]tujuanopd.TargetResponse, 0, len(ind.Target))
			for _, t := range ind.Target {
				targets = append(targets, tujuanopd.TargetResponse{
					Id: t.Id, IndikatorId: ind.KodeIndikator,
					Tahun: t.Tahun, TargetIndikator: t.Target,
					SatuanIndikator: t.Satuan, Jenis: jenisOverride,
				})
			}
			targetRenstra := make([]tujuanopd.TargetResponse, 0, len(ind.TargetRenstra))
			for _, t := range ind.TargetRenstra {
				targetRenstra = append(targetRenstra, tujuanopd.TargetResponse{
					Id: t.Id, IndikatorId: ind.KodeIndikator,
					Tahun: t.Tahun, TargetIndikator: t.Target,
					SatuanIndikator: t.Satuan, Jenis: "renstra",
				})
			}
			indikatorResponses = append(indikatorResponses, tujuanopd.IndikatorResponse{
				Id: ind.Id, KodeIndikator: ind.KodeIndikator, IdTujuanOpd: tujuan.Id,
				NamaIndikator:       ind.Indikator,
				RumusPerhitungan:    ind.RumusPerhitungan.String,
				SumberData:          ind.SumberData.String,
				DefinisiOperasional: ind.DefinisiOperasional.String,
				Jenis:               jenisOverride,
				Target:              targets,
				TargetRenstra:       targetRenstra,
			})
		}
		tujuanResp := tujuanopd.TujuanOpdResponse{
			Id: tujuan.Id, Tujuan: tujuan.Tujuan,
			TahunAwal: tujuan.TahunAwal, TahunAkhir: tujuan.TahunAkhir,
			JenisPeriode:   tujuan.JenisPeriode,
			JenisPenetapan: "penetapan_perencanaan",
			Indikator:      indikatorResponses,
		}
		mapKey := tujuan.KodeBidangUrusan
		if mapKey == "" {
			mapKey = "000"
		}
		if existing, ok := responseMap[mapKey]; ok {
			existing.TujuanOpd = append(existing.TujuanOpd, tujuanResp)
		} else {
			bu := bidangUrusanMap[tujuan.KodeBidangUrusan]
			kodeUrusan := ""
			if len(bu.KodeBidangUrusan) > 0 {
				kodeUrusan = bu.KodeBidangUrusan[:1]
			}
			responseMap[mapKey] = &tujuanopd.TujuanOpdPenetapanResponse{
				Urusan: bu.NamaUrusan, KodeUrusan: kodeUrusan,
				KodeBidangUrusan: bu.KodeBidangUrusan, NamaBidangUrusan: bu.NamaBidangUrusan,
				KodeOpd: tujuan.KodeOpd, NamaOpd: opd.NamaOpd,
				IsLock:    isLock,
				TujuanOpd: []tujuanopd.TujuanOpdResponse{tujuanResp},
			}
		}
	}
	result := make([]tujuanopd.TujuanOpdPenetapanResponse, 0, len(responseMap))
	for _, r := range responseMap {
		result = append(result, *r)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].KodeBidangUrusan < result[j].KodeBidangUrusan
	})
	return result
}


func (service *TujuanOpdServiceImpl) SetTujuanOpdLocked(
	ctx context.Context, tujuanOpdId int, locked bool,
) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)
	return service.TujuanOpdRepository.SetTujuanOpdLocked(ctx, tx, tujuanOpdId, locked)
}

// buildTujuanOpdPenetapanFromDB membangun []TujuanOpdPenetapanResponse dari hasil guard chain.
// Struktur respons sama dengan BuildTujuanOpdBidangResponse tetapi tambahan field is_lock
// dan jenis_penetapan pada TujuanOpdResponse.
func buildTujuanOpdPenetapanFromDB(
	tujuanOpds []domain.TujuanOpd,
	opd domainmaster.Opd,
	bidangUrusanMap map[string]domainmaster.BidangUrusan,
	isLock bool,
) []tujuanopd.TujuanOpdPenetapanResponse {
	responseMap := make(map[string]*tujuanopd.TujuanOpdPenetapanResponse)
	for _, tujuan := range tujuanOpds {
		indikatorResponses := make([]tujuanopd.IndikatorResponse, 0, len(tujuan.Indikator))
		for _, ind := range tujuan.Indikator {
			indikatorResponses = append(indikatorResponses, buildIndikatorResponseItem(ind, tujuan.Id, nil))
		}
		tujuanResp := tujuanopd.TujuanOpdResponse{
			Id:             tujuan.Id,
			Tujuan:         tujuan.Tujuan,
			TahunAwal:      tujuan.TahunAwal,
			TahunAkhir:     tujuan.TahunAkhir,
			JenisPeriode:   tujuan.JenisPeriode,
			JenisPenetapan: "penetapan_perencanaan",
			Indikator:      indikatorResponses,
		}
		mapKey := tujuan.KodeBidangUrusan
		if mapKey == "" {
			mapKey = "000"
		}
		if existing, ok := responseMap[mapKey]; ok {
			existing.TujuanOpd = append(existing.TujuanOpd, tujuanResp)
		} else {
			bu := bidangUrusanMap[tujuan.KodeBidangUrusan]
			kodeUrusan := ""
			if len(bu.KodeBidangUrusan) > 0 {
				kodeUrusan = bu.KodeBidangUrusan[:1]
			}
			responseMap[mapKey] = &tujuanopd.TujuanOpdPenetapanResponse{
				Urusan:           bu.NamaUrusan,
				KodeUrusan:       kodeUrusan,
				KodeBidangUrusan: bu.KodeBidangUrusan,
				NamaBidangUrusan: bu.NamaBidangUrusan,
				KodeOpd:          tujuan.KodeOpd,
				NamaOpd:          opd.NamaOpd,
				IsLock:           isLock,
				TujuanOpd:        []tujuanopd.TujuanOpdResponse{tujuanResp},
			}
		}
	}
	result := make([]tujuanopd.TujuanOpdPenetapanResponse, 0, len(responseMap))
	for _, r := range responseMap {
		result = append(result, *r)
	}
	return result
}

// validateOpdLayerJenis memastikan jenis = ranwal / rankhir / penetapan.
func validateOpdLayerJenis(jenis string) (string, error) {
	jenis = strings.TrimSpace(jenis)
	if jenis != "ranwal" && jenis != "rankhir" && jenis != "penetapan" {
		return "", fmt.Errorf("jenis layer tidak valid: '%s'. Gunakan: ranwal, rankhir, atau penetapan", jenis)
	}
	return jenis, nil
}

// CreateTargetOpdLayer — insert target baru untuk layer ranwal/rankhir/penetapan.
// Indikator renstra harus sudah ada. Gagal jika target (kode_indikator+tahun+jenis) duplikat.
func (service *TujuanOpdServiceImpl) CreateTargetOpdLayer(
	ctx context.Context,
	jenis string,
	req tujuanopd.LayerTargetBatchRequest,
) ([]tujuanopd.TargetResponse, error) {
	var err error
	jenis, err = validateOpdLayerJenis(jenis)
	if err != nil {
		return nil, err
	}
	if len(req.Targets) == 0 {
		return nil, fmt.Errorf("targets tidak boleh kosong")
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	responses := make([]tujuanopd.TargetResponse, 0, len(req.Targets))
	for _, item := range req.Targets {
		kode := strings.TrimSpace(item.KodeIndikator)
		tahun := strings.TrimSpace(item.Tahun)
		if kode == "" {
			return nil, fmt.Errorf("kode_indikator tidak boleh kosong")
		}
		if tahun == "" {
			return nil, fmt.Errorf("tahun tidak boleh kosong untuk kode_indikator %s", kode)
		}
		// Indikator renstra harus sudah ada
		_, err := service.TujuanOpdRepository.FindIndikatorByKodeIndikator(ctx, tx, kode)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("indikator '%s' tidak ditemukan di renstra", kode)
		}
		if err != nil {
			return nil, err
		}
		// Tolak duplikat
		exists, err := service.TujuanOpdRepository.TargetOpdExistsByKey(ctx, tx, kode, tahun, jenis)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf(
				"target kode_indikator '%s' tahun %s jenis %s sudah ada, gunakan update",
				kode, tahun, jenis,
			)
		}
		id := fmt.Sprintf("TRG-OPD-L-%s", uuid.New().String()[:8])
		saved, err := service.TujuanOpdRepository.CreateTargetOpdSingle(ctx, tx, domain.Target{
			Id:          id,
			IndikatorId: kode,
			Target:      item.Target,
			Satuan:      item.Satuan,
			Tahun:       tahun,
			Jenis:       jenis,
		})
		if err != nil {
			return nil, err
		}
		responses = append(responses, tujuanopd.TargetResponse{
			Id:              saved.Id,
			IndikatorId:     saved.IndikatorId,
			Tahun:           saved.Tahun,
			TargetIndikator: saved.Target,
			SatuanIndikator: saved.Satuan,
			Jenis:           saved.Jenis,
		})
	}
	return responses, nil
}

// UpdateTargetOpdLayer — update target yang sudah ada berdasarkan ID.
// Jenis target di DB harus cocok dengan parameter jenis di endpoint.
func (service *TujuanOpdServiceImpl) UpdateTargetOpdLayer(
	ctx context.Context,
	jenis string,
	req tujuanopd.LayerTargetUpdateBatchRequest,
) ([]tujuanopd.TargetResponse, error) {
	var err error
	jenis, err = validateOpdLayerJenis(jenis)
	if err != nil {
		return nil, err
	}
	if len(req.Targets) == 0 {
		return nil, fmt.Errorf("targets tidak boleh kosong")
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	responses := make([]tujuanopd.TargetResponse, 0, len(req.Targets))
	for _, item := range req.Targets {
		id := strings.TrimSpace(item.Id)
		if id == "" {
			return nil, fmt.Errorf("id target wajib diisi")
		}
		existing, err := service.TujuanOpdRepository.FindTargetOpdById(ctx, tx, id)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("target id '%s' tidak ditemukan", id)
		}
		if err != nil {
			return nil, err
		}
		if existing.Jenis != jenis {
			return nil, fmt.Errorf(
				"target id '%s' adalah jenis '%s', tidak bisa diupdate via endpoint jenis '%s'",
				id, existing.Jenis, jenis,
			)
		}
		saved, err := service.TujuanOpdRepository.UpdateTargetOpdById(ctx, tx, id, item.Target, item.Satuan)
		if err != nil {
			return nil, err
		}
		responses = append(responses, tujuanopd.TargetResponse{
			Id:              saved.Id,
			IndikatorId:     saved.IndikatorId,
			Tahun:           saved.Tahun,
			TargetIndikator: saved.Target,
			SatuanIndikator: saved.Satuan,
			Jenis:           saved.Jenis,
		})
	}
	return responses, nil
}

// DeleteTargetOpdLayer — hapus target berdasarkan kode_indikator + tahun + jenis.
func (service *TujuanOpdServiceImpl) DeleteTargetOpdLayer(
	ctx context.Context,
	kodeIndikator, tahun, jenis string,
) error {
	var err error
	jenis, err = validateOpdLayerJenis(jenis)
	if err != nil {
		return err
	}
	kodeIndikator = strings.TrimSpace(kodeIndikator)
	tahun = strings.TrimSpace(tahun)
	if kodeIndikator == "" || tahun == "" {
		return fmt.Errorf("kode_indikator dan tahun harus diisi")
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)
	return service.TujuanOpdRepository.DeleteTargetOpdByJenis(ctx, tx, kodeIndikator, tahun, jenis)
}

func (service *TujuanOpdServiceImpl) LockTujuanOpd(ctx context.Context, kodeOpd, tahun string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)
	return service.LockDataRepository.Lock(ctx, tx, lockJenisTujuanOpd, kodeOpd, tahun)
}
func (service *TujuanOpdServiceImpl) UnlockTujuanOpd(ctx context.Context, kodeOpd, tahun string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)
	return service.LockDataRepository.Unlock(ctx, tx, lockJenisTujuanOpd, kodeOpd, tahun)
}

// rankhiir builder
type TujuanOpdBidangResponseOpts struct {
	// ForceIndicatorJenisRankhir = true → kolom "jenis" di JSON selalu "rankhir";
	// nilai asli di database (ranwal/rankhir/penetapan) dipindah ke "sumber_jenis".
	ForceIndicatorJenisRankhir bool
}

// ─────────────────────────────────────────────────────────────────
// Helper murni (pure function, tidak bergantung receiver).
// ─────────────────────────────────────────────────────────────────
// buildIndikatorJenisDisplay menentukan nilai "jenis" dan "sumber_jenis" untuk JSON.
func buildIndikatorJenisDisplay(dbJenis string, opts *TujuanOpdBidangResponseOpts) (jenisJSON, sumberJenisJSON string) {
	db := strings.TrimSpace(dbJenis)
	if opts != nil && opts.ForceIndicatorJenisRankhir {
		return "rankhir", db
	}
	return db, ""
}

// buildIndikatorDualSlotResponse mengonversi domain.Indikator dengan 2-slot target ke
// IndikatorResponse dengan field terpisah berdasarkan jenis:
//   rankhir view  → TargetRanwal  + TargetRankhir
//   penetapan view→ TargetRankhir + TargetPenetapan
func buildIndikatorDualSlotResponse(ind domain.Indikator, tujuanId int) tujuanopd.IndikatorResponse {
	resp := tujuanopd.IndikatorResponse{
		Id:                  ind.Id,
		KodeIndikator:       ind.KodeIndikator,
		IdTujuanOpd:         tujuanId,
		NamaIndikator:       ind.Indikator,
		RumusPerhitungan:    ind.RumusPerhitungan.String,
		SumberData:          ind.SumberData.String,
		DefinisiOperasional: ind.DefinisiOperasional.String,
		Jenis:               "renstra",
	}
	for _, t := range ind.Target {
		tr := tujuanopd.TargetResponse{
			Id:              t.Id,
			IndikatorId:     ind.KodeIndikator,
			Tahun:           t.Tahun,
			TargetIndikator: t.Target,
			SatuanIndikator: t.Satuan,
			Jenis:           t.Jenis,
		}
		switch strings.TrimSpace(t.Jenis) {
		case "ranwal":
			resp.TargetRanwal = append(resp.TargetRanwal, tr)
		case "rankhir":
			resp.TargetRankhir = append(resp.TargetRankhir, tr)
		case "penetapan":
			resp.TargetPenetapan = append(resp.TargetPenetapan, tr)
		default:
			resp.Target = append(resp.Target, tr)
		}
	}
	return resp
}

// BuildTujuanOpdRankhirDualSlotResponse — builder untuk rankhir view (field terpisah).
func BuildTujuanOpdRankhirDualSlotResponse(
	tujuanOpds []domain.TujuanOpd,
	opd domainmaster.Opd,
	bidangUrusanMap map[string]domainmaster.BidangUrusan,
) []tujuanopd.TujuanOpdwithBidangUrusanResponse {
	responseMap := make(map[string]*tujuanopd.TujuanOpdwithBidangUrusanResponse)
	for _, tujuan := range tujuanOpds {
		indikatorResponses := make([]tujuanopd.IndikatorResponse, 0, len(tujuan.Indikator))
		for _, ind := range tujuan.Indikator {
			indikatorResponses = append(indikatorResponses, buildIndikatorDualSlotResponse(ind, tujuan.Id))
		}
		tujuanResp := tujuanopd.TujuanOpdResponse{
			Id: tujuan.Id, Tujuan: tujuan.Tujuan,
			TahunAwal: tujuan.TahunAwal, TahunAkhir: tujuan.TahunAkhir,
			JenisPeriode: tujuan.JenisPeriode, Indikator: indikatorResponses,
		}
		mapKey := tujuan.KodeBidangUrusan
		if mapKey == "" {
			mapKey = "000"
		}
		if existing, ok := responseMap[mapKey]; ok {
			existing.TujuanOpd = append(existing.TujuanOpd, tujuanResp)
		} else {
			bu := bidangUrusanMap[tujuan.KodeBidangUrusan]
			kodeUrusan := ""
			if len(bu.KodeBidangUrusan) > 0 {
				kodeUrusan = bu.KodeBidangUrusan[:1]
			}
			responseMap[mapKey] = &tujuanopd.TujuanOpdwithBidangUrusanResponse{
				Urusan: bu.NamaUrusan, KodeUrusan: kodeUrusan,
				KodeBidangUrusan: bu.KodeBidangUrusan, NamaBidangUrusan: bu.NamaBidangUrusan,
				KodeOpd: tujuan.KodeOpd, NamaOpd: opd.NamaOpd,
				TujuanOpd: []tujuanopd.TujuanOpdResponse{tujuanResp},
			}
		}
	}
	result := make([]tujuanopd.TujuanOpdwithBidangUrusanResponse, 0, len(responseMap))
	for _, r := range responseMap {
		result = append(result, *r)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].KodeBidangUrusan < result[j].KodeBidangUrusan
	})
	return result
}

// buildTujuanOpdPenetapanDualSlotResponse — builder untuk penetapan view (field terpisah).
func buildTujuanOpdPenetapanDualSlotResponse(
	tujuanOpds []domain.TujuanOpd,
	opd domainmaster.Opd,
	bidangUrusanMap map[string]domainmaster.BidangUrusan,
	isLock bool,
) []tujuanopd.TujuanOpdPenetapanResponse {
	responseMap := make(map[string]*tujuanopd.TujuanOpdPenetapanResponse)
	for _, tujuan := range tujuanOpds {
		indikatorResponses := make([]tujuanopd.IndikatorResponse, 0, len(tujuan.Indikator))
		for _, ind := range tujuan.Indikator {
			indikatorResponses = append(indikatorResponses, buildIndikatorDualSlotResponse(ind, tujuan.Id))
		}
		tujuanResp := tujuanopd.TujuanOpdResponse{
			Id: tujuan.Id, Tujuan: tujuan.Tujuan,
			TahunAwal: tujuan.TahunAwal, TahunAkhir: tujuan.TahunAkhir,
			JenisPeriode:   tujuan.JenisPeriode,
			JenisPenetapan: "penetapan_perencanaan",
			Indikator:      indikatorResponses,
		}
		mapKey := tujuan.KodeBidangUrusan
		if mapKey == "" {
			mapKey = "000"
		}
		if existing, ok := responseMap[mapKey]; ok {
			existing.TujuanOpd = append(existing.TujuanOpd, tujuanResp)
		} else {
			bu := bidangUrusanMap[tujuan.KodeBidangUrusan]
			kodeUrusan := ""
			if len(bu.KodeBidangUrusan) > 0 {
				kodeUrusan = bu.KodeBidangUrusan[:1]
			}
			responseMap[mapKey] = &tujuanopd.TujuanOpdPenetapanResponse{
				Urusan: bu.NamaUrusan, KodeUrusan: kodeUrusan,
				KodeBidangUrusan: bu.KodeBidangUrusan, NamaBidangUrusan: bu.NamaBidangUrusan,
				KodeOpd: tujuan.KodeOpd, NamaOpd: opd.NamaOpd,
				IsLock:    isLock,
				TujuanOpd: []tujuanopd.TujuanOpdResponse{tujuanResp},
			}
		}
	}
	result := make([]tujuanopd.TujuanOpdPenetapanResponse, 0, len(responseMap))
	for _, r := range responseMap {
		result = append(result, *r)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].KodeBidangUrusan < result[j].KodeBidangUrusan
	})
	return result
}

// buildIndikatorResponseItem mengonversi satu domain.Indikator + target → IndikatorResponse.
func buildIndikatorResponseItem(ind domain.Indikator, tujuanId int, opts *TujuanOpdBidangResponseOpts) tujuanopd.IndikatorResponse {
	jenisJSON, sumberJenisJSON := buildIndikatorJenisDisplay(ind.Jenis, opts)
	targets := make([]tujuanopd.TargetResponse, 0, len(ind.Target))
	for _, t := range ind.Target {
		targets = append(targets, tujuanopd.TargetResponse{
			Id:              t.Id,
			IndikatorId:     ind.KodeIndikator,
			Tahun:           t.Tahun,
			TargetIndikator: t.Target,
			SatuanIndikator: t.Satuan,
			Jenis:           t.Jenis,
		})
	}
	return tujuanopd.IndikatorResponse{
		Id:                  ind.Id,
		KodeIndikator:       ind.KodeIndikator,
		IdTujuanOpd:         tujuanId,
		NamaIndikator:       ind.Indikator,
		RumusPerhitungan:    ind.RumusPerhitungan.String,
		SumberData:          ind.SumberData.String,
		DefinisiOperasional: ind.DefinisiOperasional.String,
		Jenis:               jenisJSON,
		SumberJenis:         sumberJenisJSON,
		Target:              targets,
	}
}

// ─────────────────────────────────────────────────────────────────
// Builder utama (standalone function, dipanggil dari method service).
// ─────────────────────────────────────────────────────────────────
// BuildTujuanOpdBidangResponse membangun []TujuanOpdwithBidangUrusanResponse
// dari slice domain, opd, dan bidangUrusanMap.
// opts boleh nil (perilaku default: jenis = dari DB, sumber_jenis kosong).
func BuildTujuanOpdBidangResponse(
	tujuanOpds []domain.TujuanOpd,
	opd domainmaster.Opd,
	bidangUrusanMap map[string]domainmaster.BidangUrusan,
	opts *TujuanOpdBidangResponseOpts,
) []tujuanopd.TujuanOpdwithBidangUrusanResponse {
	responseMap := make(map[string]*tujuanopd.TujuanOpdwithBidangUrusanResponse)
	for _, tujuan := range tujuanOpds {
		indikatorResponses := make([]tujuanopd.IndikatorResponse, 0, len(tujuan.Indikator))
		for _, ind := range tujuan.Indikator {
			indikatorResponses = append(indikatorResponses, buildIndikatorResponseItem(ind, tujuan.Id, opts))
		}
		tujuanResp := tujuanopd.TujuanOpdResponse{
			Id:           tujuan.Id,
			Tujuan:       tujuan.Tujuan,
			TahunAwal:    tujuan.TahunAwal,
			TahunAkhir:   tujuan.TahunAkhir,
			JenisPeriode: tujuan.JenisPeriode,
			Indikator:    indikatorResponses,
		}
		mapKey := tujuan.KodeBidangUrusan
		if mapKey == "" {
			mapKey = "000"
		}
		if existing, ok := responseMap[mapKey]; ok {
			existing.TujuanOpd = append(existing.TujuanOpd, tujuanResp)
		} else {
			bu := bidangUrusanMap[tujuan.KodeBidangUrusan]
			kodeUrusan := ""
			if len(bu.KodeBidangUrusan) > 0 {
				kodeUrusan = bu.KodeBidangUrusan[:1]
			}
			responseMap[mapKey] = &tujuanopd.TujuanOpdwithBidangUrusanResponse{
				Urusan:           bu.NamaUrusan,
				KodeUrusan:       kodeUrusan,
				KodeBidangUrusan: bu.KodeBidangUrusan,
				NamaBidangUrusan: bu.NamaBidangUrusan,
				KodeOpd:          tujuan.KodeOpd,
				NamaOpd:          opd.NamaOpd,
				TujuanOpd:        []tujuanopd.TujuanOpdResponse{tujuanResp},
			}
		}
	}
	result := make([]tujuanopd.TujuanOpdwithBidangUrusanResponse, 0, len(responseMap))
	for _, r := range responseMap {
		result = append(result, *r)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].KodeBidangUrusan < result[j].KodeBidangUrusan
	})
	return result
}
