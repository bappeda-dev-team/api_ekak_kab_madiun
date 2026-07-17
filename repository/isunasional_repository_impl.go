package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
)

type IsuNasionalRepositoryImpl struct {
}

func NewIsuNasionalRepositoryImpl() *IsuNasionalRepositoryImpl {
	return &IsuNasionalRepositoryImpl{}
}

func (repository *IsuNasionalRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, isu domain.IsuNasional) (domain.IsuNasional, error) {

	script := `
		INSERT INTO tb_isu_nasional 
		(kode_bidang_urusan, kode_opd, isu_nasional, tahun) 
		VALUES (?, ?, ?, ?)
	`

	result, err := tx.ExecContext(
		ctx,
		script,
		isu.KodeBidangUrusan,
		isu.KodeOpd,
		isu.Isu,
		isu.Tahun,
	)
	if err != nil {
		return domain.IsuNasional{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.IsuNasional{}, err
	}

	isu.ID = int(id)

	return isu, nil
}

func (repository *IsuNasionalRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, isu domain.IsuNasional) (domain.IsuNasional, error) {

	// ================= UPDATE Isu =================
	query := `
		UPDATE tb_isu_nasional
		SET
			kode_bidang_urusan = ?,
			kode_opd = ?,
			isu_nasional = ?,
			tahun = ?
		WHERE id = ?
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		isu.KodeBidangUrusan,
		isu.KodeOpd,
		isu.Isu,
		isu.Tahun,
		isu.ID,
	)
	if err != nil {
		return domain.IsuNasional{}, err
	}

	return isu, nil
}

func (repository *IsuNasionalRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_isu_nasional WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *IsuNasionalRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.IsuNasional, error) {

	// ================= IKK =================
	query := `
		SELECT
			id,
			kode_bidang_urusan,
			kode_opd,
			isu_nasional,
			tahun
		FROM tb_isu_global
		WHERE id = ?
	`

	var result domain.IsuNasional

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&result.ID,
		&result.KodeBidangUrusan,
		&result.KodeOpd,
		&result.Isu,
		&result.Tahun,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.IsuNasional{}, errors.New("isu nasional tidak ditemukan")
		}
		return domain.IsuNasional{}, err
	}
	
	return result, nil
}

func (repository *IsuNasionalRepositoryImpl) FindAllById(ctx context.Context, tx *sql.Tx, id int) (domain.IsuNasional, error) {

	query := `
		SELECT isu.id, 
			   isu.kode_opd, 
			   od.nama_opd,
			   isu.kode_bidang_urusan, 
			   bu.nama_bidang_urusan, 
			   isu.isu_nasional, 
			   isu.tahun
		FROM tb_isu_global isu
		LEFT JOIN tb_operasional_daerah od
		ON od.kode_opd = isu.kode_opd
		LEFT JOIN tb_bidang_urusan bu
		ON bu.kode_bidang_urusan = isu.kode_bidang_urusan
	`

	args := make([]interface{}, 0)

	if id != 0 {
		query += " WHERE isu.id = ?"
		args = append(args, id)
	}

	row := tx.QueryRowContext(ctx, query, args...)

	var item domain.IsuNasional

	err := row.Scan(
		&item.ID,
		&item.KodeOpd,
		&item.NamaOpd,
		&item.KodeBidangUrusan,
		&item.NamaBidangUrusan,
		&item.Isu,
		&item.Tahun,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.IsuNasional{}, nil
		}
		return domain.IsuNasional{}, err
	}

	return item, nil
}