package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/nspkopd"
	"ekak_kabupaten_madiun/repository"

	"github.com/go-playground/validator/v10"
)

type NspkOpdServiceImpl struct {
	NspkOpdRepository repository.NspkOpdRepository
	DB                *sql.DB
	Validate          *validator.Validate
}

func NewNspkOpdServiceImpl(nspkopdRepository repository.NspkOpdRepository, db *sql.DB, validate *validator.Validate) *NspkOpdServiceImpl {
	return &NspkOpdServiceImpl{
		NspkOpdRepository: nspkopdRepository,
		DB:             db,
		Validate:       validate,
	}
}

func (service *NspkOpdServiceImpl) Create(ctx context.Context, request nspkopd.NspkRequest) (nspkopd.NspkResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return nspkopd.NspkResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return nspkopd.NspkResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	data := domain.NspkOpd{
		KodeOpd: 		request.KodeOpd,
		IdNspk:    		request.IdNspk,
		IdTujuanOpd:    request.IdTujuanOpd,
		IdSasaranOpd:   request.IdSasaranOpd,
		Tahun:   		request.Tahun,
	}

	result, err := service.NspkOpdRepository.Create(ctx, tx, data)
	if err != nil {
		return nspkopd.NspkResponse{}, err
	}

	return nspkopd.NspkResponse{
		ID:     		result.ID,
		KodeOpd: 		result.KodeOpd,
		IdNspk:    		request.IdNspk,
		IdTujuanOpd:    request.IdTujuanOpd,
		IdSasaranOpd:   request.IdSasaranOpd,
		Tahun:   		result.Tahun,
	}, nil
}

func (service *NspkOpdServiceImpl) Update(ctx context.Context, request nspkopd.NspkUpdateRequest) (nspkopd.NspkResponse, error) {

	err := service.Validate.Struct(request)
	if err != nil {
		return nspkopd.NspkResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return nspkopd.NspkResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// validasi exists
	_, err = service.NspkOpdRepository.FindById(ctx, tx, request.ID)
	if err != nil {
		return nspkopd.NspkResponse{}, err
	}

	data := domain.NspkOpd{
		ID:      		request.ID,
		KodeOpd: 		request.KodeOpd,
		IdNspk:    		request.IdNspk,
		IdTujuanOpd:    request.IdTujuanOpd,
		IdSasaranOpd:   request.IdSasaranOpd,
		Tahun:   		request.Tahun,
	}

	result, err := service.NspkOpdRepository.Update(ctx, tx, data)
	if err != nil {
		return nspkopd.NspkResponse{}, err
	}

	return nspkopd.NspkResponse{
		ID:      		result.ID,
		KodeOpd: 		result.KodeOpd,
		IdNspk:    		request.IdNspk,
		IdTujuanOpd:    request.IdTujuanOpd,
		IdSasaranOpd:   request.IdSasaranOpd,
		Tahun:   		result.Tahun,
	}, nil
}

func (service *NspkOpdServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi data exists
	_, err = service.NspkOpdRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}

	return service.NspkOpdRepository.Delete(ctx, tx, id)
}

func (service *NspkOpdServiceImpl) FindAll(ctx context.Context, kodeOpd string) ([]nspkopd.NspkFullResponse, error) {

	tx, err := service.DB.Begin()
	if err != nil {
		return []nspkopd.NspkFullResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data IKK
	results, err := service.NspkOpdRepository.FindAll(ctx, tx, kodeOpd)
	if err != nil {
		return []nspkopd.NspkFullResponse{}, err
	}

	var responses []nspkopd.NspkFullResponse
	for _, result := range results {
		responses = append(responses, nspkopd.NspkFullResponse{
			ID:      		result.ID,
			KodeOpd: 		result.KodeOpd,
			NamaOpd: 		result.NamaOpd,
			IdNspk: 		result.IdNspk,
			Nspk:    		result.NSPK,
			IdTujuanOpd:    result.IdTujuanOpd,
			TujuanOpd:      result.TujuanOpd,
			IdSasaranOpd:   result.IdSasaranOpd,
			SasaranOpd:     result.SasaranOpd,
			Tahun:   		result.Tahun,
		})
	}

	return responses, nil
}
