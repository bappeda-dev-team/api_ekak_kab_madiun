package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type LockDataPemdaRepository interface {
	IsLocked(ctx context.Context, tx *sql.Tx, jenis, tahun string) (bool, error)
	Lock(ctx context.Context, tx *sql.Tx, jenis, tahun string) error
	Unlock(ctx context.Context, tx *sql.Tx, jenis, tahun string) error
	FindByJenisTahun(ctx context.Context, tx *sql.Tx, jenis, tahun string) (domain.LockDataPemda, error)
	FindAllByJenis(ctx context.Context, tx *sql.Tx, jenis string) ([]domain.LockDataPemda, error)
}
