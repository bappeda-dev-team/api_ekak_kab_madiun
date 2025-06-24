package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web/isustrategis"
	"ekak_kabupaten_madiun/repository"
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

	// Transformasi ke bentuk response (web model)
	var responses []isustrategis.CSFResponse
	for _, csf := range csfList {
		responses = append(responses, isustrategis.CSFResponse{
			ID:                         csf.ID,
			PohonID:                    csf.PohonID,
			PernyataanKondisiStrategis: csf.PernyataanKondisiStrategis,
			AlasanKondisiStrategis:     csf.AlasanKondisiStrategis,
			DataTerukur:                csf.DataTerukur,
			KondisiTerukur:             csf.KondisiTerukur,
			KondisiWujud:               csf.KondisiWujud,
			Tahun:                      csf.Tahun,
			CreatedAt:                  csf.CreatedAt,
			UpdatedAt:                  csf.UpdatedAt,
		})
	}

	return responses, nil
}
