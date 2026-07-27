package repository

import (
	"context"
	"database/sql"
)

type LockDataRepository interface {
	// IsLocked: cek apakah kodeOpd+tahun+jenisData terkunci
	IsLocked(ctx context.Context, tx *sql.Tx, jenisData, kodeOpd, tahun string) (bool, error)
	// Lock: insert baris lock (INSERT IGNORE agar idempoten)
	Lock(ctx context.Context, tx *sql.Tx, jenisData, kodeOpd, tahun string) error
	// Unlock: hapus baris lock
	Unlock(ctx context.Context, tx *sql.Tx, jenisData, kodeOpd, tahun string) error
}
