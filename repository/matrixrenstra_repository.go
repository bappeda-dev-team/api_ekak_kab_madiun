package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type MatrixRenstraRepository interface {
	GetByKodeSubKegiatan(ctx context.Context, tx *sql.Tx, kodeOpd, tahunAwal, tahunAkhir string) ([]domain.SubKegiatanQuery, error)
	UpsertIndikator(ctx context.Context, tx *sql.Tx, indikator domain.Indikator) error
	UpsertTarget(ctx context.Context, tx *sql.Tx, target domain.Target) error
	FindIndikatorByKodeIndikator(ctx context.Context, tx *sql.Tx, kodeIndikator string) (domain.Indikator, error)
	CountKodeIndikatorByPrefix(ctx context.Context, tx *sql.Tx, prefix string) (int, error)
	DeleteIndikator(ctx context.Context, tx *sql.Tx, kodeIndikator string) error
	DeleteTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) error
	UpsertAnggaran(ctx context.Context, tx *sql.Tx, kodeSubkegiatan, kodeOpd, tahun string, pagu int64) error
	DeleteIndicatorsExcept(ctx context.Context, tx *sql.Tx, kode, kodeOpd, tahun string, keepList []string) error
}
