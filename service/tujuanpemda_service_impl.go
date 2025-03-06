package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/tujuanpemda"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type TujuanPemdaServiceImpl struct {
	TujuanPemdaRepository  repository.TujuanPemdaRepository
	PeriodeRepository      repository.PeriodeRepository
	PohonKinerjaRepository repository.PohonKinerjaRepository
	VisiPemdaRepository    repository.VisiPemdaRepository
	MisiPemdaRepository    repository.MisiPemdaRepository
	DB                     *sql.DB
}

func NewTujuanPemdaServiceImpl(tujuanPemdaRepository repository.TujuanPemdaRepository, periodeRepository repository.PeriodeRepository, pohonKinerjaRepository repository.PohonKinerjaRepository, visiPemdaRepository repository.VisiPemdaRepository, misiPemdaRepository repository.MisiPemdaRepository, DB *sql.DB) *TujuanPemdaServiceImpl {
	return &TujuanPemdaServiceImpl{
		TujuanPemdaRepository:  tujuanPemdaRepository,
		PeriodeRepository:      periodeRepository,
		PohonKinerjaRepository: pohonKinerjaRepository,
		VisiPemdaRepository:    visiPemdaRepository,
		MisiPemdaRepository:    misiPemdaRepository,
		DB:                     DB,
	}
}

func (service *TujuanPemdaServiceImpl) generateRandomId(ctx context.Context, tx *sql.Tx) int {
	rand.Seed(time.Now().UnixNano())
	for {
		// Generate random number between 10000-99999
		id := rand.Intn(90000) + 10000
		if !service.TujuanPemdaRepository.IsIdExists(ctx, tx, id) {
			return id
		}
	}
}

func generateIndikatorId() string {
	currentYear := time.Now().Format("2006")
	uuid := uuid.New().String()[:5] // Mengambil 5 karakter pertama dari UUID
	return fmt.Sprintf("IND-TJN-PMD-%s-%s", currentYear, uuid)
}

func generateTargetId() string {
	currentYear := time.Now().Format("2006")
	uuid := uuid.New().String()[:5]
	return fmt.Sprintf("TRG-TJN-PMD-%s-%s", currentYear, uuid)
}

