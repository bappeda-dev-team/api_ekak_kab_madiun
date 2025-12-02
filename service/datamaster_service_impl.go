package service

import (
	"context"
	"errors"
	"fmt"

	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web/datamaster"
	"ekak_kabupaten_madiun/repository"
)

type DataMasterServiceImpl struct {
	DataMasterRepository repository.DataMasterRepository
	DB                   *sql.DB
}

func NewDataMasterServiceImpl(dataMasterRepository repository.DataMasterRepository, DB *sql.DB) *DataMasterServiceImpl {
	return &DataMasterServiceImpl{
		DataMasterRepository: dataMasterRepository,
		DB:                   DB,
	}
}

func (service *DataMasterServiceImpl) DataRBByTahun(ctx context.Context, tahunBase int) ([]datamaster.RBResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data dari repository
	result, err := service.DataMasterRepository.DataRBByTahun(ctx, tx, tahunBase)
	if err != nil {
		return nil, err
	}

	// Convert MasterRB → RBResponse
	responses := make([]datamaster.RBResponse, 0, len(result))

	for _, rb := range result {
		// Mapping MasterRB → RBResponse
		resp := datamaster.RBResponse{
			IdRB:          rb.Id,
			JenisRB:       rb.JenisRB,
			KegiatanUtama: rb.KegiatanUtama,
			Keterangan:    rb.Keterangan,
			TahunBaseline: rb.TahunBaseline,
			TahunNext:     rb.TahunNext,
			Indikator:     make([]datamaster.IndikatorRB, 0),
		}

		// Mapping Indikator
		for _, ind := range rb.Indikator {
			indResp := datamaster.IndikatorRB{
				IdIndikator: ind.IdIndikator,
				IdRB:        ind.IdRB,
				Indikator:   ind.Indikator,
				TargetRB:    make([]datamaster.TargetRB, 0),
			}

			for _, tar := range ind.TargetRB {
				rbResp := datamaster.TargetRB{
					IdTarget:          tar.IdTarget,
					IdIndikator:       tar.IdIndikator,
					TahunBaseline:     tar.TahunBaseline,
					TargetBaseline:    tar.TargetBaseline,
					RealisasiBaseline: tar.RealisasiBaseline,
					SatuanBaseline:    tar.SatuanBaseline,
					TahunNext:         tar.TahunNext,
					TargetNext:        tar.TargetNext,
					SatuanNext:        tar.SatuanNext,
				}
				indResp.TargetRB = append(indResp.TargetRB, rbResp)
			}

			resp.Indikator = append(resp.Indikator, indResp)
		}

		responses = append(responses, resp)
	}

	return responses, nil
}

func (service *DataMasterServiceImpl) SaveRB(ctx context.Context, rb datamaster.RBRequest, userId int) (datamaster.RBResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return datamaster.RBResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// === Convert RBRequest → MasterRB ===
	entity := datamaster.ConvertRBRequestToMaster(rb, userId)

	// === 1. INSERT MASTER RB ===
	rbID, err := service.DataMasterRepository.InsertRB(ctx, tx, entity, userId)
	if err != nil {
		return datamaster.RBResponse{}, err
	}

	// === Build Response RB ===
	response := datamaster.RBResponse{
		IdRB:          int(rbID),
		JenisRB:       entity.JenisRB,
		KegiatanUtama: entity.KegiatanUtama,
		Keterangan:    entity.Keterangan,
		TahunBaseline: entity.TahunBaseline,
		TahunNext:     entity.TahunNext,
		Indikator:     []datamaster.IndikatorRB{},
	}

	// === 2. INSERT INDIKATOR ===
	for _, ind := range entity.Indikator {

		// INSERT into DB
		indikatorID, err := service.DataMasterRepository.InsertIndikator(ctx, tx, rbID, ind)
		if err != nil {
			return datamaster.RBResponse{}, err
		}

		// siapkan indikator response
		indResp := datamaster.IndikatorRB{
			IdRB:        int(rbID),
			IdIndikator: indikatorID,
			Indikator:   ind.Indikator,
			TargetRB:    []datamaster.TargetRB{},
		}

		// === 3. INSERT TARGET ===
		for _, t := range ind.TargetRB {

			// simpan target
			if err := service.DataMasterRepository.InsertTarget(ctx, tx, indikatorID, t); err != nil {
				return datamaster.RBResponse{}, err
			}

			// tambahkan ke response
			indResp.TargetRB = append(indResp.TargetRB, datamaster.TargetRB{
				IdIndikator:       indikatorID,
				TahunBaseline:     t.TahunBaseline,
				TargetBaseline:    t.TargetBaseline,
				RealisasiBaseline: t.RealisasiBaseline,
				SatuanBaseline:    t.SatuanBaseline,
				TahunNext:         t.TahunNext,
				TargetNext:        t.TargetNext,
				SatuanNext:        t.SatuanNext,
			})
		}

		// masukkan indikator ke response
		response.Indikator = append(response.Indikator, indResp)
	}

	return response, nil
}

