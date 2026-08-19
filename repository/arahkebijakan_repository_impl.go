package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
)

type ArahKebijakanRepositoryImpl struct {
}

func NewArahKebijakanRepositoryImpl() *ArahKebijakanRepositoryImpl {
	return &ArahKebijakanRepositoryImpl{}
}

func (repository *ArahKebijakanRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, ar domain.ArahKebijakanOpd) (domain.ArahKebijakanOpd, error) {

	script := `
		INSERT INTO tb_arah_kebijakan 
		(pokin_id, arah_kebijakan, kode_opd, tahun) 
		VALUES (?, ?, ?, ?)
	`

	result, err := tx.ExecContext(
		ctx,
		script,
		ar.PokinId,
		ar.Arah,
		ar.KodeOpd,
		ar.Tahun,
	)
	if err != nil {
		return domain.ArahKebijakanOpd{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.ArahKebijakanOpd{}, err
	}

	ar.ID = int(id)

	return ar, nil
}

func (repository *ArahKebijakanRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, ar domain.ArahKebijakanOpd) (domain.ArahKebijakanOpd, error) {

	// ================= UPDATE Isu =================
	query := `
		UPDATE tb_arah_kebijakan
		SET
			pokin_id = ?,
			arah_kebijakan = ?,
			kode_opd = ?,
			tahun = ?
		WHERE id = ?
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		ar.PokinId,
		ar.Arah,
		ar.KodeOpd,
		ar.Tahun,
		ar.ID,
	)
	if err != nil {
		return domain.ArahKebijakanOpd{}, err
	}

	return ar, nil
}

func (repository *ArahKebijakanRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.ArahKebijakanOpd, error) {

	// ================= IKK =================
	query := `
		SELECT
			id,
			pokin_id,
			arah_kebijakan,
			kode_opd,
			tahun
		FROM tb_master_nspk
		WHERE id = ?
	`

	var result domain.ArahKebijakanOpd

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&result.ID,
		&result.PokinId,
		&result.Arah,
		&result.KodeOpd,
		&result.Tahun,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ArahKebijakanOpd{}, errors.New("Arah kebijakan tidak ditemukan")
		}
		return domain.ArahKebijakanOpd{}, err
	}

	return result, nil
}
