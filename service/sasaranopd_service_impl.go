package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/model/web/sasaranopd"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SasaranOpdServiceImpl struct {
	sasaranOpdRepository      repository.SasaranOpdRepository
	opdRepository             repository.OpdRepository
	rencanaKinerjaRepository  repository.RencanaKinerjaRepository
	manualIndikatorRepository repository.ManualIKRepository
	pegawaiRepository         repository.PegawaiRepository
	pohonkinerjaRepository    repository.PohonKinerjaRepository
	DB                        *sql.DB
	validate                  *validator.Validate
	tujuanOpdRepository       repository.TujuanOpdRepository
}

func NewSasaranOpdServiceImpl(
	sasaranOpdRepository repository.SasaranOpdRepository,
	opdRepository repository.OpdRepository,
	rencanaKinerjaRepository repository.RencanaKinerjaRepository,
	manualIndikatorRepository repository.ManualIKRepository,
	pegawaiRepository repository.PegawaiRepository,
	pohonkinerjaRepository repository.PohonKinerjaRepository,
	tujuanOpdRepository repository.TujuanOpdRepository,
	db *sql.DB,
	validate *validator.Validate,
) *SasaranOpdServiceImpl {
	return &SasaranOpdServiceImpl{
		sasaranOpdRepository:      sasaranOpdRepository,
		opdRepository:             opdRepository,
		rencanaKinerjaRepository:  rencanaKinerjaRepository,
		manualIndikatorRepository: manualIndikatorRepository,
		pegawaiRepository:         pegawaiRepository,
		pohonkinerjaRepository:    pohonkinerjaRepository,
		tujuanOpdRepository:       tujuanOpdRepository,
		DB:                        db,
		validate:                  validate,
	}
}

func (service *SasaranOpdServiceImpl) FindAll(ctx context.Context, KodeOpd string, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]sasaranopd.SasaranOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	sasaranOpds, err := service.sasaranOpdRepository.FindAllOnly(ctx, tx, KodeOpd, tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		return nil, err
	}

	sasaranOpdIds := []int{}
	for _, so := range sasaranOpds {
		for _, sas := range so.SasaranOpd {
			sasaranOpdIds = append(sasaranOpdIds, sas.Id)
		}
	}
	indikatorSasaran, err := service.getIndikatorWithFallback(ctx, tx, sasaranOpdIds)
	if err != nil {
		log.Printf("ERROR service.getIndikatorWithFallback: %v", err)
		return nil, err
	}
	indSasaranById := make(map[int][]domain.Indikator)
	for _, ind := range indikatorSasaran {
		indSasaranById[ind.SasaranOpdId] = append(indSasaranById[ind.SasaranOpdId], ind)
	}

	for _, so := range sasaranOpds {

		for i := range so.SasaranOpd {
			id := so.SasaranOpd[i].Id
			if inds, ok := indSasaranById[id]; ok {
				so.SasaranOpd[i].Indikator = inds
			}
		}
	}

	// Sort sasaranOpds berdasarkan nama_pohon, jika sama berdasarkan id ASC
	sort.Slice(sasaranOpds, func(i, j int) bool {
		if sasaranOpds[i].NamaPohon == sasaranOpds[j].NamaPohon {
			return sasaranOpds[i].Id < sasaranOpds[j].Id
		}
		return sasaranOpds[i].NamaPohon < sasaranOpds[j].NamaPohon
	})

	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, KodeOpd)
	if err != nil {
		return nil, err
	}

	var responses []sasaranopd.SasaranOpdResponse
	for _, sasaranOpd := range sasaranOpds {
		response := sasaranopd.SasaranOpdResponse{
			IdPohon:    sasaranOpd.IdPohon,
			KodeOpd:    sasaranOpd.KodeOpd,
			NamaOpd:    opd.NamaOpd,
			NamaPohon:  sasaranOpd.NamaPohon,
			JenisPohon: sasaranOpd.JenisPohon,
			LevelPohon: sasaranOpd.LevelPohon,
			TahunPohon: sasaranOpd.TahunPohon,
			Pelaksana:  make([]sasaranopd.PelaksanaOpdResponse, 0),
			SasaranOpd: make([]sasaranopd.SasaranOpdDetailResponse, 0),
		}

		// Convert Pelaksana
		for _, pelaksana := range sasaranOpd.Pelaksana {
			response.Pelaksana = append(response.Pelaksana, sasaranopd.PelaksanaOpdResponse{
				Id:          pelaksana.Id,
				PegawaiId:   pelaksana.PegawaiId,
				Nip:         pelaksana.Nip,
				NamaPegawai: pelaksana.NamaPegawai,
			})
		}

		// Temporary slice untuk sorting sasaran
		tempSasaranResponses := make([]sasaranopd.SasaranOpdDetailResponse, 0)

		// Convert SasaranOpd
		for _, sasaran := range sasaranOpd.SasaranOpd {
			TujuanOpd, _ := service.tujuanOpdRepository.FindById(ctx, tx, sasaran.IdTujuanOpd)

			sasaranResponse := sasaranopd.SasaranOpdDetailResponse{
				Id:             strconv.Itoa(sasaran.Id),
				NamaSasaranOpd: sasaran.NamaSasaranOpd,
				IdTujuanOpd:    TujuanOpd.Id,
				NamaTujuanOpd:  TujuanOpd.Tujuan,
				TahunAwal:      sasaran.TahunAwal,
				TahunAkhir:     sasaran.TahunAkhir,
				JenisPeriode:   sasaran.JenisPeriode,
				Indikator:      make([]sasaranopd.IndikatorResponse, 0),
			}

			// Convert Indikator
			for _, indikator := range sasaran.Indikator {
				indResponse := sasaranopd.IndikatorResponse{
					Id:                  indikator.Id,
					Indikator:           indikator.Indikator,
					RumusPerhitungan:    indikator.RumusPerhitungan.String,
					DefinisiOperasional: indikator.DefinisiOperasional.String,
					SumberData:          indikator.SumberData.String,
					Target:              make([]sasaranopd.TargetResponse, 0),
				}

				// Convert Target
				for _, target := range indikator.Target {
					indResponse.Target = append(indResponse.Target, sasaranopd.TargetResponse{
						Id:     target.Id,
						Tahun:  target.Tahun,
						Target: target.Target,
						Satuan: target.Satuan,
					})
				}

				sasaranResponse.Indikator = append(sasaranResponse.Indikator, indResponse)
			}

			tempSasaranResponses = append(tempSasaranResponses, sasaranResponse)
		}

		// Sort sasaran berdasarkan nama_sasaran_opd
		sort.Slice(tempSasaranResponses, func(i, j int) bool {
			return tempSasaranResponses[i].NamaSasaranOpd < tempSasaranResponses[j].NamaSasaranOpd
		})

		// Assign sorted sasaran ke response
		response.SasaranOpd = tempSasaranResponses

		responses = append(responses, response)
	}

	return responses, nil
}

