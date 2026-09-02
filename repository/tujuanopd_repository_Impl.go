package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type TujuanOpdRepositoryImpl struct {
}

func NewTujuanOpdRepositoryImpl() *TujuanOpdRepositoryImpl {
	return &TujuanOpdRepositoryImpl{}
}

func (repository *TujuanOpdRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, tujuanOpd domain.TujuanOpd) (domain.TujuanOpd, error) {
	script := "INSERT INTO tb_tujuan_opd (kode_opd, kode_bidang_urusan, tujuan, periode_id, tahun_awal, tahun_akhir, jenis_periode) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, script,
		tujuanOpd.KodeOpd,
		tujuanOpd.KodeBidangUrusan,
		tujuanOpd.Tujuan,
		tujuanOpd.PeriodeId.Id,
		tujuanOpd.TahunAwal,
		tujuanOpd.TahunAkhir,
		tujuanOpd.JenisPeriode,
	)
	if err != nil {
		return domain.TujuanOpd{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return domain.TujuanOpd{}, err
	}
	tujuanOpd.Id = int(id)
	for _, indikator := range tujuanOpd.Indikator {
		scriptIndikator := `INSERT INTO tb_indikator_matrix
            (kode_indikator, tujuan_opd_id, indikator, rumus_perhitungan, sumber_data, definisi_operasional, jenis)
            VALUES (?, ?, ?, ?, ?, ?, ?)`
		_, err := tx.ExecContext(ctx, scriptIndikator,
			indikator.KodeIndikator,
			id,
			indikator.Indikator,
			indikator.RumusPerhitungan,
			indikator.SumberData,
			indikator.DefinisiOperasional,
			indikator.Jenis,
		)
		if err != nil {
			return domain.TujuanOpd{}, err
		}
		for _, target := range indikator.Target {
			_, err := tx.ExecContext(ctx,
				"INSERT INTO tb_target (id, indikator_id, target, satuan, tahun, jenis) VALUES (?, ?, ?, ?, ?, ?)",
				target.Id,
				indikator.KodeIndikator,
				target.Target,
				target.Satuan,
				target.Tahun,
				indikator.Jenis, // jenis target mengikuti jenis indikator (renstra)
			)
			if err != nil {
				return domain.TujuanOpd{}, err
			}
		}
	}
	return tujuanOpd, nil
}

func (repository *TujuanOpdRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, tujuanOpd domain.TujuanOpd) error {
	// 1. UPDATE header saja — TIDAK update tahun/periode
	_, err := tx.ExecContext(ctx,
		`UPDATE tb_tujuan_opd SET kode_opd=?, kode_bidang_urusan=?, tujuan=? WHERE id=?`,
		tujuanOpd.KodeOpd, tujuanOpd.KodeBidangUrusan, tujuanOpd.Tujuan, tujuanOpd.Id)
	if err != nil {
		return err
	}
	keepIndIds := make([]string, 0)
	for _, ind := range tujuanOpd.Indikator {
		kodeInd := ind.KodeIndikator
		if ind.Id != "" {
			// UPDATE indikator existing — kode_indikator TIDAK berubah
			_, err = tx.ExecContext(ctx, `
                UPDATE tb_indikator_matrix
                SET indikator=?, rumus_perhitungan=?, sumber_data=?,
                    definisi_operasional=?, jenis=?
                WHERE id=? AND tujuan_opd_id=?`,
				ind.Indikator, ind.RumusPerhitungan, ind.SumberData,
				ind.DefinisiOperasional, ind.Jenis, ind.Id, tujuanOpd.Id)
			if err != nil {
				return err
			}
			keepIndIds = append(keepIndIds, ind.Id)
		} else {
			// INSERT indikator baru
			res, err := tx.ExecContext(ctx, `
                INSERT INTO tb_indikator_matrix
                    (kode_indikator, tujuan_opd_id, indikator, rumus_perhitungan,
                     sumber_data, definisi_operasional, jenis)
                VALUES (?, ?, ?, ?, ?, ?, ?)`,
				kodeInd, tujuanOpd.Id, ind.Indikator, ind.RumusPerhitungan,
				ind.SumberData, ind.DefinisiOperasional, ind.Jenis)
			if err != nil {
				return err
			}
			newId, _ := res.LastInsertId()
			ind.Id = strconv.FormatInt(newId, 10)
			keepIndIds = append(keepIndIds, ind.Id)
		}
		keepTargetIds := make([]string, 0)
		for _, t := range ind.Target {
			if t.Id != "" {
				// UPDATE target existing — hanya renstra
				_, err = tx.ExecContext(ctx, `
                    UPDATE tb_target SET target=?, satuan=?, tahun=?
                    WHERE id=? AND indikator_id=? AND jenis='renstra'`,
					t.Target, t.Satuan, t.Tahun, t.Id, kodeInd)
				if err != nil {
					return err
				}
				keepTargetIds = append(keepTargetIds, t.Id)
			} else {
				// INSERT target baru
				uuidTrg := uuid.New().String()[:5]
				newId := fmt.Sprintf("TRG-TJN-%s", uuidTrg)
				_, err = tx.ExecContext(ctx, `
                    INSERT INTO tb_target (id, indikator_id, target, satuan, tahun, jenis)
                    VALUES (?, ?, ?, ?, ?, 'renstra')`,
					newId, kodeInd, t.Target, t.Satuan, t.Tahun)
				if err != nil {
					return err
				}
				keepTargetIds = append(keepTargetIds, newId)
			}
		}
		// Hapus target renstra yang tidak ada di request
		if len(keepTargetIds) > 0 {
			ph := strings.Repeat("?,", len(keepTargetIds))
			ph = ph[:len(ph)-1]
			args := []interface{}{kodeInd}
			for _, id := range keepTargetIds {
				args = append(args, id)
			}
			_, err = tx.ExecContext(ctx,
				fmt.Sprintf(`DELETE FROM tb_target
                             WHERE indikator_id=? AND jenis='renstra'
                               AND id NOT IN (%s)`, ph), args...)
			if err != nil {
				return err
			}
		}
	}
	// Hapus indikator renstra yang tidak ada di request
	if len(keepIndIds) > 0 {
		ph := strings.Repeat("?,", len(keepIndIds))
		ph = ph[:len(ph)-1]
		args := []interface{}{tujuanOpd.Id}
		for _, id := range keepIndIds {
			args = append(args, id)
		}
		// Hapus orphan target dulu
		_, err = tx.ExecContext(ctx,
			fmt.Sprintf(`DELETE t FROM tb_target t
                         INNER JOIN tb_indikator_matrix i ON t.indikator_id = i.kode_indikator
                         WHERE i.tujuan_opd_id=? AND i.jenis='renstra'
                           AND i.id NOT IN (%s) AND t.jenis='renstra'`, ph), args...)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx,
			fmt.Sprintf(`DELETE FROM tb_indikator_matrix
                         WHERE tujuan_opd_id=? AND jenis='renstra'
                           AND id NOT IN (%s)`, ph), args...)
		if err != nil {
			return err
		}
	}
	return nil
}
func (repository *TujuanOpdRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, tujuanOpdId int) error {
	// 1. Hapus target terkait via tb_indikator_matrix
	_, err := tx.ExecContext(ctx, `
		DELETE tg FROM tb_target tg
		INNER JOIN tb_indikator_matrix im ON tg.indikator_id = im.kode_indikator
		WHERE im.tujuan_opd_id = ?
	`, tujuanOpdId)
	if err != nil {
		return err
	}
	// 2. Hapus indikator_matrix
	_, err = tx.ExecContext(ctx,
		"DELETE FROM tb_indikator_matrix WHERE tujuan_opd_id = ?", tujuanOpdId)
	if err != nil {
		return err
	}
	// 3. Hapus tujuan_opd
	_, err = tx.ExecContext(ctx,
		"DELETE FROM tb_tujuan_opd WHERE id = ?", tujuanOpdId)
	return err
}

