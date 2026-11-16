package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/rencanaaksi"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid" // Tambahkan impor ini
)

type RencanaAksiServiceImpl struct {
	rencanaAksiRepository            repository.RencanaAksiRepository
	DB                               *sql.DB
	Validate                         *validator.Validate
	pelaksanaanRencanaAksiRepository repository.PelaksanaanRencanaAksiRepository
}

func NewRencanaAksiServiceImpl(rencanaAksiRepository repository.RencanaAksiRepository, DB *sql.DB, validate *validator.Validate, pelaksanaanRencanaAksiRepository repository.PelaksanaanRencanaAksiRepository) *RencanaAksiServiceImpl {
	return &RencanaAksiServiceImpl{
		rencanaAksiRepository:            rencanaAksiRepository,
		DB:                               DB,
		Validate:                         validate,
		pelaksanaanRencanaAksiRepository: pelaksanaanRencanaAksiRepository,
	}
}

func (service *RencanaAksiServiceImpl) Create(ctx context.Context, request rencanaaksi.RencanaAksiCreateRequest) (rencanaaksi.RencanaAksiResponse, error) {
	// Validasi request
	err := service.Validate.Struct(request)
	if err != nil {
		return rencanaaksi.RencanaAksiResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return rencanaaksi.RencanaAksiResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	if request.Urutan <= 0 {
		return rencanaaksi.RencanaAksiResponse{}, fmt.Errorf("urutan tidak boleh kurang dari atau sama dengan 0")
	}

	// Buat UUID baru dengan format yang diinginkan
	uuId := fmt.Sprintf("RENAKSI-REKIN-%s", uuid.New().String()[:5])

	// Buat objek domain.RencanaAksi dari request
	rencanaAksi := domain.RencanaAksi{
		Id:               uuId,
		RencanaKinerjaId: request.RencanaKinerjaId,
		KodeOpd:          request.KodeOpd,
		Urutan:           request.Urutan,
		NamaRencanaAksi:  request.NamaRencanaAksi,
	}

	// Panggil repository untuk membuat rencana aksi
	result, err := service.rencanaAksiRepository.Create(ctx, tx, rencanaAksi)
	if err != nil {
		tx.Rollback()
		return rencanaaksi.RencanaAksiResponse{}, err
	}

	// Buat response
	response := rencanaaksi.RencanaAksiResponse{
		Id:               result.Id,
		RencanaKinerjaId: result.RencanaKinerjaId,
		KodeOpd:          result.KodeOpd,
		Urutan:           result.Urutan,
		NamaRencanaAksi:  result.NamaRencanaAksi,
	}

	return response, nil
}

func (service *RencanaAksiServiceImpl) Update(ctx context.Context, request rencanaaksi.RencanaAksiUpdateRequest) (rencanaaksi.RencanaAksiResponse, error) {
	// Validasi request
	err := service.Validate.Struct(request)
	if err != nil {
		return rencanaaksi.RencanaAksiResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return rencanaaksi.RencanaAksiResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	if request.Urutan <= 0 {
		return rencanaaksi.RencanaAksiResponse{}, fmt.Errorf("urutan tidak boleh kurang dari atau sama dengan 0")
	}

	// Cek apakah rencana aksi dengan ID tersebut ada
	existingRencanaAksi, err := service.rencanaAksiRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return rencanaaksi.RencanaAksiResponse{}, fmt.Errorf("rencana aksi dengan ID %s tidak ditemukan", request.Id)
	}

	// Update data rencana aksi
	existingRencanaAksi.Urutan = request.Urutan
	existingRencanaAksi.NamaRencanaAksi = request.NamaRencanaAksi

	// Panggil repository untuk mengupdate rencana aksi
	updatedRencanaAksi, err := service.rencanaAksiRepository.Update(ctx, tx, existingRencanaAksi)
	if err != nil {
		return rencanaaksi.RencanaAksiResponse{}, err
	}

	// Buat response
	response := rencanaaksi.RencanaAksiResponse{
		Id:               updatedRencanaAksi.Id,
		RencanaKinerjaId: updatedRencanaAksi.RencanaKinerjaId,
		Urutan:           updatedRencanaAksi.Urutan,
		NamaRencanaAksi:  updatedRencanaAksi.NamaRencanaAksi,
	}

	return response, nil
}

func (service *RencanaAksiServiceImpl) FindAll(ctx context.Context, rencanaKinerjaId string) ([]rencanaaksi.RencanaAksiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Panggil repository untuk mendapatkan semua rencana aksi
	rencanaAksiList, err := service.rencanaAksiRepository.FindAll(ctx, tx, rencanaKinerjaId)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil daftar rencana aksi: %v", err)
	}

	// Konversi domain.RencanaAksi menjadi rencanaaksi.RencanaAksiResponse
	var rencanaAksiResponses []rencanaaksi.RencanaAksiResponse
	for _, rencanaAksi := range rencanaAksiList {
		// Ambil data pelaksanaan rencana aksi untuk setiap rencana aksi
		pelaksanaanList, err := service.pelaksanaanRencanaAksiRepository.FindByRencanaAksiId(ctx, tx, rencanaAksi.Id)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil data pelaksanaan rencana aksi: %v", err)
		}

		host := os.Getenv("host")
		port := os.Getenv("port")

		buttonActions := []web.ActionButton{
			{
				NameAction: "Find Id Rencana Aksi",
				Method:     "GET",
				Url:        fmt.Sprintf("%s:%s/detail-rencana_aksi/:rencanaaksiId", host, port),
			},
			{
				NameAction: "Update Rencana Aksi",
				Method:     "PUT",
				Url:        fmt.Sprintf("%s:%s/rencana_aksi/update/rencanaaksi/:rencanaaksiId", host, port),
			},
			{
				NameAction: "Delete Rencana Aksi",
				Method:     "DELETE",
				Url:        fmt.Sprintf("%s:%s/rencana_aksi/delete/rencanaaksi/:rencanaaksiId", host, port),
			},
			{
				NameAction: "Create Pelaksanaan Rencana Aksi",
				Method:     "POST",
				Url:        fmt.Sprintf("%s:%s/pelaksanaan_rencana_aksi/create/:rencanaAksiId", host, port),
			},
		}

		// Konversi domain.PelaksanaanRencanaAksi menjadi rencanaaksi.PelaksanaanRencanaAksiResponse
		var pelaksanaanResponses []rencanaaksi.PelaksanaanRencanaAksiResponse
		jumlahBobot := 0
		for _, pelaksanaan := range pelaksanaanList {

			buttonActions := []web.ActionButton{
				{
					NameAction: "Find Id Pelaksanaan Rencana Aksi",
					Method:     "GET",
					Url:        fmt.Sprintf("%s:%s/pelaksanaan_rencana_aksi/detail/:id", host, port),
				},
				{
					NameAction: "Update Pelaksanaan Rencana Aksi",
					Method:     "PUT",
					Url:        fmt.Sprintf("%s:%s/pelaksanaan_rencana_aksi/update/:pelaksanaanRencanaAksiId", host, port),
				},
				{
					NameAction: "Delete Pelaksanaan Rencana Aksi",
					Method:     "DELETE",
					Url:        fmt.Sprintf("%s:%s/pelaksanaan_rencana_aksi/delete/:id", host, port),
				},
			}
			pelaksanaanResponses = append(pelaksanaanResponses, rencanaaksi.PelaksanaanRencanaAksiResponse{
				Id:            pelaksanaan.Id,
				RencanaAksiId: pelaksanaan.RencanaAksiId,
				Bulan:         pelaksanaan.Bulan,
				Bobot:         pelaksanaan.Bobot,
				Action:        buttonActions,
			})
			jumlahBobot += pelaksanaan.Bobot
		}

		rencanaAksiResponses = append(rencanaAksiResponses, rencanaaksi.RencanaAksiResponse{
			Id:                     rencanaAksi.Id,
			RencanaKinerjaId:       rencanaAksi.RencanaKinerjaId,
			KodeOpd:                rencanaAksi.KodeOpd,
			Urutan:                 rencanaAksi.Urutan,
			NamaRencanaAksi:        rencanaAksi.NamaRencanaAksi,
			PelaksanaanRencanaAksi: pelaksanaanResponses,
			JumlahBobot:            jumlahBobot,
			Action:                 buttonActions,
		})
	}

	return rencanaAksiResponses, nil
}

