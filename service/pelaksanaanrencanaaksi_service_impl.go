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
	"math/rand"
)

type PelaksanaanRencanaAksiServiceImpl struct {
	PelaksanaanRencanaAksiRepository repository.PelaksanaanRencanaAksiRepository
	RencanaAksiRepository            repository.RencanaAksiRepository
	DB                               *sql.DB
}

func NewPelaksanaanRencanaAksiServiceImpl(pelaksanaanRencanaAksiRepository repository.PelaksanaanRencanaAksiRepository, rencanaAksiRepository repository.RencanaAksiRepository, DB *sql.DB) *PelaksanaanRencanaAksiServiceImpl {
	return &PelaksanaanRencanaAksiServiceImpl{
		PelaksanaanRencanaAksiRepository: pelaksanaanRencanaAksiRepository,
		RencanaAksiRepository:            rencanaAksiRepository,
		DB:                               DB,
	}
}

func (service *PelaksanaanRencanaAksiServiceImpl) Create(ctx context.Context, request rencanaaksi.PelaksanaanRencanaAksiCreateRequest) (rencanaaksi.PelaksanaanRencanaAksiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	// Validasi bulan
	if request.Bulan < 1 || request.Bulan > 12 {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("bulan harus antara 1 dan 12")
	}

	//validasi bobot
	if request.Bobot < 1 {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("bobot harus diisi lebih dari 0")
	}

	// PERBAIKAN: Lock rencana_aksi untuk mencegah race condition
	rencanaAksi, err := service.RencanaAksiRepository.FindById(ctx, tx, request.RencanaAksiId)
	if err != nil {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("gagal mendapatkan RencanaAksi: %v", err)
	}

	// Check dengan SELECT FOR UPDATE
	exists, err := service.PelaksanaanRencanaAksiRepository.ExistsByRencanaAksiIdAndBulan(ctx, tx, request.RencanaAksiId, request.Bulan)
	if err != nil {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("gagal memeriksa keberadaan bulan: %v", err)
	}
	if exists {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("bulan %d sudah ada untuk rencana aksi ini", request.Bulan)
	}

	// Periksa total bobot yang sudah ada untuk rencana kinerja ini
	totalBobot, err := service.RencanaAksiRepository.GetTotalBobotForRencanaKinerja(ctx, tx, rencanaAksi.RencanaKinerjaId)
	if err != nil {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("gagal mendapatkan total bobot: %v", err)
	}

	sisaBobot := 100 - totalBobot

	if request.Bobot > sisaBobot {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("bobot melebihi sisa yang tersedia. Sisa bobot: %d", sisaBobot)
	}

	newID := fmt.Sprintf("%05d", rand.Intn(100000))
	formattedID := fmt.Sprintf("PLKSN-RENAKSI-%s", newID)

	pelaksanaanRencanaAksi := domain.PelaksanaanRencanaAksi{
		Id:            formattedID,
		RencanaAksiId: request.RencanaAksiId,
		Bobot:         request.Bobot,
		Bulan:         request.Bulan,
	}

	result, err := service.PelaksanaanRencanaAksiRepository.Create(ctx, tx, pelaksanaanRencanaAksi)
	if err != nil {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("gagal menyimpan pelaksanaan rencana aksi: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("gagal melakukan commit transaksi: %v", err)
	}

	response := rencanaaksi.PelaksanaanRencanaAksiResponse{
		Id:            result.Id,
		RencanaAksiId: result.RencanaAksiId,
		Bobot:         result.Bobot,
		Bulan:         result.Bulan,
		BobotAvail:    sisaBobot - result.Bobot,
	}

	return response, nil
}

func (service *PelaksanaanRencanaAksiServiceImpl) Update(ctx context.Context, request rencanaaksi.PelaksanaanRencanaAksiUpdateRequest) (rencanaaksi.PelaksanaanRencanaAksiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("gagal mendapatkan PelaksanaanRencanaAksi: %v", err)
	}
	defer tx.Rollback()

	// Validasi bulan
	if request.Bulan < 1 || request.Bulan > 12 {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, web.NewBadRequestError("bulan harus antara 1 dan 12")
	}

	//validasi bobot
	if request.Bobot < 1 {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("bobot harus diisi lebih dari 0")
	}

	// Dapatkan PelaksanaanRencanaAksi yang ada
	existingPelaksanaan, err := service.PelaksanaanRencanaAksiRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return rencanaaksi.PelaksanaanRencanaAksiResponse{}, web.NewNotFoundError(fmt.Sprintf("PelaksanaanRencanaAksi dengan ID %s tidak ditemukan", request.Id))
		}
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("gagal mendapatkan PelaksanaanRencanaAksi: %v", err)
	}

	// Periksa apakah bulan sudah ada untuk rencana aksi ini (kecuali untuk bulan yang sama dengan yang sedang diupdate)
	if request.Bulan != existingPelaksanaan.Bulan {
		exists, err := service.PelaksanaanRencanaAksiRepository.ExistsByRencanaAksiIdAndBulan(ctx, tx, existingPelaksanaan.RencanaAksiId, request.Bulan)
		if err != nil {
			return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("gagal mendapatkan PelaksanaanRencanaAksi: %v", err)
		}
		if exists {
			return rencanaaksi.PelaksanaanRencanaAksiResponse{}, web.NewBadRequestError(fmt.Sprintf("bulan %d sudah ada untuk rencana aksi ini", request.Bulan))
		}
	}

	// Dapatkan RencanaAksi untuk mendapatkan RencanaKinerjaId
	rencanaAksi, err := service.RencanaAksiRepository.FindById(ctx, tx, existingPelaksanaan.RencanaAksiId)
	if err != nil {
		if err == sql.ErrNoRows {
			return rencanaaksi.PelaksanaanRencanaAksiResponse{}, web.NewNotFoundError(fmt.Sprintf("RencanaAksi dengan ID %s tidak ditemukan", existingPelaksanaan.RencanaAksiId))
		}
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("gagal mendapatkan PelaksanaanRencanaAksi: %v", err)
	}

	// Periksa total bobot yang sudah ada untuk rencana kinerja ini
	totalBobot, err := service.RencanaAksiRepository.GetTotalBobotForRencanaKinerja(ctx, tx, rencanaAksi.RencanaKinerjaId)
	if err != nil {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("gagal mendapatkan PelaksanaanRencanaAksi: %v", err)
	}

	// Hitung sisa bobot yang tersedia (tambahkan bobot yang ada saat ini karena akan diganti)
	sisaBobot := 100 - totalBobot + existingPelaksanaan.Bobot

	// Periksa apakah bobot baru melebihi sisa bobot
	if request.Bobot > sisaBobot {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, web.NewBadRequestError(fmt.Sprintf("bobot melebihi sisa yang tersedia. Sisa bobot: %d", sisaBobot))
	}

	// Update pelaksanaan rencana aksi
	updatedPelaksanaan := domain.PelaksanaanRencanaAksi{
		Id:            request.Id,
		RencanaAksiId: existingPelaksanaan.RencanaAksiId,
		Bobot:         request.Bobot,
		Bulan:         request.Bulan,
	}

	result, err := service.PelaksanaanRencanaAksiRepository.Update(ctx, tx, updatedPelaksanaan)
	if err != nil {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("gagal mendapatkan PelaksanaanRencanaAksi: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("gagal mendapatkan PelaksanaanRencanaAksi: %v", err)
	}

	response := rencanaaksi.PelaksanaanRencanaAksiResponse{
		Id:            result.Id,
		RencanaAksiId: result.RencanaAksiId,
		Bobot:         result.Bobot,
		Bulan:         result.Bulan,
		BobotAvail:    sisaBobot - result.Bobot,
	}

	return response, nil
}

