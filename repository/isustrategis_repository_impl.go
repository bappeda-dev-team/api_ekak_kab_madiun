package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/isustrategis"
	"fmt"
	"log"
	"sort"
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
		tb_pohon_kinerja.is_active,
		i.id as indikator_id,
		i.indikator as nama_indikator,
		t.id as target_id,
		t.target as target_value,
		t.satuan as target_satuan
	FROM
		tb_csf
	JOIN tb_pohon_kinerja ON tb_csf.pohon_id = tb_pohon_kinerja.id
	LEFT JOIN tb_indikator i ON tb_pohon_kinerja.id = i.pokin_id
	LEFT JOIN tb_target t ON i.id = t.indikator_id
	WHERE
		tb_csf.tahun = ?
	ORDER BY
		tb_csf.id
	`

	rows, err := tx.QueryContext(ctx, query, tahun)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	csfMap := make(map[int]*isustrategis.CSFPokin)
	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var (
			csfID, pohonID               int
			pernyataan, alasan, data     string
			kondisiTerukur, kondisiWujud string
			tahunInt                     int
			jenisPohon, namaPohon, ket   string
			levelPohon                   int
			isActive                     bool
			indikatorID                  sql.NullString
			namaIndikator                sql.NullString
			targetID                     sql.NullString
			targetValue                  sql.NullString
			targetSatuan                 sql.NullString
		)

		if err := rows.Scan(
			&csfID, &pohonID, &pernyataan, &alasan, &data,
			&kondisiTerukur, &kondisiWujud, &tahunInt,
			&jenisPohon, &levelPohon, &namaPohon, &ket, &isActive,
			&indikatorID, &namaIndikator,
			&targetID, &targetValue, &targetSatuan,
		); err != nil {
			return nil, err
		}

		csf, ok := csfMap[csfID]
		if !ok {
			csf = &isustrategis.CSFPokin{
				ID:                         csfID,
				PohonID:                    pohonID,
				PernyataanKondisiStrategis: pernyataan,
				AlasanKondisiStrategis:     alasan,
				DataTerukur:                data,
				KondisiTerukur:             kondisiTerukur,
				KondisiWujud:               kondisiWujud,
				Tahun:                      tahunInt,
				JenisPohon:                 jenisPohon,
				LevelPohon:                 levelPohon,
				Strategi:                   namaPohon,
				Keterangan:                 ket,
				IsActive:                   isActive,
			}
			csfMap[csfID] = csf
		}

		// Proses indikator jika ada
		if indikatorID.Valid && namaIndikator.Valid {
			indID := indikatorID.String
			indikator, exists := indikatorMap[indID]
			if !exists {
				indikator = &domain.Indikator{
					Id:        indID,
					PokinId:   fmt.Sprint(pohonID),
					Indikator: namaIndikator.String,
					Tahun:     fmt.Sprint(tahunInt),
					Target:    []domain.Target{},
				}
				indikatorMap[indID] = indikator
			}

			// Tambahkan target jika ada
			if targetID.Valid && targetValue.Valid && targetSatuan.Valid {
				indikator.Target = append(indikator.Target, domain.Target{
					Id:          targetID.String,
					IndikatorId: indID,
					Target:      targetValue.String,
					Satuan:      targetSatuan.String,
					Tahun:       fmt.Sprint(tahunInt),
				})
			}

			// Tambahkan ke CSF jika indikator belum pernah dimasukkan
			found := false
			for _, ind := range csf.Indikator {
				if ind.Id == indID {
					found = true
					break
				}
			}
			if !found {
				csf.Indikator = append(csf.Indikator, *indikator)
			}
		}
	}

	log.Print("[LOG] Record CSF ditemukan")
	var result []isustrategis.CSFPokin
	var keys []int
	for id := range csfMap {
		keys = append(keys, id)
	}
	sort.Ints(keys)
	for _, id := range keys {
		result = append(result, *csfMap[id])
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

func (r *CSFRepositoryImpl) UpdateCSFByPohonID(ctx context.Context, tx *sql.Tx, csf domain.CSF) (domain.CSF, error) {
	query := `
	UPDATE tb_csf
	SET
		pernyataan_kondisi_strategis = ?,
		alasan_kondisi_strategis = ?,
		data_terukur = ?,
		kondisi_terukur = ?,
		kondisi_wujud = ?
	WHERE pohon_id = ?
