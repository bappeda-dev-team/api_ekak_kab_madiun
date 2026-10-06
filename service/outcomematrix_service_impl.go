package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/outcomematrix"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type OutcomeMatrixServiceImpl struct {
	OutcomeMatrixRepository repository.OutcomeMatrixRepository
	DB                      *sql.DB
	Validate                *validator.Validate
}

func NewOutcomeMatrixServiceImpl(
	outcomeMatrixRepository repository.OutcomeMatrixRepository,
	db *sql.DB,
	validate *validator.Validate,
) *OutcomeMatrixServiceImpl {
	return &OutcomeMatrixServiceImpl{
		OutcomeMatrixRepository: outcomeMatrixRepository,
		DB:                      db,
		Validate:                validate,
	}
}

func toOutcomeMatrixResponse(data domain.OutcomeMatrix) outcomematrix.OutcomeMatrixResponse {
	return outcomematrix.OutcomeMatrixResponse{
		Id:              data.Id,
		KodeSubkegiatan: data.KodeSubkegiatan,
		Kode:            data.Kode,
		Outcome:         data.Outcome,
		CreatedAt:       data.CreatedAt,
		UpdatedAt:       data.UpdatedAt,
	}
}

func (s *OutcomeMatrixServiceImpl) Create(ctx context.Context, request outcomematrix.OutcomeMatrixCreateRequest) (outcomematrix.OutcomeMatrixResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return outcomematrix.OutcomeMatrixResponse{}, err
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return outcomematrix.OutcomeMatrixResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	result, err := s.OutcomeMatrixRepository.Create(ctx, tx, domain.OutcomeMatrix{
		KodeSubkegiatan: strings.TrimSpace(request.KodeSubkegiatan),
		Kode:            strings.TrimSpace(request.Kode),
		Outcome:         request.Outcome,
	})
	if err != nil {
		return outcomematrix.OutcomeMatrixResponse{}, err
	}
	return toOutcomeMatrixResponse(result), nil
}

func (s *OutcomeMatrixServiceImpl) Update(ctx context.Context, request outcomematrix.OutcomeMatrixUpdateRequest) (outcomematrix.OutcomeMatrixResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return outcomematrix.OutcomeMatrixResponse{}, err
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return outcomematrix.OutcomeMatrixResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	if _, err := s.OutcomeMatrixRepository.FindById(ctx, tx, request.Id); err != nil {
		return outcomematrix.OutcomeMatrixResponse{}, fmt.Errorf("outcome matrix id %d tidak ditemukan", request.Id)
	}
	result, err := s.OutcomeMatrixRepository.Update(ctx, tx, domain.OutcomeMatrix{
		Id:              request.Id,
		KodeSubkegiatan: strings.TrimSpace(request.KodeSubkegiatan),
		Kode:            strings.TrimSpace(request.Kode),
		Outcome:         request.Outcome,
	})
	if err != nil {
		return outcomematrix.OutcomeMatrixResponse{}, err
	}
	return toOutcomeMatrixResponse(result), nil
}

func (s *OutcomeMatrixServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)
	if _, err := s.OutcomeMatrixRepository.FindById(ctx, tx, id); err != nil {
		return fmt.Errorf("outcome matrix id %d tidak ditemukan", id)
	}
	return s.OutcomeMatrixRepository.Delete(ctx, tx, id)
}

func (s *OutcomeMatrixServiceImpl) FindById(ctx context.Context, id int) (outcomematrix.OutcomeMatrixResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return outcomematrix.OutcomeMatrixResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	result, err := s.OutcomeMatrixRepository.FindById(ctx, tx, id)
	if err != nil {
		return outcomematrix.OutcomeMatrixResponse{}, fmt.Errorf("outcome matrix id %d tidak ditemukan", id)
	}
	return toOutcomeMatrixResponse(result), nil
}

func (s *OutcomeMatrixServiceImpl) FindAll(ctx context.Context, kode, kodeSubkegiatan string) ([]outcomematrix.OutcomeMatrixResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	list, err := s.OutcomeMatrixRepository.FindAll(ctx, tx, strings.TrimSpace(kode), strings.TrimSpace(kodeSubkegiatan))
	if err != nil {
		return nil, err
	}
	responses := make([]outcomematrix.OutcomeMatrixResponse, 0, len(list))
	for _, item := range list {
		responses = append(responses, toOutcomeMatrixResponse(item))
	}
	return responses, nil
}

func (s *OutcomeMatrixServiceImpl) UpsertBatch(ctx context.Context, requests []outcomematrix.OutcomeMatrixBatchItemRequest) ([]outcomematrix.OutcomeMatrixResponse, error) {
	if len(requests) == 0 {
		return []outcomematrix.OutcomeMatrixResponse{}, fmt.Errorf("batch outcome matrix tidak boleh kosong")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	responses := make([]outcomematrix.OutcomeMatrixResponse, 0, len(requests))
	for i, req := range requests {
		if err := s.Validate.Struct(req); err != nil {
			return nil, fmt.Errorf("item batch ke-%d tidak valid: %w", i+1, err)
		}
		data := domain.OutcomeMatrix{
			Id:              req.Id,
			KodeSubkegiatan: strings.TrimSpace(req.KodeSubkegiatan),
			Kode:            strings.TrimSpace(req.Kode),
			Outcome:         req.Outcome,
		}
		var result domain.OutcomeMatrix
		if req.Id > 0 {
			if _, err := s.OutcomeMatrixRepository.FindById(ctx, tx, req.Id); err != nil {
				return nil, fmt.Errorf("item batch ke-%d: outcome matrix id %d tidak ditemukan", i+1, req.Id)
			}
			result, err = s.OutcomeMatrixRepository.Update(ctx, tx, data)
		} else {
			result, err = s.OutcomeMatrixRepository.Create(ctx, tx, data)
		}
		if err != nil {
			return nil, err
		}
		responses = append(responses, toOutcomeMatrixResponse(result))
	}
	return responses, nil
}
