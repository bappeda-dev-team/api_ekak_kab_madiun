package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type LockRenaksiOpdRepository interface {
	Lock(ctx context.Context, tx *sql.Tx, lock domain.LockRenaksiOpd) (domain.LockRenaksiOpd, error)
	Unlock(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string, sasaranId int, rekinId string) error
	FindByContext(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string, sasaranId int, rekinId string) (domain.LockRenaksiOpd, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string) ([]domain.LockRenaksiOpd, error)
	IsLocked(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string, sasaranId int, rekinId string) (bool, error)
}
