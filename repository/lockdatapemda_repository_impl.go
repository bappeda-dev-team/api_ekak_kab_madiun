package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
)

type LockDataPemdaRepositoryImpl struct{}

func NewLockDataPemdaRepositoryImpl() *LockDataPemdaRepositoryImpl {
	return &LockDataPemdaRepositoryImpl{}
}
func (r *LockDataPemdaRepositoryImpl) IsLocked(
	ctx context.Context, tx *sql.Tx, jenis, tahun string,
) (bool, error) {
	var count int
	err := tx.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM tb_lock_data_pemda
		WHERE jenis = ? AND tahun = ?`, jenis, tahun,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("LockDataPemdaRepository.IsLocked: %w", err)
	}
	return count > 0, nil
}
func (r *LockDataPemdaRepositoryImpl) Lock(
	ctx context.Context, tx *sql.Tx, jenis, tahun string,
) error {
	_, err := tx.ExecContext(ctx, `
		INSERT IGNORE INTO tb_lock_data_pemda (jenis, tahun) VALUES (?, ?)`,
		jenis, tahun,
	)
	if err != nil {
		return fmt.Errorf("LockDataPemdaRepository.Lock: %w", err)
	}
	return nil
}
func (r *LockDataPemdaRepositoryImpl) Unlock(
	ctx context.Context, tx *sql.Tx, jenis, tahun string,
) error {
	_, err := tx.ExecContext(ctx, `
		DELETE FROM tb_lock_data_pemda WHERE jenis = ? AND tahun = ?`,
		jenis, tahun,
	)
	if err != nil {
		return fmt.Errorf("LockDataPemdaRepository.Unlock: %w", err)
	}
	return nil
}
func (r *LockDataPemdaRepositoryImpl) FindByJenisTahun(
	ctx context.Context, tx *sql.Tx, jenis, tahun string,
) (domain.LockDataPemda, error) {
	var lock domain.LockDataPemda
	err := tx.QueryRowContext(ctx, `
		SELECT id, jenis, tahun, created_at, updated_at
		FROM tb_lock_data_pemda
		WHERE jenis = ? AND tahun = ?`, jenis, tahun,
	).Scan(&lock.Id, &lock.Jenis, &lock.Tahun, &lock.CreatedAt, &lock.UpdatedAt)
	return lock, err
}
func (r *LockDataPemdaRepositoryImpl) FindAllByJenis(
	ctx context.Context, tx *sql.Tx, jenis string,
) ([]domain.LockDataPemda, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT id, jenis, tahun, created_at, updated_at
		FROM tb_lock_data_pemda
		WHERE jenis = ?
		ORDER BY tahun`, jenis,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.LockDataPemda
	for rows.Next() {
		var lock domain.LockDataPemda
		if err := rows.Scan(&lock.Id, &lock.Jenis, &lock.Tahun, &lock.CreatedAt, &lock.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, lock)
	}
	return result, rows.Err()
}
