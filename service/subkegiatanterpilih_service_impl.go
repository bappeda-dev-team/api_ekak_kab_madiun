package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/subkegiatan"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
)

type SubKegiatanTerpilihServiceImpl struct {
	RencanaKinerjaRepository      repository.RencanaKinerjaRepository
	SubKegiatanRepository         repository.SubKegiatanRepository
	SubKegiatanTerpilihRepository repository.SubKegiatanTerpilihRepository
	opdRepository                 repository.OpdRepository
	DB                            *sql.DB
	Validate                      *validator.Validate
}

func NewSubKegiatanTerpilihServiceImpl(rencanaKinerjaRepository repository.RencanaKinerjaRepository, subKegiatanRepository repository.SubKegiatanRepository, subKegiatanTerpilihRepository repository.SubKegiatanTerpilihRepository, opdRepository repository.OpdRepository, DB *sql.DB, Validate *validator.Validate) *SubKegiatanTerpilihServiceImpl {
	return &SubKegiatanTerpilihServiceImpl{
		RencanaKinerjaRepository:      rencanaKinerjaRepository,
		SubKegiatanRepository:         subKegiatanRepository,
		SubKegiatanTerpilihRepository: subKegiatanTerpilihRepository,
		opdRepository:                 opdRepository,
		DB:                            DB,
		Validate:                      Validate,
	}
}

func (service *SubKegiatanTerpilihServiceImpl) Update(ctx context.Context, request subkegiatan.SubKegiatanTerpilihUpdateRequest) (subkegiatan.SubKegiatanTerpilihResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return subkegiatan.SubKegiatanTerpilihResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	var rencanaKinerja domain.RencanaKinerja
	if request.Id != "" {
		rencanaKinerja, err = service.RencanaKinerjaRepository.FindById(ctx, tx, request.Id, "", "")
		if err != nil {
			log.Printf("Gagal menemukan RencanaKinerja: %v", err)
			return subkegiatan.SubKegiatanTerpilihResponse{}, fmt.Errorf("gagal menemukan RencanaKinerja: %v", err)
		}
	} else {
		return subkegiatan.SubKegiatanTerpilihResponse{}, errors.New("id rencana kinerja tidak boleh kosong")
	}

	// Cek apakah data dengan kode_subkegiatan tersebut ada
	_, err = service.SubKegiatanRepository.FindByKodeSubKegiatan(ctx, tx, request.KodeSubKegiatan)
	if err != nil {
		return subkegiatan.SubKegiatanTerpilihResponse{}, errors.New("subkegiatan tidak ditemukan")
	}

	subKegiatanTerpilih := domain.SubKegiatanTerpilih{
		Id:              rencanaKinerja.Id,
		KodeSubKegiatan: request.KodeSubKegiatan,
	}

	result, err := service.SubKegiatanTerpilihRepository.Update(ctx, tx, subKegiatanTerpilih)
	if err != nil {
		return subkegiatan.SubKegiatanTerpilihResponse{}, err
	}

	return subkegiatan.SubKegiatanTerpilihResponse{
		KodeSubKegiatan: subkegiatan.SubKegiatanResponse{
			KodeSubKegiatan: result.KodeSubKegiatan,
		},
	}, nil
}

func (service *SubKegiatanTerpilihServiceImpl) FindByKodeSubKegiatan(ctx context.Context, kodeSubKegiatan string) (subkegiatan.SubKegiatanTerpilihResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return subkegiatan.SubKegiatanTerpilihResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah data dengan kode_subkegiatan tersebut ada

	result, err := service.SubKegiatanRepository.FindByKodeSubKegiatan(ctx, tx, kodeSubKegiatan)
	if err != nil {
		return subkegiatan.SubKegiatanTerpilihResponse{}, errors.New("subkegiatan tidak ditemukan")
	}

	return subkegiatan.SubKegiatanTerpilihResponse{
		KodeSubKegiatan: subkegiatan.SubKegiatanResponse{
			KodeSubKegiatan: result.KodeSubKegiatan,
			NamaSubKegiatan: result.NamaSubKegiatan,
		},
	}, nil
}

func (service *SubKegiatanTerpilihServiceImpl) Delete(ctx context.Context, id string, kodeSubKegiatan string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi: Cek apakah data dengan id dan kodeSubKegiatan ada
	_, err = service.SubKegiatanTerpilihRepository.FindByIdAndKodeSubKegiatan(ctx, tx, id, kodeSubKegiatan)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data tidak ditemukan")
		}
		return err
	}

	// Lanjutkan dengan penghapusan jika data ditemukan
	err = service.SubKegiatanTerpilihRepository.Delete(ctx, tx, id, kodeSubKegiatan)
	if err != nil {
		return err
	}

	return nil
}

