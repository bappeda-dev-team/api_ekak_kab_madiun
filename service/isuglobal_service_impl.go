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

func (service *IsuGlobalServiceImpl) FindAll(ctx context.Context, kodeOpd string) (isuglobal.IsuGlobalMasterResponse, error) {

	tx, err := service.DB.Begin()
	if err != nil {
		return isuglobal.IsuGlobalMasterResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil selection bidang urusan
	selections, err := service.IsuGlobalRepository.FindSelectionByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return isuglobal.IsuGlobalMasterResponse{}, err
	}

	// Ambil data IKK
	isus, err := service.IsuGlobalRepository.FindAll(ctx, tx, kodeOpd)
	if err != nil {
		return isuglobal.IsuGlobalMasterResponse{}, err
	}

	// Mapping selection
	selectionResponses := make([]isuglobal.BidangUrusanSelectionResponse, 0)

	for _, selection := range selections {
		selectionResponses = append(selectionResponses,
			isuglobal.BidangUrusanSelectionResponse{
				KodeBidangUrusan: selection.KodeBidangUrusan,
				NamaBidangUrusan: selection.NamaBidangUrusan,
				KodeOpd:          selection.KodeOpd,
				NamaOpd:          selection.NamaOpd,
			},
		)
	}

	// Mapping IKK
	isuResponses := make([]isuglobal.IsuGlobalFullResponse, 0)

	for _, isuData := range isus {

		isuResponses = append(isuResponses, isuglobal.IsuGlobalFullResponse{
			ID:               isuData.ID,
			KodeOpd:          isuData.KodeOpd,
			NamaOpd:          isuData.NamaOpd,
			KodeBidangUrusan: isuData.KodeBidangUrusan,
			NamaBidangUrusan: isuData.NamaBidangUrusan,
			Isu:              isuData.Isu,
			Tahun:            isuData.Tahun,
		})
	}

	return isuglobal.IsuGlobalMasterResponse{
		BidangUrusanSelections: selectionResponses,
		Isus:                   isuResponses,
	}, nil
}