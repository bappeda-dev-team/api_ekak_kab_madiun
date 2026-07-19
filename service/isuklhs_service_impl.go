package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/isuklhs"
	"ekak_kabupaten_madiun/repository"

	"github.com/go-playground/validator/v10"
)

type IsuKlhsServiceImpl struct {
	IsuKlhsRepository repository.IsuKlhsRepository
	DB                  *sql.DB
	Validate            *validator.Validate
}

func NewIsuKlhsServiceImpl(isuklhsRepository repository.IsuKlhsRepository, db *sql.DB, validate *validator.Validate) *IsuKlhsServiceImpl {
	return &IsuKlhsServiceImpl{
		IsuKlhsRepository: isuklhsRepository,
		DB:                  db,
		Validate:            validate,
	}
}

func (service *IsuKlhsServiceImpl) Create(ctx context.Context, request isuklhs.IsuKlhsRequest) (isuklhs.IsuKlhsResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return isuklhs.IsuKlhsResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return isuklhs.IsuKlhsResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	data := domain.IsuKlhs{
		KodeBidangUrusan: request.KodeBidangUrusan,
		KodeOpd:          request.KodeOpd,
		Isu:              request.Isu,
		Tahun:            request.Tahun,
	}

	result, err := service.IsuKlhsRepository.Create(ctx, tx, data)
	if err != nil {
		return isuklhs.IsuKlhsResponse{}, err
	}

	return isuklhs.IsuKlhsResponse{
		ID:                 result.ID,
		KodeBidangUrusan:   result.KodeBidangUrusan,
		KodeOpd:            result.KodeOpd,
		Isu:                result.Isu,
		Tahun:              result.Tahun,
	}, nil
}

func (service *IsuKlhsServiceImpl) Update(
	ctx context.Context,
	request isuklhs.IsuKlhsUpdateRequest,
) (isuklhs.IsuKlhsResponse, error) {

	err := service.Validate.Struct(request)
	if err != nil {
		return isuklhs.IsuKlhsResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return isuklhs.IsuKlhsResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// validasi exists
	_, err = service.IsuKlhsRepository.FindById(ctx, tx, request.ID)
	if err != nil {
		return isuklhs.IsuKlhsResponse{}, err
	}


	data := domain.IsuKlhs{
		ID:                 request.ID,
		KodeBidangUrusan:   request.KodeBidangUrusan,
		KodeOpd:            request.KodeOpd,
		Isu:                request.Isu,
		Tahun:              request.Tahun,
	}

	result, err := service.IsuKlhsRepository.Update(ctx, tx, data)
	if err != nil {
		return isuklhs.IsuKlhsResponse{}, err
	}

	return isuklhs.IsuKlhsResponse{
		ID:                 result.ID,
		KodeBidangUrusan:   result.KodeBidangUrusan,
		KodeOpd:            result.KodeOpd,
		Isu:                result.Isu,
		Tahun:              result.Tahun,
	}, nil
}

func (service *IsuKlhsServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi data exists
	_, err = service.IsuKlhsRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}

	return service.IsuKlhsRepository.Delete(ctx, tx, id)
}

func (service *IsuKlhsServiceImpl) FindAll(ctx context.Context, kodeOpd string) (isuklhs.IsuKlhsMasterResponse, error) {

	tx, err := service.DB.Begin()
	if err != nil {
		return isuklhs.IsuKlhsMasterResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil selection bidang urusan
	selections, err := service.IsuKlhsRepository.FindSelectionByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return isuklhs.IsuKlhsMasterResponse{}, err
	}

	// Ambil data IKK
	isus, err := service.IsuKlhsRepository.FindAll(ctx, tx, kodeOpd)
	if err != nil {
		return isuklhs.IsuKlhsMasterResponse{}, err
	}

	// Mapping selection
	selectionResponses := make([]isuklhs.BidangUrusanSelectionResponse, 0)

	for _, selection := range selections {
		selectionResponses = append(selectionResponses,
			isuklhs.BidangUrusanSelectionResponse{
				KodeBidangUrusan: selection.KodeBidangUrusan,
				NamaBidangUrusan: selection.NamaBidangUrusan,
				KodeOpd:          selection.KodeOpd,
				NamaOpd:          selection.NamaOpd,
			},
		)
	}

	// Mapping IKK
	isuResponses := make([]isuklhs.IsuKlhsFullResponse, 0)

	for _, isuData := range isus {

		isuResponses = append(isuResponses, isuklhs.IsuKlhsFullResponse{
			ID:               isuData.ID,
			KodeOpd:          isuData.KodeOpd,
			NamaOpd:          isuData.NamaOpd,
			KodeBidangUrusan: isuData.KodeBidangUrusan,
			NamaBidangUrusan: isuData.NamaBidangUrusan,
			Isu:              isuData.Isu,
			Tahun:            isuData.Tahun,
		})
	}

	return isuklhs.IsuKlhsMasterResponse{
		BidangUrusanSelections: selectionResponses,
		Isus:                   isuResponses,
	}, nil
}