func (service *SubKegiatanTerpilihServiceImpl) CreateRekin(ctx context.Context, request subkegiatan.SubKegiatanCreateRekinRequest) ([]subkegiatan.SubKegiatanResponse, error) {
	// Konversi single ID ke array
	idSubKegiatanArray := []string{request.IdSubKegiatan}

	tx, err := service.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah rencana kinerja dengan ID yang diberikan ada
	_, err = service.RencanaKinerjaRepository.FindById(ctx, tx, request.RekinId, "", "")
	if err != nil {
		return nil, fmt.Errorf("rencana kinerja dengan id %s tidak ditemukan: %v", request.RekinId, err)
	}

	var updatedSubKegiatans []domain.SubKegiatan

	// Proses setiap ID usulan
	for _, idSubKegiatan := range idSubKegiatanArray {
		// Cek apakah usulan dengan ID yang diberikan ada
		existingSubKegiatan, err := service.SubKegiatanRepository.FindById(ctx, tx, idSubKegiatan)
		if err != nil {
			return nil, fmt.Errorf("subkegiatan dengan id %s tidak ditemukan: %v", idSubKegiatan, err)
		}

		// Cek apakah usulan sudah memiliki rekin_id
		if existingSubKegiatan.RekinId != "" {
			return nil, fmt.Errorf("subkegiatan dengan id %s sudah memiliki rencana kinerja", idSubKegiatan)
		}

		// Update rekin_id dan status
		err = service.SubKegiatanTerpilihRepository.CreateRekin(ctx, tx, idSubKegiatan, request.RekinId)
		if err != nil {
			return nil, fmt.Errorf("gagal mengupdate rekin untuk subkegiatan %s: %v", idSubKegiatan, err)
		}

		// Ambil data usulan yang sudah diupdate
		updatedSubKegiatan, err := service.SubKegiatanRepository.FindById(ctx, tx, idSubKegiatan)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil data subkegiatan yang diupdate: %v", err)
		}

		updatedSubKegiatans = append(updatedSubKegiatans, updatedSubKegiatan)
	}

	// Konversi ke response
	responses := helper.ToSubKegiatanResponses(updatedSubKegiatans)
	return responses, nil
}

func (service *SubKegiatanTerpilihServiceImpl) DeleteSubKegiatanTerpilih(ctx context.Context, idSubKegiatan string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	err = service.SubKegiatanTerpilihRepository.DeleteSubKegiatanTerpilih(ctx, tx, idSubKegiatan)
	if err != nil {
		return fmt.Errorf("gagal menghapus subkegiatan terpilih: %v", err)
	}

	return nil
}

func (service *SubKegiatanTerpilihServiceImpl) CreateOpd(ctx context.Context, request subkegiatan.SubKegiatanOpdCreateRequest) (subkegiatan.SubKegiatanOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return subkegiatan.SubKegiatanOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	//cek subkegiatan
	kode, err := service.SubKegiatanRepository.FindByKodeSubKegiatan(ctx, tx, request.KodeSubkegiatan)
	if err != nil {
		if err == sql.ErrNoRows {
			return subkegiatan.SubKegiatanOpdResponse{}, fmt.Errorf("kode subkegiatan %s tidak ditemukan dalam database", request.KodeSubkegiatan)
		}
		return subkegiatan.SubKegiatanOpdResponse{}, fmt.Errorf("terjadi kesalahan saat mencari data kode subkegiatan: %v", err)
	}

	// Cek OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOpd)
	if err != nil {
		if err == sql.ErrNoRows {
			return subkegiatan.SubKegiatanOpdResponse{}, fmt.Errorf("OPD dengan kode %s tidak ditemukan dalam database", request.KodeOpd)
		}
		return subkegiatan.SubKegiatanOpdResponse{}, fmt.Errorf("terjadi kesalahan saat mencari data OPD: %v", err)
	}

	domain := domain.SubKegiatanOpd{
		KodeSubKegiatan: kode.KodeSubKegiatan,
		KodeOpd:         opd.KodeOpd,
		Tahun:           request.Tahun,
	}

	result, err := service.SubKegiatanTerpilihRepository.CreateOPD(ctx, tx, domain)
	if err != nil {
		return subkegiatan.SubKegiatanOpdResponse{}, err
	}

	response := subkegiatan.SubKegiatanOpdResponse{
		Id:              result.Id,
		KodeSubkegiatan: result.KodeSubKegiatan,
		KodeOpd:         result.KodeOpd,
		Tahun:           result.Tahun,
	}

	return response, nil
}

