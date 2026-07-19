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

func (service *IsuNasionalServiceImpl) FindAll(ctx context.Context, kodeOpd string) (isunasional.IsuNasionalMasterResponse, error) {

	tx, err := service.DB.Begin()
	if err != nil {
		return isunasional.IsuNasionalMasterResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil selection bidang urusan
	selections, err := service.IsuNasionalRepository.FindSelectionByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return isunasional.IsuNasionalMasterResponse{}, err
	}

	// Ambil data IKK
	isus, err := service.IsuNasionalRepository.FindAll(ctx, tx, kodeOpd)
	if err != nil {
		return isunasional.IsuNasionalMasterResponse{}, err
	}

	// Mapping selection
	selectionResponses := make([]isunasional.BidangUrusanSelectionResponse, 0)

	for _, selection := range selections {
		selectionResponses = append(selectionResponses,
			isunasional.BidangUrusanSelectionResponse{
				KodeBidangUrusan: selection.KodeBidangUrusan,
				NamaBidangUrusan: selection.NamaBidangUrusan,
				KodeOpd:          selection.KodeOpd,
				NamaOpd:          selection.NamaOpd,
			},
		)
	}

	// Mapping IKK
	isuResponses := make([]isunasional.IsuNasionalFullResponse, 0)

	for _, isuData := range isus {

		isuResponses = append(isuResponses, isunasional.IsuNasionalFullResponse{
			ID:               isuData.ID,
			KodeOpd:          isuData.KodeOpd,
			NamaOpd:          isuData.NamaOpd,
			KodeBidangUrusan: isuData.KodeBidangUrusan,
			NamaBidangUrusan: isuData.NamaBidangUrusan,
			Isu:              isuData.Isu,
			Tahun:            isuData.Tahun,
		})
	}

	return isunasional.IsuNasionalMasterResponse{
		BidangUrusanSelections: selectionResponses,
		Isus:                   isuResponses,
	}, nil
}