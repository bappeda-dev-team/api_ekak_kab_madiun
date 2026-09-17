package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/inovasirekin"
	"ekak_kabupaten_madiun/repository"
	"fmt"

	"github.com/google/uuid"
)

type InovasiRekinServiceImpl struct {
	inovasirekinRepository repository.InovasiRekinRepository
	DB                     *sql.DB
}

func NewInovasiRekinServiceImpl(inovasirekinRepository repository.InovasiRekinRepository, DB *sql.DB) *InovasiRekinServiceImpl {
	return &InovasiRekinServiceImpl{
		inovasirekinRepository: inovasirekinRepository,
		DB:                     DB,
	}
}

func (service *InovasiRekinServiceImpl) Create(ctx context.Context, request inovasirekin.InovasiRekinCreateRequest) (inovasirekin.InovasiRekinResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return inovasirekin.InovasiRekinResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Membuat UUID dengan format yang diinginkan
	randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
	uuId := fmt.Sprintf("INVS-REKIN-%s", randomDigits)

	domainInovasiRekin := domain.InovasiRekin{
		Id:           	   uuId,
		RekinId:      	   request.RekinId,
		KodeOpd:      	   request.KodeOpd,
		NamaInovasi: 	   request.NamaInovasi,
		JenisInovasiId:    request.JenisInovasiId,
		WaktuImplementasi: request.WaktuImplementasi,
		Instansi: 	       request.Instansi,
		Inovator: 		   request.Inovator,
	}

	inovasis, err := service.inovasirekinRepository.Create(ctx, tx, domainInovasiRekin)
	if err != nil {
		return inovasirekin.InovasiRekinResponse{}, err
	}

	response := helper.ToInovasiRekinResponse(inovasis) 
	return response, nil
}

func (service *InovasiRekinServiceImpl) Update(ctx context.Context, request inovasirekin.InovasiRekinUpdateRequest) (inovasirekin.InovasiRekinResponse, error) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	inovasiRekin := domain.InovasiRekin{
		Id:                request.Id,
		NamaInovasi: 	   request.NamaInovasi,
		JenisInovasiId:    request.JenisInovasiId,
		WaktuImplementasi: request.WaktuImplementasi,
		Instansi: 	       request.Instansi,
		Inovator: 		   request.Inovator,
	}

	inovasiRekin, err = service.inovasirekinRepository.Update(ctx, tx, inovasiRekin)
	if err != nil {
		return inovasirekin.InovasiRekinResponse{}, err
	}

	return helper.ToInovasiRekinResponse(inovasiRekin), nil
}

func (service *InovasiRekinServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	return service.inovasirekinRepository.Delete(ctx, tx, id)

}

func (service *InovasiRekinServiceImpl) FindAll(ctx context.Context, rekinId string) ([]inovasirekin.InovasiRekinResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback() // Hanya melakukan rollback jika belum di-commit

	inovasiRekins, err := service.inovasirekinRepository.FindAll(ctx, tx, rekinId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("rekin dengan ID %s tidak ditemukan", rekinId)
		}
		return nil, fmt.Errorf("gagal mengambil data: %v", err)
	}

	// if len(inovasiRekins) == 0 {
	// 	return nil, fmt.Errorf("tidak ada gambaran umum untuk rekin dengan ID %s", rekinId)
	// }

	// Commit transaksi jika berhasil
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("gagal melakukan commit transaksi: %v", err)
	}

	return helper.ToInovasiRekinResponses(inovasiRekins), nil
}

func (service *InovasiRekinServiceImpl) FindById(ctx context.Context, id string) (inovasirekin.InovasiRekinResponse, error) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	inovasiRekin, err := service.inovasirekinRepository.FindById(ctx, tx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return inovasirekin.InovasiRekinResponse{}, fmt.Errorf("inovasi dengan ID %s tidak ditemukan", id)
		}
		return inovasirekin.InovasiRekinResponse{}, err
	}

	return helper.ToInovasiRekinResponse(inovasiRekin), nil
}
