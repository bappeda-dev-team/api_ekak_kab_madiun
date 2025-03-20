package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
	"fmt"
)

type PeriodeRepositoryImpl struct {
}

func NewPeriodeRepositoryImpl() *PeriodeRepositoryImpl {
	return &PeriodeRepositoryImpl{}
}

func (repository *PeriodeRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, periode domain.Periode) (domain.Periode, error) {
	query := "INSERT INTO tb_periode(id, tahun_awal, tahun_akhir, jenis_periode) VALUES (?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, periode.Id, periode.TahunAwal, periode.TahunAkhir, periode.JenisPeriode)
	if err != nil {
		return periode, err
	}
	return periode, nil
}

func (repository *PeriodeRepositoryImpl) IsIdExists(ctx context.Context, tx *sql.Tx, id int) bool {
	query := "SELECT COUNT(*) FROM tb_periode WHERE id = ?"
	var count int
	err := tx.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		return true // Assume exists on error to be safe
	}
	return count > 0
}

func (repository *PeriodeRepositoryImpl) SaveTahunPeriode(ctx context.Context, tx *sql.Tx, tahunPeriode domain.TahunPeriode) error {
	query := "INSERT INTO tb_tahun_periode(id_periode, tahun) VALUES (?, ?)"
	_, err := tx.ExecContext(ctx, query, tahunPeriode.IdPeriode, tahunPeriode.Tahun)
	return err
}

func (repository *PeriodeRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, periodeId int) (domain.Periode, error) {
	query := "SELECT id, tahun_awal, tahun_akhir, jenis_periode FROM tb_periode WHERE id = ?"
	rows, err := tx.QueryContext(ctx, query, periodeId)
	if err != nil {
		return domain.Periode{}, err
	}
	defer rows.Close()

	periode := domain.Periode{}
	if rows.Next() {
		err := rows.Scan(&periode.Id, &periode.TahunAwal, &periode.TahunAkhir, &periode.JenisPeriode)
		if err != nil {
			return periode, err
		}
		return periode, nil
	}

	return periode, errors.New("periode not found")
}

