package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
	"fmt"
	"strings"
)

type NspkRepositoryImpl struct {
}

func NewNspkRepositoryImpl() *NspkRepositoryImpl {
	return &NspkRepositoryImpl{}
}

func (repository *NspkRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, nspk domain.MasterNSPK) (domain.MasterNSPK, error) {

	script := `
		INSERT INTO tb_master_nspk 
		(kode_opd, nspk, tahun) 
		VALUES (?, ?, ?)
	`

	result, err := tx.ExecContext(
		ctx,
		script,
		nspk.KodeOpd,
		nspk.NSPK,
		nspk.Tahun,
	)
	if err != nil {
		return domain.MasterNSPK{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.MasterNSPK{}, err
	}

	nspk.ID = int(id)

	return nspk, nil
}

func (repository *NspkRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, nspk domain.MasterNSPK) (domain.MasterNSPK, error) {

	// ================= UPDATE Isu =================
	query := `
		UPDATE tb_master_nspk
		SET
			kode_opd = ?,
			nspk = ?,
			tahun = ?
		WHERE id = ?
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		nspk.KodeOpd,
		nspk.NSPK,
		nspk.Tahun,
		nspk.ID,
	)
	if err != nil {
		return domain.MasterNSPK{}, err
	}

	return nspk, nil
}

func (repository *NspkRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_master_nspk WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *NspkRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.MasterNSPK, error) {

	// ================= IKK =================
	query := `
		SELECT
			id,
			kode_opd,
			nspk,
			tahun
		FROM tb_master_nspk
		WHERE id = ?
	`

	var result domain.MasterNSPK

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&result.ID,
		&result.KodeOpd,
		&result.NSPK,
		&result.Tahun,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.MasterNSPK{}, errors.New("Norma, Standar, Prosedur dan Kriteria tidak ditemukan")
		}
		return domain.MasterNSPK{}, err
	}

	return result, nil
}

func (repository *NspkRepositoryImpl) FindByIds(ctx context.Context, tx *sql.Tx, ids []int) ([]domain.MasterNSPK, error) {

	if len(ids) == 0 {
		return []domain.MasterNSPK{}, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))

	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			kode_opd,
			nspk,
			tahun
		FROM tb_master_nspk
		WHERE id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []domain.MasterNSPK{}

	for rows.Next() {
		var result domain.MasterNSPK

		err := rows.Scan(
			&result.ID,
			&result.KodeOpd,
			&result.NSPK,
			&result.Tahun,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return []domain.MasterNSPK{}, nil
	}

	return results, nil
}

func (repository *NspkRepositoryImpl) FindAllById(ctx context.Context, tx *sql.Tx, id int) (domain.MasterNSPK, error) {

	query := `
		SELECT isu.id, 
			   isu.kode_opd, 
			   od.nama_opd,
			   isu.nspk, 
			   isu.tahun
		FROM tb_master_nspk isu
		LEFT JOIN tb_operasional_daerah od
		ON od.kode_opd = isu.kode_opd
	`

	args := make([]interface{}, 0)

	if id != 0 {
		query += " WHERE isu.id = ?"
		args = append(args, id)
	}

	row := tx.QueryRowContext(ctx, query, args...)

	var item domain.MasterNSPK

	err := row.Scan(
		&item.ID,
		&item.KodeOpd,
		&item.NamaOpd,
		&item.NSPK,
		&item.Tahun,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.MasterNSPK{}, nil
		}
		return domain.MasterNSPK{}, err
	}

	return item, nil
}

func (repository *NspkRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.MasterNSPK, error) {

	query := `
		SELECT isu.id, 
			   isu.kode_opd, 
			   od.nama_opd,
			   isu.nspk, 
			   isu.tahun
		FROM tb_master_nspk isu
		LEFT JOIN tb_operasional_daerah od
		ON od.kode_opd = isu.kode_opd
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

	isuMap := make(map[int]*domain.MasterNSPK)
	isuIDs := make([]int, 0)

	for rows.Next() {
		var item domain.MasterNSPK

		err := rows.Scan(
			&item.ID,
			&item.KodeOpd,
			&item.NamaOpd,
			&item.NSPK,
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
		return []domain.MasterNSPK{}, nil
	}

	result := make([]domain.MasterNSPK, 0, len(isuMap))
	for _, v := range isuMap {
		result = append(result, *v)
	}

	return result, nil
}