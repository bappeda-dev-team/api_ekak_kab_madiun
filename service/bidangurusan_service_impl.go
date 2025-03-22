package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"ekak_kabupaten_madiun/model/web/bidangurusanresponse"
	"ekak_kabupaten_madiun/repository"
	"fmt"
)

type BidangUrusanServiceImpl struct {
	BidangUrusanRepository repository.BidangUrusanRepository
	DB                     *sql.DB
}

func NewBidangUrusanServiceImpl(bidangUrusanRepository repository.BidangUrusanRepository, db *sql.DB) *BidangUrusanServiceImpl {
	return &BidangUrusanServiceImpl{
		BidangUrusanRepository: bidangUrusanRepository,
		DB:                     db,
	}
}

func (service *BidangUrusanServiceImpl) Create(ctx context.Context, request bidangurusanresponse.BidangUrusanCreateRequest) (bidangurusanresponse.BidangUrusanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return bidangurusanresponse.BidangUrusanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	uuId := fmt.Sprintf("BID-URU-%s", request.KodeBidangUrusan)

	bidangurusan := domainmaster.BidangUrusan{
		Id:               uuId,
		KodeBidangUrusan: request.KodeBidangUrusan,
		NamaBidangUrusan: request.NamaBidangUrusan,
		Tahun:            request.Tahun,
	}

	bidangurusan = service.BidangUrusanRepository.Create(ctx, tx, bidangurusan)

	return bidangurusanresponse.BidangUrusanResponse{
		Id:               bidangurusan.Id,
		KodeBidangUrusan: bidangurusan.KodeBidangUrusan,
		NamaBidangUrusan: bidangurusan.NamaBidangUrusan,
		Tahun:            bidangurusan.Tahun,
	}, nil
}

func (service *BidangUrusanServiceImpl) Update(ctx context.Context, request bidangurusanresponse.BidangUrusanUpdateRequest) (bidangurusanresponse.BidangUrusanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return bidangurusanresponse.BidangUrusanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	bidangurusan, err := service.BidangUrusanRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return bidangurusanresponse.BidangUrusanResponse{}, err
	}

	bidangurusan = domainmaster.BidangUrusan{
		Id:               request.Id,
		KodeBidangUrusan: request.KodeBidangUrusan,
		NamaBidangUrusan: request.NamaBidangUrusan,
		Tahun:            request.Tahun,
	}

	bidangurusan = service.BidangUrusanRepository.Update(ctx, tx, bidangurusan)

	return bidangurusanresponse.BidangUrusanResponse{
		Id:               bidangurusan.Id,
		KodeBidangUrusan: bidangurusan.KodeBidangUrusan,
		NamaBidangUrusan: bidangurusan.NamaBidangUrusan,
		Tahun:            bidangurusan.Tahun,
	}, nil
}

func (service *BidangUrusanServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	return service.BidangUrusanRepository.Delete(ctx, tx, id)
}

func (service *BidangUrusanServiceImpl) FindById(ctx context.Context, id string) (bidangurusanresponse.BidangUrusanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return bidangurusanresponse.BidangUrusanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	bidangurusan, err := service.BidangUrusanRepository.FindById(ctx, tx, id)
	if err != nil {
		return bidangurusanresponse.BidangUrusanResponse{}, err
	}

	return bidangurusanresponse.BidangUrusanResponse{
		Id:               bidangurusan.Id,
		KodeBidangUrusan: bidangurusan.KodeBidangUrusan,
		NamaBidangUrusan: bidangurusan.NamaBidangUrusan,
		Tahun:            bidangurusan.Tahun,
	}, nil
}

func (service *BidangUrusanServiceImpl) FindAll(ctx context.Context) ([]bidangurusanresponse.BidangUrusanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []bidangurusanresponse.BidangUrusanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	bidangurusans, err := service.BidangUrusanRepository.FindAll(ctx, tx)
	if err != nil {
		return []bidangurusanresponse.BidangUrusanResponse{}, err
	}

	var bidangurusanResponses []bidangurusanresponse.BidangUrusanResponse
	for _, bidangurusan := range bidangurusans {
		bidangurusanResponses = append(bidangurusanResponses, bidangurusanresponse.BidangUrusanResponse{
			Id:               bidangurusan.Id,
			KodeBidangUrusan: bidangurusan.KodeBidangUrusan,
			NamaBidangUrusan: bidangurusan.NamaBidangUrusan,
			Tahun:            bidangurusan.Tahun,
		})
	}

	return bidangurusanResponses, nil
}

func (service *BidangUrusanServiceImpl) FindByKodeOpd(ctx context.Context, kodeOpd string) ([]bidangurusanresponse.BidangUrusanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []bidangurusanresponse.BidangUrusanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	bidangUrusans, err := service.BidangUrusanRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return []bidangurusanresponse.BidangUrusanResponse{}, err
	}

	var bidangUrusanResponses []bidangurusanresponse.BidangUrusanResponse
	for _, bidangUrusan := range bidangUrusans {
		bidangUrusanResponses = append(bidangUrusanResponses, bidangurusanresponse.BidangUrusanResponse{
			KodeBidangUrusan: bidangUrusan.KodeBidangUrusan,
			NamaBidangUrusan: bidangUrusan.NamaBidangUrusan,
		})
	}

	return bidangUrusanResponses, nil
}