func (service *SubKegiatanTerpilihServiceImpl) UpdateOpd(ctx context.Context, request subkegiatan.SubKegiatanOpdUpdateRequest) (subkegiatan.SubKegiatanOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return subkegiatan.SubKegiatanOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	//cek subkegiatan
	kode, err := service.SubKegiatanRepository.FindByKodeSubKegiatan(ctx, tx, request.KodeSubkegiatan)
	if err != nil {
		if err == sql.ErrNoRows {
			return subkegiatan.SubKegiatanOpdResponse{}, fmt.Errorf("kode subkegiatan %s tidak ditemukan dalam database", request.KodeSubkegiatan)
		}
		return subkegiatan.SubKegiatanOpdResponse{}, fmt.Errorf("terjadi kesalahan saat mencari data kode subkegiatan: %v", err)
	}

	// Cek OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOpd)
	if err != nil {
		if err == sql.ErrNoRows {
			return subkegiatan.SubKegiatanOpdResponse{}, fmt.Errorf("OPD dengan kode %s tidak ditemukan", request.KodeOpd)
		}
		return subkegiatan.SubKegiatanOpdResponse{}, err
	}

	domain := domain.SubKegiatanOpd{
		Id:              request.Id,
		KodeSubKegiatan: kode.KodeSubKegiatan,
		KodeOpd:         opd.KodeOpd,
		Tahun:           request.Tahun,
	}

	result, err := service.SubKegiatanTerpilihRepository.UpdateOPD(ctx, tx, domain)
	if err != nil {
		return subkegiatan.SubKegiatanOpdResponse{}, err
	}

	response := subkegiatan.SubKegiatanOpdResponse{
		Id:              result.Id,
		KodeSubkegiatan: result.KodeSubKegiatan,
		KodeOpd:         result.KodeOpd,
		Tahun:           result.Tahun,
	}
	return response, nil
}

func (service *SubKegiatanTerpilihServiceImpl) FindAllOpd(ctx context.Context, kodeOpd, tahun *string) ([]subkegiatan.SubKegiatanOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	result, err := service.SubKegiatanTerpilihRepository.FindallOpd(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}

	var responses []subkegiatan.SubKegiatanOpdResponse
	for _, sub := range result {
		kode, _ := service.SubKegiatanRepository.FindByKodeSubKegiatan(ctx, tx, sub.KodeSubKegiatan)
		opd, _ := service.opdRepository.FindByKodeOpd(ctx, tx, sub.KodeOpd)

		response := subkegiatan.SubKegiatanOpdResponse{
			Id:              sub.Id,
			KodeSubkegiatan: kode.KodeSubKegiatan,
			NamaSubkegiatan: kode.NamaSubKegiatan,
			KodeOpd:         opd.KodeOpd,
			NamaOpd:         opd.NamaOpd,
			Tahun:           sub.Tahun,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (service *SubKegiatanTerpilihServiceImpl) FindById(ctx context.Context, id int) (subkegiatan.SubKegiatanOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return subkegiatan.SubKegiatanOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	result, err := service.SubKegiatanTerpilihRepository.FindById(ctx, tx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return subkegiatan.SubKegiatanOpdResponse{}, fmt.Errorf("subkegiatan dengan id %d tidak ditemukan", id)
		}
		return subkegiatan.SubKegiatanOpdResponse{}, err
	}
	kode, _ := service.SubKegiatanRepository.FindByKodeSubKegiatan(ctx, tx, result.KodeSubKegiatan)
	opd, _ := service.opdRepository.FindByKodeOpd(ctx, tx, result.KodeOpd)

	response := subkegiatan.SubKegiatanOpdResponse{
		Id:              result.Id,
		KodeSubkegiatan: kode.KodeSubKegiatan,
		NamaSubkegiatan: kode.NamaSubKegiatan,
		KodeOpd:         result.KodeOpd,
		NamaOpd:         opd.NamaOpd,
		Tahun:           result.Tahun,
	}

	return response, nil
}

func (service *SubKegiatanTerpilihServiceImpl) DeleteOpd(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	err = service.SubKegiatanTerpilihRepository.DeleteSubOpd(ctx, tx, id)
	if err != nil {
		return err
	}

	return nil
}
