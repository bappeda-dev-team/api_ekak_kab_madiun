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

	"github.com/google/uuid"
)

type TujuanOpdServiceImpl struct {
	TujuanOpdRepository    repository.TujuanOpdRepository
	OpdRepository          repository.OpdRepository
	PeriodeRepository      repository.PeriodeRepository
	BidangUrusanRepository repository.BidangUrusanRepository
	DB                     *sql.DB
}

func NewTujuanOpdServiceImpl(tujuanOpdRepository repository.TujuanOpdRepository, opdRepository repository.OpdRepository, periodeRepository repository.PeriodeRepository, bidangUrusanRepository repository.BidangUrusanRepository, DB *sql.DB) *TujuanOpdServiceImpl {
	return &TujuanOpdServiceImpl{
		TujuanOpdRepository:    tujuanOpdRepository,
		OpdRepository:          opdRepository,
		PeriodeRepository:      periodeRepository,
		BidangUrusanRepository: bidangUrusanRepository,
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
	_, err = service.TujuanOpdRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
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
		// var kodeIndikator string
		// if indikatorReq.KodeIndikator != "" {
		// 	kodeIndikator = indikatorReq.KodeIndikator
		// } else {
		// 	uuidInd := uuid.New().String()[:5]
		// 	kodeIndikator = fmt.Sprintf("IND-TJN-%s", uuidInd)
		// }
		var kodeIndikator string
		if indikatorReq.KodeIndikator != "" {
			kodeIndikator = indikatorReq.KodeIndikator // hanya pakai ini
		} else {
			uuidInd := uuid.New().String()[:5]
			kodeIndikator = fmt.Sprintf("IND-TJN-%s", uuidInd) // generate baru!
		}
		indikatorDomain := domain.Indikator{
			KodeIndikator:       kodeIndikator,
			Jenis:               indikatorReq.Jenis, // sudah ada di UpdateRequest
			DefinisiOperasional: sql.NullString{String: indikatorReq.DefinisiOperasional, Valid: true},
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
			var targetId string
			if targetReq.Id != "" {
				targetId = targetReq.Id
			} else {
				uuidTrg := uuid.New().String()[:5]
				targetId = fmt.Sprintf("TRG-TJN-%s", uuidTrg)
			}
			targetDomain := domain.Target{
				Id:          targetId,
				IndikatorId: kodeIndikator, // FIX: pakai kodeIndikator, bukan indikatorDomain.Id
				Target:      targetReq.Target,
				Satuan:      targetReq.Satuan,
				Tahun:       targetReq.Tahun,
			}
			indikatorDomain.Target = append(indikatorDomain.Target, targetDomain)
		}
		tujuanOpd.Indikator = append(tujuanOpd.Indikator, indikatorDomain)
	}
	err = service.TujuanOpdRepository.Update(ctx, tx, tujuanOpd)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}
	return helper.ToTujuanOpdResponse(tujuanOpd), nil
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

	tujuanOpd, err := service.TujuanOpdRepository.FindByIdOnly(ctx, tx, tujuanOpdId)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}

	tujuanOpdIds := []int{tujuanOpd.Id}
	indikatorTujuan, err := service.getIndikatorWithFallback(ctx, tx, tujuanOpdIds)
	if err != nil {
		log.Printf("ERROR service.getIndikatorWithFallback: %w", err)
		return tujuanopd.TujuanOpdResponse{}, err
	}

	tujuanOpd.Indikator = indikatorTujuan

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

