package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
)

type LockDataRepositoryImpl struct{}

func NewLockDataRepositoryImpl() *LockDataRepositoryImpl {
	return &LockDataRepositoryImpl{}
}
func (r *LockDataRepositoryImpl) IsLocked(
	ctx context.Context, tx *sql.Tx, jenisData, kodeOpd, tahun string,
) (bool, error) {
	var count int
	err := tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM tb_lock_data WHERE jenis_data=? AND kode_opd=? AND tahun=?`,
		jenisData, kodeOpd, tahun,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("LockDataRepository.IsLocked: %w", err)
	}
	return count > 0, nil
}
func (r *LockDataRepositoryImpl) Lock(
	ctx context.Context, tx *sql.Tx, jenisData, kodeOpd, tahun string,
) error {
	// INSERT IGNORE → idempoten, tidak error jika sudah ada
	_, err := tx.ExecContext(ctx,
		`INSERT IGNORE INTO tb_lock_data (jenis_data, kode_opd, tahun) VALUES (?, ?, ?)`,
		jenisData, kodeOpd, tahun,
	)
	if err != nil {
		return fmt.Errorf("LockDataRepository.Lock: %w", err)
	}
	return nil
}
func (r *LockDataRepositoryImpl) Unlock(
	ctx context.Context, tx *sql.Tx, jenisData, kodeOpd, tahun string,
) error {
	_, err := tx.ExecContext(ctx,
		`DELETE FROM tb_lock_data WHERE jenis_data=? AND kode_opd=? AND tahun=?`,
		jenisData, kodeOpd, tahun,
	)
	if err != nil {
		return fmt.Errorf("LockDataRepository.Unlock: %w", err)
	}
	return nil
}

func (r *LockDataRepositoryImpl) FindAllByJenisKodeOpd(
	ctx context.Context, tx *sql.Tx, jenisData, kodeOpd string,
) ([]domain.LockData, error) {
	rows, err := tx.QueryContext(ctx,
		`SELECT id, jenis_data, kode_opd, tahun FROM tb_lock_data WHERE jenis_data=? AND kode_opd=? ORDER BY tahun`,
		jenisData, kodeOpd,
	)
	if err != nil {
		return nil, fmt.Errorf("LockDataRepository.FindAllByJenisKodeOpd: %w", err)
	}
	defer rows.Close()
	var result []domain.LockData
	for rows.Next() {
		var d domain.LockData
		if err := rows.Scan(&d.Id, &d.JenisData, &d.KodeOpd, &d.Tahun); err != nil {
			return nil, fmt.Errorf("LockDataRepository.FindAllByJenisKodeOpd scan: %w", err)
		}
		result = append(result, d)
	}
	return result, rows.Err()
}