func (service *SasaranOpdServiceImpl) FindById(ctx context.Context, id int) (*sasaranopd.SasaranOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	sasaranOpd, err := service.sasaranOpdRepository.FindByIdOnly(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	sasaranOpdIds := []int{}
	for _, sas := range sasaranOpd.SasaranOpd {
		sasaranOpdIds = append(sasaranOpdIds, sas.Id)
	}
	indikatorSasaran, err := service.getIndikatorWithFallback(ctx, tx, sasaranOpdIds)
	if err != nil {
		log.Printf("ERROR service.getIndikatorWithFallback: %v", err)
		return nil, err
	}
	indSasaranById := make(map[int][]domain.Indikator)
	for _, ind := range indikatorSasaran {
		indSasaranById[ind.SasaranOpdId] = append(indSasaranById[ind.SasaranOpdId], ind)
	}

	for i := range sasaranOpd.SasaranOpd {
		id := sasaranOpd.SasaranOpd[i].Id
		if inds, ok := indSasaranById[id]; ok {
			sasaranOpd.SasaranOpd[i].Indikator = inds
		}
	}

	response := &sasaranopd.SasaranOpdResponse{
		IdPohon:    sasaranOpd.IdPohon,
		NamaPohon:  sasaranOpd.NamaPohon,
		JenisPohon: sasaranOpd.JenisPohon,
		LevelPohon: sasaranOpd.LevelPohon,
		TahunPohon: sasaranOpd.TahunPohon,
		Pelaksana:  make([]sasaranopd.PelaksanaOpdResponse, 0),
		SasaranOpd: make([]sasaranopd.SasaranOpdDetailResponse, 0),
	}

	// Convert Pelaksana
	for _, pelaksana := range sasaranOpd.Pelaksana {
		response.Pelaksana = append(response.Pelaksana, sasaranopd.PelaksanaOpdResponse{
			Id:          pelaksana.Id,
			PegawaiId:   pelaksana.PegawaiId,
			Nip:         pelaksana.Nip,
			NamaPegawai: pelaksana.NamaPegawai,
		})
	}

	// Convert SasaranOpd
	for _, sasaran := range sasaranOpd.SasaranOpd {
		TujuanOpd, _ := service.tujuanOpdRepository.FindById(ctx, tx, sasaran.IdTujuanOpd)
		sasaranResponse := sasaranopd.SasaranOpdDetailResponse{
			Id:             strconv.Itoa(sasaran.Id),
			NamaSasaranOpd: sasaran.NamaSasaranOpd,
			IdTujuanOpd:    sasaran.IdTujuanOpd,
			NamaTujuanOpd:  TujuanOpd.Tujuan,
			TahunAwal:      sasaran.TahunAwal,
			TahunAkhir:     sasaran.TahunAkhir,
			JenisPeriode:   sasaran.JenisPeriode,
			Indikator:      make([]sasaranopd.IndikatorResponse, 0),
		}

		// Convert Indikator
		for _, indikator := range sasaran.Indikator {
			indResponse := sasaranopd.IndikatorResponse{
				Id:                  indikator.Id,
				Indikator:           indikator.Indikator,
				RumusPerhitungan:    indikator.RumusPerhitungan.String,
				DefinisiOperasional: indikator.DefinisiOperasional.String,
				SumberData:          indikator.SumberData.String,
				Target:              make([]sasaranopd.TargetResponse, 0),
			}

			// Convert Target
			for _, target := range indikator.Target {
				indResponse.Target = append(indResponse.Target, sasaranopd.TargetResponse{
					Id:     target.Id,
					Tahun:  target.Tahun,
					Target: target.Target,
					Satuan: target.Satuan,
				})
			}

			sasaranResponse.Indikator = append(sasaranResponse.Indikator, indResponse)
		}

		response.SasaranOpd = append(response.SasaranOpd, sasaranResponse)
	}

	return response, nil
}

func (service *SasaranOpdServiceImpl) Create(ctx context.Context, request sasaranopd.SasaranOpdCreateRequest) (*sasaranopd.SasaranOpdCreateResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	if err := service.validate.Struct(request); err != nil {
		return nil, err
	}

	_, err = service.pohonkinerjaRepository.FindById(ctx, tx, request.IdPohon)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("data pohon kinerja dengan ID %d tidak ditemukan", request.IdPohon)
		}
		return nil, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
	}
	// Generate ID yang lebih unik untuk sasaran OPD
	idSasaran := rand.Intn(100000)

	sasaranOpd := domain.SasaranOpdDetail{
		Id:             idSasaran,
		IdPohon:        request.IdPohon,
		NamaSasaranOpd: request.NamaSasaran,
		TahunAwal:      request.TahunAwal,
		TahunAkhir:     request.TahunAkhir,
		JenisPeriode:   request.JenisPeriode,
		IdTujuanOpd:    request.IdTujuanOpd,
		Indikator:      make([]domain.Indikator, 0),
	}

	// Proses indikator
	for _, indReq := range request.Indikator {
		kodeIndikator := fmt.Sprintf("IND-SAR-%s", uuid.New().String()[:5])

		indikator := domain.Indikator{
			KodeIndikator:       kodeIndikator,
			Indikator:           indReq.Indikator,
			Jenis:               "renstra",
			DefinisiOperasional: sql.NullString{String: indReq.DefinisiOperasional, Valid: true},
			RumusPerhitungan:    sql.NullString{String: indReq.RumusPerhitungan, Valid: true},
			SumberData:          sql.NullString{String: indReq.SumberData, Valid: true},
			Target:              make([]domain.Target, 0),
		}

		// Proses target
		for _, targetReq := range indReq.Target {
			if targetReq.Target != "" {
				targetId := fmt.Sprintf("TRG-SAR-%d-%s", uuid.New().ID()%100000, targetReq.Tahun)

				target := domain.Target{
					Id:          targetId,
					IndikatorId: kodeIndikator,
					Tahun:       targetReq.Tahun,
					Target:      targetReq.Target,
					Satuan:      targetReq.Satuan,
				}
				indikator.Target = append(indikator.Target, target)
			}
		}

		sasaranOpd.Indikator = append(sasaranOpd.Indikator, indikator)
	}

	err = service.sasaranOpdRepository.Create(ctx, tx, sasaranOpd)
	if err != nil {
		return nil, err
	}

	TujuanOpd, err := service.tujuanOpdRepository.FindById(ctx, tx, sasaranOpd.IdTujuanOpd)
	if err != nil {
		return nil, errors.New("tujuan opd tidak ditemukan")
	}

	// Buat response dengan indikator dan target
	response := &sasaranopd.SasaranOpdCreateResponse{
		IdPohon:        sasaranOpd.IdPohon,
		NamaSasaranOpd: sasaranOpd.NamaSasaranOpd,
		NamaTujuanOpd:  TujuanOpd.Tujuan,
		TahunAwal:      sasaranOpd.TahunAwal,
		TahunAkhir:     sasaranOpd.TahunAkhir,
		JenisPeriode:   sasaranOpd.JenisPeriode,
		Indikator:      make([]sasaranopd.IndikatorDetail, 0),
	}

	// Convert indikator untuk response
	for _, indikator := range sasaranOpd.Indikator {
		indResponse := sasaranopd.IndikatorDetail{
			Id:                  indikator.Id,
			Indikator:           indikator.Indikator,
			Jenis:               "renstra",
			DefinisiOperasional: indikator.DefinisiOperasional.String,
			RumusPerhitungan:    indikator.RumusPerhitungan.String,
			SumberData:          indikator.SumberData.String,
			Target:              make([]sasaranopd.TargetDetail, 0),
		}

		// Convert target untuk response
		for _, target := range indikator.Target {
			indResponse.Target = append(indResponse.Target, sasaranopd.TargetDetail{
				Id:     target.Id,
				Tahun:  target.Tahun,
				Target: target.Target,
				Satuan: target.Satuan,
			})
		}

		response.Indikator = append(response.Indikator, indResponse)
	}

	return response, nil
}