func (service *TujuanOpdServiceImpl) FindAll(ctx context.Context, kodeOpd string, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]tujuanopd.TujuanOpdwithBidangUrusanResponse, error) {
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
	tujuanOpds, err := service.TujuanOpdRepository.FindAllOnly(ctx, tx, kodeOpd, tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]tujuanopd.TujuanOpdwithBidangUrusanResponse, 0), nil
		}
		return nil, err
	}
	tujuanOpdIds := []int{}
	for _, tuj := range tujuanOpds {
		tujuanOpdIds = append(tujuanOpdIds, tuj.Id)
	}

	indikatorTujuan, err := service.getIndikatorWithFallback(ctx, tx, tujuanOpdIds)
	if err != nil {
		log.Printf("ERROR service.getIndikatorWithFallback: %v", err)
		return nil, err
	}
	indTujuanById := make(map[int][]domain.Indikator)
	for _, ind := range indikatorTujuan {
		indTujuanById[ind.TujuanOpdId] = append(indTujuanById[ind.TujuanOpdId], ind)
	}

	for i := range tujuanOpds {
		id := tujuanOpds[i].Id

		if inds, ok := indTujuanById[id]; ok {
			tujuanOpds[i].Indikator = inds
		}
	}

	// Buat map untuk mengelompokkan response berdasarkan kode_bidang_urusan
	responseMap := make(map[string]*tujuanopd.TujuanOpdwithBidangUrusanResponse)

	for _, tujuan := range tujuanOpds {
		// Ambil data bidang urusan
		bidangUrusan, err := service.BidangUrusanRepository.FindByKodeBidangUrusan(ctx, tx, tujuan.KodeBidangUrusan)
		if err != nil {
			return nil, err
		}

		tujuanResponse := tujuanopd.TujuanOpdResponse{
			Id:           tujuan.Id,
			Tujuan:       tujuan.Tujuan,
			TahunAwal:    tujuan.TahunAwal,
			TahunAkhir:   tujuan.TahunAkhir,
			JenisPeriode: tujuan.JenisPeriode,
			Indikator:    make([]tujuanopd.IndikatorResponse, 0),
		}

		// Proses indikator dan target seperti sebelumnya
		for _, indikator := range tujuan.Indikator {
			indikatorResponse := tujuanopd.IndikatorResponse{
				Id:                  indikator.Id,
				IdTujuanOpd:         tujuan.Id,
				NamaIndikator:       indikator.Indikator,
				DefinisiOperasional: indikator.DefinisiOperasional.String,
				RumusPerhitungan:    indikator.RumusPerhitungan.String,
				SumberData:          indikator.SumberData.String,
				Target:              make([]tujuanopd.TargetResponse, 0),
			}

			tahunAwalInt, _ := strconv.Atoi(tujuan.TahunAwal)
			tahunAkhirInt, _ := strconv.Atoi(tujuan.TahunAkhir)

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
						IndikatorId:     indikator.Id,
						Tahun:           tahunStr,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					}
					indikatorResponse.Target = append(indikatorResponse.Target, targetResponse)
				} else {
					targetResponse := tujuanopd.TargetResponse{
						Id:              "",
						IndikatorId:     indikator.Id,
						Tahun:           tahunStr,
						TargetIndikator: "",
						SatuanIndikator: "",
					}
					indikatorResponse.Target = append(indikatorResponse.Target, targetResponse)
				}
			}

			tujuanResponse.Indikator = append(tujuanResponse.Indikator, indikatorResponse)
		}

		// Cek apakah sudah ada entry untuk kode_bidang_urusan ini
		mapKey := tujuan.KodeBidangUrusan
		if mapKey == "" {
			mapKey = "000" // Gunakan key default untuk bidang urusan kosong
		}

		if existing, exists := responseMap[mapKey]; exists {
			// Jika sudah ada, tambahkan tujuan ke array tujuan yang ada
			existing.TujuanOpd = append(existing.TujuanOpd, tujuanResponse)
		} else {
			// Jika belum ada, buat entry baru
			kodeUrusan := ""
			if len(bidangUrusan.KodeBidangUrusan) > 0 {
				kodeUrusan = bidangUrusan.KodeBidangUrusan[:1]
			}

			responseMap[mapKey] = &tujuanopd.TujuanOpdwithBidangUrusanResponse{
				Urusan:           bidangUrusan.NamaUrusan,
				KodeUrusan:       kodeUrusan,
				KodeBidangUrusan: bidangUrusan.KodeBidangUrusan,
				NamaBidangUrusan: bidangUrusan.NamaBidangUrusan,
				KodeOpd:          tujuan.KodeOpd,
				NamaOpd:          opd.NamaOpd,
				TujuanOpd:        []tujuanopd.TujuanOpdResponse{tujuanResponse},
			}
		}
	}

	// Convert map to slice
	var responses []tujuanopd.TujuanOpdwithBidangUrusanResponse
	for _, response := range responseMap {
		responses = append(responses, *response)
	}

	// Sort responses berdasarkan kode_bidang_urusan
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].KodeBidangUrusan < responses[j].KodeBidangUrusan
	})

	if len(responses) == 0 {
		responses = make([]tujuanopd.TujuanOpdwithBidangUrusanResponse, 0)
	}

	return responses, nil
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
	tujuanOpds, err := service.TujuanOpdRepository.FindAll(ctx, tx, kodeOpd, tahunAwal, tahunAkhir, jenisPeriode)
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

