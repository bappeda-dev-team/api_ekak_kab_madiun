package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type TujuanPemdaRepositoryImpl struct {
}

func NewTujuanPemdaRepositoryImpl() *TujuanPemdaRepositoryImpl {
	return &TujuanPemdaRepositoryImpl{}
}

func (repository *TujuanPemdaRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, tujuanPemda domain.TujuanPemda) (domain.TujuanPemda, error) {
	query := "INSERT INTO tb_tujuan_pemda(id, tujuan_pemda, tematik_id, periode_id, tahun_awal_periode, tahun_akhir_periode, jenis_periode, id_visi, id_misi) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, tujuanPemda.Id, tujuanPemda.TujuanPemda, tujuanPemda.TematikId, tujuanPemda.PeriodeId, tujuanPemda.TahunAwalPeriode, tujuanPemda.TahunAkhirPeriode, tujuanPemda.JenisPeriode, tujuanPemda.IdVisi, tujuanPemda.IdMisi)
	if err != nil {
		return tujuanPemda, err
	}
	return tujuanPemda, nil
}

func (r *TujuanPemdaRepositoryImpl) CreateIndikator(ctx context.Context, tx *sql.Tx, ind domain.IndikatorPemda) (domain.IndikatorPemda, error) {
	query := `INSERT INTO tb_indikator_matrix_pemda
		(kode_indikator, tujuan_pemda_id, sasaran_pemda_id, indikator,
		 rumus_perhitungan, sumber_data, definisi_operasional, jenis)
		VALUES (?, ?, 0, ?, ?, ?, ?, ?)`
	res, err := tx.ExecContext(ctx, query,
		ind.KodeIndikator,
		ind.TujuanPemdaId,
		ind.Indikator,
		ind.RumusPerhitungan,
		ind.SumberData,
		ind.DefinisiOperasional,
		ind.Jenis,
	)
	if err != nil {
		return ind, err
	}
	dbId, _ := res.LastInsertId()
	ind.Id = int(dbId)
	return ind, nil
}

func (r *TujuanPemdaRepositoryImpl) CreateTarget(ctx context.Context, tx *sql.Tx, t domain.TargetPemda) (domain.TargetPemda, error) {
	query := `INSERT INTO tb_target_pemda
		(kode_indikator, target, satuan, tahun, jenis)
		VALUES (?, ?, ?, ?, ?)`
	res, err := tx.ExecContext(ctx, query,
		t.KodeIndikator, t.Target, t.Satuan, t.Tahun, t.Jenis)
	if err != nil {
		return t, err
	}
	dbId, _ := res.LastInsertId()
	t.Id = int(dbId)
	return t, nil
}

//logic lama
// func (repository *TujuanPemdaRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, tujuanPemda domain.TujuanPemda) (domain.TujuanPemda, error) {
// 	// Update tujuan pemda
// 	query := "UPDATE tb_tujuan_pemda SET tujuan_pemda = ?, tematik_id = ?, periode_id = ?, tahun_awal_periode = ?, tahun_akhir_periode = ?, jenis_periode = ?, id_visi = ?, id_misi = ? WHERE id = ?"
// 	_, err := tx.ExecContext(ctx, query, tujuanPemda.TujuanPemda, tujuanPemda.TematikId, tujuanPemda.PeriodeId, tujuanPemda.TahunAwalPeriode, tujuanPemda.TahunAkhirPeriode, tujuanPemda.JenisPeriode, tujuanPemda.IdVisi, tujuanPemda.IdMisi, tujuanPemda.Id)
// 	if err != nil {
// 		return tujuanPemda, err
// 	}

// 	// Hapus semua indikator lama beserta targetnya
// 	scriptDeleteOldIndicators := "DELETE FROM tb_indikator WHERE tujuan_pemda_id = ?"
// 	_, err = tx.ExecContext(ctx, scriptDeleteOldIndicators, tujuanPemda.Id)
// 	if err != nil {
// 		return tujuanPemda, err
// 	}

// 	// Insert indikator baru
// 	for _, indikator := range tujuanPemda.Indikator {
// 		scriptInsertIndikator := `
//             INSERT INTO tb_indikator
//                 (id, tujuan_pemda_id, indikator, rumus_perhitungan, sumber_data)
//             VALUES
//                 (?, ?, ?, ?, ?)`

// 		_, err := tx.ExecContext(ctx, scriptInsertIndikator,
// 			indikator.Id,
// 			tujuanPemda.Id,
// 			indikator.Indikator,
// 			indikator.RumusPerhitungan,
// 			indikator.SumberData)
// 		if err != nil {
// 			return tujuanPemda, err
// 		}

// 		// Hapus semua target lama untuk indikator ini
// 		scriptDeleteOldTargets := "DELETE FROM tb_target WHERE indikator_id = ?"
// 		_, err = tx.ExecContext(ctx, scriptDeleteOldTargets, indikator.Id)
// 		if err != nil {
// 			return tujuanPemda, err
// 		}

// 		// Insert target baru
// 		for _, target := range indikator.Target {
// 			scriptInsertTarget := `
//                 INSERT INTO tb_target
//                     (id, indikator_id, target, satuan, tahun)
//                 VALUES
//                     (?, ?, ?, ?, ?)`

// 			_, err := tx.ExecContext(ctx, scriptInsertTarget,
// 				target.Id,
// 				indikator.Id,
// 				target.Target,
// 				target.Satuan,
// 				target.Tahun)
// 			if err != nil {
// 				return tujuanPemda, err
// 			}
// 		}
// 	}

// 	return tujuanPemda, nil
// }