func (service *SasaranOpdServiceImpl) Update(ctx context.Context, request sasaranopd.SasaranOpdUpdateRequest) (*sasaranopd.SasaranOpdCreateResponse, error) {
	if err := service.validate.Struct(request); err != nil {
		return nil, err
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	// Validasi tahun
	tahunAwalInt, err := strconv.Atoi(request.TahunAwal)
	if err != nil {
		return nil, errors.New("format tahun awal tidak valid")
	}
	tahunAkhirInt, err := strconv.Atoi(request.TahunAkhir)
	if err != nil {
		return nil, errors.New("format tahun akhir tidak valid")
	}
	if tahunAkhirInt < tahunAwalInt {
		return nil, errors.New("tahun akhir tidak boleh lebih kecil dari tahun awal")
	}
	// Validasi sasaran OPD ada
	existingSasaran, err := service.sasaranOpdRepository.FindByIdSasaran(ctx, tx, request.IdSasaranOpd)
	if err != nil {
		return nil, errors.New("sasaran opd tidak ditemukan")
	}
	// Validasi tujuan OPD ada
	tujuanOpd, err := service.tujuanOpdRepository.FindById(ctx, tx, request.IdTujuanOpd)
	if err != nil {
		return nil, errors.New("tujuan opd tidak ditemukan")
	}
	// Bangun daftar indikator domain + response
	var indikatorList []domain.Indikator
	var indikatorResponses []sasaranopd.IndikatorDetail
	for _, indikatorReq := range request.Indikator {
		kodeIndikator := helper.GenerateID("IND-SAS")
		if len(indikatorReq.Target) == 0 {
			return nil, fmt.Errorf("indikator '%s' harus memiliki minimal 1 target", indikatorReq.Indikator)
		}
		var targetList []domain.Target
		var targetResponses []sasaranopd.TargetDetail
		for _, targetReq := range indikatorReq.Target {
			// Tentukan ID target
			targetId := helper.GenerateID("TRG-SAS")
			targetList = append(targetList, domain.Target{
				Id:          targetId,
				IndikatorId: kodeIndikator, // referensi ke kode_indikator
				Tahun:       targetReq.Tahun,
				Target:      targetReq.Target,
				Satuan:      targetReq.Satuan,
			})
			targetResponses = append(targetResponses, sasaranopd.TargetDetail{
				Id:     targetId,
				Tahun:  targetReq.Tahun,
				Target: targetReq.Target,
				Satuan: targetReq.Satuan,
			})
		}
		indikatorList = append(indikatorList, domain.Indikator{
			Id:            kodeIndikator,
			KodeIndikator: kodeIndikator,
			SasaranOpdId:  request.IdSasaranOpd,
			Indikator:     indikatorReq.Indikator,
			Jenis:         "renstra", // ← hardcode
			DefinisiOperasional: sql.NullString{
				String: indikatorReq.DefinisiOperasional,
				Valid:  true,
			},
			RumusPerhitungan: sql.NullString{
				String: indikatorReq.RumusPerhitungan,
				Valid:  true,
			},
			SumberData: sql.NullString{
				String: indikatorReq.SumberData,
				Valid:  true,
			},
			Target: targetList,
		})
		indikatorResponses = append(indikatorResponses, sasaranopd.IndikatorDetail{
			Id:                  kodeIndikator,
			KodeIndikator:       kodeIndikator,
			Jenis:               "renstra",
			DefinisiOperasional: indikatorReq.DefinisiOperasional,
			Indikator:           indikatorReq.Indikator,
			RumusPerhitungan:    indikatorReq.RumusPerhitungan,
			SumberData:          indikatorReq.SumberData,
			Target:              targetResponses,
		})
	}
	sasaranOpdUpdate := domain.SasaranOpdDetail{
		Id:             request.IdSasaranOpd,
		IdPohon:        existingSasaran.IdPohon,
		NamaSasaranOpd: request.NamaSasaran,
		IdTujuanOpd:    request.IdTujuanOpd,
		TahunAwal:      request.TahunAwal,
		TahunAkhir:     request.TahunAkhir,
		JenisPeriode:   request.JenisPeriode,
		Indikator:      indikatorList,
	}
	updatedSasaranOpd, err := service.sasaranOpdRepository.Update(ctx, tx, sasaranOpdUpdate)
	if err != nil {
		return nil, err
	}
	return &sasaranopd.SasaranOpdCreateResponse{
		IdPohon:        updatedSasaranOpd.IdPohon,
		NamaSasaranOpd: updatedSasaranOpd.NamaSasaranOpd,
		NamaTujuanOpd:  tujuanOpd.Tujuan,
		TahunAwal:      updatedSasaranOpd.TahunAwal,
		TahunAkhir:     updatedSasaranOpd.TahunAkhir,
		JenisPeriode:   updatedSasaranOpd.JenisPeriode,
		Indikator:      indikatorResponses,
	}, nil
}

func (service *SasaranOpdServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	err = service.sasaranOpdRepository.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	return nil
}

func (service *SasaranOpdServiceImpl) FindByIdPokin(ctx context.Context, idPokin int, tahun string) (*sasaranopd.SasaranOpdResponse, error) {
	// Start transaction
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	// Ensure transaction is handled properly
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	sasaranOpd, err := service.sasaranOpdRepository.FindByIdPokin(ctx, tx, idPokin, tahun)
	if err != nil {
		return nil, fmt.Errorf("error finding sasaran opd: %v", err)
	}

	// Commit transaction before converting response
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	response := &sasaranopd.SasaranOpdResponse{
		IdPohon:    sasaranOpd.IdPohon,
		NamaPohon:  sasaranOpd.NamaPohon,
		JenisPohon: sasaranOpd.JenisPohon,
		LevelPohon: sasaranOpd.LevelPohon,
		TahunPohon: sasaranOpd.TahunPohon,
		Pelaksana:  make([]sasaranopd.PelaksanaOpdResponse, 0),
		SasaranOpd: make([]sasaranopd.SasaranOpdDetailResponse, 0),
	}

	// Convert response outside of transaction
	for _, pelaksana := range sasaranOpd.Pelaksana {
		response.Pelaksana = append(response.Pelaksana, sasaranopd.PelaksanaOpdResponse{
			Id:          pelaksana.Id,
			PegawaiId:   pelaksana.PegawaiId,
			Nip:         pelaksana.Nip,
			NamaPegawai: pelaksana.NamaPegawai,
		})
	}

	for _, sasaran := range sasaranOpd.SasaranOpd {
		sasaranResponse := sasaranopd.SasaranOpdDetailResponse{
			Id:             strconv.Itoa(sasaran.Id),
			NamaSasaranOpd: sasaran.NamaSasaranOpd,
			TahunAwal:      sasaran.TahunAwal,
			TahunAkhir:     sasaran.TahunAkhir,
			JenisPeriode:   sasaran.JenisPeriode,
			Indikator:      make([]sasaranopd.IndikatorResponse, 0),
		}

		// Convert dan urutkan indikator
		for _, indikator := range sasaran.Indikator {
			indResponse := sasaranopd.IndikatorResponse{
				Id:               indikator.Id,
				Indikator:        indikator.Indikator,
				RumusPerhitungan: indikator.RumusPerhitungan.String,
				SumberData:       indikator.SumberData.String,
				Target:           make([]sasaranopd.TargetResponse, 0),
			}

			// Convert dan urutkan target
			for _, target := range indikator.Target {
				indResponse.Target = append(indResponse.Target, sasaranopd.TargetResponse{
					Id:     target.Id,
					Tahun:  target.Tahun,
					Target: target.Target,
					Satuan: target.Satuan,
				})
			}

			// Urutkan target berdasarkan tahun
			sort.Slice(indResponse.Target, func(i, j int) bool {
				return indResponse.Target[i].Tahun < indResponse.Target[j].Tahun
			})

			sasaranResponse.Indikator = append(sasaranResponse.Indikator, indResponse)
		}

		// Urutkan indikator berdasarkan nama
		sort.Slice(sasaranResponse.Indikator, func(i, j int) bool {
			return sasaranResponse.Indikator[i].Indikator < sasaranResponse.Indikator[j].Indikator
		})

		response.SasaranOpd = append(response.SasaranOpd, sasaranResponse)
	}

	return response, nil
}

func (service *SasaranOpdServiceImpl) FindIdPokinSasaran(ctx context.Context, id int) (pohonkinerja.PohonKinerjaOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data pohon kinerja dengan sasaran
	pokin, err := service.sasaranOpdRepository.FindIdPokinSasaran(ctx, tx, id)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}

	// Ambil data OPD jika ada
	var namaOpd string
	if pokin.KodeOpd != "" {
		opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokin.KodeOpd)
		if err == nil {
			namaOpd = opd.NamaOpd
		}
	}

	// Konversi indikator dan target ke response
	var indikatorResponses []pohonkinerja.IndikatorResponse
	for _, ind := range pokin.Indikator {
		var targetResponses []pohonkinerja.TargetResponse

		// Cari target yang ada di database
		var existingTarget *domain.Target
		for _, t := range ind.Target {
			if t.Id != fmt.Sprintf("TRG-%s-%s", ind.Id, pokin.Tahun) {
				existingTarget = &t
				break
			}
		}

		// Jika ada target di database, gunakan itu
		if existingTarget != nil {
			targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
				Id:              existingTarget.Id,
				IndikatorId:     existingTarget.IndikatorId,
				TargetIndikator: existingTarget.Target,
				SatuanIndikator: existingTarget.Satuan,
				TahunSasaran:    pokin.Tahun,
			})
		} else {
			// Jika tidak ada target di database, gunakan target default
			targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
				Id:              fmt.Sprintf("TRG-%s-%s", ind.Id, pokin.Tahun),
				IndikatorId:     ind.Id,
				TargetIndikator: "-",
				SatuanIndikator: "-",
				TahunSasaran:    pokin.Tahun,
			})
		}

		indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
			Id:            ind.Id,
			IdPokin:       fmt.Sprint(pokin.Id),
			NamaIndikator: ind.Indikator,
			Target:        targetResponses,
		})
	}

	response := pohonkinerja.PohonKinerjaOpdResponse{
		Id:         pokin.Id,
		Parent:     strconv.Itoa(pokin.Parent),
		NamaPohon:  pokin.NamaPohon,
		JenisPohon: pokin.JenisPohon,
		LevelPohon: pokin.LevelPohon,
		KodeOpd:    pokin.KodeOpd,
		NamaOpd:    namaOpd,
		Keterangan: pokin.Keterangan,
		Tahun:      pokin.Tahun,
		Status:     pokin.Status,

		Indikator: indikatorResponses,
	}

	return response, nil
}

