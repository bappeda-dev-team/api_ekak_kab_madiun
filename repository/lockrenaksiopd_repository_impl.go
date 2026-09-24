package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
)

type LockRenaksiOpdRepositoryImpl struct{}

func NewLockRenaksiOpdRepositoryImpl() *LockRenaksiOpdRepositoryImpl {
	return &LockRenaksiOpdRepositoryImpl{}
}

func (repository *LockRenaksiOpdRepositoryImpl) Lock(ctx context.Context, tx *sql.Tx, lock domain.LockRenaksiOpd) (domain.LockRenaksiOpd, error) {
	query := `
		INSERT INTO tb_lock_renaksi_opd (
			kode_opd, tahun, sasaran_id, rekin_id, aksi_kegiatan, sub_kegiatan,
			anggaran, nama_pemilik, tw1, tw2, tw3, tw4
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			aksi_kegiatan = VALUES(aksi_kegiatan),
			sub_kegiatan = VALUES(sub_kegiatan),
			anggaran = VALUES(anggaran),
			nama_pemilik = VALUES(nama_pemilik),
			tw1 = VALUES(tw1),
			tw2 = VALUES(tw2),
			tw3 = VALUES(tw3),
			tw4 = VALUES(tw4),
			id = LAST_INSERT_ID(id)
	`
	_, err := tx.ExecContext(ctx, query,
		lock.KodeOpd,
		lock.Tahun,
		lock.SasaranId,
		lock.RekinId,
		lock.AksiKegiatan,
		lock.SubKegiatan,
		lock.Anggaran,
		lock.NamaPemilik,
		lock.Tw1,
		lock.Tw2,
		lock.Tw3,
		lock.Tw4,
	)
	if err != nil {
		return domain.LockRenaksiOpd{}, fmt.Errorf("LockRenaksiOpdRepository.Lock: %w", err)
	}
	return repository.FindByContext(ctx, tx, lock.KodeOpd, lock.Tahun, lock.SasaranId, lock.RekinId)
}

func (repository *LockRenaksiOpdRepositoryImpl) Unlock(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string, sasaranId int, rekinId string) error {
	result, err := tx.ExecContext(ctx, `
		DELETE FROM tb_lock_renaksi_opd
		WHERE kode_opd = ? AND tahun = ? AND sasaran_id = ? AND rekin_id = ?
	`, kodeOpd, tahun, sasaranId, rekinId)
	if err != nil {
		return fmt.Errorf("LockRenaksiOpdRepository.Unlock: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("LockRenaksiOpdRepository.Unlock RowsAffected: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (repository *LockRenaksiOpdRepositoryImpl) FindByContext(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string, sasaranId int, rekinId string) (domain.LockRenaksiOpd, error) {
	var lock domain.LockRenaksiOpd
	err := tx.QueryRowContext(ctx, `
		SELECT id, kode_opd, tahun, sasaran_id, rekin_id, aksi_kegiatan, sub_kegiatan,
		       anggaran, nama_pemilik, tw1, tw2, tw3, tw4
		FROM tb_lock_renaksi_opd
		WHERE kode_opd = ? AND tahun = ? AND sasaran_id = ? AND rekin_id = ?
	`, kodeOpd, tahun, sasaranId, rekinId).Scan(
		&lock.Id,
		&lock.KodeOpd,
		&lock.Tahun,
		&lock.SasaranId,
		&lock.RekinId,
		&lock.AksiKegiatan,
		&lock.SubKegiatan,
		&lock.Anggaran,
		&lock.NamaPemilik,
		&lock.Tw1,
		&lock.Tw2,
		&lock.Tw3,
		&lock.Tw4,
	)
	if err != nil {
		return domain.LockRenaksiOpd{}, fmt.Errorf("LockRenaksiOpdRepository.FindByContext: %w", err)
	}
	return lock, nil
}

func (repository *LockRenaksiOpdRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string) ([]domain.LockRenaksiOpd, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT id, kode_opd, tahun, sasaran_id, rekin_id, aksi_kegiatan, sub_kegiatan,
		       anggaran, nama_pemilik, tw1, tw2, tw3, tw4
		FROM tb_lock_renaksi_opd
		WHERE kode_opd = ? AND tahun = ?
		ORDER BY id ASC
	`, kodeOpd, tahun)
	if err != nil {
		return nil, fmt.Errorf("LockRenaksiOpdRepository.FindAll: %w", err)
	}
	defer rows.Close()

	result := make([]domain.LockRenaksiOpd, 0)
	for rows.Next() {
		var lock domain.LockRenaksiOpd
		if err := rows.Scan(
			&lock.Id,
			&lock.KodeOpd,
			&lock.Tahun,
			&lock.SasaranId,
			&lock.RekinId,
			&lock.AksiKegiatan,
			&lock.SubKegiatan,
			&lock.Anggaran,
			&lock.NamaPemilik,
			&lock.Tw1,
			&lock.Tw2,
			&lock.Tw3,
			&lock.Tw4,
		); err != nil {
			return nil, fmt.Errorf("LockRenaksiOpdRepository.FindAll scan: %w", err)
		}
		result = append(result, lock)
	}
	return result, rows.Err()
}

func (repository *LockRenaksiOpdRepositoryImpl) IsLocked(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string, sasaranId int, rekinId string) (bool, error) {
	var count int
	err := tx.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM tb_lock_renaksi_opd
		WHERE kode_opd = ? AND tahun = ? AND sasaran_id = ? AND rekin_id = ?
	`, kodeOpd, tahun, sasaranId, rekinId).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("LockRenaksiOpdRepository.IsLocked: %w", err)
	}
	return count > 0, nil
}
