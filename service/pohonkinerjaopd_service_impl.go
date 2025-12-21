package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"ekak_kabupaten_madiun/model/web/opdmaster"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type PohonKinerjaOpdServiceImpl struct {
	pohonKinerjaOpdRepository repository.PohonKinerjaRepository
	opdRepository             repository.OpdRepository
	pegawaiRepository         repository.PegawaiRepository
	tujuanOpdRepository       repository.TujuanOpdRepository
	crosscuttingOpdRepository repository.CrosscuttingOpdRepository
	reviewRepository          repository.ReviewRepository
	DB                        *sql.DB
	Validate                  *validator.Validate
	ProgramUnggulanRepository repository.ProgramUnggulanRepository
	RedisClient               *redis.Client
}

func NewPohonKinerjaOpdServiceImpl(pohonKinerjaOpdRepository repository.PohonKinerjaRepository, opdRepository repository.OpdRepository, pegawaiRepository repository.PegawaiRepository, tujuanOpdRepository repository.TujuanOpdRepository, crosscuttingOpdRepository repository.CrosscuttingOpdRepository, reviewRepository repository.ReviewRepository, DB *sql.DB, validate *validator.Validate, programUnggulanRepository repository.ProgramUnggulanRepository, redisClient *redis.Client) *PohonKinerjaOpdServiceImpl {
	return &PohonKinerjaOpdServiceImpl{
		pohonKinerjaOpdRepository: pohonKinerjaOpdRepository,
		opdRepository:             opdRepository,
		pegawaiRepository:         pegawaiRepository,
		tujuanOpdRepository:       tujuanOpdRepository,
		crosscuttingOpdRepository: crosscuttingOpdRepository,
		reviewRepository:          reviewRepository,
		DB:                        DB,
		Validate:                  validate,
		ProgramUnggulanRepository: programUnggulanRepository,
		RedisClient:               redisClient,
	}
}

func (service *PohonKinerjaOpdServiceImpl) Create(ctx context.Context, request pohonkinerja.PohonKinerjaCreateRequest) (pohonkinerja.PohonKinerjaOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi request
	if request.NamaPohon == "" {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("nama program tidak boleh kosong")
	}

	// Validasi kode OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOpd)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("kode opd tidak ditemukan")
	}
	if opd.KodeOpd == "" {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("kode opd tidak valid")
	}

	// Validasi dan persiapan data pelaksana
	var pelaksanaList []domain.PelaksanaPokin
	var pelaksanaResponses []pohonkinerja.PelaksanaOpdResponse

	for _, pelaksanaReq := range request.PelaksanaId {
		// Generate ID untuk pelaksana_pokin
		pelaksanaId := fmt.Sprintf("PLKS-%s", uuid.New().String()[:8])

		// Validasi setiap pelaksana
		pegawaiPelaksana, err := service.pegawaiRepository.FindById(ctx, tx, pelaksanaReq.PegawaiId)
		if err != nil {
			return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("pelaksana tidak ditemukan")
		}
		if pegawaiPelaksana.Id == "" {
			return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("pelaksana tidak ditemukan")
		}

		// Tambahkan ke list pelaksana
		pelaksanaList = append(pelaksanaList, domain.PelaksanaPokin{
			Id:        pelaksanaId,
			PegawaiId: pelaksanaReq.PegawaiId,
		})

		// Siapkan response pelaksana
		pelaksanaResponses = append(pelaksanaResponses, pohonkinerja.PelaksanaOpdResponse{
			Id:          pelaksanaId,
			PegawaiId:   pegawaiPelaksana.Id,
			NamaPegawai: pegawaiPelaksana.NamaPegawai,
		})
	}

	// Validasi dan persiapan data indikator dan target
	var indikatorList []domain.Indikator
	var indikatorResponses []pohonkinerja.IndikatorResponse

	for _, indikatorReq := range request.Indikator {
		// Generate ID untuk indikator
		indikatorId := fmt.Sprintf("IND-%s", uuid.New().String()[:8])

		var targetList []domain.Target
		var targetResponses []pohonkinerja.TargetResponse

		// Proses target untuk setiap indikator
		for _, targetReq := range indikatorReq.Target {
			targetId := fmt.Sprintf("TRG-%s", uuid.New().String()[:8])

			target := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorId,
				Target:      targetReq.Target,
				Satuan:      targetReq.Satuan,
				Tahun:       request.Tahun,
			}
			targetList = append(targetList, target)

			targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
				Id:              targetId,
				IndikatorId:     indikatorId,
				TargetIndikator: targetReq.Target,
				SatuanIndikator: targetReq.Satuan,
			})
		}

		indikator := domain.Indikator{
			Id:        indikatorId,
			Indikator: indikatorReq.NamaIndikator,
			Target:    targetList,
		}
		indikatorList = append(indikatorList, indikator)

		indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
			Id:            indikatorId,
			NamaIndikator: indikatorReq.NamaIndikator,
			Target:        targetResponses,
		})
	}

	// Persiapkan data tagging
	var taggingList []domain.TaggingPokin
	var taggingResponses []pohonkinerja.TaggingResponse

	for _, tagging := range request.TaggingPokin {
		var keteranganList []domain.KeteranganTagging
		var keteranganResponses []pohonkinerja.KeteranganTaggingResponse
		for _, keterangan := range tagging.KeteranganTaggingProgram {
			// Ambil detail program unggulan
			programUnggulan, err := service.ProgramUnggulanRepository.FindByKodeProgramUnggulan(ctx, tx, keterangan.KodeProgramUnggulan)
			if err != nil {
				continue
			}

			// Tambahkan ke list domain
			keteranganList = append(keteranganList, domain.KeteranganTagging{
				KodeProgramUnggulan: keterangan.KodeProgramUnggulan,
				Tahun:               keterangan.Tahun,
			})

			// Tambahkan ke response
			keteranganResponses = append(keteranganResponses, pohonkinerja.KeteranganTaggingResponse{
				KodeProgramUnggulan: keterangan.KodeProgramUnggulan,
				RencanaImplementasi: programUnggulan.KeteranganProgramUnggulan,
				Tahun:               request.Tahun, // Tambahkan tahun ke response
			})
		}

		// Tambahkan ke list domain
		taggingList = append(taggingList, domain.TaggingPokin{
			NamaTagging:              tagging.NamaTagging,
			KeteranganTaggingProgram: keteranganList,
		})

		// Tambahkan ke response
		taggingResponses = append(taggingResponses, pohonkinerja.TaggingResponse{
			NamaTagging:              tagging.NamaTagging,
			KeteranganTaggingProgram: keteranganResponses,
		})
	}

	pohonKinerja := domain.PohonKinerja{
		NamaPohon:    request.NamaPohon,
		Parent:       request.Parent,
		JenisPohon:   request.JenisPohon,
		LevelPohon:   request.LevelPohon,
		KodeOpd:      request.KodeOpd,
		Keterangan:   request.Keterangan,
		Tahun:        request.Tahun,
		Status:       request.Status,
		Pelaksana:    pelaksanaList,
		Indikator:    indikatorList,
		TaggingPokin: taggingList,
	}

	result, err := service.pohonKinerjaOpdRepository.Create(ctx, tx, pohonKinerja)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}

	// Update tagging responses dengan ID yang sudah di-generate
	for i, tagging := range result.TaggingPokin {
		if i < len(taggingResponses) {
			taggingResponses[i].Id = tagging.Id
			taggingResponses[i].IdPokin = tagging.IdPokin

			// Update keterangan responses dengan ID yang sudah di-generate
			for j, keterangan := range tagging.KeteranganTaggingProgram {
				if j < len(taggingResponses[i].KeteranganTaggingProgram) {
					taggingResponses[i].KeteranganTaggingProgram[j].Id = keterangan.Id
					taggingResponses[i].KeteranganTaggingProgram[j].IdTagging = keterangan.IdTagging
				}
			}
		}
	}

	countReview, err := service.reviewRepository.CountReviewByPohonKinerja(ctx, tx, result.Id)
	helper.PanicIfError(err)

	response := pohonkinerja.PohonKinerjaOpdResponse{
		Id:          result.Id,
		Parent:      strconv.Itoa(result.Parent),
		NamaPohon:   result.NamaPohon,
		JenisPohon:  result.JenisPohon,
		LevelPohon:  result.LevelPohon,
		KodeOpd:     result.KodeOpd,
		NamaOpd:     opd.NamaOpd,
		Keterangan:  result.Keterangan,
		Tahun:       result.Tahun,
		Status:      result.Status,
		CountReview: countReview,
		Pelaksana:   pelaksanaResponses,
		Indikator:   indikatorResponses,
		Tagging:     taggingResponses,
	}

	return response, nil
}

func (service *PohonKinerjaOpdServiceImpl) Update(ctx context.Context, request pohonkinerja.PohonKinerjaUpdateRequest) (pohonkinerja.PohonKinerjaOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi request
	if request.NamaPohon == "" {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("nama program tidak boleh kosong")
	}

	// Validasi kode OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOpd)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("kode opd tidak ditemukan")
	}
	if opd.KodeOpd == "" {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("kode opd tidak valid")
	}

	// Cek apakah ini adalah pohon kinerja yang di-clone
	cloneFrom, err := service.pohonKinerjaOpdRepository.CheckCloneFrom(ctx, tx, request.Id)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}

	// Jika ini adalah pohon kinerja yang di-clone, tidak boleh diupdate
	if cloneFrom != 0 {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("tidak dapat mengupdate pohon kinerja yang merupakan hasil clone")
	}

	// Dapatkan semua pohon kinerja yang terkait (asli dan clone)
	var pokinsToUpdate []domain.PohonKinerja

	// Tambahkan pohon kinerja yang sedang diupdate
	existingPokin, err := service.pohonKinerjaOpdRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("data pohon kinerja tidak ditemukan")
	}
	pokinsToUpdate = append(pokinsToUpdate, existingPokin)

	// Cari pohon kinerja yang merupakan clone dari yang sedang diupdate
	clonedPokins, err := service.pohonKinerjaOpdRepository.FindPokinByCloneFrom(ctx, tx, request.Id)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}
	pokinsToUpdate = append(pokinsToUpdate, clonedPokins...)

	// Persiapkan data pelaksana
	var pelaksanaList []domain.PelaksanaPokin
	var pelaksanaResponses []pohonkinerja.PelaksanaOpdResponse

	for _, pelaksanaReq := range request.PelaksanaId {
		pelaksanaId := fmt.Sprintf("PLKS-%s", uuid.New().String()[:8])
		pegawaiPelaksana, err := service.pegawaiRepository.FindById(ctx, tx, pelaksanaReq.PegawaiId)
		if err != nil {
			return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("pelaksana tidak ditemukan")
		}

		pelaksanaList = append(pelaksanaList, domain.PelaksanaPokin{
			Id:        pelaksanaId,
			PegawaiId: pelaksanaReq.PegawaiId,
		})

		pelaksanaResponses = append(pelaksanaResponses, pohonkinerja.PelaksanaOpdResponse{
			Id:          pelaksanaId,
			PegawaiId:   pegawaiPelaksana.Id,
			NamaPegawai: pegawaiPelaksana.NamaPegawai,
		})
	}

	// Persiapkan response indikator di luar loop
	var indikatorResponses []pohonkinerja.IndikatorResponse

	// Update untuk setiap pohon kinerja (asli dan clone)
	var updatedPokin domain.PohonKinerja
	for _, pokin := range pokinsToUpdate {
		var indikatorList []domain.Indikator

		for _, indikatorReq := range request.Indikator {
			var indikatorId string
			var cloneFromIndikator string

			if pokin.Id == request.Id {
				// Untuk pohon asli, gunakan ID dari request
				indikatorId = indikatorReq.Id
				if indikatorId == "" {
					indikatorId = fmt.Sprintf("IND-%s", uuid.New().String()[:8])
				}
				cloneFromIndikator = ""
			} else {
				// Untuk pohon clone, cari ID indikator yang sudah ada berdasarkan clone_from
				existingIndikator, err := service.pohonKinerjaOpdRepository.FindIndikatorByCloneFrom(ctx, tx, pokin.Id, indikatorReq.Id)
				if err == nil && existingIndikator.Id != "" {
					// Gunakan ID yang sudah ada jika ditemukan
					indikatorId = existingIndikator.Id
				} else {
					// Buat ID baru jika belum ada
					indikatorId = fmt.Sprintf("IND-%s", uuid.New().String()[:8])
				}
				cloneFromIndikator = indikatorReq.Id
			}

			var targetList []domain.Target
			var targetResponses []pohonkinerja.TargetResponse

			for _, targetReq := range indikatorReq.Target {
				var targetId string
				var cloneFromTarget string

				if pokin.Id == request.Id {
					// Untuk pohon asli, gunakan ID dari request
					targetId = targetReq.Id
					if targetId == "" {
						targetId = fmt.Sprintf("TRG-%s", uuid.New().String()[:8])
					}
					cloneFromTarget = ""
				} else {
					// Untuk pohon clone, cari ID target yang sudah ada berdasarkan clone_from
					existingTarget, err := service.pohonKinerjaOpdRepository.FindTargetByCloneFrom(ctx, tx, indikatorId, targetReq.Id)
					if err == nil && existingTarget.Id != "" {
						// Gunakan ID yang sudah ada jika ditemukan
						targetId = existingTarget.Id
					} else {
						// Buat ID baru jika belum ada
						targetId = fmt.Sprintf("TRG-%s", uuid.New().String()[:8])
					}
					cloneFromTarget = targetReq.Id
				}

				target := domain.Target{
					Id:          targetId,
					IndikatorId: indikatorId,
					Target:      targetReq.Target,
					Satuan:      targetReq.Satuan,
					Tahun:       request.Tahun,
					CloneFrom:   cloneFromTarget,
				}
				targetList = append(targetList, target)

				if pokin.Id == request.Id {
					targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
						Id:              targetId,
						IndikatorId:     indikatorId,
						TargetIndikator: targetReq.Target,
						SatuanIndikator: targetReq.Satuan,
					})
				}
			}

			indikator := domain.Indikator{
				Id:        indikatorId,
				PokinId:   fmt.Sprint(pokin.Id),
				Indikator: indikatorReq.NamaIndikator,
				Tahun:     request.Tahun,
				Target:    targetList,
				CloneFrom: cloneFromIndikator,
			}
			indikatorList = append(indikatorList, indikator)

			if pokin.Id == request.Id {
				indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
					Id:            indikatorId,
					IdPokin:       fmt.Sprint(pokin.Id),
					NamaIndikator: indikatorReq.NamaIndikator,
					Target:        targetResponses,
				})
			}
		}

		pohonKinerjaUpdate := domain.PohonKinerja{
			Id:                     pokin.Id,
			NamaPohon:              request.NamaPohon,
			Parent:                 request.Parent,
			JenisPohon:             request.JenisPohon,
			LevelPohon:             request.LevelPohon,
			KodeOpd:                request.KodeOpd,
			Keterangan:             request.Keterangan,
			Tahun:                  request.Tahun,
			Status:                 pokin.Status,
			CloneFrom:              pokin.CloneFrom,
			Pelaksana:              pelaksanaList,
			Indikator:              indikatorList,
			KeteranganCrosscutting: pokin.KeteranganCrosscutting,
			UpdatedBy:              request.UpdatedBy,
		}

		result, err := service.pohonKinerjaOpdRepository.Update(ctx, tx, pohonKinerjaUpdate)
		if err != nil {
			return pohonkinerja.PohonKinerjaOpdResponse{}, err
		}

		if pokin.Id == request.Id {
			updatedPokin = result
		}
	}

	countReview, err := service.reviewRepository.CountReviewByPohonKinerja(ctx, tx, updatedPokin.Id)
	helper.PanicIfError(err)

	// Proses tagging
	var taggingList []domain.TaggingPokin
	var taggingResponses []pohonkinerja.TaggingResponse

	// Update untuk tagging asli
	for _, taggingReq := range request.TaggingPokin {
		var keteranganList []domain.KeteranganTagging
		for _, keteranganReq := range taggingReq.KeteranganTaggingProgram {
			keteranganList = append(keteranganList, domain.KeteranganTagging{
				KodeProgramUnggulan: keteranganReq.KodeProgramUnggulan,
				Tahun:               keteranganReq.Tahun,
			})
		}

		tagging := domain.TaggingPokin{
			Id:                       taggingReq.Id,
			IdPokin:                  existingPokin.Id,
			NamaTagging:              taggingReq.NamaTagging,
			KeteranganTaggingProgram: keteranganList,
		}
		taggingList = append(taggingList, tagging)
	}

	// Update tagging untuk pohon asli
	taggingResults, err := service.pohonKinerjaOpdRepository.UpdateTagging(ctx, tx, existingPokin.Id, taggingList)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}

	// Cari dan update tagging untuk pohon yang di-clone
	clonedPokins, err = service.pohonKinerjaOpdRepository.FindPokinByCloneFrom(ctx, tx, request.Id)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}

	// Update tagging untuk setiap pohon yang di-clone
	for _, clonedPokin := range clonedPokins {
		var clonedTaggingList []domain.TaggingPokin

		// Ambil tagging yang ada di pohon yang di-clone
		existingClonedTaggings, err := service.pohonKinerjaOpdRepository.FindTaggingByPokinId(ctx, tx, clonedPokin.Id)
		if err != nil {
			return pohonkinerja.PohonKinerjaOpdResponse{}, err
		}

		// Buat map untuk mempermudah pencarian tagging berdasarkan clone_from
		clonedTaggingMap := make(map[int]domain.TaggingPokin)
		for _, tag := range existingClonedTaggings {
			clonedTaggingMap[tag.CloneFrom] = tag
		}

		// Update setiap tagging yang sesuai
		for _, originalTagging := range taggingResults {
			if clonedTagging, exists := clonedTaggingMap[originalTagging.Id]; exists {
				// Update tagging yang sudah ada
				var keteranganList []domain.KeteranganTagging
				for _, keterangan := range originalTagging.KeteranganTaggingProgram {
					keteranganList = append(keteranganList, domain.KeteranganTagging{
						KodeProgramUnggulan: keterangan.KodeProgramUnggulan,
						Tahun:               keterangan.Tahun,
					})
				}

				clonedTagging.NamaTagging = originalTagging.NamaTagging
				clonedTagging.KeteranganTaggingProgram = keteranganList
				clonedTaggingList = append(clonedTaggingList, clonedTagging)
			} else {
				// Buat tagging baru jika belum ada
				var keteranganList []domain.KeteranganTagging
				for _, keterangan := range originalTagging.KeteranganTaggingProgram {
					keteranganList = append(keteranganList, domain.KeteranganTagging{
						KodeProgramUnggulan: keterangan.KodeProgramUnggulan,
					})
				}

				newClonedTagging := domain.TaggingPokin{
					IdPokin:                  clonedPokin.Id,
					NamaTagging:              originalTagging.NamaTagging,
					KeteranganTaggingProgram: keteranganList,
					CloneFrom:                originalTagging.Id,
				}
				clonedTaggingList = append(clonedTaggingList, newClonedTagging)
			}
		}

		// Update tagging untuk pohon yang di-clone
		_, err = service.pohonKinerjaOpdRepository.UpdateTagging(ctx, tx, clonedPokin.Id, clonedTaggingList)
		if err != nil {
			return pohonkinerja.PohonKinerjaOpdResponse{}, err
		}
	}

	// Konversi ke response
	for _, tagging := range taggingResults {
		var keteranganResponses []pohonkinerja.KeteranganTaggingResponse
		for _, keterangan := range tagging.KeteranganTaggingProgram {
			programUnggulan, err := service.ProgramUnggulanRepository.FindByKodeProgramUnggulan(ctx, tx, keterangan.KodeProgramUnggulan)
			if err != nil {
				continue
			}
			keteranganResponses = append(keteranganResponses, pohonkinerja.KeteranganTaggingResponse{
				Id:                  keterangan.Id,
				IdTagging:           keterangan.IdTagging,
				KodeProgramUnggulan: keterangan.KodeProgramUnggulan,
				RencanaImplementasi: programUnggulan.KeteranganProgramUnggulan,
				Tahun:               keterangan.Tahun,
			})
		}

		taggingResponses = append(taggingResponses, pohonkinerja.TaggingResponse{
			Id:                       tagging.Id,
			IdPokin:                  tagging.IdPokin,
			NamaTagging:              tagging.NamaTagging,
			KeteranganTaggingProgram: keteranganResponses,
		})
	}

	return pohonkinerja.PohonKinerjaOpdResponse{
		Id:                     updatedPokin.Id,
		Parent:                 strconv.Itoa(updatedPokin.Parent),
		NamaPohon:              updatedPokin.NamaPohon,
		JenisPohon:             updatedPokin.JenisPohon,
		LevelPohon:             updatedPokin.LevelPohon,
		KodeOpd:                updatedPokin.KodeOpd,
		NamaOpd:                opd.NamaOpd,
		Keterangan:             updatedPokin.Keterangan,
		Tahun:                  updatedPokin.Tahun,
		CountReview:            countReview,
		Status:                 updatedPokin.Status,
		Pelaksana:              pelaksanaResponses,
		Indikator:              indikatorResponses,
		Tagging:                taggingResponses,
		KeteranganCrosscutting: updatedPokin.KeteranganCrosscutting,
		UpdatedBy:              updatedPokin.UpdatedBy,
	}, nil
}

