package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"strings"
)

type PermasalahanRekinRepositoryImpl struct {
}

func NewPermasalahanRekinRepositoryImpl() *PermasalahanRekinRepositoryImpl {
	return &PermasalahanRekinRepositoryImpl{}
}

func (repository *PermasalahanRekinRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, permasalahan domain.PermasalahanRekin) (domain.PermasalahanRekin, error) {
	script := "INSERT INTO tb_permasalahan (id, rekin_id, permasalahan, penyebab_internal, penyebab_eksternal, jenis_permasalahan) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, permasalahan.Id, permasalahan.RekinId, permasalahan.Permasalahan, permasalahan.PenyebabInternal, permasalahan.PenyebabEksternal, permasalahan.JenisPermasalahan)
	if err != nil {
		return domain.PermasalahanRekin{}, err
	}
	return permasalahan, nil
}

func (repository *PermasalahanRekinRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, permasalahan domain.PermasalahanRekin) (domain.PermasalahanRekin, error) {
	script := "UPDATE tb_permasalahan SET permasalahan = ?, penyebab_internal = ?, penyebab_eksternal = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, permasalahan.Permasalahan, permasalahan.PenyebabInternal, permasalahan.PenyebabEksternal, permasalahan.Id)
	if err != nil {
		return domain.PermasalahanRekin{}, err
	}
	return permasalahan, nil
}

func (repository *PermasalahanRekinRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_permasalahan WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *PermasalahanRekinRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, rekinId *string) ([]domain.PermasalahanRekin, error) {
	script := "SELECT id, rekin_id, permasalahan, penyebab_internal, penyebab_eksternal, jenis_permasalahan, created_at FROM tb_permasalahan WHERE 1=1"

	var args []interface{}

	if rekinId != nil {
		script += " AND rekin_id = ?"
		args = append(args, *rekinId)
	}

	script += " order by created_at asc"

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return []domain.PermasalahanRekin{}, err
	}
	defer rows.Close()

	var permasalahanList []domain.PermasalahanRekin
	for rows.Next() {
		var permasalahan domain.PermasalahanRekin
		err := rows.Scan(&permasalahan.Id, &permasalahan.RekinId, &permasalahan.Permasalahan, &permasalahan.PenyebabInternal, &permasalahan.PenyebabEksternal, &permasalahan.JenisPermasalahan, &permasalahan.CreatedAt)
		if err != nil {
			return []domain.PermasalahanRekin{}, err
		}
		permasalahanList = append(permasalahanList, permasalahan)
	}
	return permasalahanList, nil
}

func (repository *PermasalahanRekinRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.PermasalahanRekin, error) {
	script := "SELECT id, rekin_id, permasalahan, penyebab_internal, penyebab_eksternal, jenis_permasalahan FROM tb_permasalahan WHERE id = ?"
	var permasalahan domain.PermasalahanRekin
	err := tx.QueryRowContext(ctx, script, id).Scan(&permasalahan.Id, &permasalahan.RekinId, &permasalahan.Permasalahan, &permasalahan.PenyebabInternal, &permasalahan.PenyebabEksternal, &permasalahan.JenisPermasalahan)
	if err != nil {
		return domain.PermasalahanRekin{}, err
	}
	return permasalahan, nil
}

func (repository *PermasalahanRekinRepositoryImpl) FindByRekinIds(ctx context.Context, tx *sql.Tx, rekinIds []string) ([]domain.PermasalahanRekin, error) {
	if len(rekinIds) == 0 {
		return []domain.PermasalahanRekin{}, nil
	}
	// Buat placeholder untuk IN clause
	placeholders := make([]string, len(rekinIds))
	args := make([]any, len(rekinIds))
	for i, id := range rekinIds {
		placeholders[i] = "?"
		args[i] = id
	}

	script := fmt.Sprintf(`
        SELECT
           pr.id,
           pr.rekin_id,
           pr.permasalahan,
           pr.penyebab_internal,
           pr.penyebab_eksternal,
           pr.jenis_permasalahan
        FROM tb_permasalahan pr
        WHERE pr.rekin_id IN (%s)
        `, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return []domain.PermasalahanRekin{}, err
	}
	defer rows.Close()

	var permasalahans []domain.PermasalahanRekin
	for rows.Next() {
		var permasalahan domain.PermasalahanRekin
		err := rows.Scan(
			&permasalahan.Id,
			&permasalahan.RekinId,
			&permasalahan.Permasalahan,
			&permasalahan.PenyebabInternal,
			&permasalahan.PenyebabEksternal,
			&permasalahan.JenisPermasalahan,
		)
		if err != nil {
			return []domain.PermasalahanRekin{}, err
		}
		permasalahans = append(permasalahans, permasalahan)
	}

	return permasalahans, nil
}

func (repository *PermasalahanRekinRepositoryImpl) BatchCreate(ctx context.Context, tx *sql.Tx, permasalahans []domain.PermasalahanRekin) error {
	if len(permasalahans) == 0 {
		return nil
	}

	query := `
		INSERT INTO tb_permasalahan (
			rekin_id, permasalahan, penyebab_internal, penyebab_eksternal, jenis_permasalahan
		) VALUES
	`

	var placeholders []string
	var values []any

	for _, item := range permasalahans {
		placeholders = append(placeholders, "(?,?,?,?,?)")

		values = append(values,
			item.RekinId,
			item.Permasalahan,
			item.PenyebabInternal,
			item.PenyebabEksternal,
			item.JenisPermasalahan,
		)
	}

	query += strings.Join(placeholders, ",")

	_, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("gagal batch insert permasalahan: %w", err)
	}

	return nil
}