func (service *PelaksanaanRencanaAksiServiceImpl) FindById(ctx context.Context, id string) (rencanaaksi.PelaksanaanRencanaAksiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	pelaksanaan, err := service.PelaksanaanRencanaAksiRepository.FindById(ctx, tx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("pelaksanaan rencana aksi dengan ID %s tidak ditemukan", id)
		}
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("gagal mendapatkan pelaksanaan rencana aksi: %v", err)
	}

	// Dapatkan RencanaAksi untuk mendapatkan RencanaKinerjaId
	rencanaAksi, err := service.RencanaAksiRepository.FindById(ctx, tx, pelaksanaan.RencanaAksiId)
	if err != nil {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("gagal mendapatkan RencanaAksi: %v", err)
	}

	// Hitung total bobot dan sisa bobot
	totalBobot, err := service.RencanaAksiRepository.GetTotalBobotForRencanaKinerja(ctx, tx, rencanaAksi.RencanaKinerjaId)
	if err != nil {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("gagal mendapatkan total bobot: %v", err)
	}

	sisaBobot := 100 - totalBobot + pelaksanaan.Bobot

	err = tx.Commit()
	if err != nil {
		return rencanaaksi.PelaksanaanRencanaAksiResponse{}, fmt.Errorf("gagal melakukan commit transaksi: %v", err)
	}

	response := rencanaaksi.PelaksanaanRencanaAksiResponse{
		Id:            pelaksanaan.Id,
		RencanaAksiId: pelaksanaan.RencanaAksiId,
		Bobot:         pelaksanaan.Bobot,
		Bulan:         pelaksanaan.Bulan,
		BobotAvail:    sisaBobot - pelaksanaan.Bobot,
	}

	return response, nil
}

func (service *PelaksanaanRencanaAksiServiceImpl) FindByRencanaAksiId(ctx context.Context, rencanaAksiId string) ([]rencanaaksi.PelaksanaanRencanaAksiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	pelaksanaanList, err := service.PelaksanaanRencanaAksiRepository.FindByRencanaAksiId(ctx, tx, rencanaAksiId)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan daftar pelaksanaan rencana aksi: %v", err)
	}

	// Dapatkan RencanaAksi untuk mendapatkan RencanaKinerjaId
	rencanaAksi, err := service.RencanaAksiRepository.FindById(ctx, tx, rencanaAksiId)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan RencanaAksi: %v", err)
	}

	// Hitung total bobot dan sisa bobot
	totalBobot, err := service.RencanaAksiRepository.GetTotalBobotForRencanaKinerja(ctx, tx, rencanaAksi.RencanaKinerjaId)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan total bobot: %v", err)
	}

	sisaBobot := 100 - totalBobot

	var responseList []rencanaaksi.PelaksanaanRencanaAksiResponse
	for _, pelaksanaan := range pelaksanaanList {
		response := rencanaaksi.PelaksanaanRencanaAksiResponse{
			Id:            pelaksanaan.Id,
			RencanaAksiId: pelaksanaan.RencanaAksiId,
			Bobot:         pelaksanaan.Bobot,
			Bulan:         pelaksanaan.Bulan,
			BobotAvail:    sisaBobot + pelaksanaan.Bobot,
		}
		responseList = append(responseList, response)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("gagal melakukan commit transaksi: %v", err)
	}

	return responseList, nil
}

func (service *PelaksanaanRencanaAksiServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	err = service.PelaksanaanRencanaAksiRepository.Delete(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus pelaksanaan rencana aksi: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("gagal melakukan commit transaksi: %v", err)
	}

	return nil
}