func (r *TujuanPemdaRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, tp domain.TujuanPemda) (domain.TujuanPemda, error) {
	_, err := tx.ExecContext(ctx,
		`UPDATE tb_tujuan_pemda
		 SET tujuan_pemda=?, tematik_id=?, periode_id=?,
		     id_visi=?, id_misi=?
		 WHERE id=?`,
		tp.TujuanPemda, tp.TematikId, tp.PeriodeId,
		tp.IdVisi, tp.IdMisi, tp.Id)
	if err != nil {
		return tp, err
	}
	keepIndIds := make([]int, 0, len(tp.IndikatorPemda))
	for _, ind := range tp.IndikatorPemda {
		kodeInd := ind.KodeIndikator
		if ind.Id > 0 {
			// UPDATE indikator existing — id tetap
			_, err = tx.ExecContext(ctx, `
				UPDATE tb_indikator_matrix_pemda
				SET indikator=?, rumus_perhitungan=?, sumber_data=?,
				    definisi_operasional=?, jenis=?
				WHERE id=? AND tujuan_pemda_id=?`,
				ind.Indikator, ind.RumusPerhitungan, ind.SumberData,
				ind.DefinisiOperasional, ind.Jenis, ind.Id, tp.Id)
			if err != nil {
				return tp, err
			}
			keepIndIds = append(keepIndIds, ind.Id)
		} else {
			// INSERT indikator baru
			res, err := tx.ExecContext(ctx, `
				INSERT INTO tb_indikator_matrix_pemda
					(kode_indikator, tujuan_pemda_id, sasaran_pemda_id, indikator,
					 rumus_perhitungan, sumber_data, definisi_operasional, jenis)
				VALUES (?, ?, 0, ?, ?, ?, ?, ?)`,
				kodeInd, tp.Id,
				ind.Indikator, ind.RumusPerhitungan,
				ind.SumberData, ind.DefinisiOperasional, ind.Jenis)
			if err != nil {
				return tp, err
			}
			newId, _ := res.LastInsertId()
			ind.Id = int(newId)
			keepIndIds = append(keepIndIds, ind.Id)
		}
		keepTargetIds := make([]int, 0, len(ind.Target))
		for _, t := range ind.Target {
			if t.Id > 0 {
				// UPDATE target existing — id tetap
				_, err = tx.ExecContext(ctx, `
				UPDATE tb_target_pemda
				SET target=?, satuan=?, tahun=?
				WHERE id=? AND kode_indikator=?
				  AND (jenis = 'renstra' OR jenis = '' OR jenis IS NULL)`,
					t.Target, t.Satuan, t.Tahun, t.Id, kodeInd)
				if err != nil {
					return tp, err
				}
				keepTargetIds = append(keepTargetIds, t.Id)
			} else {
				// INSERT target baru
				res, err := tx.ExecContext(ctx, `
					INSERT INTO tb_target_pemda (kode_indikator, target, satuan, tahun, jenis)
					VALUES (?, ?, ?, ?, ?)`,
					kodeInd, t.Target, t.Satuan, t.Tahun, t.Jenis)
				if err != nil {
					return tp, err
				}
				newId, _ := res.LastInsertId()
				keepTargetIds = append(keepTargetIds, int(newId))
			}
		}
		// Hapus target lama yang tidak ada di request (per indikator)
		if len(keepTargetIds) > 0 {
			placeholders := strings.Repeat("?,", len(keepTargetIds))
			placeholders = placeholders[:len(placeholders)-1]
			args := make([]interface{}, 0, len(keepTargetIds)+1)
			args = append(args, kodeInd)
			for _, id := range keepTargetIds {
				args = append(args, id)
			}
			_, err = tx.ExecContext(ctx,
				fmt.Sprintf(`DELETE FROM tb_target_pemda
                 WHERE kode_indikator=?
                   AND (jenis = 'renstra' OR jenis = '' OR jenis IS NULL)
                   AND id NOT IN (%s)`, placeholders),
				args...)
			if err != nil {
				return tp, err
			}
		}
	}
	// Hapus indikator lama yang tidak ada di request
	if len(keepIndIds) > 0 {
		placeholders := strings.Repeat("?,", len(keepIndIds))
		placeholders = placeholders[:len(placeholders)-1]
		args := make([]interface{}, 0, len(keepIndIds)+1)
		args = append(args, tp.Id)
		for _, id := range keepIndIds {
			args = append(args, id)
		}
		// Hapus target orphan dulu
		_, err = tx.ExecContext(ctx,
			fmt.Sprintf(`DELETE tg FROM tb_target_pemda tg
						 INNER JOIN tb_indikator_matrix_pemda i ON tg.kode_indikator = i.kode_indikator
						 WHERE i.tujuan_pemda_id=? AND i.jenis='renstra' AND i.id NOT IN (%s)
						   AND (tg.jenis = 'renstra' OR tg.jenis = '' OR tg.jenis IS NULL)`, placeholders),
			args...)
		if err != nil {
			return tp, err
		}
		_, err = tx.ExecContext(ctx,
			fmt.Sprintf(`DELETE FROM tb_indikator_matrix_pemda
			             WHERE tujuan_pemda_id=? AND jenis='renstra' AND id NOT IN (%s)`, placeholders),
			args...)
		if err != nil {
			return tp, err
		}
	}
	return tp, nil
}

func (r *TujuanPemdaRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, tujuanPemdaId int) error {
	// 1. Hapus target pemda
	_, err := tx.ExecContext(ctx, `
		DELETE tg FROM tb_target_pemda tg
		INNER JOIN tb_indikator_matrix_pemda i ON tg.kode_indikator = i.kode_indikator
		WHERE i.tujuan_pemda_id = ?`, tujuanPemdaId)
	if err != nil {
		return err
	}
	// 2. Hapus indikator matrix
	_, err = tx.ExecContext(ctx,
		"DELETE FROM tb_indikator_matrix_pemda WHERE tujuan_pemda_id = ?", tujuanPemdaId)
	if err != nil {
		return err
	}
	// 3. Hapus tujuan pemda
	_, err = tx.ExecContext(ctx,
		"DELETE FROM tb_tujuan_pemda WHERE id = ?", tujuanPemdaId)
	return err
}

