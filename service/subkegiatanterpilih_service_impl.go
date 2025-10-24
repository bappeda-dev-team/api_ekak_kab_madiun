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
	"strings"

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
	kodeSubKegiatanArray := []string{request.KodeSubKegiatan}

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

	// Proses setiap kdde
	for _, kodeSubKegiatan := range kodeSubKegiatanArray {
		// Cek apakah usulan dengan ID yang diberikan ada
		existingSubKegiatan, err := service.SubKegiatanRepository.FindByKodeSubKegiatan(ctx, tx, kodeSubKegiatan)
		if err != nil {
			return nil, fmt.Errorf("subkegiatan dengan kode %s tidak ditemukan: %v", kodeSubKegiatan, err)
		}

		// Cek apakah usulan sudah memiliki rekin_id
		if existingSubKegiatan.RekinId != "" {
			return nil, fmt.Errorf("subkegiatan dengan kode %s sudah memiliki rencana kinerja", kodeSubKegiatan)
		}

		// Update rekin_id dan status
		err = service.SubKegiatanTerpilihRepository.CreateRekin(ctx, tx, existingSubKegiatan.Id, request.RekinId, existingSubKegiatan.KodeSubKegiatan)
		if err != nil {
			return nil, fmt.Errorf("gagal mengupdate rekin untuk subkegiatan %s: %v", existingSubKegiatan.KodeSubKegiatan, err)
		}

		// Ambil data usulan yang sudah diupdate
		updatedSubKegiatan, err := service.SubKegiatanRepository.FindByKodeSubKegiatan(ctx, tx, kodeSubKegiatan)
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

// func (service *SubKegiatanTerpilihServiceImpl) CreateOpd(ctx context.Context, request subkegiatan.SubKegiatanOpdCreateRequest) (subkegiatan.SubKegiatanOpdResponse, error) {
// 	tx, err := service.DB.Begin()
// 	if err != nil {
// 		return subkegiatan.SubKegiatanOpdResponse{}, err
// 	}
// 	defer helper.CommitOrRollback(tx)

// 	//cek subkegiatan
// 	kode, err := service.SubKegiatanRepository.FindByKodeSubKegiatan(ctx, tx, request.KodeSubkegiatan)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return subkegiatan.SubKegiatanOpdResponse{}, fmt.Errorf("kode subkegiatan %s tidak ditemukan dalam database", request.KodeSubkegiatan)
// 		}
// 		return subkegiatan.SubKegiatanOpdResponse{}, fmt.Errorf("terjadi kesalahan saat mencari data kode subkegiatan: %v", err)
// 	}

// 	// Cek OPD
// 	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOpd)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return subkegiatan.SubKegiatanOpdResponse{}, fmt.Errorf("OPD dengan kode %s tidak ditemukan dalam database", request.KodeOpd)
// 		}
// 		return subkegiatan.SubKegiatanOpdResponse{}, fmt.Errorf("terjadi kesalahan saat mencari data OPD: %v", err)
// 	}

// 	domain := domain.SubKegiatanOpd{
// 		KodeSubKegiatan: kode.KodeSubKegiatan,
// 		KodeOpd:         opd.KodeOpd,
// 		Tahun:           request.Tahun,
// 	}

// 	result, err := service.SubKegiatanTerpilihRepository.CreateOPD(ctx, tx, domain)
// 	if err != nil {
// 		return subkegiatan.SubKegiatanOpdResponse{}, err
// 	}

// 	response := subkegiatan.SubKegiatanOpdResponse{
// 		Id:              result.Id,
// 		KodeSubkegiatan: result.KodeSubKegiatan,
// 		KodeOpd:         result.KodeOpd,
// 		Tahun:           result.Tahun,
// 	}

// 	return response, nil
// }

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

func (service *SubKegiatanTerpilihServiceImpl) FindAllSubkegiatanByBidangUrusanOpd(ctx context.Context, kodeOpd string) ([]subkegiatan.SubKegiatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi format kode OPD
	if !isValidKodeOpd(kodeOpd) {
		return nil, fmt.Errorf("format kode OPD tidak valid")
	}

	// Cek apakah OPD ada
	_, err = service.opdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("OPD dengan kode %s tidak ditemukan", kodeOpd)
		}
		return nil, err
	}

	result, err := service.SubKegiatanTerpilihRepository.FindAllSubkegiatanByBidangUrusanOpd(ctx, tx, kodeOpd)
	if err != nil {
		return nil, err
	}

	var responses []subkegiatan.SubKegiatanResponse
	for _, sub := range result {
		response := subkegiatan.SubKegiatanResponse{
			KodeSubKegiatan: sub.KodeSubKegiatan,
			NamaSubKegiatan: sub.NamaSubKegiatan,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// Helper function untuk validasi format kode OPD
func isValidKodeOpd(kodeOpd string) bool {
	parts := strings.Split(kodeOpd, ".")
	return len(parts) == 8 // Format: 5.01.5.05.0.00.01.0000
}

func (service *SubKegiatanTerpilihServiceImpl) CreateOpdMultiple(ctx context.Context, request subkegiatan.SubKegiatanOpdMultipleCreateRequest) (subkegiatan.SubKegiatanOpdMultipleResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return subkegiatan.SubKegiatanOpdMultipleResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Cek OPD sekali saja
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOpd)
	if err != nil {
		if err == sql.ErrNoRows {
			return subkegiatan.SubKegiatanOpdMultipleResponse{}, fmt.Errorf("OPD dengan kode %s tidak ditemukan dalam database", request.KodeOpd)
		}
		return subkegiatan.SubKegiatanOpdMultipleResponse{}, fmt.Errorf("terjadi kesalahan saat mencari data OPD: %v", err)
	}

	var successItems []subkegiatan.SubKegiatanOpdResponse
	var skippedItems []subkegiatan.SubKegiatanOpdResponse
	successCount := 0
	skippedCount := 0

	// Proses setiap subkegiatan
	for _, kodeSubkegiatan := range request.KodeSubkegiatan {
		// Cek apakah subkegiatan ada di database
		subkegiatanData, err := service.SubKegiatanRepository.FindByKodeSubKegiatan(ctx, tx, kodeSubkegiatan)
		if err != nil {
			if err == sql.ErrNoRows {
				// Skip jika subkegiatan tidak ditemukan
				skippedItems = append(skippedItems, subkegiatan.SubKegiatanOpdResponse{
					KodeSubkegiatan: kodeSubkegiatan,
					KodeOpd:         request.KodeOpd,
					NamaOpd:         opd.NamaOpd,
					Tahun:           request.Tahun,
					Status:          "not_found",
				})
				skippedCount++
				continue
			}
			return subkegiatan.SubKegiatanOpdMultipleResponse{}, fmt.Errorf("terjadi kesalahan saat mencari data subkegiatan %s: %v", kodeSubkegiatan, err)
		}

		// Cek apakah kombinasi sudah ada
		exists, err := service.SubKegiatanTerpilihRepository.CheckExists(ctx, tx, kodeSubkegiatan, request.KodeOpd, request.Tahun)
		if err != nil {
			return subkegiatan.SubKegiatanOpdMultipleResponse{}, fmt.Errorf("terjadi kesalahan saat mengecek duplikasi: %v", err)
		}

		if exists {
			// Skip jika sudah ada
			skippedItems = append(skippedItems, subkegiatan.SubKegiatanOpdResponse{
				KodeSubkegiatan: kodeSubkegiatan,
				NamaSubkegiatan: subkegiatanData.NamaSubKegiatan,
				KodeOpd:         request.KodeOpd,
				NamaOpd:         opd.NamaOpd,
				Tahun:           request.Tahun,
				Status:          "already_exists",
			})
			skippedCount++
			continue
		}

		// Buat domain untuk insert
		domain := domain.SubKegiatanOpd{
			KodeSubKegiatan: kodeSubkegiatan,
			KodeOpd:         request.KodeOpd,
			Tahun:           request.Tahun,
		}

		// Insert ke database
		result, err := service.SubKegiatanTerpilihRepository.CreateOPD(ctx, tx, domain)
		if err != nil {
			return subkegiatan.SubKegiatanOpdMultipleResponse{}, fmt.Errorf("terjadi kesalahan saat menyimpan subkegiatan %s: %v", kodeSubkegiatan, err)
		}

		// Tambahkan ke success items
		successItems = append(successItems, subkegiatan.SubKegiatanOpdResponse{
			Id:              result.Id,
			KodeSubkegiatan: result.KodeSubKegiatan,
			NamaSubkegiatan: subkegiatanData.NamaSubKegiatan,
			KodeOpd:         result.KodeOpd,
			NamaOpd:         opd.NamaOpd,
			Tahun:           result.Tahun,
			Status:          "created",
		})
		successCount++
	}

	opd, _ = service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOpd)

	// Buat response
	response := subkegiatan.SubKegiatanOpdMultipleResponse{
		SuccessCount:   successCount,
		TotalRequested: len(request.KodeSubkegiatan),
		SkippedCount:   skippedCount,
		SuccessItems:   successItems,
		SkippedItems:   skippedItems,
		Message:        fmt.Sprintf("Berhasil menambahkan %d dari %d subkegiatan untuk OPD %s", successCount, len(request.KodeSubkegiatan), opd.NamaOpd),
	}

	return response, nil
}
