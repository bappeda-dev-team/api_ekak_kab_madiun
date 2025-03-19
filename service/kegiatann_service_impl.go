package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"ekak_kabupaten_madiun/model/web/kegiatan"
	"ekak_kabupaten_madiun/repository"
	"fmt"

	"github.com/google/uuid"
)

type KegiatanServiceImpl struct {
	KegiatanRepository repository.KegiatanRepository
	DB                 *sql.DB
}

func NewKegiatanServiceImpl(kegiatanRepository repository.KegiatanRepository, DB *sql.DB) *KegiatanServiceImpl {
	return &KegiatanServiceImpl{
		KegiatanRepository: kegiatanRepository,
		DB:                 DB,
	}
}

func (service *KegiatanServiceImpl) Create(ctx context.Context, request kegiatan.KegiatanCreateRequest) (kegiatan.KegiatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return kegiatan.KegiatanResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	uuidKegiatan := fmt.Sprintf("KGT-%s", uuid.New().String()[:5])

	kegiatans := domainmaster.Kegiatan{
		Id:           uuidKegiatan,
		KodeKegiatan: request.KodeKegiatan,
		NamaKegiatan: request.NamaKegiatan,
	}

	var indikators []domain.Indikator
	for _, indikatorRequest := range request.Indikator {
		uuidIndikator := fmt.Sprintf("IND-KGT-%s", uuid.New().String()[:5])
		indikator := domain.Indikator{
			Id:         uuidIndikator,
			KegiatanId: kegiatans.Id,
			Indikator:  indikatorRequest.Indikator,
			Tahun:      indikatorRequest.Tahun,
		}

		var targets []domain.Target
		for _, targetRequest := range indikatorRequest.Target {
			uuidTarget := fmt.Sprintf("TRGT-KGT-%s", uuid.New().String()[:5])
			target := domain.Target{
				Id:          uuidTarget,
				IndikatorId: indikator.Id,
				Target:      targetRequest.Target,
				Satuan:      targetRequest.Satuan,
				Tahun:       targetRequest.Tahun,
			}
			targets = append(targets, target)
		}
		indikator.Target = targets
		indikators = append(indikators, indikator)
	}
	kegiatans.Indikator = indikators

	result, err := service.KegiatanRepository.Create(ctx, tx, kegiatans)
	if err != nil {
		return kegiatan.KegiatanResponse{}, fmt.Errorf("gagal membuat kegiatan: %v", err)
	}

	// Konversi langsung ke KegiatanResponse
	var indikatorResponses []kegiatan.IndikatorResponse
	for _, ind := range result.Indikator {
		var targetResponses []kegiatan.TargetResponse
		for _, t := range ind.Target {
			targetResponses = append(targetResponses, kegiatan.TargetResponse{
				Id:          t.Id,
				IndikatorId: t.IndikatorId,
				Target:      t.Target,
				Satuan:      t.Satuan,
				Tahun:       t.Tahun,
			})
		}

		indikatorResponses = append(indikatorResponses, kegiatan.IndikatorResponse{
			Id:         ind.Id,
			KegiatanId: ind.KegiatanId,
			Indikator:  ind.Indikator,
			Tahun:      ind.Tahun,
			Target:     targetResponses,
		})
	}

	return kegiatan.KegiatanResponse{
		Id:           result.Id,
		KodeKegiatan: result.KodeKegiatan,
		NamaKegiatan: result.NamaKegiatan,
		Indikator:    indikatorResponses,
	}, nil
}

func (service *KegiatanServiceImpl) Update(ctx context.Context, request kegiatan.KegiatanUpdateRequest) (kegiatan.KegiatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return kegiatan.KegiatanResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah kegiatan exists
	_, err = service.KegiatanRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return kegiatan.KegiatanResponse{}, fmt.Errorf("kegiatan tidak ditemukan: %v", err)
	}

	// Update data kegiatan
	kegiatans := domainmaster.Kegiatan{
		Id:           request.Id,
		KodeKegiatan: request.KodeKegiatan,
		NamaKegiatan: request.NamaKegiatan,
	}

	var indikators []domain.Indikator
	for _, indikator := range request.Indikator {
		// Generate ID baru jika indikator baru
		indikatorId := indikator.Id
		if indikatorId == "" {
			indikatorId = fmt.Sprintf("IND-KGT-%s", uuid.New().String()[:5])
		}

		var targets []domain.Target
		for _, target := range indikator.Target {
			// Generate ID baru jika target baru
			targetId := target.Id
			if targetId == "" {
				targetId = fmt.Sprintf("TRGT-KGT-%s", uuid.New().String()[:5])
			}

			targetDomain := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorId,
				Target:      target.Target,
				Satuan:      target.Satuan,
				Tahun:       target.Tahun,
			}
			targets = append(targets, targetDomain)
		}

		indikatorDomain := domain.Indikator{
			Id:         indikatorId,
			KegiatanId: request.Id,
			Indikator:  indikator.Indikator,
			Tahun:      indikator.Tahun,
			Target:     targets,
		}
		indikators = append(indikators, indikatorDomain)
	}
	kegiatans.Indikator = indikators

	result, err := service.KegiatanRepository.Update(ctx, tx, kegiatans)
	if err != nil {
		return kegiatan.KegiatanResponse{}, fmt.Errorf("gagal mengupdate kegiatan: %v", err)
	}

	// Konversi langsung ke KegiatanResponse
	var indikatorResponses []kegiatan.IndikatorResponse
	for _, ind := range result.Indikator {
		var targetResponses []kegiatan.TargetResponse
		for _, t := range ind.Target {
			targetResponses = append(targetResponses, kegiatan.TargetResponse{
				Id:          t.Id,
				IndikatorId: t.IndikatorId,
				Target:      t.Target,
				Satuan:      t.Satuan,
				Tahun:       t.Tahun,
			})
		}

		indikatorResponses = append(indikatorResponses, kegiatan.IndikatorResponse{
			Id:         ind.Id,
			KegiatanId: ind.KegiatanId,
			Indikator:  ind.Indikator,
			Tahun:      ind.Tahun,
			Target:     targetResponses,
		})
	}

	return kegiatan.KegiatanResponse{
		Id:           result.Id,
		KodeKegiatan: result.KodeKegiatan,
		NamaKegiatan: result.NamaKegiatan,
		Indikator:    indikatorResponses,
	}, nil
}

