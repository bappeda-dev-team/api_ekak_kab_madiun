package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type SasaranPemdaRepository interface {
	// ── CRUD ─────────────────────────────────────────────────────
	Create(ctx context.Context, tx *sql.Tx, sasaranPemda domain.SasaranPemda) (domain.SasaranPemda, error)
	Update(ctx context.Context, tx *sql.Tx, sasaranPemda domain.SasaranPemda) (domain.SasaranPemda, error)
	Delete(ctx context.Context, tx *sql.Tx, sasaranPemdaId int) error
	DeleteIndikator(ctx context.Context, tx *sql.Tx, sasaranPemdaId int) error
	// ── FIND ─────────────────────────────────────────────────────
	FindById(ctx context.Context, tx *sql.Tx, sasaranPemdaId int) (domain.SasaranPemda, error)
	FindAll(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.SasaranPemda, error)
	FindAllWithPokin(ctx context.Context, tx *sql.Tx, tahunAwal, tahunAkhir, jenisPeriode string) ([]domain.PohonKinerjaWithSasaran, error)
	FindAllByTahun(ctx context.Context, tx *sql.Tx, tahun, jenisPeriode, jenis string) ([]domain.SasaranPemda, error)
	FindIndikatorByKode(ctx context.Context, tx *sql.Tx, kodeIndikator string) (domain.IndikatorPemda, error)
	FindTargetLayerById(ctx context.Context, tx *sql.Tx, id int) (domain.TargetPemda, error)
	FindRanwalByTahun(ctx context.Context, tx *sql.Tx, tahun, jenisPeriode string) ([]domain.SasaranPemda, error)
	// ── TARGET LAYER ─────────────────────────────────────────────
	CreateTargetLayer(ctx context.Context, tx *sql.Tx, target domain.TargetPemda) (domain.TargetPemda, error)
	UpdateTargetLayerById(ctx context.Context, tx *sql.Tx, id int, target, satuan string) (domain.TargetPemda, error)
	UpsertTargetPemda(ctx context.Context, tx *sql.Tx, t domain.TargetPemda) (domain.TargetPemda, error)
	// ── UTILS ────────────────────────────────────────────────────
	IsIdExists(ctx context.Context, tx *sql.Tx, id int) bool
	IsSubtemaIdExists(ctx context.Context, tx *sql.Tx, subtemaId int) bool
	UpdatePeriode(ctx context.Context, tx *sql.Tx, sasaranPemda domain.SasaranPemda) (domain.SasaranPemda, error)
}