`
	if csf.PohonID == 0 {
		return domain.CSF{}, fmt.Errorf("[ERROR] POHON ID TIDAK BOLEH 0")
	}
	result, err := tx.ExecContext(ctx, query,
		csf.PernyataanKondisiStrategis,
		csf.AlasanKondisiStrategis,
		csf.DataTerukur,
		csf.KondisiTerukur,
		csf.KondisiWujud,
		csf.PohonID,
	)
	if err != nil {
		return domain.CSF{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.CSF{}, err
	}
	log.Printf("[LOG] ROW CSF UPDATED: %d", rowsAffected)

	return csf, nil
}

func (repository *CSFRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, csfId int) (isustrategis.CSFPokin, error) {
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
		tb_pohon_kinerja.is_active,
		i.id as indikator_id,
		i.indikator as nama_indikator,
		t.id as target_id,
		t.target as target_value,
		t.satuan as target_satuan
	FROM
		tb_csf
	JOIN tb_pohon_kinerja ON tb_csf.pohon_id = tb_pohon_kinerja.id
	LEFT JOIN tb_indikator i ON tb_pohon_kinerja.id = i.pokin_id
	LEFT JOIN tb_target t ON i.id = t.indikator_id
	WHERE
		tb_csf.id = ?
	ORDER BY i.id, t.id
	`

	rows, err := tx.QueryContext(ctx, query, csfId)
	if err != nil {
		return isustrategis.CSFPokin{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var csf *isustrategis.CSFPokin
	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var (
			csfID, pohonID               int
			pernyataan, alasan, data     string
			kondisiTerukur, kondisiWujud string
			tahunInt                     int
			jenisPohon, namaPohon, ket   string
			levelPohon                   int
			isActive                     bool
			indikatorID                  sql.NullString
			namaIndikator                sql.NullString
			targetID                     sql.NullString
			targetValue                  sql.NullString
			targetSatuan                 sql.NullString
		)

		if err := rows.Scan(
			&csfID, &pohonID, &pernyataan, &alasan, &data,
			&kondisiTerukur, &kondisiWujud, &tahunInt,
			&jenisPohon, &levelPohon, &namaPohon, &ket, &isActive,
			&indikatorID, &namaIndikator,
			&targetID, &targetValue, &targetSatuan,
		); err != nil {
			return isustrategis.CSFPokin{}, err
		}

		if csf == nil {
			csf = &isustrategis.CSFPokin{
				ID:                         csfID,
				PohonID:                    pohonID,
				PernyataanKondisiStrategis: pernyataan,
				AlasanKondisiStrategis:     alasan,
				DataTerukur:                data,
				KondisiTerukur:             kondisiTerukur,
				KondisiWujud:               kondisiWujud,
				Tahun:                      tahunInt,
				JenisPohon:                 jenisPohon,
				LevelPohon:                 levelPohon,
				Strategi:                   namaPohon,
				Keterangan:                 ket,
				IsActive:                   isActive,
				Indikator:                  []domain.Indikator{},
			}
		}

		if indikatorID.Valid && namaIndikator.Valid {
			indID := indikatorID.String
			indikator, exists := indikatorMap[indID]
			if !exists {
				indikator = &domain.Indikator{
					Id:        indID,
					PokinId:   fmt.Sprint(pohonID),
					Indikator: namaIndikator.String,
					Tahun:     fmt.Sprint(tahunInt),
					Target:    []domain.Target{},
				}
				indikatorMap[indID] = indikator
			}

			if targetID.Valid && targetValue.Valid && targetSatuan.Valid {
				indikator.Target = append(indikator.Target, domain.Target{
					Id:          targetID.String,
					IndikatorId: indID,
					Target:      targetValue.String,
					Satuan:      targetSatuan.String,
					Tahun:       fmt.Sprint(tahunInt),
				})
			}
		}
	}

	// Akumulasi indikator yang ditemukan
	for _, indikator := range indikatorMap {
		if csf != nil && indikator != nil { // Add check for csf != nil
			csf.Indikator = append(csf.Indikator, *indikator)
		}
	}

	if csf == nil {
		return isustrategis.CSFPokin{}, sql.ErrNoRows
	}

	return *csf, nil
}
