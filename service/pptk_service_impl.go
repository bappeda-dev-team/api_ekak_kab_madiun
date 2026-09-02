package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/pptk"
	"ekak_kabupaten_madiun/repository"

	"github.com/go-playground/validator/v10"
)

type PptkServiceImpl struct {
	PptkRepository                  repository.PptkRepository
	DB                              *sql.DB
	Validate                        *validator.Validate
}

func NewPptkServiceImpl(pptkRepository repository.PptkRepository, db *sql.DB, validate *validator.Validate) *PptkServiceImpl {
	return &PptkServiceImpl{
		PptkRepository:                  pptkRepository,
		DB:                              db,
		Validate:                        validate,
	}
}

func (service *PptkServiceImpl) Create(ctx context.Context, request pptk.PptkCreateRequest) (pptk.PptkResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return pptk.PptkResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return pptk.PptkResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	datapptk := domain.Pptk{
		Nip:                     		 request.Nip,
		KodeOpd:       			 		 request.KodeOpd,
		Tahun:                   		 request.Tahun,
		KodeSubKegiatan:                 request.KodeSubKegiatan,
		NipAtasan:                       request.NipAtasan,
		NonAktifAt:                      request.NonAktifAt,
	}

	result, err := service.PptkRepository.Create(ctx, tx, datapptk)
	if err != nil {
		return pptk.PptkResponse{}, err
	}

	newData, err := service.PptkRepository.FindById(ctx, tx, result.Id)
	if err != nil {
		return pptk.PptkResponse{}, err
	}

	return pptk.PptkResponse{
		Id:                 newData.Id,
		Nip:                newData.Nip,
		NamaPegawai:        newData.NamaPegawai,
		KodeOpd:       	 	newData.KodeOpd,
		Tahun: 				newData.Tahun,
		KodeSubKegiatan:    newData.KodeSubKegiatan,
		NipAtasan:          newData.NipAtasan,
		NamaAtasan:         newData.NamaAtasan,
		AktifAt:            newData.AktifAt,
		NonAktifAt: 	    newData.NonAktifAt,
	}, nil
}

func (service *PptkServiceImpl) Update(ctx context.Context, request pptk.PptkUpdateRequest) (pptk.PptkResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return pptk.PptkResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return pptk.PptkResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi data exists
	_, err = service.PptkRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return pptk.PptkResponse{}, err
	}

	datapptk := domain.Pptk{
		Id:                              request.Id,
		Nip:                     		 request.Nip,
		KodeOpd:       			 		 request.KodeOpd,
		Tahun:                   		 request.Tahun,
		KodeSubKegiatan:                 request.KodeSubKegiatan,
		NipAtasan:                       request.NipAtasan,
		NonAktifAt:                      request.NonAktifAt,
	}

	result, err := service.PptkRepository.Update(ctx, tx, datapptk)
	if err != nil {
		return pptk.PptkResponse{}, err
	}

	updateData, err := service.PptkRepository.FindById(ctx, tx, result.Id)
	if err != nil {
		return pptk.PptkResponse{}, err
	}

	return pptk.PptkResponse{
		Id:                 updateData.Id,
		Nip:                updateData.Nip,
		NamaPegawai:        updateData.NamaPegawai,
		KodeOpd:       	 	updateData.KodeOpd,
		Tahun: 				updateData.Tahun,
		KodeSubKegiatan:    updateData.KodeSubKegiatan,
		NipAtasan:          updateData.NipAtasan,
		NamaAtasan:         updateData.NamaAtasan,
		AktifAt:            updateData.AktifAt,
		NonAktifAt: 	    updateData.NonAktifAt,
	}, nil
}

func (service *PptkServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi data exists
	_, err = service.PptkRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}

	return service.PptkRepository.Delete(ctx, tx, id)
}

func (service *PptkServiceImpl) FindById(ctx context.Context, id int) (pptk.PptkResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pptk.PptkResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	result, err := service.PptkRepository.FindById(ctx, tx, id)
	if err != nil {
		return pptk.PptkResponse{}, err
	}

	return pptk.PptkResponse{
		Id:                 result.Id,
		Nip:                result.Nip,
		NamaPegawai:        result.NamaPegawai,
		KodeOpd:       	 	result.KodeOpd,
		Tahun: 				result.Tahun,
		KodeSubKegiatan:    result.KodeSubKegiatan,
		NipAtasan:          result.NipAtasan,
		NamaAtasan:         result.NamaAtasan,
		AktifAt:            result.AktifAt,
		NonAktifAt: 	    result.NonAktifAt,
	}, nil
}

func (service *PptkServiceImpl) FindAll(ctx context.Context, kodeSubkegiatan string, kodeOpd string, tahun string) ([]pptk.PptkResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	results, err := service.PptkRepository.FindAll(ctx, tx, kodeSubkegiatan, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}

	var responses []pptk.PptkResponse
	for _, result := range results {
		responses = append(responses, pptk.PptkResponse{
			Id:                 result.Id,
			Nip:                result.Nip,
			KodeOpd:       	 	result.KodeOpd,
			Tahun: 				result.Tahun,
			KodeSubKegiatan:    result.KodeSubKegiatan,
			NipAtasan:          result.NipAtasan,
			AktifAt:            result.AktifAt,
			NonAktifAt: 	    result.NonAktifAt,
		})
	}

	return responses, nil
}
// func (service *PptkServiceImpl) FindAllByNip(ctx context.Context, nip string, tahun string) ([]pptk.PptkResponse, error) {
// 	tx, err := service.DB.Begin()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer helper.CommitOrRollback(tx)

// 	results, err := service.PptkRepository.FindAllByNip(ctx, tx, nip, tahun)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var responses []pptk.PptkResponse
// 	for _, result := range results {
// 		responses = append(responses, pptk.PptkResponse{
// 			Id:                 result.Id,
// 			Nip:                result.Nip,
// 			KodeOpd:       	 	result.KodeOpd,
// 			Tahun: 				result.Tahun,
// 			KodeSubKegiatan:    result.KodeSubKegiatan,
// 			NipAtasan:          result.NipAtasan,
// 			AktifAt:            result.AktifAt,
// 			NonAktifAt: 	    result.NonAktifAt,
// 		})
// 	}

// 	return responses, nil
// }