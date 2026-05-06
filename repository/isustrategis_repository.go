package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/isustrategis"
	"ekak_kabupaten_madiun/model/web/strategic"
	"ekak_kabupaten_madiun/model/web/strategicarahkebijakan"
)

type CSFRepository interface {
	FindByTahun(ctx context.Context, tx *sql.Tx, tahun string) ([]isustrategis.CSFPokin, error)
	IsuFindByTahun(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]strategic.IsuStrategiOpd, error)
	IsuFindBetweenTahun(ctx context.Context, tx *sql.Tx, tahunAwal string, tahunAkhir string) ([]strategicarahkebijakan.IsuStrategiPemda, error)
	CreateCsf(ctx context.Context, tx *sql.Tx, csf domain.CSF) error
	UpdateCSFByPohonID(ctx context.Context, tx *sql.Tx, csf domain.CSF) (domain.CSF, error)
	FindById(ctx context.Context, tx *sql.Tx, csfId int) (isustrategis.CSFPokin, error)
}
