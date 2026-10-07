package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"strings"
)

type OutcomeMatrixRepositoryImpl struct {
}

func NewOutcomeMatrixRepositoryImpl() *OutcomeMatrixRepositoryImpl {
	return &OutcomeMatrixRepositoryImpl{}
}

func (r *OutcomeMatrixRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, data domain.OutcomeMatrix) (domain.OutcomeMatrix, error) {
	script := `
		INSERT INTO tb_outcome_matrix (kode_opd, kode, jenis, outcome)
		VALUES (?, ?, ?, ?)
	`
	result, err := tx.ExecContext(ctx, script, data.KodeOpd, data.Kode, data.Jenis, data.Outcome)
	if err != nil {
		return domain.OutcomeMatrix{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return domain.OutcomeMatrix{}, err
	}
	return r.FindById(ctx, tx, int(id))
}

func (r *OutcomeMatrixRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, data domain.OutcomeMatrix) (domain.OutcomeMatrix, error) {
	script := `
		UPDATE tb_outcome_matrix
		SET kode_opd = ?, kode = ?, jenis = ?, outcome = ?
		WHERE id = ?
	`
	_, err := tx.ExecContext(ctx, script, data.KodeOpd, data.Kode, data.Jenis, data.Outcome, data.Id)
	if err != nil {
		return domain.OutcomeMatrix{}, err
	}
	return r.FindById(ctx, tx, data.Id)
}

func (r *OutcomeMatrixRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	_, err := tx.ExecContext(ctx, "DELETE FROM tb_outcome_matrix WHERE id = ?", id)
	return err
}

func (r *OutcomeMatrixRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.OutcomeMatrix, error) {
	script := `
		SELECT id, kode_opd, kode, COALESCE(jenis, ''), outcome, created_at, updated_at
		FROM tb_outcome_matrix
		WHERE id = ?
	`
	var data domain.OutcomeMatrix
	err := tx.QueryRowContext(ctx, script, id).Scan(
		&data.Id, &data.KodeOpd, &data.Kode, &data.Jenis, &data.Outcome, &data.CreatedAt, &data.UpdatedAt,
	)
	if err != nil {
		return domain.OutcomeMatrix{}, err
	}
	return data, nil
}

func (r *OutcomeMatrixRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kode, kodeOpd, jenis string) ([]domain.OutcomeMatrix, error) {
	script := `
		SELECT id, kode_opd, kode, COALESCE(jenis, ''), outcome, created_at, updated_at
		FROM tb_outcome_matrix
		WHERE (? = '' OR kode = ?)
		  AND (? = '' OR kode_opd = ?)
		  AND (? = '' OR jenis = ?)
		ORDER BY id ASC
	`
	rows, err := tx.QueryContext(ctx, script, kode, kode, kodeOpd, kodeOpd, jenis, jenis)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOutcomeMatrixRows(rows)
}

func (r *OutcomeMatrixRepositoryImpl) FindByKodeAndKodeOpd(ctx context.Context, tx *sql.Tx, kodes []string, kodeOpd string) ([]domain.OutcomeMatrix, error) {
	if kodeOpd == "" || len(kodes) == 0 {
		return []domain.OutcomeMatrix{}, nil
	}
	placeholders := strings.Repeat("?,", len(kodes))
	placeholders = placeholders[:len(placeholders)-1]
	args := make([]interface{}, 0, len(kodes)+1)
	for _, k := range kodes {
		args = append(args, k)
	}
	args = append(args, kodeOpd)
	script := fmt.Sprintf(`
		SELECT id, kode_opd, kode, COALESCE(jenis, ''), outcome, created_at, updated_at
		FROM tb_outcome_matrix
		WHERE kode IN (%s)
		  AND kode_opd = ?
		ORDER BY id ASC
	`, placeholders)
	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOutcomeMatrixRows(rows)
}

func scanOutcomeMatrixRows(rows *sql.Rows) ([]domain.OutcomeMatrix, error) {
	result := make([]domain.OutcomeMatrix, 0)
	for rows.Next() {
		var data domain.OutcomeMatrix
		if err := rows.Scan(
			&data.Id, &data.KodeOpd, &data.Kode, &data.Jenis, &data.Outcome, &data.CreatedAt, &data.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, data)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
