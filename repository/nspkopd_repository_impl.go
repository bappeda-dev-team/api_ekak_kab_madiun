package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
)

type NspkOpdRepositoryImpl struct {
}

func NewNspkOpdRepositoryImpl() *NspkOpdRepositoryImpl {
	return &NspkOpdRepositoryImpl{}
}

func (repository *NspkOpdRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, nspk domain.NspkOpd) (domain.NspkOpd, error) {

	script := `
		INSERT INTO tb_nspk_opd 
		(kode_opd, id_nspk, id_tujuan_opd, id_sasaran_opd, tahun) 
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := tx.ExecContext(
		ctx,
		script,
		nspk.KodeOpd,
		nspk.IdNspk,
		nspk.IdTujuanOpd,
		nspk.IdSasaranOpd,
		nspk.Tahun,
	)
	if err != nil {
		return domain.NspkOpd{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.NspkOpd{}, err
	}

	nspk.ID = int(id)

	return nspk, nil
}

func (repository *NspkOpdRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, nspk domain.NspkOpd) (domain.NspkOpd, error) {

	// ================= UPDATE Isu =================
	query := `
		UPDATE tb_nspk_opd
		SET
			kode_opd = ?,
			id_nspk = ?,
			id_tujuan_opd = ?,
			id_sasaran_opd = ?,
			tahun = ?
		WHERE id = ?
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		nspk.KodeOpd,
		nspk.IdNspk,
		nspk.IdTujuanOpd,
		nspk.IdSasaranOpd,
		nspk.Tahun,
		nspk.ID,
	)
	if err != nil {
		return domain.NspkOpd{}, err
	}

	return nspk, nil
}

func (repository *NspkOpdRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_nspk_opd WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *NspkOpdRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.NspkOpd, error) {

	// ================= IKK =================
	query := `
		SELECT
			id,
			kode_opd,
			id_nspk,
			id_tujuan_opd,
			id_sasaran_opd,
			tahun
		FROM tb_nspk_opd
		WHERE id = ?
	`

	var result domain.NspkOpd

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&result.ID,
		&result.KodeOpd,
		&result.IdNspk,
		&result.IdTujuanOpd,
		&result.IdSasaranOpd,
		&result.Tahun,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.NspkOpd{}, errors.New("Norma, Standar, Prosedur dan Kriteria opd tidak ditemukan")
		}
		return domain.NspkOpd{}, err
	}

	return result, nil
}

func (repository *NspkOpdRepositoryImpl) FindAllById(ctx context.Context, tx *sql.Tx, id int) (domain.NspkOpd, error) {

	query := `
		SELECT isu.id, 
			   isu.kode_opd, 
			   od.nama_opd,
			   isu.id_nspk, 
			   mn.nspk, 
			   isu.id_tujuan_opd, 
			   to.tujuan, 
			   isu.id_sasaran_opd, 
			   so.nama_sasaran_opd, 
			   isu.tahun
		FROM tb_nspk_opd isu

		LEFT JOIN tb_operasional_daerah od
		ON od.kode_opd = isu.kode_opd

		LEFT JOIN tb_master_nspk mn
		ON mn.id = isu.id_nspk

		LEFT JOIN tb_tujuan_opd to
		ON to.id = isu.id_tujuan_opd

		LEFT JOIN tb_sasaran_opd so
		ON so.id = isu.id_sasaran_opd
	`

	args := make([]interface{}, 0)

	if id != 0 {
		query += " WHERE isu.id = ?"
		args = append(args, id)
	}

	row := tx.QueryRowContext(ctx, query, args...)

	var item domain.NspkOpd

	err := row.Scan(
		&item.ID,
		&item.KodeOpd,
		&item.NamaOpd,
		&item.IdNspk,
		&item.NSPK,
		&item.IdTujuanOpd,
		&item.TujuanOpd,
		&item.IdSasaranOpd,
		&item.SasaranOpd,
		&item.Tahun,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.NspkOpd{}, nil
		}
		return domain.NspkOpd{}, err
	}

	return item, nil
}

func (repository *NspkOpdRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.NspkOpd, error) {

	query := `
		SELECT isu.id, 
			   isu.kode_opd, 
			   od.nama_opd,
			   isu.id_nspk, 
			   mn.nspk, 
			   isu.id_tujuan_opd, 
			   tto.tujuan, 
			   isu.id_sasaran_opd, 
			   so.nama_sasaran_opd, 
			   isu.tahun
		FROM tb_nspk_opd isu

		LEFT JOIN tb_operasional_daerah od
		ON od.kode_opd = isu.kode_opd

		LEFT JOIN tb_master_nspk mn
		ON mn.id = isu.id_nspk

		LEFT JOIN tb_tujuan_opd tto
		ON tto.id = isu.id_tujuan_opd

		LEFT JOIN tb_sasaran_opd so
		ON so.id = isu.id_sasaran_opd
	`

	args := make([]interface{}, 0)

	if kodeOpd != "" {
		query += " WHERE isu.kode_opd = ?"
		args = append(args, kodeOpd)
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	isuMap := make(map[int]*domain.NspkOpd)
	isuIDs := make([]int, 0)

	for rows.Next() {
		var item domain.NspkOpd

		err := rows.Scan(
			&item.ID,
			&item.KodeOpd,
			&item.NamaOpd,
			&item.IdNspk,
			&item.NSPK,
			&item.IdTujuanOpd,
			&item.TujuanOpd,
			&item.IdSasaranOpd,
			&item.SasaranOpd,
			&item.Tahun,
		)
		if err != nil {
			return nil, err
		}

		copyItem := item
		isuMap[item.ID] = &copyItem
		isuIDs = append(isuIDs, item.ID)
	}

	if len(isuIDs) == 0 {
		return []domain.NspkOpd{}, nil
	}

	result := make([]domain.NspkOpd, 0, len(isuMap))
	for _, v := range isuMap {
		result = append(result, *v)
	}

	return result, nil
}
