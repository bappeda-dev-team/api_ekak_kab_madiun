package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/isustrategis"
)

type CSFRepositoryImpl struct{}

func NewCSFRepositoryImpl() CSFRepository {
	return &CSFRepositoryImpl{}
}

func (repository *CSFRepositoryImpl) FindByTahun(ctx context.Context, tx *sql.Tx, tahun string) ([]isustrategis.CSFPokin, error) {
	query := `
	SELECT
		tb_csf.id,
		tb_csf.pohon_id,
		tb_csf.pernyataan_kondisi_strategis,
		tb_csf.alasan_kondisi_strategis,
		tb_csf.data_terukur,
		tb_csf.kondisi_terukur,
		tb_csf.kondisi_wujud,
		tb_csf.tahun,
		tb_pohon_kinerja.jenis_pohon,
		tb_pohon_kinerja.level_pohon,
		tb_pohon_kinerja.nama_pohon,
		tb_pohon_kinerja.keterangan,
        tb_pohon_kinerja.is_active
	FROM
		tb_csf
	JOIN tb_pohon_kinerja ON tb_csf.pohon_id = tb_pohon_kinerja.id
	WHERE
		tb_csf.tahun = ?
	`

	rows, err := tx.QueryContext(ctx, query, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []isustrategis.CSFPokin

	for rows.Next() {
		var csf isustrategis.CSFPokin
		err := rows.Scan(
			&csf.ID,
			&csf.PohonID,
			&csf.PernyataanKondisiStrategis,
			&csf.AlasanKondisiStrategis,
			&csf.DataTerukur,
			&csf.KondisiTerukur,
			&csf.KondisiWujud,
			&csf.Tahun,
			&csf.JenisPohon,
			&csf.LevelPohon,
			&csf.Strategi,
			&csf.Keterangan,
			&csf.IsActive,
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

func (r *CSFRepositoryImpl) CreateCsf(ctx context.Context, tx *sql.Tx, csf domain.CSF) error {
	query := `
		INSERT INTO tb_csf 
			(pohon_id, pernyataan_kondisi_strategis, alasan_kondisi_strategis, data_terukur, kondisi_terukur, kondisi_wujud, tahun)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := tx.ExecContext(ctx, query,
		csf.PohonID,
		csf.PernyataanKondisiStrategis,
		csf.AlasanKondisiStrategis,
		csf.DataTerukur,
		csf.KondisiTerukur,
		csf.KondisiWujud,
		csf.Tahun,
	)
	if err != nil {
		return err
	}
	return nil
}