func (service *TujuanPemdaServiceImpl) Create(ctx context.Context, request tujuanpemda.TujuanPemdaCreateRequest) (tujuanpemda.TujuanPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi periode
	periode, err := service.PeriodeRepository.FindById(ctx, tx, request.PeriodeId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("periode tidak ditemukan: %v", err)
	}

	//validasi pohon kinerja = tematik
	err = service.PohonKinerjaRepository.ValidatePokinLevel(ctx, tx, request.TematikId, 0, "tujuan pemda")
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	TujuanPemdaId, err := service.PohonKinerjaRepository.FindById(ctx, tx, request.TematikId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	// Validasi tahun target untuk setiap indikator
	tahunAwal, _ := strconv.Atoi(periode.TahunAwal)
	tahunAkhir, _ := strconv.Atoi(periode.TahunAkhir)

	for _, indikatorReq := range request.Indikator {
		tahunMap := make(map[string]bool)
		for _, targetReq := range indikatorReq.Target {
			targetTahun, _ := strconv.Atoi(targetReq.Tahun)

			// Validasi rentang tahun
			if targetTahun < tahunAwal || targetTahun > tahunAkhir {
				return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf(
					"tahun target %d harus berada dalam rentang periode %d-%d",
					targetTahun, tahunAwal, tahunAkhir,
				)
			}

			// Validasi duplikasi tahun
			if tahunMap[targetReq.Tahun] {
				return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf(
					"duplikasi tahun %s pada indikator %s",
					targetReq.Tahun, indikatorReq.Indikator,
				)
			}
			tahunMap[targetReq.Tahun] = true
		}
	}

	visiPemda, err := service.VisiPemdaRepository.FindById(ctx, tx, request.IdVisi)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	misiPemda, err := service.MisiPemdaRepository.FindById(ctx, tx, request.IdMisi)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	tujuanPemda := domain.TujuanPemda{
		Id:                service.generateRandomId(ctx, tx),
		TujuanPemda:       request.TujuanPemda,
		IdVisi:            request.IdVisi,
		IdMisi:            request.IdMisi,
		TematikId:         request.TematikId,
		PeriodeId:         request.PeriodeId,
		TahunAwalPeriode:  periode.TahunAwal,
		TahunAkhirPeriode: periode.TahunAkhir,
		JenisPeriode:      periode.JenisPeriode,
	}
	tujuanPemda, err = service.TujuanPemdaRepository.Create(ctx, tx, tujuanPemda)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	var indikatorResponses []tujuanpemda.IndikatorResponse
	for _, indikatorRequest := range request.Indikator {
		indikator := domain.Indikator{
			Id:            generateIndikatorId(),
			TujuanPemdaId: tujuanPemda.Id,
			Indikator:     indikatorRequest.Indikator,
			RumusPerhitungan: sql.NullString{
				String: indikatorRequest.RumusPerhitungan,
				Valid:  true,
			},
			SumberData: sql.NullString{
				String: indikatorRequest.SumberData,
				Valid:  true,
			},
		}

		indikator, err = service.TujuanPemdaRepository.CreateIndikator(ctx, tx, indikator)
		if err != nil {
			return tujuanpemda.TujuanPemdaResponse{}, err
		}

		var targetResponses []tujuanpemda.TargetResponse
		for _, targetRequest := range indikatorRequest.Target {
			target := domain.Target{
				Id:          generateTargetId(),
				IndikatorId: indikator.Id,
				Target:      targetRequest.Target,
				Satuan:      targetRequest.Satuan,
				Tahun:       targetRequest.Tahun,
			}

			err = service.TujuanPemdaRepository.CreateTarget(ctx, tx, target)
			if err != nil {
				return tujuanpemda.TujuanPemdaResponse{}, err
			}

			targetResponses = append(targetResponses, tujuanpemda.TargetResponse{
				Id:     target.Id,
				Target: target.Target,
				Satuan: target.Satuan,
				Tahun:  target.Tahun,
			})
		}

		indikatorResponses = append(indikatorResponses, tujuanpemda.IndikatorResponse{
			Id:               indikator.Id,
			Indikator:        indikator.Indikator,
			RumusPerhitungan: indikator.RumusPerhitungan.String,
			SumberData:       indikator.SumberData.String,
			Target:           targetResponses,
		})
	}

	return tujuanpemda.TujuanPemdaResponse{
		Id:          tujuanPemda.Id,
		Visi:        visiPemda.Visi,
		Misi:        misiPemda.Misi,
		TujuanPemda: tujuanPemda.TujuanPemda,
		TematikId:   tujuanPemda.TematikId,
		NamaTematik: TujuanPemdaId.NamaPohon,
		Periode: tujuanpemda.PeriodeResponse{
			TahunAwal:    periode.TahunAwal,
			TahunAkhir:   periode.TahunAkhir,
			JenisPeriode: periode.JenisPeriode,
		},
		Indikator: indikatorResponses,
	}, nil
}

func (service *TujuanPemdaServiceImpl) Update(ctx context.Context, request tujuanpemda.TujuanPemdaUpdateRequest) (tujuanpemda.TujuanPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Convert request.Id string ke int
	tujuanPemdaId := request.Id

	// Validasi tujuan pemda exists
	tujuanPemda, err := service.TujuanPemdaRepository.FindById(ctx, tx, tujuanPemdaId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	// Ambil data periode
	periode, err := service.PeriodeRepository.FindById(ctx, tx, request.PeriodeId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("periode tidak ditemukan: %v", err)
	}

	_, err = service.VisiPemdaRepository.FindById(ctx, tx, tujuanPemda.IdVisi)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	_, err = service.MisiPemdaRepository.FindById(ctx, tx, tujuanPemda.IdMisi)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	// Validasi tahun target untuk setiap indikator
	tahunAwal, _ := strconv.Atoi(periode.TahunAwal)
	tahunAkhir, _ := strconv.Atoi(periode.TahunAkhir)

	for _, indikatorReq := range request.Indikator {
		tahunMap := make(map[string]bool)
		for _, targetReq := range indikatorReq.Target {
			targetTahun, _ := strconv.Atoi(targetReq.Tahun)

			// Validasi rentang tahun
			if targetTahun < tahunAwal || targetTahun > tahunAkhir {
				return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf(
					"tahun target %d harus berada dalam rentang periode %d-%d",
					targetTahun, tahunAwal, tahunAkhir,
				)
			}

			// Validasi duplikasi tahun
			if tahunMap[targetReq.Tahun] {
				return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf(
					"duplikasi tahun %s pada indikator %s",
					targetReq.Tahun, indikatorReq.Indikator,
				)
			}
			tahunMap[targetReq.Tahun] = true
		}
	}

	// Update data tujuan pemda
	tujuanPemda.TujuanPemda = request.TujuanPemda
	tujuanPemda.IdVisi = request.IdVisi
	tujuanPemda.IdMisi = request.IdMisi
	tujuanPemda.PeriodeId = request.PeriodeId
	tujuanPemda.TahunAwalPeriode = periode.TahunAwal
	tujuanPemda.TahunAkhirPeriode = periode.TahunAkhir
	tujuanPemda.JenisPeriode = periode.JenisPeriode

	// Proses indikator
	var indikators []domain.Indikator
	for _, indikatorReq := range request.Indikator {
		var targets []domain.Target

		// Proses target untuk setiap indikator
		for _, targetReq := range indikatorReq.Target {
			var targetId string
			if targetReq.Id != "" {
				targetId = targetReq.Id
			} else {
				targetId = generateTargetId()
			}

			target := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorReq.Id,
				Target:      targetReq.Target,
				Satuan:      targetReq.Satuan,
				Tahun:       targetReq.Tahun,
			}
			targets = append(targets, target)
		}

		// Buat atau gunakan ID indikator yang ada
		var indikatorId string
		if indikatorReq.Id != "" {
			indikatorId = indikatorReq.Id
		} else {
			indikatorId = generateIndikatorId()
		}

		indikator := domain.Indikator{
			Id:            indikatorId,
			TujuanPemdaId: tujuanPemdaId,
			Indikator:     indikatorReq.Indikator,
			RumusPerhitungan: sql.NullString{
				String: indikatorReq.RumusPerhitungan,
				Valid:  true,
			},
			SumberData: sql.NullString{
				String: indikatorReq.SumberData,
				Valid:  true,
			},
			Target: targets,
		}
		indikators = append(indikators, indikator)
	}

	tujuanPemda.Indikator = indikators

	// Simpan semua perubahan
	result, err := service.TujuanPemdaRepository.Update(ctx, tx, tujuanPemda)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	return service.toTujuanPemdaResponse(result), nil
}
func (service *TujuanPemdaServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Hapus indikator (cascade akan menghapus target)
	err = service.TujuanPemdaRepository.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	// Hapus tujuan pemda
	return service.TujuanPemdaRepository.Delete(ctx, tx, id)
}

func (service *TujuanPemdaServiceImpl) FindById(ctx context.Context, tujuanPemdaId int) (tujuanpemda.TujuanPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	tujuanPemda, err := service.TujuanPemdaRepository.FindById(ctx, tx, tujuanPemdaId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, tujuanPemda.TematikId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
	}

	visiPemda, err := service.VisiPemdaRepository.FindById(ctx, tx, tujuanPemda.IdVisi)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	misiPemda, err := service.MisiPemdaRepository.FindById(ctx, tx, tujuanPemda.IdMisi)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	var indikatorResponses []tujuanpemda.IndikatorResponse
	for _, indikator := range tujuanPemda.Indikator {
		indikatorResponse := tujuanpemda.IndikatorResponse{
			Id:               indikator.Id,
			Indikator:        indikator.Indikator,
			RumusPerhitungan: indikator.RumusPerhitungan.String,
			SumberData:       indikator.SumberData.String,
			Target:           []tujuanpemda.TargetResponse{},
		}

		// Tambahkan target hanya jika tahun awal dan akhir valid
		if tujuanPemda.Periode.TahunAwal != "" && tujuanPemda.Periode.TahunAwal != "Pilih periode" {
			// Urutkan target berdasarkan tahun
			sort.Slice(indikator.Target, func(i, j int) bool {
				return indikator.Target[i].Tahun < indikator.Target[j].Tahun
			})

			indikatorResponse.Target = convertToTargetResponses(indikator.Target)
		}

		indikatorResponses = append(indikatorResponses, indikatorResponse)
	}

	return tujuanpemda.TujuanPemdaResponse{
		Id:          tujuanPemda.Id,
		IdVisi:      tujuanPemda.IdVisi,
		IdMisi:      tujuanPemda.IdMisi,
		Visi:        visiPemda.Visi,
		Misi:        misiPemda.Misi,
		TujuanPemda: tujuanPemda.TujuanPemda,
		TematikId:   tujuanPemda.TematikId,
		NamaTematik: pokinData.NamaPohon,
		JenisPohon:  tujuanPemda.JenisPohon,
		PeriodeId:   tujuanPemda.PeriodeId,
		Periode: tujuanpemda.PeriodeResponse{
			TahunAwal:    tujuanPemda.Periode.TahunAwal,
			TahunAkhir:   tujuanPemda.Periode.TahunAkhir,
			JenisPeriode: tujuanPemda.Periode.JenisPeriode,
		},
		Indikator: indikatorResponses,
	}, nil
}

func (service *TujuanPemdaServiceImpl) FindAll(ctx context.Context, tahun string, jenisPeriode string) ([]tujuanpemda.TujuanPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []tujuanpemda.TujuanPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	tujuanPemdaList, err := service.TujuanPemdaRepository.FindAll(ctx, tx, tahun, jenisPeriode)
	if err != nil {
		return []tujuanpemda.TujuanPemdaResponse{}, err
	}

	tujuanPemdaResponses := make([]tujuanpemda.TujuanPemdaResponse, 0, len(tujuanPemdaList))
	for _, tujuanPemda := range tujuanPemdaList {
		pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, tujuanPemda.TematikId)
		if err != nil {
			return []tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
		}

		tujuanPemdaResponses = append(tujuanPemdaResponses, tujuanpemda.TujuanPemdaResponse{
			Id:          tujuanPemda.Id,
			TujuanPemda: tujuanPemda.TujuanPemda,
			TematikId:   tujuanPemda.TematikId,
			NamaTematik: pokinData.NamaPohon,
			Periode: tujuanpemda.PeriodeResponse{
				TahunAwal:    tujuanPemda.Periode.TahunAwal,
				TahunAkhir:   tujuanPemda.Periode.TahunAkhir,
				JenisPeriode: tujuanPemda.Periode.JenisPeriode,
			},
		})
	}

	return tujuanPemdaResponses, nil
}

