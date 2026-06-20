package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
)

type IkmRepositoryImpl struct {
}

func NewIkmRepositoryImpl() *IkmRepositoryImpl {
	return &IkmRepositoryImpl{}
}

func (r *IkmRepositoryImpl) FindById(
	ctx context.Context,
	tx *sql.Tx,
	ikmId string,
) (domain.IndikatorIkm, error) {

	query := `
		SELECT
			id,
			indikator,
			kode_bidang_urusan,
			nama_bidang_urusan,
			is_active,
			definisi_operasional,
			rumus_perhitungan,
			sumber_data,
			jenis,
			tahun_awal,
			tahun_akhir,
			created_at,
			updated_at
		FROM tb_indikator_ikm
		WHERE id = ?
	`

	var res domain.IndikatorIkm

	err := tx.QueryRowContext(ctx, query, ikmId).Scan(
		&res.Id,
		&res.Indikator,
		&res.KodeBidangUrusan,
		&res.NamaBidangUrusan,
		&res.IsActive,
		&res.DefinisiOperasional,
		&res.RumusPerhitungan,
		&res.SumberData,
		&res.Jenis,
		&res.TahunAwal,
		&res.TahunAkhir,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.IndikatorIkm{}, fmt.Errorf("ikm tidak ditemukan (id=%s)", ikmId)
		}
		return domain.IndikatorIkm{}, fmt.Errorf("FindById: %w", err)
	}

	return res, nil
}

func (r *IkmRepositoryImpl) FindAllByPeriode(
	ctx context.Context,
	tx *sql.Tx,
	tahunAwal, tahunAkhir string,
) ([]domain.IndikatorIkm, error) {

	query := `
		SELECT
			id,
			indikator,
			kode_bidang_urusan,
			nama_bidang_urusan,
			is_active,
			definisi_operasional,
			rumus_perhitungan,
			sumber_data,
			jenis,
			tahun_awal,
			tahun_akhir,
			created_at,
			updated_at
		FROM tb_indikator_ikm
		WHERE tahun_awal >= ? AND tahun_akhir <= ?
	`

	rows, err := tx.QueryContext(ctx, query, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, fmt.Errorf("FindAllByPeriode: %w", err)
	}
	defer rows.Close()

	var result []domain.IndikatorIkm

	for rows.Next() {
		var res domain.IndikatorIkm

		err := rows.Scan(
			&res.Id,
			&res.Indikator,
			&res.KodeBidangUrusan,
			&res.NamaBidangUrusan,
			&res.IsActive,
			&res.DefinisiOperasional,
			&res.RumusPerhitungan,
			&res.SumberData,
			&res.Jenis,
			&res.TahunAwal,
			&res.TahunAkhir,
			&res.CreatedAt,
			&res.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, res)
	}

	return result, nil
}

func (r *IkmRepositoryImpl) ExistsById(
	ctx context.Context,
	tx *sql.Tx,
	ikmId string,
) (bool, error) {

	query := `SELECT 1 FROM tb_indikator_ikm WHERE id = ? LIMIT 1`

	var dummy int
	err := tx.QueryRowContext(ctx, query, ikmId).Scan(&dummy)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("ExistsById: %w", err)
	}

	return true, nil
}

func (r *IkmRepositoryImpl) Create(
	ctx context.Context,
	tx *sql.Tx,
	req domain.IndikatorIkm,
) (domain.IndikatorIkm, error) {

	query := `
		INSERT INTO tb_indikator_ikm (
			id,
			indikator,
			kode_bidang_urusan,
			nama_bidang_urusan,
			is_active,
			definisi_operasional,
			rumus_perhitungan,
			sumber_data,
			jenis,
			tahun_awal,
			tahun_akhir
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := tx.ExecContext(ctx, query,
		req.Id,
		req.Indikator,
		req.KodeBidangUrusan,
		req.NamaBidangUrusan,
		req.IsActive,
		req.DefinisiOperasional,
		req.RumusPerhitungan,
		req.SumberData,
		req.Jenis,
		req.TahunAwal,
		req.TahunAkhir,
	)

	if err != nil {
		return domain.IndikatorIkm{}, fmt.Errorf("Create: %w", err)
	}

	return req, nil
}

func (r *IkmRepositoryImpl) Update(
	ctx context.Context,
	tx *sql.Tx,
	req domain.IndikatorIkm,
	ikmId string,
) (domain.IndikatorIkm, error) {

	query := `
		UPDATE tb_indikator_ikm
		SET
			indikator = ?,
			kode_bidang_urusan = ?,
			nama_bidang_urusan = ?,
			is_active = ?,
			definisi_operasional = ?,
			rumus_perhitungan = ?,
			sumber_data = ?,
			jenis = ?,
			tahun_awal = ?,
			tahun_akhir = ?,
			updated_at = NOW()
		WHERE id = ?
	`

	_, err := tx.ExecContext(ctx, query,
		req.Indikator,
		req.KodeBidangUrusan,
		req.NamaBidangUrusan,
		req.IsActive,
		req.DefinisiOperasional,
		req.RumusPerhitungan,
		req.SumberData,
		req.Jenis,
		req.TahunAwal,
		req.TahunAkhir,
		ikmId,
	)

	if err != nil {
		return domain.IndikatorIkm{}, fmt.Errorf("Update: %w", err)
	}

	req.Id = ikmId
	return req, nil
}

func (r *IkmRepositoryImpl) Delete(
	ctx context.Context,
	tx *sql.Tx,
	ikmId string,
) error {

	query := `DELETE FROM tb_indikator_ikm WHERE id = ?`

	_, err := tx.ExecContext(ctx, query, ikmId)
	if err != nil {
		return fmt.Errorf("Delete: %w", err)
	}

	return nil
}