func (service *PohonKinerjaOpdServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// 1. Cek apakah pohon kinerja dengan ID tersebut ada
	_, err = service.pohonKinerjaOpdRepository.FindById(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("pohon kinerja tidak ditemukan: %v", err)
	}

	// 2. Lakukan penghapusan dengan fungsi baru
	err = service.pohonKinerjaOpdRepository.Delete(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus pohon kinerja: %v", err)
	}

	return nil
}

func (service *PohonKinerjaOpdServiceImpl) FindById(ctx context.Context, id int) (pohonkinerja.PohonKinerjaOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// 1. Ambil data pohon kinerja
	pokin, err := service.pohonKinerjaOpdRepository.FindById(ctx, tx, id)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}

	// 2. Validasi data pohon kinerja
	if pokin.Id == 0 {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("data tidak ditemukan")
	}

	// 3. Ambil data OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokin.KodeOpd)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("data opd tidak ditemukan")
	}

	// 4. Ambil data pelaksana
	pelaksanaList, err := service.pohonKinerjaOpdRepository.FindPelaksanaPokin(ctx, tx, fmt.Sprint(pokin.Id))
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("gagal mengambil data pelaksana")
	}

	// 5. Proses data pelaksana
	var pelaksanaResponses []pohonkinerja.PelaksanaOpdResponse
	for _, pelaksana := range pelaksanaList {
		pegawaiPelaksana, err := service.pegawaiRepository.FindById(ctx, tx, pelaksana.PegawaiId)
		if err != nil {
			continue // Skip jika pegawai tidak ditemukan
		}

		pelaksanaResponses = append(pelaksanaResponses, pohonkinerja.PelaksanaOpdResponse{
			Id:          pelaksana.Id,
			PegawaiId:   pegawaiPelaksana.Id,
			Nip:         pegawaiPelaksana.Nip,
			NamaPegawai: pegawaiPelaksana.NamaPegawai,
		})
	}

	// 6. Ambil data indikator dan target
	var indikatorResponses []pohonkinerja.IndikatorResponse
	indikatorList, err := service.pohonKinerjaOpdRepository.FindIndikatorByPokinId(ctx, tx, fmt.Sprint(pokin.Id))
	if err == nil {
		for _, indikator := range indikatorList {
			// Ambil target untuk setiap indikator
			targetList, err := service.pohonKinerjaOpdRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
			if err != nil {
				continue
			}

			var targetResponses []pohonkinerja.TargetResponse
			for _, target := range targetList {
				targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
					Id:              target.Id,
					IndikatorId:     target.IndikatorId,
					TargetIndikator: target.Target,
					SatuanIndikator: target.Satuan,
				})
			}

			indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
				Id:            indikator.Id,
				IdPokin:       indikator.PokinId,
				NamaIndikator: indikator.Indikator,
				Target:        targetResponses,
			})
		}
	}

	// Tambahkan: Ambil data tagging
	var taggingResponses []pohonkinerja.TaggingResponse
	taggingList, err := service.pohonKinerjaOpdRepository.FindTaggingByPokinId(ctx, tx, pokin.Id)
	if err == nil {
		for _, tagging := range taggingList {
			// Konversi keterangan program ke response
			var keteranganResponses []pohonkinerja.KeteranganTaggingResponse
			for _, keterangan := range tagging.KeteranganTaggingProgram {
				programUnggulan, err := service.ProgramUnggulanRepository.FindByKodeProgramUnggulan(ctx, tx, keterangan.KodeProgramUnggulan)
				if err != nil {
					continue
				}
				keteranganResponses = append(keteranganResponses, pohonkinerja.KeteranganTaggingResponse{
					Id:                  keterangan.Id,
					IdTagging:           keterangan.IdTagging,
					KodeProgramUnggulan: keterangan.KodeProgramUnggulan,
					RencanaImplementasi: programUnggulan.KeteranganProgramUnggulan,
					Tahun:               keterangan.Tahun,
				})
			}

			taggingResponses = append(taggingResponses, pohonkinerja.TaggingResponse{
				Id:                       tagging.Id,
				IdPokin:                  tagging.IdPokin,
				NamaTagging:              tagging.NamaTagging,
				KeteranganTaggingProgram: keteranganResponses,
			})
		}
	}

	// Susun response
	response := pohonkinerja.PohonKinerjaOpdResponse{
		Id:         pokin.Id,
		Parent:     strconv.Itoa(pokin.Parent),
		NamaPohon:  pokin.NamaPohon,
		JenisPohon: pokin.JenisPohon,
		LevelPohon: pokin.LevelPohon,
		KodeOpd:    pokin.KodeOpd,
		NamaOpd:    opd.NamaOpd,
		Keterangan: pokin.Keterangan,
		Tahun:      pokin.Tahun,
		Status:     pokin.Status,
		Pelaksana:  pelaksanaResponses,
		Indikator:  indikatorResponses,
		Tagging:    taggingResponses,
	}

	return response, nil
}

//find all lama
// func (service *PohonKinerjaOpdServiceImpl) FindAll(ctx context.Context, kodeOpd, tahun string) (pohonkinerja.PohonKinerjaOpdAllResponse, error) {
// 	tx, err := service.DB.Begin()
// 	if err != nil {
// 		return pohonkinerja.PohonKinerjaOpdAllResponse{}, err
// 	}
// 	defer helper.CommitOrRollback(tx)

// 	// Validasi OPD
// 	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
// 	if err != nil {
// 		return pohonkinerja.PohonKinerjaOpdAllResponse{}, errors.New("kode opd tidak ditemukan")
// 	}

// 	// Inisialisasi response dasar
// 	response := pohonkinerja.PohonKinerjaOpdAllResponse{
// 		KodeOpd:    kodeOpd,
// 		NamaOpd:    opd.NamaOpd,
// 		Tahun:      tahun,
// 		TujuanOpd:  make([]pohonkinerja.TujuanOpdResponse, 0),
// 		Strategics: make([]pohonkinerja.StrategicOpdResponse, 0),
// 	}

// 	// Ambil data tujuan OPD
// 	tujuanOpds, err := service.tujuanOpdRepository.FindTujuanOpdByTahun(ctx, tx, kodeOpd, tahun, "RPJMD")
// 	if err != nil {
// 		log.Printf("Error getting tujuan OPD: %v", err)
// 		// Kembalikan response dengan array kosong jika terjadi error
// 		return response, nil
// 	}

// 	// Konversi tujuan OPD ke format response
// 	for _, tujuan := range tujuanOpds {
// 		indikators, err := service.tujuanOpdRepository.FindIndikatorByTujuanOpdId(ctx, tx, tujuan.Id)
// 		if err != nil {
// 			log.Printf("Error getting indikator for tujuan ID %d: %v", tujuan.Id, err)
// 			continue
// 		}

// 		var indikatorResponses []pohonkinerja.IndikatorTujuanResponse
// 		for _, indikator := range indikators {
// 			indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorTujuanResponse{
// 				Indikator: indikator.Indikator,
// 			})
// 		}

// 		response.TujuanOpd = append(response.TujuanOpd, pohonkinerja.TujuanOpdResponse{
// 			Id:        tujuan.Id,
// 			KodeOpd:   tujuan.KodeOpd,
// 			Tujuan:    tujuan.Tujuan,
// 			Indikator: indikatorResponses,
// 		})
// 	}

// 	// Ambil data pohon kinerja
// 	pokins, err := service.pohonKinerjaOpdRepository.FindAll(ctx, tx, kodeOpd, tahun)
// 	if err != nil {
// 		// Kembalikan response dengan data yang sudah ada jika terjadi error
// 		return response, nil
// 	}

// 	// Jika tidak ada data pohon kinerja, kembalikan response dengan array kosong
// 	if len(pokins) == 0 {
// 		return response, nil
// 	}

// 	// Proses data pohon kinerja seperti sebelumnya
// 	pohonMap := make(map[int]map[int][]domain.PohonKinerja)
// 	pelaksanaMap := make(map[int][]pohonkinerja.PelaksanaOpdResponse)
// 	indikatorMap := make(map[int][]pohonkinerja.IndikatorResponse)

// 	// Kelompokkan data dan ambil data pelaksana & indikator
// 	maxLevel := 0
// 	for _, p := range pokins {
// 		if p.LevelPohon >= 4 {
// 			// Update max level jika ditemukan level yang lebih tinggi
// 			if p.LevelPohon > maxLevel {
// 				maxLevel = p.LevelPohon
// 			}

// 			// Inisialisasi map untuk level jika belum ada
// 			if pohonMap[p.LevelPohon] == nil {
// 				pohonMap[p.LevelPohon] = make(map[int][]domain.PohonKinerja)
// 			}

// 			p.NamaOpd = opd.NamaOpd
// 			pohonMap[p.LevelPohon][p.Parent] = append(
// 				pohonMap[p.LevelPohon][p.Parent],
// 				p,
// 			)

// 			// Ambil data pelaksana
// 			pelaksanaList, err := service.pohonKinerjaOpdRepository.FindPelaksanaPokin(ctx, tx, fmt.Sprint(p.Id))
// 			if err == nil {
// 				var pelaksanaResponses []pohonkinerja.PelaksanaOpdResponse
// 				for _, pelaksana := range pelaksanaList {
// 					pegawaiPelaksana, err := service.pegawaiRepository.FindById(ctx, tx, pelaksana.PegawaiId)
// 					if err != nil {
// 						continue
// 					}
// 					pelaksanaResponses = append(pelaksanaResponses, pohonkinerja.PelaksanaOpdResponse{
// 						Id:          pelaksana.Id,
// 						PegawaiId:   pegawaiPelaksana.Id,
// 						NamaPegawai: pegawaiPelaksana.NamaPegawai,
// 					})
// 				}
// 				pelaksanaMap[p.Id] = pelaksanaResponses
// 			}

// 			// Ambil data indikator dan target
// 			indikatorList, err := service.pohonKinerjaOpdRepository.FindIndikatorByPokinId(ctx, tx, fmt.Sprint(p.Id))
// 			if err == nil {
// 				var indikatorResponses []pohonkinerja.IndikatorResponse
// 				for _, indikator := range indikatorList {
// 					// Ambil target untuk setiap indikator
// 					targetList, err := service.pohonKinerjaOpdRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
// 					if err != nil {
// 						continue
// 					}

// 					var targetResponses []pohonkinerja.TargetResponse
// 					for _, target := range targetList {
// 						targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 							Id:              target.Id,
// 							IndikatorId:     target.IndikatorId,
// 							TargetIndikator: target.Target,
// 							SatuanIndikator: target.Satuan,
// 						})
// 					}

// 					indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
// 						Id:            indikator.Id,
// 						IdPokin:       indikator.PokinId,
// 						NamaIndikator: indikator.Indikator,
// 						Target:        targetResponses,
// 					})
// 				}
// 				indikatorMap[p.Id] = indikatorResponses
// 			}
// 		}
// 	}

// 	// Build response untuk strategic (level 4)
// 	if strategicList := pohonMap[4]; len(strategicList) > 0 {
// 		for _, strategicsByParent := range strategicList {
// 			sort.Slice(strategicsByParent, func(i, j int) bool {
// 				return strategicsByParent[i].Id < strategicsByParent[j].Id
// 			})

// 			for _, strategic := range strategicsByParent {
// 				strategicResp := service.buildStrategicResponse(ctx, tx, pohonMap, strategic, pelaksanaMap, indikatorMap)
// 				response.Strategics = append(response.Strategics, strategicResp)
// 			}
// 		}

// 		// Urutkan strategics berdasarkan Id
// 		sort.Slice(response.Strategics, func(i, j int) bool {
// 			return response.Strategics[i].Id < response.Strategics[j].Id
// 		})
// 	}

// 	return response, nil
// }

// findall baru
// func (service *PohonKinerjaOpdServiceImpl) FindAll(ctx context.Context, kodeOpd, tahun string) (pohonkinerja.PohonKinerjaOpdAllResponse, error) {
// 	tx, err := service.DB.Begin()
// 	if err != nil {
// 		return pohonkinerja.PohonKinerjaOpdAllResponse{}, err
// 	}
// 	defer helper.CommitOrRollback(tx)

// 	// Validasi OPD
// 	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
// 	if err != nil {
// 		return pohonkinerja.PohonKinerjaOpdAllResponse{}, errors.New("kode opd tidak ditemukan")
// 	}

// 	// Inisialisasi response dasar
// 	response := pohonkinerja.PohonKinerjaOpdAllResponse{
// 		KodeOpd:    kodeOpd,
// 		NamaOpd:    opd.NamaOpd,
// 		Tahun:      tahun,
// 		TujuanOpd:  make([]pohonkinerja.TujuanOpdResponse, 0),
// 		Strategics: make([]pohonkinerja.StrategicOpdResponse, 0),
// 	}

// 	// Ambil data tujuan OPD
// 	tujuanOpds, err := service.tujuanOpdRepository.FindTujuanOpdByTahun(ctx, tx, kodeOpd, tahun, "RPJMD")
// 	if err != nil {
// 		log.Printf("Error getting tujuan OPD: %v", err)
// 		return response, nil
// 	}

// 	// Proses data tujuan OPD
// 	for _, tujuan := range tujuanOpds {
// 		indikators, err := service.tujuanOpdRepository.FindIndikatorByTujuanOpdId(ctx, tx, tujuan.Id)
// 		if err != nil {
// 			log.Printf("Error getting indikator for tujuan ID %d: %v", tujuan.Id, err)
// 			continue
// 		}

// 		var indikatorResponses []pohonkinerja.IndikatorTujuanResponse
// 		for _, indikator := range indikators {
// 			targets, err := service.tujuanOpdRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id, tahun)
// 			if err != nil {
// 				log.Printf("Error getting targets for indikator ID %s: %v", indikator.Id, err)
// 				continue
// 			}

// 			var targetResponses []pohonkinerja.TargetTujuanResponse
// 			for _, target := range targets {
// 				targetResponses = append(targetResponses, pohonkinerja.TargetTujuanResponse{
// 					Tahun:  target.Tahun,
// 					Target: target.Target,
// 					Satuan: target.Satuan,
// 				})
// 			}

// 			indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorTujuanResponse{
// 				Indikator: indikator.Indikator,
// 				Target:    targetResponses,
// 			})
// 		}

// 		response.TujuanOpd = append(response.TujuanOpd, pohonkinerja.TujuanOpdResponse{
// 			Id:        tujuan.Id,
// 			KodeOpd:   tujuan.KodeOpd,
// 			Tujuan:    tujuan.Tujuan,
// 			Indikator: indikatorResponses,
// 		})
// 	}

// 	// Ambil data pohon kinerja
// 	pokins, err := service.pohonKinerjaOpdRepository.FindAll(ctx, tx, kodeOpd, tahun)
// 	if err != nil {
// 		return response, nil
// 	}

// 	if len(pokins) == 0 {
// 		return response, nil
// 	}

// 	// Proses data pohon kinerja
// 	pohonMap := make(map[int]map[int][]domain.PohonKinerja)
// 	pelaksanaMap := make(map[int][]pohonkinerja.PelaksanaOpdResponse)
// 	indikatorMap := make(map[int][]pohonkinerja.IndikatorResponse)

// 	// Kelompokkan data dan ambil data pelaksana & indikator
// 	maxLevel := 0
// 	for _, p := range pokins {
// 		if p.LevelPohon >= 4 {
// 			if p.LevelPohon > maxLevel {
// 				maxLevel = p.LevelPohon
// 			}

// 			if pohonMap[p.LevelPohon] == nil {
// 				pohonMap[p.LevelPohon] = make(map[int][]domain.PohonKinerja)
// 			}

// 			p.NamaOpd = opd.NamaOpd
// 			pohonMap[p.LevelPohon][p.Parent] = append(
// 				pohonMap[p.LevelPohon][p.Parent],
// 				p,
// 			)

// 			// Ambil data pelaksana
// 			pelaksanaList, err := service.pohonKinerjaOpdRepository.FindPelaksanaPokin(ctx, tx, fmt.Sprint(p.Id))
// 			if err == nil {
// 				var pelaksanaResponses []pohonkinerja.PelaksanaOpdResponse
// 				for _, pelaksana := range pelaksanaList {
// 					pegawaiPelaksana, err := service.pegawaiRepository.FindById(ctx, tx, pelaksana.PegawaiId)
// 					if err != nil {
// 						continue
// 					}
// 					pelaksanaResponses = append(pelaksanaResponses, pohonkinerja.PelaksanaOpdResponse{
// 						Id:          pelaksana.Id,
// 						PegawaiId:   pegawaiPelaksana.Id,
// 						NamaPegawai: pegawaiPelaksana.NamaPegawai,
// 					})
// 				}
// 				pelaksanaMap[p.Id] = pelaksanaResponses
// 			}

// 			// Ambil data indikator
// 			indikatorList, err := service.pohonKinerjaOpdRepository.FindIndikatorByPokinId(ctx, tx, fmt.Sprint(p.Id))
// 			if err == nil {
// 				var indikatorResponses []pohonkinerja.IndikatorResponse
// 				for _, indikator := range indikatorList {
// 					targetList, err := service.pohonKinerjaOpdRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
// 					if err != nil {
// 						continue
// 					}

// 					var targetResponses []pohonkinerja.TargetResponse
// 					for _, target := range targetList {
// 						targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 							Id:              target.Id,
// 							IndikatorId:     target.IndikatorId,
// 							TargetIndikator: target.Target,
// 							SatuanIndikator: target.Satuan,
// 						})
// 					}

// 					indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
// 						Id:            indikator.Id,
// 						IdPokin:       indikator.PokinId,
// 						NamaIndikator: indikator.Indikator,
// 						Target:        targetResponses,
// 					})
// 				}
// 				indikatorMap[p.Id] = indikatorResponses
// 			}
// 		}
// 	}

// 	// Proses khusus untuk level 4 (Strategic)
// 	if strategicList := pohonMap[4]; len(strategicList) > 0 {
// 		var allStrategics []domain.PohonKinerja
// 		processedIds := make(map[int]bool)

// 		// Kumpulkan semua strategic
// 		for _, strategicsByParent := range strategicList {
// 			for _, strategic := range strategicsByParent {
// 				if !processedIds[strategic.Id] {
// 					allStrategics = append(allStrategics, strategic)
// 					processedIds[strategic.Id] = true
// 				}
// 			}
// 		}