func (repository *TujuanOpdRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, tujuanOpdId int) (domain.TujuanOpd, error) {
	script := `
		SELECT
			t.id,
			t.kode_opd,
			COALESCE(t.kode_bidang_urusan, '') AS kode_bidang_urusan,
			t.tujuan,
			t.tahun_awal,
			t.tahun_akhir,
			t.jenis_periode,
			im.id                              AS indikator_id,
			im.kode_indikator,
			COALESCE(im.indikator, '')         AS indikator_nama,
			COALESCE(im.rumus_perhitungan, '') AS rumus_perhitungan,
			COALESCE(im.definisi_operasional,'') AS definisi_operasional,
			COALESCE(im.sumber_data, '')       AS sumber_data,
			COALESCE(im.jenis, '')             AS indikator_jenis,
			COALESCE(tg.id, '')                AS target_id,
			COALESCE(tg.target, '')            AS target_value,
			COALESCE(tg.satuan, '')            AS satuan,
			COALESCE(tg.tahun, '')             AS tahun_target
		FROM tb_tujuan_opd t
		LEFT JOIN tb_indikator_matrix im
			ON t.id = im.tujuan_opd_id AND im.jenis = 'renstra'
		LEFT JOIN tb_target tg
			ON im.kode_indikator = tg.indikator_id
			AND (tg.jenis = 'renstra' OR tg.jenis = '')
		WHERE t.id = ?
		ORDER BY im.id ASC, CAST(tg.tahun AS SIGNED) ASC
	`
	rows, err := tx.QueryContext(ctx, script, tujuanOpdId)
	if err != nil {
		return domain.TujuanOpd{}, err
	}
	defer rows.Close()
	var tujuanOpd *domain.TujuanOpd
	indikatorMap := make(map[string]*domain.Indikator)
	indikatorOrder := []string{}
	for rows.Next() {
		var (
			id                  int
			kodeOpd             string
			kodeBidangUrusan    string
			tujuan              string
			tahunAwal           string
			tahunAkhir          string
			jenisPeriode        string
			indikatorId         sql.NullString
			kodeIndikator       sql.NullString
			indikatorNama       sql.NullString
			rumusPerhitungan    sql.NullString
			definisiOperasional sql.NullString
			sumberData          sql.NullString
			indikatorJenis      sql.NullString
			targetId            sql.NullString
			targetValue         sql.NullString
			satuan              sql.NullString
			tahunTarget         sql.NullString
		)
		err := rows.Scan(
			&id, &kodeOpd, &kodeBidangUrusan, &tujuan,
			&tahunAwal, &tahunAkhir, &jenisPeriode,
			&indikatorId, &kodeIndikator, &indikatorNama,
			&rumusPerhitungan, &definisiOperasional, &sumberData, &indikatorJenis,
			&targetId, &targetValue, &satuan, &tahunTarget,
		)
		if err != nil {
			return domain.TujuanOpd{}, err
		}
		if tujuanOpd == nil {
			tujuanOpd = &domain.TujuanOpd{
				Id:               id,
				KodeOpd:          kodeOpd,
				KodeBidangUrusan: kodeBidangUrusan,
				Tujuan:           tujuan,
				TahunAwal:        tahunAwal,
				TahunAkhir:       tahunAkhir,
				JenisPeriode:     jenisPeriode,
				Indikator:        []domain.Indikator{},
			}
		}
		if indikatorId.Valid {
			if _, exists := indikatorMap[indikatorId.String]; !exists {
				indikatorMap[indikatorId.String] = &domain.Indikator{
					Id:                  indikatorId.String,
					KodeIndikator:       kodeIndikator.String,
					Indikator:           indikatorNama.String,
					RumusPerhitungan:    rumusPerhitungan,
					SumberData:          sumberData,
					DefinisiOperasional: definisiOperasional,
					Jenis:               indikatorJenis.String,
					Target:              []domain.Target{},
				}
				indikatorOrder = append(indikatorOrder, indikatorId.String)
			}
			if targetId.Valid && targetId.String != "" && tahunTarget.Valid {
				target := domain.Target{
					Id:          targetId.String,
					IndikatorId: kodeIndikator.String,
					Target:      targetValue.String,
					Satuan:      satuan.String,
					Tahun:       tahunTarget.String,
				}
				indikatorMap[indikatorId.String].Target = append(
					indikatorMap[indikatorId.String].Target, target,
				)
			}
		}
	}
	if err := rows.Err(); err != nil {
		return domain.TujuanOpd{}, err
	}
	if tujuanOpd == nil {
		return domain.TujuanOpd{}, fmt.Errorf("tujuan opd dengan id %d tidak ditemukan", tujuanOpdId)
	}
	// Generate slot target lengkap per indikator
	for _, indId := range indikatorOrder {
		ind := indikatorMap[indId]
		tahunAwalInt, _ := strconv.Atoi(tujuanOpd.TahunAwal)
		tahunAkhirInt, _ := strconv.Atoi(tujuanOpd.TahunAkhir)
		existingTargets := make(map[string]domain.Target)
		for _, t := range ind.Target {
			existingTargets[t.Tahun] = t
		}
		var completeTargets []domain.Target
		for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
			tahunStr := strconv.Itoa(tahun)
			if t, ok := existingTargets[tahunStr]; ok {
				completeTargets = append(completeTargets, t)
			} else {
				completeTargets = append(completeTargets, domain.Target{
					Id: "", IndikatorId: ind.KodeIndikator,
					Target: "", Satuan: "", Tahun: tahunStr,
				})
			}
		}
		ind.Target = completeTargets
		tujuanOpd.Indikator = append(tujuanOpd.Indikator, *ind)
	}
	return *tujuanOpd, nil
}