func (service *RencanaAksiServiceImpl) FindById(ctx context.Context, id string) (rencanaaksi.RencanaAksiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return rencanaaksi.RencanaAksiResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Panggil repository untuk mendapatkan rencana aksi berdasarkan ID
	rencanaAksi, err := service.rencanaAksiRepository.FindById(ctx, tx, id)
	if err != nil {
		return rencanaaksi.RencanaAksiResponse{}, fmt.Errorf("gagal mengambil rencana aksi: %v", err)
	}

	// Ambil data pelaksanaan rencana aksi untuk rencana aksi ini
	pelaksanaanList, err := service.pelaksanaanRencanaAksiRepository.FindByRencanaAksiId(ctx, tx, rencanaAksi.Id)
	if err != nil {
		return rencanaaksi.RencanaAksiResponse{}, fmt.Errorf("gagal mengambil data pelaksanaan rencana aksi: %v", err)
	}

	// Konversi domain.PelaksanaanRencanaAksi menjadi rencanaaksi.PelaksanaanRencanaAksiResponse
	var pelaksanaanResponses []rencanaaksi.PelaksanaanRencanaAksiResponse
	jumlahBobot := 0
	for _, pelaksanaan := range pelaksanaanList {
		pelaksanaanResponses = append(pelaksanaanResponses, rencanaaksi.PelaksanaanRencanaAksiResponse{
			Id:            pelaksanaan.Id,
			RencanaAksiId: pelaksanaan.RencanaAksiId,
			Bulan:         pelaksanaan.Bulan,
			Bobot:         pelaksanaan.Bobot,
		})
		jumlahBobot += pelaksanaan.Bobot
	}

	// Buat response
	response := rencanaaksi.RencanaAksiResponse{
		Id:                     rencanaAksi.Id,
		RencanaKinerjaId:       rencanaAksi.RencanaKinerjaId,
		KodeOpd:                rencanaAksi.KodeOpd,
		Urutan:                 rencanaAksi.Urutan,
		NamaRencanaAksi:        rencanaAksi.NamaRencanaAksi,
		PelaksanaanRencanaAksi: pelaksanaanResponses,
		JumlahBobot:            jumlahBobot,
	}

	return response, nil
}