// 		// Urutkan strategic
// 		sort.Slice(allStrategics, func(i, j int) bool {
// 			if allStrategics[i].Status == "pokin dari pemda" && allStrategics[j].Status != "pokin dari pemda" {
// 				return true
// 			}
// 			if allStrategics[i].Status != "pokin dari pemda" && allStrategics[j].Status == "pokin dari pemda" {
// 				return false
// 			}
// 			return allStrategics[i].Id < allStrategics[j].Id
// 		})

// 		// Build response untuk setiap strategic
// 		for _, strategic := range allStrategics {
// 			strategicResp := service.buildStrategicResponse(ctx, tx, pohonMap, strategic, pelaksanaMap, indikatorMap)
// 			response.Strategics = append(response.Strategics, strategicResp)
// 		}
// 	}

// 	return response, nil
// }

// findall with redis
// func (service *PohonKinerjaOpdServiceImpl) FindAll(ctx context.Context, kodeOpd, tahun string) (pohonkinerja.PohonKinerjaOpdAllResponse, error) {
// 	// Log start time untuk menghitung response time
// 	startTime := time.Now()
// 	serviceName := "PohonKinerjaOpdService.FindAll"
// 	log.Printf("[%s] [START] [%s] kodeOpd=%s, tahun=%s",
// 		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, kodeOpd, tahun)

// 	// Generate cache key
// 	cacheKey := helper.GenerateCacheKey(helper.CacheKeyPohonKinerjaOpdAll, kodeOpd, tahun)

// 	// Coba ambil dari cache terlebih dahulu
// 	var response pohonkinerja.PohonKinerjaOpdAllResponse
// 	err := helper.GetFromCache(ctx, service.RedisClient, cacheKey, &response)
// 	if err == nil {
// 		// Cache hit, return data dari cache
// 		responseTime := time.Since(startTime)
// 		log.Printf("[%s] [CACHE HIT] [%s] kodeOpd=%s, tahun=%s, responseTime=%v",
// 			time.Now().Format("2006-01-02 15:04:05.000"), serviceName, kodeOpd, tahun, responseTime)
// 		return response, nil
// 	}

// 	// Cache miss, ambil data dari database
// 	log.Printf("[%s] [CACHE MISS] [%s] kodeOpd=%s, tahun=%s, mengambil dari database",
// 		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, kodeOpd, tahun)

// 	tx, err := service.DB.Begin()
// 	if err != nil {
// 		log.Printf("[%s] [ERROR] [%s] Error starting transaction: %v",
// 			time.Now().Format("2006-01-02 15:04:05.000"), serviceName, err)
// 		return pohonkinerja.PohonKinerjaOpdAllResponse{}, err
// 	}
// 	defer helper.CommitOrRollback(tx)

// 	// Validasi OPD
// 	opdStartTime := time.Now()
// 	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
// 	if err != nil {
// 		log.Printf("[%s] [ERROR] [%s] OPD not found: kode_opd=%s, error=%v, duration=%v",
// 			time.Now().Format("2006-01-02 15:04:05.000"), serviceName, kodeOpd, err, time.Since(opdStartTime))
// 		return pohonkinerja.PohonKinerjaOpdAllResponse{}, errors.New("kode opd tidak ditemukan")
// 	}
// 	log.Printf("[%s] [QUERY] [%s] FindByKodeOpd duration=%v",
// 		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, time.Since(opdStartTime))

// 	// Inisialisasi response dasar
// 	response = pohonkinerja.PohonKinerjaOpdAllResponse{
// 		KodeOpd:    kodeOpd,
// 		NamaOpd:    opd.NamaOpd,
// 		Tahun:      tahun,
// 		TujuanOpd:  make([]pohonkinerja.TujuanOpdResponse, 0),
// 		Strategics: make([]pohonkinerja.StrategicOpdResponse, 0),
// 	}

// 	// Ambil data tujuan OPD
// 	tujuanStartTime := time.Now()
// 	tujuanOpds, err := service.tujuanOpdRepository.FindTujuanOpdByTahun(ctx, tx, kodeOpd, tahun, "RPJMD")
// 	if err != nil {
// 		log.Printf("[%s] [WARNING] [%s] Error getting tujuan OPD: %v, duration=%v",
// 			time.Now().Format("2006-01-02 15:04:05.000"), serviceName, err, time.Since(tujuanStartTime))
// 		helper.SetToCache(ctx, service.RedisClient, cacheKey, response, helper.PohonKinerjaCacheTTL)
// 		return response, nil
// 	}
// 	log.Printf("[%s] [QUERY] [%s] FindTujuanOpdByTahun count=%d, duration=%v",
// 		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, len(tujuanOpds), time.Since(tujuanStartTime))

// 	// Proses data tujuan OPD (ini bisa dioptimasi lebih lanjut dengan batch queries)
// 	tujuanProcessStartTime := time.Now()
// 	for _, tujuan := range tujuanOpds {
// 		indikators, err := service.tujuanOpdRepository.FindIndikatorByTujuanOpdId(ctx, tx, tujuan.Id)
// 		if err != nil {
// 			log.Printf("[%s] [WARNING] [%s] Error getting indikator for tujuan ID %d: %v",
// 				time.Now().Format("2006-01-02 15:04:05.000"), serviceName, tujuan.Id, err)
// 			continue
// 		}

// 		var indikatorResponses []pohonkinerja.IndikatorTujuanResponse
// 		for _, indikator := range indikators {
// 			targets, err := service.tujuanOpdRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id, tahun)
// 			if err != nil {
// 				log.Printf("[%s] [WARNING] [%s] Error getting targets for indikator ID %s: %v",
// 					time.Now().Format("2006-01-02 15:04:05.000"), serviceName, indikator.Id, err)
// 				continue
// 			}

// 			var targetResponses []pohonkinerja.TargetTujuanResponse
// 			for _, target := range targets {
// 				targetResponses = append(targetResponses, pohonkinerja.TargetTujuanResponse{
// 					Tahun:  target.Tahun,
// 					Target: target.Target,
// 					Satuan: target.Satuan,
// 				})
// 			}

// 			indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorTujuanResponse{
// 				Indikator: indikator.Indikator,
// 				Target:    targetResponses,
// 			})
// 		}

// 		response.TujuanOpd = append(response.TujuanOpd, pohonkinerja.TujuanOpdResponse{
// 			Id:        tujuan.Id,
// 			KodeOpd:   tujuan.KodeOpd,
// 			Tujuan:    tujuan.Tujuan,
// 			Indikator: indikatorResponses,
// 		})
// 	}
// 	log.Printf("[%s] [PROCESS] [%s] Process tujuan OPD duration=%v",
// 		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, time.Since(tujuanProcessStartTime))

// 	// Ambil data pohon kinerja
// 	pokinStartTime := time.Now()
// 	pokins, err := service.pohonKinerjaOpdRepository.FindAll(ctx, tx, kodeOpd, tahun)
// 	if err != nil {
// 		log.Printf("[%s] [ERROR] [%s] Error getting pohon kinerja: %v, duration=%v",
// 			time.Now().Format("2006-01-02 15:04:05.000"), serviceName, err, time.Since(pokinStartTime))
// 		helper.SetToCache(ctx, service.RedisClient, cacheKey, response, helper.PohonKinerjaCacheTTL)
// 		return response, nil
// 	}
// 	log.Printf("[%s] [QUERY] [%s] FindAll pohon kinerja count=%d, duration=%v",
// 		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, len(pokins), time.Since(pokinStartTime))

// 	if len(pokins) == 0 {
// 		// Simpan ke cache meskipun hasil kosong
// 		helper.SetToCache(ctx, service.RedisClient, cacheKey, response, helper.PohonKinerjaCacheTTL)
// 		responseTime := time.Since(startTime)
// 		log.Printf("[%s] [END] [%s] kodeOpd=%s, tahun=%s, totalResponseTime=%v (empty result)",
// 			time.Now().Format("2006-01-02 15:04:05.000"), serviceName, kodeOpd, tahun, responseTime)
// 		return response, nil
// 	}

// 	// OPTIMASI: Proses data pohon kinerja dengan batch queries
// 	pohonMap := make(map[int]map[int][]domain.PohonKinerja)
// 	indikatorMap := make(map[int][]pohonkinerja.IndikatorResponse)

// 	// Kumpulkan semua pokin IDs untuk batch queries
// 	var pokinIds []int
// 	maxLevel := 0
// 	for _, p := range pokins {
// 		if p.LevelPohon >= 4 {
// 			if p.LevelPohon > maxLevel {
// 				maxLevel = p.LevelPohon
// 			}

// 			if pohonMap[p.LevelPohon] == nil {
// 				pohonMap[p.LevelPohon] = make(map[int][]domain.PohonKinerja)
// 			}

// 			p.NamaOpd = opd.NamaOpd
// 			pohonMap[p.LevelPohon][p.Parent] = append(
// 				pohonMap[p.LevelPohon][p.Parent],
// 				p,
// 			)
// 			pokinIds = append(pokinIds, p.Id)
// 		}
// 	}

// 	log.Printf("[%s] [INFO] [%s] Processing %d pokin entries, maxLevel=%d",
// 		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, len(pokinIds), maxLevel)

// 	// OPTIMASI: Batch fetch semua pelaksana dalam 1 query (sudah include nama pegawai dari JOIN)
// 	pelaksanaMap := make(map[int][]pohonkinerja.PelaksanaOpdResponse)
// 	if len(pokinIds) > 0 {
// 		pelaksanaBatchStartTime := time.Now()
// 		pelaksanaBatch, err := service.pohonKinerjaOpdRepository.FindPelaksanaPokinBatch(ctx, tx, pokinIds)
// 		if err == nil {
// 			// FindPelaksanaPokinBatch sudah JOIN dengan tb_pegawai, jadi langsung gunakan nama_pegawai yang sudah ada
// 			// Tidak perlu fetch pegawai lagi secara individual!
// 			for pokinId, pelaksanaList := range pelaksanaBatch {
// 				var pelaksanaResponses []pohonkinerja.PelaksanaOpdResponse
// 				for _, pelaksana := range pelaksanaList {
// 					pelaksanaResponses = append(pelaksanaResponses, pohonkinerja.PelaksanaOpdResponse{
// 						Id:          pelaksana.Id,
// 						PegawaiId:   pelaksana.PegawaiId,
// 						NamaPegawai: pelaksana.NamaPegawai, // Sudah ada dari JOIN di repository
// 					})
// 				}
// 				pelaksanaMap[pokinId] = pelaksanaResponses
// 			}
// 			log.Printf("[%s] [QUERY] [%s] FindPelaksanaPokinBatch count=%d, duration=%v",
// 				time.Now().Format("2006-01-02 15:04:05.000"), serviceName, len(pelaksanaBatch), time.Since(pelaksanaBatchStartTime))
// 		} else {
// 			log.Printf("[%s] [ERROR] [%s] Error batch fetching pelaksana: %v",
// 				time.Now().Format("2006-01-02 15:04:05.000"), serviceName, err)
// 		}
// 	}

// 	// OPTIMASI: Batch fetch semua indikator dan target dalam 1 query
// 	if len(pokinIds) > 0 {
// 		indikatorBatchStartTime := time.Now()
// 		indikatorBatch, err := service.pohonKinerjaOpdRepository.FindIndikatorByPokinIdsBatch(ctx, tx, pokinIds)
// 		if err == nil {
// 			// Kumpulkan semua indikator IDs untuk batch fetch target
// 			var allIndikatorIds []string
// 			for _, indikatorList := range indikatorBatch {
// 				for _, indikator := range indikatorList {
// 					allIndikatorIds = append(allIndikatorIds, indikator.Id)
// 				}
// 			}

// 			// Batch fetch semua target dalam 1 query
// 			var targetBatch map[string][]domain.Target
// 			targetBatchStartTime := time.Now()
// 			if len(allIndikatorIds) > 0 {
// 				targetBatch, err = service.pohonKinerjaOpdRepository.FindTargetByIndikatorIdsBatch(ctx, tx, allIndikatorIds)
// 				if err != nil {
// 					log.Printf("[%s] [ERROR] [%s] Error batch fetching targets: %v",
// 						time.Now().Format("2006-01-02 15:04:05.000"), serviceName, err)
// 					targetBatch = make(map[string][]domain.Target)
// 				} else {
// 					log.Printf("[%s] [QUERY] [%s] FindTargetByIndikatorIdsBatch count=%d, duration=%v",
// 						time.Now().Format("2006-01-02 15:04:05.000"), serviceName, len(allIndikatorIds), time.Since(targetBatchStartTime))
// 				}
// 			} else {
// 				targetBatch = make(map[string][]domain.Target)
// 			}

// 			// Build indikator responses dengan target
// 			for pokinId, indikatorList := range indikatorBatch {
// 				var indikatorResponses []pohonkinerja.IndikatorResponse
// 				for _, indikator := range indikatorList {
// 					targetList := targetBatch[indikator.Id]
// 					var targetResponses []pohonkinerja.TargetResponse
// 					for _, target := range targetList {
// 						targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 							Id:              target.Id,
// 							IndikatorId:     target.IndikatorId,
// 							TargetIndikator: target.Target,
// 							SatuanIndikator: target.Satuan,
// 						})
// 					}

// 					indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
// 						Id:            indikator.Id,
// 						IdPokin:       fmt.Sprint(pokinId),
// 						NamaIndikator: indikator.Indikator,
// 						Target:        targetResponses,
// 					})
// 				}
// 				indikatorMap[pokinId] = indikatorResponses
// 			}
// 			log.Printf("[%s] [QUERY] [%s] FindIndikatorByPokinIdsBatch count=%d, duration=%v",
// 				time.Now().Format("2006-01-02 15:04:05.000"), serviceName, len(indikatorBatch), time.Since(indikatorBatchStartTime))
// 		} else {
// 			log.Printf("[%s] [ERROR] [%s] Error batch fetching indikator: %v",
// 				time.Now().Format("2006-01-02 15:04:05.000"), serviceName, err)
// 		}
// 	}

// 	// OPTIMASI: Pre-fetch semua data yang dibutuhkan untuk build response
// 	prefetchStartTime := time.Now()

// 	// Pre-fetch semua tagging untuk semua pokin
// 	taggingMap := make(map[int][]domain.TaggingPokin)
// 	for _, pokinId := range pokinIds {
// 		taggingList, err := service.pohonKinerjaOpdRepository.FindTaggingByPokinId(ctx, tx, pokinId)
// 		if err == nil {
// 			taggingMap[pokinId] = taggingList
// 		}
// 	}

// 	// Pre-fetch semua review untuk semua pokin
// 	reviewMap := make(map[int][]domain.Review)
// 	reviewCountMap := make(map[int]int)
// 	reviewNipSet := make(map[string]bool)
// 	for _, pokinId := range pokinIds {
// 		count, err := service.reviewRepository.CountReviewByPohonKinerja(ctx, tx, pokinId)
// 		if err == nil {
// 			reviewCountMap[pokinId] = count
// 		}

// 		reviews, err := service.reviewRepository.FindByPohonKinerja(ctx, tx, pokinId)
// 		if err == nil {
// 			reviewMap[pokinId] = reviews
// 			for _, review := range reviews {
// 				if review.CreatedBy != "" {
// 					reviewNipSet[review.CreatedBy] = true
// 				}
// 			}
// 		}
// 	}

// 	// Batch fetch semua pegawai untuk review (hanya yang belum ada di pelaksana)
// 	var reviewNips []string
// 	for nip := range reviewNipSet {
// 		reviewNips = append(reviewNips, nip)
// 	}
// 	reviewPegawaiMap := make(map[string]*domainmaster.Pegawai)
// 	for _, nip := range reviewNips {
// 		pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, nip)
// 		if err == nil {
// 			reviewPegawaiMap[nip] = &pegawai
// 		}
// 	}

// 	// Pre-fetch semua program unggulan yang dibutuhkan
// 	programUnggulanMap := make(map[string]*domain.ProgramUnggulan)
// 	for _, taggingList := range taggingMap {
// 		for _, tagging := range taggingList {
// 			for _, keterangan := range tagging.KeteranganTaggingProgram {
// 				if _, exists := programUnggulanMap[keterangan.KodeProgramUnggulan]; !exists {
// 					program, err := service.ProgramUnggulanRepository.FindByKodeProgramUnggulan(ctx, tx, keterangan.KodeProgramUnggulan)
// 					if err == nil {
// 						programUnggulanMap[keterangan.KodeProgramUnggulan] = &program
// 					}
// 				}
// 			}
// 		}
// 	}

// 	log.Printf("[%s] [QUERY] [%s] Pre-fetch tagging/review/program duration=%v",
// 		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, time.Since(prefetchStartTime))

// 	// Proses khusus untuk level 4 (Strategic)
// 	buildResponseStartTime := time.Now()
// 	if strategicList := pohonMap[4]; len(strategicList) > 0 {
// 		var allStrategics []domain.PohonKinerja
// 		processedIds := make(map[int]bool)

// 		// Kumpulkan semua strategic
// 		for _, strategicsByParent := range strategicList {
// 			for _, strategic := range strategicsByParent {
// 				if !processedIds[strategic.Id] {
// 					allStrategics = append(allStrategics, strategic)
// 					processedIds[strategic.Id] = true
// 				}
// 			}
// 		}

// 		// Urutkan strategic
// 		sort.Slice(allStrategics, func(i, j int) bool {
// 			if allStrategics[i].Status == "pokin dari pemda" && allStrategics[j].Status != "pokin dari pemda" {
// 				return true
// 			}
// 			if allStrategics[i].Status != "pokin dari pemda" && allStrategics[j].Status == "pokin dari pemda" {
// 				return false
// 			}
// 			return allStrategics[i].Id < allStrategics[j].Id
// 		})

// 		// Build response untuk setiap strategic menggunakan pre-fetched data
// 		for _, strategic := range allStrategics {
// 			strategicResp := service.buildStrategicResponseOptimized(
// 				ctx, tx, pohonMap, strategic, pelaksanaMap, indikatorMap,
// 				taggingMap[strategic.Id], reviewMap[strategic.Id], reviewCountMap[strategic.Id],
// 				reviewPegawaiMap, programUnggulanMap,
// 			)
// 			response.Strategics = append(response.Strategics, strategicResp)
// 		}
// 	}
// 	log.Printf("[%s] [PROCESS] [%s] Build strategic responses duration=%v",
// 		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, time.Since(buildResponseStartTime))

// 	// Simpan ke cache setelah berhasil mengambil data
// 	cacheSetStartTime := time.Now()
// 	helper.SetToCache(ctx, service.RedisClient, cacheKey, response, helper.PohonKinerjaCacheTTL)
// 	log.Printf("[%s] [CACHE] [%s] Set cache duration=%v",
// 		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, time.Since(cacheSetStartTime))

// 	responseTime := time.Since(startTime)
// 	log.Printf("[%s] [END] [%s] kodeOpd=%s, tahun=%s, totalResponseTime=%v, strategicsCount=%d",
// 		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, kodeOpd, tahun, responseTime, len(response.Strategics))

