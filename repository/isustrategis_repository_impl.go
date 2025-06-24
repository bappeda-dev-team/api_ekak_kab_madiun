package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/isustrategis"
)

type CSFRepositoryImpl struct{}

func NewCSFRepositoryImpl() CSFRepository {
	return &CSFRepositoryImpl{}
}

func (repository *CSFRepositoryImpl) FindByTahun(ctx context.Context, tx *sql.Tx, tahun string) ([]isustrategis.CSF, error) {
	query := `
		SELECT id, pohon_id, pernyataan_kondisi_strategis,
		       alasan_kondisi_strategis, data_terukur,
		       kondisi_terukur, kondisi_wujud, tahun,
		       created_at, updated_at
		FROM tb_csf
		WHERE tahun = ?
	`

	rows, err := tx.QueryContext(ctx, query, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []isustrategis.CSF

	for rows.Next() {
		var csf isustrategis.CSF
		err := rows.Scan(
			&csf.ID,
			&csf.PohonID,
			&csf.PernyataanKondisiStrategis,
			&csf.AlasanKondisiStrategis,
			&csf.DataTerukur,
			&csf.KondisiTerukur,
			&csf.KondisiWujud,
			&csf.Tahun,
			&csf.CreatedAt,
			&csf.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, csf)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
