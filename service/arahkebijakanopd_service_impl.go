package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/arahkebijakan"
	"ekak_kabupaten_madiun/repository"

	"github.com/go-playground/validator/v10"
)

type ArahKebijakanServiceImpl struct {
	ArahKebijakanRepository repository.ArahKebijakanRepository
	DB                      *sql.DB
	Validate                *validator.Validate
}

func NewArahKebijakanServiceImpl(arahkebijakanRepository repository.ArahKebijakanRepository, db *sql.DB, validate *validator.Validate) *ArahKebijakanServiceImpl {
	return &ArahKebijakanServiceImpl{
		ArahKebijakanRepository: arahkebijakanRepository,
		DB:             db,
		Validate:       validate,
	}
}

func (service *ArahKebijakanServiceImpl) Create(ctx context.Context, request arahkebijakan.ArahKebijakanRequest) (arahkebijakan.ArahKebijakanResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return arahkebijakan.ArahKebijakanResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return arahkebijakan.ArahKebijakanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	data := domain.ArahKebijakanOpd{
		PokinId: request.PokinId,
		Arah: request.Arah,
		KodeOpd: request.KodeOpd,
		Tahun:   request.Tahun,
	}

	result, err := service.ArahKebijakanRepository.Create(ctx, tx, data)
	if err != nil {
		return arahkebijakan.ArahKebijakanResponse{}, err
	}

	return arahkebijakan.ArahKebijakanResponse{
		ID:      result.ID,
		PokinId: result.PokinId,
		Arah: result.Arah,
		KodeOpd: result.KodeOpd,
		Tahun:   result.Tahun,
	}, nil
}

func (service *ArahKebijakanServiceImpl) Update(ctx context.Context, request arahkebijakan.ArahKebijakanUpdateRequest) (arahkebijakan.ArahKebijakanResponse, error) {

	err := service.Validate.Struct(request)
	if err != nil {
		return arahkebijakan.ArahKebijakanResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return arahkebijakan.ArahKebijakanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// validasi exists
	_, err = service.ArahKebijakanRepository.FindById(ctx, tx, request.ID)
	if err != nil {
		return arahkebijakan.ArahKebijakanResponse{}, err
	}

	data := domain.ArahKebijakanOpd{
		ID:      request.ID,
		PokinId: request.PokinId,
		Arah: request.Arah,
		KodeOpd: request.KodeOpd,
		Tahun:   request.Tahun,
	}

	result, err := service.ArahKebijakanRepository.Update(ctx, tx, data)
	if err != nil {
		return arahkebijakan.ArahKebijakanResponse{}, err
	}

	return arahkebijakan.ArahKebijakanResponse{
		ID:      result.ID,
		PokinId: result.PokinId,
		Arah: result.Arah,
		KodeOpd: result.KodeOpd,
		Tahun:   result.Tahun,
	}, nil
}