//		return response, nil
//	}
//
// CEK OPTIMASI
func (service *PohonKinerjaOpdServiceImpl) FindAll(ctx context.Context, kodeOpd, tahun string) (pohonkinerja.PohonKinerjaOpdAllResponse, error) {
	// Log start time
	startTime := time.Now()
	serviceName := "PohonKinerjaOpdService.FindAll"
	log.Printf("[%s] [START] [%s] kodeOpd=%s, tahun=%s",
		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, kodeOpd, tahun)

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("[%s] [ERROR] [%s] Error starting transaction: %v",
			time.Now().Format("2006-01-02 15:04:05.000"), serviceName, err)
		return pohonkinerja.PohonKinerjaOpdAllResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi OPD
	opdStartTime := time.Now()
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		log.Printf("[%s] [ERROR] [%s] OPD not found: kode_opd=%s, error=%v, duration=%v",
			time.Now().Format("2006-01-02 15:04:05.000"), serviceName, kodeOpd, err, time.Since(opdStartTime))
		return pohonkinerja.PohonKinerjaOpdAllResponse{}, errors.New("kode opd tidak ditemukan")
	}
	log.Printf("[%s] [QUERY] [%s] FindByKodeOpd duration=%v",
		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, time.Since(opdStartTime))

	// Inisialisasi response dasar
	response := pohonkinerja.PohonKinerjaOpdAllResponse{
		KodeOpd:    kodeOpd,
		NamaOpd:    opd.NamaOpd,
		Tahun:      tahun,
		TujuanOpd:  make([]pohonkinerja.TujuanOpdResponse, 0),
		Strategics: make([]pohonkinerja.StrategicOpdResponse, 0),
	}

	// Ambil data tujuan OPD
	tujuanStartTime := time.Now()
	tujuanOpds, err := service.tujuanOpdRepository.FindTujuanOpdByTahun(ctx, tx, kodeOpd, tahun, "RPJMD")
	if err != nil {
		log.Printf("[%s] [WARNING] [%s] Error getting tujuan OPD: %v, duration=%v",
			time.Now().Format("2006-01-02 15:04:05.000"), serviceName, err, time.Since(tujuanStartTime))
	} else {
		// Proses data tujuan OPD dengan batch
		tujuanProcessStartTime := time.Now()
		if len(tujuanOpds) > 0 {
			var tujuanIds []int
			for _, tujuan := range tujuanOpds {
				tujuanIds = append(tujuanIds, tujuan.Id)
			}

			// Batch fetch semua indikator untuk semua tujuan
			indikatorTujuanMap := make(map[int][]domain.Indikator)
			if len(tujuanIds) > 0 {
				indikatorBatch, err := service.tujuanOpdRepository.FindIndikatorByTujuanOpdIdsBatch(ctx, tx, tujuanIds)
				if err == nil {
					indikatorTujuanMap = indikatorBatch
				}
			}

			// Kumpulkan semua indikator IDs untuk batch fetch target
			var allTujuanIndikatorIds []string
			for _, indikators := range indikatorTujuanMap {
				for _, indikator := range indikators {
					allTujuanIndikatorIds = append(allTujuanIndikatorIds, indikator.Id)
				}
			}

			// Batch fetch semua target untuk semua indikator tujuan
			targetTujuanBatch := make(map[string][]domain.Target)
			if len(allTujuanIndikatorIds) > 0 {
				targetBatch, err := service.pohonKinerjaOpdRepository.FindTargetByIndikatorIdsBatch(ctx, tx, allTujuanIndikatorIds)
				if err == nil {
					// Filter target berdasarkan tahun
					for indikatorId, targets := range targetBatch {
						filteredTargets := make([]domain.Target, 0, len(targets))
						for _, target := range targets {
							if target.Tahun == tahun {
								filteredTargets = append(filteredTargets, target)
							}
						}
						if len(filteredTargets) > 0 {
							targetTujuanBatch[indikatorId] = filteredTargets
						}
					}
				}
			}

			// Build response tujuan OPD
			tujuanResponses := make([]pohonkinerja.TujuanOpdResponse, 0, len(tujuanOpds))
			for _, tujuan := range tujuanOpds {
				indikators := indikatorTujuanMap[tujuan.Id]
				indikatorResponses := make([]pohonkinerja.IndikatorTujuanResponse, 0, len(indikators))
				for _, indikator := range indikators {
					targets := targetTujuanBatch[indikator.Id]
					targetResponses := make([]pohonkinerja.TargetTujuanResponse, 0, len(targets))
					for _, target := range targets {
						targetResponses = append(targetResponses, pohonkinerja.TargetTujuanResponse{
							Tahun:  target.Tahun,
							Target: target.Target,
							Satuan: target.Satuan,
						})
					}
					indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorTujuanResponse{
						Indikator: indikator.Indikator,
						Target:    targetResponses,
					})
				}
				tujuanResponses = append(tujuanResponses, pohonkinerja.TujuanOpdResponse{
					Id:        tujuan.Id,
					KodeOpd:   tujuan.KodeOpd,
					Tujuan:    tujuan.Tujuan,
					Indikator: indikatorResponses,
				})
			}
			response.TujuanOpd = tujuanResponses
		}
		log.Printf("[%s] [PROCESS] [%s] Process tujuan OPD duration=%v",
			time.Now().Format("2006-01-02 15:04:05.000"), serviceName, time.Since(tujuanProcessStartTime))
	}
	log.Printf("[%s] [QUERY] [%s] FindTujuanOpdByTahun count=%d, duration=%v",
		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, len(tujuanOpds), time.Since(tujuanStartTime))

	// Ambil data pohon kinerja
	pokinStartTime := time.Now()
	pokins, err := service.pohonKinerjaOpdRepository.FindAll(ctx, tx, kodeOpd, tahun)
	if err != nil {
		log.Printf("[%s] [ERROR] [%s] Error getting pohon kinerja: %v, duration=%v",
			time.Now().Format("2006-01-02 15:04:05.000"), serviceName, err, time.Since(pokinStartTime))
		return response, nil
	}
	log.Printf("[%s] [QUERY] [%s] FindAll pohon kinerja count=%d, duration=%v",
		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, len(pokins), time.Since(pokinStartTime))

	if len(pokins) == 0 {
		log.Printf("[%s] [INFO] [%s] pohon_kinerja kosong",
			time.Now().Format("2006-01-02 15:04:05.000"), serviceName)
		responseTime := time.Since(startTime)
		log.Printf("[%s] [END] [%s] kodeOpd=%s, tahun=%s, totalResponseTime=%v (empty result)",
			time.Now().Format("2006-01-02 15:04:05.000"), serviceName, kodeOpd, tahun, responseTime)
		return response, nil
	}

	// Kumpulkan semua pokin IDs untuk batch queries
	pokinIds := make([]int, 0, len(pokins))
	for _, p := range pokins {
		if p.LevelPohon >= 4 {
			pokinIds = append(pokinIds, p.Id)
		}
	}

	// Batch fetch pelaksana
	pelaksanaBatchStartTime := time.Now()
	pelaksanas, err := service.pohonKinerjaOpdRepository.FindPelaksanaPokinBatch(ctx, tx, pokinIds)
	if err != nil {
		log.Printf("[%s] [ERROR] [%s] FindPelaksanaPokinBatch gagal | err=%v",
			time.Now().Format("2006-01-02 15:04:05.000"), serviceName, err)
		return response, err
	}
	pelaksanaMap := make(map[int][]pohonkinerja.PelaksanaOpdResponse)
	for pokinId, pelaksanaList := range pelaksanas {
		pelaksanaResponses := make([]pohonkinerja.PelaksanaOpdResponse, 0, len(pelaksanaList))
		for _, p := range pelaksanaList {
			pelaksanaResponses = append(pelaksanaResponses, pohonkinerja.PelaksanaOpdResponse{
				Id:          p.Id,
				PegawaiId:   p.PegawaiId,
				NamaPegawai: p.NamaPegawai,
			})
		}
		pelaksanaMap[pokinId] = pelaksanaResponses
	}
	log.Printf("[%s] [QUERY] [%s] FindPelaksanaPokinBatch count=%d, duration=%v",
		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, len(pelaksanas), time.Since(pelaksanaBatchStartTime))

	// Batch fetch indikator
	indikatorBatchStartTime := time.Now()
	indikatorBatch, err := service.pohonKinerjaOpdRepository.FindIndikatorByPokinIdsBatch(ctx, tx, pokinIds)
	if err != nil {
		log.Printf("[%s] [ERROR] [%s] FindIndikatorByPokinIdsBatch gagal | err=%v",
			time.Now().Format("2006-01-02 15:04:05.000"), serviceName, err)
		return response, err
	}

	// Kumpulkan semua indikator IDs untuk batch fetch target
	var allIndikatorIds []string
	for _, indikatorList := range indikatorBatch {
		for _, indikator := range indikatorList {
			allIndikatorIds = append(allIndikatorIds, indikator.Id)
		}
	}

	// Batch fetch semua target untuk semua indikator
	var targetBatch map[string][]domain.Target
	if len(allIndikatorIds) > 0 {
		targetBatch, err = service.pohonKinerjaOpdRepository.FindTargetByIndikatorIdsBatch(ctx, tx, allIndikatorIds)
		if err != nil {
			log.Printf("[%s] [ERROR] [%s] FindTargetByIndikatorIdsBatch gagal | err=%v",
				time.Now().Format("2006-01-02 15:04:05.000"), serviceName, err)
			targetBatch = make(map[string][]domain.Target)
		}
	} else {
		targetBatch = make(map[string][]domain.Target)
	}

	// Build indikator map dengan target
	indikatorMap := make(map[int][]pohonkinerja.IndikatorResponse)
	for pokinId, indikatorList := range indikatorBatch {
		indikatorResponses := make([]pohonkinerja.IndikatorResponse, 0, len(indikatorList))
		for _, indikator := range indikatorList {
			// Ambil target dari batch yang sudah di-fetch
			targetList := targetBatch[indikator.Id]
			targetResponses := make([]pohonkinerja.TargetResponse, 0, len(targetList))
			for _, target := range targetList {
				targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
					Id:              target.Id,
					IndikatorId:     target.IndikatorId,
					TargetIndikator: target.Target,
					SatuanIndikator: target.Satuan,
					// TahunSasaran:    target.Tahun,
				})
			}

			indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
				Id:            indikator.Id,
				IdPokin:       fmt.Sprint(pokinId), // Konversi int pokinId ke string
				NamaIndikator: indikator.Indikator,
				Target:        targetResponses,
			})
		}
		indikatorMap[pokinId] = indikatorResponses
	}
	log.Printf("[%s] [QUERY] [%s] FindIndikatorByPokinIdsBatch count=%d, duration=%v",
		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, len(indikatorBatch), time.Since(indikatorBatchStartTime))

	// Batch fetch tematik berdasarkan clone_from
	// Batch fetch tematik berdasarkan clone_from - OPTIMASI: Skip jika tidak ada atau buat optional
	tematikBatchStartTime := time.Now()
	tematikCloneFromSet := make(map[int]bool)
	for _, p := range pokins {
		if p.LevelPohon >= 4 && p.CloneFrom > 0 {
			tematikCloneFromSet[p.CloneFrom] = true
		}
	}

	// OPTIMASI: Batasi jumlah cloneFromIds untuk menghindari query yang terlalu kompleks
	var tematikCloneFromIds []int
	maxTematikIds := 100 // Batasi maksimal 100 IDs untuk menghindari query yang terlalu lambat
	count := 0
	for cloneFromId := range tematikCloneFromSet {
		if count >= maxTematikIds {
			break
		}
		tematikCloneFromIds = append(tematikCloneFromIds, cloneFromId)
		count++
	}

	tematikMap := make(map[int]*domain.PohonKinerja)
	if len(tematikCloneFromIds) > 0 {
		// OPTIMASI: Gunakan context dengan timeout untuk menghindari query yang terlalu lama
		tematikCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		tematikBatch, err := service.pohonKinerjaOpdRepository.FindTematikByCloneFromBatch(tematikCtx, tx, tematikCloneFromIds)
		if err == nil {
			tematikMap = tematikBatch
		} else {
			// Log error tapi jangan fail seluruh request
			log.Printf("[%s] [WARNING] [%s] FindTematikByCloneFromBatch error (non-fatal): %v",
				time.Now().Format("2006-01-02 15:04:05.000"), serviceName, err)
		}
	}
	log.Printf("[%s] [QUERY] [%s] FindTematikByCloneFromBatch count=%d, duration=%v",
		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, len(tematikMap), time.Since(tematikBatchStartTime))

	// Proses data pohon kinerja ke map
	pohonMap := make(map[int]map[int][]domain.PohonKinerja)
	for _, p := range pokins {
		if p.LevelPohon >= 4 {
			if pohonMap[p.LevelPohon] == nil {
				pohonMap[p.LevelPohon] = make(map[int][]domain.PohonKinerja)
			}
			p.NamaOpd = opd.NamaOpd
			pohonMap[p.LevelPohon][p.Parent] = append(pohonMap[p.LevelPohon][p.Parent], p)
		}
	}

	// Batch fetch tagging
	taggingBatchStartTime := time.Now()
	taggings, err := service.pohonKinerjaOpdRepository.FindTaggingByPokinIdsBatch(ctx, tx, pokinIds)
	if err != nil {
		log.Printf("[%s] [ERROR] [%s] FindTaggingByPokinIdsBatch gagal | err=%v",
			time.Now().Format("2006-01-02 15:04:05.000"), serviceName, err)
		return response, err
	}
	taggingMap := make(map[int][]pohonkinerja.TaggingResponse)
	for pokinId, tagList := range taggings {
		taggingResponses := make([]pohonkinerja.TaggingResponse, 0, len(tagList))
		for _, tag := range tagList {
			keteranganResponses := make([]pohonkinerja.KeteranganTaggingResponse, 0, len(tag.KeteranganTaggingProgram))
			for _, keterangan := range tag.KeteranganTaggingProgram {
				keteranganResponses = append(keteranganResponses, pohonkinerja.KeteranganTaggingResponse{
					Id:                  keterangan.Id,
					IdTagging:           keterangan.IdTagging,
					KodeProgramUnggulan: keterangan.KodeProgramUnggulan,
					RencanaImplementasi: keterangan.RencanaImplementasi,
				})
			}
			taggingResponses = append(taggingResponses, pohonkinerja.TaggingResponse{
				Id:                       tag.Id,
				IdPokin:                  tag.IdPokin,
				NamaTagging:              tag.NamaTagging,
				KeteranganTaggingProgram: keteranganResponses,
			})
		}
		taggingMap[pokinId] = taggingResponses
	}
	log.Printf("[%s] [QUERY] [%s] FindTaggingByPokinIdsBatch count=%d, duration=%v",
		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, len(taggings), time.Since(taggingBatchStartTime))

	// Batch fetch review
	reviewBatchStartTime := time.Now()
	reviews, err := service.reviewRepository.FindByPokinIdBatch(ctx, tx, pokinIds)
	if err != nil {
		log.Printf("[%s] [ERROR] [%s] FindByPokinIdBatch gagal | err=%v",
			time.Now().Format("2006-01-02 15:04:05.000"), serviceName, err)
		return response, err
	}

	// OPTIMASI: Pre-count reviews per pokinId untuk pre-allocation
	reviewCountPerPokin := make(map[int]int)
	for _, review := range reviews {
		reviewCountPerPokin[review.IdPohonKinerja]++
	}

	// Pre-allocate reviewMap dengan capacity yang tepat
	reviewMap := make(map[int][]pohonkinerja.ReviewResponse, len(reviewCountPerPokin))
	for pokinId, count := range reviewCountPerPokin {
		reviewMap[pokinId] = make([]pohonkinerja.ReviewResponse, 0, count)
	}

	// Fill reviewMap
	for _, review := range reviews {
		reviewMap[review.IdPohonKinerja] = append(reviewMap[review.IdPohonKinerja], pohonkinerja.ReviewResponse{
			Id:             review.Id,
			IdPohonKinerja: review.IdPohonKinerja,
			Review:         review.Review,
			Keterangan:     review.Keterangan,
			CreatedBy:      review.CreatedBy,
			NamaPegawai:    review.NamaReviewer,
			JenisPokin:     review.Jenis_pokin,
		})
	}
	log.Printf("[%s] [QUERY] [%s] FindByPokinIdBatch count=%d, duration=%v",
		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, len(reviews), time.Since(reviewBatchStartTime))

	// Build response untuk strategic (level 4)
	buildResponseStartTime := time.Now()
	strategicList := pohonMap[4]
	if len(strategicList) > 0 {
		allStrategics := flattenAndSort(strategicList)

		for _, strategic := range allStrategics {
			strategicResp := buildStrategicOnly(
				strategic,
				taggingMap,
				pelaksanaMap,
				indikatorMap,
				reviewMap,
				tematikMap,
			)

			// Append tactical (level 5)
			if tacticalsByParent, ok := pohonMap[5][strategic.Id]; ok {
				for _, tactical := range tacticalsByParent {
					tacticalResp := buildTacticalOnly(
						tactical,
						taggingMap,
						pelaksanaMap,
						indikatorMap,
						reviewMap,
						tematikMap,
					)

					// Lanjut append operational
					appendOperationals(&tacticalResp, pohonMap, taggingMap, pelaksanaMap, indikatorMap, reviewMap, tematikMap)

					strategicResp.Tacticals = append(strategicResp.Tacticals, tacticalResp)
				}
			}

			response.Strategics = append(response.Strategics, strategicResp)
		}
	}
	log.Printf("[%s] [PROCESS] [%s] Build strategic responses duration=%v",
		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, time.Since(buildResponseStartTime))

	responseTime := time.Since(startTime)
	log.Printf("[%s] [END] [%s] kodeOpd=%s, tahun=%s, totalResponseTime=%v, strategicsCount=%d",
		time.Now().Format("2006-01-02 15:04:05.000"), serviceName, kodeOpd, tahun, responseTime, len(response.Strategics))

	return response, nil
}

// Helper function untuk flatten dan sort strategic
func flattenAndSort(nodesByParent map[int][]domain.PohonKinerja) []domain.PohonKinerja {
	var result []domain.PohonKinerja
	seen := make(map[int]bool)

	for _, nodes := range nodesByParent {
		for _, n := range nodes {
			if !seen[n.Id] {
				result = append(result, n)
				seen[n.Id] = true
			}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Status == "pokin dari pemda" && result[j].Status != "pokin dari pemda" {
			return true
		}
		if result[i].Status != "pokin dari pemda" && result[j].Status == "pokin dari pemda" {
			return false
		}
		return result[i].Id < result[j].Id
	})

	return result
}

// Helper function untuk build strategic only
func buildStrategicOnly(
	strategic domain.PohonKinerja,
	taggingMap map[int][]pohonkinerja.TaggingResponse,
	pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse,
	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
	reviewMap map[int][]pohonkinerja.ReviewResponse,
	tematikMap map[int]*domain.PohonKinerja,
) pohonkinerja.StrategicOpdResponse {
	var keteranganCrosscutting *string
	if strategic.KeteranganCrosscutting != nil && *strategic.KeteranganCrosscutting != "" {
		keteranganCrosscutting = strategic.KeteranganCrosscutting
	}
	// Cari tematik dari pre-fetched map
	var idTematik *int
	var namaTematik *string
	if strategic.CloneFrom > 0 {
		if tematik, exists := tematikMap[strategic.CloneFrom]; exists && tematik != nil {
			idTematik = &tematik.Id
			namaTematik = &tematik.NamaPohon
		} else {
			// Jika tematik tidak ditemukan tapi CloneFrom > 0, tetap set idTematik
			idTematik = &strategic.CloneFrom
		}
	}

	reviewPokin := reviewMap[strategic.Id]
	countReview := len(reviewPokin)

	strategicResp := pohonkinerja.StrategicOpdResponse{
		Id:                     strategic.Id,
		Parent:                 nil,
		Strategi:               strategic.NamaPohon,
		JenisPohon:             strategic.JenisPohon,
		LevelPohon:             strategic.LevelPohon,
		Keterangan:             strategic.Keterangan,
		KeteranganCrosscutting: keteranganCrosscutting,
		Status:                 strategic.Status,
		IsActive:               strategic.IsActive,
		IdTematik:              idTematik,
		NamaTematik:            namaTematik,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: strategic.KodeOpd,
			NamaOpd: strategic.NamaOpd,
		},
		Tagging:     taggingMap[strategic.Id],
		Pelaksana:   pelaksanaMap[strategic.Id],
		Indikator:   indikatorMap[strategic.Id],
		Review:      reviewPokin,
		CountReview: countReview,
	}
	return strategicResp
}

