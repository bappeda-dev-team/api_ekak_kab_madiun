package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"strings"
)

type GambaranUmumRepositoryImpl struct {
}

func NewGambaranUmumRepositoryImpl() *GambaranUmumRepositoryImpl {
	return &GambaranUmumRepositoryImpl{}
}

func (repository *GambaranUmumRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, gambaranUmum domain.GambaranUmum) (domain.GambaranUmum, error) {
	query := "INSERT INTO tb_gambaran_umum (id, rekin_id, kode_opd, urutan, gambaran_umum) VALUES (?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, gambaranUmum.Id, gambaranUmum.RekinId, gambaranUmum.KodeOpd, gambaranUmum.Urutan, gambaranUmum.GambaranUmum)
	if err != nil {
		return domain.GambaranUmum{}, err
	}
	return gambaranUmum, nil
}

func (repository *GambaranUmumRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, gambaranUmum domain.GambaranUmum) (domain.GambaranUmum, error) {
	query := "UPDATE tb_gambaran_umum SET gambaran_umum = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, gambaranUmum.GambaranUmum, gambaranUmum.Id)
	if err != nil {
		return domain.GambaranUmum{}, err
	}
	return gambaranUmum, nil
}

func (repository *GambaranUmumRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	query := "DELETE FROM tb_gambaran_umum WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *GambaranUmumRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domain.GambaranUmum, error) {
	query := "SELECT id, rekin_id, kode_opd, urutan, gambaran_umum FROM tb_gambaran_umum WHERE id = ? ORDER BY urutan ASC"
	row := tx.QueryRowContext(ctx, query, id)
	var gambaranUmum domain.GambaranUmum
	err := row.Scan(&gambaranUmum.Id, &gambaranUmum.RekinId, &gambaranUmum.KodeOpd, &gambaranUmum.Urutan, &gambaranUmum.GambaranUmum)
	if err != nil {
		return domain.GambaranUmum{}, err
	}
	return gambaranUmum, nil
}

func (repository *GambaranUmumRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.GambaranUmum, error) {
	query := "SELECT id, rekin_id, kode_opd, urutan, gambaran_umum FROM tb_gambaran_umum WHERE rekin_id = ? ORDER BY urutan ASC"
	rows, err := tx.QueryContext(ctx, query, rekinId)
	if err != nil {
		return []domain.GambaranUmum{}, err
	}
	defer rows.Close()

	var gambaranUmumList []domain.GambaranUmum
	for rows.Next() {
		var gambaranUmum domain.GambaranUmum
		err := rows.Scan(&gambaranUmum.Id, &gambaranUmum.RekinId, &gambaranUmum.KodeOpd, &gambaranUmum.Urutan, &gambaranUmum.GambaranUmum)
		if err != nil {
			return []domain.GambaranUmum{}, err
		}

		gambaranUmumList = append(gambaranUmumList, gambaranUmum)
	}

	err = rows.Err()
	if err != nil {
		return []domain.GambaranUmum{}, err
	}

	return gambaranUmumList, nil
}

func (repository *GambaranUmumRepositoryImpl) GetLastUrutanByRekinId(ctx context.Context, tx *sql.Tx, rekinId string) (int, error) {
	SQL := "SELECT COALESCE(MAX(urutan), 0) FROM tb_gambaran_umum WHERE rekin_id = ?"
	var lastUrutan int
	err := tx.QueryRowContext(ctx, SQL, rekinId).Scan(&lastUrutan)
	if err != nil {
		return 0, err
	}
	return lastUrutan, nil
}

func (repository *GambaranUmumRepositoryImpl) FindByRekinIds(
	ctx context.Context,
	tx *sql.Tx,
	rekinIds []string,
) ([]domain.GambaranUmum, error) {

	if len(rekinIds) == 0 {
		return []domain.GambaranUmum{}, nil
	}

	placeholders := make([]string, len(rekinIds))
	args := make([]any, len(rekinIds))

	for i, id := range rekinIds {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT
			id, rekin_id, kode_opd, urutan, gambaran_umum
		FROM tb_gambaran_umum
		WHERE rekin_id IN (%s)
		ORDER BY urutan ASC
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.GambaranUmum

	for rows.Next() {
		var item domain.GambaranUmum
		err := rows.Scan(
			&item.Id,
			&item.RekinId,
			&item.KodeOpd,
			&item.Urutan,
			&item.GambaranUmum,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (repository *GambaranUmumRepositoryImpl) BatchCreate(
	ctx context.Context,
	tx *sql.Tx,
	items []domain.GambaranUmum,
) error {

	if len(items) == 0 {
		return nil
	}

	query := `
		INSERT INTO tb_gambaran_umum (
			id, rekin_id, kode_opd, urutan, gambaran_umum
		) VALUES
	`

	var placeholders []string
	var values []any

	for _, item := range items {
		placeholders = append(placeholders, "(?,?,?,?,?)")

		values = append(values,
			item.Id,
			item.RekinId,
			item.KodeOpd,
			item.Urutan,
			item.GambaranUmum,
		)
	}

	query += strings.Join(placeholders, ",")

	_, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("gagal batch insert gambaran umum: %w", err)
	}

	return nil
}