func (service *SasaranOpdServiceImpl) FindByTahun(ctx context.Context, kodeOpd string, tahun string, jenisPeriode string) ([]sasaranopd.SasaranOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi tahun
	tahunInt, err := strconv.Atoi(tahun)
	if err != nil {
		return nil, fmt.Errorf("format tahun tidak valid")
	}

	// Ambil data
	sasaranOpds, err := service.sasaranOpdRepository.FindByTahun(ctx, tx, kodeOpd, tahun, jenisPeriode)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]sasaranopd.SasaranOpdResponse, 0), nil
		}
		return nil, err
	}

	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, sasaranOpds[0].KodeOpd)
	if err != nil {
		return nil, err
	}

	var responses []sasaranopd.SasaranOpdResponse
	for _, sasaranOpd := range sasaranOpds {
		// Validasi tahun pokin terhadap range sasaran
		// tahunPokinInt, _ := strconv.Atoi(sasaranOpd.TahunPohon)

		response := sasaranopd.SasaranOpdResponse{
			IdPohon:    sasaranOpd.IdPohon,
			KodeOpd:    sasaranOpd.KodeOpd,
			NamaOpd:    opd.NamaOpd,
			NamaPohon:  sasaranOpd.NamaPohon,
			JenisPohon: sasaranOpd.JenisPohon,
			LevelPohon: sasaranOpd.LevelPohon,
			TahunPohon: sasaranOpd.TahunPohon,
			Pelaksana:  make([]sasaranopd.PelaksanaOpdResponse, 0),
			SasaranOpd: make([]sasaranopd.SasaranOpdDetailResponse, 0),
		}

		// Convert Pelaksana
		for _, pelaksana := range sasaranOpd.Pelaksana {
			response.Pelaksana = append(response.Pelaksana, sasaranopd.PelaksanaOpdResponse{
				Id:          pelaksana.Id,
				PegawaiId:   pelaksana.PegawaiId,
				Nip:         pelaksana.Nip,
				NamaPegawai: pelaksana.NamaPegawai,
			})
		}

		// Convert SasaranOpd
		for _, sasaran := range sasaranOpd.SasaranOpd {
			tahunAwalInt, _ := strconv.Atoi(sasaran.TahunAwal)
			tahunAkhirInt, _ := strconv.Atoi(sasaran.TahunAkhir)

			// Validasi tahun parameter dalam range sasaran
			if tahunInt < tahunAwalInt || tahunInt > tahunAkhirInt {
				continue
			}

			TujuanOpd, _ := service.tujuanOpdRepository.FindById(ctx, tx, sasaran.IdTujuanOpd)

			sasaranResponse := sasaranopd.SasaranOpdDetailResponse{
				Id:             strconv.Itoa(sasaran.Id),
				NamaSasaranOpd: sasaran.NamaSasaranOpd,
				IdTujuanOpd:    TujuanOpd.Id,
				NamaTujuanOpd:  TujuanOpd.Tujuan,
				TahunAwal:      sasaran.TahunAwal,
				TahunAkhir:     sasaran.TahunAkhir,
				JenisPeriode:   sasaran.JenisPeriode,
				Indikator:      make([]sasaranopd.IndikatorResponse, 0),
			}

			// Convert Indikator
			for _, indikator := range sasaran.Indikator {
				indResponse := sasaranopd.IndikatorResponse{
					Id:               indikator.Id,
					Indikator:        indikator.Indikator,
					RumusPerhitungan: indikator.RumusPerhitungan.String,
					SumberData:       indikator.SumberData.String,
					Target:           make([]sasaranopd.TargetResponse, 0),
				}

				// Hanya ambil target untuk tahun yang diminta
				for _, target := range indikator.Target {
					if target.Tahun == tahun {
						indResponse.Target = append(indResponse.Target, sasaranopd.TargetResponse{
							Id:     target.Id,
							Tahun:  target.Tahun,
							Target: target.Target,
							Satuan: target.Satuan,
						})
					}
				}

				if len(indResponse.Target) > 0 {
					sasaranResponse.Indikator = append(sasaranResponse.Indikator, indResponse)
				}
			}

			if len(sasaranResponse.Indikator) > 0 {
				response.SasaranOpd = append(response.SasaranOpd, sasaranResponse)
			}
		}

		if len(response.SasaranOpd) > 0 {
			responses = append(responses, response)
		}
	}

	return responses, nil
}

