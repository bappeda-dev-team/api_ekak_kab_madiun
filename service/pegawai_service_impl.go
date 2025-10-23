package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"ekak_kabupaten_madiun/model/web/pegawai"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PegawaiServiceImpl struct {
	pegawaiRepository repository.PegawaiRepository
	opdRepository     repository.OpdRepository
	DB                *sql.DB
}

func NewPegawaiServiceImpl(pegawaiRepository repository.PegawaiRepository, opdRepository repository.OpdRepository, DB *sql.DB) *PegawaiServiceImpl {
	return &PegawaiServiceImpl{
		pegawaiRepository: pegawaiRepository,
		opdRepository:     opdRepository,
		DB:                DB,
	}
}

func (service *PegawaiServiceImpl) Create(ctx context.Context, request pegawai.PegawaiCreateRequest) (pegawai.PegawaiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pegawai.PegawaiResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// ✅ VALIDASI NIP TIDAK DUPLIKAT
	existingPegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, request.Nip)
	if err == nil {
		// Jika tidak ada error, berarti NIP sudah ada
		return pegawai.PegawaiResponse{}, fmt.Errorf("NIP %s sudah digunakan oleh pegawai %s", request.Nip, existingPegawai.NamaPegawai)
	}
	// Jika error adalah sql.ErrNoRows, berarti NIP belum ada (OK)
	if err != sql.ErrNoRows {
		// Jika error lain, return error
		return pegawai.PegawaiResponse{}, fmt.Errorf("gagal validasi NIP: %v", err)
	}

	// Generate UUID dan timestamp untuk ID yang lebih unik
	currentTime := time.Now().Format("20060102")
	uuid := uuid.New().String()
	pegawaiId := fmt.Sprintf("PEG-%s-%s", currentTime, uuid[:8])

	// Debug log untuk memastikan ID ter-generate
	fmt.Printf("Generated ID: %s\n", pegawaiId)

	pegawaiDomain := domainmaster.Pegawai{
		Id:          pegawaiId,
		NamaPegawai: request.NamaPegawai,
		Nip:         request.Nip,
		KodeOpd:     helper.EmptyStringIfNull(request.KodeOpd),
	}

	pegawais, err := service.pegawaiRepository.Create(ctx, tx, pegawaiDomain)
	if err != nil {
		return pegawai.PegawaiResponse{}, err
	}

	return helper.ToPegawaiResponse(pegawais), nil
}

func (service *PegawaiServiceImpl) Update(ctx context.Context, request pegawai.PegawaiUpdateRequest) (pegawai.PegawaiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pegawai.PegawaiResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data pegawai yang akan diupdate
	pegawaiData, err := service.pegawaiRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return pegawai.PegawaiResponse{}, err
	}

	// ✅ VALIDASI NIP TIDAK DUPLIKAT (hanya jika NIP berubah)
	if pegawaiData.Nip != request.Nip {
		existingPegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, request.Nip)
		if err == nil {
			// Jika tidak ada error, berarti NIP sudah digunakan oleh pegawai lain
			return pegawai.PegawaiResponse{}, fmt.Errorf("NIP %s sudah digunakan oleh pegawai %s", request.Nip, existingPegawai.NamaPegawai)
		}
		// Jika error adalah sql.ErrNoRows, berarti NIP belum ada (OK)
		if err != sql.ErrNoRows {
			// Jika error lain, return error
			return pegawai.PegawaiResponse{}, fmt.Errorf("gagal validasi NIP: %v", err)
		}
	}

	pegawaiData.NamaPegawai = request.NamaPegawai
	pegawaiData.Nip = request.Nip
	pegawaiData.KodeOpd = helper.EmptyStringIfNull(request.KodeOpd)

	updatedPegawai := service.pegawaiRepository.Update(ctx, tx, pegawaiData)
	return helper.ToPegawaiResponse(updatedPegawai), nil
}

func (service *PegawaiServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Tambahkan validasi jika id tidak ada
	pegawais, err := service.pegawaiRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}
	if pegawais.Id == "" {
		return fmt.Errorf("pegawai dengan id %s tidak ditemukan", id)
	}

	err = service.pegawaiRepository.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	return nil
}

func (service *PegawaiServiceImpl) FindById(ctx context.Context, id string) (pegawai.PegawaiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pegawai.PegawaiResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	pegawais, err := service.pegawaiRepository.FindById(ctx, tx, id)
	if err != nil {
		return pegawai.PegawaiResponse{}, err
	}

	// Tambahkan nama OPD jika pegawai memiliki kodeOpd
	if pegawais.KodeOpd != "" {
		opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pegawais.KodeOpd)
		if err == nil {
			pegawais.NamaOpd = opd.NamaOpd
		}
	}

	return helper.ToPegawaiResponse(pegawais), nil
}

func (service *PegawaiServiceImpl) FindAll(ctx context.Context, kodeOpd string) ([]pegawai.PegawaiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []pegawai.PegawaiResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	pegawais, err := service.pegawaiRepository.FindAll(ctx, tx, kodeOpd)
	if err != nil {
		return []pegawai.PegawaiResponse{}, err
	}

	for i := range pegawais {
		if pegawais[i].KodeOpd != "" {
			opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pegawais[i].KodeOpd)
			if err == nil {
				pegawais[i].NamaOpd = opd.NamaOpd
			}
		}
	}

	return helper.ToPegawaiResponses(pegawais), nil
}
