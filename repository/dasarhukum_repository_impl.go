package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"strings"
)

type DasarHukumRepositoryImpl struct {
}

func NewDasarHukumRepositoryImpl() *DasarHukumRepositoryImpl {
	return &DasarHukumRepositoryImpl{}
}

func (repository *DasarHukumRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, dasarHukum domain.DasarHukum) (domain.DasarHukum, error) {
	script := "INSERT INTO tb_dasar_hukum (id, rekin_id, kode_opd, urutan, peraturan_terkait, uraian) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, dasarHukum.Id, dasarHukum.RekinId, dasarHukum.KodeOpd, dasarHukum.Urutan, dasarHukum.PeraturanTerkait, dasarHukum.Uraian)
	if err != nil {
		return domain.DasarHukum{}, err
	}

	return dasarHukum, nil
}

func (repository *DasarHukumRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, dasarHukum domain.DasarHukum) (domain.DasarHukum, error) {
	script := "UPDATE tb_dasar_hukum SET urutan = ?, peraturan_terkait = ?, uraian = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, dasarHukum.Urutan, dasarHukum.PeraturanTerkait, dasarHukum.Uraian, dasarHukum.Id)
	if err != nil {
		return domain.DasarHukum{}, err
	}

	return dasarHukum, nil
}
func (repository *DasarHukumRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.DasarHukum, error) {
	script := "SELECT id, rekin_id, kode_opd, urutan, peraturan_terkait, uraian FROM tb_dasar_hukum WHERE 1=1"
	var params []interface{}

	if rekinId != "" {
		script += " AND rekin_id = ?"
		params = append(params, rekinId)
	}
	script += " ORDER BY urutan ASC"
	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return []domain.DasarHukum{}, err
	}
	defer rows.Close()

	var dasarHukum []domain.DasarHukum
	for rows.Next() {
		var dh domain.DasarHukum
		err = rows.Scan(&dh.Id, &dh.RekinId, &dh.KodeOpd, &dh.Urutan, &dh.PeraturanTerkait, &dh.Uraian)
		if err != nil {
			return []domain.DasarHukum{}, err
		}
		dasarHukum = append(dasarHukum, dh)
	}

	return dasarHukum, nil
}

func (repository *DasarHukumRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	script := "DELETE FROM tb_dasar_hukum WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *DasarHukumRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domain.DasarHukum, error) {
	script := "SELECT id, rekin_id, kode_opd, urutan, peraturan_terkait, uraian FROM tb_dasar_hukum WHERE id = ? ORDER BY urutan ASC"
	var dh domain.DasarHukum
	err := tx.QueryRowContext(ctx, script, id).Scan(&dh.Id, &dh.RekinId, &dh.KodeOpd, &dh.Urutan, &dh.PeraturanTerkait, &dh.Uraian)
	if err != nil {
		return domain.DasarHukum{}, err
	}
	return dh, nil
}

func (repository *DasarHukumRepositoryImpl) GetLastUrutan(ctx context.Context, tx *sql.Tx) (int, error) {
	script := "SELECT COALESCE(MAX(urutan), 0) FROM tb_dasar_hukum"
	var lastUrutan int
	err := tx.QueryRowContext(ctx, script).Scan(&lastUrutan)
	if err != nil {
		return 0, err
	}
	return lastUrutan, nil
}

func (r *DasarHukumRepositoryImpl) GetLastUrutanByRekinId(ctx context.Context, tx *sql.Tx, rekinId string) (int, error) {
	SQL := "SELECT COALESCE(MAX(urutan), 0) FROM tb_dasar_hukum WHERE rekin_id = ?"
	var lastUrutan int
	err := tx.QueryRowContext(ctx, SQL, rekinId).Scan(&lastUrutan)
	if err != nil {
		return 0, err
	}
	return lastUrutan, nil
}

func (repoistory *DasarHukumRepositoryImpl) FindByRekinIds(ctx context.Context, tx *sql.Tx, rekinIds []string) ([]domain.DasarHukum, error) {
	if len(rekinIds) == 0 {
		return []domain.DasarHukum{}, nil
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
           dh.id,
           dh.rekin_id,
           dh.kode_opd,
           dh.urutan,
           dh.peraturan_terkait,
           dh.uraian
        FROM tb_dasar_hukum dh
        WHERE dh.rekin_id IN (%s)
        `, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return []domain.DasarHukum{}, err
	}
	defer rows.Close()

	var dasarHukums []domain.DasarHukum
	for rows.Next() {
		var dasarHukum domain.DasarHukum
		err := rows.Scan(
			&dasarHukum.Id,
			&dasarHukum.RekinId,
			&dasarHukum.KodeOpd,
			&dasarHukum.Urutan,
			&dasarHukum.PeraturanTerkait,
			&dasarHukum.Uraian,
		)
		if err != nil {
			return []domain.DasarHukum{}, err
		}
		dasarHukums = append(dasarHukums, dasarHukum)
	}

	return dasarHukums, nil
}

func (repository *DasarHukumRepositoryImpl) BatchCreate(ctx context.Context, tx *sql.Tx, dasarHukums []domain.DasarHukum) error {
	if len(dasarHukums) == 0 {
		return nil
	}

	query := `
		INSERT INTO tb_dasar_hukum (
			id, rekin_id, urutan, peraturan_terkait, uraian, kode_opd
		) VALUES
	`

	var placeholders []string
	var values []any

	for _, item := range dasarHukums {
		placeholders = append(placeholders, "(?,?,?,?,?,?)")

		values = append(values,
			item.Id,
			item.RekinId,
			item.Urutan,
			item.PeraturanTerkait,
			item.Uraian,
			item.KodeOpd,
		)
	}

	query += strings.Join(placeholders, ",")

	_, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("gagal batch insert dasar hukum: %w", err)
	}

	return nil
}