// Helper function untuk build tactical only
func buildTacticalOnly(
	tactical domain.PohonKinerja,
	taggingMap map[int][]pohonkinerja.TaggingResponse,
	pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse,
	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
	reviewMap map[int][]pohonkinerja.ReviewResponse,
	tematikMap map[int]*domain.PohonKinerja,
) pohonkinerja.TacticalOpdResponse {
	var keteranganCrosscutting *string
	if tactical.KeteranganCrosscutting != nil && *tactical.KeteranganCrosscutting != "" {
		keteranganCrosscutting = tactical.KeteranganCrosscutting
	}
	// Cari tematik dari pre-fetched map
	var idTematik *int
	var namaTematik *string
	if tactical.CloneFrom > 0 {
		if tematik, exists := tematikMap[tactical.CloneFrom]; exists && tematik != nil {
			idTematik = &tematik.Id
			namaTematik = &tematik.NamaPohon
		} else {
			// Jika tematik tidak ditemukan tapi CloneFrom > 0, tetap set idTematik
			idTematik = &tactical.CloneFrom
		}
	}

	reviewPokin := reviewMap[tactical.Id]
	countReview := len(reviewPokin)

	tacticalResp := pohonkinerja.TacticalOpdResponse{
		Id:                     tactical.Id,
		Parent:                 tactical.Parent,
		Strategi:               tactical.NamaPohon,
		JenisPohon:             tactical.JenisPohon,
		LevelPohon:             tactical.LevelPohon,
		Keterangan:             tactical.Keterangan,
		KeteranganCrosscutting: keteranganCrosscutting,
		Status:                 tactical.Status,
		IsActive:               tactical.IsActive,
		IdTematik:              idTematik,
		NamaTematik:            namaTematik,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: tactical.KodeOpd,
			NamaOpd: tactical.NamaOpd,
		},
		Tagging:     taggingMap[tactical.Id],
		Pelaksana:   pelaksanaMap[tactical.Id],
		Indikator:   indikatorMap[tactical.Id],
		Review:      reviewPokin,
		CountReview: countReview,
	}
	return tacticalResp
}

// Helper function untuk build operational only
func buildOperationalOnly(
	operational domain.PohonKinerja,
	taggingMap map[int][]pohonkinerja.TaggingResponse,
	pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse,
	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
	reviewMap map[int][]pohonkinerja.ReviewResponse,
	tematikMap map[int]*domain.PohonKinerja,
) pohonkinerja.OperationalOpdResponse {
	var keteranganCrosscutting *string
	if operational.KeteranganCrosscutting != nil && *operational.KeteranganCrosscutting != "" {
		keteranganCrosscutting = operational.KeteranganCrosscutting
	}
	// Cari tematik dari pre-fetched map
	// Cari tematik dari pre-fetched map
	var idTematik *int
	var namaTematik *string
	if operational.CloneFrom > 0 {
		if tematik, exists := tematikMap[operational.CloneFrom]; exists && tematik != nil {
			idTematik = &tematik.Id
			namaTematik = &tematik.NamaPohon
		} else {
			// Jika tematik tidak ditemukan tapi CloneFrom > 0, tetap set idTematik
			idTematik = &operational.CloneFrom
		}
	}

	reviewPokin := reviewMap[operational.Id]
	countReview := len(reviewPokin)

	operationalResp := pohonkinerja.OperationalOpdResponse{
		Id:                     operational.Id,
		Parent:                 operational.Parent,
		Strategi:               operational.NamaPohon,
		JenisPohon:             operational.JenisPohon,
		LevelPohon:             operational.LevelPohon,
		Keterangan:             operational.Keterangan,
		KeteranganCrosscutting: keteranganCrosscutting,
		Status:                 operational.Status,
		IsActive:               operational.IsActive,
		IdTematik:              idTematik,
		NamaTematik:            namaTematik,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: operational.KodeOpd,
			NamaOpd: operational.NamaOpd,
		},
		Tagging:     taggingMap[operational.Id],
		Pelaksana:   pelaksanaMap[operational.Id],
		Indikator:   indikatorMap[operational.Id],
		Review:      reviewPokin,
		CountReview: countReview,
	}
	return operationalResp
}

// Helper function untuk build operational N only
func buildOperationalNOnly(
	operationalN domain.PohonKinerja,
	taggingMap map[int][]pohonkinerja.TaggingResponse,
	pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse,
	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
	reviewMap map[int][]pohonkinerja.ReviewResponse,
) pohonkinerja.OperationalNOpdResponse {
	var keteranganCrosscutting *string
	if operationalN.KeteranganCrosscutting != nil && *operationalN.KeteranganCrosscutting != "" {
		keteranganCrosscutting = operationalN.KeteranganCrosscutting
	}

	reviewPokin := reviewMap[operationalN.Id]
	countReview := len(reviewPokin)

	operationalNResp := pohonkinerja.OperationalNOpdResponse{
		Id:                     operationalN.Id,
		Parent:                 operationalN.Parent,
		Strategi:               operationalN.NamaPohon,
		JenisPohon:             operationalN.JenisPohon,
		LevelPohon:             operationalN.LevelPohon,
		Keterangan:             operationalN.Keterangan,
		KeteranganCrosscutting: keteranganCrosscutting,
		Status:                 operationalN.Status,
		IsActive:               operationalN.IsActive,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: operationalN.KodeOpd,
			NamaOpd: operationalN.NamaOpd,
		},
		Tagging:     taggingMap[operationalN.Id],
		Pelaksana:   pelaksanaMap[operationalN.Id],
		Indikator:   indikatorMap[operationalN.Id],
		Review:      reviewPokin,
		CountReview: countReview,
	}
	return operationalNResp
}

// Helper function untuk append operationals
func appendOperationals(
	tacticalResp *pohonkinerja.TacticalOpdResponse,
	pohonMap map[int]map[int][]domain.PohonKinerja,
	taggingMap map[int][]pohonkinerja.TaggingResponse,
	pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse,
	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
	reviewMap map[int][]pohonkinerja.ReviewResponse,
	tematikMap map[int]*domain.PohonKinerja,
) {
	operationals, ok := pohonMap[6][tacticalResp.Id]
	if !ok {
		return
	}

	// Sort operationals
	sort.Slice(operationals, func(i, j int) bool {
		if operationals[i].Status == "pokin dari pemda" && operationals[j].Status != "pokin dari pemda" {
			return true
		}
		if operationals[i].Status != "pokin dari pemda" && operationals[j].Status == "pokin dari pemda" {
			return false
		}
		return operationals[i].Id < operationals[j].Id
	})

	for _, operational := range operationals {
		opResp := buildOperationalOnly(
			operational,
			taggingMap,
			pelaksanaMap,
			indikatorMap,
			reviewMap,
			tematikMap,
		)

		// Lanjut append operational N
		appendOperationalN(&opResp, pohonMap, taggingMap, pelaksanaMap, indikatorMap, reviewMap)

		tacticalResp.Operationals = append(tacticalResp.Operationals, opResp)
	}
}

func appendOperationalN(
	operationalResp *pohonkinerja.OperationalOpdResponse,
	pohonMap map[int]map[int][]domain.PohonKinerja,
	taggingMap map[int][]pohonkinerja.TaggingResponse,
	pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse,
	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
	reviewMap map[int][]pohonkinerja.ReviewResponse,
) {
	nextLevel := operationalResp.LevelPohon + 1
	children, ok := pohonMap[nextLevel][operationalResp.Id]
	if !ok {
		return
	}

	// Sort children
	sort.Slice(children, func(i, j int) bool {
		if children[i].Status == "pokin dari pemda" && children[j].Status != "pokin dari pemda" {
			return true
		}
		if children[i].Status != "pokin dari pemda" && children[j].Status == "pokin dari pemda" {
			return false
		}
		return children[i].Id < children[j].Id
	})

	for _, child := range children {
		childResp := buildOperationalNOnly(
			child,
			taggingMap,
			pelaksanaMap,
			indikatorMap,
			reviewMap,
		)

		// Recursive untuk level berikutnya jika ada (gunakan function terpisah untuk OperationalNOpdResponse)
		appendOperationalNRecursive(&childResp, pohonMap, taggingMap, pelaksanaMap, indikatorMap, reviewMap)

		operationalResp.Childs = append(operationalResp.Childs, childResp)
	}
}

// Helper function untuk append operational N recursive (untuk OperationalNOpdResponse)
func appendOperationalNRecursive(
	operationalNResp *pohonkinerja.OperationalNOpdResponse,
	pohonMap map[int]map[int][]domain.PohonKinerja,
	taggingMap map[int][]pohonkinerja.TaggingResponse,
	pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse,
	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
	reviewMap map[int][]pohonkinerja.ReviewResponse,
) {
	nextLevel := operationalNResp.LevelPohon + 1
	children, ok := pohonMap[nextLevel][operationalNResp.Id]
	if !ok {
		return
	}

	// Sort children
	sort.Slice(children, func(i, j int) bool {
		if children[i].Status == "pokin dari pemda" && children[j].Status != "pokin dari pemda" {
			return true
		}
		if children[i].Status != "pokin dari pemda" && children[j].Status == "pokin dari pemda" {
			return false
		}
		return children[i].Id < children[j].Id
	})

	for _, child := range children {
		childResp := buildOperationalNOnly(
			child,
			taggingMap,
			pelaksanaMap,
			indikatorMap,
			reviewMap,
		)

		// Recursive untuk level berikutnya jika ada
		appendOperationalNRecursive(&childResp, pohonMap, taggingMap, pelaksanaMap, indikatorMap, reviewMap)

		operationalNResp.Childs = append(operationalNResp.Childs, childResp)
	}
}

// Versi optimized buildStrategicResponse yang menggunakan pre-fetched data
func (service *PohonKinerjaOpdServiceImpl) buildStrategicResponseOptimized(
	ctx context.Context,
	pohonMap map[int]map[int][]domain.PohonKinerja,
	strategic domain.PohonKinerja,
	pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse,
	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
	taggingList []domain.TaggingPokin,
	reviews []domain.Review,
	countReview int,
	reviewPegawaiMap map[string]*domainmaster.Pegawai,
	programUnggulanMap map[string]*domain.ProgramUnggulan,
	taggingMap map[int][]domain.TaggingPokin,
	reviewMap map[int][]domain.Review,
	reviewCountMap map[int]int,
	tematikMap map[int]*domain.PohonKinerja,
) pohonkinerja.StrategicOpdResponse {
	var keteranganCrosscutting *string
	if strategic.KeteranganCrosscutting != nil && *strategic.KeteranganCrosscutting != "" {
		keteranganCrosscutting = strategic.KeteranganCrosscutting
	}

	// Build tagging responses menggunakan pre-fetched data
	var taggingResponses []pohonkinerja.TaggingResponse
	for _, tagging := range taggingList {
		var keteranganResponses []pohonkinerja.KeteranganTaggingResponse
		for _, keterangan := range tagging.KeteranganTaggingProgram {
			if program, exists := programUnggulanMap[keterangan.KodeProgramUnggulan]; exists {
				keteranganResponses = append(keteranganResponses, pohonkinerja.KeteranganTaggingResponse{
					Id:                  keterangan.Id,
					IdTagging:           keterangan.IdTagging,
					KodeProgramUnggulan: keterangan.KodeProgramUnggulan,
					RencanaImplementasi: program.KeteranganProgramUnggulan,
				})
			}
		}

		taggingResponses = append(taggingResponses, pohonkinerja.TaggingResponse{
			Id:                       tagging.Id,
			IdPokin:                  tagging.IdPokin,
			NamaTagging:              tagging.NamaTagging,
			KeteranganTaggingProgram: keteranganResponses,
		})
	}

	// Build review responses menggunakan pre-fetched data
	var reviewResponses []pohonkinerja.ReviewResponse
	for _, review := range reviews {
		pegawai, exists := reviewPegawaiMap[review.CreatedBy]
		if exists {
			reviewResponses = append(reviewResponses, pohonkinerja.ReviewResponse{
				Id:             review.Id,
				IdPohonKinerja: review.IdPohonKinerja,
				Review:         review.Review,
				Keterangan:     review.Keterangan,
				CreatedBy:      pegawai.NamaPegawai,
			})
		}
	}

	// Cari tematik dari pre-fetched map
	var idTematik *int
	var namaTematik *string
	if strategic.CloneFrom > 0 {
		if tematik, exists := tematikMap[strategic.CloneFrom]; exists && tematik != nil {
			idTematik = &tematik.Id
			namaTematik = &tematik.NamaPohon
		}
	}

	strategicResp := pohonkinerja.StrategicOpdResponse{
		Id:                     strategic.Id,
		Parent:                 nil,
		Strategi:               strategic.NamaPohon,
		JenisPohon:             strategic.JenisPohon,
		LevelPohon:             strategic.LevelPohon,
		Keterangan:             strategic.Keterangan,
		KeteranganCrosscutting: keteranganCrosscutting,
		Status:                 strategic.Status,
		IsActive:               strategic.IsActive,
		IdTematik:              idTematik,
		NamaTematik:            namaTematik,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: strategic.KodeOpd,
			NamaOpd: strategic.NamaOpd,
		},
		Tagging:     taggingResponses,
		Pelaksana:   pelaksanaMap[strategic.Id],
		Indikator:   indikatorMap[strategic.Id],
		Review:      reviewResponses,
		CountReview: countReview,
	}

	// Build tactical (level 5)
	if tacticalList := pohonMap[5][strategic.Id]; len(tacticalList) > 0 {
		var tacticals []pohonkinerja.TacticalOpdResponse
		sort.Slice(tacticalList, func(i, j int) bool {
			if tacticalList[i].Status == "pokin dari pemda" && tacticalList[j].Status != "pokin dari pemda" {
				return true
			}
			if tacticalList[i].Status != "pokin dari pemda" && tacticalList[j].Status == "pokin dari pemda" {
				return false
			}
			return tacticalList[i].Id < tacticalList[j].Id
		})

		for _, tactical := range tacticalList {
			tacticalResp := service.buildTacticalResponseOptimized(
				ctx, pohonMap, tactical, pelaksanaMap, indikatorMap,
				taggingMap[tactical.Id], reviewMap[tactical.Id], reviewCountMap[tactical.Id],
				reviewPegawaiMap, programUnggulanMap, idTematik, namaTematik,
				taggingMap, reviewMap, reviewCountMap,
			)
			tacticals = append(tacticals, tacticalResp)
		}
		strategicResp.Tacticals = tacticals
	}

	return strategicResp
}

// Versi optimized buildTacticalResponse
func (service *PohonKinerjaOpdServiceImpl) buildTacticalResponseOptimized(
	ctx context.Context,
	pohonMap map[int]map[int][]domain.PohonKinerja,
	tactical domain.PohonKinerja,
	pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse,
	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
	taggingList []domain.TaggingPokin,
	reviews []domain.Review,
	countReview int,
	reviewPegawaiMap map[string]*domainmaster.Pegawai,
	programUnggulanMap map[string]*domain.ProgramUnggulan,
	idTematik *int,
	namaTematik *string,
	taggingMap map[int][]domain.TaggingPokin,
	reviewMap map[int][]domain.Review,
	reviewCountMap map[int]int,
) pohonkinerja.TacticalOpdResponse {
	var keteranganCrosscutting *string
	if tactical.KeteranganCrosscutting != nil && *tactical.KeteranganCrosscutting != "" {
		keteranganCrosscutting = tactical.KeteranganCrosscutting
	}

	// Build tagging dan review menggunakan pre-fetched data
	var taggingResponses []pohonkinerja.TaggingResponse
	for _, tagging := range taggingList {
		var keteranganResponses []pohonkinerja.KeteranganTaggingResponse
		for _, keterangan := range tagging.KeteranganTaggingProgram {
			if program, exists := programUnggulanMap[keterangan.KodeProgramUnggulan]; exists {
				keteranganResponses = append(keteranganResponses, pohonkinerja.KeteranganTaggingResponse{
					Id:                  keterangan.Id,
					IdTagging:           keterangan.IdTagging,
					KodeProgramUnggulan: keterangan.KodeProgramUnggulan,
					RencanaImplementasi: program.KeteranganProgramUnggulan,
				})
			}
		}
		taggingResponses = append(taggingResponses, pohonkinerja.TaggingResponse{
			Id:                       tagging.Id,
			IdPokin:                  tagging.IdPokin,
			NamaTagging:              tagging.NamaTagging,
			KeteranganTaggingProgram: keteranganResponses,
		})
	}

	var reviewResponses []pohonkinerja.ReviewResponse
	for _, review := range reviews {
		pegawai, exists := reviewPegawaiMap[review.CreatedBy]
		if exists {
			reviewResponses = append(reviewResponses, pohonkinerja.ReviewResponse{
				Id:             review.Id,
				IdPohonKinerja: review.IdPohonKinerja,
				Review:         review.Review,
				Keterangan:     review.Keterangan,
				CreatedBy:      pegawai.NamaPegawai,
			})
		}
	}

	tacticalResp := pohonkinerja.TacticalOpdResponse{
		Id:                     tactical.Id,
		Parent:                 tactical.Parent,
		Strategi:               tactical.NamaPohon,
		JenisPohon:             tactical.JenisPohon,
		LevelPohon:             tactical.LevelPohon,
		Keterangan:             tactical.Keterangan,
		KeteranganCrosscutting: keteranganCrosscutting,
		Status:                 tactical.Status,
		IsActive:               tactical.IsActive,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: tactical.KodeOpd,
			NamaOpd: tactical.NamaOpd,
		},
		IdTematik:   idTematik,
		NamaTematik: namaTematik,
		Pelaksana:   pelaksanaMap[tactical.Id],
		Tagging:     taggingResponses,
		Indikator:   indikatorMap[tactical.Id],
		Review:      reviewResponses,
		CountReview: countReview,
	}

	// Build operational (level 6)
	if operationalList := pohonMap[6][tactical.Id]; len(operationalList) > 0 {
		var operationals []pohonkinerja.OperationalOpdResponse
		sort.Slice(operationalList, func(i, j int) bool {
			if operationalList[i].Status == "pokin dari pemda" && operationalList[j].Status != "pokin dari pemda" {
				return true
			}
			if operationalList[i].Status != "pokin dari pemda" && operationalList[j].Status == "pokin dari pemda" {
				return false
			}
			return operationalList[i].Id < operationalList[j].Id
		})

		for _, operational := range operationalList {
			operationalResp := service.buildOperationalResponseOptimized(
				ctx, pohonMap, operational, pelaksanaMap, indikatorMap,
				taggingMap[operational.Id], reviewMap[operational.Id], reviewCountMap[operational.Id],
				reviewPegawaiMap, programUnggulanMap, idTematik, namaTematik,
				taggingMap, reviewMap, reviewCountMap,
			)
			operationals = append(operationals, operationalResp)
		}
		tacticalResp.Operationals = operationals
	}

	return tacticalResp
}

