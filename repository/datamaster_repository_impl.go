package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/datamaster"
	"fmt"
)

type DataMasterRepositoryImpl struct {
}

func NewDataMasterRepositoryImpl() *DataMasterRepositoryImpl {
	return &DataMasterRepositoryImpl{}
}

func (repository *DataMasterRepositoryImpl) DataRBByTahun(ctx context.Context, tx *sql.Tx, tahun int) ([]datamaster.MasterRB, error) {

	query := `
		SELECT
			rb.id,
			rb.jenis_rb,
			rb.kegiatan_utama,
			rb.keterangan,
			rb.tahun_baseline,
			rb.tahun_next,

			ind.id AS indikator_id,
			ind.indikator,

			tar.id AS target_id,
			tar.tahun,
			tar.target,
			tar.realisasi,
			tar.satuan

		FROM datamaster_rb rb
		LEFT JOIN tb_indikator ind ON ind.id_rb = rb.id
		LEFT JOIN tb_target tar ON tar.indikator_id = ind.id
		WHERE rb.tahun_next = ? AND rb.is_active = 1
		ORDER BY rb.id, ind.id, tar.id;
	`

	rows, err := tx.QueryContext(ctx, query, tahun)
	if err != nil {
		return nil, fmt.Errorf("error querying rb: %v", err)
	}
	defer rows.Close()

	rbMap := make(map[int]*datamaster.MasterRB)
	indikatorIndexMap := make(map[int]map[string]int)

	for rows.Next() {

		var (
			rb datamaster.MasterRB

			indikatorRb datamaster.IndikatorRB
			targetRb    datamaster.TargetRB

			// FROM datamaster_rb (beberapa bisa NULL: keterangan)
			keterangan sql.NullString

			// FROM tb_indikator (bisa NULL SEMUA)
			indikatorId       sql.NullString
			indikatorNullable sql.NullString

			// FROM tb_target (bisa NULL SEMUA)
			targetId          sql.NullString
			targetYear        sql.NullInt64
			targetValue       sql.NullInt64
			realisasiNullable sql.NullFloat64
			satuanNullable    sql.NullString
		)

		err := rows.Scan(
			&rb.Id,
			&rb.JenisRB,
			&rb.KegiatanUtama,
			&keterangan,
			&rb.TahunBaseline,
			&rb.TahunNext,

			&indikatorId,
			&indikatorNullable,

			&targetId,
			&targetYear,
			&targetValue,
			&realisasiNullable,
			&satuanNullable,
		)

		if err != nil {
			return nil, err
		}

		// =============================
		// Handle fields datamaster_rb
		// =============================
		if keterangan.Valid {
			rb.Keterangan = keterangan.String
		} else {
			rb.Keterangan = ""
		}

		// Tambahkan RB jika belum ada
		if rbMap[rb.Id] == nil {
			rb.Indikator = []datamaster.IndikatorRB{}
			rbMap[rb.Id] = &rb
			indikatorIndexMap[rb.Id] = make(map[string]int)
		}

		master := rbMap[rb.Id]

		// =============================
		// Jika indikator NULL → RB tetap muncul
		// =============================
		if !indikatorId.Valid {
			continue
		}

		indikatorRb.IdRB = rb.Id
		indikatorRb.IdIndikator = indikatorId.String

		if indikatorNullable.Valid {
			indikatorRb.Indikator = indikatorNullable.String
		} else {
			indikatorRb.Indikator = ""
		}

		indikatorRb.TargetRB = []datamaster.TargetRB{}

		// cek apakah indikator sudah ada
		idxMap := indikatorIndexMap[rb.Id]
		indIndex, exists := idxMap[indikatorRb.IdIndikator]
		if !exists {
			master.Indikator = append(master.Indikator, indikatorRb)
			indIndex = len(master.Indikator) - 1
			idxMap[indikatorRb.IdIndikator] = indIndex
		}

		// =============================
		// Jika target NULL → lanjut
		// =============================
		if !targetId.Valid {
			continue
		}

		targetRb.IdTarget = targetId.String
		targetRb.IdIndikator = indikatorId.String

		// Tahun target
		if targetYear.Valid {
			year := int(targetYear.Int64)
			if year == master.TahunBaseline {
				targetRb.TahunBaseline = year
			}
			if year == master.TahunNext {
				targetRb.TahunNext = year
			}
		}

		// Nilai target
		if targetValue.Valid {
			if targetYear.Int64 == int64(master.TahunBaseline) {
				targetRb.TargetBaseline = int(targetValue.Int64)
			}
			if targetYear.Int64 == int64(master.TahunNext) {
				targetRb.TargetNext = int(targetValue.Int64)
			}
		}

		// Realisasi
		if realisasiNullable.Valid {
			targetRb.RealisasiBaseline = float32(realisasiNullable.Float64)
		}

		// Satuan
		if satuanNullable.Valid {
			targetRb.SatuanBaseline = satuanNullable.String
			targetRb.SatuanNext = satuanNullable.String
		}

		// Masukkan target ke indikator
		master.Indikator[indIndex].TargetRB =
			append(master.Indikator[indIndex].TargetRB, targetRb)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Convert map → slice
	result := make([]datamaster.MasterRB, 0, len(rbMap))
	for _, v := range rbMap {
		result = append(result, *v)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	return result, nil
}
