package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
)

type CrosscuttingOpdRepository interface {
	CreateCrosscutting(ctx context.Context, tx *sql.Tx, pokin domain.PohonKinerja, parentId int) (domain.PohonKinerja, error)
	UpdateCrosscutting(ctx context.Context, tx *sql.Tx, pokin domain.PohonKinerja) (domain.PohonKinerja, error)
	DeleteCrosscutting(ctx context.Context, tx *sql.Tx, pokinId int, nipPegawai string) error
	FindAllCrosscutting(ctx context.Context, tx *sql.Tx, crosscuttingId int) ([]domain.PohonKinerja, error)
	ValidateKodeOpdChange(ctx context.Context, tx *sql.Tx, id int) error
	FindTargetByIndikatorIds(ctx context.Context, tx *sql.Tx, indikatorIds []string) ([]domain.Target, error)
	FindIndikatorByPokinId(ctx context.Context, tx *sql.Tx, pokinIds []int) ([]domain.Indikator, error)
	ApproveOrRejectCrosscutting(ctx context.Context, tx *sql.Tx, crosscuttingId int, request pohonkinerja.CrosscuttingApproveRequest) error
	DeleteUnused(ctx context.Context, tx *sql.Tx, crosscuttingId int) error
	FindPokinByCrosscuttingStatus(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.Crosscutting, error)
	FindOPDCrosscuttingFrom(ctx context.Context, tx *sql.Tx, crosscuttingTo int) (string, error)
	// DeleteCrosscuttingExisting(ctx context.Context, tx *sql.Tx, crosscuttingId int) error
	FindCrosscuttingByPokinIdsBatch(ctx context.Context, tx *sql.Tx, pokinIds []int) (map[int][]domain.Crosscutting, error)

	//crosscutting legacy untuk delete
	FixPokinStatusAfterExistingUnlink(ctx context.Context, tx *sql.Tx, pokinId int) error
	FixPokinStatusAfterExistingDelete(ctx context.Context, tx *sql.Tx, pokinId int) error

	// Plan A: jika ref tunggal → hapus pohon kinerja + child
	DeleteCrosscuttingDiterima(ctx context.Context, tx *sql.Tx, crosscuttingId int) error
	// Plan B: jika ref tunggal → hanya lepas tautan, pohon kinerja tidak dihapus
	UnlinkCrosscuttingDiterima(ctx context.Context, tx *sql.Tx, crosscuttingId int) error
}
