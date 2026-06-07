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
		var opdResponseList []programunggulan.OpdProgramUnggulanResponse
		for _, opd := range result.OpdList {
			opdResponseList = append(opdResponseList, programunggulan.OpdProgramUnggulanResponse{
				Id:      opd.Id,
				KodeOpd: opd.KodeOpd,
				NamaOpd: opd.NamaOpd,
			})
		}
		responses = append(responses, programunggulan.ProgramUnggulanResponse{
			Id:                        result.Id,
			NamaTagging:               result.NamaTagging,
			KodeProgramUnggulan:       result.KodeProgramUnggulan,
			KeteranganProgramUnggulan: result.KeteranganProgramUnggulan,
			Keterangan:                result.Keterangan,
			TahunAwal:                 result.TahunAwal,
			TahunAkhir:                result.TahunAkhir,
			IsActive:                  result.IsActive,
			OpdList:                   opdResponseList,
			TahunTerpakai:             result.TahunTerpakai,
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

func (service *ProgramUnggulanServiceImpl) FindByTahun(ctx context.Context, tahun string, kodeOpd string) ([]programunggulan.ProgramUnggulanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	_, err = strconv.Atoi(tahun)
	if err != nil {
		return nil, errors.New("format tahun tidak valid")
	}
	if kodeOpd == "" {
		return nil, errors.New("kode_opd tidak boleh kosong")
	}
	results, err := service.ProgramUnggulanRepository.FindByTahunAndKodeOpd(ctx, tx, tahun, kodeOpd)
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

func (service *ProgramUnggulanServiceImpl) CreateOpdProgramUnggulan(ctx context.Context, request programunggulan.CreateOpdProgramUnggulanRequest) (programunggulan.CreateOpdProgramUnggulanResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return programunggulan.CreateOpdProgramUnggulanResponse{}, err
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return programunggulan.CreateOpdProgramUnggulanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	_, err = service.ProgramUnggulanRepository.FindByKodeProgramUnggulan(ctx, tx, request.KodeProgramUnggulan)
	if err != nil {
		return programunggulan.CreateOpdProgramUnggulanResponse{}, errors.New("program unggulan tidak ditemukan")
	}
	kodeOpds := dedupeStrings(request.KodeOpd)
	// Guard: cek data yang sudah ada
	existingOpds, err := service.ProgramUnggulanRepository.FindOpdByKodeProgramUnggulanAndKodeOpds(ctx, tx, request.KodeProgramUnggulan, kodeOpds)
	if err != nil {
		return programunggulan.CreateOpdProgramUnggulanResponse{}, err
	}
	// Filter hanya kode_opd yang belum ada
	newKodeOpds := filterNewKodeOpds(kodeOpds, existingOpds)
	// Insert hanya yang baru
	if len(newKodeOpds) > 0 {
		err = service.ProgramUnggulanRepository.CreateOpdProgramUnggulan(ctx, tx, request.KodeProgramUnggulan, newKodeOpds)
		if err != nil {
			return programunggulan.CreateOpdProgramUnggulanResponse{}, err
		}
	}
	// Ambil semua data (existing + baru) untuk response
	opdList, err := service.ProgramUnggulanRepository.FindOpdByKodeProgramUnggulanAndKodeOpds(ctx, tx, request.KodeProgramUnggulan, kodeOpds)
	if err != nil {
		return programunggulan.CreateOpdProgramUnggulanResponse{}, err
	}
	var opdResponses []programunggulan.OpdProgramUnggulanResponse
	for _, opd := range opdList {
		opdResponses = append(opdResponses, programunggulan.OpdProgramUnggulanResponse{
			Id:      opd.Id,
			KodeOpd: opd.KodeOpd,
			NamaOpd: opd.NamaOpd,
		})
	}
	return programunggulan.CreateOpdProgramUnggulanResponse{
		KodeProgramUnggulan: request.KodeProgramUnggulan,
		OpdList:             opdResponses,
	}, nil
}
func filterNewKodeOpds(requested []string, existing []domain.OpdProgramUnggulan) []string {
	existingMap := make(map[string]struct{}, len(existing))
	for _, opd := range existing {
		existingMap[opd.KodeOpd] = struct{}{}
	}
	var result []string
	for _, kodeOpd := range requested {
		if _, exists := existingMap[kodeOpd]; exists {
			continue // skip, sudah ada di database
		}
		result = append(result, kodeOpd)
	}
	return result
}

func dedupeStrings(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		if _, exists := seen[item]; exists {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}

func (service *ProgramUnggulanServiceImpl) DeleteOpdProgramUnggulan(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)
	opdList, err := service.ProgramUnggulanRepository.FindOpdProgramUnggulanById(ctx, tx, id)
	if err != nil {
		return err
	}
	if len(opdList) == 0 {
		return errors.New("opd program unggulan tidak ditemukan")
	}
	return service.ProgramUnggulanRepository.DeleteOpdProgramUnggulan(ctx, tx, id)
}
