package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type LockDataRepository interface {
	// IsLocked: cek apakah kodeOpd+tahun+jenisData terkunci
	IsLocked(ctx context.Context, tx *sql.Tx, jenisData, kodeOpd, tahun string) (bool, error)
	// Lock: insert baris lock (INSERT IGNORE agar idempoten)
	Lock(ctx context.Context, tx *sql.Tx, jenisData, kodeOpd, tahun string) error
	// Unlock: hapus baris lock
	Unlock(ctx context.Context, tx *sql.Tx, jenisData, kodeOpd, tahun string) error
	// FindAllByJenisKodeOpd: ambil semua tahun yg di-lock untuk jenis+kodeOpd tertentu
	FindAllByJenisKodeOpd(ctx context.Context, tx *sql.Tx, jenisData, kodeOpd string) ([]domain.LockData, error)
}