func (repository *TujuanOpdRepositoryImpl) FindIndikatorByTujuanId(ctx context.Context, tx *sql.Tx, tujuanOpdId int) ([]domain.Indikator, error) {
	script := `SELECT id, indikator
               FROM tb_indikator
               WHERE tujuan_opd_id = ?`

	rows, err := tx.QueryContext(ctx, script, tujuanOpdId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indikators []domain.Indikator
	for rows.Next() {
		var indikator domain.Indikator
		err := rows.Scan(&indikator.Id, &indikator.Indikator)
		if err != nil {
			return nil, err
		}
		indikators = append(indikators, indikator)
	}

	return indikators, nil
}

func (repository *TujuanOpdRepositoryImpl) FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string, tahun string) ([]domain.Target, error) {
	script := `
        SELECT id, target, satuan, tahun
        FROM tb_target
        WHERE indikator_id = ?
        AND tahun = ?
        ORDER BY tahun ASC
    `

	rows, err := tx.QueryContext(ctx, script, indikatorId, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []domain.Target
	for rows.Next() {
		var target domain.Target
		err := rows.Scan(
			&target.Id,
			&target.Target,
			&target.Satuan,
			&target.Tahun,
		)
		if err != nil {
			return nil, err
		}
		target.IndikatorId = indikatorId
		targets = append(targets, target)
	}

	return targets, nil
}

func (repository *TujuanOpdRepositoryImpl) FindAll(
	ctx context.Context, tx *sql.Tx,
	kodeOpd, tahunAwal, tahunAkhir, jenisPeriode string,
) ([]domain.TujuanOpd, error) {
	// Delegasi ke FindAllByPeriod dengan jenis='renstra' — sudah benar
	return repository.FindAllByPeriod(ctx, tx, kodeOpd, tahunAwal, tahunAkhir, jenisPeriode, "renstra")
}

func (repository *TujuanOpdRepositoryImpl) FindTujuanOpdByTahun(
	ctx context.Context, tx *sql.Tx,
	kodeOpd, tahun, jenisPeriode string,
) ([]domain.TujuanOpd, error) {
	// Delegasi ke FindAllByTahun dengan jenis='renstra' sebagai base
	return repository.FindAllByTahun(ctx, tx, kodeOpd, tahun, jenisPeriode, "renstra")
}
func (repository *TujuanOpdRepositoryImpl) FindTujuanOpdByTahunByStrategicArahKebijakan(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string, jenisPeriode string) ([]domain.TujuanOpd, error) {
	scriptTujuan := `
		SELECT
			MIN(t.id) as id,
			t.kode_opd,
			t.tujuan,
			MIN(t.tahun_awal) as tahun_awal,
			MAX(t.tahun_akhir) as tahun_akhir
		FROM tb_tujuan_opd t
		WHERE t.kode_opd = ?
			AND CAST(? AS SIGNED) BETWEEN CAST(t.tahun_awal AS SIGNED) AND CAST(t.tahun_akhir AS SIGNED)
			AND t.jenis_periode = ?
		GROUP BY
			t.kode_opd,
			t.tujuan
		ORDER BY id ASC;
	`

	rows, err := tx.QueryContext(ctx,
		scriptTujuan,
		kodeOpd,
		tahun,
		jenisPeriode,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.TujuanOpd

	for rows.Next() {
		var tujuan domain.TujuanOpd

		err := rows.Scan(
			&tujuan.Id,
			&tujuan.KodeOpd,
			&tujuan.Tujuan,
			&tujuan.TahunAwal,
			&tujuan.TahunAkhir,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, tujuan)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return []domain.TujuanOpd{}, nil
	}

	return result, nil
}

// Perbaikan pada FindIndikatorByTujuanOpdId untuk menyertakan rumus_perhitungan dan sumber_data
func (repository *TujuanOpdRepositoryImpl) FindIndikatorByTujuanOpdId(
	ctx context.Context, tx *sql.Tx, tujuanOpdId int,
) ([]domain.Indikator, error) {
	script := `
		SELECT
			id,
			kode_indikator,
			COALESCE(indikator, '')             AS indikator,
			COALESCE(rumus_perhitungan, '')     AS rumus_perhitungan,
			COALESCE(sumber_data, '')           AS sumber_data,
			COALESCE(definisi_operasional, '')  AS definisi_operasional,
			COALESCE(jenis, '')                 AS jenis
		FROM tb_indikator_matrix
		WHERE tujuan_opd_id = ?
		AND jenis = 'renstra'
		ORDER BY id ASC
	`
	rows, err := tx.QueryContext(ctx, script, tujuanOpdId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var indikators []domain.Indikator
	for rows.Next() {
		var (
			ind                 domain.Indikator
			rumusPerhitungan    string
			sumberData          string
			definisiOperasional string
			jenis               string
		)
		err := rows.Scan(
			&ind.Id, &ind.KodeIndikator, &ind.Indikator,
			&rumusPerhitungan, &sumberData, &definisiOperasional, &jenis,
		)
		if err != nil {
			return nil, err
		}
		ind.TujuanOpdId = tujuanOpdId
		ind.RumusPerhitungan = sql.NullString{String: rumusPerhitungan, Valid: rumusPerhitungan != ""}
		ind.SumberData = sql.NullString{String: sumberData, Valid: sumberData != ""}
		ind.DefinisiOperasional = sql.NullString{String: definisiOperasional, Valid: definisiOperasional != ""}
		ind.Jenis = jenis
		indikators = append(indikators, ind)
	}
	return indikators, rows.Err()
}

func (repository *TujuanOpdRepositoryImpl) FindTujuanOpdForCascadingOpd(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string, jenisPeriode string) ([]domain.TujuanOpd, error) {
	script := `
		SELECT
			t.id,
			t.kode_opd,
			t.tujuan,
			t.tahun_awal,
			t.tahun_akhir,
			t.jenis_periode,
			t.kode_bidang_urusan
		FROM tb_tujuan_opd t
		INNER JOIN tb_periode p ON
			t.tahun_awal = p.tahun_awal
			AND t.tahun_akhir = p.tahun_akhir
			AND t.jenis_periode = p.jenis_periode
		WHERE t.kode_opd = ?
		AND CAST(? AS SIGNED) BETWEEN CAST(p.tahun_awal AS SIGNED) AND CAST(p.tahun_akhir AS SIGNED)
		AND p.jenis_periode = ?
		ORDER BY t.id ASC
    `

	rows, err := tx.QueryContext(ctx, script, kodeOpd, tahun, jenisPeriode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tujuanOpds []domain.TujuanOpd
	for rows.Next() {
		var tujuanOpd domain.TujuanOpd
		err := rows.Scan(
			&tujuanOpd.Id,
			&tujuanOpd.KodeOpd,
			&tujuanOpd.Tujuan,
			&tujuanOpd.TahunAwal,
			&tujuanOpd.TahunAkhir,
			&tujuanOpd.JenisPeriode,
			&tujuanOpd.KodeBidangUrusan,
		)
		if err != nil {
			return nil, err
		}
		tujuanOpds = append(tujuanOpds, tujuanOpd)
	}

	// Untuk setiap tujuan OPD, ambil indikatornya
	for i := range tujuanOpds {
		indikators, err := repository.FindIndikatorByTujuanOpdId(ctx, tx, tujuanOpds[i].Id)
		if err != nil {
			return nil, err
		}
		tujuanOpds[i].Indikator = indikators

		// Generate target untuk setiap indikator
		for j := range tujuanOpds[i].Indikator {
			tahunAwalInt, _ := strconv.Atoi(tujuanOpds[i].TahunAwal)
			tahunAkhirInt, _ := strconv.Atoi(tujuanOpds[i].TahunAkhir)

			var targets []domain.Target
			for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
				targets = append(targets, domain.Target{
					Id:          "",
					IndikatorId: tujuanOpds[i].Indikator[j].Id,
					Target:      "",
					Satuan:      "",
					Tahun:       strconv.Itoa(tahun),
				})
			}
			tujuanOpds[i].Indikator[j].Target = targets
		}
	}

	return tujuanOpds, nil
}

func (repository *TujuanOpdRepositoryImpl) FindIndikatorByTujuanOpdIdsBatch(
	ctx context.Context, tx *sql.Tx, tujuanOpdIds []int,
) (map[int][]domain.Indikator, error) {
	if len(tujuanOpdIds) == 0 {
		return make(map[int][]domain.Indikator), nil
	}
	placeholders := make([]string, len(tujuanOpdIds))
	args := make([]interface{}, len(tujuanOpdIds))
	for i, id := range tujuanOpdIds {
		placeholders[i] = "?"
		args[i] = id
	}
	script := fmt.Sprintf(`
		SELECT
			id,
			kode_indikator,
			tujuan_opd_id,
			COALESCE(indikator, '')             AS indikator,
			COALESCE(rumus_perhitungan, '')     AS rumus_perhitungan,
			COALESCE(sumber_data, '')           AS sumber_data,
			COALESCE(definisi_operasional, '')  AS definisi_operasional,
			COALESCE(jenis, '')                 AS jenis
		FROM tb_indikator_matrix
		WHERE tujuan_opd_id IN (%s)
		AND jenis = 'renstra'
		ORDER BY tujuan_opd_id ASC, id ASC
	`, strings.Join(placeholders, ","))
	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[int][]domain.Indikator)
	for rows.Next() {
		var (
			ind                 domain.Indikator
			tujuanOpdId         int
			rumusPerhitungan    string
			sumberData          string
			definisiOperasional string
			jenis               string
		)
		err := rows.Scan(
			&ind.Id, &ind.KodeIndikator, &tujuanOpdId, &ind.Indikator,
			&rumusPerhitungan, &sumberData, &definisiOperasional, &jenis,
		)
		if err != nil {
			return nil, err
		}
		ind.RumusPerhitungan = sql.NullString{String: rumusPerhitungan, Valid: rumusPerhitungan != ""}
		ind.SumberData = sql.NullString{String: sumberData, Valid: sumberData != ""}
		ind.DefinisiOperasional = sql.NullString{String: definisiOperasional, Valid: definisiOperasional != ""}
		ind.Jenis = jenis
		result[tujuanOpdId] = append(result[tujuanOpdId], ind)
	}
	return result, rows.Err()
}

// renstra new
func (repository *TujuanOpdRepositoryImpl) FindAllByPeriod(
	ctx context.Context, tx *sql.Tx,
	kodeOpd, tahunAwal, tahunAkhir, jenisPeriode, jenisIndikator string,
) ([]domain.TujuanOpd, error) {
	jenisClause := ""
	var finalArgs []interface{}
	if jenisIndikator != "" {
		jenisClause = "AND im.jenis = ?"
		finalArgs = append(finalArgs, jenisIndikator) // 1: im.jenis = ?
	}
	finalArgs = append(finalArgs, tahunAwal, tahunAkhir)                        // 2,3: BETWEEN target
	finalArgs = append(finalArgs, kodeOpd, tahunAwal, tahunAkhir, jenisPeriode) // 4,5,6,7: WHERE tujuan
	query := fmt.Sprintf(`
        SELECT
            t.id,
            t.kode_opd,
            COALESCE(t.kode_bidang_urusan, '')         AS kode_bidang_urusan,
            t.tujuan,
            t.tahun_awal,
            t.tahun_akhir,
            t.jenis_periode,
            im.id                                      AS indikator_id,
            im.kode_indikator                           AS kode_indikator,
            im.indikator,
            COALESCE(im.rumus_perhitungan, '')         AS rumus_perhitungan,
            COALESCE(im.sumber_data, '')               AS sumber_data,
            COALESCE(im.definisi_operasional, '')      AS definisi_operasional,
            COALESCE(im.jenis, '')                     AS indikator_jenis,
            tg.id                                      AS target_id,
            tg.target                                  AS target_value,
            tg.satuan,
            tg.tahun                                   AS tahun_target
        FROM tb_tujuan_opd t
        LEFT JOIN tb_indikator_matrix im
            ON t.id = im.tujuan_opd_id %s
        LEFT JOIN tb_target tg
            ON im.kode_indikator = tg.indikator_id
            AND CAST(tg.tahun AS SIGNED) BETWEEN CAST(? AS SIGNED) AND CAST(? AS SIGNED)
        WHERE t.kode_opd      = ?
          AND t.tahun_awal    = ?
          AND t.tahun_akhir   = ?
          AND t.jenis_periode = ?
        ORDER BY t.id, im.id, tg.tahun
    `, jenisClause)
	rows, err := tx.QueryContext(ctx, query, finalArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tujuanOpdMap := make(map[int]*domain.TujuanOpd)
	indikatorSeen := make(map[string]bool) // key: "tujuanId-indikatorId"
	tujuanOrder := []int{}
	for rows.Next() {
		var (
			tujuanId            int
			kodeOpdData         string
			kodeBidangUrusan    string
			tujuan              string
			tahunAwalData       string
			tahunAkhirData      string
			jenisPeriodeData    string
			indikatorId         sql.NullString
			kodeIndikator       sql.NullString
			indikatorNama       sql.NullString
			rumusPerhitungan    sql.NullString
			sumberData          sql.NullString
			definisiOperasional sql.NullString // NEW
			indikatorJenis      sql.NullString // im.jenis
			targetId            sql.NullString
			targetValue         sql.NullString
			satuan              sql.NullString
			tahunTarget         sql.NullString
		)
		err := rows.Scan(
			&tujuanId,
			&kodeOpdData,
			&kodeBidangUrusan,
			&tujuan,
			&tahunAwalData,
			&tahunAkhirData,
			&jenisPeriodeData,
			&indikatorId,
			&kodeIndikator,
			&indikatorNama,
			&rumusPerhitungan,
			&sumberData,
			&definisiOperasional, // NEW
			&indikatorJenis,
			&targetId,
			&targetValue,
			&satuan,
			&tahunTarget,
		)
		if err != nil {
			return nil, err
		}
		if _, exists := tujuanOpdMap[tujuanId]; !exists {
			tujuanOpdMap[tujuanId] = &domain.TujuanOpd{
				Id:               tujuanId,
				KodeOpd:          kodeOpdData,
				KodeBidangUrusan: kodeBidangUrusan,
				Tujuan:           tujuan,
				TahunAwal:        tahunAwalData,
				TahunAkhir:       tahunAkhirData,
				JenisPeriode:     jenisPeriodeData,
				Indikator:        []domain.Indikator{},
			}
			tujuanOrder = append(tujuanOrder, tujuanId)
		}
		if indikatorId.Valid {
			mapKey := fmt.Sprintf("%d-%s", tujuanId, indikatorId.String)
			if !indikatorSeen[mapKey] {
				indikatorSeen[mapKey] = true
				tujuanOpdMap[tujuanId].Indikator = append(tujuanOpdMap[tujuanId].Indikator, domain.Indikator{
					Id:                  indikatorId.String,
					KodeIndikator:       kodeIndikator.String,
					Indikator:           indikatorNama.String,
					RumusPerhitungan:    rumusPerhitungan,
					SumberData:          sumberData,
					DefinisiOperasional: definisiOperasional, // NEW
					Jenis:               indikatorJenis.String,
					TujuanOpdId:         tujuanId,
					Target:              []domain.Target{},
				})
			}
			if targetId.Valid && tahunTarget.Valid {
				target := domain.Target{
					Id:          targetId.String,
					IndikatorId: kodeIndikator.String,
					Target:      targetValue.String,
					Satuan:      satuan.String,
					Tahun:       tahunTarget.String,
				}
				for idx := range tujuanOpdMap[tujuanId].Indikator {
					if tujuanOpdMap[tujuanId].Indikator[idx].Id == indikatorId.String {
						tujuanOpdMap[tujuanId].Indikator[idx].Target = append(
							tujuanOpdMap[tujuanId].Indikator[idx].Target, target,
						)
						break
					}
				}
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Renstra: generate slot target untuk setiap tahun dalam range
	var result []domain.TujuanOpd
	for _, id := range tujuanOrder {
		tujuanOpd := tujuanOpdMap[id]
		for i := range tujuanOpd.Indikator {
			tahunAwalInt, _ := strconv.Atoi(tujuanOpd.TahunAwal)
			tahunAkhirInt, _ := strconv.Atoi(tujuanOpd.TahunAkhir)
			existingTargets := make(map[string]domain.Target)
			for _, t := range tujuanOpd.Indikator[i].Target {
				if t.Id != "" {
					existingTargets[t.Tahun] = t
				}
			}
			var completeTargets []domain.Target
			for year := tahunAwalInt; year <= tahunAkhirInt; year++ {
				yearStr := strconv.Itoa(year)
				if t, ok := existingTargets[yearStr]; ok {
					completeTargets = append(completeTargets, t)
				} else {
					completeTargets = append(completeTargets, domain.Target{
						Id:          "",
						IndikatorId: tujuanOpd.Indikator[i].KodeIndikator,
						Target:      "",
						Satuan:      "",
						Tahun:       yearStr,
					})
				}
			}
			tujuanOpd.Indikator[i].Target = completeTargets
		}
		result = append(result, *tujuanOpd)
	}
	if len(result) == 0 {
		return make([]domain.TujuanOpd, 0), nil
	}
	return result, nil
}

// FindAllByTahun mengambil tujuan OPD dengan indikator default dari renstra,
// dan target sesuai jenisIndikator (renstra/ranwal/rankhir/penetapan) untuk tahun yang diminta.
// Metadata indikator (nama, rumus, sumber data, dll) selalu bersumber dari layer renstra.
// Catatan: untuk jenis="renstra", data lama dengan jenis=” juga diambil (backward compat).
func (repository *TujuanOpdRepositoryImpl) FindAllByTahun(
	ctx context.Context, tx *sql.Tx,
	kodeOpd, tahun, jenisPeriode, jenisIndikator string,
) ([]domain.TujuanOpd, error) {
	// Untuk layer renstra, ambil juga target dengan jenis='' (data sebelum migrasi).
	var (
		targetJenisClause string
		finalArgs         []interface{}
	)
	if jenisIndikator == "renstra" {
		targetJenisClause = "(tg.jenis = 'renstra' OR tg.jenis = '')"
		finalArgs = []interface{}{tahun, kodeOpd, jenisPeriode, tahun, tahun}
	} else {
		targetJenisClause = "tg.jenis = ?"
		finalArgs = []interface{}{tahun, jenisIndikator, kodeOpd, jenisPeriode, tahun, tahun}
	}

	query := fmt.Sprintf(`
    SELECT
        t.id,
        t.kode_opd,
        COALESCE(t.kode_bidang_urusan, '')      AS kode_bidang_urusan,
        t.tujuan,
        t.tahun_awal,
        t.tahun_akhir,
        t.jenis_periode,
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
        im_tg.tahun_target
    FROM tb_tujuan_opd t
    LEFT JOIN (
        SELECT
            im.id                                  AS indikator_id,
            im.kode_indikator                      AS kode_indikator,
            im.tujuan_opd_id,
            COALESCE(im.indikator, '')             AS indikator,
            COALESCE(im.rumus_perhitungan, '')     AS rumus_perhitungan,
            COALESCE(im.sumber_data, '')           AS sumber_data,
            COALESCE(im.definisi_operasional, '')  AS definisi_operasional,
            COALESCE(im.jenis, '')                 AS indikator_jenis,
            tg.id                                  AS target_id,
            tg.target                              AS target_value,
            tg.satuan,
            tg.tahun                               AS tahun_target
        FROM tb_indikator_matrix im
        LEFT JOIN tb_target tg
            ON im.kode_indikator = tg.indikator_id
            AND tg.tahun = ?
            AND %s
        WHERE im.jenis = 'renstra'
    ) im_tg ON t.id = im_tg.tujuan_opd_id
    WHERE t.kode_opd      = ?
      AND t.jenis_periode  = ?
      AND CAST(t.tahun_awal  AS SIGNED) <= CAST(? AS SIGNED)
      AND CAST(t.tahun_akhir AS SIGNED) >= CAST(? AS SIGNED)
    ORDER BY t.id, im_tg.indikator_id`, targetJenisClause)
	rows, err := tx.QueryContext(ctx, query, finalArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tujuanOpdMap := make(map[int]*domain.TujuanOpd)
	indikatorSeen := make(map[string]bool) // key: "tujuanId-indikatorId"
	tujuanOrder := []int{}
	for rows.Next() {
		var (
			tujuanId            int
			kodeOpdData         string
			kodeBidangUrusan    string
			tujuan              string
			tahunAwalData       string
			tahunAkhirData      string
			jenisPeriodeData    string
			indikatorId         sql.NullString
			kodeIndikator       sql.NullString
			indikatorNama       sql.NullString
			rumusPerhitungan    sql.NullString
			sumberData          sql.NullString
			definisiOperasional sql.NullString // NEW
			indikatorJenis      sql.NullString // im.jenis
			targetId            sql.NullString
			targetValue         sql.NullString
			satuan              sql.NullString
			tahunTarget         sql.NullString
		)
		err := rows.Scan(
			&tujuanId,
			&kodeOpdData,
			&kodeBidangUrusan,
			&tujuan,
			&tahunAwalData,
			&tahunAkhirData,
			&jenisPeriodeData,
			&indikatorId,
			&kodeIndikator,
			&indikatorNama,
			&rumusPerhitungan,
			&sumberData,
			&definisiOperasional, // NEW
			&indikatorJenis,
			&targetId,
			&targetValue,
			&satuan,
			&tahunTarget,
		)
		if err != nil {
			return nil, err
		}
		if _, exists := tujuanOpdMap[tujuanId]; !exists {
			tujuanOpdMap[tujuanId] = &domain.TujuanOpd{
				Id:               tujuanId,
				KodeOpd:          kodeOpdData,
				KodeBidangUrusan: kodeBidangUrusan,
				Tujuan:           tujuan,
				TahunAwal:        tahunAwalData,
				TahunAkhir:       tahunAkhirData,
				JenisPeriode:     jenisPeriodeData,
				Indikator:        []domain.Indikator{},
			}
			tujuanOrder = append(tujuanOrder, tujuanId)
		}
		if indikatorId.Valid {
			mapKey := fmt.Sprintf("%d-%s", tujuanId, indikatorId.String)
			if !indikatorSeen[mapKey] {
				indikatorSeen[mapKey] = true
				tujuanOpdMap[tujuanId].Indikator = append(tujuanOpdMap[tujuanId].Indikator, domain.Indikator{
					Id:                  indikatorId.String,
					KodeIndikator:       kodeIndikator.String,
					Indikator:           indikatorNama.String,
					RumusPerhitungan:    rumusPerhitungan,
					SumberData:          sumberData,
					DefinisiOperasional: definisiOperasional, // NEW
					Jenis:               indikatorJenis.String,
					TujuanOpdId:         tujuanId,
					Target:              []domain.Target{},
				})
			}
			if targetId.Valid {
				target := domain.Target{
					Id:          targetId.String,
					IndikatorId: kodeIndikator.String,
					Target:      targetValue.String,
					Satuan:      satuan.String,
					Tahun:       tahunTarget.String,
				}
				for idx := range tujuanOpdMap[tujuanId].Indikator {
					if tujuanOpdMap[tujuanId].Indikator[idx].Id == indikatorId.String {
						tujuanOpdMap[tujuanId].Indikator[idx].Target = append(
							tujuanOpdMap[tujuanId].Indikator[idx].Target, target,
						)
						break
					}
				}
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Ranwal / Rankhir: tepat 1 slot target per indikator untuk tahun yang diminta
	var result []domain.TujuanOpd
	for _, id := range tujuanOrder {
		tujuanOpd := tujuanOpdMap[id]
		for i := range tujuanOpd.Indikator {
			if len(tujuanOpd.Indikator[i].Target) == 0 {
				// Belum ada data target di DB → buat slot kosong
				tujuanOpd.Indikator[i].Target = []domain.Target{{
					Id:          "",
					IndikatorId: tujuanOpd.Indikator[i].KodeIndikator,
					Target:      "",
					Satuan:      "",
					Tahun:       tahun,
				}}
			}
		}
		result = append(result, *tujuanOpd)
	}
	if len(result) == 0 {
		return make([]domain.TujuanOpd, 0), nil
	}
	return result, nil
}

func (repository *TujuanOpdRepositoryImpl) FindBidangUrusanBatch(
	ctx context.Context, tx *sql.Tx,
	kodeBidangUrusanList []string,
) (map[string]domainmaster.BidangUrusan, error) {
	result := make(map[string]domainmaster.BidangUrusan)
	if len(kodeBidangUrusanList) == 0 {
		return result, nil
	}
	// Deduplicate & filter kosong
	uniqueSet := make(map[string]struct{})
	for _, k := range kodeBidangUrusanList {
		if k != "" {
			uniqueSet[k] = struct{}{}
		}
	}
	if len(uniqueSet) == 0 {
		return result, nil
	}
	placeholders := make([]string, 0, len(uniqueSet))
	args := make([]interface{}, 0, len(uniqueSet))
	for k := range uniqueSet {
		placeholders = append(placeholders, "?")
		args = append(args, k)
	}
	// JOIN ke tb_urusan menggunakan digit pertama kode_bidang_urusan
	// Sama persis dengan pola FindByKodeBidangUrusan yang sudah jalan
	query := fmt.Sprintf(`
        SELECT
            COALESCE(bu.id, '')                  AS id,
            COALESCE(bu.kode_bidang_urusan, '')   AS kode_bidang_urusan,
            COALESCE(bu.nama_bidang_urusan, '')   AS nama_bidang_urusan,
            COALESCE(u.kode_urusan, '')           AS kode_urusan,
            COALESCE(u.nama_urusan, '')           AS nama_urusan
        FROM tb_bidang_urusan bu
        LEFT JOIN tb_urusan u
            ON LEFT(bu.kode_bidang_urusan, 1) = u.kode_urusan
        WHERE bu.kode_bidang_urusan IN (%s)
    `, strings.Join(placeholders, ","))
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var bu domainmaster.BidangUrusan
		if err := rows.Scan(
			&bu.Id,
			&bu.KodeBidangUrusan,
			&bu.NamaBidangUrusan,
			&bu.KodeUrusan,
			&bu.NamaUrusan,
		); err != nil {
			return nil, err
		}
		result[bu.KodeBidangUrusan] = bu
	}
	return result, rows.Err()
}

func (r *TujuanOpdRepositoryImpl) CreateRenjaIndikator(
	ctx context.Context, tx *sql.Tx,
	tujuanOpdId int, indikators []domain.Indikator,
) error {
	for _, ind := range indikators {
		_, err := tx.ExecContext(ctx, `
            INSERT INTO tb_indikator_matrix
                (kode_indikator, tujuan_opd_id, indikator, rumus_perhitungan,
                 sumber_data, definisi_operasional, jenis)
            VALUES (?, ?, ?, ?, ?, ?, ?)`,
			ind.KodeIndikator, tujuanOpdId,
			ind.Indikator, ind.RumusPerhitungan.String,
			ind.SumberData.String, ind.DefinisiOperasional.String, ind.Jenis,
		)
		if err != nil {
			return err
		}
		// 1 target per indikator; jenis target = jenis indikator (ranwal/rankhir/penetapan)
		t := ind.Target[0]
		_, err = tx.ExecContext(ctx,
			"INSERT INTO tb_target (id, indikator_id, target, satuan, tahun, jenis) VALUES (?, ?, ?, ?, ?, ?)",
			t.Id, ind.KodeIndikator, t.Target, t.Satuan, t.Tahun, ind.Jenis,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// UPDATE: hanya UPDATE, kode_indikator wajib ada
func (r *TujuanOpdRepositoryImpl) UpdateRenjaIndikator(
	ctx context.Context, tx *sql.Tx,
	indikators []domain.Indikator,
) error {
	for _, ind := range indikators {
		_, err := tx.ExecContext(ctx, `
            UPDATE tb_indikator_matrix
            SET indikator = ?, rumus_perhitungan = ?, sumber_data = ?,
                definisi_operasional = ?, jenis = ?
            WHERE kode_indikator = ?`,
			ind.Indikator, ind.RumusPerhitungan.String,
			ind.SumberData.String, ind.DefinisiOperasional.String,
			ind.Jenis, ind.KodeIndikator,
		)
		if err != nil {
			return err
		}
		// DELETE target lama dengan tahun+jenis yang sama, lalu INSERT baru
		t := ind.Target[0]
		_, err = tx.ExecContext(ctx,
			"DELETE FROM tb_target WHERE indikator_id = ? AND tahun = ? AND jenis = ?",
			ind.KodeIndikator, t.Tahun, ind.Jenis,
		)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx,
			"INSERT INTO tb_target (id, indikator_id, target, satuan, tahun, jenis) VALUES (?, ?, ?, ?, ?, ?)",
			t.Id, ind.KodeIndikator, t.Target, t.Satuan, t.Tahun, ind.Jenis,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TujuanOpdRepositoryImpl) DeleteIndikatorTargetRenja(ctx context.Context, tx *sql.Tx, indikatorId string) error {
	// 1. Hapus Child terlebih dahulu (tb_target)
	_, err := tx.ExecContext(ctx, "DELETE FROM tb_target WHERE indikator_id = ?", indikatorId)
	if err != nil {
		return err
	}

	// 2. Hapus Parent (tb_indikator_matrix)
	_, err = tx.ExecContext(ctx, "DELETE FROM tb_indikator_matrix WHERE kode_indikator = ?", indikatorId)
	if err != nil {
		return err
	}

	return nil
}

func (r *TujuanOpdRepositoryImpl) FindIndikatorByKodeIndikator(
	ctx context.Context, tx *sql.Tx, kodeIndikator string,
) (domain.Indikator, error) {
	row := tx.QueryRowContext(ctx, `
        SELECT kode_indikator,
               COALESCE(indikator, ''),
               COALESCE(rumus_perhitungan, ''),
               COALESCE(sumber_data, ''),
               COALESCE(definisi_operasional, ''),
               COALESCE(jenis, '')
        FROM tb_indikator_matrix
        WHERE kode_indikator = ?`,
		kodeIndikator,
	)
	var indikator domain.Indikator
	err := row.Scan(
		&indikator.KodeIndikator, // ← tidak scan id sama sekali
		&indikator.Indikator,
		&indikator.RumusPerhitungan,
		&indikator.SumberData,
		&indikator.DefinisiOperasional,
		&indikator.Jenis,
	)
	if err != nil {
		return domain.Indikator{}, err
	}
	return indikator, nil
}

// FindAllByTahunForPokin identik dengan FindAllByTahun namun diperuntukkan
// khusus pemanggil dari pohon kinerja. Indikator default dari renstra,
// target difilter berdasarkan jenisIndikator dengan fallback renstra backward compat.
func (repository *TujuanOpdRepositoryImpl) FindAllByTahunForPokin(ctx context.Context, tx *sql.Tx, kodeOpd, tahun, jenisPeriode, jenisIndikator string) ([]domain.TujuanOpd, error) {
	var (
		targetJenisClause string
		finalArgs         []interface{}
	)
	if jenisIndikator == "renstra" {
		targetJenisClause = "(tg.jenis = 'renstra' OR tg.jenis = '')"
		finalArgs = []interface{}{tahun, kodeOpd, jenisPeriode, tahun, tahun}
	} else {
		targetJenisClause = "tg.jenis = ?"
		finalArgs = []interface{}{tahun, jenisIndikator, kodeOpd, jenisPeriode, tahun, tahun}
	}

	query := fmt.Sprintf(`
    SELECT
        t.id,
        t.kode_opd,
        COALESCE(t.kode_bidang_urusan, '')      AS kode_bidang_urusan,
        t.tujuan,
        t.tahun_awal,
        t.tahun_akhir,
        t.jenis_periode,
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
        im_tg.tahun_target
    FROM tb_tujuan_opd t
    JOIN (
        SELECT
            im.id                                  AS indikator_id,
            im.kode_indikator                      AS kode_indikator,
            im.tujuan_opd_id,
            COALESCE(im.indikator, '')             AS indikator,
            COALESCE(im.rumus_perhitungan, '')     AS rumus_perhitungan,
            COALESCE(im.sumber_data, '')           AS sumber_data,
            COALESCE(im.definisi_operasional, '')  AS definisi_operasional,
            COALESCE(im.jenis, '')                 AS indikator_jenis,
            tg.id                                  AS target_id,
            tg.target                              AS target_value,
            tg.satuan,
            tg.tahun                               AS tahun_target
        FROM tb_indikator_matrix im
        LEFT JOIN tb_target tg
            ON im.kode_indikator = tg.indikator_id
            AND tg.tahun = ?
            AND %s
        WHERE im.jenis = 'renstra'
    ) im_tg ON t.id = im_tg.tujuan_opd_id
    WHERE t.kode_opd      = ?
      AND t.jenis_periode  = ?
      AND CAST(t.tahun_awal  AS SIGNED) <= CAST(? AS SIGNED)
      AND CAST(t.tahun_akhir AS SIGNED) >= CAST(? AS SIGNED)
    ORDER BY t.id, im_tg.indikator_id`, targetJenisClause)
	rows, err := tx.QueryContext(ctx, query, finalArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tujuanOpdMap := make(map[int]*domain.TujuanOpd)
	indikatorSeen := make(map[string]bool) // key: "tujuanId-indikatorId"
	tujuanOrder := []int{}
	for rows.Next() {
		var (
			tujuanId            int
			kodeOpdData         string
			kodeBidangUrusan    string
			tujuan              string
			tahunAwalData       string
			tahunAkhirData      string
			jenisPeriodeData    string
			indikatorId         sql.NullString
			kodeIndikator       sql.NullString
			indikatorNama       sql.NullString
			rumusPerhitungan    sql.NullString
			sumberData          sql.NullString
			definisiOperasional sql.NullString // NEW
			indikatorJenis      sql.NullString // im.jenis
			targetId            sql.NullString
			targetValue         sql.NullString
			satuan              sql.NullString
			tahunTarget         sql.NullString
		)
		err := rows.Scan(
			&tujuanId,
			&kodeOpdData,
			&kodeBidangUrusan,
			&tujuan,
			&tahunAwalData,
			&tahunAkhirData,
			&jenisPeriodeData,
			&indikatorId,
			&kodeIndikator,
			&indikatorNama,
			&rumusPerhitungan,
			&sumberData,
			&definisiOperasional, // NEW
			&indikatorJenis,
			&targetId,
			&targetValue,
			&satuan,
			&tahunTarget,
		)
		if err != nil {
			return nil, err
		}
		if _, exists := tujuanOpdMap[tujuanId]; !exists {
			tujuanOpdMap[tujuanId] = &domain.TujuanOpd{
				Id:               tujuanId,
				KodeOpd:          kodeOpdData,
				KodeBidangUrusan: kodeBidangUrusan,
				Tujuan:           tujuan,
				TahunAwal:        tahunAwalData,
				TahunAkhir:       tahunAkhirData,
				JenisPeriode:     jenisPeriodeData,
				Indikator:        []domain.Indikator{},
			}
			tujuanOrder = append(tujuanOrder, tujuanId)
		}
		if indikatorId.Valid {
			mapKey := fmt.Sprintf("%d-%s", tujuanId, indikatorId.String)
			if !indikatorSeen[mapKey] {
				indikatorSeen[mapKey] = true
				tujuanOpdMap[tujuanId].Indikator = append(tujuanOpdMap[tujuanId].Indikator, domain.Indikator{
					Id:                  indikatorId.String,
					KodeIndikator:       kodeIndikator.String,
					Indikator:           indikatorNama.String,
					RumusPerhitungan:    rumusPerhitungan,
					SumberData:          sumberData,
					DefinisiOperasional: definisiOperasional, // NEW
					Jenis:               indikatorJenis.String,
					TujuanOpdId:         tujuanId,
					Target:              []domain.Target{},
				})
			}
			if targetId.Valid {
				target := domain.Target{
					Id:          targetId.String,
					IndikatorId: kodeIndikator.String,
					Target:      targetValue.String,
					Satuan:      satuan.String,
					Tahun:       tahunTarget.String,
				}
				for idx := range tujuanOpdMap[tujuanId].Indikator {
					if tujuanOpdMap[tujuanId].Indikator[idx].Id == indikatorId.String {
						tujuanOpdMap[tujuanId].Indikator[idx].Target = append(
							tujuanOpdMap[tujuanId].Indikator[idx].Target, target,
						)
						break
					}
				}
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Ranwal / Rankhir: tepat 1 slot target per indikator untuk tahun yang diminta
	var result []domain.TujuanOpd
	for _, id := range tujuanOrder {
		tujuanOpd := tujuanOpdMap[id]
		for i := range tujuanOpd.Indikator {
			if len(tujuanOpd.Indikator[i].Target) == 0 {
				// Belum ada data target di DB → buat slot kosong
				tujuanOpd.Indikator[i].Target = []domain.Target{{
					Id:          "",
					IndikatorId: tujuanOpd.Indikator[i].KodeIndikator,
					Target:      "",
					Satuan:      "",
					Tahun:       tahun,
				}}
			}
		}
		result = append(result, *tujuanOpd)
	}
	if len(result) == 0 {
		return make([]domain.TujuanOpd, 0), nil
	}
	return result, nil
}

func (repository *TujuanOpdRepositoryImpl) SetTujuanOpdLocked(ctx context.Context, tx *sql.Tx, id int, locked bool) error {
	val := 0
	if locked {
		val = 1
	}
	_, err := tx.ExecContext(ctx,
		`UPDATE tb_tujuan_opd SET is_locked = ? WHERE id = ?`, val, id)
	return err
}

func (repository *TujuanOpdRepositoryImpl) FindAllByTahunDualJenis(ctx context.Context, tx *sql.Tx, kodeOpd, tahun, jenisPeriode, jenisOverride string) ([]domain.TujuanOpd, error) {
	// Query 1: ambil tujuan + indikator renstra + target renstra (semua tahun dalam periode)
	queryRenstra := `
		SELECT
			t.id,
			t.kode_opd,
			COALESCE(t.kode_bidang_urusan, '')      AS kode_bidang_urusan,
			t.tujuan,
			t.tahun_awal,
			t.tahun_akhir,
			t.jenis_periode,
			im.id                                   AS indikator_id,
			im.kode_indikator,
			COALESCE(im.indikator, '')              AS indikator_nama,
			COALESCE(im.rumus_perhitungan, '')      AS rumus_perhitungan,
			COALESCE(im.sumber_data, '')            AS sumber_data,
			COALESCE(im.definisi_operasional, '')   AS definisi_operasional,
			COALESCE(im.jenis, '')                  AS indikator_jenis,
			tg.id                                   AS target_id,
			COALESCE(tg.target, '')                 AS target_value,
			COALESCE(tg.satuan, '')                 AS satuan,
			COALESCE(tg.tahun, '')                  AS tahun_target
		FROM tb_tujuan_opd t
		LEFT JOIN tb_indikator_matrix im
			ON t.id = im.tujuan_opd_id AND im.jenis = 'renstra'
		LEFT JOIN tb_target tg
			ON im.kode_indikator = tg.indikator_id
			AND (tg.jenis = 'renstra' OR tg.jenis = '')
		WHERE t.kode_opd     = ?
		  AND t.jenis_periode = ?
		  AND CAST(t.tahun_awal  AS SIGNED) <= CAST(? AS SIGNED)
		  AND CAST(t.tahun_akhir AS SIGNED) >= CAST(? AS SIGNED)
		ORDER BY t.id, im.id, CAST(tg.tahun AS SIGNED)
	`
	rowsRenstra, err := tx.QueryContext(ctx, queryRenstra,
		kodeOpd, jenisPeriode, tahun, tahun,
	)
	if err != nil {
		return nil, fmt.Errorf("FindAllByTahunDualJenis renstra: %w", err)
	}
	defer rowsRenstra.Close()
	type indikatorKey struct {
		tujuanId int
		indId    string
	}
	tujuanMap := make(map[int]*domain.TujuanOpd)
	tujuanOrder := []int{}
	indRenstraTargets := make(map[indikatorKey][]domain.Target) // kode_indikator → []Target renstra
	for rowsRenstra.Next() {
		var (
			tId                             int
			kodeOpdData                     string
			kodeBidangUrusan                string
			tujuanTeks                      string
			tahunAwal, tahunAkhir, jenisPer string
			indikatorId                     sql.NullString
			kodeIndikator                   sql.NullString
			indikatorNama                   sql.NullString
			rumusPerhitungan                sql.NullString
			sumberData                      sql.NullString
			definisiOperasional             sql.NullString
			indikatorJenis                  sql.NullString
			targetId                        sql.NullString
			targetValue                     sql.NullString
			satuan                          sql.NullString
			tahunTarget                     sql.NullString
		)
		if err := rowsRenstra.Scan(
			&tId, &kodeOpdData, &kodeBidangUrusan, &tujuanTeks,
			&tahunAwal, &tahunAkhir, &jenisPer,
			&indikatorId, &kodeIndikator, &indikatorNama,
			&rumusPerhitungan, &sumberData, &definisiOperasional, &indikatorJenis,
			&targetId, &targetValue, &satuan, &tahunTarget,
		); err != nil {
			return nil, err
		}
		if _, exists := tujuanMap[tId]; !exists {
			tujuanMap[tId] = &domain.TujuanOpd{
				Id: tId, KodeOpd: kodeOpdData,
				KodeBidangUrusan: kodeBidangUrusan, Tujuan: tujuanTeks,
				TahunAwal: tahunAwal, TahunAkhir: tahunAkhir, JenisPeriode: jenisPer,
				Indikator: []domain.Indikator{},
			}
			tujuanOrder = append(tujuanOrder, tId)
		}
		if indikatorId.Valid {
			key := indikatorKey{tId, indikatorId.String}
			// Tambahkan indikator (hanya sekali per indikator_id)
			found := false
			for _, ind := range tujuanMap[tId].Indikator {
				if ind.Id == indikatorId.String {
					found = true
					break
				}
			}
			if !found {
				tujuanMap[tId].Indikator = append(tujuanMap[tId].Indikator, domain.Indikator{
					Id:                  indikatorId.String,
					KodeIndikator:       kodeIndikator.String,
					Indikator:           indikatorNama.String,
					RumusPerhitungan:    rumusPerhitungan,
					SumberData:          sumberData,
					DefinisiOperasional: definisiOperasional,
					Jenis:               indikatorJenis.String,
					TujuanOpdId:         tId,
					Target:              []domain.Target{},
				})
			}
			// Simpan target renstra ke map
			if targetId.Valid && tahunTarget.Valid {
				indRenstraTargets[key] = append(indRenstraTargets[key], domain.Target{
					Id:          targetId.String,
					IndikatorId: kodeIndikator.String,
					Target:      targetValue.String,
					Satuan:      satuan.String,
					Tahun:       tahunTarget.String,
					Jenis:       "renstra",
				})
			}
		}
	}
	if err := rowsRenstra.Err(); err != nil {
		return nil, err
	}
	// Query 2: ambil target jenisOverride (rankhir/penetapan) untuk tahun yang diminta
	queryOverride := `
		SELECT
			im.id          AS indikator_id,
			im.tujuan_opd_id,
			tg.id          AS target_id,
			COALESCE(tg.target, '')  AS target_value,
			COALESCE(tg.satuan, '')  AS satuan,
			COALESCE(tg.tahun, '')   AS tahun_target
		FROM tb_indikator_matrix im
		INNER JOIN tb_tujuan_opd t ON im.tujuan_opd_id = t.id
		LEFT JOIN tb_target tg
			ON im.kode_indikator = tg.indikator_id
			AND tg.tahun = ?
			AND tg.jenis = ?
		WHERE t.kode_opd     = ?
		  AND t.jenis_periode = ?
		  AND CAST(t.tahun_awal  AS SIGNED) <= CAST(? AS SIGNED)
		  AND CAST(t.tahun_akhir AS SIGNED) >= CAST(? AS SIGNED)
		  AND im.jenis = 'renstra'
	`
	rowsOverride, err := tx.QueryContext(ctx, queryOverride,
		tahun, jenisOverride,
		kodeOpd, jenisPeriode, tahun, tahun,
	)
	if err != nil {
		return nil, fmt.Errorf("FindAllByTahunDualJenis override: %w", err)
	}
	defer rowsOverride.Close()
	// Map: indikatorKey → target override (untuk tahun yang diminta)
	type overrideVal struct {
		targetId    string
		targetValue string
		satuan      string
	}
	overrideMap := make(map[indikatorKey]overrideVal)
	for rowsOverride.Next() {
		var (
			indId       sql.NullString
			tujuanId    int
			targetId    sql.NullString
			targetValue sql.NullString
			satuan      sql.NullString
			tahunTg     sql.NullString
		)
		if err := rowsOverride.Scan(
			&indId, &tujuanId, &targetId, &targetValue, &satuan, &tahunTg,
		); err != nil {
			return nil, err
		}
		if indId.Valid {
			overrideMap[indikatorKey{tujuanId, indId.String}] = overrideVal{
				targetId:    targetId.String,
				targetValue: targetValue.String,
				satuan:      satuan.String,
			}
		}
	}
	if err := rowsOverride.Err(); err != nil {
		return nil, err
	}
	// Gabungkan: generate slot renstra lengkap + set target override
	var result []domain.TujuanOpd
	for _, tId := range tujuanOrder {
		tujuanOpd := tujuanMap[tId]
		tahunAwalInt, _ := strconv.Atoi(tujuanOpd.TahunAwal)
		tahunAkhirInt, _ := strconv.Atoi(tujuanOpd.TahunAkhir)
		for i := range tujuanOpd.Indikator {
			ind := &tujuanOpd.Indikator[i]
			key := indikatorKey{tId, ind.Id}
			// Build renstra slots (full period)
			renstraRaw := indRenstraTargets[key]
			renstraByTahun := make(map[string]domain.Target)
			for _, t := range renstraRaw {
				renstraByTahun[t.Tahun] = t
			}
			var renstraSlots []domain.Target
			for yr := tahunAwalInt; yr <= tahunAkhirInt; yr++ {
				yrStr := strconv.Itoa(yr)
				if t, ok := renstraByTahun[yrStr]; ok {
					renstraSlots = append(renstraSlots, t)
				} else {
					renstraSlots = append(renstraSlots, domain.Target{
						Id: "", IndikatorId: ind.KodeIndikator,
						Target: "", Satuan: "", Tahun: yrStr, Jenis: "renstra",
					})
				}
			}
			ind.TargetRenstra = renstraSlots
			// Build target override (1 slot untuk tahun yang diminta)
			ov := overrideMap[key]
			ind.Target = []domain.Target{{
				Id:          ov.targetId,
				IndikatorId: ind.KodeIndikator,
				Target:      ov.targetValue,
				Satuan:      ov.satuan,
				Tahun:       tahun,
				Jenis:       jenisOverride,
			}}
		}
		result = append(result, *tujuanOpd)
	}
	if len(result) == 0 {
		return make([]domain.TujuanOpd, 0), nil
	}
	return result, nil
}

func (repository *TujuanOpdRepositoryImpl) FindIndikatorTargetsByTujuanIds(
	ctx context.Context,
	tx *sql.Tx,
	tujuanIds []int,
) ([]domain.Indikator, error) {

	if len(tujuanIds) == 0 {
		return []domain.Indikator{}, nil
	}

	// buat placeholder (?, ?, ?, ...)
	placeholders := make([]string, len(tujuanIds))
	args := make([]any, len(tujuanIds))

	for i, id := range tujuanIds {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT ind.id, ind.indikator, ind.tujuan_opd_id,
                     ind.rumus_perhitungan, ind.sumber_data,
                     tar.id, tar.target, tar.satuan, tar.tahun, tar.jenis
		FROM tb_indikator ind
                LEFT JOIN tb_target tar ON tar.indikator_id = ind.id
		WHERE ind.tujuan_opd_id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var (
			indId         string
			indikatorName string
			tujuanOpdId   int

			rumusPerhitunganNS sql.NullString
			sumberDataNS       sql.NullString
			tarIdNS            sql.NullString
			targetNS           sql.NullString
			satuanNS           sql.NullString
			tahunNS            sql.NullString
			jenisNS            sql.NullString
		)

		err := rows.Scan(
			&indId,
			&indikatorName,
			&tujuanOpdId,
			&rumusPerhitunganNS,
			&sumberDataNS,
			&tarIdNS,
			&targetNS,
			&satuanNS,
			&tahunNS,
			&jenisNS,
		)
		if err != nil {
			return nil, err
		}

		// 🔹 ambil / buat indikator
		ind, exists := indikatorMap[indId]
		if !exists {
			ind = &domain.Indikator{
				Id:               indId,
				KodeIndikator:    indId,
				Indikator:        indikatorName,
				TujuanOpdId:      tujuanOpdId,
				RumusPerhitungan: rumusPerhitunganNS,
				SumberData:       sumberDataNS,
				Target:           []domain.Target{},
			}
			indikatorMap[indId] = ind
		}

		// 🔹 kalau ada target, append
		if tarIdNS.Valid {
			target := domain.Target{
				Id: tarIdNS.String,
			}

			if targetNS.Valid {
				target.Target = targetNS.String
			}
			if satuanNS.Valid {
				target.Satuan = satuanNS.String
			}
			if tahunNS.Valid {
				target.Tahun = tahunNS.String
			}
			if jenisNS.Valid {
				target.Jenis = jenisNS.String
			}

			ind.Target = append(ind.Target, target)
		}
	}

	result := make([]domain.Indikator, 0, len(indikatorMap))
	for _, v := range indikatorMap {
		result = append(result, *v)
	}

	return result, nil
}

func (repository *TujuanOpdRepositoryImpl) FindByIdOnly(
	ctx context.Context,
	tx *sql.Tx,
	tujuanOpdId int,
) (domain.TujuanOpd, error) {

	script := `
		SELECT
			t.id,
			t.kode_opd,
			COALESCE(t.kode_bidang_urusan, '') as kode_bidang_urusan,
			t.tujuan,
			t.tahun_awal,
			t.tahun_akhir,
			t.jenis_periode
		FROM tb_tujuan_opd t
		WHERE t.id = ?
	`

	var result domain.TujuanOpd

	err := tx.QueryRowContext(ctx, script, tujuanOpdId).Scan(
		&result.Id,
		&result.KodeOpd,
		&result.KodeBidangUrusan,
		&result.Tujuan,
		&result.TahunAwal,
		&result.TahunAkhir,
		&result.JenisPeriode,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.TujuanOpd{}, err
		}
		return domain.TujuanOpd{}, err
	}

	return result, nil
}

func (repository *TujuanOpdRepositoryImpl) FindIndikatorTargetsRenstraByTujuanIds(
	ctx context.Context,
	tx *sql.Tx,
	tujuanIds []int,
) ([]domain.Indikator, error) {

	if len(tujuanIds) == 0 {
		return []domain.Indikator{}, nil
	}

	// buat placeholder (?, ?, ?, ...)
	placeholders := make([]string, len(tujuanIds))
	args := make([]any, len(tujuanIds))

	for i, id := range tujuanIds {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT ind.kode_indikator, ind.indikator, ind.tujuan_opd_id,
                     ind.rumus_perhitungan, ind.sumber_data,
                     ind.definisi_operasional,
                     tar.id, tar.target, tar.satuan, tar.tahun, tar.jenis
		FROM tb_indikator_matrix ind
                LEFT JOIN tb_target tar ON tar.indikator_id = ind.kode_indikator
		WHERE ind.tujuan_opd_id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var (
			indId                 string
			indikatorName         string
			tujuanOpdId           int
			rumusPerhitunganNS    sql.NullString
			sumberDataNS          sql.NullString
			definisiOperasionalNS sql.NullString

			tarIdNS  sql.NullString
			targetNS sql.NullString
			satuanNS sql.NullString
			tahunNS  sql.NullString
			jenisNS  sql.NullString
		)

		err := rows.Scan(
			&indId,
			&indikatorName,
			&tujuanOpdId,
			&rumusPerhitunganNS,
			&sumberDataNS,
			&definisiOperasionalNS,
			&tarIdNS,
			&targetNS,
			&satuanNS,
			&tahunNS,
			&jenisNS,
		)
		if err != nil {
			return nil, err
		}

		// 🔹 ambil / buat indikator
		ind, exists := indikatorMap[indId]
		if !exists {
			ind = &domain.Indikator{
				Id:                  indId,
				KodeIndikator:       indId,
				Indikator:           indikatorName,
				TujuanOpdId:         tujuanOpdId,
				RumusPerhitungan:    rumusPerhitunganNS,
				DefinisiOperasional: definisiOperasionalNS,
				SumberData:          sumberDataNS,
				Target:              []domain.Target{},
			}
			indikatorMap[indId] = ind
		}

		// 🔹 kalau ada target, append
		if tarIdNS.Valid {
			target := domain.Target{
				Id: tarIdNS.String,
			}

			if targetNS.Valid {
				target.Target = targetNS.String
			}
			if satuanNS.Valid {
				target.Satuan = satuanNS.String
			}
			if tahunNS.Valid {
				target.Tahun = tahunNS.String
			}
			if jenisNS.Valid {
				target.Jenis = jenisNS.String
			}

			ind.Target = append(ind.Target, target)
		}
	}

	result := make([]domain.Indikator, 0, len(indikatorMap))
	for _, v := range indikatorMap {
		result = append(result, *v)
	}

	return result, nil
}

func (repository *TujuanOpdRepositoryImpl) FindAllOnly(ctx context.Context, tx *sql.Tx, kodeOpd string, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]domain.TujuanOpd, error) {
	scriptTujuan := `
        SELECT
            t.id,
            t.kode_opd,
            COALESCE(t.kode_bidang_urusan, '') as kode_bidang_urusan,
            t.tujuan,
            t.tahun_awal,
            t.tahun_akhir,
            t.jenis_periode
        FROM tb_tujuan_opd t
        WHERE t.kode_opd = ?
        AND t.tahun_awal = ?
        AND t.tahun_akhir = ?
        AND t.jenis_periode = ?
        ORDER BY t.id ASC
    `

	rows, err := tx.QueryContext(ctx, scriptTujuan,
		kodeOpd,
		tahunAwal,
		tahunAkhir,
		jenisPeriode,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tujuanOpdMap := make(map[int]*domain.TujuanOpd)

	for rows.Next() {
		var (
			tujuanId         int
			kodeOpd          string
			kodeBidangUrusan string
			tujuan           string
			// periodeId        int
			tahunAwalData    string
			tahunAkhirData   string
			jenisPeriodeData string
		)

		err := rows.Scan(
			&tujuanId,
			&kodeOpd,
			&kodeBidangUrusan,
			&tujuan,
			// &periodeId,
			&tahunAwalData,
			&tahunAkhirData,
			&jenisPeriodeData,
		)
		if err != nil {
			return nil, err
		}

		// Buat atau ambil TujuanOpd
		if _, exists := tujuanOpdMap[tujuanId]; !exists {
			tujuanOpdMap[tujuanId] = &domain.TujuanOpd{
				Id:               tujuanId,
				KodeOpd:          kodeOpd,
				KodeBidangUrusan: kodeBidangUrusan,
				Tujuan:           tujuan,
				TahunAwal:        tahunAwalData,
				TahunAkhir:       tahunAkhirData,
				JenisPeriode:     jenisPeriodeData,
				Indikator:        []domain.Indikator{},
			}
		}

	}

	// Perbaikan pada bagian generate target
	var result []domain.TujuanOpd
	for _, tujuanOpd := range tujuanOpdMap {
		result = append(result, *tujuanOpd)
	}

	if len(result) == 0 {
		return make([]domain.TujuanOpd, 0), nil
	}

	return result, nil
}

func (repository *TujuanOpdRepositoryImpl) DeleteIndikatorByIds(ctx context.Context, tx *sql.Tx, indikatorIds []string) error {
	if len(indikatorIds) == 0 {
		return nil
	}

	// Build query dengan IN clause
	placeholders := make([]string, len(indikatorIds))
	args := make([]interface{}, len(indikatorIds))
	for i, id := range indikatorIds {
		placeholders[i] = "?"
		args[i] = id
	}

	script := fmt.Sprintf(`
		DELETE FROM tb_indikator_matrix
		WHERE kode_indikator IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

// ─────────────────────────────────────────────────────────────────
// Target-only CRUD helpers (ranwal / rankhir / penetapan)
// Indikator tetap dari renstra; hanya baris di tb_target yang dikelola.
// ─────────────────────────────────────────────────────────────────

func (repository *TujuanOpdRepositoryImpl) TargetOpdExistsByKey(
	ctx context.Context, tx *sql.Tx,
	kodeIndikator, tahun, jenis string,
) (bool, error) {
	var count int
	err := tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM tb_target WHERE indikator_id = ? AND tahun = ? AND jenis = ?`,
		kodeIndikator, tahun, jenis,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repository *TujuanOpdRepositoryImpl) CreateTargetOpdSingle(
	ctx context.Context, tx *sql.Tx,
	t domain.Target,
) (domain.Target, error) {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO tb_target (id, indikator_id, target, satuan, tahun, jenis) VALUES (?, ?, ?, ?, ?, ?)`,
		t.Id, t.IndikatorId, t.Target, t.Satuan, t.Tahun, t.Jenis,
	)
	if err != nil {
		return domain.Target{}, err
	}
	return t, nil
}

func (repository *TujuanOpdRepositoryImpl) FindTargetOpdById(
	ctx context.Context, tx *sql.Tx,
	id string,
) (domain.Target, error) {
	var t domain.Target
	err := tx.QueryRowContext(ctx,
		`SELECT id, indikator_id,
		        COALESCE(target,'')  AS target,
		        COALESCE(satuan,'')  AS satuan,
		        COALESCE(tahun,'')   AS tahun,
		        COALESCE(jenis,'')   AS jenis
		 FROM tb_target WHERE id = ?`,
		id,
	).Scan(&t.Id, &t.IndikatorId, &t.Target, &t.Satuan, &t.Tahun, &t.Jenis)
	if err != nil {
		return domain.Target{}, err
	}
	return t, nil
}

func (repository *TujuanOpdRepositoryImpl) UpdateTargetOpdById(
	ctx context.Context, tx *sql.Tx,
	id, targetValue, satuan string,
) (domain.Target, error) {
	_, err := tx.ExecContext(ctx,
		`UPDATE tb_target SET target = ?, satuan = ? WHERE id = ?`,
		targetValue, satuan, id,
	)
	if err != nil {
		return domain.Target{}, err
	}
	return repository.FindTargetOpdById(ctx, tx, id)
}

func (repository *TujuanOpdRepositoryImpl) DeleteTargetOpdByJenis(
	ctx context.Context, tx *sql.Tx,
	kodeIndikator, tahun, jenis string,
) error {
	_, err := tx.ExecContext(ctx,
		`DELETE FROM tb_target WHERE indikator_id = ? AND tahun = ? AND jenis = ?`,
		kodeIndikator, tahun, jenis,
	)
	return err
}
