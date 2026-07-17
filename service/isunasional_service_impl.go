package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/isunasional"
	"ekak_kabupaten_madiun/repository"

	"github.com/go-playground/validator/v10"
)

type IsuNasionalServiceImpl struct {
	IsuNasionalRepository repository.IsuNasionalRepository
	DB                    *sql.DB
	Validate              *validator.Validate
}

func NewIsuNasionalServiceImpl(isunasionalRepository repository.IsuNasionalRepository, db *sql.DB, validate *validator.Validate) *IsuNasionalServiceImpl {
	return &IsuNasionalServiceImpl{
		IsuNasionalRepository: isunasionalRepository,
		DB:                db,
		Validate:          validate,
	}
}

func (service *IsuNasionalServiceImpl) Create(ctx context.Context, request isunasional.IsuNasionalRequest) (isunasional.IsuNasionalResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return isunasional.IsuNasionalResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return isunasional.IsuNasionalResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	data := domain.IsuNasional{
		KodeBidangUrusan: request.KodeBidangUrusan,
		KodeOpd:          request.KodeOpd,
		Isu:              request.Isu,
		Tahun:            request.Tahun,
	}

	result, err := service.IsuNasionalRepository.Create(ctx, tx, data)
	if err != nil {
		return isunasional.IsuNasionalResponse{}, err
	}

	return isunasional.IsuNasionalResponse{
		ID:                 result.ID,
		KodeBidangUrusan:   result.KodeBidangUrusan,
		KodeOpd:            result.KodeOpd,
		Isu:                result.Isu,
		Tahun:              result.Tahun,
	}, nil
}

func (service *IsuNasionalServiceImpl) Update(
	ctx context.Context,
	request isunasional.IsuNasionalUpdateRequest,
) (isunasional.IsuNasionalResponse, error) {

	err := service.Validate.Struct(request)
	if err != nil {
		return isunasional.IsuNasionalResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return isunasional.IsuNasionalResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// validasi exists
	_, err = service.IsuNasionalRepository.FindById(ctx, tx, request.ID)
	if err != nil {
		return isunasional.IsuNasionalResponse{}, err
	}


	data := domain.IsuNasional{
		ID:                 request.ID,
		KodeBidangUrusan:   request.KodeBidangUrusan,
		KodeOpd:            request.KodeOpd,
		Isu:                request.Isu,
		Tahun:              request.Tahun,
	}

	result, err := service.IsuNasionalRepository.Update(ctx, tx, data)
	if err != nil {
		return isunasional.IsuNasionalResponse{}, err
	}

	return isunasional.IsuNasionalResponse{
		ID:                 result.ID,
		KodeBidangUrusan:   result.KodeBidangUrusan,
		KodeOpd:            result.KodeOpd,
		Isu:                result.Isu,
		Tahun:              result.Tahun,
	}, nil
}

func (service *IsuNasionalServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi data exists
	_, err = service.IsuNasionalRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}

	return service.IsuNasionalRepository.Delete(ctx, tx, id)
}