func (s *SasaranOpdServiceImpl) FindSasaranRenstra(
	ctx context.Context, kodeOpd, tahunAwal, tahunAkhir, jenisPeriode string,
) ([]sasaranopd.SasaranOpdResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	sasaranOpds, err := s.sasaranOpdRepository.FindAllOnly(ctx, tx, kodeOpd, tahunAwal, tahunAkhir, jenisPeriode)

	sasaranOpdIds := []int{}
	for _, so := range sasaranOpds {
		for _, sas := range so.SasaranOpd {
			sasaranOpdIds = append(sasaranOpdIds, sas.Id)
		}
	}
	indikatorSasaran, err := s.getIndikatorWithFallback(ctx, tx, sasaranOpdIds)
	if err != nil {
		log.Printf("ERROR service.getIndikatorWithFallback: %v", err)
		return nil, err
	}
	indSasaranById := make(map[int][]domain.Indikator)
	for _, ind := range indikatorSasaran {
		indSasaranById[ind.SasaranOpdId] = append(indSasaranById[ind.SasaranOpdId], ind)
	}

	for _, so := range sasaranOpds {

		for i := range so.SasaranOpd {
			id := so.SasaranOpd[i].Id
			if inds, ok := indSasaranById[id]; ok {
				so.SasaranOpd[i].Indikator = inds
			}
		}
	}

	return s.buildSasaranResponse(ctx, tx, kodeOpd, sasaranOpds)
}

