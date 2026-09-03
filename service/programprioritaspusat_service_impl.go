package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/programprioritaspusat"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ProgramPrioritasPusatServiceImpl struct {
	ProgramPrioritasPusatRepository repository.ProgramPrioritasPusatRepository
	DB                        *sql.DB
	Validate                  *validator.Validate
}

func NewProgramPrioritasPusatServiceImpl(programPrioritasPusatRepository repository.ProgramPrioritasPusatRepository, db *sql.DB, validate *validator.Validate) *ProgramPrioritasPusatServiceImpl {
	return &ProgramPrioritasPusatServiceImpl{
		ProgramPrioritasPusatRepository: programPrioritasPusatRepository,
		DB:                        db,
		Validate:                  validate,
	}
}

func (service *ProgramPrioritasPusatServiceImpl) Create(ctx context.Context, request programprioritaspusat.ProgramPrioritasPusatCreateRequest) (programprioritaspusat.ProgramPrioritasPusatResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return programprioritaspusat.ProgramPrioritasPusatResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return programprioritaspusat.ProgramPrioritasPusatResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	kodeProgram := fmt.Sprintf("PRG-UNG-%s", uuid.New().String()[:6])

	programPrioritasPusat := domain.ProgramPrioritasPusat{
		NamaTagging:               request.NamaTagging,
		KodeProgramPrioritasPusat:       kodeProgram,
		KeteranganProgramPrioritasPusat: &request.KeteranganProgramPrioritasPusat,
		Keterangan:                &request.Keterangan,
		TahunAwal:                 request.TahunAwal,
		TahunAkhir:                request.TahunAkhir,
	}

	result, err := service.ProgramPrioritasPusatRepository.Create(ctx, tx, programPrioritasPusat)
	if err != nil {
		return programprioritaspusat.ProgramPrioritasPusatResponse{}, err
	}

	return programprioritaspusat.ProgramPrioritasPusatResponse{
		Id:                        result.Id,
		NamaTagging:               result.NamaTagging,
		KodeProgramPrioritasPusat:       result.KodeProgramPrioritasPusat,
		KeteranganProgramPrioritasPusat: result.KeteranganProgramPrioritasPusat,
		Keterangan:                result.Keterangan,
		TahunAwal:                 result.TahunAwal,
		TahunAkhir:                result.TahunAkhir,
	}, nil
}

func (service *ProgramPrioritasPusatServiceImpl) Update(ctx context.Context, request programprioritaspusat.ProgramPrioritasPusatUpdateRequest) (programprioritaspusat.ProgramPrioritasPusatResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return programprioritaspusat.ProgramPrioritasPusatResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return programprioritaspusat.ProgramPrioritasPusatResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi data exists
	_, err = service.ProgramPrioritasPusatRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return programprioritaspusat.ProgramPrioritasPusatResponse{}, err
	}

	programPrioritasPusat := domain.ProgramPrioritasPusat{
		Id:                        request.Id,
		NamaTagging:               request.NamaTagging,
		KeteranganProgramPrioritasPusat: &request.KeteranganProgramPrioritasPusat,
		Keterangan:                &request.Keterangan,
		TahunAwal:                 request.TahunAwal,
		TahunAkhir:                request.TahunAkhir,
	}

	result, err := service.ProgramPrioritasPusatRepository.Update(ctx, tx, programPrioritasPusat)
	if err != nil {
		return programprioritaspusat.ProgramPrioritasPusatResponse{}, err
	}

	updateData, err := service.ProgramPrioritasPusatRepository.FindById(ctx, tx, result.Id)
	if err != nil {
		return programprioritaspusat.ProgramPrioritasPusatResponse{}, err
	}

	return programprioritaspusat.ProgramPrioritasPusatResponse{
		Id:                        updateData.Id,
		NamaTagging:               updateData.NamaTagging,
		KodeProgramPrioritasPusat:       updateData.KodeProgramPrioritasPusat,
		KeteranganProgramPrioritasPusat: updateData.KeteranganProgramPrioritasPusat,
		Keterangan:                updateData.Keterangan,
		TahunAwal:                 updateData.TahunAwal,
		TahunAkhir:                updateData.TahunAkhir,
	}, nil
}

func (service *ProgramPrioritasPusatServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi data exists
	_, err = service.ProgramPrioritasPusatRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}

	return service.ProgramPrioritasPusatRepository.Delete(ctx, tx, id)
}

func (service *ProgramPrioritasPusatServiceImpl) FindById(ctx context.Context, id int) (programprioritaspusat.ProgramPrioritasPusatResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return programprioritaspusat.ProgramPrioritasPusatResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	result, err := service.ProgramPrioritasPusatRepository.FindById(ctx, tx, id)
	if err != nil {
		return programprioritaspusat.ProgramPrioritasPusatResponse{}, err
	}

	return programprioritaspusat.ProgramPrioritasPusatResponse{
		Id:                        result.Id,
		NamaTagging:               result.NamaTagging,
		KodeProgramPrioritasPusat:       result.KodeProgramPrioritasPusat,
		KeteranganProgramPrioritasPusat: result.KeteranganProgramPrioritasPusat,
		Keterangan:                result.Keterangan,
		TahunAwal:                 result.TahunAwal,
		TahunAkhir:                result.TahunAkhir,
	}, nil
}