// Fungsi helper untuk konversi target
func convertToTargetResponses(targets []domain.Target) []tujuanpemda.TargetResponse {
	var targetResponses []tujuanpemda.TargetResponse
	for _, target := range targets {
		targetResponses = append(targetResponses, tujuanpemda.TargetResponse{
			Id:     target.Id,
			Target: target.Target,
			Satuan: target.Satuan,
			Tahun:  target.Tahun,
		})
	}
	return targetResponses
}

// Fungsi helper untuk konversi ke response
func (service *TujuanPemdaServiceImpl) toTujuanPemdaResponse(tujuanPemda domain.TujuanPemda) tujuanpemda.TujuanPemdaResponse {
	var indikatorResponses []tujuanpemda.IndikatorResponse
	for _, indikator := range tujuanPemda.Indikator {
		var targetResponses []tujuanpemda.TargetResponse
		for _, target := range indikator.Target {
			targetResponses = append(targetResponses, tujuanpemda.TargetResponse{
				Id:     target.Id,
				Target: target.Target,
				Satuan: target.Satuan,
				Tahun:  target.Tahun,
			})
		}

		indikatorResponses = append(indikatorResponses, tujuanpemda.IndikatorResponse{
			Id:               indikator.Id,
			Indikator:        indikator.Indikator,
			RumusPerhitungan: indikator.RumusPerhitungan.String,
			SumberData:       indikator.SumberData.String,
			Target:           targetResponses,
		})
	}

	return tujuanpemda.TujuanPemdaResponse{
		Id:          tujuanPemda.Id,
		TujuanPemda: tujuanPemda.TujuanPemda,
		IdVisi:      tujuanPemda.IdVisi,
		IdMisi:      tujuanPemda.IdMisi,
		Visi:        tujuanPemda.Visi,
		Misi:        tujuanPemda.Misi,
		TematikId:   tujuanPemda.TematikId,
		NamaTematik: tujuanPemda.NamaTematik,
		Periode: tujuanpemda.PeriodeResponse{
			TahunAwal:    tujuanPemda.Periode.TahunAwal,
			TahunAkhir:   tujuanPemda.Periode.TahunAkhir,
			JenisPeriode: tujuanPemda.Periode.JenisPeriode,
		},
		Indikator: indikatorResponses,
	}
}