func (repository *PeriodeRepositoryImpl) FindOverlappingPeriodes(ctx context.Context, tx *sql.Tx, tahunAwal, tahunAkhir, jenisPeriode string) ([]domain.Periode, error) {
	query := `
		SELECT id, tahun_awal, tahun_akhir, jenis_periode 
		FROM tb_periode 
		WHERE jenis_periode = ? AND (
			(tahun_awal <= ? AND tahun_akhir >= ?) 
			OR (tahun_awal <= ? AND tahun_akhir >= ?)
			OR (tahun_awal >= ? AND tahun_akhir <= ?)
		)`

	rows, err := tx.QueryContext(ctx, query, jenisPeriode, tahunAkhir, tahunAwal, tahunAkhir, tahunAwal, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var periodes []domain.Periode
	for rows.Next() {
		periode := domain.Periode{}
		err := rows.Scan(&periode.Id, &periode.TahunAwal, &periode.TahunAkhir, &periode.JenisPeriode)
		if err != nil {
			return nil, err
		}
		periodes = append(periodes, periode)
	}
	return periodes, nil
}

func (repository *PeriodeRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, periode domain.Periode) (domain.Periode, error) {
	query := "UPDATE tb_periode SET tahun_awal = ?, tahun_akhir = ?, jenis_periode = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, periode.TahunAwal, periode.TahunAkhir, periode.JenisPeriode, periode.Id)
	if err != nil {
		return periode, err
	}
	return periode, nil
}

func (repository *PeriodeRepositoryImpl) DeleteTahunPeriode(ctx context.Context, tx *sql.Tx, periodeId int) error {
	query := "DELETE FROM tb_tahun_periode WHERE id_periode = ?"
	_, err := tx.ExecContext(ctx, query, periodeId)
	return err
}

func (repository *PeriodeRepositoryImpl) FindByTahun(ctx context.Context, tx *sql.Tx, tahun string) (domain.Periode, error) {
	query := `
		SELECT p.id, p.tahun_awal, p.tahun_akhir, p.jenis_periode
		FROM tb_periode p
		JOIN tb_tahun_periode tp ON p.id = tp.id_periode
		WHERE tp.tahun = ?
		LIMIT 1`

	rows, err := tx.QueryContext(ctx, query, tahun)
	if err != nil {
		return domain.Periode{}, err
	}
	defer rows.Close()

	if rows.Next() {
		periode := domain.Periode{}
		err := rows.Scan(&periode.Id, &periode.TahunAwal, &periode.TahunAkhir, &periode.JenisPeriode)
		if err != nil {
			return periode, err
		}
		return periode, nil
	}

	return domain.Periode{}, errors.New("periode not found")
}

func (repository *PeriodeRepositoryImpl) FindOverlappingPeriodesExcludeCurrent(ctx context.Context, tx *sql.Tx, currentId int, tahunAwal, tahunAkhir, jenisPeriode string) ([]domain.Periode, error) {
	query := `
		SELECT id, tahun_awal, tahun_akhir, jenis_periode 
		FROM tb_periode 
		WHERE id != ? AND jenis_periode = ? AND (
			(tahun_awal <= ? AND tahun_akhir >= ?) 
			OR (tahun_awal <= ? AND tahun_akhir >= ?)
			OR (tahun_awal >= ? AND tahun_akhir <= ?)
		)`

	rows, err := tx.QueryContext(ctx, query, currentId, jenisPeriode, tahunAkhir, tahunAwal, tahunAkhir, tahunAwal, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var periodes []domain.Periode
	for rows.Next() {
		periode := domain.Periode{}
		err := rows.Scan(&periode.Id, &periode.TahunAwal, &periode.TahunAkhir, &periode.JenisPeriode)
		if err != nil {
			return nil, err
		}
		periodes = append(periodes, periode)
	}
	return periodes, nil
}

func (repository *PeriodeRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, jenis_periode string) ([]domain.Periode, error) {
	query := "SELECT id, tahun_awal, tahun_akhir, jenis_periode FROM tb_periode WHERE 1=1"

	var params []interface{}

	if jenis_periode != "" {
		query += " AND jenis_periode = ?"
		params = append(params, jenis_periode)
	}

	query += " ORDER BY id ASC"

	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var periodes []domain.Periode
	for rows.Next() {
		periode := domain.Periode{}
		err := rows.Scan(&periode.Id, &periode.TahunAwal, &periode.TahunAkhir, &periode.JenisPeriode)
		if err != nil {
			return nil, err
		}
		periodes = append(periodes, periode)
	}
	return periodes, nil
}

func (repository *PeriodeRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, periodeId int) error {
	queryDeleteTahun := "DELETE FROM tb_tahun_periode WHERE id_periode = ?"
	_, err := tx.ExecContext(ctx, queryDeleteTahun, periodeId)
	if err != nil {
		return err
	}

	// Hapus periode
	queryDeletePeriode := "DELETE FROM tb_periode WHERE id = ?"
	result, err := tx.ExecContext(ctx, queryDeletePeriode, periodeId)
	if err != nil {
		return err
	}

	// Periksa apakah data berhasil dihapus
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("periode tidak ditemukan")
	}

	return nil

}

func (repository *PeriodeRepositoryImpl) FindRPJMDByTahun(ctx context.Context, tx *sql.Tx, tahun string) (domain.Periode, error) {
	query := `
        SELECT id, tahun_awal, tahun_akhir 
        FROM tb_periode 
        WHERE jenis_periode = 'RPJMD'
        AND ? BETWEEN tahun_awal AND tahun_akhir
    `

	rows, err := tx.QueryContext(ctx, query, tahun)
	if err != nil {
		return domain.Periode{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var periode domain.Periode
		err = rows.Scan(&periode.Id, &periode.TahunAwal, &periode.TahunAkhir)
		if err != nil {
			return domain.Periode{}, err
		}
		return periode, nil
	}

	return domain.Periode{}, fmt.Errorf("periode RPJMD tidak ditemukan")
}
