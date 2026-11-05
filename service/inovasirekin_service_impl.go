package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/inovasirekin"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type InovasiRekinServiceImpl struct {
	inovasirekinRepository repository.InovasiRekinRepository
	rencanaKinerjaRepository         repository.RencanaKinerjaRepository
	rincianBelanjaRepository repository.RincianBelanjaRepository
	DB                     *sql.DB
}

func NewInovasiRekinServiceImpl(inovasirekinRepository repository.InovasiRekinRepository, rencanaKinerjaRepository repository.RencanaKinerjaRepository, rincianBelanjaRepository repository.RincianBelanjaRepository, DB *sql.DB) *InovasiRekinServiceImpl {
	return &InovasiRekinServiceImpl{
		inovasirekinRepository: inovasirekinRepository,
		rencanaKinerjaRepository:         rencanaKinerjaRepository,
		rincianBelanjaRepository: rincianBelanjaRepository,
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
		Kebaruan: 		   request.Kebaruan,
		AsalInovasi: 	   request.AsalInovasi,
		Tahun: 	           request.Tahun,
		NipInovator: 	   request.NipInovator,
		PegawaiId: 	       request.PegawaiId,
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
		Kebaruan: 		   request.Kebaruan,
		AsalInovasi: 	   request.AsalInovasi,
		NipInovator: 	   request.NipInovator,
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
func (service *InovasiRekinServiceImpl) FindAllKodeOpdTahun(ctx context.Context, kodeOpd string, tahun string) ([]inovasirekin.InovasiLaporanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback() // Hanya melakukan rollback jika belum di-commit

	inovasiRekins, err := service.inovasirekinRepository.FindAllKodeOpdTahun(ctx, tx, kodeOpd, tahun)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("rekin dengan ID %s tidak ditemukan", kodeOpd)
		}
		return nil, fmt.Errorf("gagal mengambil data: %v", err)
	}

	var rekinIds []string
	for _, rencana := range inovasiRekins {
		rekinIds = append(rekinIds, rencana.RekinId)
	}

	totalAnggaranByRekin, err := service.rincianBelanjaRepository.TotalAnggaranByIdRekins(ctx, tx, rekinIds)
	if err != nil {
		log.Printf("Gagal mengambil total anggaran: %v", err)
		return nil, fmt.Errorf("gagal mengambil total anggaran: %v", err)
	}

	var responses []inovasirekin.InovasiLaporanResponse
	for _, rencana := range inovasiRekins {
		log.Printf("Memproses RencanaKinerja dengan ID: %s", rencana.Id)

		indikators, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rencana.RekinId)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Gagal mencari Indikator: %v", err)
			return nil, fmt.Errorf("gagal mencari Indikator: %v", err)
		}

		var indikatorResponses []inovasirekin.IndikatorResponse
		for _, indikator := range indikators {
			targets, err := service.rencanaKinerjaRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
			if err != nil && err != sql.ErrNoRows {
				log.Printf("Gagal mencari Target: %v", err)
				return nil, fmt.Errorf("gagal mencari Target: %v", err)
			}

			var targetResponses []inovasirekin.TargetResponse
			for _, target := range targets {
				targetResponses = append(targetResponses, inovasirekin.TargetResponse{
					Id:              target.Id,
					IndikatorId:     target.IndikatorId,
					TargetIndikator: target.Target,
					SatuanIndikator: target.Satuan,
				})
			}

			indikatorResponses = append(indikatorResponses, inovasirekin.IndikatorResponse{
				Id:               indikator.Id,
				RencanaKinerjaId: indikator.RencanaKinerjaId,
				NamaIndikator:    indikator.Indikator,
				Target:           targetResponses,
			})
		}
		
		paguAnggaran := totalAnggaranByRekin[rencana.RekinId]

		responses = append(responses, inovasirekin.InovasiLaporanResponse{
			Id:           	 	rencana.Id,
			RekinId:      	 	rencana.RekinId,
			NamaRencanaKinerja: rencana.NamaRencanaKinerja,
			Indikator:          indikatorResponses,
			KodeOpd:      	 	rencana.KodeOpd,
			NamaOpd:      	 	rencana.NamaOpd,
			NamaInovasi:  	 	rencana.NamaInovasi,
			JenisInovasiId:  	rencana.JenisInovasiId,
			JenisInovasi:  		rencana.JenisInovasi,
			WaktuImplementasi:  rencana.WaktuImplementasi,
			Instansi:           rencana.Instansi,
			Inovator:           rencana.Inovator,
			Kebaruan:           rencana.Kebaruan,
			AsalInovasi:        rencana.AsalInovasi,
			Tahun:              rencana.Tahun,
			NipInovator:        rencana.NipInovator,
			NamaNipInovator:    rencana.NamaNipInovator,
			Level:              rencana.Level,
			PegawaiId:          rencana.PegawaiId,
			NamaPegawai:        rencana.NamaPegawai,
			NamaSubKegiatan:    rencana.NamaSubKegiatan,
			PaguAnggaran:       paguAnggaran,
		})
		log.Printf("RencanaKinerja Response ditambahkan untuk ID: %s", rencana.Id)
	}
	
	return responses, nil

	// return helper.ToInovasiLaporanResponses(inovasiRekins), nil
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
