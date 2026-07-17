package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
)

type IsuKlhsRepositoryImpl struct {
}

func NewIsuKlhsRepositoryImpl() *IsuKlhsRepositoryImpl {
	return &IsuKlhsRepositoryImpl{}
}

func (repository *IsuKlhsRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, isu domain.IsuKlhs) (domain.IsuKlhs, error) {

	script := `
		INSERT INTO tb_isu_klhs 
		(kode_bidang_urusan, kode_opd, isu_klhs, tahun) 
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
		return domain.IsuKlhs{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.IsuKlhs{}, err
	}

	isu.ID = int(id)

	return isu, nil
}

func (repository *IsuKlhsRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, isu domain.IsuKlhs) (domain.IsuKlhs, error) {

	// ================= UPDATE Isu =================
	query := `
		UPDATE tb_isu_klhs
		SET
			kode_bidang_urusan = ?,
			kode_opd = ?,
			isu_klhs = ?,
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
		return domain.IsuKlhs{}, err
	}

	return isu, nil
}

func (repository *IsuKlhsRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_isu_klhs WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *IsuKlhsRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.IsuKlhs, error) {

	// ================= IKK =================
	query := `
		SELECT
			id,
			kode_bidang_urusan,
			kode_opd,
			isu_klhs,
			tahun
		FROM tb_isu_global
		WHERE id = ?
	`

	var result domain.IsuKlhs

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&result.ID,
		&result.KodeBidangUrusan,
		&result.KodeOpd,
		&result.Isu,
		&result.Tahun,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.IsuKlhs{}, errors.New("isu KLHS tidak ditemukan")
		}
		return domain.IsuKlhs{}, err
	}
	
	return result, nil
}

func (repository *IsuKlhsRepositoryImpl) FindAllById(ctx context.Context, tx *sql.Tx, id int) (domain.IsuKlhs, error) {

	query := `
		SELECT isu.id, 
			   isu.kode_opd, 
			   od.nama_opd,
			   isu.kode_bidang_urusan, 
			   bu.nama_bidang_urusan, 
			   isu.isu_klhs, 
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

	var item domain.IsuKlhs

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
			return domain.IsuKlhs{}, nil
		}
		return domain.IsuKlhs{}, err
	}

	return item, nil
}