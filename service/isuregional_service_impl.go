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

func (service *IsuRegionalServiceImpl) FindAll(ctx context.Context, kodeOpd string) (isuregional.IsuRegionalMasterResponse, error) {

	tx, err := service.DB.Begin()
	if err != nil {
		return isuregional.IsuRegionalMasterResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil selection bidang urusan
	selections, err := service.IsuRegionalRepository.FindSelectionByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return isuregional.IsuRegionalMasterResponse{}, err
	}

	// Ambil data IKK
	isus, err := service.IsuRegionalRepository.FindAll(ctx, tx, kodeOpd)
	if err != nil {
		return isuregional.IsuRegionalMasterResponse{}, err
	}

	// Mapping selection
	selectionResponses := make([]isuregional.BidangUrusanSelectionResponse, 0)

	for _, selection := range selections {
		selectionResponses = append(selectionResponses,
			isuregional.BidangUrusanSelectionResponse{
				KodeBidangUrusan: selection.KodeBidangUrusan,
				NamaBidangUrusan: selection.NamaBidangUrusan,
				KodeOpd:          selection.KodeOpd,
				NamaOpd:          selection.NamaOpd,
			},
		)
	}

	// Mapping IKK
	isuResponses := make([]isuregional.IsuRegionalFullResponse, 0)

	for _, isuData := range isus {

		isuResponses = append(isuResponses, isuregional.IsuRegionalFullResponse{
			ID:               isuData.ID,
			KodeOpd:          isuData.KodeOpd,
			NamaOpd:          isuData.NamaOpd,
			KodeBidangUrusan: isuData.KodeBidangUrusan,
			NamaBidangUrusan: isuData.NamaBidangUrusan,
			Isu:              isuData.Isu,
			Tahun:            isuData.Tahun,
		})
	}

	return isuregional.IsuRegionalMasterResponse{
		BidangUrusanSelections: selectionResponses,
		Isus:                   isuResponses,
	}, nil
}