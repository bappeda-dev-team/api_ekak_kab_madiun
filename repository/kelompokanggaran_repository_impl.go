package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
)

type KelompokAnggaranRepositoryImpl struct {
}

func NewKelompokAnggaranRepositoryImpl() *KelompokAnggaranRepositoryImpl {
	return &KelompokAnggaranRepositoryImpl{}
}

func (repository *KelompokAnggaranRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, ka domain.KelompokAnggaran) (domain.KelompokAnggaran, error) {
	if repository.isKodeAnggaranExists(ctx, tx, ka.KodeKelompok) {
		return ka, errors.New("kelompok anggaran sudah ada")
	}
	script := "INSERT INTO tb_kelompok_anggaran (tahun, kelompok, kode_kelompok) VALUES (?,?,?)"
	result, err := tx.ExecContext(ctx, script, ka.Tahun, ka.Kelompok, ka.KodeKelompok)
	if err != nil {
		return ka, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	ka.Id = int(id)

	return ka, nil
}

func (repository *KelompokAnggaranRepositoryImpl) isKodeAnggaranExists(ctx context.Context, tx *sql.Tx, kodeAnggaran string) bool {
	SQL := "SELECT COUNT(*) FROM tb_kelompok_anggaran WHERE kode_kelompok = ?"
	var count int
	err := tx.QueryRowContext(ctx, SQL, kodeAnggaran).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func (repository *KelompokAnggaranRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, ka domain.KelompokAnggaran) (domain.KelompokAnggaran, error) {
	script := "UPDATE tb_kelompok_anggaran SET tahun = ?, kelompok = ?, kode_kelompok = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, ka.Tahun, ka.Kelompok, ka.KodeKelompok, ka.Id)
	if err != nil {
		return ka, err
	}
	return ka, nil
}

func (repository *KelompokAnggaranRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.KelompokAnggaran {
	SQL := "SELECT id, tahun, kelompok, kode_kelompok FROM tb_kelompok_anggaran"
	rows, err := tx.QueryContext(ctx, SQL)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var ka []domain.KelompokAnggaran
	for rows.Next() {
		KelompokAnggaran := domain.KelompokAnggaran{}
		err := rows.Scan(&KelompokAnggaran.Id, &KelompokAnggaran.Tahun, &KelompokAnggaran.Kelompok, &KelompokAnggaran.KodeKelompok)
		if err != nil {
			panic(err)
		}
		ka = append(ka, KelompokAnggaran)
	}
	return ka
}

func (repository *KelompokAnggaranRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domain.KelompokAnggaran, error) {
	SQL := "SELECT id, tahun, kelompok, kode_kelompok FROM tb_kelompok_anggaran WHERE id = ?"
	rows, err := tx.QueryContext(ctx, SQL, id)
	if err != nil {
		return domain.KelompokAnggaran{}, err
	}
	defer rows.Close()

	kelompok := domain.KelompokAnggaran{}
	if rows.Next() {
		err := rows.Scan(&kelompok.Id, &kelompok.Tahun, &kelompok.Kelompok, &kelompok.KodeKelompok)
		if err != nil {
			return domain.KelompokAnggaran{}, err
		}
		return kelompok, nil
	} else {
		return domain.KelompokAnggaran{}, errors.New("kelompok anggaran tidak ditemukan")
	}
}

func (repository *KelompokAnggaranRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	SQL := "DELETE FROM tb_kelompok_anggaran WHERE id = ?"
	result, err := tx.ExecContext(ctx, SQL, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("kelompok anggaran tidak ditemukan")
	}

	return nil
}