func (service *DataMasterServiceImpl) UpdateRB(ctx context.Context, rb datamaster.RBRequest, userId int, rbId int) (datamaster.RBResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return datamaster.RBResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// === 1. Cek apakah RB ada ===
	existingRB, err := service.DataMasterRepository.FindRBById(ctx, tx, rbId)
	if err != nil {
		return datamaster.RBResponse{}, errors.New("rb_not_found")
	}

	// === 2. Convert RBRequest → Entity ===
	entity := datamaster.ConvertRBRequestToMaster(rb, userId)
	entity.Id = existingRB.Id // penting

	// === 3. Update master RB ===
	if err := service.DataMasterRepository.UpdateRB(ctx, tx, entity, rbId); err != nil {
		return datamaster.RBResponse{}, err
	}

	// === 4. Hapus semua indikator & target lama ===
	if err := service.DataMasterRepository.DeleteAllIndikatorAndTargetByRB(ctx, tx, rbId); err != nil {
		return datamaster.RBResponse{}, err
	}

	// === 5. Insert indikator & target baru ===
	response := datamaster.RBResponse{
		IdRB:          rbId,
		JenisRB:       entity.JenisRB,
		KegiatanUtama: entity.KegiatanUtama,
		Keterangan:    entity.Keterangan,
		TahunBaseline: entity.TahunBaseline,
		TahunNext:     entity.TahunNext,
		Indikator:     []datamaster.IndikatorRB{},
	}

	for _, ind := range entity.Indikator {

		// Insert indikator baru
		indikatorID, err := service.DataMasterRepository.InsertIndikator(ctx, tx, int64(rbId), ind)
		if err != nil {
			return datamaster.RBResponse{}, err
		}

		indikatorRes := datamaster.IndikatorRB{
			IdRB:        rbId,
			IdIndikator: indikatorID,
			Indikator:   ind.Indikator,
			TargetRB:    []datamaster.TargetRB{},
		}

		for _, t := range ind.TargetRB {
			if err := service.DataMasterRepository.InsertTarget(ctx, tx, indikatorID, t); err != nil {
				return datamaster.RBResponse{}, err
			}

			indikatorRes.TargetRB = append(indikatorRes.TargetRB, datamaster.TargetRB{
				IdIndikator:       indikatorID,
				TahunBaseline:     t.TahunBaseline,
				TargetBaseline:    t.TargetBaseline,
				RealisasiBaseline: t.RealisasiBaseline,
				SatuanBaseline:    t.SatuanBaseline,
				TahunNext:         t.TahunNext,
				TargetNext:        t.TargetNext,
				SatuanNext:        t.SatuanNext,
			})
		}

		response.Indikator = append(response.Indikator, indikatorRes)
	}

	return response, nil
}

func (s *DataMasterServiceImpl) DeleteRB(ctx context.Context, rbId int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	err = s.DataMasterRepository.DeleteRB(ctx, tx, rbId)
	if err != nil {
		return err
	}
	return nil
}