func (service *KegiatanServiceImpl) FindById(ctx context.Context, kegiatanId string) (kegiatan.KegiatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return kegiatan.KegiatanResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Mengambil data kegiatan
	result, err := service.KegiatanRepository.FindById(ctx, tx, kegiatanId)
	if err != nil {
		return kegiatan.KegiatanResponse{}, fmt.Errorf("gagal mengambil data kegiatan: %v", err)
	}

	// Mengambil semua indikator untuk kegiatan ini
	indikators, err := service.KegiatanRepository.FindIndikatorByKegiatanId(ctx, tx, result.Id)
	if err != nil {
		return kegiatan.KegiatanResponse{}, fmt.Errorf("gagal mengambil data indikator: %v", err)
	}

	var indikatorResponses []kegiatan.IndikatorResponse
	for _, indikator := range indikators {
		// Mengambil semua target untuk setiap indikator
		targets, err := service.KegiatanRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
		if err != nil {
			return kegiatan.KegiatanResponse{}, fmt.Errorf("gagal mengambil data target: %v", err)
		}

		var targetResponses []kegiatan.TargetResponse
		for _, target := range targets {
			targetResponse := kegiatan.TargetResponse{
				Id:          target.Id,
				IndikatorId: target.IndikatorId,
				Target:      target.Target,
				Satuan:      target.Satuan,
				Tahun:       target.Tahun,
			}
			targetResponses = append(targetResponses, targetResponse)
		}

		indikatorResponse := kegiatan.IndikatorResponse{
			Id:         indikator.Id,
			KegiatanId: indikator.KegiatanId,
			Indikator:  indikator.Indikator,
			Tahun:      indikator.Tahun,
			Target:     targetResponses,
		}
		indikatorResponses = append(indikatorResponses, indikatorResponse)
	}

	return kegiatan.KegiatanResponse{
		Id:           result.Id,
		KodeKegiatan: result.KodeKegiatan,
		NamaKegiatan: result.NamaKegiatan,
		Indikator:    indikatorResponses,
	}, nil
}

func (service *KegiatanServiceImpl) FindAll(ctx context.Context) ([]kegiatan.KegiatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Mengambil semua kegiatan
	results, err := service.KegiatanRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data kegiatan: %v", err)
	}

	var kegiatanResponses []kegiatan.KegiatanResponse

	for _, keg := range results {
		// Mengambil semua indikator untuk kegiatan ini
		indikators, err := service.KegiatanRepository.FindIndikatorByKegiatanId(ctx, tx, keg.Id)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil data indikator untuk kegiatan %s: %v", keg.Id, err)
		}

		var indikatorResponses []kegiatan.IndikatorResponse
		for _, indikator := range indikators {
			// Mengambil semua target untuk setiap indikator
			targets, err := service.KegiatanRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
			if err != nil {
				return nil, fmt.Errorf("gagal mengambil data target untuk indikator %s: %v", indikator.Id, err)
			}

			var targetResponses []kegiatan.TargetResponse
			for _, target := range targets {
				targetResponse := kegiatan.TargetResponse{
					Id:          target.Id,
					IndikatorId: target.IndikatorId,
					Target:      target.Target,
					Satuan:      target.Satuan,
					Tahun:       target.Tahun,
				}
				targetResponses = append(targetResponses, targetResponse)
			}

			indikatorResponse := kegiatan.IndikatorResponse{
				Id:         indikator.Id,
				KegiatanId: indikator.KegiatanId,
				Indikator:  indikator.Indikator,
				Tahun:      indikator.Tahun,
				Target:     targetResponses,
			}
			indikatorResponses = append(indikatorResponses, indikatorResponse)
		}

		kegiatanResponse := kegiatan.KegiatanResponse{
			Id:           keg.Id,
			KodeKegiatan: keg.KodeKegiatan,
			NamaKegiatan: keg.NamaKegiatan,
			Indikator:    indikatorResponses,
		}
		kegiatanResponses = append(kegiatanResponses, kegiatanResponse)
	}

	return kegiatanResponses, nil
}

func (service *KegiatanServiceImpl) Delete(ctx context.Context, kegiatanId string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	err = service.KegiatanRepository.Delete(ctx, tx, kegiatanId)
	if err != nil {
		return err
	}

	return nil
}
