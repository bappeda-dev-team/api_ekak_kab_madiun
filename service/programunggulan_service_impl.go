package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/programunggulan"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ProgramUnggulanServiceImpl struct {
	ProgramUnggulanRepository repository.ProgramUnggulanRepository
	DB                        *sql.DB
	Validate                  *validator.Validate
}

func NewProgramUnggulanServiceImpl(programUnggulanRepository repository.ProgramUnggulanRepository, db *sql.DB, validate *validator.Validate) *ProgramUnggulanServiceImpl {
	return &ProgramUnggulanServiceImpl{
		ProgramUnggulanRepository: programUnggulanRepository,
		DB:                        db,
		Validate:                  validate,
	}
}

func (service *ProgramUnggulanServiceImpl) Create(ctx context.Context, request programunggulan.ProgramUnggulanCreateRequest) (programunggulan.ProgramUnggulanResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return programunggulan.ProgramUnggulanResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return programunggulan.ProgramUnggulanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	kodeProgram := fmt.Sprintf("PRG-UNG-%s", uuid.New().String()[:6])

	programUnggulan := domain.ProgramUnggulan{
		NamaTagging:               request.NamaTagging,
		KodeProgramUnggulan:       kodeProgram,
		KeteranganProgramUnggulan: &request.KeteranganProgramUnggulan,
		Keterangan:                &request.Keterangan,
		TahunAwal:                 request.TahunAwal,
		TahunAkhir:                request.TahunAkhir,
	}

	result, err := service.ProgramUnggulanRepository.Create(ctx, tx, programUnggulan)
	if err != nil {
		return programunggulan.ProgramUnggulanResponse{}, err
	}

	return programunggulan.ProgramUnggulanResponse{
		Id:                        result.Id,
		NamaTagging:               result.NamaTagging,
		KodeProgramUnggulan:       result.KodeProgramUnggulan,
		KeteranganProgramUnggulan: result.KeteranganProgramUnggulan,
		Keterangan:                result.Keterangan,
		TahunAwal:                 result.TahunAwal,
		TahunAkhir:                result.TahunAkhir,
	}, nil
}

func (service *ProgramUnggulanServiceImpl) Update(ctx context.Context, request programunggulan.ProgramUnggulanUpdateRequest) (programunggulan.ProgramUnggulanResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return programunggulan.ProgramUnggulanResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return programunggulan.ProgramUnggulanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi data exists
	_, err = service.ProgramUnggulanRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return programunggulan.ProgramUnggulanResponse{}, err
	}

	programUnggulan := domain.ProgramUnggulan{
		Id:                        request.Id,
		NamaTagging:               request.NamaTagging,
		KeteranganProgramUnggulan: &request.KeteranganProgramUnggulan,
		Keterangan:                &request.Keterangan,
		TahunAwal:                 request.TahunAwal,
		TahunAkhir:                request.TahunAkhir,
	}

	result, err := service.ProgramUnggulanRepository.Update(ctx, tx, programUnggulan)
	if err != nil {
		return programunggulan.ProgramUnggulanResponse{}, err
	}

	return programunggulan.ProgramUnggulanResponse{
		Id:                        result.Id,
		NamaTagging:               result.NamaTagging,
		KodeProgramUnggulan:       result.KodeProgramUnggulan,
		KeteranganProgramUnggulan: result.KeteranganProgramUnggulan,
		Keterangan:                result.Keterangan,
		TahunAwal:                 result.TahunAwal,
		TahunAkhir:                result.TahunAkhir,
	}, nil
}

func (service *ProgramUnggulanServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi data exists
	_, err = service.ProgramUnggulanRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}

	return service.ProgramUnggulanRepository.Delete(ctx, tx, id)
}

func (service *ProgramUnggulanServiceImpl) FindById(ctx context.Context, id int) (programunggulan.ProgramUnggulanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return programunggulan.ProgramUnggulanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	result, err := service.ProgramUnggulanRepository.FindById(ctx, tx, id)
	if err != nil {
		return programunggulan.ProgramUnggulanResponse{}, err
	}

	return programunggulan.ProgramUnggulanResponse{
		Id:                        result.Id,
		NamaTagging:               result.NamaTagging,
		KodeProgramUnggulan:       result.KodeProgramUnggulan,
		KeteranganProgramUnggulan: result.KeteranganProgramUnggulan,
		Keterangan:                result.Keterangan,
		TahunAwal:                 result.TahunAwal,
		TahunAkhir:                result.TahunAkhir,
	}, nil
}

