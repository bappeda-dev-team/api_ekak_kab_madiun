package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type MatrixRenjaRepository interface {
	GetRenja(ctx context.Context, tx *sql.Tx, kodeOpd, tahun, jenisPagu string) ([]domain.SubKegiatanQuery, error)
	GetRenjaRankhir(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string) ([]domain.SubKegiatanQuery, error)
	UpsertBatchIndikatorRenja(ctx context.Context, tx *sql.Tx, items []domain.Indikator) error
	CountIndikatorMatrixByScope(ctx context.Context, tx *sql.Tx, kode, kodeOpd, tahun, jenis, prefix string) (int, error)
	FindIndikatorRenjaByKode(ctx context.Context, tx *sql.Tx, kodeIndikator string) (domain.Indikator, error)
	UpsertAnggaran(ctx context.Context, tx *sql.Tx, kodeSubkegiatan, kodeOpd, tahun string, pagu int64) error
	DeleteIndicatorsExcept(ctx context.Context, tx *sql.Tx, kode, kodeOpd, tahun, jenis string, keepList []string) error
	GetRenjaPenetapan(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string, jenisPagu string) ([]domain.SubKegiatanQuery, error)

	// CRUD target per layer (rankhir / penetapan)
	UpsertTargetRenjaLayer(ctx context.Context, tx *sql.Tx, kodeIndikator, tahun, target, satuan, jenis string) error
	DeleteTargetRenjaLayer(ctx context.Context, tx *sql.Tx, kodeIndikator, tahun, jenis string) error
	FindTargetRenjaLayer(ctx context.Context, tx *sql.Tx, kodeIndikator, tahun, jenis string) (domain.Target, error)

	// CRUD pagu per jenis (ranwal / rankhir / penetapan) — generik
	UpsertPaguRenjaByJenis(ctx context.Context, tx *sql.Tx, kodeSubkegiatan, kodeOpd, tahun, jenis string, pagu int64) error

	// V2 — sama dengan GetRenja / GetRenjaRankhir / GetRenjaPenetapan tetapi khusus level Program:
	// indikator dari tb_master_indikator_program_kegiatan, target dari tb_target_program_kegiatan.
	GetRenjaV2(ctx context.Context, tx *sql.Tx, kodeOpd, tahun, jenisPagu string) ([]domain.SubKegiatanQuery, error)
	GetRenjaRankhirV2(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string) ([]domain.SubKegiatanQuery, error)
	GetRenjaPenetapanV2(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string, jenisPagu string) ([]domain.SubKegiatanQuery, error)
}
