package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
	"fmt"
	"strings"
)

type JenisInovasiRepositoryImpl struct {
}

func NewJenisInovasiRepositoryImpl() *JenisInovasiRepositoryImpl {
	return &JenisInovasiRepositoryImpl{}
}

func (repository *JenisInovasiRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, ji domain.JenisInovasi) (domain.JenisInovasi, error) {

	script := `
		INSERT INTO tb_jenis_inovasi 
		(jenis) 
		VALUES (?)
	`

	result, err := tx.ExecContext(
		ctx,
		script,
		ji.Jenis,
	)
	if err != nil {
		return domain.JenisInovasi{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.JenisInovasi{}, err
	}

	ji.ID = int(id)

	return ji, nil
}

func (repository *JenisInovasiRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, ji domain.JenisInovasi) (domain.JenisInovasi, error) {

	// ================= UPDATE Isu =================
	query := `
		UPDATE tb_jenis_inovasi
		SET
			jenis = ?
		WHERE id = ?
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		ji.Jenis,
		ji.ID,
	)
	if err != nil {
		return domain.JenisInovasi{}, err
	}

	return ji, nil
}

func (repository *JenisInovasiRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_jenis_inovasi WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *JenisInovasiRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.JenisInovasi, error) {

	// ================= IKK =================
	query := `
		SELECT
			id,
			jenis
		FROM tb_jenis_inovasi
		WHERE id = ?
	`

	var result domain.JenisInovasi

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&result.ID,
		&result.Jenis,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.JenisInovasi{}, errors.New("Jenis Inovasi tidak ditemukan")
		}
		return domain.JenisInovasi{}, err
	}

	return result, nil
}

func (repository *JenisInovasiRepositoryImpl) FindByIds(ctx context.Context, tx *sql.Tx, ids []int) ([]domain.JenisInovasi, error) {

	if len(ids) == 0 {
		return []domain.JenisInovasi{}, nil
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
			jenis
		FROM tb_jenis_inovasi
		WHERE id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []domain.JenisInovasi{}

	for rows.Next() {
		var result domain.JenisInovasi

		err := rows.Scan(
			&result.ID,
			&result.Jenis,
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
		return []domain.JenisInovasi{}, nil
	}

	return results, nil
}

func (repository *JenisInovasiRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domain.JenisInovasi, error) {

	query := `
		SELECT id, 
			   jenis
		FROM tb_jenis_inovasi
	`

	args := make([]interface{}, 0)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	isuMap := make(map[int]*domain.JenisInovasi)
	isuIDs := make([]int, 0)

	for rows.Next() {
		var item domain.JenisInovasi

		err := rows.Scan(
			&item.ID,
			&item.Jenis,
		)
		if err != nil {
			return nil, err
		}

		copyItem := item
		isuMap[item.ID] = &copyItem
		isuIDs = append(isuIDs, item.ID)
	}

	if len(isuIDs) == 0 {
		return []domain.JenisInovasi{}, nil
	}

	result := make([]domain.JenisInovasi, 0, len(isuMap))
	for _, v := range isuMap {
		result = append(result, *v)
	}

	return result, nil
}