func (service *TujuanOpdServiceImpl) FindTujuanOpdByTahun(ctx context.Context, kodeOpd string, tahun string, jenisPeriode string) ([]tujuanopd.TujuanOpdwithBidangUrusanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi tahun
	if len(tahun) != 4 {
		return nil, fmt.Errorf("format tahun tidak valid")
	}
	if _, err := strconv.Atoi(tahun); err != nil {
		return nil, fmt.Errorf("tahun harus berupa angka")
	}

	// Ambil data OPD
	opd, err := service.OpdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return nil, err
	}

	// Ambil tujuan OPD berdasarkan tahun
	tujuanOpds, err := service.TujuanOpdRepository.FindTujuanOpdByTahun(ctx, tx, kodeOpd, tahun, jenisPeriode)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]tujuanopd.TujuanOpdwithBidangUrusanResponse, 0), nil
		}
		return nil, err
	}

	// Buat map untuk mengelompokkan response berdasarkan kode_bidang_urusan
	responseMap := make(map[string]*tujuanopd.TujuanOpdwithBidangUrusanResponse)

	for _, tujuan := range tujuanOpds {
		// Ambil data bidang urusan
		bidangUrusan, err := service.BidangUrusanRepository.FindByKodeBidangUrusan(ctx, tx, tujuan.KodeBidangUrusan)
		if err != nil {
			return nil, err
		}

		tujuanResponse := tujuanopd.TujuanOpdResponse{
			Id:           tujuan.Id,
			Tujuan:       tujuan.Tujuan,
			TahunAwal:    tujuan.TahunAwal,
			TahunAkhir:   tujuan.TahunAkhir,
			JenisPeriode: tujuan.JenisPeriode,
			Indikator:    make([]tujuanopd.IndikatorResponse, 0),
		}

		// Proses indikator
		for _, indikator := range tujuan.Indikator {
			indikatorResponse := tujuanopd.IndikatorResponse{
				Id:               indikator.Id,
				IdTujuanOpd:      tujuan.Id,
				NamaIndikator:    indikator.Indikator,
				RumusPerhitungan: indikator.RumusPerhitungan.String,
				SumberData:       indikator.SumberData.String,
				Target:           make([]tujuanopd.TargetResponse, 0),
			}

			// Proses target untuk tahun yang diminta
			for _, target := range indikator.Target {
				if target.Tahun == tahun {
					targetResponse := tujuanopd.TargetResponse{
						Id:              target.Id,
						IndikatorId:     indikator.Id,
						Tahun:           target.Tahun,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					}
					indikatorResponse.Target = append(indikatorResponse.Target, targetResponse)
				}
			}

			// Hanya tambahkan indikator jika ada target
			if len(indikatorResponse.Target) > 0 {
				tujuanResponse.Indikator = append(tujuanResponse.Indikator, indikatorResponse)
			}
		}

		// Cek apakah sudah ada entry untuk kode_bidang_urusan ini
		mapKey := tujuan.KodeBidangUrusan
		if mapKey == "" {
			mapKey = "000" // Gunakan key default untuk bidang urusan kosong
		}

		if existing, exists := responseMap[mapKey]; exists {
			// Jika sudah ada dan tujuan memiliki indikator, tambahkan ke array yang ada
			if len(tujuanResponse.Indikator) > 0 {
				existing.TujuanOpd = append(existing.TujuanOpd, tujuanResponse)
			}
		} else {
			// Jika belum ada dan tujuan memiliki indikator, buat entry baru
			if len(tujuanResponse.Indikator) > 0 {
				// Ambil data urusan berdasarkan kode urusan dari bidang urusan
				var kodeUrusan string
				if len(bidangUrusan.KodeBidangUrusan) > 0 {
					kodeUrusan = bidangUrusan.KodeBidangUrusan[:1]
				}

				responseMap[mapKey] = &tujuanopd.TujuanOpdwithBidangUrusanResponse{
					Urusan:           bidangUrusan.NamaUrusan, // Menggunakan NamaUrusan dari BidangUrusan
					KodeUrusan:       kodeUrusan,
					KodeBidangUrusan: bidangUrusan.KodeBidangUrusan,
					NamaBidangUrusan: bidangUrusan.NamaBidangUrusan,
					KodeOpd:          tujuan.KodeOpd,
					NamaOpd:          opd.NamaOpd,
					TujuanOpd:        []tujuanopd.TujuanOpdResponse{tujuanResponse},
				}
			}
		}
	}

	// Convert map to slice
	var responses []tujuanopd.TujuanOpdwithBidangUrusanResponse
	for _, response := range responseMap {
		responses = append(responses, *response)
	}

	// Sort responses berdasarkan kode_bidang_urusan
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].KodeBidangUrusan < responses[j].KodeBidangUrusan
	})

	if len(responses) == 0 {
		responses = make([]tujuanopd.TujuanOpdwithBidangUrusanResponse, 0)
	}

	return responses, nil
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