func (s *SasaranOpdServiceImpl) FindSasaranRanwal(
	ctx context.Context, kodeOpd, tahun, jenisPeriode string,
) ([]sasaranopd.SasaranOpdResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	sasaranOpds, err := s.sasaranOpdRepository.FindSasaranByTahun(ctx, tx, kodeOpd, tahun, jenisPeriode, "ranwal")
	if err != nil {
		return nil, err
	}
	return s.buildSasaranResponse(ctx, tx, kodeOpd, sasaranOpds)
}
func (s *SasaranOpdServiceImpl) FindSasaranRankhir(
	ctx context.Context, kodeOpd, tahun, jenisPeriode string,
) ([]sasaranopd.SasaranOpdResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	sasaranOpds, err := s.sasaranOpdRepository.FindSasaranByTahun(ctx, tx, kodeOpd, tahun, jenisPeriode, "rankhir")
	if err != nil {
		return nil, err
	}
	return s.buildSasaranResponse(ctx, tx, kodeOpd, sasaranOpds)
}

func (s *SasaranOpdServiceImpl) FindSasaranPenetapan(
	ctx context.Context, kodeOpd, tahun, jenisPeriode string,
) ([]sasaranopd.SasaranOpdResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	sasaranOpds, err := s.sasaranOpdRepository.FindSasaranByTahun(ctx, tx, kodeOpd, tahun, jenisPeriode, "penetapan")
	if err != nil {
		return nil, err
	}
	return s.buildSasaranResponse(ctx, tx, kodeOpd, sasaranOpds)
}