func (service *RencanaAksiServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Periksa apakah rencana aksi dengan ID tersebut ada
	_, err = service.rencanaAksiRepository.FindById(ctx, tx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("rencana aksi dengan ID %s tidak ditemukan", id)
		}
		return fmt.Errorf("gagal memeriksa rencana aksi: %v", err)
	}

	// Panggil repository untuk menghapus rencana aksi
	err = service.rencanaAksiRepository.Delete(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus rencana aksi: %v", err)
	}

	return nil
}

func (service *RencanaAksiServiceImpl) FindByRekin(ctx context.Context, rekinId string) ([]rencanaaksi.RencanaAksiResponse, error) {
	if strings.TrimSpace(rekinId) == "" {
		return nil, fmt.Errorf("rencana_kinerja_id wajib diisi")
	}
	return service.FindAll(ctx, rekinId)
}

func (service *RencanaAksiServiceImpl) FindByPegawai(ctx context.Context, nip string, kodeOpd string, tahun string) ([]rencanaaksi.RencanaAksiPegawaiResponse, error) {
	if strings.TrimSpace(nip) == "" {
		return nil, fmt.Errorf("nip wajib diisi")
	}
	if strings.TrimSpace(kodeOpd) == "" {
		return nil, fmt.Errorf("kode_opd wajib diisi")
	}
	if strings.TrimSpace(tahun) == "" {
		return nil, fmt.Errorf("tahun wajib diisi")
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	rencanaAksiList, err := service.rencanaAksiRepository.FindByPegawai(ctx, tx, nip, kodeOpd, tahun)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil rencana aksi pegawai: %v", err)
	}

	order := make([]string, 0)
	groupMap := make(map[string]*rencanaaksi.RencanaAksiPegawaiResponse)

	for _, rencanaAksi := range rencanaAksiList {
		group, exists := groupMap[rencanaAksi.RencanaKinerjaId]
		if !exists {
			group = &rencanaaksi.RencanaAksiPegawaiResponse{
				RencanaKinerjaId:   rencanaAksi.RencanaKinerjaId,
				NamaRencanaKinerja: rencanaAksi.NamaRencanaKinerja,
				Nip:                rencanaAksi.PegawaiId,
				KodeOpd:            rencanaAksi.KodeOpd,
				Tahun:              rencanaAksi.Tahun,
				RencanaAksi:        []rencanaaksi.RencanaAksiResponse{},
			}
			groupMap[rencanaAksi.RencanaKinerjaId] = group
			order = append(order, rencanaAksi.RencanaKinerjaId)
		}

		pelaksanaanResponses := make([]rencanaaksi.PelaksanaanRencanaAksiResponse, 0, len(rencanaAksi.Pelaksanaan))
		jumlahBobot := 0
		for _, pelaksanaan := range rencanaAksi.Pelaksanaan {
			pelaksanaanResponses = append(pelaksanaanResponses, rencanaaksi.PelaksanaanRencanaAksiResponse{
				Id:            pelaksanaan.Id,
				RencanaAksiId: pelaksanaan.RencanaAksiId,
				Bulan:         pelaksanaan.Bulan,
				Bobot:         pelaksanaan.Bobot,
			})
			jumlahBobot += pelaksanaan.Bobot
		}

		group.RencanaAksi = append(group.RencanaAksi, rencanaaksi.RencanaAksiResponse{
			Id:                     rencanaAksi.Id,
			RencanaKinerjaId:       rencanaAksi.RencanaKinerjaId,
			KodeOpd:                rencanaAksi.KodeOpd,
			Urutan:                 rencanaAksi.Urutan,
			NamaRencanaAksi:        rencanaAksi.NamaRencanaAksi,
			PelaksanaanRencanaAksi: pelaksanaanResponses,
			JumlahBobot:            jumlahBobot,
		})
		group.JumlahBobot += jumlahBobot
	}

	results := make([]rencanaaksi.RencanaAksiPegawaiResponse, 0, len(order))
	for _, rekinId := range order {
		results = append(results, *groupMap[rekinId])
	}

	return results, nil
}