// Versi optimized buildOperationalResponse
func (service *PohonKinerjaOpdServiceImpl) buildOperationalResponseOptimized(
	ctx context.Context,
	pohonMap map[int]map[int][]domain.PohonKinerja,
	operational domain.PohonKinerja,
	pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse,
	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
	taggingList []domain.TaggingPokin,
	reviews []domain.Review,
	countReview int,
	reviewPegawaiMap map[string]*domainmaster.Pegawai,
	programUnggulanMap map[string]*domain.ProgramUnggulan,
	idTematik *int,
	namaTematik *string,
	taggingMap map[int][]domain.TaggingPokin,
	reviewMap map[int][]domain.Review,
	reviewCountMap map[int]int,
) pohonkinerja.OperationalOpdResponse {
	var keteranganCrosscutting *string
	if operational.KeteranganCrosscutting != nil && *operational.KeteranganCrosscutting != "" {
		keteranganCrosscutting = operational.KeteranganCrosscutting
	}

	// Build tagging dan review menggunakan pre-fetched data
	var taggingResponses []pohonkinerja.TaggingResponse
	for _, tagging := range taggingList {
		var keteranganResponses []pohonkinerja.KeteranganTaggingResponse
		for _, keterangan := range tagging.KeteranganTaggingProgram {
			if program, exists := programUnggulanMap[keterangan.KodeProgramUnggulan]; exists {
				keteranganResponses = append(keteranganResponses, pohonkinerja.KeteranganTaggingResponse{
					Id:                  keterangan.Id,
					IdTagging:           keterangan.IdTagging,
					KodeProgramUnggulan: keterangan.KodeProgramUnggulan,
					RencanaImplementasi: program.KeteranganProgramUnggulan,
				})
			}
		}
		taggingResponses = append(taggingResponses, pohonkinerja.TaggingResponse{
			Id:                       tagging.Id,
			IdPokin:                  tagging.IdPokin,
			NamaTagging:              tagging.NamaTagging,
			KeteranganTaggingProgram: keteranganResponses,
		})
	}

	var reviewResponses []pohonkinerja.ReviewResponse
	for _, review := range reviews {
		pegawai, exists := reviewPegawaiMap[review.CreatedBy]
		if exists {
			reviewResponses = append(reviewResponses, pohonkinerja.ReviewResponse{
				Id:             review.Id,
				IdPohonKinerja: review.IdPohonKinerja,
				Review:         review.Review,
				Keterangan:     review.Keterangan,
				CreatedBy:      pegawai.NamaPegawai,
			})
		}
	}

	operationalResp := pohonkinerja.OperationalOpdResponse{
		Id:                     operational.Id,
		Parent:                 operational.Parent,
		Strategi:               operational.NamaPohon,
		JenisPohon:             operational.JenisPohon,
		LevelPohon:             operational.LevelPohon,
		Keterangan:             operational.Keterangan,
		KeteranganCrosscutting: keteranganCrosscutting,
		Status:                 operational.Status,
		IsActive:               operational.IsActive,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: operational.KodeOpd,
			NamaOpd: operational.NamaOpd,
		},
		IdTematik:   idTematik,
		NamaTematik: namaTematik,
		Pelaksana:   pelaksanaMap[operational.Id],
		Tagging:     taggingResponses,
		Indikator:   indikatorMap[operational.Id],
		Review:      reviewResponses,
		CountReview: countReview,
	}

	// Build operational-n untuk level > 6
	nextLevel := operational.LevelPohon + 1
	if operationalNList := pohonMap[nextLevel][operational.Id]; len(operationalNList) > 0 {
		var childs []pohonkinerja.OperationalNOpdResponse
		sort.Slice(operationalNList, func(i, j int) bool {
			if operationalNList[i].Status == "pokin dari pemda" && operationalNList[j].Status != "pokin dari pemda" {
				return true
			}
			if operationalNList[i].Status != "pokin dari pemda" && operationalNList[j].Status == "pokin dari pemda" {
				return false
			}
			return operationalNList[i].Id < operationalNList[j].Id
		})

		for _, opN := range operationalNList {
			childResp := service.buildOperationalNResponseOptimized(
				ctx, pohonMap, opN, pelaksanaMap, indikatorMap,
				taggingMap[opN.Id], reviewMap[opN.Id], reviewCountMap[opN.Id],
				reviewPegawaiMap, programUnggulanMap,
				taggingMap, reviewMap, reviewCountMap,
			)
			childs = append(childs, childResp)
		}
		operationalResp.Childs = childs
	}

	return operationalResp
}

// Versi optimized buildOperationalNResponse
func (service *PohonKinerjaOpdServiceImpl) buildOperationalNResponseOptimized(
	ctx context.Context,
	pohonMap map[int]map[int][]domain.PohonKinerja,
	operationalN domain.PohonKinerja,
	pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse,
	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
	taggingList []domain.TaggingPokin,
	reviews []domain.Review,
	countReview int,
	reviewPegawaiMap map[string]*domainmaster.Pegawai,
	programUnggulanMap map[string]*domain.ProgramUnggulan,
	taggingMap map[int][]domain.TaggingPokin,
	reviewMap map[int][]domain.Review,
	reviewCountMap map[int]int,
) pohonkinerja.OperationalNOpdResponse {
	var keteranganCrosscutting *string
	if operationalN.KeteranganCrosscutting != nil && *operationalN.KeteranganCrosscutting != "" {
		keteranganCrosscutting = operationalN.KeteranganCrosscutting
	}

	// Build tagging dan review menggunakan pre-fetched data
	var taggingResponses []pohonkinerja.TaggingResponse
	for _, tagging := range taggingList {
		var keteranganResponses []pohonkinerja.KeteranganTaggingResponse
		for _, keterangan := range tagging.KeteranganTaggingProgram {
			if program, exists := programUnggulanMap[keterangan.KodeProgramUnggulan]; exists {
				keteranganResponses = append(keteranganResponses, pohonkinerja.KeteranganTaggingResponse{
					Id:                  keterangan.Id,
					IdTagging:           keterangan.IdTagging,
					KodeProgramUnggulan: keterangan.KodeProgramUnggulan,
					RencanaImplementasi: program.KeteranganProgramUnggulan,
				})
			}
		}
		taggingResponses = append(taggingResponses, pohonkinerja.TaggingResponse{
			Id:                       tagging.Id,
			IdPokin:                  tagging.IdPokin,
			NamaTagging:              tagging.NamaTagging,
			KeteranganTaggingProgram: keteranganResponses,
		})
	}

	var reviewResponses []pohonkinerja.ReviewResponse
	for _, review := range reviews {
		pegawai, exists := reviewPegawaiMap[review.CreatedBy]
		if exists {
			reviewResponses = append(reviewResponses, pohonkinerja.ReviewResponse{
				Id:             review.Id,
				IdPohonKinerja: review.IdPohonKinerja,
				Review:         review.Review,
				Keterangan:     review.Keterangan,
				CreatedBy:      pegawai.NamaPegawai,
			})
		}
	}

	operationalNResp := pohonkinerja.OperationalNOpdResponse{
		Id:                     operationalN.Id,
		Parent:                 operationalN.Parent,
		Strategi:               operationalN.NamaPohon,
		JenisPohon:             operationalN.JenisPohon,
		LevelPohon:             operationalN.LevelPohon,
		Keterangan:             operationalN.Keterangan,
		KeteranganCrosscutting: keteranganCrosscutting,
		Status:                 operationalN.Status,
		IsActive:               operationalN.IsActive,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: operationalN.KodeOpd,
			NamaOpd: operationalN.NamaOpd,
		},
		Tagging:     taggingResponses,
		Pelaksana:   pelaksanaMap[operationalN.Id],
		Indikator:   indikatorMap[operationalN.Id],
		Review:      reviewResponses,
		CountReview: countReview,
	}

	// Build child nodes secara rekursif
	nextLevel := operationalN.LevelPohon + 1
	if nextOperationalNList := pohonMap[nextLevel][operationalN.Id]; len(nextOperationalNList) > 0 {
		var childs []pohonkinerja.OperationalNOpdResponse
		sort.Slice(nextOperationalNList, func(i, j int) bool {
			if nextOperationalNList[i].Status == "pokin dari pemda" && nextOperationalNList[j].Status != "pokin dari pemda" {
				return true
			}
			if nextOperationalNList[i].Status != "pokin dari pemda" && nextOperationalNList[j].Status == "pokin dari pemda" {
				return false
			}
			return nextOperationalNList[i].Id < nextOperationalNList[j].Id
		})

		for _, child := range nextOperationalNList {
			childResp := service.buildOperationalNResponseOptimized(
				ctx, pohonMap, child, pelaksanaMap, indikatorMap,
				taggingMap[child.Id], reviewMap[child.Id], reviewCountMap[child.Id],
				reviewPegawaiMap, programUnggulanMap,
				taggingMap, reviewMap, reviewCountMap,
			)
			childs = append(childs, childResp)
		}
		operationalNResp.Childs = childs
	}

	return operationalNResp
}

//end optimasi

func (service *PohonKinerjaOpdServiceImpl) FindStrategicNoParent(ctx context.Context, kodeOpd, tahun string) ([]pohonkinerja.StrategicOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi kode OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return nil, errors.New("kode opd tidak ditemukan")
	}

	// Ambil data strategic dengan level pohon 4
	pokins, err := service.pohonKinerjaOpdRepository.FindStrategicNoParent(ctx, tx, 4, 0, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}

	// Urutkan data berdasarkan ID
	sort.Slice(pokins, func(i, j int) bool {
		return pokins[i].Id < pokins[j].Id
	})

	// Konversi ke response format
	var strategics []pohonkinerja.StrategicOpdResponse
	for _, pokin := range pokins {
		strategic := pohonkinerja.StrategicOpdResponse{
			Id: pokin.Id,
			KodeOpd: opdmaster.OpdResponseForAll{
				KodeOpd: kodeOpd,
				NamaOpd: opd.NamaOpd,
			},
			Strategi:   pokin.NamaPohon,
			Keterangan: pokin.Keterangan,
		}
		strategics = append(strategics, strategic)
	}

	return strategics, nil
}

func (service *PohonKinerjaOpdServiceImpl) DeletePelaksana(ctx context.Context, pelaksanaId string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)
	return service.pohonKinerjaOpdRepository.DeletePelaksanaPokin(ctx, tx, pelaksanaId)
}

// Tambahkan fungsi helper untuk membangun OperationalN response
func (service *PohonKinerjaOpdServiceImpl) buildOperationalNResponse(ctx context.Context, tx *sql.Tx, pohonMap map[int]map[int][]domain.PohonKinerja, operationalN domain.PohonKinerja, pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse, indikatorMap map[int][]pohonkinerja.IndikatorResponse) pohonkinerja.OperationalNOpdResponse {
	var keteranganCrosscutting *string
	if operationalN.KeteranganCrosscutting != nil && *operationalN.KeteranganCrosscutting != "" {
		keteranganCrosscutting = operationalN.KeteranganCrosscutting
	}
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, operationalN.KodeOpd)
	if err == nil {
		operationalN.NamaOpd = opd.NamaOpd
	}

	//tagging
	taggingList, err := service.pohonKinerjaOpdRepository.FindTaggingByPokinId(ctx, tx, operationalN.Id)
	var taggingResponses []pohonkinerja.TaggingResponse
	if err == nil {
		for _, tagging := range taggingList {
			var keteranganResponses []pohonkinerja.KeteranganTaggingResponse
			for _, keterangan := range tagging.KeteranganTaggingProgram {
				programUnggulan, err := service.ProgramUnggulanRepository.FindByKodeProgramUnggulan(ctx, tx, keterangan.KodeProgramUnggulan)
				if err != nil {
					continue
				}
				keteranganResponses = append(keteranganResponses, pohonkinerja.KeteranganTaggingResponse{
					Id:                  keterangan.Id,
					IdTagging:           keterangan.IdTagging,
					KodeProgramUnggulan: keterangan.KodeProgramUnggulan,
					RencanaImplementasi: programUnggulan.KeteranganProgramUnggulan,
				})
			}

			taggingResponses = append(taggingResponses, pohonkinerja.TaggingResponse{
				Id:                       tagging.Id,
				IdPokin:                  tagging.IdPokin,
				NamaTagging:              tagging.NamaTagging,
				KeteranganTaggingProgram: keteranganResponses,
			})
		}
	}

	//review
	countReview, err := service.reviewRepository.CountReviewByPohonKinerja(ctx, tx, operationalN.Id)
	helper.PanicIfError(err)

	reviews, err := service.reviewRepository.FindByPohonKinerja(ctx, tx, operationalN.Id)
	var reviewResponses []pohonkinerja.ReviewResponse
	if err == nil {
		for _, review := range reviews {
			pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, reviews[0].CreatedBy)
			if err != nil {
				return pohonkinerja.OperationalNOpdResponse{}
			}
			reviewResponses = append(reviewResponses, pohonkinerja.ReviewResponse{
				Id:             review.Id,
				IdPohonKinerja: review.IdPohonKinerja,
				Review:         review.Review,
				Keterangan:     review.Keterangan,
				CreatedBy:      pegawai.NamaPegawai,
			})
		}
	}
	operationalNResp := pohonkinerja.OperationalNOpdResponse{
		Id:                     operationalN.Id,
		Parent:                 operationalN.Parent,
		Strategi:               operationalN.NamaPohon,
		JenisPohon:             operationalN.JenisPohon,
		LevelPohon:             operationalN.LevelPohon,
		Keterangan:             operationalN.Keterangan,
		KeteranganCrosscutting: keteranganCrosscutting,
		Status:                 operationalN.Status,
		IsActive:               operationalN.IsActive,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: operationalN.KodeOpd,
			NamaOpd: operationalN.NamaOpd,
		},
		Tagging:     taggingResponses,
		Pelaksana:   pelaksanaMap[operationalN.Id],
		Indikator:   indikatorMap[operationalN.Id],
		Review:      reviewResponses,
		CountReview: countReview,
	}

	// Build child nodes secara rekursif
	nextLevel := operationalN.LevelPohon + 1
	if nextOperationalNList := pohonMap[nextLevel][operationalN.Id]; len(nextOperationalNList) > 0 {
		var childs []pohonkinerja.OperationalNOpdResponse
		sort.Slice(nextOperationalNList, func(i, j int) bool {
			return nextOperationalNList[i].Id < nextOperationalNList[j].Id
		})

		for _, nextOpN := range nextOperationalNList {
			childResp := service.buildOperationalNResponse(ctx, tx, pohonMap, nextOpN, pelaksanaMap, indikatorMap)
			childs = append(childs, childResp)
		}
		operationalNResp.Childs = childs
	}

	return operationalNResp
}

// Helper functions untuk membangun response
func (service *PohonKinerjaOpdServiceImpl) buildStrategicResponse(ctx context.Context, tx *sql.Tx, pohonMap map[int]map[int][]domain.PohonKinerja, strategic domain.PohonKinerja, pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse, indikatorMap map[int][]pohonkinerja.IndikatorResponse) pohonkinerja.StrategicOpdResponse {
	var keteranganCrosscutting *string
	if strategic.KeteranganCrosscutting != nil && *strategic.KeteranganCrosscutting != "" {
		keteranganCrosscutting = strategic.KeteranganCrosscutting
	}
	//tagging
	var taggingResponses []pohonkinerja.TaggingResponse
	taggingList, err := service.pohonKinerjaOpdRepository.FindTaggingByPokinId(ctx, tx, strategic.Id)
	if err == nil {
		for _, tagging := range taggingList {
			var keteranganResponses []pohonkinerja.KeteranganTaggingResponse
			for _, keterangan := range tagging.KeteranganTaggingProgram {

				programUnggulan, err := service.ProgramUnggulanRepository.FindByKodeProgramUnggulan(ctx, tx, keterangan.KodeProgramUnggulan)
				if err != nil {
					continue
				}
				keteranganResponses = append(keteranganResponses, pohonkinerja.KeteranganTaggingResponse{
					Id:                  keterangan.Id,
					IdTagging:           keterangan.IdTagging,
					KodeProgramUnggulan: keterangan.KodeProgramUnggulan,
					RencanaImplementasi: programUnggulan.KeteranganProgramUnggulan,
				})
			}

			taggingResponses = append(taggingResponses, pohonkinerja.TaggingResponse{
				Id:                       tagging.Id,
				IdPokin:                  tagging.IdPokin,
				NamaTagging:              tagging.NamaTagging,
				KeteranganTaggingProgram: keteranganResponses,
			})
		}
	}
	//review
	countReview, err := service.reviewRepository.CountReviewByPohonKinerja(ctx, tx, strategic.Id)
	helper.PanicIfError(err)

	reviews, err := service.reviewRepository.FindByPohonKinerja(ctx, tx, strategic.Id)
	var reviewResponses []pohonkinerja.ReviewResponse
	if err == nil {
		for _, review := range reviews {
			pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, reviews[0].CreatedBy)
			if err != nil {
				return pohonkinerja.StrategicOpdResponse{}
			}
			reviewResponses = append(reviewResponses, pohonkinerja.ReviewResponse{
				Id:             review.Id,
				IdPohonKinerja: review.IdPohonKinerja,
				Review:         review.Review,
				Keterangan:     review.Keterangan,
				CreatedBy:      pegawai.NamaPegawai,
			})
		}
	}

	// Cari tematik jika ada clone_from
	var idTematik *int
	var namaTematik *string
	if strategic.CloneFrom > 0 {
		tematik, err := service.pohonKinerjaOpdRepository.FindTematikByCloneFrom(ctx, tx, strategic.CloneFrom)
		if err == nil && tematik != nil {
			idTematik = &tematik.Id
			namaTematik = &tematik.NamaPohon
		}
	}

	strategicResp := pohonkinerja.StrategicOpdResponse{
		Id:                     strategic.Id,
		Parent:                 nil,
		Strategi:               strategic.NamaPohon,
		JenisPohon:             strategic.JenisPohon,
		LevelPohon:             strategic.LevelPohon,
		Keterangan:             strategic.Keterangan,
		KeteranganCrosscutting: keteranganCrosscutting,
		Status:                 strategic.Status,
		IsActive:               strategic.IsActive,
		IdTematik:              idTematik,
		NamaTematik:            namaTematik,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: strategic.KodeOpd,
			NamaOpd: strategic.NamaOpd,
		},
		Tagging:     taggingResponses,
		Pelaksana:   pelaksanaMap[strategic.Id],
		Indikator:   indikatorMap[strategic.Id],
		Review:      reviewResponses,
		CountReview: countReview,
	}

	// Build tactical (level 5)
	if tacticalList := pohonMap[5][strategic.Id]; len(tacticalList) > 0 {
		var tacticals []pohonkinerja.TacticalOpdResponse
		sort.Slice(tacticalList, func(i, j int) bool {
			if tacticalList[i].Status == "pokin dari pemda" && tacticalList[j].Status != "pokin dari pemda" {
				return true
			}
			if tacticalList[i].Status != "pokin dari pemda" && tacticalList[j].Status == "pokin dari pemda" {
				return false
			}
			return tacticalList[i].Id < tacticalList[j].Id
		})

		for _, tactical := range tacticalList {
			tacticalResp := service.buildTacticalResponse(ctx, tx, pohonMap, tactical, pelaksanaMap, indikatorMap)
			tacticalResp.IdTematik = idTematik
			tacticalResp.NamaTematik = namaTematik
			tacticals = append(tacticals, tacticalResp)
		}
		strategicResp.Tacticals = tacticals
	}

	return strategicResp
}