func (service *ProgramUnggulanServiceImpl) FindAll(ctx context.Context, tahunAwal string, tahunAkhir string) ([]programunggulan.ProgramUnggulanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	results, err := service.ProgramUnggulanRepository.FindAll(ctx, tx, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}

	var responses []programunggulan.ProgramUnggulanResponse
	for _, result := range results {
		responses = append(responses, programunggulan.ProgramUnggulanResponse{
			Id:                        result.Id,
			NamaTagging:               result.NamaTagging,
			KodeProgramUnggulan:       result.KodeProgramUnggulan,
			KeteranganProgramUnggulan: result.KeteranganProgramUnggulan,
			Keterangan:                result.Keterangan,
			TahunAwal:                 result.TahunAwal,
			TahunAkhir:                result.TahunAkhir,
			IsActive:                  result.IsActive,
		})
	}

	return responses, nil
}

func (service *ProgramUnggulanServiceImpl) FindByKodeProgramUnggulan(ctx context.Context, kodeProgramUnggulan string) (programunggulan.ProgramUnggulanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return programunggulan.ProgramUnggulanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	result, err := service.ProgramUnggulanRepository.FindByKodeProgramUnggulan(ctx, tx, kodeProgramUnggulan)
	if err != nil {
		return programunggulan.ProgramUnggulanResponse{}, err
	}

	return programunggulan.ProgramUnggulanResponse{
		Id:                        result.Id,
		NamaTagging:               result.NamaTagging,
		KodeProgramUnggulan:       result.KodeProgramUnggulan,
		KeteranganProgramUnggulan: result.KeteranganProgramUnggulan,
		Keterangan:                result.Keterangan,
		TahunAwal:                 result.TahunAwal,
		TahunAkhir:                result.TahunAkhir,
	}, nil
}

func (service *ProgramUnggulanServiceImpl) FindByTahun(ctx context.Context, tahun string) ([]programunggulan.ProgramUnggulanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi format tahun
	_, err = strconv.Atoi(tahun)
	if err != nil {
		return nil, errors.New("format tahun tidak valid")
	}

	results, err := service.ProgramUnggulanRepository.FindByTahun(ctx, tx, tahun)
	if err != nil {
		return nil, err
	}

	var responses []programunggulan.ProgramUnggulanResponse
	for _, result := range results {
		responses = append(responses, programunggulan.ProgramUnggulanResponse{
			Id:                        result.Id,
			NamaTagging:               result.NamaTagging,
			KodeProgramUnggulan:       result.KodeProgramUnggulan,
			KeteranganProgramUnggulan: result.KeteranganProgramUnggulan,
			Keterangan:                result.Keterangan,
			TahunAwal:                 result.TahunAwal,
			TahunAkhir:                result.TahunAkhir,
		})
	}

	return responses, nil
}

func (service *ProgramUnggulanServiceImpl) FindUnusedByTahun(ctx context.Context, tahun string) ([]programunggulan.ProgramUnggulanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi format tahun
	_, err = strconv.Atoi(tahun)
	if err != nil {
		return nil, errors.New("format tahun tidak valid")
	}

	results, err := service.ProgramUnggulanRepository.FindUnusedByTahun(ctx, tx, tahun)
	if err != nil {
		return nil, err
	}

	var responses []programunggulan.ProgramUnggulanResponse
	for _, result := range results {
		responses = append(responses, programunggulan.ProgramUnggulanResponse{
			Id:                        result.Id,
			NamaTagging:               result.NamaTagging,
			KodeProgramUnggulan:       result.KodeProgramUnggulan,
			KeteranganProgramUnggulan: result.KeteranganProgramUnggulan,
			Keterangan:                result.Keterangan,
			TahunAwal:                 result.TahunAwal,
			TahunAkhir:                result.TahunAkhir,
		})
	}

	return responses, nil
}

func (service *ProgramUnggulanServiceImpl) FindByIdTerkait(ctx context.Context, request programunggulan.FindByIdTerkaitRequest) ([]programunggulan.ProgramUnggulanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []programunggulan.ProgramUnggulanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	results, err := service.ProgramUnggulanRepository.FindByIdTerkait(ctx, tx, request.Ids)
	if err != nil {
		return []programunggulan.ProgramUnggulanResponse{}, err
	}

	var responses []programunggulan.ProgramUnggulanResponse
	for _, result := range results {
		responses = append(responses, programunggulan.ProgramUnggulanResponse{
			Id:                        result.Id,
			NamaTagging:               result.NamaTagging,
			KodeProgramUnggulan:       result.KodeProgramUnggulan,
			KeteranganProgramUnggulan: result.KeteranganProgramUnggulan,
			Keterangan:                result.Keterangan,
			TahunAwal:                 result.TahunAwal,
			TahunAkhir:                result.TahunAkhir,
		})
	}

	return responses, nil
}