// ─────────────────────────────────────────────────────────────────
// HELPER: bangun response TujuanOpdwithBidangUrusanResponse
//
//	dari domain, opd, dan bidangUrusanMap (sudah di-batch)
//
// ─────────────────────────────────────────────────────────────────
func (service *TujuanOpdServiceImpl) buildTujuanOpdResponse(
	tujuanOpds []domain.TujuanOpd,
	opd domainmaster.Opd,
	bidangUrusanMap map[string]domainmaster.BidangUrusan,
) []tujuanopd.TujuanOpdwithBidangUrusanResponse {
	responseMap := make(map[string]*tujuanopd.TujuanOpdwithBidangUrusanResponse)
	for _, tujuan := range tujuanOpds {
		tujuanResponse := tujuanopd.TujuanOpdResponse{
			Id:           tujuan.Id,
			Tujuan:       tujuan.Tujuan,
			TahunAwal:    tujuan.TahunAwal,
			TahunAkhir:   tujuan.TahunAkhir,
			JenisPeriode: tujuan.JenisPeriode,
			Indikator:    make([]tujuanopd.IndikatorResponse, 0),
		}
		for _, indikator := range tujuan.Indikator {
			indikatorResponse := tujuanopd.IndikatorResponse{
				Id:                  indikator.Id,
				KodeIndikator:       indikator.KodeIndikator,
				IdTujuanOpd:         tujuan.Id,
				NamaIndikator:       indikator.Indikator,
				RumusPerhitungan:    indikator.RumusPerhitungan.String,
				SumberData:          indikator.SumberData.String,
				DefinisiOperasional: indikator.DefinisiOperasional.String,
				Jenis:               indikator.Jenis,
				Target:              make([]tujuanopd.TargetResponse, 0),
			}
			for _, target := range indikator.Target {
				indikatorResponse.Target = append(indikatorResponse.Target, tujuanopd.TargetResponse{
					Id:              target.Id,
					IndikatorId:     indikator.KodeIndikator,
					Tahun:           target.Tahun,
					TargetIndikator: target.Target,
					SatuanIndikator: target.Satuan,
				})
			}
			tujuanResponse.Indikator = append(tujuanResponse.Indikator, indikatorResponse)
		}
		// Gunakan bidangUrusanMap hasil batch — tidak ada query loop
		mapKey := tujuan.KodeBidangUrusan
		if mapKey == "" {
			mapKey = "000"
		}
		if existing, exists := responseMap[mapKey]; exists {
			existing.TujuanOpd = append(existing.TujuanOpd, tujuanResponse)
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
				TujuanOpd:        []tujuanopd.TujuanOpdResponse{tujuanResponse},
			}
		}
	}
	var responses []tujuanopd.TujuanOpdwithBidangUrusanResponse
	for _, r := range responseMap {
		responses = append(responses, *r)
	}
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].KodeBidangUrusan < responses[j].KodeBidangUrusan
	})
	if len(responses) == 0 {
		return make([]tujuanopd.TujuanOpdwithBidangUrusanResponse, 0)
	}
	return responses
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
	return service.buildTujuanOpdResponse(tujuanOpds, opd, bidangUrusanMap), nil
}

