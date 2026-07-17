package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/isuregional"
	"ekak_kabupaten_madiun/repository"

	"github.com/go-playground/validator/v10"
)

type IsuRegionalServiceImpl struct {
	IsuRegionalRepository repository.IsuRegionalRepository
	DB                    *sql.DB
	Validate              *validator.Validate
}

func NewIsuRegionalServiceImpl(isuregionalRepository repository.IsuRegionalRepository, db *sql.DB, validate *validator.Validate) *IsuRegionalServiceImpl {
	return &IsuRegionalServiceImpl{
		IsuRegionalRepository: isuregionalRepository,
		DB:                    db,
		Validate:              validate,
	}
}

func (service *IsuRegionalServiceImpl) Create(ctx context.Context, request isuregional.IsuRegionalRequest) (isuregional.IsuRegionalResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return isuregional.IsuRegionalResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return isuregional.IsuRegionalResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	data := domain.IsuRegional{
		KodeBidangUrusan: request.KodeBidangUrusan,
		KodeOpd:          request.KodeOpd,
		Isu:              request.Isu,
		Tahun:            request.Tahun,
	}

	result, err := service.IsuRegionalRepository.Create(ctx, tx, data)
	if err != nil {
		return isuregional.IsuRegionalResponse{}, err
	}

	return isuregional.IsuRegionalResponse{
		ID:                 result.ID,
		KodeBidangUrusan:   result.KodeBidangUrusan,
		KodeOpd:            result.KodeOpd,
		Isu:                result.Isu,
		Tahun:              result.Tahun,
	}, nil
}

func (service *IsuRegionalServiceImpl) Update(
	ctx context.Context,
	request isuregional.IsuRegionalUpdateRequest,
) (isuregional.IsuRegionalResponse, error) {

	err := service.Validate.Struct(request)
	if err != nil {
		return isuregional.IsuRegionalResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return isuregional.IsuRegionalResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// validasi exists
	_, err = service.IsuRegionalRepository.FindById(ctx, tx, request.ID)
	if err != nil {
		return isuregional.IsuRegionalResponse{}, err
	}


	data := domain.IsuRegional{
		ID:                 request.ID,
		KodeBidangUrusan:   request.KodeBidangUrusan,
		KodeOpd:            request.KodeOpd,
		Isu:                request.Isu,
		Tahun:              request.Tahun,
	}

	result, err := service.IsuRegionalRepository.Update(ctx, tx, data)
	if err != nil {
		return isuregional.IsuRegionalResponse{}, err
	}

	return isuregional.IsuRegionalResponse{
		ID:                 result.ID,
		KodeBidangUrusan:   result.KodeBidangUrusan,
		KodeOpd:            result.KodeOpd,
		Isu:                result.Isu,
		Tahun:              result.Tahun,
	}, nil
}

func (service *IsuRegionalServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi data exists
	_, err = service.IsuRegionalRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}

	return service.IsuRegionalRepository.Delete(ctx, tx, id)
}