func (r *TujuanPemdaRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, tujuanPemdaId int) (domain.TujuanPemda, error) {
	query := `
		SELECT
			tp.id,
			tp.tujuan_pemda,
			tp.tematik_id,
			tp.periode_id,
			tp.tahun_awal_periode,
			tp.tahun_akhir_periode,
			tp.jenis_periode,
			tp.id_visi,
			tp.id_misi,
			COALESCE(pk.jenis_pohon, '') AS jenis_pohon,
			i.id                         AS indikator_db_id,
			i.kode_indikator,
			i.indikator                  AS indikator_text,
			i.rumus_perhitungan,
			i.sumber_data,
			i.definisi_operasional,
			i.jenis                      AS indikator_jenis,
			t.id                         AS target_db_id,
			t.target                     AS target_value,
			t.satuan,
			t.tahun                      AS target_tahun,
			t.jenis                      AS target_jenis
		FROM tb_tujuan_pemda tp
		LEFT JOIN tb_pohon_kinerja pk ON tp.tematik_id = pk.id
		LEFT JOIN tb_indikator_matrix_pemda i
			ON tp.id = i.tujuan_pemda_id AND i.jenis = 'renstra'
		LEFT JOIN tb_target_pemda t
			ON t.kode_indikator = i.kode_indikator AND t.jenis = 'renstra'
		WHERE tp.id = ?
		ORDER BY tp.id, i.kode_indikator, CAST(t.tahun AS SIGNED)`
	rows, err := tx.QueryContext(ctx, query, tujuanPemdaId)
	if err != nil {
		return domain.TujuanPemda{}, fmt.Errorf("error querying tujuan pemda: %v", err)
	}
	defer rows.Close()
	var result domain.TujuanPemda
	firstRow := true
	indikatorMap := make(map[string]*domain.IndikatorPemda)
	for rows.Next() {
		var (
			id, tematikId, periodeId, idVisi, idMisi             int
			tujuanPemdaText, tahunAwal, tahunAkhir, jenisPeriode string
			jenisPohon                                           string
			indikatorDbId                                        sql.NullInt64
			kodeIndikator, indikatorText                         sql.NullString
			rumusPerhitungan, sumberData, definisiOp             sql.NullString
			indikatorJenis                                       sql.NullString
			targetDbId                                           sql.NullInt64
			targetValue, targetSatuan, targetTahun, targetJenis  sql.NullString
		)
		err := rows.Scan(
			&id, &tujuanPemdaText, &tematikId, &periodeId,
			&tahunAwal, &tahunAkhir, &jenisPeriode,
			&idVisi, &idMisi, &jenisPohon,
			&indikatorDbId, &kodeIndikator, &indikatorText,
			&rumusPerhitungan, &sumberData, &definisiOp, &indikatorJenis,
			&targetDbId, &targetValue, &targetSatuan, &targetTahun, &targetJenis,
		)
		if err != nil {
			return domain.TujuanPemda{}, fmt.Errorf("error scanning row: %v", err)
		}
		if firstRow {
			result = domain.TujuanPemda{
				Id:                id,
				TujuanPemda:       tujuanPemdaText,
				TematikId:         tematikId,
				JenisPohon:        jenisPohon,
				PeriodeId:         periodeId,
				IdVisi:            idVisi,
				IdMisi:            idMisi,
				TahunAwalPeriode:  tahunAwal,
				TahunAkhirPeriode: tahunAkhir,
				JenisPeriode:      jenisPeriode,
				Periode: domain.Periode{
					TahunAwal:    tahunAwal,
					TahunAkhir:   tahunAkhir,
					JenisPeriode: jenisPeriode,
				},
				IndikatorPemda: []domain.IndikatorPemda{},
			}
			firstRow = false
		}
		if indikatorDbId.Valid && kodeIndikator.Valid {
			ind, exists := indikatorMap[kodeIndikator.String]
			if !exists {
				newInd := domain.IndikatorPemda{
					Id:                  int(indikatorDbId.Int64),
					KodeIndikator:       kodeIndikator.String,
					TujuanPemdaId:       id,
					Indikator:           indikatorText,
					RumusPerhitungan:    rumusPerhitungan,
					SumberData:          sumberData,
					DefinisiOperasional: definisiOp,
					Jenis:               indikatorJenis.String,
					Target:              []domain.TargetPemda{},
				}
				// Buat placeholder target untuk setiap tahun dalam periode
				tahunAwalInt, errA := strconv.Atoi(tahunAwal)
				tahunAkhirInt, errB := strconv.Atoi(tahunAkhir)
				if errA == nil && errB == nil {
					for y := tahunAwalInt; y <= tahunAkhirInt; y++ {
						newInd.Target = append(newInd.Target, domain.TargetPemda{
							Id:            0,
							KodeIndikator: kodeIndikator.String,
							Target:        "-",
							Satuan:        "-",
							Tahun:         strconv.Itoa(y),
							Jenis:         "renstra",
						})
					}
				}
				result.IndikatorPemda = append(result.IndikatorPemda, newInd)
				indikatorMap[kodeIndikator.String] = &result.IndikatorPemda[len(result.IndikatorPemda)-1]
				ind = &result.IndikatorPemda[len(result.IndikatorPemda)-1]
			}
			// Isi target yang sebenarnya
			if targetDbId.Valid && targetValue.Valid && targetTahun.Valid {
				for i := range ind.Target {
					if ind.Target[i].Tahun == targetTahun.String {
						ind.Target[i] = domain.TargetPemda{
							Id:            int(targetDbId.Int64),
							KodeIndikator: kodeIndikator.String,
							Target:        targetValue.String,
							Satuan:        targetSatuan.String,
							Tahun:         targetTahun.String,
							Jenis:         targetJenis.String,
						}
						break
					}
				}
			}
		}
	}
	if err = rows.Err(); err != nil {
		return domain.TujuanPemda{}, fmt.Errorf("error iterating rows: %v", err)
	}
	if result.Id == 0 {
		return domain.TujuanPemda{}, fmt.Errorf("tujuan pemda dengan id %d tidak ditemukan", tujuanPemdaId)
	}
	return result, nil
}