// ─────────────────────────────────────────────────────────────────
// GET /tujuan_opd/renja_ranwal/:kode_opd/:tahun
// jenis indikator hardcode = "ranwal"
// target: 1 slot untuk tahun yang diminta
// ─────────────────────────────────────────────────────────────────
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
	tujuanOpds, err := service.TujuanOpdRepository.FindAllByTahun(
		ctx, tx, kodeOpd, tahun, jenisPeriode, "ranwal",
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

// ─────────────────────────────────────────────────────────────────
// GET /tujuan_opd/rankhir/:kode_opd/:tahun
// jenis indikator hardcode = "rankhir"
// target: 1 slot untuk tahun yang diminta
// ─────────────────────────────────────────────────────────────────
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
	tujuanOpds, err := service.TujuanOpdRepository.FindAllByTahun(
		ctx, tx, kodeOpd, tahun, jenisPeriode, "rankhir",
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

func (s *TujuanOpdServiceImpl) getIndikatorWithFallback(
	ctx context.Context,
	tx *sql.Tx,
	tujuanIds []int,
) ([]domain.Indikator, error) {

	indikatorBaru, err := s.TujuanOpdRepository.
		FindIndikatorTargetsRenstraByTujuanIds(ctx, tx, tujuanIds)
	if err != nil {
		return nil, err
	}

	indikatorLama, err := s.TujuanOpdRepository.
		FindIndikatorTargetsByTujuanIds(ctx, tx, tujuanIds)
	if err != nil {
		return nil, err
	}

	return mergeIndikator(indikatorBaru, indikatorLama), nil
}
func mergeIndikator(
	indikatorBaru []domain.Indikator,
	indikatorLama []domain.Indikator,
) []domain.Indikator {

	// fallback kalau data baru kosong
	if len(indikatorBaru) == 0 {
		log.Println("USING INDIKATOR LAMA")
		return indikatorLama
	}

	log.Println("USING INDIKATOR BARU")

	return indikatorBaru
}

func isEmptyIndikator(ind domain.Indikator) bool {
	return !ind.DefinisiOperasional.Valid &&
		len(ind.Target) == 0
}

func fallbackNullString(newVal, oldVal sql.NullString) sql.NullString {
	if !newVal.Valid || newVal.String == "" {
		if oldVal.Valid {
			return oldVal
		}
	}
	return newVal
}
