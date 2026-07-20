package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
)

type IsuRegionalRepositoryImpl struct {
}

func NewIsuRegionalRepositoryImpl() *IsuRegionalRepositoryImpl {
	return &IsuRegionalRepositoryImpl{}
}

func (repository *IsuRegionalRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, isu domain.IsuRegional) (domain.IsuRegional, error) {

	script := `
		INSERT INTO tb_isu_regional 
		(kode_bidang_urusan, kode_opd, isu_regional, tahun) 
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
		return domain.IsuRegional{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.IsuRegional{}, err
	}

	isu.ID = int(id)

	return isu, nil
}

func (repository *IsuRegionalRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, isu domain.IsuRegional) (domain.IsuRegional, error) {

	// ================= UPDATE Isu =================
	query := `
		UPDATE tb_isu_regional
		SET
			kode_bidang_urusan = ?,
			kode_opd = ?,
			isu_regional = ?,
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
		return domain.IsuRegional{}, err
	}

	return isu, nil
}

func (repository *IsuRegionalRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_isu_regional WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *IsuRegionalRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.IsuRegional, error) {

	// ================= IKK =================
	query := `
		SELECT
			id,
			kode_bidang_urusan,
			kode_opd,
			isu_regional,
			tahun
		FROM tb_isu_regional
		WHERE id = ?
	`

	var result domain.IsuRegional

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&result.ID,
		&result.KodeBidangUrusan,
		&result.KodeOpd,
		&result.Isu,
		&result.Tahun,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.IsuRegional{}, errors.New("isu regional tidak ditemukan")
		}
		return domain.IsuRegional{}, err
	}
	
	return result, nil
}

func (repository *IsuRegionalRepositoryImpl) FindAllById(ctx context.Context, tx *sql.Tx, id int) (domain.IsuRegional, error) {

	query := `
		SELECT isu.id, 
			   isu.kode_opd, 
			   od.nama_opd,
			   isu.kode_bidang_urusan, 
			   bu.nama_bidang_urusan, 
			   isu.isu_regional, 
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

	var item domain.IsuRegional

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
			return domain.IsuRegional{}, nil
		}
		return domain.IsuRegional{}, err
	}

	return item, nil
}

func (repository *IsuRegionalRepositoryImpl) FindSelectionByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.BidangUrusanSelection, error) {

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

func (repository *IsuRegionalRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.IsuRegional, error) {

	query := `
		SELECT isu.id, 
			   isu.kode_opd, 
			   od.nama_opd,
			   isu.kode_bidang_urusan, 
			   bu.nama_bidang_urusan, 
			   isu.isu_regional, 
			   isu.tahun
		FROM tb_isu_regional isu
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

	isuMap := make(map[int]*domain.IsuRegional)
	isuIDs := make([]int, 0)

	for rows.Next() {
		var item domain.IsuRegional

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
		return []domain.IsuRegional{}, nil
	}

	result := make([]domain.IsuRegional, 0, len(isuMap))
	for _, v := range isuMap {
		result = append(result, *v)
	}

	return result, nil
}