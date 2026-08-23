package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
)

type TujuanOpdRepository interface {
	Create(ctx context.Context, tx *sql.Tx, tujuanOpd domain.TujuanOpd) (domain.TujuanOpd, error)
	Update(ctx context.Context, tx *sql.Tx, tujuanOpd domain.TujuanOpd) error
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.TujuanOpd, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]domain.TujuanOpd, error)
	FindIndikatorByTujuanId(ctx context.Context, tx *sql.Tx, tujuanOpdId int) ([]domain.Indikator, error)
	FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string, tahun string) ([]domain.Target, error)
	FindTujuanOpdByTahun(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string, jenisPeriode string) ([]domain.TujuanOpd, error)
	FindIndikatorByTujuanOpdId(ctx context.Context, tx *sql.Tx, tujuanOpdId int) ([]domain.Indikator, error)
	FindTujuanOpdForCascadingOpd(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string, jenisPeriode string) ([]domain.TujuanOpd, error)
	FindIndikatorByTujuanOpdIdsBatch(ctx context.Context, tx *sql.Tx, tujuanOpdIds []int) (map[int][]domain.Indikator, error)

	//trnstra
	// FindAll untuk renstra (range tahun)
	FindAllByPeriod(ctx context.Context, tx *sql.Tx, kodeOpd, tahunAwal, tahunAkhir, jenisPeriode, jenisIndikator string) ([]domain.TujuanOpd, error)
	// FindAll untuk ranwal/rankhir (single tahun)
	FindAllByTahun(ctx context.Context, tx *sql.Tx, kodeOpd, tahun, jenisPeriode, jenisIndikator string) ([]domain.TujuanOpd, error)
	// Batch fetch untuk optimasi
	FindBidangUrusanBatch(ctx context.Context, tx *sql.Tx, kodeBidangUrusanList []string) (map[string]domainmaster.BidangUrusan, error)
	CreateRenjaIndikator(ctx context.Context, tx *sql.Tx, tujuanOpdId int, indikators []domain.Indikator) error
	UpdateRenjaIndikator(ctx context.Context, tx *sql.Tx, indikators []domain.Indikator) error
	DeleteIndikatorTargetRenja(ctx context.Context, tx *sql.Tx, indikatorId string) error
	FindIndikatorByKodeIndikator(ctx context.Context, tx *sql.Tx, kodeIndikator string) (domain.Indikator, error)
	FindAllByTahunForPokin(ctx context.Context, tx *sql.Tx, kodeOpd, tahun, jenisPeriode, jenisIndikator string) ([]domain.TujuanOpd, error)
	FindAllByTahunDualJenis(ctx context.Context, tx *sql.Tx, kodeOpd, tahun, jenisPeriode, jenisOverride string) ([]domain.TujuanOpd, error)

	SetTujuanOpdLocked(ctx context.Context, tx *sql.Tx, id int, locked bool) error

	// FindTujuanOpdWithIndikatorByIdsBatch mengambil data tujuan OPD lengkap (dengan
	// indikator dan target renstra) secara batch berdasarkan list id tujuan_opd.
	FindTujuanOpdWithIndikatorByIdsBatch(ctx context.Context, tx *sql.Tx, ids []int) (map[int]domain.TujuanOpd, error)

	// ── Target-only CRUD (ranwal/rankhir/penetapan) ──────────────
	TargetOpdExistsByKey(ctx context.Context, tx *sql.Tx, kodeIndikator, tahun, jenis string) (bool, error)
	CreateTargetOpdSingle(ctx context.Context, tx *sql.Tx, t domain.Target) (domain.Target, error)
	FindTargetOpdById(ctx context.Context, tx *sql.Tx, id string) (domain.Target, error)
	UpdateTargetOpdById(ctx context.Context, tx *sql.Tx, id, targetValue, satuan string) (domain.Target, error)
	DeleteTargetOpdByJenis(ctx context.Context, tx *sql.Tx, kodeIndikator, tahun, jenis string) error
}
