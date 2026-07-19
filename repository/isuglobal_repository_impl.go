package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
)

type IsuGlobalRepositoryImpl struct {
}

func NewIsuGlobalRepositoryImpl() *IsuGlobalRepositoryImpl {
	return &IsuGlobalRepositoryImpl{}
}

func (repository *IsuGlobalRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, isu domain.IsuGlobal) (domain.IsuGlobal, error) {

	script := `
		INSERT INTO tb_isu_global 
		(kode_bidang_urusan, kode_opd, isu_global, tahun) 
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
		return domain.IsuGlobal{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.IsuGlobal{}, err
	}

	isu.ID = int(id)

	return isu, nil
}

func (repository *IsuGlobalRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, isu domain.IsuGlobal) (domain.IsuGlobal, error) {

	// ================= UPDATE IKK =================
	query := `
		UPDATE tb_isu_global
		SET
			kode_bidang_urusan = ?,
			kode_opd = ?,
			isu_global = ?,
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
		return domain.IsuGlobal{}, err
	}

	return isu, nil
}

func (repository *IsuGlobalRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_isu_global WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *IsuGlobalRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.IsuGlobal, error) {

	// ================= IKK =================
	query := `
		SELECT
			id,
			kode_bidang_urusan,
			kode_opd,
			isu_global,
			tahun
		FROM tb_isu_global
		WHERE id = ?
	`

	var result domain.IsuGlobal

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&result.ID,
		&result.KodeBidangUrusan,
		&result.KodeOpd,
		&result.Isu,
		&result.Tahun,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.IsuGlobal{}, errors.New("isu global tidak ditemukan")
		}
		return domain.IsuGlobal{}, err
	}
	
	return result, nil
}

func (repository *IsuGlobalRepositoryImpl) FindAllById(ctx context.Context, tx *sql.Tx, id int) (domain.IsuGlobal, error) {

	query := `
		SELECT isu.id, 
			   isu.kode_opd, 
			   od.nama_opd,
			   isu.kode_bidang_urusan, 
			   bu.nama_bidang_urusan, 
			   isu.isu_global, 
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

	var item domain.IsuGlobal

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
			return domain.IsuGlobal{}, nil
		}
		return domain.IsuGlobal{}, err
	}

	return item, nil
}

func (repository *IsuGlobalRepositoryImpl) FindSelectionByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.BidangUrusanSelection, error) {

	query := `SELECT
				bu.kode_bidang_urusan,
				COALESCE(bu.nama_bidang_urusan, '') AS nama_bidang_urusan,
				od.kode_opd,
				COALESCE(od.nama_opd, '') AS nama_opd
			FROM tb_bidang_urusan bu
			CROSS JOIN tb_operasional_daerah od
			WHERE od.kode_opd = ?`

	rows, err := tx.QueryContext(ctx, query, kodeOpd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	selections := make([]domain.BidangUrusanSelection, 0)

	for rows.Next() {
		var selection domain.BidangUrusanSelection

		err := rows.Scan(
			&selection.KodeBidangUrusan,
			&selection.NamaBidangUrusan,
			&selection.KodeOpd,
			&selection.NamaOpd,
		)

		if err != nil {
			return nil, err
		}

		selections = append(selections, selection)
	}

	return selections, nil
}

func (repository *IsuGlobalRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.IsuGlobal, error) {

	query := `
		SELECT isu.id, 
			   isu.kode_opd, 
			   od.nama_opd,
			   isu.kode_bidang_urusan, 
			   bu.nama_bidang_urusan, 
			   isu.isu_global, 
			   isu.tahun, 
		FROM tb_isu_global isu
		LEFT JOIN tb_operasional_daerah od
		ON od.kode_opd = isu.kode_opd
		LEFT JOIN tb_bidang_urusan bu
		ON bu.kode_bidang_urusan = isu.kode_bidang_urusan
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

	isuMap := make(map[int]*domain.IsuGlobal)
	isuIDs := make([]int, 0)

	for rows.Next() {
		var item domain.IsuGlobal

		err := rows.Scan(
			&item.ID,
			&item.KodeOpd,
			&item.NamaOpd,
			&item.KodeBidangUrusan,
			&item.NamaBidangUrusan,
			&item.Isu,
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
		return []domain.IsuGlobal{}, nil
	}

	result := make([]domain.IsuGlobal, 0, len(isuMap))
	for _, v := range isuMap {
		result = append(result, *v)
	}

	return result, nil
}