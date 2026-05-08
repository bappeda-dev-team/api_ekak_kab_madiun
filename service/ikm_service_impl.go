package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/indikator"
	"ekak_kabupaten_madiun/repository"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type IkmServiceImpl struct {
	IkmRepository repository.IkmRepository
	DB            *sql.DB
	Validate      *validator.Validate
}

func NewIkmServiceImpl(
	ikmRepository repository.IkmRepository,
	db *sql.DB,
	validate *validator.Validate,
) *IkmServiceImpl {
	return &IkmServiceImpl{
		IkmRepository: ikmRepository,
		DB:            db,
		Validate:      validate,
	}
}

func (s *IkmServiceImpl) FindById(
	ctx context.Context,
	ikmId string,
) (indikator.IkmResponse, error) {

	tx, err := s.DB.Begin()
	if err != nil {
		return indikator.IkmResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	data, err := s.IkmRepository.FindById(ctx, tx, ikmId)
	if err != nil {
		return indikator.IkmResponse{}, err
	}

	return toIkmResponse(data), nil
}

func (s *IkmServiceImpl) FindAllByPeriode(
	ctx context.Context,
	tahunAwal, tahunAkhir string,
) ([]indikator.IkmResponse, error) {

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	list, err := s.IkmRepository.
		FindAllByPeriode(ctx, tx, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}

	result := make([]indikator.IkmResponse, 0)

	for _, d := range list {
		result = append(result, toIkmResponse(d))
	}

	return result, nil
}

func (s *IkmServiceImpl) Create(
	ctx context.Context,
	request indikator.IkmRequest,
) (indikator.IkmResponse, error) {

	err := s.Validate.Struct(request)
	if err != nil {
		return indikator.IkmResponse{}, err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return indikator.IkmResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	id := uuid.New().String()

	domainData := domain.IndikatorIkm{
		Id:                  id,
		Indikator:           request.Indikator,
		KodeBidangUrusan:    request.KodeBidangUrusan,
		NamaBidangUrusan:    request.NamaBidangUrusan,
		IsActive:            true,
		DefinisiOperasional: request.DefinisiOperasional,
		RumusPerhitungan:    request.RumusPerhitungan,
		SumberData:          request.SumberData,
		Jenis:               request.Jenis,
		TahunAwal:           request.TahunAwal,
		TahunAkhir:          request.TahunAkhir,
	}

	res, err := s.IkmRepository.Create(ctx, tx, domainData)
	if err != nil {
		return indikator.IkmResponse{}, err
	}

	return toIkmResponse(res), nil
}

func (s *IkmServiceImpl) Update(
	ctx context.Context,
	request indikator.IkmRequest,
	ikmId string,
) (indikator.IkmResponse, error) {

	err := s.Validate.Struct(request)
	if err != nil {
		return indikator.IkmResponse{}, err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return indikator.IkmResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	exists, err := s.IkmRepository.ExistsById(ctx, tx, ikmId)
	if err != nil {
		return indikator.IkmResponse{}, err
	}
	if !exists {
		return indikator.IkmResponse{}, errors.New("ikm tidak ditemukan")
	}

	domainData := domain.IndikatorIkm{
		Indikator:           request.Indikator,
		KodeBidangUrusan:    request.KodeBidangUrusan,
		NamaBidangUrusan:    request.NamaBidangUrusan,
		IsActive:            request.IsActive,
		DefinisiOperasional: request.DefinisiOperasional,
		RumusPerhitungan:    request.RumusPerhitungan,
		SumberData:          request.SumberData,
		Jenis:               request.Jenis,
		TahunAwal:           request.TahunAwal,
		TahunAkhir:          request.TahunAkhir,
	}

	res, err := s.IkmRepository.Update(ctx, tx, domainData, ikmId)
	if err != nil {
		return indikator.IkmResponse{}, err
	}

	return toIkmResponse(res), nil
}

func (s *IkmServiceImpl) Delete(
	ctx context.Context,
	ikmId string,
) error {

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	exists, err := s.IkmRepository.ExistsById(ctx, tx, ikmId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("ikm tidak ditemukan")
	}

	return s.IkmRepository.Delete(ctx, tx, ikmId)
}

func toIkmResponse(d domain.IndikatorIkm) indikator.IkmResponse {
	return indikator.IkmResponse{
		Id:                  d.Id,
		Indikator:           d.Indikator,
		KodeBidangUrusan:    d.KodeBidangUrusan,
		NamaBidangUrusan:    d.NamaBidangUrusan,
		IsActive:            d.IsActive,
		DefinisiOperasional: d.DefinisiOperasional,
		RumusPerhitungan:    d.RumusPerhitungan,
		SumberData:          d.SumberData,
		Jenis:               d.Jenis,
		TahunAwal:           d.TahunAwal,
		TahunAkhir:          d.TahunAkhir,
	}
}