func (service *PohonKinerjaOpdServiceImpl) buildTacticalResponse(ctx context.Context, tx *sql.Tx, pohonMap map[int]map[int][]domain.PohonKinerja, tactical domain.PohonKinerja, pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse, indikatorMap map[int][]pohonkinerja.IndikatorResponse) pohonkinerja.TacticalOpdResponse {
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, tactical.KodeOpd)
	if err == nil {
		tactical.NamaOpd = opd.NamaOpd
	}
	var keteranganCrosscutting *string
	if tactical.KeteranganCrosscutting != nil && *tactical.KeteranganCrosscutting != "" {
		keteranganCrosscutting = tactical.KeteranganCrosscutting
	}
	//tagging
	var taggingResponses []pohonkinerja.TaggingResponse
	taggingList, err := service.pohonKinerjaOpdRepository.FindTaggingByPokinId(ctx, tx, tactical.Id)
	if err == nil {
		for _, tagging := range taggingList {
			var keteranganResponses []pohonkinerja.KeteranganTaggingResponse
			for _, keterangan := range tagging.KeteranganTaggingProgram {
				programUnggulan, err := service.ProgramUnggulanRepository.FindByKodeProgramUnggulan(ctx, tx, keterangan.KodeProgramUnggulan)
				if err != nil {
					continue
				}
				keteranganResponses = append(keteranganResponses, pohonkinerja.KeteranganTaggingResponse{
					Id:                  keterangan.Id,
					IdTagging:           keterangan.IdTagging,
					KodeProgramUnggulan: keterangan.KodeProgramUnggulan,
					RencanaImplementasi: programUnggulan.KeteranganProgramUnggulan,
				})
			}

			taggingResponses = append(taggingResponses, pohonkinerja.TaggingResponse{
				Id:                       tagging.Id,
				IdPokin:                  tagging.IdPokin,
				NamaTagging:              tagging.NamaTagging,
				KeteranganTaggingProgram: keteranganResponses,
			})
		}
	}
	//review
	countReview, err := service.reviewRepository.CountReviewByPohonKinerja(ctx, tx, tactical.Id)
	helper.PanicIfError(err)
	reviews, err := service.reviewRepository.FindByPohonKinerja(ctx, tx, tactical.Id)
	var reviewResponses []pohonkinerja.ReviewResponse
	if err == nil {
		for _, review := range reviews {
			pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, reviews[0].CreatedBy)
			if err != nil {
				return pohonkinerja.TacticalOpdResponse{}
			}
			reviewResponses = append(reviewResponses, pohonkinerja.ReviewResponse{
				Id:             review.Id,
				IdPohonKinerja: review.IdPohonKinerja,
				Review:         review.Review,
				Keterangan:     review.Keterangan,
				CreatedBy:      pegawai.NamaPegawai,
			})
		}
	}
	// Cari tematik jika ada clone_from
	var idTematik *int
	var namaTematik *string
	if tactical.CloneFrom > 0 {
		tematik, err := service.pohonKinerjaOpdRepository.FindTematikByCloneFrom(ctx, tx, tactical.CloneFrom)
		if err == nil && tematik != nil {
			idTematik = &tematik.Id
			namaTematik = &tematik.NamaPohon
		}
	}
	tacticalResp := pohonkinerja.TacticalOpdResponse{
		Id:                     tactical.Id,
		Parent:                 tactical.Parent,
		Strategi:               tactical.NamaPohon,
		JenisPohon:             tactical.JenisPohon,
		LevelPohon:             tactical.LevelPohon,
		Keterangan:             tactical.Keterangan,
		KeteranganCrosscutting: keteranganCrosscutting,
		Status:                 tactical.Status,
		IsActive:               tactical.IsActive,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: tactical.KodeOpd,
			NamaOpd: tactical.NamaOpd,
		},
		IdTematik:   idTematik,
		NamaTematik: namaTematik,
		Pelaksana:   pelaksanaMap[tactical.Id],
		Tagging:     taggingResponses,
		Indikator:   indikatorMap[tactical.Id],
		Review:      reviewResponses,
		CountReview: countReview,
	}

	// Build operational (level 6)
	if operationalList := pohonMap[6][tactical.Id]; len(operationalList) > 0 {
		var operationals []pohonkinerja.OperationalOpdResponse
		sort.Slice(operationalList, func(i, j int) bool {
			if operationalList[i].Status == "pokin dari pemda" && operationalList[j].Status != "pokin dari pemda" {
				return true
			}
			if operationalList[i].Status != "pokin dari pemda" && operationalList[j].Status == "pokin dari pemda" {
				return false
			}
			return operationalList[i].Id < operationalList[j].Id
		})

		for _, operational := range operationalList {
			operationalResp := service.buildOperationalResponse(ctx, tx, pohonMap, operational, pelaksanaMap, indikatorMap)
			operationalResp.IdTematik = idTematik
			operationalResp.NamaTematik = namaTematik
			operationals = append(operationals, operationalResp)
		}
		tacticalResp.Operationals = operationals
	}

	return tacticalResp
}

func (service *PohonKinerjaOpdServiceImpl) buildOperationalResponse(ctx context.Context, tx *sql.Tx, pohonMap map[int]map[int][]domain.PohonKinerja, operational domain.PohonKinerja, pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse, indikatorMap map[int][]pohonkinerja.IndikatorResponse) pohonkinerja.OperationalOpdResponse {
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, operational.KodeOpd)
	if err == nil {
		operational.NamaOpd = opd.NamaOpd
	}
	var keteranganCrosscutting *string
	if operational.KeteranganCrosscutting != nil && *operational.KeteranganCrosscutting != "" {
		keteranganCrosscutting = operational.KeteranganCrosscutting
	}
	//review
	countReview, err := service.reviewRepository.CountReviewByPohonKinerja(ctx, tx, operational.Id)
	helper.PanicIfError(err)

	//tagging
	var taggingResponses []pohonkinerja.TaggingResponse
	taggingList, err := service.pohonKinerjaOpdRepository.FindTaggingByPokinId(ctx, tx, operational.Id)
	if err == nil {
		for _, tagging := range taggingList {
			var keteranganResponses []pohonkinerja.KeteranganTaggingResponse
			for _, keterangan := range tagging.KeteranganTaggingProgram {
				programUnggulan, err := service.ProgramUnggulanRepository.FindByKodeProgramUnggulan(ctx, tx, keterangan.KodeProgramUnggulan)
				if err != nil {
					continue
				}
				keteranganResponses = append(keteranganResponses, pohonkinerja.KeteranganTaggingResponse{
					Id:                  keterangan.Id,
					IdTagging:           keterangan.IdTagging,
					KodeProgramUnggulan: keterangan.KodeProgramUnggulan,
					RencanaImplementasi: programUnggulan.KeteranganProgramUnggulan,
				})
			}

			taggingResponses = append(taggingResponses, pohonkinerja.TaggingResponse{
				Id:                       tagging.Id,
				IdPokin:                  tagging.IdPokin,
				NamaTagging:              tagging.NamaTagging,
				KeteranganTaggingProgram: keteranganResponses,
			})
		}
	}
	reviews, err := service.reviewRepository.FindByPohonKinerja(ctx, tx, operational.Id)
	var reviewResponses []pohonkinerja.ReviewResponse
	if err == nil {
		for _, review := range reviews {
			pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, reviews[0].CreatedBy)
			if err != nil {
				return pohonkinerja.OperationalOpdResponse{}
			}
			reviewResponses = append(reviewResponses, pohonkinerja.ReviewResponse{
				Id:             review.Id,
				IdPohonKinerja: review.IdPohonKinerja,
				Review:         review.Review,
				Keterangan:     review.Keterangan,
				CreatedBy:      pegawai.NamaPegawai,
			})
		}
	}
	// Cari tematik jika ada clone_from
	var idTematik *int
	var namaTematik *string
	if operational.CloneFrom > 0 {
		tematik, err := service.pohonKinerjaOpdRepository.FindTematikByCloneFrom(ctx, tx, operational.CloneFrom)
		if err == nil && tematik != nil {
			idTematik = &tematik.Id
			namaTematik = &tematik.NamaPohon
		}
	}
	operationalResp := pohonkinerja.OperationalOpdResponse{
		Id:                     operational.Id,
		Parent:                 operational.Parent,
		Strategi:               operational.NamaPohon,
		JenisPohon:             operational.JenisPohon,
		LevelPohon:             operational.LevelPohon,
		Keterangan:             operational.Keterangan,
		KeteranganCrosscutting: keteranganCrosscutting,
		Status:                 operational.Status,
		IsActive:               operational.IsActive,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: operational.KodeOpd,
			NamaOpd: operational.NamaOpd,
		},
		IdTematik:   idTematik,
		NamaTematik: namaTematik,
		Pelaksana:   pelaksanaMap[operational.Id],
		Tagging:     taggingResponses,
		Indikator:   indikatorMap[operational.Id],
		Review:      reviewResponses,
		CountReview: countReview,
	}

	// Build operational-n untuk level > 6
	nextLevel := operational.LevelPohon + 1
	if operationalNList := pohonMap[nextLevel][operational.Id]; len(operationalNList) > 0 {
		var childs []pohonkinerja.OperationalNOpdResponse
		// Ubah pengurutan untuk operational-n
		sort.Slice(operationalNList, func(i, j int) bool {
			// Prioritaskan status "pokin dari pemda"
			if operationalNList[i].Status == "pokin dari pemda" && operationalNList[j].Status != "pokin dari pemda" {
				return true
			}
			if operationalNList[i].Status != "pokin dari pemda" && operationalNList[j].Status == "pokin dari pemda" {
				return false
			}
			return operationalNList[i].Id < operationalNList[j].Id
		})

		for _, opN := range operationalNList {
			childResp := service.buildOperationalNResponse(ctx, tx, pohonMap, opN, pelaksanaMap, indikatorMap)
			childs = append(childs, childResp)
		}
		operationalResp.Childs = childs
	}

	return operationalResp
}

func (service *PohonKinerjaOpdServiceImpl) FindPokinByPelaksana(ctx context.Context, nip string, tahun string) ([]pohonkinerja.PohonKinerjaOpdResponse, error) {
	log.Printf("Memulai proses FindPokinByPelaksana untuk NIP: %s", nip)

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// ✅ VALIDASI PEGAWAI BERDASARKAN NIP
	pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, nip) // ✅ GUNAKAN FindByNip
	if err != nil {
		log.Printf("Pegawai dengan NIP %s tidak ditemukan: %v", nip, err)
		return nil, fmt.Errorf("pegawai dengan NIP %s tidak ditemukan: %v", nip, err)
	}

	// ✅ AMBIL DATA POHON KINERJA BERDASARKAN NIP
	pokinList, err := service.pohonKinerjaOpdRepository.FindPokinByPelaksana(ctx, tx, nip, tahun) // ✅ PARAMETER NIP
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Tidak ada pohon kinerja untuk NIP: %s", nip)
			return []pohonkinerja.PohonKinerjaOpdResponse{}, nil
		}
		log.Printf("Gagal mengambil data pohon kinerja: %v", err)
		return nil, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
	}

	var responses []pohonkinerja.PohonKinerjaOpdResponse
	for _, pokin := range pokinList {
		// Ambil data OPD
		opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokin.KodeOpd)
		if err != nil {
			log.Printf("Gagal mengambil data OPD: %v", err)
			return nil, fmt.Errorf("gagal mengambil data OPD: %v", err)
		}

		indikators, err := service.pohonKinerjaOpdRepository.FindIndikatorByPokinId(ctx, tx, fmt.Sprint(pokin.Id))
		var indikatorResponses []pohonkinerja.IndikatorResponse
		if err == nil {
			for _, indikator := range indikators {
				var targetResponses []pohonkinerja.TargetResponse
				for _, target := range indikator.Target {
					targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
						Id:              target.Id,
						IndikatorId:     target.IndikatorId,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					})
				}

				indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
					Id:            indikator.Id,
					IdPokin:       indikator.PokinId,
					NamaIndikator: indikator.Indikator,
					Target:        targetResponses,
				})
			}
		}

		// ✅ BUAT RESPONSE PELAKSANA BERDASARKAN DATA YANG SUDAH ADA
		pelaksanaResponse := pohonkinerja.PelaksanaOpdResponse{
			Id:             pokin.Pelaksana[0].Id,
			PohonKinerjaId: fmt.Sprint(pokin.Id),
			PegawaiId:      pokin.Pelaksana[0].PegawaiId,   // ✅ GUNAKAN ID PEGAWAI DARI DATABASE
			NamaPegawai:    pokin.Pelaksana[0].NamaPegawai, // ✅ GUNAKAN NAMA DARI DATABASE
		}

		responses = append(responses, pohonkinerja.PohonKinerjaOpdResponse{
			Id:         pokin.Id,
			Parent:     fmt.Sprint(pokin.Parent),
			NamaPohon:  pokin.NamaPohon,
			JenisPohon: pokin.JenisPohon,
			LevelPohon: pokin.LevelPohon,
			KodeOpd:    opd.KodeOpd,
			NamaOpd:    opd.NamaOpd,
			Keterangan: pokin.Keterangan,
			Tahun:      pokin.Tahun,
			Indikator:  indikatorResponses,
			Pelaksana:  []pohonkinerja.PelaksanaOpdResponse{pelaksanaResponse},
		})
	}

	log.Printf("Berhasil mengambil %d pohon kinerja untuk NIP %s (%s)", len(responses), nip, pegawai.NamaPegawai)
	return responses, nil
}

// func (service *PohonKinerjaOpdServiceImpl) buildCrosscuttingResponse(ctx context.Context, tx *sql.Tx, pokinId int, pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse, indikatorMap map[int][]pohonkinerja.IndikatorResponse) []pohonkinerja.CrosscuttingOpdResponse {
// 	// Ambil data crosscutting berdasarkan crosscutting_from
// 	crosscuttings, err := service.crosscuttingOpdRepository.FindAllCrosscutting(ctx, tx, pokinId)
// 	if err != nil {
// 		log.Printf("Error getting crosscutting data: %v", err)
// 		return nil
// 	}

// 	var crosscuttingResponses []pohonkinerja.CrosscuttingOpdResponse
// 	for _, crosscutting := range crosscuttings {
// 		// Filter status crosscutting yang akan ditampilkan
// 		if crosscutting.Status != "crosscutting_disetujui" &&
// 			crosscutting.Status != "crosscutting_disetujui_existing" {
// 			continue
// 		}

// 		// Ambil data OPD
// 		opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, crosscutting.KodeOpd)
// 		if err != nil {
// 			log.Printf("Gagal mengambil data OPD: %v", err)
// 			continue
// 		}

// 		// Ambil data indikator untuk crosscutting
// 		indikatorList, err := service.pohonKinerjaOpdRepository.FindIndikatorByPokinId(ctx, tx, fmt.Sprint(crosscutting.Id))
// 		if err == nil {
// 			var indikatorResponses []pohonkinerja.IndikatorResponse
// 			for _, indikator := range indikatorList {
// 				// Ambil target untuk setiap indikator
// 				targetList, err := service.pohonKinerjaOpdRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
// 				if err != nil {
// 					continue
// 				}

// 				var targetResponses []pohonkinerja.TargetResponse
// 				for _, target := range targetList {
// 					targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 						Id:              target.Id,
// 						IndikatorId:     target.IndikatorId,
// 						TargetIndikator: target.Target,
// 						SatuanIndikator: target.Satuan,
// 					})
// 				}

// 				indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
// 					Id:            indikator.Id,
// 					IdPokin:       fmt.Sprint(crosscutting.Id),
// 					NamaIndikator: indikator.Indikator,
// 					Target:        targetResponses,
// 				})
// 			}
// 			indikatorMap[crosscutting.Id] = indikatorResponses
// 		}

// 		// Ambil data pelaksana untuk crosscutting
// 		pelaksanaList, err := service.pohonKinerjaOpdRepository.FindPelaksanaPokin(ctx, tx, fmt.Sprint(crosscutting.Id))
// 		if err == nil {
// 			var pelaksanaResponses []pohonkinerja.PelaksanaOpdResponse
// 			for _, pelaksana := range pelaksanaList {
// 				pegawai, err := service.pegawaiRepository.FindById(ctx, tx, pelaksana.PegawaiId)
// 				if err != nil {
// 					continue
// 				}
// 				pelaksanaResponses = append(pelaksanaResponses, pohonkinerja.PelaksanaOpdResponse{
// 					Id:             pelaksana.Id,
// 					PohonKinerjaId: fmt.Sprint(crosscutting.Id),
// 					PegawaiId:      pelaksana.PegawaiId,
// 					NamaPegawai:    pegawai.NamaPegawai,
// 				})
// 			}
// 			pelaksanaMap[crosscutting.Id] = pelaksanaResponses
// 		}

// 		// Jika status disetujui_existing, ambil data pohon kinerja yang di-crosscut
// 		var namaPohon string
// 		var jenisPohon string
// 		var levelPohon int

// 		if crosscutting.Status == "crosscutting_disetujui_existing" && crosscutting.CrosscuttingTo != 0 {
// 			pokinExisting, err := service.pohonKinerjaOpdRepository.FindById(ctx, tx, crosscutting.CrosscuttingTo)
// 			if err == nil {
// 				namaPohon = pokinExisting.NamaPohon
// 				jenisPohon = pokinExisting.JenisPohon
// 				levelPohon = pokinExisting.LevelPohon
// 			}
// 		} else {
// 			namaPohon = crosscutting.NamaPohon
// 			jenisPohon = crosscutting.JenisPohon
// 			levelPohon = crosscutting.LevelPohon
// 		}

// 		crosscuttingResp := pohonkinerja.CrosscuttingOpdResponse{
// 			Id:         crosscutting.Id,
// 			Parent:     pokinId,
// 			NamaPohon:  namaPohon,
// 			JenisPohon: jenisPohon,
// 			LevelPohon: levelPohon,
// 			Keterangan: crosscutting.Keterangan,
// 			Status:     crosscutting.Status,
// 			KodeOpd:    crosscutting.KodeOpd,
// 			NamaOpd:    opd.NamaOpd,
// 			Tahun:      crosscutting.Tahun,
// 			Pelaksana:  pelaksanaMap[crosscutting.Id],
// 			Indikator:  indikatorMap[crosscutting.Id],
// 		}

// 		crosscuttingResponses = append(crosscuttingResponses, crosscuttingResp)
// 	}

// 	return crosscuttingResponses
// }

func (service *PohonKinerjaOpdServiceImpl) DeletePokinPemdaInOpd(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// 1. Cek apakah pohon kinerja dengan ID tersebut ada
	_, err = service.pohonKinerjaOpdRepository.FindById(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("pohon kinerja tidak ditemukan: %v", err)
	}

	// 2. Cek apakah ini adalah pohon kinerja yang di-clone dan dapatkan ID aslinya
	cloneFrom, err := service.pohonKinerjaOpdRepository.CheckCloneFrom(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("gagal memeriksa clone_from: %v", err)
	}

	if cloneFrom == 0 {
		return fmt.Errorf("pohon kinerja ini bukan merupakan hasil clone dari pemda")
	}

	// 3. Hapus data clone saat ini
	err = service.pohonKinerjaOpdRepository.Delete(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus pohon kinerja clone: %v", err)
	}

	// 4. Update status pohon kinerja asli (yang di-clone) menjadi "ditolak"
	err = service.pohonKinerjaOpdRepository.UpdatePokinStatusFromApproved(ctx, tx, cloneFrom)
	if err != nil {
		return fmt.Errorf("gagal mengupdate status pohon kinerja asli: %v", err)
	}

	// 5. Dapatkan dan update status semua child dari pohon kinerja asli
	originalHierarchy, err := service.pohonKinerjaOpdRepository.FindPokinAdminByIdHierarki(ctx, tx, cloneFrom)
	if err != nil {
		return fmt.Errorf("gagal mendapatkan hierarki pohon kinerja asli: %v", err)
	}

	// 6. Update status untuk semua child dari pohon kinerja asli
	for _, originalPokin := range originalHierarchy {
		if originalPokin.Id != cloneFrom { // Skip pohon kinerja utama karena sudah diupdate
			err = service.pohonKinerjaOpdRepository.UpdatePokinStatusFromApproved(ctx, tx, originalPokin.Id)
			if err != nil {
				log.Printf("Warning: gagal mengupdate status child pohon kinerja asli ID %d: %v", originalPokin.Id, err)
				continue // Lanjutkan ke child berikutnya meskipun ada error
			}
		}
	}

	return nil
}

func (service *PohonKinerjaOpdServiceImpl) UpdateParent(ctx context.Context, pohonKinerja pohonkinerja.PohonKinerjaUpdateRequest) (pohonkinerja.PohonKinerjaOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	pokin := domain.PohonKinerja{
		Id:     pohonKinerja.Id,
		Parent: pohonKinerja.Parent,
	}

	pokin, err = service.pohonKinerjaOpdRepository.UpdateParent(ctx, tx, pokin)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}

	return pohonkinerja.PohonKinerjaOpdResponse{
		Id:     pokin.Id,
		Parent: fmt.Sprint(pokin.Parent),
	}, nil
}

