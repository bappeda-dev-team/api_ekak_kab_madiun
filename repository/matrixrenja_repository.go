package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type MatrixRenjaRepository interface {
	GetRenja(ctx context.Context, tx *sql.Tx, kodeOpd, tahun, jenisIndikator, jenisPagu string) ([]domain.SubKegiatanQuery, error)
	GetRenjaRankhir(ctx context.Context, tx *sql.Tx, kodeOpd, tahun, jenisIndikator string) ([]domain.SubKegiatanQuery, error)
	UpsertBatchIndikatorRenja(ctx context.Context, tx *sql.Tx, items []domain.Indikator) error
	CountKodeIndikatorByPrefix(ctx context.Context, tx *sql.Tx, prefix string) (int, error)
	FindIndikatorRenjaByKode(ctx context.Context, tx *sql.Tx, kodeIndikator string) (domain.Indikator, error)
	UpsertAnggaran(ctx context.Context, tx *sql.Tx, kodeSubkegiatan, kodeOpd, tahun string, pagu int64) error
	DeleteIndicatorsExcept(ctx context.Context, tx *sql.Tx, kode, kodeOpd, tahun, jenis string, keepList []string) error
	GetRenjaPenetapan(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string, jenisPagu string) ([]domain.SubKegiatanQuery, error)
}
