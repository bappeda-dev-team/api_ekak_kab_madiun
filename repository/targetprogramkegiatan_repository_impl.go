package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
)

type TargetProgramKegiatanRepositoryImpl struct{}

func NewTargetProgramKegiatanRepositoryImpl() *TargetProgramKegiatanRepositoryImpl {
	return &TargetProgramKegiatanRepositoryImpl{}
}

func (r *TargetProgramKegiatanRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, t domain.TargetProgramKegiatan) (domain.TargetProgramKegiatan, error) {
	res, err := tx.ExecContext(ctx, `
		INSERT INTO tb_target_program_kegiatan (kode_indikator, kode_opd, jenis, target, satuan, tahun)
		VALUES (?, ?, ?, ?, ?, ?)`,
		t.KodeIndikator, t.KodeOpd, t.Jenis, t.Target, t.Satuan, t.Tahun,
	)
	if err != nil {
		return domain.TargetProgramKegiatan{}, fmt.Errorf("gagal membuat target program kegiatan: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.TargetProgramKegiatan{}, err
	}
	t.Id = int(id)
	return t, nil
}

func (r *TargetProgramKegiatanRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, t domain.TargetProgramKegiatan) (domain.TargetProgramKegiatan, error) {
	_, err := tx.ExecContext(ctx, `
		UPDATE tb_target_program_kegiatan
		SET target = ?, satuan = ?
		WHERE id = ?`,
		t.Target, t.Satuan, t.Id,
	)
	if err != nil {
		return domain.TargetProgramKegiatan{}, fmt.Errorf("gagal mengupdate target program kegiatan id=%d: %w", t.Id, err)
	}
	return t, nil
}

func (r *TargetProgramKegiatanRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	_, err := tx.ExecContext(ctx, `DELETE FROM tb_target_program_kegiatan WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus target program kegiatan id=%d: %w", id, err)
	}
	return nil
}

func (r *TargetProgramKegiatanRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.TargetProgramKegiatan, error) {
	var t domain.TargetProgramKegiatan
	err := tx.QueryRowContext(ctx, `
		SELECT id, kode_indikator, kode_opd, jenis, target, satuan, tahun
		FROM tb_target_program_kegiatan
		WHERE id = ?`, id,
	).Scan(&t.Id, &t.KodeIndikator, &t.KodeOpd, &t.Jenis, &t.Target, &t.Satuan, &t.Tahun)
	if err == sql.ErrNoRows {
		return domain.TargetProgramKegiatan{}, fmt.Errorf("target program kegiatan dengan id=%d tidak ditemukan", id)
	}
	if err != nil {
		return domain.TargetProgramKegiatan{}, fmt.Errorf("gagal mencari target program kegiatan: %w", err)
	}
	return t, nil
}

func (r *TargetProgramKegiatanRepositoryImpl) FindByKodeIndikatorAndKodeOpd(ctx context.Context, tx *sql.Tx, kodeIndikator, kodeOpd string) ([]domain.TargetProgramKegiatan, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT id, kode_indikator, kode_opd, jenis, target, satuan, tahun
		FROM tb_target_program_kegiatan
		WHERE kode_indikator = ? AND kode_opd = ?
		ORDER BY jenis ASC, tahun ASC`,
		kodeIndikator, kodeOpd,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal query target program kegiatan: %w", err)
	}
	defer rows.Close()

	var result []domain.TargetProgramKegiatan
	for rows.Next() {
		var t domain.TargetProgramKegiatan
		if err := rows.Scan(&t.Id, &t.KodeIndikator, &t.KodeOpd, &t.Jenis, &t.Target, &t.Satuan, &t.Tahun); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, rows.Err()
}