func (service *PohonKinerjaOpdServiceImpl) FindidPokinWithAllTema(ctx context.Context, id int) (pohonkinerja.PohonKinerjaAdminResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	pokins, err := service.pohonKinerjaOpdRepository.FindidPokinWithAllTema(ctx, tx, id)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponse{}, err
	}

	// Temukan node target dan level parentnya
	var targetPokin domain.PohonKinerja
	var parentLevel int
	var strategicId int // Untuk menyimpan ID dari node Strategic

	for _, pokin := range pokins {
		if pokin.Id == id {
			targetPokin = pokin
			// Jika level 4, cari level parentnya
			if pokin.LevelPohon == 4 {
				for _, p := range pokins {
					if p.Id == pokin.Parent {
						parentLevel = p.LevelPohon
						break
					}
				}
				strategicId = pokin.Id
			} else if pokin.LevelPohon == 5 || pokin.LevelPohon == 6 {
				// Untuk level 5 dan 6, cari node Strategic (level 4) di ancestors
				for _, p := range pokins {
					if p.LevelPohon == 4 {
						strategicId = p.Id
						// Cari level parent dari Strategic
						for _, grandParent := range pokins {
							if grandParent.Id == p.Parent {
								parentLevel = grandParent.LevelPohon
								break
							}
						}
						break
					}
				}
			}
			break
		}
	}

	// Validasi level
	if targetPokin.LevelPohon < 4 || targetPokin.LevelPohon > 6 {
		return pohonkinerja.PohonKinerjaAdminResponse{}, fmt.Errorf("ID harus merujuk ke level Strategic (4), Tactical (5), atau Operational (6)")
	}

	// Helper functions (sama seperti sebelumnya)
	createIndikatorResponse := func(pokin domain.PohonKinerja) []pohonkinerja.IndikatorResponse {
		var indikators []pohonkinerja.IndikatorResponse
		for _, ind := range pokin.Indikator {
			var targets []pohonkinerja.TargetResponse
			for _, t := range ind.Target {
				targets = append(targets, pohonkinerja.TargetResponse{
					Id:              t.Id,
					IndikatorId:     t.IndikatorId,
					TargetIndikator: t.Target,
					SatuanIndikator: t.Satuan,
				})
			}
			indikators = append(indikators, pohonkinerja.IndikatorResponse{
				Id:            ind.Id,
				IdPokin:       ind.PokinId,
				NamaIndikator: ind.Indikator,
				Target:        targets,
			})
		}
		return indikators
	}

	// Buat responses untuk setiap level yang diperlukan
	var strategicResp pohonkinerja.StrategicResponse
	var tacticalResp *pohonkinerja.TacticalResponse
	var operationalResp *pohonkinerja.OperationalResponse

	// Bangun response berdasarkan level target
	for _, pokin := range pokins {
		if pokin.Id == strategicId {
			// Buat Strategic Response
			strategicResp = pohonkinerja.StrategicResponse{
				Id:         pokin.Id,
				Parent:     pokin.Parent,
				Strategi:   pokin.NamaPohon,
				JenisPohon: pokin.JenisPohon,
				LevelPohon: pokin.LevelPohon,
				Keterangan: pokin.Keterangan,
				Status:     pokin.Status,
				Indikators: createIndikatorResponse(pokin),
				Childs:     []interface{}{},
			}
		} else if pokin.LevelPohon == 5 && targetPokin.LevelPohon >= 5 {
			// Buat Tactical Response
			tacticalResp = &pohonkinerja.TacticalResponse{
				Id:         pokin.Id,
				Parent:     pokin.Parent,
				Strategi:   pokin.NamaPohon,
				JenisPohon: pokin.JenisPohon,
				LevelPohon: pokin.LevelPohon,
				Keterangan: &pokin.Keterangan,
				Status:     pokin.Status,
				Indikators: createIndikatorResponse(pokin),
				Childs:     []interface{}{},
			}
			strategicResp.Childs = append(strategicResp.Childs, tacticalResp)
		} else if pokin.LevelPohon == 6 && targetPokin.LevelPohon == 6 {
			// Buat Operational Response
			operationalResp = &pohonkinerja.OperationalResponse{
				Id:         pokin.Id,
				Parent:     pokin.Parent,
				Strategi:   pokin.NamaPohon,
				JenisPohon: pokin.JenisPohon,
				LevelPohon: pokin.LevelPohon,
				Keterangan: &pokin.Keterangan,
				Status:     pokin.Status,
				Indikators: createIndikatorResponse(pokin),
				Childs:     []interface{}{},
			}
			if tacticalResp != nil {
				tacticalResp.Childs = append(tacticalResp.Childs, operationalResp)
			}
		}
	}

	// Bangun hierarki dari Tematik ke bawah
	var tematikResp pohonkinerja.TematikResponse
	for _, pokin := range pokins {
		if pokin.LevelPohon == 0 { // Tematik
			parentInt := pokin.Parent
			tematikResp = pohonkinerja.TematikResponse{
				Id:         pokin.Id,
				Parent:     &parentInt,
				Tema:       pokin.NamaPohon,
				JenisPohon: pokin.JenisPohon,
				LevelPohon: pokin.LevelPohon,
				Keterangan: pokin.Keterangan,
				Indikators: createIndikatorResponse(pokin),
				Child:      []interface{}{},
			}

			if parentLevel == 0 {
				tematikResp.Child = append(tematikResp.Child, strategicResp)
			}
		} else if pokin.LevelPohon <= parentLevel {
			switch pokin.LevelPohon {
			case 1: // Subtematik
				subtematikResp := pohonkinerja.SubtematikResponse{
					Id:         pokin.Id,
					Parent:     pokin.Parent,
					Tema:       pokin.NamaPohon,
					JenisPohon: pokin.JenisPohon,
					LevelPohon: pokin.LevelPohon,
					Keterangan: pokin.Keterangan,
					Indikators: createIndikatorResponse(pokin),
					Child:      []interface{}{},
				}
				if parentLevel == 1 {
					subtematikResp.Child = append(subtematikResp.Child, strategicResp)
				}
				tematikResp.Child = append(tematikResp.Child, subtematikResp)

			case 2: // SubSubTematik
				subsubtematikResp := pohonkinerja.SubSubTematikResponse{
					Id:         pokin.Id,
					Parent:     pokin.Parent,
					Tema:       pokin.NamaPohon,
					JenisPohon: pokin.JenisPohon,
					LevelPohon: pokin.LevelPohon,
					Keterangan: pokin.Keterangan,
					Indikators: createIndikatorResponse(pokin),
					Child:      []interface{}{},
				}
				if parentLevel == 2 {
					subsubtematikResp.Child = append(subsubtematikResp.Child, strategicResp)
				}
				for i := range tematikResp.Child {
					if sub, ok := tematikResp.Child[i].(pohonkinerja.SubtematikResponse); ok && sub.Id == pokin.Parent {
						sub.Child = append(sub.Child, subsubtematikResp)
						tematikResp.Child[i] = sub
					}
				}

			case 3: // SuperSubTematik
				supersubtematikResp := pohonkinerja.SuperSubTematikResponse{
					Id:         pokin.Id,
					Parent:     pokin.Parent,
					Tema:       pokin.NamaPohon,
					JenisPohon: pokin.JenisPohon,
					LevelPohon: pokin.LevelPohon,
					Keterangan: pokin.Keterangan,
					Indikators: createIndikatorResponse(pokin),
					Childs:     []interface{}{},
				}
				if parentLevel == 3 {
					supersubtematikResp.Childs = append(supersubtematikResp.Childs, strategicResp)
				}
				// Tambahkan ke parent yang sesuai
				for i := range tematikResp.Child {
					if sub, ok := tematikResp.Child[i].(pohonkinerja.SubtematikResponse); ok {
						for j := range sub.Child {
							if subsub, ok := sub.Child[j].(pohonkinerja.SubSubTematikResponse); ok && subsub.Id == pokin.Parent {
								subsub.Child = append(subsub.Child, supersubtematikResp)
								sub.Child[j] = subsub
								tematikResp.Child[i] = sub
							}
						}
					}
				}
			}
		}
	}

	response := pohonkinerja.PohonKinerjaAdminResponse{
		Tahun:   targetPokin.Tahun,
		Tematik: []pohonkinerja.TematikResponse{tematikResp},
	}

	return response, nil
}
func (service *PohonKinerjaOpdServiceImpl) CloneByKodeOpdAndTahun(ctx context.Context, request pohonkinerja.PohonKinerjaCloneRequest) error {
	// Validasi request
	err := service.Validate.Struct(request)
	if err != nil {
		return err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Lakukan cloning
	err = service.pohonKinerjaOpdRepository.ClonePokinOpd(ctx, tx, request.KodeOpd, request.TahunSumber, request.TahunTujuan)
	if err != nil {
		return fmt.Errorf("gagal melakukan cloning: %v", err)
	}

	return nil
}

func (service *PohonKinerjaOpdServiceImpl) CheckPokinExistsByTahun(ctx context.Context, kodeOpd string, tahun string) (bool, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return false, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	exists := service.pohonKinerjaOpdRepository.IsExistsByTahun(ctx, tx, kodeOpd, tahun)

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("gagal melakukan commit transaksi: %v", err)
	}

	return exists, nil
}

func (service *PohonKinerjaOpdServiceImpl) CountPokinPemda(ctx context.Context, kodeOpd, tahun string) (pohonkinerja.CountPokinPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.CountPokinPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return pohonkinerja.CountPokinPemdaResponse{}, errors.New("kode opd tidak ditemukan")
	}

	// Hitung jumlah pokin pemda per level
	countByLevel, err := service.pohonKinerjaOpdRepository.CountPokinPemdaByLevel(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return pohonkinerja.CountPokinPemdaResponse{}, err
	}

	// Siapkan response
	response := pohonkinerja.CountPokinPemdaResponse{
		KodeOpd: kodeOpd,
		NamaOpd: opd.NamaOpd,
		Tahun:   tahun,
	}

	// Definisikan level yang harus selalu ada dalam response
	requiredLevels := []struct {
		Level      int
		JenisPohon string
	}{
		{4, "Strategic"},
		{5, "Tactical"},
		{6, "Operational"},
	}

	var details []pohonkinerja.LevelDetail
	totalPemda := 0

	// Tambahkan semua level yang required
	for _, req := range requiredLevels {
		count := countByLevel[req.Level] // Akan return 0 jika tidak ada data
		details = append(details, pohonkinerja.LevelDetail{
			Level:       req.Level,
			JenisPohon:  req.JenisPohon,
			JumlahPemda: count,
		})
		totalPemda += count
	}

	// Tambahkan level operational-n (> 6) jika ada
	for level, count := range countByLevel {
		if level > 6 {
			details = append(details, pohonkinerja.LevelDetail{
				Level:       level,
				JenisPohon:  fmt.Sprintf("Operational-%d", level-6),
				JumlahPemda: count,
			})
			totalPemda += count
		}
	}

	// Urutkan detail berdasarkan level
	sort.Slice(details, func(i, j int) bool {
		return details[i].Level < details[j].Level
	})

	response.TotalPemda = totalPemda
	response.DetailLevel = details

	return response, nil
}

func (service *PohonKinerjaOpdServiceImpl) FindPokinAtasan(ctx context.Context, id int) (pohonkinerja.PokinAtasanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PokinAtasanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi ID pokin
	err = service.pohonKinerjaOpdRepository.ValidatePokinId(ctx, tx, id)
	if err != nil {
		return pohonkinerja.PokinAtasanResponse{}, err
	}

	// Ambil data pokin atasan dan pegawainya
	pokinAtasan, pegawaiList, err := service.pohonKinerjaOpdRepository.FindPokinAtasan(ctx, tx, id)
	if err != nil {
		return pohonkinerja.PokinAtasanResponse{}, err
	}

	// Transform ke response
	var pegawaiResponses []pohonkinerja.PegawaiResponse
	for _, pegawai := range pegawaiList {
		pegawaiResponses = append(pegawaiResponses, pohonkinerja.PegawaiResponse{
			IdPegawai:   pegawai.Id,
			NipPegawai:  pegawai.Nip,
			NamaPegawai: pegawai.NamaPegawai,
		})
	}

	response := pohonkinerja.PokinAtasanResponse{
		Id:        pokinAtasan.Id,
		NamaPohon: pokinAtasan.NamaPohon,
		Pegawai:   pegawaiResponses,
	}

	return response, nil
}

func (service *PohonKinerjaOpdServiceImpl) ControlPokinOpd(ctx context.Context, kodeOpd, tahun string) (pohonkinerja.ControlPokinOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.ControlPokinOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data per level
	dataPerLevel, err := service.pohonKinerjaOpdRepository.ControlPokinOpdByLevel(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return pohonkinerja.ControlPokinOpdResponse{}, err
	}

	// Cari level maksimum yang ada di data
	maxLevel := 6 // Minimal sampai Operational (level 6)
	for level := range dataPerLevel {
		if level > maxLevel {
			maxLevel = level
		}
	}

	// Map nama level
	levelNames := map[int]string{
		4: "Strategic",
		5: "Tactical",
		6: "Operational",
	}

	// Build response data per level
	var responseData []pohonkinerja.ControlPokinOpdData
	var totalPokin, totalPelaksana, totalPokinAdaPelaksana, totalPokinTanpaPelaksana int
	var totalRencanaKinerja, totalPokinAdaRekin, totalPokinTanpaRekin int // ← TAMBAH VARIABEL BARU

	// Iterasi dari level 4 sampai maxLevel
	for level := 4; level <= maxLevel; level++ {
		var namaLevel string

		// Tentukan nama level
		if level >= 7 {
			// Level 7+ adalah Operational N (Operational 1, 2, 3, dst)
			operationalN := level - 6 // Level 7 = Operational 1, Level 8 = Operational 2, dst
			namaLevel = fmt.Sprintf("Operational %d", operationalN)
		} else {
			namaLevel = levelNames[level]
		}

		if data, exists := dataPerLevel[level]; exists {
			// Hitung persentase pelaksana = (pokin ada pelaksana / total pokin) * 100
			persentasePelaksana := 0.0
			if data.JumlahPokin > 0 {
				persentasePelaksana = (float64(data.JumlahPokinAdaPelaksana) / float64(data.JumlahPokin)) * 100
			}

			// ✅ Hitung persentase cascading = (pokin ada rekin / total pokin) * 100
			persentaseCascading := 0.0
			if data.JumlahPokin > 0 {
				persentaseCascading = (float64(data.JumlahPokinAdaRekin) / float64(data.JumlahPokin)) * 100
			}

			responseData = append(responseData, pohonkinerja.ControlPokinOpdData{
				LevelPohon:                level,
				NamaLevel:                 namaLevel,
				JumlahPokin:               data.JumlahPokin,
				JumlahPelaksana:           data.JumlahPelaksana,
				JumlahPokinAdaPelaksana:   data.JumlahPokinAdaPelaksana,
				JumlahPokinTanpaPelaksana: data.JumlahPokinTanpaPelaksana,
				JumlahRencanaKinerja:      data.JumlahRencanaKinerja,
				JumlahPokinAdaRekin:       data.JumlahPokinAdaRekin,
				JumlahPokinTanpaRekin:     data.JumlahPokinTanpaRekin,
				Persentase:                fmt.Sprintf("%.0f%%", persentasePelaksana),
				PersentaseCascading:       fmt.Sprintf("%.0f%%", persentaseCascading),
			})

			// Akumulasi total
			totalPokin += data.JumlahPokin
			totalPelaksana += data.JumlahPelaksana
			totalPokinAdaPelaksana += data.JumlahPokinAdaPelaksana
			totalPokinTanpaPelaksana += data.JumlahPokinTanpaPelaksana
			totalRencanaKinerja += data.JumlahRencanaKinerja
			totalPokinAdaRekin += data.JumlahPokinAdaRekin
			totalPokinTanpaRekin += data.JumlahPokinTanpaRekin
		}
	}

	// Hitung persentase total pelaksana
	persentaseTotalPelaksana := 0.0
	if totalPokin > 0 {
		persentaseTotalPelaksana = (float64(totalPokinAdaPelaksana) / float64(totalPokin)) * 100
	}

	// ✅ Hitung persentase total cascading
	persentaseTotalCascading := 0.0
	if totalPokin > 0 {
		persentaseTotalCascading = (float64(totalPokinAdaRekin) / float64(totalPokin)) * 100
	}

	response := pohonkinerja.ControlPokinOpdResponse{
		Data: responseData,
		Total: pohonkinerja.ControlPokinOpdTotal{
			TotalPokin:               totalPokin,
			TotalPelaksana:           totalPelaksana,
			TotalPokinAdaPelaksana:   totalPokinAdaPelaksana,
			TotalPokinTanpaPelaksana: totalPokinTanpaPelaksana,
			TotalRencanaKinerja:      totalRencanaKinerja,
			TotalPokinAdaRekin:       totalPokinAdaRekin,
			TotalPokinTanpaRekin:     totalPokinTanpaRekin,
			Persentase:               fmt.Sprintf("%.0f%%", persentaseTotalPelaksana),
			PersentaseCascading:      fmt.Sprintf("%.0f%%", persentaseTotalCascading),
		},
	}

	return response, nil
}

type CachedLeaderboard struct {
	Data      []pohonkinerja.LeaderboardPokinResponse
	Timestamp time.Time
}

var (
	leaderboardCache     = make(map[string]CachedLeaderboard)
	leaderboardCacheLock sync.RWMutex
	cacheExpiry          = 5 * time.Minute
)

func (service *PohonKinerjaOpdServiceImpl) LeaderboardPokinOpd(ctx context.Context, tahun string) ([]pohonkinerja.LeaderboardPokinResponse, error) {
	cacheKey := fmt.Sprintf("leaderboard_%s", tahun)

	// Cek cache terlebih dahulu
	leaderboardCacheLock.RLock()
	if cached, exists := leaderboardCache[cacheKey]; exists {
		if time.Since(cached.Timestamp) < cacheExpiry {
			leaderboardCacheLock.RUnlock()
			return cached.Data, nil
		}
	}
	leaderboardCacheLock.RUnlock()

	// Jika cache tidak ada atau expired, ambil dari database
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	leaderboardData, err := service.pohonKinerjaOpdRepository.LeaderboardPokinOpd(ctx, tx, tahun)
	if err != nil {
		return nil, err
	}

	var response []pohonkinerja.LeaderboardPokinResponse
	for _, data := range leaderboardData {
		var tematikItems []pohonkinerja.LeaderboardTematikItem
		for _, name := range data.TematikNames {
			tematikItems = append(tematikItems, pohonkinerja.LeaderboardTematikItem{
				Nama: name,
			})
		}

		response = append(response, pohonkinerja.LeaderboardPokinResponse{
			KodeOpd:             data.KodeOpd,
			NamaOpd:             data.NamaOpd,
			Tematik:             tematikItems,
			PersentaseCascading: fmt.Sprintf("%.0f%%", data.PersentaseCascading),
		})
	}

	//  Simpan ke cache
	leaderboardCacheLock.Lock()
	leaderboardCache[cacheKey] = CachedLeaderboard{
		Data:      response,
		Timestamp: time.Now(),
	}
	leaderboardCacheLock.Unlock()

	return response, nil
}

// Tambahkan fungsi untuk clear cache (dipanggil saat ada update data)
func (service *PohonKinerjaOpdServiceImpl) ClearLeaderboardCache(tahun string) {
	cacheKey := fmt.Sprintf("leaderboard_%s", tahun)
	leaderboardCacheLock.Lock()
	delete(leaderboardCache, cacheKey)
	leaderboardCacheLock.Unlock()
}