func (service *TujuanPemdaServiceImpl) UpdatePeriode(ctx context.Context, request tujuanpemda.TujuanPemdaUpdateRequest) (tujuanpemda.TujuanPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi tujuan pemda exists
	exists := service.TujuanPemdaRepository.IsIdExists(ctx, tx, request.Id)
	if !exists {
		return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("periode pemda dengan id %d tidak ditemukan", request.Id)
	}

	// Validasi periode exists jika periode_id tidak 0
	if request.PeriodeId != 0 {
		exists = service.PeriodeRepository.IsIdExists(ctx, tx, request.PeriodeId)
		if !exists {
			return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("periode dengan id %d tidak ditemukan", request.PeriodeId)
		}
	}

	// Persiapkan domain object untuk update (hanya id dan periode_id)
	tujuanPemda := domain.TujuanPemda{
		Id:        request.Id,
		PeriodeId: request.PeriodeId,
	}

	// Update periode
	result, err := service.TujuanPemdaRepository.UpdatePeriode(ctx, tx, tujuanPemda)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	// Ambil data pohon kinerja untuk response
	pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, result.TematikId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
	}

	// Buat response
	return tujuanpemda.TujuanPemdaResponse{
		Id:          result.Id,
		TujuanPemda: result.TujuanPemda,
		TematikId:   result.TematikId,
		NamaTematik: pokinData.NamaPohon,
		Periode: tujuanpemda.PeriodeResponse{
			TahunAwal:  result.Periode.TahunAwal,
			TahunAkhir: result.Periode.TahunAkhir,
		},
	}, nil
}

