package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/ikk"
	"ekak_kabupaten_madiun/repository"

	"github.com/go-playground/validator/v10"
)

type IkkServiceImpl struct {
	IkkRepository repository.IkkRepository
	DB            *sql.DB
	Validate      *validator.Validate
}

func NewIkkServiceImpl(ikkRepository repository.IkkRepository, db *sql.DB, validate *validator.Validate) *IkkServiceImpl {
	return &IkkServiceImpl{
		IkkRepository: ikkRepository,
		DB:            db,
		Validate:      validate,
	}
}

func (service *IkkServiceImpl) Create(ctx context.Context, request ikk.IkkRequest) (ikk.IkkResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return ikk.IkkResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return ikk.IkkResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	data := domain.IndikatorIkk{
		KodeBidangUrusan:          request.KodeBidangUrusan,
		Jenis: 					   request.Jenis,
		NamaIndikator:             request.NamaIndikator,
		Target:                    request.Target,
		Satuan:                	   request.Satuan,
		Keterangan:                request.Keterangan,
	}

	result, err := service.IkkRepository.Create(ctx, tx, data)
	if err != nil {
		return ikk.IkkResponse{}, err
	}

	return ikk.IkkResponse{
		ID:                         result.ID,
		KodeBidangUrusan:           result.KodeBidangUrusan,
		Jenis:       				result.Jenis,
		NamaIndikator: 				result.NamaIndikator,
		Target:                		result.Target,
		Satuan:                 	result.Satuan,
		Keterangan:                 result.Keterangan,
	}, nil
}

func (service *IkkServiceImpl) Update(ctx context.Context, request ikk.IkkUpdateRequest) (ikk.IkkResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return ikk.IkkResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return ikk.IkkResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi data exists
	_, err = service.IkkRepository.FindById(ctx, tx, request.ID)
	if err != nil {
		return ikk.IkkResponse{}, err
	}

	data := domain.IndikatorIkk{
		ID:                        		request.ID,
		KodeBidangUrusan:               request.KodeBidangUrusan,
		Jenis: 							request.Jenis,
		NamaIndikator:                	request.NamaIndikator,
		Target:                 		request.Target,
		Satuan:                			request.Satuan,
		Keterangan:                		request.Keterangan,
	}

	result, err := service.IkkRepository.Update(ctx, tx, data)
	if err != nil {
		return ikk.IkkResponse{}, err
	}

	updateData, err := service.IkkRepository.FindById(ctx, tx, result.ID)
	if err != nil {
		return ikk.IkkResponse{}, err
	}

	return ikk.IkkResponse{
		ID:                        updateData.ID,
		KodeBidangUrusan:          updateData.KodeBidangUrusan,
		Jenis:       			   updateData.Jenis,
		NamaIndikator: 			   updateData.NamaIndikator,
		Target:                	   updateData.Target,
		Satuan:                    updateData.Satuan,
		Keterangan:                updateData.Keterangan,
	}, nil
}

func (service *IkkServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi data exists
	_, err = service.IkkRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}

	return service.IkkRepository.Delete(ctx, tx, id)
}

func (service *IkkServiceImpl) FindById(ctx context.Context, id int) (ikk.IkkResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return ikk.IkkResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	result, err := service.IkkRepository.FindById(ctx, tx, id)
	if err != nil {
		return ikk.IkkResponse{}, err
	}

	return ikk.IkkResponse{
		ID:                        result.ID,
		KodeBidangUrusan:          result.KodeBidangUrusan,
		Jenis:       			   result.Jenis,
		NamaIndikator: 			   result.NamaIndikator,
		Target:                    result.Target,
		Satuan:                    result.Satuan,
		Keterangan:                result.Keterangan,
	}, nil
}

func (service *IkkServiceImpl) FindByKodeOpd(ctx context.Context, levelPohon int, kodeOpd string) ([]ikk.IkkFullResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []ikk.IkkFullResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	mapping := map[int]string{
		5: "outcome",
		6: "output",
	}

	jenis := mapping[levelPohon]

	if jenis == "" {
		return []ikk.IkkFullResponse{}, nil
	}

	bidangUrusans, err := service.IkkRepository.FindByKodeOpd(ctx, tx, jenis, kodeOpd)
	if err != nil {
		return []ikk.IkkFullResponse{}, err
	}

	var bidangUrusanResponses []ikk.IkkFullResponse
	for _, bidangUrusan := range bidangUrusans {
		bidangUrusanResponses = append(bidangUrusanResponses, ikk.IkkFullResponse{
			ID: bidangUrusan.ID,
			KodeBidangUrusan: bidangUrusan.KodeBidangUrusan,
			NamaOpd: bidangUrusan.NamaOpd,
			Jenis: bidangUrusan.Jenis,
			NamaIndikator: bidangUrusan.NamaIndikator,
			Target: bidangUrusan.Target,
			Satuan: bidangUrusan.Satuan,
			Keterangan: bidangUrusan.Keterangan,
		})
	}

	return bidangUrusanResponses, nil
}
