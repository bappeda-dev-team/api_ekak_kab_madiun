package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/kelompokanggarans"
	"ekak_kabupaten_madiun/repository"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type KelompokAnggaranServiceImpl struct {
	KelompokAnggaranRepository repository.KelompokAnggaranRepository
	DB                         *sql.DB
	Validate                   *validator.Validate
}

func NewKelompokAnggaranServiceImpl(kelompokAnggaranRepository repository.KelompokAnggaranRepository, DB *sql.DB, validate *validator.Validate) *KelompokAnggaranServiceImpl {
	return &KelompokAnggaranServiceImpl{
		KelompokAnggaranRepository: kelompokAnggaranRepository,
		DB:                         DB,
		Validate:                   validate,
	}
}

func (service *KelompokAnggaranServiceImpl) Create(ctx context.Context, request kelompokanggarans.KelompokAnggaranCreateRequest) (kelompokanggarans.KelompokAnggaranResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return kelompokanggarans.KelompokAnggaranResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return kelompokanggarans.KelompokAnggaranResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	kodeAnggaran := fmt.Sprintf("%v_%v", request.Tahun, request.Kelompok)

	kelompokAnggaran := domain.KelompokAnggaran{
		Tahun:        request.Tahun,
		Kelompok:     request.Kelompok,
		KodeKelompok: kodeAnggaran,
	}

	result, err := service.KelompokAnggaranRepository.Create(ctx, tx, kelompokAnggaran)
	if err != nil {
		return kelompokanggarans.KelompokAnggaranResponse{}, err
	}

	return kelompokanggarans.KelompokAnggaranResponse{
		Id:           result.Id,
		Tahun:        result.Tahun,
		Kelompok:     result.Kelompok,
		KodeKelompok: result.KodeKelompok,
	}, nil
}

func (service *KelompokAnggaranServiceImpl) Update(ctx context.Context, request kelompokanggarans.KelompokAnggaranUpdateRequest) (kelompokanggarans.KelompokAnggaranResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return kelompokanggarans.KelompokAnggaranResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return kelompokanggarans.KelompokAnggaranResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	kelompokAnggaran := domain.KelompokAnggaran{
		Id:           request.Id,
		Tahun:        request.Tahun,
		Kelompok:     request.Kelompok,
		KodeKelompok: request.KodeKelompok,
	}

	result, err := service.KelompokAnggaranRepository.Update(ctx, tx, kelompokAnggaran)
	if err != nil {
		return kelompokanggarans.KelompokAnggaranResponse{}, err
	}

	return kelompokanggarans.KelompokAnggaranResponse{
		Id:           result.Id,
		Tahun:        result.Tahun,
		Kelompok:     result.Kelompok,
		KodeKelompok: result.KodeKelompok,
	}, nil
}

func (service *KelompokAnggaranServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	return service.KelompokAnggaranRepository.Delete(ctx, tx, id)
}

func (service *KelompokAnggaranServiceImpl) FindById(ctx context.Context, id string) (kelompokanggarans.KelompokAnggaranResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return kelompokanggarans.KelompokAnggaranResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	result, err := service.KelompokAnggaranRepository.FindById(ctx, tx, id)
	if err != nil {
		return kelompokanggarans.KelompokAnggaranResponse{}, err
	}

	return kelompokanggarans.KelompokAnggaranResponse{
		Id:           result.Id,
		Tahun:        result.Tahun,
		Kelompok:     result.Kelompok,
		KodeKelompok: result.KodeKelompok,
	}, nil
}

func (service *KelompokAnggaranServiceImpl) FindAll(ctx context.Context) ([]kelompokanggarans.KelompokAnggaranResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	results := service.KelompokAnggaranRepository.FindAll(ctx, tx)

	var kelompokAnggarans []kelompokanggarans.KelompokAnggaranResponse
	for _, result := range results {
		tahunView := result.Tahun
		if result.Kelompok != "murni" {
			tahunView = fmt.Sprintf("%s %s", result.Tahun, result.Kelompok)
		}
		kelompokAnggarans = append(kelompokAnggarans, kelompokanggarans.KelompokAnggaranResponse{
			Id:           result.Id,
			Tahun:        result.Tahun,
			Kelompok:     result.Kelompok,
			KodeKelompok: result.KodeKelompok,
			TahunView:    tahunView,
		})
	}

	return kelompokAnggarans, nil
}
