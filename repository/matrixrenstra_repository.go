package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type MatrixRenstraRepository interface {
	GetByKodeSubKegiatan(ctx context.Context, tx *sql.Tx, kodeOpd string, tahunAwal string, tahunAkhir string) ([]domain.SubKegiatanQuery, error)
	SaveIndikator(ctx context.Context, tx *sql.Tx, indikator domain.Indikator) error
	SaveTarget(ctx context.Context, tx *sql.Tx, target domain.Target) error
	FindIndikatorById(ctx context.Context, tx *sql.Tx, indikatorId string) (domain.Indikator, error)
	UpdateIndikator(ctx context.Context, tx *sql.Tx, indikator domain.Indikator) error
	UpdateTarget(ctx context.Context, tx *sql.Tx, target domain.Target) error
	DeleteIndikator(ctx context.Context, tx *sql.Tx, indikatorId string) error
	DeleteTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) error
}
