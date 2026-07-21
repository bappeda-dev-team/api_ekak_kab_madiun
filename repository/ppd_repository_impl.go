package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
)

type PpdRepositoryImpl struct {
}

func NewPpdRepositoryImpl() *PpdRepositoryImpl {
	return &PpdRepositoryImpl{}
}

func (repository *PpdRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, ppd domain.PotensiPerangkatDaerah) (domain.PotensiPerangkatDaerah, error) {

	script := `
		INSERT INTO tb_ppd 
		(kode_bidang_urusan, kode_opd, potensi, tahun) 
		VALUES (?, ?, ?, ?)
	`

	result, err := tx.ExecContext(
		ctx,
		script,
		ppd.KodeBidangUrusan,
		ppd.KodeOpd,
		ppd.Potensi,
		ppd.Tahun,
	)
	if err != nil {
		return domain.PotensiPerangkatDaerah{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.PotensiPerangkatDaerah{}, err
	}

	ppd.ID = int(id)

	return ppd, nil
}

func (repository *PpdRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, ppd domain.PotensiPerangkatDaerah) (domain.PotensiPerangkatDaerah, error) {

	// ================= UPDATE Isu =================
	query := `
		UPDATE tb_ppd
		SET
			kode_bidang_urusan = ?,
			kode_opd = ?,
			potensi = ?,
			tahun = ?
		WHERE id = ?
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		ppd.KodeBidangUrusan,
		ppd.KodeOpd,
		ppd.Potensi,
		ppd.Tahun,
		ppd.ID,
	)
	if err != nil {
		return domain.PotensiPerangkatDaerah{}, err
	}

	return ppd, nil
}

func (repository *PpdRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_ppd WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *PpdRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.PotensiPerangkatDaerah, error) {

	// ================= IKK =================
	query := `
		SELECT
			id,
			kode_bidang_urusan,
			kode_opd,
			potensi,
			tahun
		FROM tb_ppd
		WHERE id = ?
	`

	var result domain.PotensiPerangkatDaerah

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&result.ID,
		&result.KodeBidangUrusan,
		&result.KodeOpd,
		&result.Potensi,
		&result.Tahun,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.PotensiPerangkatDaerah{}, errors.New("Potensi perangkat daerah tidak ditemukan")
		}
		return domain.PotensiPerangkatDaerah{}, err
	}
	
	return result, nil
}

func (repository *PpdRepositoryImpl) FindAllById(ctx context.Context, tx *sql.Tx, id int) (domain.PotensiPerangkatDaerah, error) {

	query := `
		SELECT isu.id, 
			   isu.kode_opd, 
			   od.nama_opd,
			   isu.kode_bidang_urusan, 
			   bu.nama_bidang_urusan, 
			   isu.potensi, 
			   isu.tahun
		FROM tb_ppd isu
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

	var item domain.PotensiPerangkatDaerah

	err := row.Scan(
		&item.ID,
		&item.KodeOpd,
		&item.NamaOpd,
		&item.KodeBidangUrusan,
		&item.NamaBidangUrusan,
		&item.Potensi,
		&item.Tahun,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.PotensiPerangkatDaerah{}, nil
		}
		return domain.PotensiPerangkatDaerah{}, err
	}

	return item, nil
}

func (repository *PpdRepositoryImpl) FindSelectionByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.BidangUrusanSelection, error) {

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

func (repository *PpdRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.PotensiPerangkatDaerah, error) {

	query := `
		SELECT isu.id, 
			   isu.kode_opd, 
			   od.nama_opd,
			   isu.kode_bidang_urusan, 
			   bu.nama_bidang_urusan, 
			   isu.potensi, 
			   isu.tahun
		FROM tb_ppd isu
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

	isuMap := make(map[int]*domain.PotensiPerangkatDaerah)
	isuIDs := make([]int, 0)

	for rows.Next() {
		var item domain.PotensiPerangkatDaerah

		err := rows.Scan(
			&item.ID,
			&item.KodeOpd,
			&item.NamaOpd,
			&item.KodeBidangUrusan,
			&item.NamaBidangUrusan,
			&item.Potensi,
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
		return []domain.PotensiPerangkatDaerah{}, nil
	}

	result := make([]domain.PotensiPerangkatDaerah, 0, len(isuMap))
	for _, v := range isuMap {
		result = append(result, *v)
	}

	return result, nil
}