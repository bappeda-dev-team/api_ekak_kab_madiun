package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/ppd"
	"ekak_kabupaten_madiun/repository"

	"github.com/go-playground/validator/v10"
)

type PpdServiceImpl struct {
	PpdRepository repository.PpdRepository
	DB                  *sql.DB
	Validate            *validator.Validate
}

func NewPpdServiceImpl(ppdRepository repository.PpdRepository, db *sql.DB, validate *validator.Validate) *PpdServiceImpl {
	return &PpdServiceImpl{
		PpdRepository: ppdRepository,
		DB:                  db,
		Validate:            validate,
	}
}

func (service *PpdServiceImpl) Create(ctx context.Context, request ppd.PpdRequest) (ppd.PpdResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return ppd.PpdResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return ppd.PpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	data := domain.PotensiPerangkatDaerah{
		KodeBidangUrusan: request.KodeBidangUrusan,
		KodeOpd:          request.KodeOpd,
		Potensi:              request.Potensi,
		Tahun:            request.Tahun,
	}

	result, err := service.PpdRepository.Create(ctx, tx, data)
	if err != nil {
		return ppd.PpdResponse{}, err
	}

	return ppd.PpdResponse{
		ID:                 result.ID,
		KodeBidangUrusan:   result.KodeBidangUrusan,
		KodeOpd:            result.KodeOpd,
		Potensi:                result.Potensi,
		Tahun:              result.Tahun,
	}, nil
}

func (service *PpdServiceImpl) Update(
	ctx context.Context,
	request ppd.PpdUpdateRequest,
) (ppd.PpdResponse, error) {

	err := service.Validate.Struct(request)
	if err != nil {
		return ppd.PpdResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return ppd.PpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// validasi exists
	_, err = service.PpdRepository.FindById(ctx, tx, request.ID)
	if err != nil {
		return ppd.PpdResponse{}, err
	}


	data := domain.PotensiPerangkatDaerah{
		ID:                 request.ID,
		KodeBidangUrusan:   request.KodeBidangUrusan,
		KodeOpd:            request.KodeOpd,
		Potensi:                request.Potensi,
		Tahun:              request.Tahun,
	}

	result, err := service.PpdRepository.Update(ctx, tx, data)
	if err != nil {
		return ppd.PpdResponse{}, err
	}

	return ppd.PpdResponse{
		ID:                 result.ID,
		KodeBidangUrusan:   result.KodeBidangUrusan,
		KodeOpd:            result.KodeOpd,
		Potensi:                result.Potensi,
		Tahun:              result.Tahun,
	}, nil
}

func (service *PpdServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi data exists
	_, err = service.PpdRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}

	return service.PpdRepository.Delete(ctx, tx, id)
}

func (service *PpdServiceImpl) FindAll(ctx context.Context, kodeOpd string) (ppd.PpdMasterResponse, error) {

	tx, err := service.DB.Begin()
	if err != nil {
		return ppd.PpdMasterResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil selection bidang urusan
	selections, err := service.PpdRepository.FindSelectionByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return ppd.PpdMasterResponse{}, err
	}

	// Ambil data IKK
	isus, err := service.PpdRepository.FindAll(ctx, tx, kodeOpd)
	if err != nil {
		return ppd.PpdMasterResponse{}, err
	}

	// Mapping selection
	selectionResponses := make([]ppd.BidangUrusanSelectionResponse, 0)

	for _, selection := range selections {
		selectionResponses = append(selectionResponses,
			ppd.BidangUrusanSelectionResponse{
				KodeBidangUrusan: selection.KodeBidangUrusan,
				NamaBidangUrusan: selection.NamaBidangUrusan,
				KodeOpd:          selection.KodeOpd,
				NamaOpd:          selection.NamaOpd,
			},
		)
	}

	// Mapping IKK
	isuResponses := make([]ppd.PpdFullResponse, 0)

	for _, isuData := range isus {

		isuResponses = append(isuResponses, ppd.PpdFullResponse{
			ID:               isuData.ID,
			KodeOpd:          isuData.KodeOpd,
			NamaOpd:          isuData.NamaOpd,
			KodeBidangUrusan: isuData.KodeBidangUrusan,
			NamaBidangUrusan: isuData.NamaBidangUrusan,
			Potensi:              isuData.Potensi,
			Tahun:            isuData.Tahun,
		})
	}

	return ppd.PpdMasterResponse{
		BidangUrusanSelections: selectionResponses,
		Ppds:                   isuResponses,
	}, nil
}