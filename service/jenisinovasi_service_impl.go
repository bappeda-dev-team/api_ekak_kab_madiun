package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/jenisinovasi"
	"ekak_kabupaten_madiun/repository"

	"github.com/go-playground/validator/v10"
)

type JenisInovasiServiceImpl struct {
	JenisInovasiRepository repository.JenisInovasiRepository
	DB             		   *sql.DB
	Validate       		   *validator.Validate
}

func NewJenisInovasiServiceImpl(jenisinovasiRepository repository.JenisInovasiRepository, db *sql.DB, validate *validator.Validate) *JenisInovasiServiceImpl {
	return &JenisInovasiServiceImpl{
		JenisInovasiRepository: jenisinovasiRepository,
		DB:             		db,
		Validate:       		validate,
	}
}

func (service *JenisInovasiServiceImpl) Create(ctx context.Context, request jenisinovasi.JenisInovasiRequest) (jenisinovasi.JenisInovasiResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return jenisinovasi.JenisInovasiResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return jenisinovasi.JenisInovasiResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	data := domain.JenisInovasi{
		Jenis:          request.Jenis,
	}

	result, err := service.JenisInovasiRepository.Create(ctx, tx, data)
	if err != nil {
		return jenisinovasi.JenisInovasiResponse{}, err
	}

	return jenisinovasi.JenisInovasiResponse{
		ID:               result.ID,
		Jenis:            result.Jenis,
	}, nil
}

func (service *JenisInovasiServiceImpl) Update(ctx context.Context, request jenisinovasi.JenisInovasiUpdateRequest) (jenisinovasi.JenisInovasiResponse, error) {

	err := service.Validate.Struct(request)
	if err != nil {
		return jenisinovasi.JenisInovasiResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return jenisinovasi.JenisInovasiResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// validasi exists
	_, err = service.JenisInovasiRepository.FindById(ctx, tx, request.ID)
	if err != nil {
		return jenisinovasi.JenisInovasiResponse{}, err
	}

	data := domain.JenisInovasi{
		ID:             request.ID,
		Jenis:          request.Jenis,
	}

	result, err := service.JenisInovasiRepository.Update(ctx, tx, data)
	if err != nil {
		return jenisinovasi.JenisInovasiResponse{}, err
	}

	return jenisinovasi.JenisInovasiResponse{
		ID:               result.ID,
		Jenis:          result.Jenis,
	}, nil
}

func (service *JenisInovasiServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi data exists
	_, err = service.JenisInovasiRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}

	return service.JenisInovasiRepository.Delete(ctx, tx, id)
}

func (service *JenisInovasiServiceImpl) FindAll(ctx context.Context) ([]jenisinovasi.JenisInovasiResponse, error) {

	tx, err := service.DB.Begin()
	if err != nil {
		return []jenisinovasi.JenisInovasiResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data IKK
	results, err := service.JenisInovasiRepository.FindAll(ctx, tx)
	if err != nil {
		return []jenisinovasi.JenisInovasiResponse{}, err
	}

	var responses []jenisinovasi.JenisInovasiResponse
	for _, result := range results {
		responses = append(responses, jenisinovasi.JenisInovasiResponse{
			ID:                  result.ID,
			Jenis:               result.Jenis,
		})
	}

	return responses, nil
}

