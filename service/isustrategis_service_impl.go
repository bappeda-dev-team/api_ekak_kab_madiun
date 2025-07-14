package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web/isustrategis"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"strconv"
)

type CSFServiceImpl struct {
	CSFRepository repository.CSFRepository
	DB            *sql.DB
}

func NewCSFService(csfRepository repository.CSFRepository, db *sql.DB) CSFService {
	return &CSFServiceImpl{
		CSFRepository: csfRepository,
		DB:            db,
	}
}

func (service *CSFServiceImpl) AllCsfsByTahun(ctx context.Context, tahun string) ([]isustrategis.CSFResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Buat instance repository pohon
	pokinRepo := repository.NewPohonKinerjaRepositoryImpl()

	// Ambil data dari CSFRepository
	csfList, err := service.CSFRepository.AllCsfByTahun(ctx, tx, tahun, pokinRepo)
	if err != nil {
		return nil, err
	}

	// Konversi tahun ke int
	tahunInt, err := strconv.Atoi(tahun)
	if err != nil {
		return nil, fmt.Errorf("tahun tidak valid: %v", err)
	}

	var responses []isustrategis.CSFResponse
	for _, pokin := range csfList {
		// Konversi indikator
		var indikatorResponses []pohonkinerja.IndikatorResponse
		for _, ind := range pokin.Indikator {
			var targetResponses []pohonkinerja.TargetResponse
			for _, t := range ind.Target {
				targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
					Id:              t.Id,
					IndikatorId:     t.IndikatorId,
					TargetIndikator: t.Target,
					SatuanIndikator: t.Satuan,
				})
			}

			indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
				Id:            ind.Id,
				NamaIndikator: ind.Indikator,
				Target:        targetResponses,
			})
		}

		// Bangun response CSF
		csfResp := isustrategis.CSFResponse{
			ID:                         pokin.Id, // belum ada ID, default
			PohonID:                    pokin.Id,
			PernyataanKondisiStrategis: "",
			AlasanKondisiStrategis:     "",
			DataTerukur:                "",
			KondisiTerukur:             "",
			KondisiWujud:               "",
			Tahun:                      tahunInt,
			JenisPohon:                 pokin.JenisPohon,
			LevelPohon:                 pokin.LevelPohon,
			Strategi:                   pokin.Strategi,
			NamaPohon:                  pokin.NamaPohon,
			Keterangan:                 pokin.Keterangan,
			IsActive:                   pokin.IsActive,
			Indikators:                 indikatorResponses,
		}

		// Jika CSF ada, isi field-field terkait
		if pokin.CSF != nil {
			csfResp.PernyataanKondisiStrategis = pokin.CSF.PernyataanKondisiStrategis
			csfResp.AlasanKondisiStrategis = pokin.CSF.AlasanKondisiStrategis
			csfResp.DataTerukur = pokin.CSF.DataTerukur
			csfResp.KondisiTerukur = pokin.CSF.KondisiTerukur
			csfResp.KondisiWujud = pokin.CSF.KondisiWujud
			csfResp.Tahun = pokin.CSF.Tahun
		}

		responses = append(responses, csfResp)
	}

	return responses, nil
}

func (service *CSFServiceImpl) FindByTahun(ctx context.Context, tahun string) ([]isustrategis.CSFResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data dari repository
	csfList, err := service.CSFRepository.FindByTahun(ctx, tx, tahun)
	if err != nil {
		return nil, err
	}

	var responses []isustrategis.CSFResponse

	for _, csf := range csfList {
		var indikatorResponses []pohonkinerja.IndikatorResponse

		for _, ind := range csf.Indikator {
			var targetResponses []pohonkinerja.TargetResponse
			for _, t := range ind.Target {
				targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
					Id:              t.Id,
					IndikatorId:     t.IndikatorId,
					TargetIndikator: t.Target,
					SatuanIndikator: t.Satuan,
				})
			}

			indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
				Id:            ind.Id,
				NamaIndikator: ind.Indikator,
				Target:        targetResponses,
			})
		}

		response := isustrategis.CSFResponse{
			ID:                         csf.ID,
			PohonID:                    csf.PohonID,
			PernyataanKondisiStrategis: csf.PernyataanKondisiStrategis,
			AlasanKondisiStrategis:     csf.AlasanKondisiStrategis,
			DataTerukur:                csf.DataTerukur,
			KondisiTerukur:             csf.KondisiTerukur,
			KondisiWujud:               csf.KondisiWujud,
			Tahun:                      csf.Tahun,
			JenisPohon:                 csf.JenisPohon,
			LevelPohon:                 csf.LevelPohon,
			Strategi:                   csf.Strategi,
			Keterangan:                 csf.Keterangan,
			IsActive:                   csf.IsActive,
			Indikators:                 indikatorResponses,
		}

		responses = append(responses, response)
	}

	return responses, nil
}

func (service *CSFServiceImpl) FindById(ctx context.Context, csfId int) (isustrategis.CSFResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return isustrategis.CSFResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data dari repository
	csf, err := service.CSFRepository.FindById(ctx, tx, csfId)
	if err != nil {
		return isustrategis.CSFResponse{}, err
	}

	var indikatorResponses []pohonkinerja.IndikatorResponse
	for _, ind := range csf.Indikator {
		var targetResponses []pohonkinerja.TargetResponse
		for _, t := range ind.Target {
			targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
				Id:              t.Id,
				IndikatorId:     t.IndikatorId,
				TargetIndikator: t.Target,
				SatuanIndikator: t.Satuan,
			})
		}
		indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
			Id:            ind.Id,
			NamaIndikator: ind.Indikator,
			Target:        targetResponses,
		})
	}

	response := isustrategis.CSFResponse{
		ID:                         csf.ID,
		PohonID:                    csf.PohonID,
		PernyataanKondisiStrategis: csf.PernyataanKondisiStrategis,
		AlasanKondisiStrategis:     csf.AlasanKondisiStrategis,
		DataTerukur:                csf.DataTerukur,
		KondisiTerukur:             csf.KondisiTerukur,
		KondisiWujud:               csf.KondisiWujud,
		Tahun:                      csf.Tahun,
		JenisPohon:                 csf.JenisPohon,
		LevelPohon:                 csf.LevelPohon,
		Strategi:                   csf.Strategi,
		NamaPohon:                  csf.Strategi,
		Keterangan:                 csf.Keterangan,
		IsActive:                   csf.IsActive,
		Indikators:                 indikatorResponses,
	}

	return response, nil
}
