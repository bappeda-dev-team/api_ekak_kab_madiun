package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"ekak_kabupaten_madiun/model/web/lembaga"
	"ekak_kabupaten_madiun/repository"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type LembagaServiceImpl struct {
	LembagaRepository repository.LembagaRepository
	DB                *sql.DB
	Validator         *validator.Validate
}

func NewLembagaServiceImpl(lembagaRepository repository.LembagaRepository, DB *sql.DB, validator *validator.Validate) *LembagaServiceImpl {
	return &LembagaServiceImpl{
		LembagaRepository: lembagaRepository,
		DB:                DB,
		Validator:         validator,
	}
}

func (service *LembagaServiceImpl) Create(ctx context.Context, request lembaga.LembagaCreateRequest) (lembaga.LembagaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return lembaga.LembagaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	err = service.Validator.Struct(request)
	if err != nil {
		return lembaga.LembagaResponse{}, err
	}

	// Generate ID unik dengan prefix
	uuid := uuid.New().String()
	lembagaId := fmt.Sprintf("LMBG-%v", uuid[:4])

	// Buat objek domain lembaga
	lembagaDomain := domainmaster.Lembaga{
		Id:                 lembagaId,
		KodeLembaga:        request.KodeLembaga,
		NamaLembaga:        request.NamaLembaga,
		NamaKepalaPemda:    request.NamaKepalaPemda,
		NipKepalaPemda:     request.NipKepalaPemda,
		JabatanKepalaPemda: request.JabatanKepalaPemda,
		IsActive:           true, // Tambahkan respons isactive true
	}

	result := service.LembagaRepository.Create(ctx, tx, lembagaDomain)

	response := lembaga.LembagaResponse{
		Id:                 result.Id,
		KodeLembaga:        result.KodeLembaga,
		NamaLembaga:        result.NamaLembaga,
		NamaKepalaPemda:    result.NamaKepalaPemda,
		NipKepalaPemda:     result.NipKepalaPemda,
		JabatanKepalaPemda: result.JabatanKepalaPemda,
		IsActive:           result.IsActive, // Tambahkan respons isactive true
	}

	return response, nil
}

func (service *LembagaServiceImpl) Update(ctx context.Context, request lembaga.LembagaUpdateRequest) (lembaga.LembagaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return lembaga.LembagaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	err = service.Validator.Struct(request)
	if err != nil {
		return lembaga.LembagaResponse{}, err
	}

	// Cek apakah data exists
	_, err = service.LembagaRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return lembaga.LembagaResponse{}, err
	}

	lembagaDomain := domainmaster.Lembaga{
		Id:                 request.Id,
		KodeLembaga:        request.KodeLembaga,
		NamaLembaga:        request.NamaLembaga,
		NamaKepalaPemda:    request.NamaKepalaPemda,
		NipKepalaPemda:     request.NipKepalaPemda,
		JabatanKepalaPemda: request.JabatanKepalaPemda,
		IsActive:           request.IsActive,
	}

	result := service.LembagaRepository.Update(ctx, tx, lembagaDomain)

	response := lembaga.LembagaResponse{
		Id:                 result.Id,
		KodeLembaga:        result.KodeLembaga,
		NamaLembaga:        result.NamaLembaga,
		NamaKepalaPemda:    result.NamaKepalaPemda,
		NipKepalaPemda:     result.NipKepalaPemda,
		JabatanKepalaPemda: result.JabatanKepalaPemda,
		IsActive:           result.IsActive,
	}

	return response, nil
}

func (service *LembagaServiceImpl) FindById(ctx context.Context, id string) (lembaga.LembagaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return lembaga.LembagaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	result, err := service.LembagaRepository.FindById(ctx, tx, id)
	if err != nil {
		return lembaga.LembagaResponse{}, err
	}

	response := lembaga.LembagaResponse{
		Id:                 result.Id,
		KodeLembaga:        result.KodeLembaga,
		NamaLembaga:        result.NamaLembaga,
		NamaKepalaPemda:    result.NamaKepalaPemda,
		NipKepalaPemda:     result.NipKepalaPemda,
		JabatanKepalaPemda: result.JabatanKepalaPemda,
		IsActive:           result.IsActive,
	}

	return response, nil
}

func (service *LembagaServiceImpl) FindAll(ctx context.Context) ([]lembaga.LembagaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []lembaga.LembagaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	result, err := service.LembagaRepository.FindAll(ctx, tx)
	if err != nil {
		return []lembaga.LembagaResponse{}, err
	}

	response := []lembaga.LembagaResponse{}
	for _, value := range result {
		response = append(response, lembaga.LembagaResponse{
			Id:                 value.Id,
			KodeLembaga:        value.KodeLembaga,
			NamaLembaga:        value.NamaLembaga,
			NamaKepalaPemda:    value.NamaKepalaPemda,
			NipKepalaPemda:     value.NipKepalaPemda,
			JabatanKepalaPemda: value.JabatanKepalaPemda,
			IsActive:           value.IsActive,
		})
	}
	return response, nil
}

func (service *LembagaServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	return service.LembagaRepository.Delete(ctx, tx, id)
}
