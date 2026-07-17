package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/isuglobal"
	"ekak_kabupaten_madiun/repository"

	"github.com/go-playground/validator/v10"
)

type IsuGlobalServiceImpl struct {
	IsuGlobalRepository repository.IsuGlobalRepository
	DB            *sql.DB
	Validate      *validator.Validate
}

func NewIsuGlobalServiceImpl(isuglobalRepository repository.IsuGlobalRepository, db *sql.DB, validate *validator.Validate) *IsuGlobalServiceImpl {
	return &IsuGlobalServiceImpl{
		IsuGlobalRepository: isuglobalRepository,
		DB:            db,
		Validate:      validate,
	}
}

func (service *IsuGlobalServiceImpl) Create(ctx context.Context, request isuglobal.IsuGlobalRequest) (isuglobal.IsuGlobalResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return isuglobal.IsuGlobalResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return isuglobal.IsuGlobalResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	data := domain.IsuGlobal{
		KodeBidangUrusan: request.KodeBidangUrusan,
		KodeOpd:          request.KodeOpd,
		Isu:              request.Isu,
		Tahun:            request.Tahun,
	}

	result, err := service.IsuGlobalRepository.Create(ctx, tx, data)
	if err != nil {
		return isuglobal.IsuGlobalResponse{}, err
	}

	return isuglobal.IsuGlobalResponse{
		ID:                 result.ID,
		KodeBidangUrusan:   result.KodeBidangUrusan,
		KodeOpd:            result.KodeOpd,
		Isu:                result.Isu,
		Tahun:              result.Tahun,
	}, nil
}

func (service *IsuGlobalServiceImpl) Update(
	ctx context.Context,
	request isuglobal.IsuGlobalUpdateRequest,
) (isuglobal.IsuGlobalResponse, error) {

	err := service.Validate.Struct(request)
	if err != nil {
		return isuglobal.IsuGlobalResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return isuglobal.IsuGlobalResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// validasi exists
	_, err = service.IsuGlobalRepository.FindById(ctx, tx, request.ID)
	if err != nil {
		return isuglobal.IsuGlobalResponse{}, err
	}


	data := domain.IsuGlobal{
		ID:                 request.ID,
		KodeBidangUrusan:   request.KodeBidangUrusan,
		KodeOpd:            request.KodeOpd,
		Isu:                request.Isu,
		Tahun:              request.Tahun,
	}

	result, err := service.IsuGlobalRepository.Update(ctx, tx, data)
	if err != nil {
		return isuglobal.IsuGlobalResponse{}, err
	}

	return isuglobal.IsuGlobalResponse{
		ID:                 result.ID,
		KodeBidangUrusan:   result.KodeBidangUrusan,
		KodeOpd:            result.KodeOpd,
		Isu:                result.Isu,
		Tahun:              result.Tahun,
	}, nil
}

func (service *IsuGlobalServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi data exists
	_, err = service.IsuGlobalRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}

	return service.IsuGlobalRepository.Delete(ctx, tx, id)
}