func (service *TujuanPemdaServiceImpl) FindAllWithPokin(ctx context.Context, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]tujuanpemda.TujuanPemdaWithPokinResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	tujuanPemdaList, err := service.TujuanPemdaRepository.FindAllWithPokin(ctx, tx, tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		return nil, err
	}

	tahunAwalInt, err := strconv.Atoi(tahunAwal)
	if err != nil {
		return nil, fmt.Errorf("format tahun awal tidak valid: %v", err)
	}
	tahunAkhirInt, err := strconv.Atoi(tahunAkhir)
	if err != nil {
		return nil, fmt.Errorf("format tahun akhir tidak valid: %v", err)
	}

	// Buat map untuk mengelompokkan berdasarkan PokinId
	pokinMap := make(map[int]tujuanpemda.TujuanPemdaWithPokinResponse)

	for _, item := range tujuanPemdaList {
		tahunPokinInt, err := strconv.Atoi(item.TahunPokin)
		if err != nil {
			continue // Skip jika tahun pokin tidak valid
		}

		// Skip jika tahun pokin di luar range periode
		if tahunPokinInt < tahunAwalInt || tahunPokinInt > tahunAkhirInt {
			continue
		}
		pokinResp, exists := pokinMap[item.PokinId]
		if !exists {
			pokinResp = tujuanpemda.TujuanPemdaWithPokinResponse{
				PokinId:     item.PokinId,
				NamaPohon:   item.NamaPohon,
				JenisPohon:  item.JenisPohon,
				LevelPohon:  item.LevelPohon,
				Keterangan:  item.Keterangan,
				TahunPokin:  item.TahunPokin,
				IsActive:    item.IsActive,
				TujuanPemda: make([]tujuanpemda.TujuanPemdaResponse, 0),
			}
		}

		// Proses setiap tujuan pemda
		for _, tujuanPemda := range item.TujuanPemda {
			var indikatorResponses []tujuanpemda.IndikatorResponse

			// Ambil data visi dengan penanganan default
			visiPemda, err := service.VisiPemdaRepository.FindById(ctx, tx, tujuanPemda.IdVisi)
			if err != nil {
				// Jika visi tidak ditemukan, gunakan nilai default
				visiPemda = domain.VisiPemda{
					Id:   0,
					Visi: "Belum ada visi",
				}
			}

			// Ambil data misi dengan penanganan default
			misiPemda, err := service.MisiPemdaRepository.FindById(ctx, tx, tujuanPemda.IdMisi)
			if err != nil {
				// Jika misi tidak ditemukan, gunakan nilai default
				misiPemda = domain.MisiPemda{
					Id:   0,
					Misi: "Belum ada misi",
				}
			}

			// Proses indikator
			for _, indikator := range tujuanPemda.Indikator {
				targetResponses := make([]tujuanpemda.TargetResponse, 0)

				// Buat map untuk target yang ada
				targetMap := make(map[string]domain.Target)
				for _, target := range indikator.Target {
					targetMap[target.Tahun] = target
				}

				// Generate target untuk semua tahun dalam periode
				tahunAwalInt, _ := strconv.Atoi(tahunAwal)
				tahunAkhirInt, _ := strconv.Atoi(tahunAkhir)

				for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
					tahunStr := strconv.Itoa(tahun)
					if target, exists := targetMap[tahunStr]; exists {
						targetResponses = append(targetResponses, tujuanpemda.TargetResponse{
							Id:     target.Id,
							Target: target.Target,
							Satuan: target.Satuan,
							Tahun:  target.Tahun,
						})
					} else {
						targetResponses = append(targetResponses, tujuanpemda.TargetResponse{
							Id:     "-",
							Target: "-",
							Satuan: "-",
							Tahun:  tahunStr,
						})
					}
				}

				indikatorResponses = append(indikatorResponses, tujuanpemda.IndikatorResponse{
					Id:               indikator.Id,
					Indikator:        indikator.Indikator,
					RumusPerhitungan: indikator.RumusPerhitungan.String,
					SumberData:       indikator.SumberData.String,
					Target:           targetResponses,
				})
			}

			tujuanPemdaResponse := tujuanpemda.TujuanPemdaResponse{
				Id:          tujuanPemda.Id,
				IdVisi:      visiPemda.Id,
				Visi:        visiPemda.Visi,
				IdMisi:      misiPemda.Id,
				Misi:        misiPemda.Misi,
				TujuanPemda: tujuanPemda.TujuanPemda,
				Periode: tujuanpemda.PeriodeResponse{
					TahunAwal:    tahunAwal,
					TahunAkhir:   tahunAkhir,
					JenisPeriode: jenisPeriode,
				},
				Indikator: indikatorResponses,
			}

			pokinResp.TujuanPemda = append(pokinResp.TujuanPemda, tujuanPemdaResponse)
		}

		pokinMap[item.PokinId] = pokinResp
	}

	// Konversi map ke slice dan sorting
	result := make([]tujuanpemda.TujuanPemdaWithPokinResponse, 0, len(pokinMap))
	for _, pokinResp := range pokinMap {
		// Sort indikator
		for i := range pokinResp.TujuanPemda {
			if len(pokinResp.TujuanPemda[i].Indikator) > 0 {
				sort.Slice(pokinResp.TujuanPemda[i].Indikator, func(x, y int) bool {
					return pokinResp.TujuanPemda[i].Indikator[x].Id < pokinResp.TujuanPemda[i].Indikator[y].Id
				})
			}
		}

		// Sort tujuan pemda
		if len(pokinResp.TujuanPemda) > 0 {
			sort.Slice(pokinResp.TujuanPemda, func(i, j int) bool {
				return pokinResp.TujuanPemda[i].Id < pokinResp.TujuanPemda[j].Id
			})
		}

		result = append(result, pokinResp)
	}

	// Sort hasil akhir berdasarkan PokinId
	sort.Slice(result, func(i, j int) bool {
		return result[i].PokinId < result[j].PokinId
	})

	return result, nil
}

