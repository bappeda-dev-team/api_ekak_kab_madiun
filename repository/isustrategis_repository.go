package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/isustrategis"
)

type CSFRepository interface {
	AllCsfByTahun(ctx context.Context, tx *sql.Tx, tahun string, repository PohonKinerjaRepository) ([]domain.PohonKinerja, error)
	FindByTahun(ctx context.Context, tx *sql.Tx, tahun string) ([]isustrategis.CSFPokin, error)
	CreateCsf(ctx context.Context, tx *sql.Tx, csf domain.CSF) error
	UpdateCSFByPohonID(ctx context.Context, tx *sql.Tx, csf domain.CSF) (domain.CSF, error)
	FindById(ctx context.Context, tx *sql.Tx, csfId int) (isustrategis.CSFPokin, error)
}
