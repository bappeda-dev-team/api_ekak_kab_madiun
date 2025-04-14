package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/rincianbelanja"
	"ekak_kabupaten_madiun/repository"
	"errors"
)

type RincianBelanjaServiceImpl struct {
	rincianBelanjaRepository repository.RincianBelanjaRepository
	DB                       *sql.DB
}

func NewRincianBelanjaServiceImpl(rincianBelanjaRepository repository.RincianBelanjaRepository, DB *sql.DB) *RincianBelanjaServiceImpl {
	return &RincianBelanjaServiceImpl{
		rincianBelanjaRepository: rincianBelanjaRepository,
		DB:                       DB,
	}
}

func (service *RincianBelanjaServiceImpl) Create(ctx context.Context, request rincianbelanja.RincianBelanjaCreateRequest) (rincianbelanja.RencanaAksiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return rincianbelanja.RencanaAksiResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi request
	if request.RenaksiId == "" {
		return rincianbelanja.RencanaAksiResponse{}, errors.New("renaksi_id tidak boleh kosong")
	}
	if request.Anggaran < 0 {
		return rincianbelanja.RencanaAksiResponse{}, errors.New("anggaran tidak boleh negatif")
	}

	// Konversi request ke domain model
	rincianBelanja := domain.RincianBelanja{
		RenaksiId: request.RenaksiId,
		Anggaran:  int64(request.Anggaran),
	}

	// Simpan ke database
	result, err := service.rincianBelanjaRepository.Create(ctx, tx, rincianBelanja)
	if err != nil {
		return rincianbelanja.RencanaAksiResponse{}, err
	}

	// Konversi domain model ke response
	response := rincianbelanja.RencanaAksiResponse{
		RenaksiId: result.RenaksiId,
		Renaksi:   result.Renaksi,
		Anggaran:  int(result.Anggaran),
	}

	return response, nil
}

func (service *RincianBelanjaServiceImpl) Update(ctx context.Context, request rincianbelanja.RincianBelanjaUpdateRequest) (rincianbelanja.RencanaAksiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return rincianbelanja.RencanaAksiResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah data exists
	existing, err := service.rincianBelanjaRepository.FindByRenaksiId(ctx, tx, request.RenaksiId)
	if err != nil {
		return rincianbelanja.RencanaAksiResponse{}, err
	}
	if existing.RenaksiId == "" {
		return rincianbelanja.RencanaAksiResponse{}, errors.New("rincian belanja tidak ditemukan")
	}

	// Konversi request ke domain model
	rincianBelanja := domain.RincianBelanja{
		RenaksiId: request.RenaksiId,
		Anggaran:  int64(request.Anggaran),
	}

	// Update ke database
	result, err := service.rincianBelanjaRepository.Update(ctx, tx, rincianBelanja)
	if err != nil {
		return rincianbelanja.RencanaAksiResponse{}, err
	}

	// Konversi domain model ke response
	response := rincianbelanja.RencanaAksiResponse{
		RenaksiId: result.RenaksiId,
		Renaksi:   result.Renaksi,
		Anggaran:  int(result.Anggaran),
	}

	return response, nil
}
func (service *RincianBelanjaServiceImpl) FindRincianBelanjaAsn(ctx context.Context, pegawaiId string, tahun string) []rincianbelanja.RincianBelanjaAsnResponse {
	tx, err := service.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer helper.CommitOrRollback(tx)

	rincianBelanjaList, err := service.rincianBelanjaRepository.FindRincianBelanjaAsn(ctx, tx, pegawaiId, tahun)
	if err != nil {
		panic(err)
	}

	var responses []rincianbelanja.RincianBelanjaAsnResponse
	for _, rb := range rincianBelanjaList {
		var rencanaKinerjaResponses []rincianbelanja.RincianBelanjaResponse

		for _, rk := range rb.RencanaKinerja {
			var rencanaAksiResponses []rincianbelanja.RencanaAksiResponse

			// Pastikan slice RencanaAksi tidak nil
			if rk.RencanaAksi != nil {
				for _, ra := range rk.RencanaAksi {
					rencanaAksiResponses = append(rencanaAksiResponses, rincianbelanja.RencanaAksiResponse{
						RenaksiId: ra.RenaksiId,
						Renaksi:   ra.Renaksi,
						Anggaran:  int(ra.Anggaran),
					})
				}
			}

			// Jika tidak ada rencana aksi, inisialisasi dengan slice kosong
			if rencanaAksiResponses == nil {
				rencanaAksiResponses = make([]rincianbelanja.RencanaAksiResponse, 0)
			}

			rencanaKinerjaResponses = append(rencanaKinerjaResponses, rincianbelanja.RincianBelanjaResponse{
				RencanaKinerja: rk.RencanaKinerja,
				RencanaAksi:    rencanaAksiResponses,
			})
		}

		responses = append(responses, rincianbelanja.RincianBelanjaAsnResponse{
			PegawaiId:       rb.PegawaiId,
			NamaPegawai:     rb.NamaPegawai,
			KodeSubkegiatan: rb.KodeSubkegiatan,
			NamaSubkegiatan: rb.NamaSubkegiatan,
			TotalAnggaran:   rb.TotalAnggaran,
			RincianBelanja:  rencanaKinerjaResponses,
		})
	}

	return responses
}
