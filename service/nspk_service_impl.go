package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/masternspk"
	"ekak_kabupaten_madiun/repository"

	"github.com/go-playground/validator/v10"
)

type NspkServiceImpl struct {
	NspkRepository repository.NspkRepository
	DB             *sql.DB
	Validate       *validator.Validate
}

func NewNspkServiceImpl(nspkRepository repository.NspkRepository, db *sql.DB, validate *validator.Validate) *NspkServiceImpl {
	return &NspkServiceImpl{
		NspkRepository: nspkRepository,
		DB:             db,
		Validate:       validate,
	}
}

func (service *NspkServiceImpl) Create(ctx context.Context, request masternspk.NspkRequest) (masternspk.NspkResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return masternspk.NspkResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return masternspk.NspkResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	data := domain.MasterNSPK{
		KodeOpd:          request.KodeOpd,
		NSPK:             request.Nspk,
		Tahun:            request.Tahun,
	}

	result, err := service.NspkRepository.Create(ctx, tx, data)
	if err != nil {
		return masternspk.NspkResponse{}, err
	}

	return masternspk.NspkResponse{
		ID:               result.ID,
		KodeOpd:          result.KodeOpd,
		Nspk:             result.NSPK,
		Tahun:            result.Tahun,
	}, nil
}

func (service *NspkServiceImpl) Update(ctx context.Context, request masternspk.NspkUpdateRequest) (masternspk.NspkResponse, error) {

	err := service.Validate.Struct(request)
	if err != nil {
		return masternspk.NspkResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return masternspk.NspkResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// validasi exists
	_, err = service.NspkRepository.FindById(ctx, tx, request.ID)
	if err != nil {
		return masternspk.NspkResponse{}, err
	}

	data := domain.MasterNSPK{
		ID:               request.ID,
		KodeOpd:          request.KodeOpd,
		NSPK:             request.Nspk,
		Tahun:            request.Tahun,
	}

	result, err := service.NspkRepository.Update(ctx, tx, data)
	if err != nil {
		return masternspk.NspkResponse{}, err
	}

	return masternspk.NspkResponse{
		ID:               result.ID,
		KodeOpd:          result.KodeOpd,
		Nspk:             result.NSPK,
		Tahun:            result.Tahun,
	}, nil
}

func (service *NspkServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi data exists
	_, err = service.NspkRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}

	return service.NspkRepository.Delete(ctx, tx, id)
}

func (service *NspkServiceImpl) FindAll(ctx context.Context, kodeOpd string) ([]masternspk.NspkFullResponse, error) {

	tx, err := service.DB.Begin()
	if err != nil {
		return []masternspk.NspkFullResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data IKK
	results, err := service.NspkRepository.FindAll(ctx, tx, kodeOpd)
	if err != nil {
		return []masternspk.NspkFullResponse{}, err
	}

	var responses []masternspk.NspkFullResponse
	for _, result := range results {
		responses = append(responses, masternspk.NspkFullResponse{
			ID:                        result.ID,
			KodeOpd:               result.KodeOpd,
			NamaOpd:       result.NamaOpd,
			Nspk: result.NSPK,
			Tahun:                result.Tahun,
		})
	}

	return responses, nil
}