func (repository *TujuanPemdaRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, tahun string, jenisPeriode string) ([]domain.TujuanPemda, error) {
	query := `
        SELECT DISTINCT
            tp.id,
            tp.tujuan_pemda,
            tp.tematik_id,
            tp.periode_id,
            p.tahun_awal,
            p.tahun_akhir,
            p.jenis_periode
        FROM 
            tb_tujuan_pemda tp
            INNER JOIN tb_periode p ON tp.periode_id = p.id
            INNER JOIN tb_pohon_kinerja pk ON tp.tematik_id = pk.id
        WHERE 
            CAST(? AS SIGNED) BETWEEN CAST(p.tahun_awal AS SIGNED) AND CAST(p.tahun_akhir AS SIGNED)
            AND p.jenis_periode = ?
            AND pk.is_active = true
            AND pk.level_pohon = 0
        ORDER BY 
            tp.id`

	rows, err := tx.QueryContext(ctx, query, tahun, jenisPeriode)
	if err != nil {
		return nil, fmt.Errorf("error querying tujuan pemda: %v", err)
	}
	defer rows.Close()

	var result []domain.TujuanPemda

	for rows.Next() {
		var tujuanPemda domain.TujuanPemda
		err := rows.Scan(
			&tujuanPemda.Id,
			&tujuanPemda.TujuanPemda,
			&tujuanPemda.TematikId,
			&tujuanPemda.PeriodeId,
			&tujuanPemda.Periode.TahunAwal,
			&tujuanPemda.Periode.TahunAkhir,
			&tujuanPemda.Periode.JenisPeriode,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		result = append(result, tujuanPemda)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return result, nil
}
func (repository *TujuanPemdaRepositoryImpl) FindAllBetweenTahun(ctx context.Context, tx *sql.Tx, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]domain.TujuanPemda, error) {
	query := `
        SELECT DISTINCT
            tp.id,
            tp.tujuan_pemda,
            tp.tematik_id,
            tp.periode_id,
            p.tahun_awal,
            p.tahun_akhir,
            p.jenis_periode
        FROM 
            tb_tujuan_pemda tp
            INNER JOIN tb_periode p ON tp.periode_id = p.id
            INNER JOIN tb_pohon_kinerja pk ON tp.tematik_id = pk.id
        WHERE 
            p.tahun_awal <= ?
			AND p.tahun_akhir >= ?
            AND p.jenis_periode = ?
            AND pk.is_active = true
            AND pk.level_pohon = 0
        ORDER BY 
            tp.id`

	rows, err := tx.QueryContext(ctx, query, tahunAkhir, tahunAwal, jenisPeriode)
	if err != nil {
		return nil, fmt.Errorf("error querying tujuan pemda: %v", err)
	}
	defer rows.Close()

	var result []domain.TujuanPemda

	for rows.Next() {
		var tujuanPemda domain.TujuanPemda
		err := rows.Scan(
			&tujuanPemda.Id,
			&tujuanPemda.TujuanPemda,
			&tujuanPemda.TematikId,
			&tujuanPemda.PeriodeId,
			&tujuanPemda.Periode.TahunAwal,
			&tujuanPemda.Periode.TahunAkhir,
			&tujuanPemda.Periode.JenisPeriode,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		result = append(result, tujuanPemda)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return result, nil
}

func (repository *TujuanPemdaRepositoryImpl) DeleteIndikator(ctx context.Context, tx *sql.Tx, tujuanPemdaId int) error {
	query := "DELETE FROM tb_indikator WHERE tujuan_pemda_id = ?"
	_, err := tx.ExecContext(ctx, query, tujuanPemdaId)
	return err
}

func (repository *TujuanPemdaRepositoryImpl) IsIdExists(ctx context.Context, tx *sql.Tx, id int) bool {
	query := "SELECT COUNT(*) FROM tb_tujuan_pemda WHERE id = ?"
	var count int
	err := tx.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		return true
	}
	return count > 0
}

func (repository *TujuanPemdaRepositoryImpl) UpdatePeriode(ctx context.Context, tx *sql.Tx, tujuanPemda domain.TujuanPemda) (domain.TujuanPemda, error) {
	// Update hanya periode_id
	query := "UPDATE tb_tujuan_pemda SET periode_id = ? WHERE id = ?"
	result, err := tx.ExecContext(ctx, query, tujuanPemda.PeriodeId, tujuanPemda.Id)
	if err != nil {
		return domain.TujuanPemda{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.TujuanPemda{}, err
	}

	if rowsAffected == 0 {
		return domain.TujuanPemda{}, fmt.Errorf("periode dengan id %d sudah digunakan", tujuanPemda.PeriodeId)
	}

	// Ambil data terbaru setelah update
	query = `
        SELECT 
            tp.id,
            tp.tujuan_pemda,
            tp.tematik_id,
            tp.periode_id,
            COALESCE(p.tahun_awal, 'Pilih periode') as tahun_awal,
            COALESCE(p.tahun_akhir, 'Pilih periode') as tahun_akhir
        FROM 
            tb_tujuan_pemda tp
            LEFT JOIN tb_periode p ON tp.periode_id = p.id
        WHERE tp.id = ?`

	var updatedTujuanPemda domain.TujuanPemda
	err = tx.QueryRowContext(ctx, query, tujuanPemda.Id).Scan(
		&updatedTujuanPemda.Id,
		&updatedTujuanPemda.TujuanPemda,
		&updatedTujuanPemda.TematikId,
		&updatedTujuanPemda.PeriodeId,
		&updatedTujuanPemda.Periode.TahunAwal,
		&updatedTujuanPemda.Periode.TahunAkhir,
	)
	if err != nil {
		return domain.TujuanPemda{}, fmt.Errorf("gagal mengambil data setelah update: %v", err)
	}

	return updatedTujuanPemda, nil
}

func (r *TujuanPemdaRepositoryImpl) findAllWithPokinQuery(ctx context.Context, tx *sql.Tx, tahunAwal, tahunAkhir, jenisPeriode, jenisFilter string) ([]domain.TujuanPemdaWithPokin, error) {
	// Validasi periode
	var exists bool
	err := tx.QueryRowContext(ctx, `
		SELECT EXISTS (SELECT 1 FROM tb_periode
			WHERE tahun_awal=? AND tahun_akhir=? AND jenis_periode=?)`,
		tahunAwal, tahunAkhir, jenisPeriode).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("periode %s-%s %s tidak ditemukan", tahunAwal, tahunAkhir, jenisPeriode)
	}
	// Kondisi filter jenis untuk JOIN
	indJenisClause := ""
	tgJenisClause := ""
	if jenisFilter != "" {
		indJenisClause = fmt.Sprintf(" AND i.jenis = '%s'", jenisFilter)
		tgJenisClause = fmt.Sprintf(" AND tg.jenis = '%s'", jenisFilter)
	}
	query := fmt.Sprintf(`
		SELECT
			pk.id          AS pokin_id,
			pk.nama_pohon,
			pk.jenis_pohon,
			pk.level_pohon,
			pk.is_active,
			pk.kode_opd,
			pk.keterangan,
			pk.tahun        AS tahun_pokin,
			tp.id           AS tujuan_id,
			tp.tujuan_pemda,
			tp.id_visi,
			tp.id_misi,
			tp.tahun_awal_periode,
			tp.tahun_akhir_periode,
			tp.jenis_periode,
			i.id            AS indikator_db_id,
			i.kode_indikator,
			i.indikator     AS indikator_text,
			i.rumus_perhitungan,
			i.sumber_data,
			i.definisi_operasional,
			i.jenis         AS indikator_jenis,
			tg.id           AS target_db_id,
			tg.target,
			tg.satuan,
			tg.tahun        AS target_tahun,
			tg.jenis        AS target_jenis
		FROM tb_pohon_kinerja pk
		LEFT JOIN tb_tujuan_pemda tp
			ON pk.id = tp.tematik_id
			AND tp.tahun_awal_periode = ?
			AND tp.tahun_akhir_periode = ?
			AND tp.jenis_periode = ?
		LEFT JOIN tb_indikator_matrix_pemda i
			ON tp.id = i.tujuan_pemda_id%s
		LEFT JOIN tb_target_pemda tg
			ON tg.kode_indikator = i.kode_indikator
			AND CAST(tg.tahun AS SIGNED) BETWEEN CAST(? AS SIGNED) AND CAST(? AS SIGNED)%s
		WHERE pk.level_pohon = 0
		  AND CAST(pk.tahun AS SIGNED) BETWEEN CAST(? AS SIGNED) AND CAST(? AS SIGNED)
		ORDER BY pk.id, tp.id, i.kode_indikator, tg.tahun`,
		indJenisClause, tgJenisClause)
	rows, err := tx.QueryContext(ctx, query,
		tahunAwal, tahunAkhir, jenisPeriode,
		tahunAwal, tahunAkhir,
		tahunAwal, tahunAkhir)
	if err != nil {
		return nil, fmt.Errorf("error querying data: %v", err)
	}
	defer rows.Close()
	pokinMap := make(map[int]*domain.TujuanPemdaWithPokin)
	indikatorMap := make(map[string]*domain.IndikatorPemda) // key: kode_indikator
	for rows.Next() {
		var (
			pokinId                                             int
			namaPohon, jenisPohon, kodeOpd, keterangan, tahunPk string
			levelPohon                                          int
			isActive                                            bool
			tujuanId                                            sql.NullInt64
			tujuanPemdaText                                     sql.NullString
			idVisi, idMisi                                      sql.NullInt64
			tahunAwalP, tahunAkhirP, jenisPeriodeP              sql.NullString
			indikatorDbId                                       sql.NullInt64
			kodeIndikator, indikatorText                        sql.NullString
			rumusPerhitungan, sumberData, definisiOp            sql.NullString
			indikatorJenis                                      sql.NullString
			targetDbId                                          sql.NullInt64
			targetValue, targetSatuan, targetTahun, targetJenis sql.NullString
		)
		err := rows.Scan(
			&pokinId, &namaPohon, &jenisPohon, &levelPohon,
			&isActive, &kodeOpd, &keterangan, &tahunPk,
			&tujuanId, &tujuanPemdaText,
			&idVisi, &idMisi,
			&tahunAwalP, &tahunAkhirP, &jenisPeriodeP,
			&indikatorDbId, &kodeIndikator, &indikatorText,
			&rumusPerhitungan, &sumberData, &definisiOp, &indikatorJenis,
			&targetDbId, &targetValue, &targetSatuan, &targetTahun, &targetJenis,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		// Inisialisasi pokin
		pokin, ok := pokinMap[pokinId]
		if !ok {
			pokin = &domain.TujuanPemdaWithPokin{
				PokinId:     pokinId,
				NamaPohon:   namaPohon,
				JenisPohon:  jenisPohon,
				LevelPohon:  levelPohon,
				IsActive:    isActive,
				KodeOpd:     kodeOpd,
				Keterangan:  keterangan,
				TahunPokin:  tahunPk,
				TujuanPemda: []domain.TujuanPemda{},
			}
			pokinMap[pokinId] = pokin
		}
		if !tujuanId.Valid || !tujuanPemdaText.Valid {
			continue
		}
		// Cari atau buat TujuanPemda
		var existingTP *domain.TujuanPemda
		for i := range pokin.TujuanPemda {
			if pokin.TujuanPemda[i].Id == int(tujuanId.Int64) {
				existingTP = &pokin.TujuanPemda[i]
				break
			}
		}
		if existingTP == nil {
			newTP := domain.TujuanPemda{
				Id:                int(tujuanId.Int64),
				TujuanPemda:       tujuanPemdaText.String,
				TematikId:         pokinId,
				IsActive:          isActive,
				IdVisi:            int(idVisi.Int64),
				IdMisi:            int(idMisi.Int64),
				TahunAwalPeriode:  tahunAwalP.String,
				TahunAkhirPeriode: tahunAkhirP.String,
				JenisPeriode:      jenisPeriodeP.String,
				IndikatorPemda:    []domain.IndikatorPemda{},
			}
			pokin.TujuanPemda = append(pokin.TujuanPemda, newTP)
			existingTP = &pokin.TujuanPemda[len(pokin.TujuanPemda)-1]
		}
		if !indikatorDbId.Valid || !kodeIndikator.Valid {
			continue
		}
		// Cari atau buat IndikatorPemda
		ind, indExists := indikatorMap[kodeIndikator.String]
		if !indExists {
			newInd := domain.IndikatorPemda{
				Id:                  int(indikatorDbId.Int64),
				KodeIndikator:       kodeIndikator.String,
				TujuanPemdaId:       int(tujuanId.Int64),
				Indikator:           indikatorText,
				RumusPerhitungan:    rumusPerhitungan,
				SumberData:          sumberData,
				DefinisiOperasional: definisiOp,
				Jenis:               indikatorJenis.String,
				Target:              []domain.TargetPemda{},
			}
			// Placeholder target per tahun
			tpAwal, _ := strconv.Atoi(tahunAwal)
			tpAkhir, _ := strconv.Atoi(tahunAkhir)
			for y := tpAwal; y <= tpAkhir; y++ {
				newInd.Target = append(newInd.Target, domain.TargetPemda{
					Id:            0,
					KodeIndikator: kodeIndikator.String,
					Target:        "-",
					Satuan:        "-",
					Tahun:         strconv.Itoa(y),
					Jenis:         jenisFilter,
				})
			}
			existingTP.IndikatorPemda = append(existingTP.IndikatorPemda, newInd)
			indikatorMap[kodeIndikator.String] = &existingTP.IndikatorPemda[len(existingTP.IndikatorPemda)-1]
			ind = &existingTP.IndikatorPemda[len(existingTP.IndikatorPemda)-1]
		}
		// Isi target yang sesungguhnya
		if targetDbId.Valid && targetValue.Valid && targetTahun.Valid {
			for i := range ind.Target {
				if ind.Target[i].Tahun == targetTahun.String {
					ind.Target[i] = domain.TargetPemda{
						Id:            int(targetDbId.Int64),
						KodeIndikator: kodeIndikator.String,
						Target:        targetValue.String,
						Satuan:        targetSatuan.String,
						Tahun:         targetTahun.String,
						Jenis:         targetJenis.String,
					}
					break
				}
			}
		}
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}
	result := make([]domain.TujuanPemdaWithPokin, 0, len(pokinMap))
	for _, p := range pokinMap {
		result = append(result, *p)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].PokinId < result[j].PokinId
	})
	return result, nil
}

// FindAllWithPokin — semua jenis indikator + target (tanpa filter jenis)
func (r *TujuanPemdaRepositoryImpl) FindAllWithPokin(ctx context.Context, tx *sql.Tx,
	tahunAwal, tahunAkhir, jenisPeriode string) ([]domain.TujuanPemdaWithPokin, error) {
	return r.findAllWithPokinQuery(ctx, tx, tahunAwal, tahunAkhir, jenisPeriode, "")
}

// FindAllWithPokinRenstra — khusus jenis='renstra' (5 tahunan)
func (r *TujuanPemdaRepositoryImpl) FindAllWithPokinRenstra(ctx context.Context, tx *sql.Tx,
	tahunAwal, tahunAkhir, jenisPeriode string) ([]domain.TujuanPemdaWithPokin, error) {
	return r.findAllWithPokinQuery(ctx, tx, tahunAwal, tahunAkhir, jenisPeriode, "renstra")
}
func (repository *TujuanPemdaRepositoryImpl) IsPokinIdExists(ctx context.Context, tx *sql.Tx, pokinId int) (bool, error) {
	query := "SELECT COUNT(*) FROM tb_tujuan_pemda WHERE tematik_id = ?"
	var count int
	err := tx.QueryRowContext(ctx, query, pokinId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *TujuanPemdaRepositoryImpl) TargetPemdaExistsByKey(
	ctx context.Context, tx *sql.Tx,
	kodeIndikator, tahun, jenis string,
) (bool, error) {
	var count int
	err := tx.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM tb_target_pemda
		WHERE kode_indikator = ? AND tahun = ? AND jenis = ?`,
		kodeIndikator, tahun, jenis,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *TujuanPemdaRepositoryImpl) FindTargetPemdaById(
	ctx context.Context, tx *sql.Tx, id int,
) (domain.TargetPemda, error) {
	query := `
		SELECT id, kode_indikator, target, satuan, tahun, jenis
		FROM tb_target_pemda
		WHERE id = ?`
	var t domain.TargetPemda
	err := tx.QueryRowContext(ctx, query, id).Scan(
		&t.Id, &t.KodeIndikator, &t.Target, &t.Satuan, &t.Tahun, &t.Jenis,
	)
	return t, err
}

// UpdateTargetPemda — hanya update target & satuan, jenis TIDAK diubah
func (r *TujuanPemdaRepositoryImpl) UpdateTargetPemda(
	ctx context.Context, tx *sql.Tx, id int, target, satuan string,
) (domain.TargetPemda, error) {
	result, err := tx.ExecContext(ctx, `
		UPDATE tb_target_pemda
		SET target = ?, satuan = ?
		WHERE id = ?`,
		target, satuan, id,
	)
	if err != nil {
		return domain.TargetPemda{}, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return domain.TargetPemda{}, err
	}
	if rows == 0 {
		return domain.TargetPemda{}, fmt.Errorf("target id %d tidak ditemukan", id)
	}
	// Kembalikan data lengkap (jenis tetap dari DB)
	return r.FindTargetPemdaById(ctx, tx, id)
}

// ─────────────────────────────────────────────────────────────────
// FUNC BARU: FindAllWithPokinByTargetJenis
// Sama seperti FindAllWithPokinRenstra, tapi target di-filter per jenis layer.
// Indikator SELALU dari jenis 'renstra' (metadata tidak berubah antar layer).
// ─────────────────────────────────────────────────────────────────
func (r *TujuanPemdaRepositoryImpl) FindAllWithPokinByTargetJenis(
	ctx context.Context, tx *sql.Tx,
	tahunAwal, tahunAkhir, jenisPeriode, targetJenis string,
) ([]domain.TujuanPemdaWithPokin, error) {
	// Validasi periode dulu
	var exists bool
	err := tx.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM tb_periode
			WHERE tahun_awal = ? AND tahun_akhir = ? AND jenis_periode = ?
		)`, tahunAwal, tahunAkhir, jenisPeriode).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("periode %s-%s %s tidak ditemukan", tahunAwal, tahunAkhir, jenisPeriode)
	}
	query := `
		SELECT
			pk.id,
			pk.nama_pohon,
			pk.jenis_pohon,
			pk.level_pohon,
			pk.is_active,
			pk.kode_opd,
			pk.keterangan,
			pk.tahun AS tahun_pokin,
			tp.id           AS tujuan_id,
			tp.tujuan_pemda,
			tp.id_visi,
			tp.id_misi,
			tp.tahun_awal_periode,
			tp.tahun_akhir_periode,
			tp.jenis_periode,
			i.id            AS indikator_id,
			i.kode_indikator,
			i.indikator,
			i.rumus_perhitungan,
			i.sumber_data,
			i.definisi_operasional,
			i.jenis         AS indikator_jenis,
			tg.id           AS target_id,
			tg.target,
			tg.satuan,
			tg.tahun        AS target_tahun,
			tg.jenis        AS target_jenis
		FROM tb_pohon_kinerja pk
		LEFT JOIN tb_tujuan_pemda tp
			ON pk.id = tp.tematik_id
			AND tp.tahun_awal_periode  = ?
			AND tp.tahun_akhir_periode = ?
			AND tp.jenis_periode       = ?
		LEFT JOIN tb_indikator_matrix_pemda i
			ON tp.id = i.tujuan_pemda_id
			-- indikator SELALU dari renstra; '' dan NULL untuk backward compat data lama
			AND (i.jenis = 'renstra' OR i.jenis = '' OR i.jenis IS NULL)
		LEFT JOIN tb_target_pemda tg
			ON tg.kode_indikator = i.kode_indikator
			AND CAST(tg.tahun AS SIGNED) BETWEEN CAST(? AS SIGNED) AND CAST(? AS SIGNED)
			AND tg.jenis = ?
		WHERE pk.level_pohon = 0
		  AND CAST(pk.tahun AS SIGNED) BETWEEN CAST(? AS SIGNED) AND CAST(? AS SIGNED)
		ORDER BY pk.id, tp.id, i.kode_indikator, tg.tahun`
	rows, err := tx.QueryContext(ctx, query,
		tahunAwal, tahunAkhir, jenisPeriode, // tp filter
		tahunAwal, tahunAkhir, targetJenis, // tg filter
		tahunAwal, tahunAkhir, // pk filter
	)
	if err != nil {
		return nil, fmt.Errorf("FindAllWithPokinByTargetJenis query error: %v", err)
	}
	defer rows.Close()
	// Scan dan build hasil
	pokinMap := make(map[int]*domain.TujuanPemdaWithPokin)
	// key: "tujuanId:kodeIndikator" — hindari collision antar tujuan pemda
	indikatorMap := make(map[string]*domain.IndikatorPemda)
	for rows.Next() {
		var (
			pokinId                                               int
			namaPohon, jenisPohon, kodeOpd, keterangan, tahunPk   string
			levelPohon                                            int
			isActive                                              bool
			tujuanId                                              sql.NullInt64
			tujuanText                                            sql.NullString
			idVisi, idMisi                                        sql.NullInt64
			tahunAwalP, tahunAkhirP, jenisPeriodeP                sql.NullString
			indikatorId                                           sql.NullInt64
			kodeIndikator, indikatorText                          sql.NullString
			rumus, sumber, definisi, indikatorJenis               sql.NullString
			targetId                                              sql.NullInt64
			targetVal, targetSatuan, targetTahun, targetJenisScan sql.NullString
		)
		if err := rows.Scan(
			&pokinId, &namaPohon, &jenisPohon, &levelPohon, &isActive, &kodeOpd, &keterangan, &tahunPk,
			&tujuanId, &tujuanText, &idVisi, &idMisi, &tahunAwalP, &tahunAkhirP, &jenisPeriodeP,
			&indikatorId, &kodeIndikator, &indikatorText, &rumus, &sumber, &definisi, &indikatorJenis,
			&targetId, &targetVal, &targetSatuan, &targetTahun, &targetJenisScan,
		); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		// Inisialisasi pokin jika belum ada
		if _, ok := pokinMap[pokinId]; !ok {
			pokinMap[pokinId] = &domain.TujuanPemdaWithPokin{
				PokinId:     pokinId,
				NamaPohon:   namaPohon,
				JenisPohon:  jenisPohon,
				LevelPohon:  levelPohon,
				IsActive:    isActive,
				KodeOpd:     kodeOpd,
				Keterangan:  keterangan,
				TahunPokin:  tahunPk,
				TujuanPemda: []domain.TujuanPemda{},
			}
		}
		pokin := pokinMap[pokinId]
		// Lewati baris yang tidak punya tujuan
		if !tujuanId.Valid {
			continue
		}
		// Cari atau buat TujuanPemda di dalam pokin
		tpId := int(tujuanId.Int64)
		var currentTP *domain.TujuanPemda
		for i := range pokin.TujuanPemda {
			if pokin.TujuanPemda[i].Id == tpId {
				currentTP = &pokin.TujuanPemda[i]
				break
			}
		}
		if currentTP == nil {
			pokin.TujuanPemda = append(pokin.TujuanPemda, domain.TujuanPemda{
				Id:                tpId,
				TujuanPemda:       tujuanText.String,
				TematikId:         pokinId,
				IdVisi:            int(idVisi.Int64),
				IdMisi:            int(idMisi.Int64),
				TahunAwalPeriode:  tahunAwalP.String,
				TahunAkhirPeriode: tahunAkhirP.String,
				JenisPeriode:      jenisPeriodeP.String,
				IndikatorPemda:    []domain.IndikatorPemda{},
			})
			currentTP = &pokin.TujuanPemda[len(pokin.TujuanPemda)-1]
		}
		// Lewati baris yang tidak punya indikator
		if !indikatorId.Valid || !kodeIndikator.Valid {
			continue
		}
		// key unik per (tujuan, kode_indikator)
		indKey := fmt.Sprintf("%d:%s", tpId, kodeIndikator.String)
		currentInd, indExists := indikatorMap[indKey]
		if !indExists {
			// Buat indikator baru dengan placeholder target untuk setiap tahun
			newInd := domain.IndikatorPemda{
				Id:                  int(indikatorId.Int64),
				KodeIndikator:       kodeIndikator.String,
				TujuanPemdaId:       tpId,
				Indikator:           indikatorText,
				RumusPerhitungan:    rumus,
				SumberData:          sumber,
				DefinisiOperasional: definisi,
				Jenis:               indikatorJenis.String,
				Target:              []domain.TargetPemda{},
			}
			// Buat slot placeholder per tahun (akan diisi oleh data real di bawah)
			tpAwal, _ := strconv.Atoi(tahunAwal)
			tpAkhir, _ := strconv.Atoi(tahunAkhir)
			for y := tpAwal; y <= tpAkhir; y++ {
				newInd.Target = append(newInd.Target, domain.TargetPemda{
					Id:            0,
					KodeIndikator: kodeIndikator.String,
					Target:        "-",
					Satuan:        "-",
					Tahun:         strconv.Itoa(y),
					Jenis:         targetJenis,
				})
			}
			currentTP.IndikatorPemda = append(currentTP.IndikatorPemda, newInd)
			currentInd = &currentTP.IndikatorPemda[len(currentTP.IndikatorPemda)-1]
			indikatorMap[indKey] = currentInd
		}
		// Isi target real (jika ada)
		if targetId.Valid && targetVal.Valid && targetTahun.Valid {
			trimTahun := strings.TrimSpace(targetTahun.String)
			for i := range currentInd.Target {
				if strings.TrimSpace(currentInd.Target[i].Tahun) == trimTahun {
					currentInd.Target[i] = domain.TargetPemda{
						Id:            int(targetId.Int64),
						KodeIndikator: kodeIndikator.String,
						Target:        targetVal.String,
						Satuan:        targetSatuan.String,
						Tahun:         trimTahun,
						Jenis:         targetJenisScan.String,
					}
					break
				}
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	// Ubah map ke slice, urutkan per pokin id
	result := make([]domain.TujuanPemdaWithPokin, 0, len(pokinMap))
	for _, p := range pokinMap {
		result = append(result, *p)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].PokinId < result[j].PokinId
	})
	return result, nil
}

// ─────────────────────────────────────────────────────────────────
// FUNC BARU: FindIndikatorPemdaByKode
// Dipakai untuk validasi sebelum upsert target layer.
// Cukup cek apakah indikator renstra-nya ada.
// ─────────────────────────────────────────────────────────────────
func (r *TujuanPemdaRepositoryImpl) FindIndikatorPemdaByKode(
	ctx context.Context, tx *sql.Tx, kodeIndikator string,
) (domain.IndikatorPemda, error) {
	query := `
		SELECT id, tujuan_pemda_id, kode_indikator, indikator,
		       rumus_perhitungan, sumber_data, definisi_operasional, jenis
		FROM tb_indikator_matrix_pemda
		WHERE kode_indikator = ?
		  AND (jenis = 'renstra' OR jenis = '' OR jenis IS NULL)
		LIMIT 1`
	var ind domain.IndikatorPemda
	err := tx.QueryRowContext(ctx, query, kodeIndikator).Scan(
		&ind.Id,
		&ind.TujuanPemdaId,
		&ind.KodeIndikator,
		&ind.Indikator,
		&ind.RumusPerhitungan,
		&ind.SumberData,
		&ind.DefinisiOperasional,
		&ind.Jenis,
	)
	if err != nil {
		return domain.IndikatorPemda{}, err // sql.ErrNoRows jika tidak ada
	}
	return ind, nil
}

// ─────────────────────────────────────────────────────────────────
// FUNC BARU: UpsertTargetPemda
// Insert jika belum ada, update jika sudah ada.
// Key unik: kode_indikator + tahun + jenis.
// ─────────────────────────────────────────────────────────────────
func (r *TujuanPemdaRepositoryImpl) UpsertTargetPemda(
	ctx context.Context, tx *sql.Tx, t domain.TargetPemda,
) (domain.TargetPemda, error) {
	// Cek apakah sudah ada
	var existingId int
	err := tx.QueryRowContext(ctx, `
		SELECT id FROM tb_target_pemda
		WHERE kode_indikator = ? AND tahun = ? AND jenis = ?
		LIMIT 1`,
		t.KodeIndikator, t.Tahun, t.Jenis,
	).Scan(&existingId)
	if err == sql.ErrNoRows {
		// Belum ada → insert baru
		return r.CreateTarget(ctx, tx, t)
	}
	if err != nil {
		return t, err
	}
	// Sudah ada → update saja
	_, err = tx.ExecContext(ctx, `
		UPDATE tb_target_pemda
		SET target = ?, satuan = ?
		WHERE id = ?`,
		t.Target, t.Satuan, existingId,
	)
	t.Id = existingId
	return t, err
}

// ─────────────────────────────────────────────────────────────────
// FindAllByTahun
// Langsung dari tb_tujuan_pemda (TANPA wrapper pohon kinerja).
// Indikator selalu dari renstra, target sesuai targetJenis untuk 1 tahun.
// ─────────────────────────────────────────────────────────────────
func (r *TujuanPemdaRepositoryImpl) FindAllByTahun(
	ctx context.Context, tx *sql.Tx,
	tahun, jenisPeriode, targetJenis string,
) ([]domain.TujuanPemda, error) {
	var targetJenisClause string
	var args []interface{}
	if targetJenis == "renstra" {
		targetJenisClause = "(tg.jenis = 'renstra' OR tg.jenis = '' OR tg.jenis IS NULL)"
		args = []interface{}{tahun, jenisPeriode, tahun, tahun}
	} else {
		targetJenisClause = "tg.jenis = ?"
		args = []interface{}{tahun, targetJenis, jenisPeriode, tahun, tahun}
	}
	query := fmt.Sprintf(`
		SELECT
			tp.id,
			tp.tujuan_pemda,
			tp.tematik_id,
			tp.id_visi,
			tp.id_misi,
			tp.tahun_awal_periode,
			tp.tahun_akhir_periode,
			tp.jenis_periode,
			im_tg.indikator_id,
			im_tg.kode_indikator,
			im_tg.indikator,
			im_tg.rumus_perhitungan,
			im_tg.sumber_data,
			im_tg.definisi_operasional,
			im_tg.indikator_jenis,
			im_tg.target_id,
			im_tg.target_value,
			im_tg.satuan,
			im_tg.tahun_target,
			im_tg.target_jenis
		FROM tb_tujuan_pemda tp
		INNER JOIN tb_pohon_kinerja pk
			ON tp.tematik_id = pk.id
			AND pk.is_active = true
			AND pk.level_pohon = 0
		LEFT JOIN (
			SELECT
				im.id             AS indikator_id,
				im.kode_indikator,
				im.tujuan_pemda_id,
				im.indikator,
				im.rumus_perhitungan,
				im.sumber_data,
				im.definisi_operasional,
				COALESCE(im.jenis, 'renstra') AS indikator_jenis,
				tg.id             AS target_id,
				tg.target         AS target_value,
				tg.satuan,
				tg.tahun          AS tahun_target,
				tg.jenis          AS target_jenis
			FROM tb_indikator_matrix_pemda im
			LEFT JOIN tb_target_pemda tg
				ON im.kode_indikator = tg.kode_indikator
				AND tg.tahun = ?
				AND %s
			WHERE im.jenis = 'renstra' OR im.jenis = '' OR im.jenis IS NULL
		) im_tg ON tp.id = im_tg.tujuan_pemda_id
		WHERE tp.jenis_periode = ?
		  AND CAST(tp.tahun_awal_periode  AS SIGNED) <= CAST(? AS SIGNED)
		  AND CAST(tp.tahun_akhir_periode AS SIGNED) >= CAST(? AS SIGNED)
		ORDER BY tp.id, im_tg.indikator_id`, targetJenisClause)
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("FindAllByTahun error: %v", err)
	}
	defer rows.Close()
	tujuanMap := make(map[int]*domain.TujuanPemda)
	tujuanOrder := []int{}
	indikatorSeen := make(map[string]bool) // key: "tujuanId-indikatorId"
	for rows.Next() {
		var (
			tujuanId, tematikId, idVisi, idMisi                int
			tujuanText, tahunAwal, tahunAkhir, jenisPeriodeVal string
			indikatorId                                        sql.NullInt64
			kodeIndikator, indikatorText                       sql.NullString
			rumus, sumber, definisi, indikatorJenis            sql.NullString
			targetId                                           sql.NullInt64
			targetVal, satuan, tahunTarget, targetJenisVal     sql.NullString
		)
		if err := rows.Scan(
			&tujuanId, &tujuanText, &tematikId, &idVisi, &idMisi,
			&tahunAwal, &tahunAkhir, &jenisPeriodeVal,
			&indikatorId, &kodeIndikator, &indikatorText,
			&rumus, &sumber, &definisi, &indikatorJenis,
			&targetId, &targetVal, &satuan, &tahunTarget, &targetJenisVal,
		); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		if _, ok := tujuanMap[tujuanId]; !ok {
			tujuanMap[tujuanId] = &domain.TujuanPemda{
				Id:                tujuanId,
				TujuanPemda:       tujuanText,
				TematikId:         tematikId,
				IdVisi:            idVisi,
				IdMisi:            idMisi,
				TahunAwalPeriode:  tahunAwal,
				TahunAkhirPeriode: tahunAkhir,
				JenisPeriode:      jenisPeriodeVal,
				IndikatorPemda:    []domain.IndikatorPemda{},
			}
			tujuanOrder = append(tujuanOrder, tujuanId)
		}
		if !indikatorId.Valid {
			continue
		}
		indKey := fmt.Sprintf("%d-%d", tujuanId, indikatorId.Int64)
		if !indikatorSeen[indKey] {
			indikatorSeen[indKey] = true
			tujuanMap[tujuanId].IndikatorPemda = append(tujuanMap[tujuanId].IndikatorPemda, domain.IndikatorPemda{
				Id:                  int(indikatorId.Int64),
				KodeIndikator:       kodeIndikator.String,
				TujuanPemdaId:       tujuanId,
				Indikator:           indikatorText,
				RumusPerhitungan:    rumus,
				SumberData:          sumber,
				DefinisiOperasional: definisi,
				Jenis:               indikatorJenis.String,
				Target:              []domain.TargetPemda{},
			})
		}
		if targetId.Valid {
			tg := domain.TargetPemda{
				Id:            int(targetId.Int64),
				KodeIndikator: kodeIndikator.String,
				Target:        targetVal.String,
				Satuan:        satuan.String,
				Tahun:         tahunTarget.String,
				Jenis:         targetJenisVal.String,
			}
			for i := range tujuanMap[tujuanId].IndikatorPemda {
				if tujuanMap[tujuanId].IndikatorPemda[i].Id == int(indikatorId.Int64) {
					tujuanMap[tujuanId].IndikatorPemda[i].Target = append(
						tujuanMap[tujuanId].IndikatorPemda[i].Target, tg,
					)
					break
				}
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Pastikan setiap indikator punya 1 slot target untuk tahun yang diminta
	result := make([]domain.TujuanPemda, 0, len(tujuanOrder))
	for _, id := range tujuanOrder {
		tp := tujuanMap[id]
		for i := range tp.IndikatorPemda {
			if len(tp.IndikatorPemda[i].Target) == 0 {
				tp.IndikatorPemda[i].Target = []domain.TargetPemda{{
					Id:            0,
					KodeIndikator: tp.IndikatorPemda[i].KodeIndikator,
					Target:        "-",
					Satuan:        "-",
					Tahun:         tahun,
					Jenis:         targetJenis,
				}}
			}
		}
		result = append(result, *tp)
	}
	return result, nil
}