func (service *ProgramPrioritasPusatServiceImpl) FindAll(ctx context.Context, tahunAwal string, tahunAkhir string) ([]programprioritaspusat.ProgramPrioritasPusatResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	results, err := service.ProgramPrioritasPusatRepository.FindAll(ctx, tx, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}

	var responses []programprioritaspusat.ProgramPrioritasPusatResponse
	for _, result := range results {
		responses = append(responses, programprioritaspusat.ProgramPrioritasPusatResponse{
			Id:                        result.Id,
			NamaTagging:               result.NamaTagging,
			KodeProgramPrioritasPusat:       result.KodeProgramPrioritasPusat,
			KeteranganProgramPrioritasPusat: result.KeteranganProgramPrioritasPusat,
			Keterangan:                result.Keterangan,
			TahunAwal:                 result.TahunAwal,
			TahunAkhir:                result.TahunAkhir,
			IsActive:                  result.IsActive,
		})
	}

	return responses, nil
}

func (service *ProgramPrioritasPusatServiceImpl) FindByKodeProgramPrioritasPusat(ctx context.Context, kodeProgramPrioritasPusat string) (programprioritaspusat.ProgramPrioritasPusatResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return programprioritaspusat.ProgramPrioritasPusatResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	result, err := service.ProgramPrioritasPusatRepository.FindByKodeProgramPrioritasPusat(ctx, tx, kodeProgramPrioritasPusat)
	if err != nil {
		return programprioritaspusat.ProgramPrioritasPusatResponse{}, err
	}

	return programprioritaspusat.ProgramPrioritasPusatResponse{
		Id:                        result.Id,
		NamaTagging:               result.NamaTagging,
		KodeProgramPrioritasPusat:       result.KodeProgramPrioritasPusat,
		KeteranganProgramPrioritasPusat: result.KeteranganProgramPrioritasPusat,
		Keterangan:                result.Keterangan,
		TahunAwal:                 result.TahunAwal,
		TahunAkhir:                result.TahunAkhir,
	}, nil
}

func (service *ProgramPrioritasPusatServiceImpl) FindByTahun(ctx context.Context, tahun string) ([]programprioritaspusat.ProgramPrioritasPusatResponse, error) {
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

	results, err := service.ProgramPrioritasPusatRepository.FindByTahun(ctx, tx, tahun)
	if err != nil {
		return nil, err
	}

	var responses []programprioritaspusat.ProgramPrioritasPusatResponse
	for _, result := range results {
		responses = append(responses, programprioritaspusat.ProgramPrioritasPusatResponse{
			Id:                        result.Id,
			NamaTagging:               result.NamaTagging,
			KodeProgramPrioritasPusat:       result.KodeProgramPrioritasPusat,
			KeteranganProgramPrioritasPusat: result.KeteranganProgramPrioritasPusat,
			Keterangan:                result.Keterangan,
			TahunAwal:                 result.TahunAwal,
			TahunAkhir:                result.TahunAkhir,
		})
	}

	return responses, nil
}

func (service *ProgramPrioritasPusatServiceImpl) FindUnusedByTahun(ctx context.Context, tahun string) ([]programprioritaspusat.ProgramPrioritasPusatResponse, error) {
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

	results, err := service.ProgramPrioritasPusatRepository.FindUnusedByTahun(ctx, tx, tahun)
	if err != nil {
		return nil, err
	}

	var responses []programprioritaspusat.ProgramPrioritasPusatResponse
	for _, result := range results {
		responses = append(responses, programprioritaspusat.ProgramPrioritasPusatResponse{
			Id:                        result.Id,
			NamaTagging:               result.NamaTagging,
			KodeProgramPrioritasPusat:       result.KodeProgramPrioritasPusat,
			KeteranganProgramPrioritasPusat: result.KeteranganProgramPrioritasPusat,
			Keterangan:                result.Keterangan,
			TahunAwal:                 result.TahunAwal,
			TahunAkhir:                result.TahunAkhir,
		})
	}

	return responses, nil
}

func (service *ProgramPrioritasPusatServiceImpl) FindByIdTerkait(ctx context.Context, request programprioritaspusat.FindByIdTerkaitRequest) ([]programprioritaspusat.ProgramPrioritasPusatResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []programprioritaspusat.ProgramPrioritasPusatResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	results, err := service.ProgramPrioritasPusatRepository.FindByIdTerkait(ctx, tx, request.Ids)
	if err != nil {
		return []programprioritaspusat.ProgramPrioritasPusatResponse{}, err
	}

	var responses []programprioritaspusat.ProgramPrioritasPusatResponse
	for _, result := range results {
		responses = append(responses, programprioritaspusat.ProgramPrioritasPusatResponse{
			Id:                        result.Id,
			NamaTagging:               result.NamaTagging,
			KodeProgramPrioritasPusat:       result.KodeProgramPrioritasPusat,
			KeteranganProgramPrioritasPusat: result.KeteranganProgramPrioritasPusat,
			Keterangan:                result.Keterangan,
			TahunAwal:                 result.TahunAwal,
			TahunAkhir:                result.TahunAkhir,
		})
	}

	return responses, nil
}