// Helper bersama untuk build response (menghindari duplikasi)
func (s *SasaranOpdServiceImpl) buildSasaranResponse(
	ctx context.Context, tx *sql.Tx,
	kodeOpd string, sasaranOpds []domain.SasaranOpd,
) ([]sasaranopd.SasaranOpdResponse, error) {
	if len(sasaranOpds) == 0 {
		return []sasaranopd.SasaranOpdResponse{}, nil
	}
	opd, _ := s.opdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	// ── Batch fetch tujuan_opd (hindari N+1) ──────────────────
	tujuanCache := make(map[int]domain.TujuanOpd)
	for _, so := range sasaranOpds {
		for _, sas := range so.SasaranOpd {
			if sas.IdTujuanOpd == 0 {
				continue
			}
			if _, exists := tujuanCache[sas.IdTujuanOpd]; !exists {
				tujuan, err := s.tujuanOpdRepository.FindById(ctx, tx, sas.IdTujuanOpd)
				if err == nil {
					tujuanCache[sas.IdTujuanOpd] = tujuan
				}
			}
		}
	}
	var responses []sasaranopd.SasaranOpdResponse
	for _, so := range sasaranOpds {
		resp := sasaranopd.SasaranOpdResponse{
			IdPohon: so.IdPohon, KodeOpd: so.KodeOpd, NamaOpd: opd.NamaOpd,
			NamaPohon: so.NamaPohon, JenisPohon: so.JenisPohon,
			LevelPohon: so.LevelPohon, TahunPohon: so.TahunPohon,
			Pelaksana:  []sasaranopd.PelaksanaOpdResponse{},
			SasaranOpd: []sasaranopd.SasaranOpdDetailResponse{},
		}
		for _, pl := range so.Pelaksana {
			resp.Pelaksana = append(resp.Pelaksana, sasaranopd.PelaksanaOpdResponse{
				Id: pl.Id, PegawaiId: pl.PegawaiId, Nip: pl.Nip, NamaPegawai: pl.NamaPegawai,
			})
		}
		for _, sas := range so.SasaranOpd {
			tujuan := tujuanCache[sas.IdTujuanOpd] // pakai cache, tidak re-query
			sasResp := sasaranopd.SasaranOpdDetailResponse{
				Id: strconv.Itoa(sas.Id), NamaSasaranOpd: sas.NamaSasaranOpd,
				IdTujuanOpd: tujuan.Id, NamaTujuanOpd: tujuan.Tujuan,
				TahunAwal: sas.TahunAwal, TahunAkhir: sas.TahunAkhir,
				JenisPeriode: sas.JenisPeriode,
				Indikator:    []sasaranopd.IndikatorResponse{},
			}
			for _, ind := range sas.Indikator {
				indResp := sasaranopd.IndikatorResponse{
					Id: ind.KodeIndikator, KodeIndikator: ind.KodeIndikator,
					Jenis:               ind.Jenis,
					DefinisiOperasional: ind.DefinisiOperasional.String,
					Indikator:           ind.Indikator,
					RumusPerhitungan:    ind.RumusPerhitungan.String,
					SumberData:          ind.SumberData.String,
					Target:              []sasaranopd.TargetResponse{},
				}
				for _, t := range ind.Target {
					indResp.Target = append(indResp.Target, sasaranopd.TargetResponse{
						Id: t.Id, Tahun: t.Tahun, Target: t.Target, Satuan: t.Satuan,
					})
				}
				sasResp.Indikator = append(sasResp.Indikator, indResp)
			}
			resp.SasaranOpd = append(resp.SasaranOpd, sasResp)
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

func (service *SasaranOpdServiceImpl) CreateRenjaIndikator(
	ctx context.Context,
	sasaranOpdId int,
	jenis string,
	requests []sasaranopd.IndikatorCreateRequest,
) ([]sasaranopd.IndikatorResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	_, err = service.sasaranOpdRepository.FindById(ctx, tx, sasaranOpdId) // ← lowercase
	if err != nil {
		return nil, fmt.Errorf("sasaran opd id %d tidak ditemukan", sasaranOpdId)
	}
	var indikatorDomains []domain.Indikator
	var responses []sasaranopd.IndikatorResponse
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
		kodeIndikator := fmt.Sprintf("IND-SAR-%s", uuid.New().String()[:5])
		targetId := fmt.Sprintf("TRG-SAR-%s", uuid.New().String()[:5])
		ind := domain.Indikator{
			KodeIndikator:       kodeIndikator,
			Jenis:               jenis,
			DefinisiOperasional: sql.NullString{String: req.DefinisiOperasional, Valid: true},
			Indikator:           req.Indikator,
			RumusPerhitungan:    sql.NullString{String: req.RumusPerhitungan, Valid: true},
			SumberData:          sql.NullString{String: req.SumberData, Valid: true},
			Target: []domain.Target{{
				Id: targetId, IndikatorId: kodeIndikator,
				Target: req.Target[0].Target, Satuan: req.Target[0].Satuan, Tahun: req.Target[0].Tahun,
			}},
		}
		indikatorDomains = append(indikatorDomains, ind)
		responses = append(responses, sasaranopd.IndikatorResponse{
			Id:                  kodeIndikator,
			KodeIndikator:       kodeIndikator,
			Indikator:           req.Indikator,
			RumusPerhitungan:    req.RumusPerhitungan,
			SumberData:          req.SumberData,
			DefinisiOperasional: req.DefinisiOperasional,
			Jenis:               jenis,
			Target: []sasaranopd.TargetResponse{{
				Id:     targetId,
				Tahun:  req.Target[0].Tahun,
				Target: req.Target[0].Target,
				Satuan: req.Target[0].Satuan,
			}},
		})
	}
	if err := service.sasaranOpdRepository.CreateRenjaIndikator(ctx, tx, sasaranOpdId, indikatorDomains); err != nil {
		return nil, err
	}
	return responses, nil
}

func (service *SasaranOpdServiceImpl) UpdateRenjaIndikator(ctx context.Context, kodeIndikator string, jenis string, request sasaranopd.IndikatorUpdateRequest) (sasaranopd.IndikatorResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return sasaranopd.IndikatorResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	_, err = service.sasaranOpdRepository.FindIndikatorByKodeIndikator(ctx, tx, kodeIndikator)
	if err != nil {
		if err == sql.ErrNoRows {
			return sasaranopd.IndikatorResponse{}, fmt.Errorf("indikator dengan kode %s tidak ditemukan", kodeIndikator)
		}
		return sasaranopd.IndikatorResponse{}, err
	}
	if request.Indikator == "" {
		return sasaranopd.IndikatorResponse{}, fmt.Errorf("nama indikator tidak boleh kosong")
	}
	if len(request.Target) != 1 {
		return sasaranopd.IndikatorResponse{}, fmt.Errorf("harus memiliki tepat 1 target")
	}
	if request.Target[0].Target == "" {
		return sasaranopd.IndikatorResponse{}, fmt.Errorf("nilai target tidak boleh kosong")
	}
	if request.Target[0].Tahun == "" {
		return sasaranopd.IndikatorResponse{}, fmt.Errorf("tahun target tidak boleh kosong")
	}
	targetId := request.Target[0].Id
	if targetId == "" {
		targetId = fmt.Sprintf("TRG-SAR-%s", uuid.New().String()[:5])
	}
	ind := domain.Indikator{
		KodeIndikator:       kodeIndikator,
		Jenis:               jenis,
		DefinisiOperasional: sql.NullString{String: request.DefinisiOperasional, Valid: true},
		Indikator:           request.Indikator,
		RumusPerhitungan:    sql.NullString{String: request.RumusPerhitungan, Valid: true},
		SumberData:          sql.NullString{String: request.SumberData, Valid: true},
		Target: []domain.Target{{
			Id: targetId, IndikatorId: kodeIndikator,
			Target: request.Target[0].Target, Satuan: request.Target[0].Satuan, Tahun: request.Target[0].Tahun,
		}},
	}
	if err := service.sasaranOpdRepository.UpdateRenjaIndikator(ctx, tx, []domain.Indikator{ind}); err != nil {
		return sasaranopd.IndikatorResponse{}, err
	}
	return sasaranopd.IndikatorResponse{
		Id:                  kodeIndikator,
		KodeIndikator:       kodeIndikator,
		Indikator:           request.Indikator,
		RumusPerhitungan:    request.RumusPerhitungan,
		SumberData:          request.SumberData,
		DefinisiOperasional: request.DefinisiOperasional,
		Jenis:               jenis,
		Target: []sasaranopd.TargetResponse{{
			Id:     targetId,
			Tahun:  request.Target[0].Tahun,
			Target: request.Target[0].Target,
			Satuan: request.Target[0].Satuan,
		}},
	}, nil
}

func (service *SasaranOpdServiceImpl) DeleteRenjaIndikator(ctx context.Context, kodeIndikator string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)
	_, err = service.sasaranOpdRepository.FindIndikatorByKodeIndikator(ctx, tx, kodeIndikator) // ← lowercase
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("kode indikator %s tidak ditemukan", kodeIndikator)
		}
		return err
	}
	return service.sasaranOpdRepository.DeleteIndikatorTargetRenja(ctx, tx, kodeIndikator) // ← lowercase
}

func (s *SasaranOpdServiceImpl) getIndikatorWithFallback(
	ctx context.Context,
	tx *sql.Tx,
	sasaranIds []int,
) ([]domain.Indikator, error) {

	indikatorBaru, err := s.sasaranOpdRepository.
		FindIndikatorTargetsRenstraBySasaranIds(ctx, tx, sasaranIds)
	if err != nil {
		return nil, err
	}

	indikatorLama, err := s.sasaranOpdRepository.
		FindIndikatorTargetsBySasaranIds(ctx, tx, sasaranIds)
	if err != nil {
		return nil, err
	}

	return mergeIndikator(indikatorBaru, indikatorLama), nil
}
