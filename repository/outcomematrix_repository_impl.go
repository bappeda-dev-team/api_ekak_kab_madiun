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
		INSERT INTO tb_outcome_matrix (kode_subkegiatan, kode, outcome)
		VALUES (?, ?, ?)
	`
	result, err := tx.ExecContext(ctx, script, data.KodeSubkegiatan, data.Kode, data.Outcome)
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
		SET kode_subkegiatan = ?, kode = ?, outcome = ?
		WHERE id = ?
	`
	_, err := tx.ExecContext(ctx, script, data.KodeSubkegiatan, data.Kode, data.Outcome, data.Id)
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
		SELECT id, kode_subkegiatan, kode, outcome, created_at, updated_at
		FROM tb_outcome_matrix
		WHERE id = ?
	`
	var data domain.OutcomeMatrix
	err := tx.QueryRowContext(ctx, script, id).Scan(
		&data.Id, &data.KodeSubkegiatan, &data.Kode, &data.Outcome, &data.CreatedAt, &data.UpdatedAt,
	)
	if err != nil {
		return domain.OutcomeMatrix{}, err
	}
	return data, nil
}

func (r *OutcomeMatrixRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kode, kodeSubkegiatan string) ([]domain.OutcomeMatrix, error) {
	script := `
		SELECT id, kode_subkegiatan, kode, outcome, created_at, updated_at
		FROM tb_outcome_matrix
		WHERE (? = '' OR kode = ?)
		  AND (? = '' OR kode_subkegiatan = ?)
		ORDER BY id ASC
	`
	rows, err := tx.QueryContext(ctx, script, kode, kode, kodeSubkegiatan, kodeSubkegiatan)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOutcomeMatrixRows(rows)
}

func (r *OutcomeMatrixRepositoryImpl) FindByKodes(ctx context.Context, tx *sql.Tx, kodes []string, kodeSubkegiatans []string) ([]domain.OutcomeMatrix, error) {
	if len(kodes) == 0 && len(kodeSubkegiatans) == 0 {
		return []domain.OutcomeMatrix{}, nil
	}
	var conditions []string
	var args []interface{}
	if len(kodes) > 0 {
		placeholders := strings.Repeat("?,", len(kodes))
		placeholders = placeholders[:len(placeholders)-1]
		conditions = append(conditions, fmt.Sprintf("kode IN (%s)", placeholders))
		for _, k := range kodes {
			args = append(args, k)
		}
	}
	if len(kodeSubkegiatans) > 0 {
		placeholders := strings.Repeat("?,", len(kodeSubkegiatans))
		placeholders = placeholders[:len(placeholders)-1]
		conditions = append(conditions, fmt.Sprintf("kode_subkegiatan IN (%s)", placeholders))
		for _, k := range kodeSubkegiatans {
			args = append(args, k)
		}
	}
	script := `
		SELECT id, kode_subkegiatan, kode, outcome, created_at, updated_at
		FROM tb_outcome_matrix
		WHERE ` + strings.Join(conditions, " OR ") + `
		ORDER BY id ASC
	`
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
			&data.Id, &data.KodeSubkegiatan, &data.Kode, &data.Outcome, &data.CreatedAt, &data.UpdatedAt,
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