func (service *TujuanPemdaServiceImpl) FindPokinWithPeriode(ctx context.Context, pokinId int, jenisPeriode string) (tujuanpemda.PokinWithPeriodeResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanpemda.PokinWithPeriodeResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi pokin ID
	err = service.PohonKinerjaRepository.ValidatePokinId(ctx, tx, pokinId)
	if err != nil {
		return tujuanpemda.PokinWithPeriodeResponse{}, err
	}

	// Validasi jenis periode
	if jenisPeriode == "" {
		return tujuanpemda.PokinWithPeriodeResponse{}, fmt.Errorf("jenis periode tidak boleh kosong")
	}

	// Ambil data pokin dengan periode
	pokin, periode, err := service.PohonKinerjaRepository.FindPokinWithPeriode(ctx, tx, pokinId, jenisPeriode)
	if err != nil {
		return tujuanpemda.PokinWithPeriodeResponse{}, err
	}

	// Transform ke response
	response := tujuanpemda.PokinWithPeriodeResponse{
		Id:         pokin.Id,
		NamaPohon:  pokin.NamaPohon,
		JenisPohon: pokin.JenisPohon,
		LevelPohon: pokin.LevelPohon,
		Tahun:      pokin.Tahun,
		Status:     pokin.Status,
		Periode: tujuanpemda.PokinPeriodeResponse{
			Id:         periode.Id,
			TahunAwal:  periode.TahunAwal,
			TahunAkhir: periode.TahunAkhir,
		},
		Indikator: []tujuanpemda.PokinIndikatorResponse{},
	}

	// Transform indikator dan target
	for _, indikator := range pokin.Indikator {
		indikatorResponse := tujuanpemda.PokinIndikatorResponse{
			Id:               indikator.Id,
			Indikator:        indikator.Indikator,
			RumusPerhitungan: indikator.RumusPerhitungan.String,
			SumberData:       indikator.SumberData.String,
			Target:           []tujuanpemda.PokinTargetResponse{},
		}

		for _, target := range indikator.Target {
			targetResponse := tujuanpemda.PokinTargetResponse{
				Id:     target.Id,
				Target: target.Target,
				Satuan: target.Satuan,
				Tahun:  target.Tahun,
			}
			indikatorResponse.Target = append(indikatorResponse.Target, targetResponse)
		}

		response.Indikator = append(response.Indikator, indikatorResponse)
	}

	return response, nil
}
