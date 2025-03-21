package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"ekak_kabupaten_madiun/model/web/bidangurusanresponse"
	"ekak_kabupaten_madiun/model/web/urusanrespon"
	"ekak_kabupaten_madiun/repository"
	"fmt"
)

type UrusanServiceImpl struct {
	UrusanRepository repository.UrusanRepository
	DB               *sql.DB
}

func NewUrusanServiceImpl(urusanRepository repository.UrusanRepository, DB *sql.DB) *UrusanServiceImpl {
	return &UrusanServiceImpl{
		UrusanRepository: urusanRepository,
		DB:               DB,
	}
}

func (service *UrusanServiceImpl) Create(ctx context.Context, request urusanrespon.UrusanCreateRequest) (urusanrespon.UrusanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return urusanrespon.UrusanResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	uuId := fmt.Sprintf("URU-%s", request.KodeUrusan)

	domainUrusan := domainmaster.Urusan{
		Id:         uuId,
		KodeUrusan: request.KodeUrusan,
		NamaUrusan: request.NamaUrusan,
	}

	urusan, err := service.UrusanRepository.Create(ctx, tx, domainUrusan)
	if err != nil {
		return urusanrespon.UrusanResponse{}, err
	}

	return urusanrespon.UrusanResponse{
		Id:         urusan.Id,
		KodeUrusan: urusan.KodeUrusan,
		NamaUrusan: urusan.NamaUrusan,
	}, nil
}

func (service *UrusanServiceImpl) Update(ctx context.Context, request urusanrespon.UrusanUpdateRequest) (urusanrespon.UrusanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return urusanrespon.UrusanResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	domainUrusan := domainmaster.Urusan{
		Id:         request.Id,
		KodeUrusan: request.KodeUrusan,
		NamaUrusan: request.NamaUrusan,
	}

	urusan, err := service.UrusanRepository.Update(ctx, tx, domainUrusan)
	if err != nil {
		return urusanrespon.UrusanResponse{}, err
	}

	return urusanrespon.UrusanResponse{
		Id:         urusan.Id,
		KodeUrusan: urusan.KodeUrusan,
		NamaUrusan: urusan.NamaUrusan,
	}, nil
}

func (service *UrusanServiceImpl) FindById(ctx context.Context, id string) (urusanrespon.UrusanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return urusanrespon.UrusanResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	urusan, err := service.UrusanRepository.FindById(ctx, tx, id)
	if err != nil {
		return urusanrespon.UrusanResponse{}, err
	}

	return urusanrespon.UrusanResponse{
		Id:         urusan.Id,
		KodeUrusan: urusan.KodeUrusan,
		NamaUrusan: urusan.NamaUrusan,
	}, nil
}

func (service *UrusanServiceImpl) FindAll(ctx context.Context) ([]urusanrespon.UrusanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []urusanrespon.UrusanResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	urusans, err := service.UrusanRepository.FindAll(ctx, tx)
	if err != nil {
		return []urusanrespon.UrusanResponse{}, err
	}

	var urusanResponses []urusanrespon.UrusanResponse
	for _, urusan := range urusans {
		urusanResponses = append(urusanResponses, urusanrespon.UrusanResponse{
			Id:         urusan.Id,
			KodeUrusan: urusan.KodeUrusan,
			NamaUrusan: urusan.NamaUrusan,
		})
	}

	return urusanResponses, nil
}

func (service *UrusanServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah urusan dengan ID tersebut ada
	_, err = service.UrusanRepository.FindById(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("urusan dengan id %s tidak ditemukan", id)
	}

	return service.UrusanRepository.Delete(ctx, tx, id)
}

func (service *UrusanServiceImpl) FindByKodeOpd(ctx context.Context, kodeOpd string) ([]urusanrespon.UrusanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []urusanrespon.UrusanResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	urusans, err := service.UrusanRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return []urusanrespon.UrusanResponse{}, err
	}

	var urusanResponses []urusanrespon.UrusanResponse
	for _, urusan := range urusans {
		urusanResponses = append(urusanResponses, urusanrespon.UrusanResponse{
			Id:         urusan.Id,
			KodeUrusan: urusan.KodeUrusan,
			NamaUrusan: urusan.NamaUrusan,
		})
	}

	return urusanResponses, nil
}

func (service *UrusanServiceImpl) FindUrusanAndBidangByKodeOpd(ctx context.Context, kodeOpd string) ([]urusanrespon.UrusanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Debug: cetak kode OPD
	fmt.Printf("Service - Mencari dengan kode OPD: %s\n", kodeOpd)

	// Panggil repository untuk mendapatkan data
	urusans, err := service.UrusanRepository.FindUrusanAndBidangByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data urusan dan bidang: %v", err)
	}

	// Debug: cetak hasil dari repository
	fmt.Printf("Service - Data dari repository:\n")
	for _, u := range urusans {
		fmt.Printf("Urusan: %s - %s\n", u.KodeUrusan, u.NamaUrusan)
		fmt.Printf("Jumlah Bidang Urusan: %d\n", len(u.BidangUrusan))
		for _, b := range u.BidangUrusan {
			fmt.Printf("  Bidang: %s - %s\n", b.KodeBidangUrusan, b.NamaBidangUrusan)
		}
	}

	// Konversi ke response
	var response []urusanrespon.UrusanResponse
	for _, urusan := range urusans {
		bidangResponses := make([]bidangurusanresponse.BidangUrusanResponse, 0)

		// Debug: cetak bidang urusan sebelum konversi
		fmt.Printf("Service - Converting bidang urusan for urusan %s\n", urusan.KodeUrusan)

		for _, bidang := range urusan.BidangUrusan {
			bidangResponses = append(bidangResponses, bidangurusanresponse.BidangUrusanResponse{
				KodeBidangUrusan: bidang.KodeBidangUrusan,
				NamaBidangUrusan: bidang.NamaBidangUrusan,
			})
			// Debug: cetak setiap bidang urusan yang dikonversi
			fmt.Printf("  Added bidang: %s - %s\n", bidang.KodeBidangUrusan, bidang.NamaBidangUrusan)
		}

		response = append(response, urusanrespon.UrusanResponse{
			Id:           urusan.Id,
			KodeUrusan:   urusan.KodeUrusan,
			NamaUrusan:   urusan.NamaUrusan,
			BidangUrusan: bidangResponses,
		})
	}

	// Debug: cetak hasil akhir
	fmt.Printf("Service - Final response:\n")
	for _, r := range response {
		fmt.Printf("Response Urusan: %s - %s\n", r.KodeUrusan, r.NamaUrusan)
		fmt.Printf("Response Jumlah Bidang: %d\n", len(r.BidangUrusan))
	}

	return response, nil
}
