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

type SasaranPemdaRepositoryImpl struct{}

func NewSasaranPemdaRepositoryImpl() *SasaranPemdaRepositoryImpl {
	return &SasaranPemdaRepositoryImpl{}
}

// ═══════════════════════════════════════════════════════════════════
// INSERT HELPER
// ═══════════════════════════════════════════════════════════════════
func (r *SasaranPemdaRepositoryImpl) insertIndikator(
	ctx context.Context, tx *sql.Tx, ind domain.IndikatorPemda, sasaranPemdaId int,
) (int, error) {
	res, err := tx.ExecContext(ctx,
		`INSERT INTO tb_indikator_matrix_pemda
			(kode_indikator, sasaran_pemda_id, tujuan_pemda_id,
			 indikator, rumus_perhitungan, sumber_data, definisi_operasional, jenis)
		 VALUES (?, ?, 0, ?, ?, ?, ?, 'renstra')`,
		ind.KodeIndikator, sasaranPemdaId,
		ind.Indikator, ind.RumusPerhitungan, ind.SumberData, ind.DefinisiOperasional,
	)
	if err != nil {
		return 0, err
	}
	newId, _ := res.LastInsertId()
	for _, t := range ind.Target {
		if t.Target == "" && t.Satuan == "" {
			continue
		}
		if err := r.insertTarget(ctx, tx, t, ind.KodeIndikator); err != nil {
			return 0, err
		}
	}
	return int(newId), nil
}
func (r *SasaranPemdaRepositoryImpl) insertTarget(
	ctx context.Context, tx *sql.Tx, t domain.TargetPemda, kodeIndikator string,
) error {
	jenis := t.Jenis
	if jenis == "" {
		jenis = "renstra"
	}
	_, err := tx.ExecContext(ctx,
		`INSERT INTO tb_target_pemda (kode_indikator, target, satuan, tahun, jenis)
		 VALUES (?, ?, ?, ?, ?)`,
		kodeIndikator, t.Target, t.Satuan, t.Tahun, jenis,
	)
	return err
}

// ═══════════════════════════════════════════════════════════════════
// CREATE
// ═══════════════════════════════════════════════════════════════════
func (r *SasaranPemdaRepositoryImpl) Create(
	ctx context.Context, tx *sql.Tx, sp domain.SasaranPemda,
) (domain.SasaranPemda, error) {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO tb_sasaran_pemda
			(id, tujuan_pemda_id, subtema_id, sasaran_pemda, periode_id, tahun_awal, tahun_akhir, jenis_periode)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		sp.Id, sp.TujuanPemdaId, sp.SubtemaId, sp.SasaranPemda,
		sp.PeriodeId, sp.TahunAwal, sp.TahunAkhir, sp.JenisPeriode,
	)
	if err != nil {
		return sp, err
	}
	for _, ind := range sp.Indikator {
		if _, err := r.insertIndikator(ctx, tx, ind, sp.Id); err != nil {
			return sp, err
		}
	}
	return sp, nil
}

// ═══════════════════════════════════════════════════════════════════
// UPDATE — upsert per id (int) untuk indikator dan target
//
//	rankhir & penetapan TIDAK disentuh
//
// ═══════════════════════════════════════════════════════════════════
func (r *SasaranPemdaRepositoryImpl) Update(
	ctx context.Context, tx *sql.Tx, sp domain.SasaranPemda,
) (domain.SasaranPemda, error) {
	_, err := tx.ExecContext(ctx,
		`UPDATE tb_sasaran_pemda
		 SET sasaran_pemda=?, tujuan_pemda_id=?, subtema_id=?
		 WHERE id=?`,
		sp.SasaranPemda, sp.TujuanPemdaId, sp.SubtemaId, sp.Id,
	)
	if err != nil {
		return sp, err
	}
	keepIndIds := make([]int, 0, len(sp.Indikator))
	for _, ind := range sp.Indikator {
		if ind.Id > 0 {
			// UPDATE indikator existing — id tetap, kode_indikator tidak boleh berubah
			_, err = tx.ExecContext(ctx,
				`UPDATE tb_indikator_matrix_pemda
				 SET indikator=?, rumus_perhitungan=?, sumber_data=?, definisi_operasional=?
				 WHERE id=? AND sasaran_pemda_id=?`,
				ind.Indikator, ind.RumusPerhitungan, ind.SumberData, ind.DefinisiOperasional,
				ind.Id, sp.Id,
			)
			if err != nil {
				return sp, fmt.Errorf("update indikator id %d: %w", ind.Id, err)
			}
			keepIndIds = append(keepIndIds, ind.Id)
		} else {
			// INSERT indikator baru
			newId, err := r.insertIndikator(ctx, tx, ind, sp.Id)
			if err != nil {
				return sp, fmt.Errorf("insert indikator baru: %w", err)
			}
			keepIndIds = append(keepIndIds, newId)
			continue // target sudah di-insert oleh insertIndikator
		}
		// Upsert target renstra per indikator (hanya jika update indikator)
		keepTargetIds := make([]int, 0, len(ind.Target))
		for _, t := range ind.Target {
			if t.Id > 0 {
				_, err = tx.ExecContext(ctx,
					`UPDATE tb_target_pemda
					 SET target=?, satuan=?, tahun=?
					 WHERE id=? AND kode_indikator=?
					   AND (jenis='renstra' OR jenis='' OR jenis IS NULL)`,
					t.Target, t.Satuan, t.Tahun, t.Id, ind.KodeIndikator,
				)
				if err != nil {
					return sp, fmt.Errorf("update target id %d: %w", t.Id, err)
				}
				keepTargetIds = append(keepTargetIds, t.Id)
			} else {
				if t.Target == "" && t.Satuan == "" {
					continue
				}
				res, err := tx.ExecContext(ctx,
					`INSERT INTO tb_target_pemda (kode_indikator, target, satuan, tahun, jenis)
					 VALUES (?, ?, ?, ?, 'renstra')`,
					ind.KodeIndikator, t.Target, t.Satuan, t.Tahun,
				)
				if err != nil {
					return sp, fmt.Errorf("insert target baru: %w", err)
				}
				newId, _ := res.LastInsertId()
				keepTargetIds = append(keepTargetIds, int(newId))
			}
		}
		// Hapus target renstra yang tidak ada di request
		if len(keepTargetIds) > 0 {
			ph := strings.Repeat("?,", len(keepTargetIds))
			ph = ph[:len(ph)-1]
			args := make([]interface{}, 0, len(keepTargetIds)+1)
			args = append(args, ind.KodeIndikator)
			for _, id := range keepTargetIds {
				args = append(args, id)
			}
			_, err = tx.ExecContext(ctx, fmt.Sprintf(
				`DELETE FROM tb_target_pemda
				 WHERE kode_indikator=?
				   AND (jenis='renstra' OR jenis='' OR jenis IS NULL)
				   AND id NOT IN (%s)`, ph), args...,
			)
		} else {
			_, err = tx.ExecContext(ctx,
				`DELETE FROM tb_target_pemda
				 WHERE kode_indikator=?
				   AND (jenis='renstra' OR jenis='' OR jenis IS NULL)`,
				ind.KodeIndikator,
			)
		}
		if err != nil {
			return sp, fmt.Errorf("hapus target lama: %w", err)
		}
	}
	// Hapus indikator yang tidak ada di request
	if len(keepIndIds) > 0 {
		ph := strings.Repeat("?,", len(keepIndIds))
		ph = ph[:len(ph)-1]
		args := make([]interface{}, 0, len(keepIndIds)+1)
		args = append(args, sp.Id)
		for _, id := range keepIndIds {
			args = append(args, id)
		}
		// hapus target renstra dari indikator yang akan dihapus
		_, err = tx.ExecContext(ctx, fmt.Sprintf(
			`DELETE tg FROM tb_target_pemda tg
			 INNER JOIN tb_indikator_matrix_pemda i ON tg.kode_indikator=i.kode_indikator
			 WHERE i.sasaran_pemda_id=?
			   AND (tg.jenis='renstra' OR tg.jenis='' OR tg.jenis IS NULL)
			   AND i.id NOT IN (%s)`, ph), args...,
		)
		if err != nil {
			return sp, fmt.Errorf("hapus target orphan: %w", err)
		}
		_, err = tx.ExecContext(ctx, fmt.Sprintf(
			`DELETE FROM tb_indikator_matrix_pemda
			 WHERE sasaran_pemda_id=? AND id NOT IN (%s)`, ph), args...,
		)
	} else {
		_, err = tx.ExecContext(ctx,
			`DELETE tg FROM tb_target_pemda tg
			 INNER JOIN tb_indikator_matrix_pemda i ON tg.kode_indikator=i.kode_indikator
			 WHERE i.sasaran_pemda_id=?
			   AND (tg.jenis='renstra' OR tg.jenis='' OR tg.jenis IS NULL)`, sp.Id,
		)
		if err == nil {
			_, err = tx.ExecContext(ctx,
				`DELETE FROM tb_indikator_matrix_pemda WHERE sasaran_pemda_id=?`, sp.Id,
			)
		}
	}
	return sp, err
}

// ═══════════════════════════════════════════════════════════════════
// DELETE
// ═══════════════════════════════════════════════════════════════════
func (r *SasaranPemdaRepositoryImpl) Delete(
	ctx context.Context, tx *sql.Tx, sasaranPemdaId int,
) error {
	if _, err := tx.ExecContext(ctx,
		`DELETE tg FROM tb_target_pemda tg
		 INNER JOIN tb_indikator_matrix_pemda i ON tg.kode_indikator=i.kode_indikator
		 WHERE i.sasaran_pemda_id=?`, sasaranPemdaId,
	); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM tb_indikator_matrix_pemda WHERE sasaran_pemda_id=?`, sasaranPemdaId,
	); err != nil {
		return err
	}
	_, err := tx.ExecContext(ctx,
		`DELETE FROM tb_sasaran_pemda WHERE id=?`, sasaranPemdaId,
	)
	return err
}
func (r *SasaranPemdaRepositoryImpl) DeleteIndikator(
	ctx context.Context, tx *sql.Tx, sasaranPemdaId int,
) error {
	_, err := tx.ExecContext(ctx,
		`DELETE FROM tb_indikator_matrix_pemda WHERE sasaran_pemda_id=?`, sasaranPemdaId,
	)
	return err
}

// ═══════════════════════════════════════════════════════════════════
// FIND BY ID
// ═══════════════════════════════════════════════════════════════════
func (r *SasaranPemdaRepositoryImpl) FindById(
	ctx context.Context, tx *sql.Tx, sasaranPemdaId int,
) (domain.SasaranPemda, error) {
	query := `
		SELECT
			sp.id, sp.tujuan_pemda_id, sp.subtema_id, sp.sasaran_pemda, sp.periode_id,
			COALESCE(p.tahun_awal,'')           AS tahun_awal,
			COALESCE(p.tahun_akhir,'')          AS tahun_akhir,
			COALESCE(p.jenis_periode,'')        AS jenis_periode,
			COALESCE(pk.jenis_pohon,'')         AS jenis_pohon,
			COALESCE(i.id, 0)                   AS indikator_db_id,
			COALESCE(i.kode_indikator,'')       AS kode_indikator,
			COALESCE(i.indikator,'')            AS indikator_text,
			COALESCE(i.rumus_perhitungan,'')    AS rumus,
			COALESCE(i.sumber_data,'')          AS sumber,
			COALESCE(i.definisi_operasional,'') AS definisi,
			COALESCE(t.id, 0)                   AS target_db_id,
			COALESCE(t.target,'')               AS target_value,
			COALESCE(t.satuan,'')               AS target_satuan,
			COALESCE(t.tahun,'')                AS target_tahun,
			COALESCE(t.jenis,'renstra')         AS target_jenis
		FROM tb_sasaran_pemda sp
		LEFT JOIN tb_periode p ON sp.periode_id=p.id
		LEFT JOIN tb_pohon_kinerja pk ON sp.subtema_id=pk.id
		LEFT JOIN tb_indikator_matrix_pemda i
			ON sp.id=i.sasaran_pemda_id
			AND (i.jenis='renstra' OR i.jenis='' OR i.jenis IS NULL)
		LEFT JOIN tb_target_pemda t
			ON t.kode_indikator=i.kode_indikator
			AND (t.jenis='renstra' OR t.jenis='' OR t.jenis IS NULL)
		WHERE sp.id=?
		ORDER BY sp.id, i.id, CAST(t.tahun AS SIGNED)`
	rows, err := tx.QueryContext(ctx, query, sasaranPemdaId)
	if err != nil {
		return domain.SasaranPemda{}, fmt.Errorf("FindById: %w", err)
	}
	defer rows.Close()
	var result domain.SasaranPemda
	firstRow := true
	indMap := make(map[int]*domain.IndikatorPemda) // key: indikator DB id
	for rows.Next() {
		var (
			id, tujuanPemdaId, subtemaId, periodeId                      int
			sasaranText, tahunAwal, tahunAkhir, jenisPeriode, jenisPohon string
			indDbId                                                      int
			kodeIndikator, indText, rumus, sumber, definisi              string
			targetDbId                                                   int
			targetValue, targetSatuan, targetTahun, targetJenis          string
		)
		if err := rows.Scan(
			&id, &tujuanPemdaId, &subtemaId, &sasaranText,
			&periodeId, &tahunAwal, &tahunAkhir, &jenisPeriode, &jenisPohon,
			&indDbId, &kodeIndikator, &indText, &rumus, &sumber, &definisi,
			&targetDbId, &targetValue, &targetSatuan, &targetTahun, &targetJenis,
		); err != nil {
			return domain.SasaranPemda{}, fmt.Errorf("FindById scan: %w", err)
		}
		if firstRow {
			result = domain.SasaranPemda{
				Id: id, TujuanPemdaId: tujuanPemdaId,
				SubtemaId: subtemaId, SasaranPemda: sasaranText,
				PeriodeId: periodeId, JenisPohon: jenisPohon,
				TahunAwal: tahunAwal, TahunAkhir: tahunAkhir, JenisPeriode: jenisPeriode,
				Periode: domain.Periode{
					TahunAwal: tahunAwal, TahunAkhir: tahunAkhir, JenisPeriode: jenisPeriode,
				},
				Indikator: []domain.IndikatorPemda{},
			}
			firstRow = false
		}
		if indDbId == 0 || kodeIndikator == "" {
			continue
		}
		if _, exists := indMap[indDbId]; !exists {
			ind := domain.IndikatorPemda{
				Id: indDbId, SasaranPemdaId: id,
				KodeIndikator:       kodeIndikator,
				Indikator:           sql.NullString{String: indText, Valid: true},
				RumusPerhitungan:    sql.NullString{String: rumus, Valid: rumus != ""},
				SumberData:          sql.NullString{String: sumber, Valid: sumber != ""},
				DefinisiOperasional: sql.NullString{String: definisi, Valid: definisi != ""},
				Jenis:               "renstra",
				Target:              []domain.TargetPemda{},
			}
			// buat slot placeholder per tahun
			if tahunAwal != "" && tahunAkhir != "" {
				awal, _ := strconv.Atoi(tahunAwal)
				akhir, _ := strconv.Atoi(tahunAkhir)
				for y := awal; y <= akhir; y++ {
					ind.Target = append(ind.Target, domain.TargetPemda{
						KodeIndikator: kodeIndikator,
						Tahun:         strconv.Itoa(y),
						Jenis:         "renstra",
					})
				}
			}
			result.Indikator = append(result.Indikator, ind)
			indMap[indDbId] = &result.Indikator[len(result.Indikator)-1]
		}
		if targetDbId > 0 && targetTahun != "" {
			cur := indMap[indDbId]
			awal, _ := strconv.Atoi(tahunAwal)
			tInt, _ := strconv.Atoi(targetTahun)
			idx := tInt - awal
			if idx >= 0 && idx < len(cur.Target) {
				cur.Target[idx] = domain.TargetPemda{
					Id: targetDbId, KodeIndikator: kodeIndikator,
					Target: targetValue, Satuan: targetSatuan,
					Tahun: targetTahun, Jenis: targetJenis,
				}
			}
		}
	}
	if err := rows.Err(); err != nil {
		return domain.SasaranPemda{}, err
	}
	if result.Id == 0 {
		return domain.SasaranPemda{}, fmt.Errorf("sasaran pemda id %d tidak ditemukan", sasaranPemdaId)
	}
	sort.Slice(result.Indikator, func(i, j int) bool {
		return result.Indikator[i].Id < result.Indikator[j].Id
	})
	return result, nil
}

// ═══════════════════════════════════════════════════════════════════
// FIND ALL
// ═══════════════════════════════════════════════════════════════════
func (r *SasaranPemdaRepositoryImpl) FindAll(
	ctx context.Context, tx *sql.Tx, tahun string,
) ([]domain.SasaranPemda, error) {
	query := `
		SELECT
			sp.id, sp.subtema_id, sp.sasaran_pemda, sp.periode_id,
			COALESCE(p.tahun_awal,''), COALESCE(p.tahun_akhir,''),
			COALESCE(i.id, 0),
			COALESCE(i.kode_indikator,''),
			COALESCE(i.indikator,''),
			COALESCE(t.id, 0),
			COALESCE(t.target,''), COALESCE(t.satuan,''), COALESCE(t.tahun,'')
		FROM tb_sasaran_pemda sp
		LEFT JOIN tb_periode p ON sp.periode_id=p.id
		LEFT JOIN tb_indikator_matrix_pemda i
			ON sp.id=i.sasaran_pemda_id
			AND (i.jenis='renstra' OR i.jenis='' OR i.jenis IS NULL)
		LEFT JOIN tb_target_pemda t
			ON t.kode_indikator=i.kode_indikator
			AND (t.jenis='renstra' OR t.jenis='' OR t.jenis IS NULL)
		ORDER BY sp.id, i.id, CAST(t.tahun AS SIGNED)`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	sasaranMap := make(map[int]*domain.SasaranPemda)
	indMap := make(map[string]*domain.IndikatorPemda) // key: "spId:indId"
	for rows.Next() {
		var (
			id, subtemaId, periodeId               int
			sasaranText, tahunAwal, tahunAkhir     string
			indDbId                                int
			kodeIndikator, indText                 string
			targetDbId                             int
			targetValue, targetSatuan, targetTahun string
		)
		if err := rows.Scan(
			&id, &subtemaId, &sasaranText, &periodeId,
			&tahunAwal, &tahunAkhir,
			&indDbId, &kodeIndikator, &indText,
			&targetDbId, &targetValue, &targetSatuan, &targetTahun,
		); err != nil {
			return nil, err
		}
		if _, exists := sasaranMap[id]; !exists {
			sasaranMap[id] = &domain.SasaranPemda{
				Id: id, SubtemaId: subtemaId, SasaranPemda: sasaranText,
				PeriodeId: periodeId,
				Periode:   domain.Periode{TahunAwal: tahunAwal, TahunAkhir: tahunAkhir},
				Indikator: []domain.IndikatorPemda{},
			}
		}
		sp := sasaranMap[id]
		if indDbId == 0 || kodeIndikator == "" {
			continue
		}
		mapKey := fmt.Sprintf("%d:%d", id, indDbId)
		if _, exists := indMap[mapKey]; !exists {
			ind := domain.IndikatorPemda{
				Id: indDbId, SasaranPemdaId: id,
				KodeIndikator: kodeIndikator,
				Indikator:     sql.NullString{String: indText, Valid: true},
				Target:        []domain.TargetPemda{},
			}
			if tahunAwal != "" && tahunAkhir != "" {
				awal, _ := strconv.Atoi(tahunAwal)
				akhir, _ := strconv.Atoi(tahunAkhir)
				for y := awal; y <= akhir; y++ {
					ind.Target = append(ind.Target, domain.TargetPemda{
						KodeIndikator: kodeIndikator,
						Tahun:         strconv.Itoa(y),
					})
				}
			}
			sp.Indikator = append(sp.Indikator, ind)
			indMap[mapKey] = &sp.Indikator[len(sp.Indikator)-1]
		}
		if targetDbId > 0 && targetTahun != "" {
			cur := indMap[mapKey]
			for i := range cur.Target {
				if cur.Target[i].Tahun == targetTahun {
					cur.Target[i] = domain.TargetPemda{
						Id: targetDbId, KodeIndikator: kodeIndikator,
						Target: targetValue, Satuan: targetSatuan, Tahun: targetTahun,
					}
					break
				}
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	result := make([]domain.SasaranPemda, 0, len(sasaranMap))
	for _, sp := range sasaranMap {
		result = append(result, *sp)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Id < result[j].Id })
	return result, nil
}

// ═══════════════════════════════════════════════════════════════════
// FIND ALL BY TAHUN (dual layer)
// ═══════════════════════════════════════════════════════════════════
func (r *SasaranPemdaRepositoryImpl) FindAllByTahun(
	ctx context.Context, tx *sql.Tx, tahun, jenisPeriode, jenis string,
) ([]domain.SasaranPemda, error) {
	var jenisClause string
	var args []interface{}
	if jenis == "renstra" {
		jenisClause = "(t.jenis='renstra' OR t.jenis='' OR t.jenis IS NULL)"
		args = []interface{}{tahun, tahun, jenisPeriode}
	} else {
		jenisClause = "t.jenis=?"
		args = []interface{}{tahun, jenis, tahun, jenisPeriode}
	}
	query := fmt.Sprintf(`
		SELECT
			sp.id, sp.tujuan_pemda_id, sp.subtema_id, sp.sasaran_pemda, sp.periode_id,
			COALESCE(p.tahun_awal,''), COALESCE(p.tahun_akhir,''), COALESCE(p.jenis_periode,''),
			COALESCE(i.id, 0),
			COALESCE(i.kode_indikator,''),
			COALESCE(i.indikator,''),
			COALESCE(i.rumus_perhitungan,''), COALESCE(i.sumber_data,''),
			COALESCE(i.definisi_operasional,''),
			COALESCE(t.id, 0),
			COALESCE(t.target,''), COALESCE(t.satuan,''), COALESCE(t.tahun,''),
			COALESCE(t.jenis,'renstra')
		FROM tb_sasaran_pemda sp
		INNER JOIN tb_periode p ON sp.periode_id=p.id
		LEFT JOIN tb_indikator_matrix_pemda i
			ON sp.id=i.sasaran_pemda_id
			AND (i.jenis='renstra' OR i.jenis='' OR i.jenis IS NULL)
		LEFT JOIN tb_target_pemda t
			ON t.kode_indikator=i.kode_indikator
			AND t.tahun=?
			AND %s
		WHERE CAST(? AS SIGNED) BETWEEN CAST(p.tahun_awal AS SIGNED) AND CAST(p.tahun_akhir AS SIGNED)
		  AND p.jenis_periode=?
		ORDER BY sp.id, i.id`, jenisClause)
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("FindAllByTahun: %w", err)
	}
	defer rows.Close()
	sasaranMap := make(map[int]*domain.SasaranPemda)
	indMap := make(map[string]*domain.IndikatorPemda)
	for rows.Next() {
		var (
			spId, tujuanPemdaId, subtemaId, periodeId           int
			sasaranText, tahunAwal, tahunAkhir, jenisPeriodeCol string
			indDbId                                             int
			kodeIndikator, indText, rumus, sumber, definisi     string
			targetDbId                                          int
			targetValue, targetSatuan, targetTahun, targetJenis string
		)
		if err := rows.Scan(
			&spId, &tujuanPemdaId, &subtemaId, &sasaranText,
			&periodeId, &tahunAwal, &tahunAkhir, &jenisPeriodeCol,
			&indDbId, &kodeIndikator, &indText, &rumus, &sumber, &definisi,
			&targetDbId, &targetValue, &targetSatuan, &targetTahun, &targetJenis,
		); err != nil {
			return nil, err
		}
		if _, exists := sasaranMap[spId]; !exists {
			sasaranMap[spId] = &domain.SasaranPemda{
				Id: spId, TujuanPemdaId: tujuanPemdaId,
				SubtemaId: subtemaId, SasaranPemda: sasaranText,
				PeriodeId: periodeId,
				Periode:   domain.Periode{TahunAwal: tahunAwal, TahunAkhir: tahunAkhir, JenisPeriode: jenisPeriodeCol},
				Indikator: []domain.IndikatorPemda{},
			}
		}
		sp := sasaranMap[spId]
		if indDbId == 0 || kodeIndikator == "" {
			continue
		}
		mapKey := fmt.Sprintf("%d:%d", spId, indDbId)
		if _, exists := indMap[mapKey]; !exists {
			ind := domain.IndikatorPemda{
				Id: indDbId, SasaranPemdaId: spId,
				KodeIndikator:       kodeIndikator,
				Indikator:           sql.NullString{String: indText, Valid: true},
				RumusPerhitungan:    sql.NullString{String: rumus, Valid: rumus != ""},
				SumberData:          sql.NullString{String: sumber, Valid: sumber != ""},
				DefinisiOperasional: sql.NullString{String: definisi, Valid: definisi != ""},
				Target:              []domain.TargetPemda{},
			}
			sp.Indikator = append(sp.Indikator, ind)
			indMap[mapKey] = &sp.Indikator[len(sp.Indikator)-1]
		}
		if targetDbId > 0 {
			cur := indMap[mapKey]
			cur.Target = []domain.TargetPemda{{
				Id: targetDbId, KodeIndikator: kodeIndikator,
				Target: targetValue, Satuan: targetSatuan,
				Tahun: targetTahun, Jenis: targetJenis,
			}}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	result := make([]domain.SasaranPemda, 0, len(sasaranMap))
	for _, sp := range sasaranMap {
		result = append(result, *sp)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Id < result[j].Id })
	return result, nil
}

// ═══════════════════════════════════════════════════════════════════
// FIND ALL WITH POKIN
// ═══════════════════════════════════════════════════════════════════
func (r *SasaranPemdaRepositoryImpl) FindAllWithPokin(
	ctx context.Context, tx *sql.Tx, tahunAwal, tahunAkhir, jenisPeriode string,
) ([]domain.PohonKinerjaWithSasaran, error) {
	query := `
	WITH RECURSIVE pohon_hierarchy AS (
		SELECT pk.id, pk.nama_pohon, pk.parent, pk.level_pohon, pk.jenis_pohon,
		       pk.keterangan, pk.is_active, pk.tahun,
		       CAST(pk.id AS CHAR(50)) AS path, pk.id AS root_id, pk.nama_pohon AS root_nama
		FROM tb_pohon_kinerja pk
		WHERE pk.level_pohon=0
		  AND CAST(pk.tahun AS SIGNED) BETWEEN CAST(? AS SIGNED) AND CAST(? AS SIGNED)
		UNION ALL
		SELECT c.id, c.nama_pohon, c.parent, c.level_pohon, c.jenis_pohon,
		       c.keterangan, c.is_active, c.tahun,
		       CONCAT(ph.path,',',c.id), ph.root_id, ph.root_nama
		FROM tb_pohon_kinerja c
		JOIN pohon_hierarchy ph ON c.parent=ph.id
		WHERE CAST(c.tahun AS SIGNED) BETWEEN CAST(? AS SIGNED) AND CAST(? AS SIGNED)
	)
	SELECT
		pk.id, pk.nama_pohon, pk.jenis_pohon, pk.level_pohon,
		pk.keterangan,
		CASE
			WHEN sp.id IS NULL THEN pk.is_active
			ELSE CASE WHEN pk.is_active = true AND tematik.is_active = true THEN true ELSE false END
		END AS is_active,
		pk.tahun,
		pk.root_id, pk.root_nama,
		sp.id, sp.sasaran_pemda,
		sp.tahun_awal, sp.tahun_akhir, sp.jenis_periode,
		COALESCE(i.id, 0),
		COALESCE(i.kode_indikator,''),
		COALESCE(i.indikator,''),
		COALESCE(i.rumus_perhitungan,''), COALESCE(i.sumber_data,''),
		COALESCE(t.id, 0),
		COALESCE(t.target,''), COALESCE(t.satuan,''), COALESCE(t.tahun,'')
	FROM pohon_hierarchy pk
	LEFT JOIN tb_sasaran_pemda sp
		ON pk.id=sp.subtema_id
		AND sp.tahun_awal=? AND sp.tahun_akhir=? AND sp.jenis_periode=?
	LEFT JOIN tb_tujuan_pemda tp ON sp.tujuan_pemda_id=tp.id
	LEFT JOIN tb_pohon_kinerja tematik ON tp.tematik_id=tematik.id
	LEFT JOIN tb_indikator_matrix_pemda i
		ON sp.id=i.sasaran_pemda_id
		AND (i.jenis='renstra' OR i.jenis='' OR i.jenis IS NULL)
	LEFT JOIN tb_target_pemda t
		ON t.kode_indikator=i.kode_indikator
		AND (t.jenis='renstra' OR t.jenis='' OR t.jenis IS NULL)
		AND CAST(t.tahun AS SIGNED) BETWEEN CAST(? AS SIGNED) AND CAST(? AS SIGNED)
	WHERE pk.level_pohon BETWEEN 1 AND 3
	ORDER BY pk.root_id, pk.id, sp.id, i.id, CAST(t.tahun AS SIGNED)`
	rows, err := tx.QueryContext(ctx, query,
		tahunAwal, tahunAkhir,
		tahunAwal, tahunAkhir,
		tahunAwal, tahunAkhir, jenisPeriode,
		tahunAwal, tahunAkhir,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tematikMap := make(map[int]*domain.PohonKinerjaWithSasaran)
	for rows.Next() {
		var (
			subtematikId                              int
			namaSubtematik, jenisPohon, keterangan    string
			levelPohon                                int
			isActive                                  bool
			pohonTahun                                string
			tematikId                                 int
			namaTematik                               string
			idSasaran                                 sql.NullInt64
			sasaranText                               sql.NullString
			tahunAwalSp, tahunAkhirSp, jenisPeriodeSp sql.NullString
			indDbId                                   int
			kodeIndikator, indText, rumus, sumber     string
			targetDbId                                int
			targetValue, targetSatuan, targetTahun    string
		)
		if err := rows.Scan(
			&subtematikId, &namaSubtematik, &jenisPohon, &levelPohon,
			&keterangan, &isActive, &pohonTahun,
			&tematikId, &namaTematik,
			&idSasaran, &sasaranText,
			&tahunAwalSp, &tahunAkhirSp, &jenisPeriodeSp,
			&indDbId, &kodeIndikator, &indText, &rumus, &sumber,
			&targetDbId, &targetValue, &targetSatuan, &targetTahun,
		); err != nil {
			return nil, err
		}
		if _, exists := tematikMap[tematikId]; !exists {
			tematikMap[tematikId] = &domain.PohonKinerjaWithSasaran{
				TematikId: tematikId, NamaTematik: namaTematik, Tahun: pohonTahun,
				Subtematik: []domain.SubtematikWithSasaran{},
			}
		}
		tematik := tematikMap[tematikId]
		var foundSub *domain.SubtematikWithSasaran
		for i := range tematik.Subtematik {
			if tematik.Subtematik[i].SubtematikId == subtematikId {
				foundSub = &tematik.Subtematik[i]
				break
			}
		}
		if foundSub == nil {
			tematik.Subtematik = append(tematik.Subtematik, domain.SubtematikWithSasaran{
				SubtematikId: subtematikId, NamaSubtematik: namaSubtematik,
				JenisPohon: jenisPohon, LevelPohon: levelPohon,
				Tahun: pohonTahun, IsActive: isActive,
				SasaranPemdaList: []domain.SasaranPemdaDetail{},
			})
			foundSub = &tematik.Subtematik[len(tematik.Subtematik)-1]
		}
		if !idSasaran.Valid {
			continue
		}
		var foundSasaran *domain.SasaranPemdaDetail
		for i := range foundSub.SasaranPemdaList {
			if foundSub.SasaranPemdaList[i].Id == int(idSasaran.Int64) {
				foundSasaran = &foundSub.SasaranPemdaList[i]
				break
			}
		}
		if foundSasaran == nil {
			foundSub.SasaranPemdaList = append(foundSub.SasaranPemdaList, domain.SasaranPemdaDetail{
				Id: int(idSasaran.Int64), SasaranPemda: sasaranText.String,
				Periode: domain.Periode{
					TahunAwal: tahunAwal, TahunAkhir: tahunAkhir, JenisPeriode: jenisPeriode,
				},
				Indikator: []domain.IndikatorDetail{},
			})
			foundSasaran = &foundSub.SasaranPemdaList[len(foundSub.SasaranPemdaList)-1]
		}
		if indDbId == 0 || kodeIndikator == "" {
			continue
		}
		var foundInd *domain.IndikatorDetail
		for i := range foundSasaran.Indikator {
			if foundSasaran.Indikator[i].Id == indDbId {
				foundInd = &foundSasaran.Indikator[i]
				break
			}
		}
		if foundInd == nil {
			newInd := domain.IndikatorDetail{
				Id: indDbId, KodeIndikator: kodeIndikator, Indikator: indText,
				RumusPerhitungan: sql.NullString{String: rumus, Valid: rumus != ""},
				SumberData:       sql.NullString{String: sumber, Valid: sumber != ""},
				Target:           []domain.TargetDetail{},
			}
			awal, _ := strconv.Atoi(tahunAwal)
			akhir, _ := strconv.Atoi(tahunAkhir)
			for y := awal; y <= akhir; y++ {
				newInd.Target = append(newInd.Target, domain.TargetDetail{
					KodeIndikator: kodeIndikator, Tahun: strconv.Itoa(y),
				})
			}
			foundSasaran.Indikator = append(foundSasaran.Indikator, newInd)
			foundInd = &foundSasaran.Indikator[len(foundSasaran.Indikator)-1]
		}
		if targetDbId > 0 && targetTahun != "" {
			for i := range foundInd.Target {
				if foundInd.Target[i].Tahun == targetTahun {
					foundInd.Target[i] = domain.TargetDetail{
						Id: targetDbId, KodeIndikator: kodeIndikator,
						Target: targetValue, Satuan: targetSatuan, Tahun: targetTahun,
					}
					break
				}
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	result := make([]domain.PohonKinerjaWithSasaran, 0, len(tematikMap))
	for _, t := range tematikMap {
		result = append(result, *t)
	}
	return result, nil
}

// ═══════════════════════════════════════════════════════════════════
// UTILS
// ═══════════════════════════════════════════════════════════════════
func (r *SasaranPemdaRepositoryImpl) IsIdExists(ctx context.Context, tx *sql.Tx, id int) bool {
	var c int
	_ = tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM tb_sasaran_pemda WHERE id=?`, id,
	).Scan(&c)
	return c > 0
}
func (r *SasaranPemdaRepositoryImpl) IsSubtemaIdExists(ctx context.Context, tx *sql.Tx, subtemaId int) bool {
	var c int
	_ = tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM tb_sasaran_pemda WHERE subtema_id=?`, subtemaId,
	).Scan(&c)
	return c > 0
}
func (r *SasaranPemdaRepositoryImpl) UpdatePeriode(
	ctx context.Context, tx *sql.Tx, sp domain.SasaranPemda,
) (domain.SasaranPemda, error) {
	res, err := tx.ExecContext(ctx,
		`UPDATE tb_sasaran_pemda SET periode_id=? WHERE id=?`, sp.PeriodeId, sp.Id,
	)
	if err != nil {
		return domain.SasaranPemda{}, err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return domain.SasaranPemda{}, fmt.Errorf("sasaran pemda id %d tidak ditemukan", sp.Id)
	}
	return r.FindById(ctx, tx, sp.Id)
}

// ═══════════════════════════════════════════════════════════════════
// INDIKATOR & TARGET LAYER
// ═══════════════════════════════════════════════════════════════════
func (r *SasaranPemdaRepositoryImpl) FindIndikatorByKode(
	ctx context.Context, tx *sql.Tx, kodeIndikator string,
) (domain.IndikatorPemda, error) {
	var ind domain.IndikatorPemda
	var indText string
	err := tx.QueryRowContext(ctx,
		`SELECT id, kode_indikator, sasaran_pemda_id, COALESCE(indikator,'')
		 FROM tb_indikator_matrix_pemda
		 WHERE kode_indikator=?
		   AND (jenis='renstra' OR jenis='' OR jenis IS NULL)
		 LIMIT 1`, kodeIndikator,
	).Scan(&ind.Id, &ind.KodeIndikator, &ind.SasaranPemdaId, &indText)
	if err != nil {
		return domain.IndikatorPemda{}, err
	}
	ind.Indikator = sql.NullString{String: indText, Valid: true}
	return ind, nil
}
func (r *SasaranPemdaRepositoryImpl) FindTargetLayerById(
	ctx context.Context, tx *sql.Tx, id int,
) (domain.TargetPemda, error) {
	var t domain.TargetPemda
	err := tx.QueryRowContext(ctx,
		`SELECT id, kode_indikator,
			COALESCE(target,''), COALESCE(satuan,''), COALESCE(tahun,''),
			COALESCE(jenis,'renstra')
		 FROM tb_target_pemda WHERE id=?`, id,
	).Scan(&t.Id, &t.KodeIndikator, &t.Target, &t.Satuan, &t.Tahun, &t.Jenis)
	return t, err
}
func (r *SasaranPemdaRepositoryImpl) CreateTargetLayer(
	ctx context.Context, tx *sql.Tx, t domain.TargetPemda,
) (domain.TargetPemda, error) {
	jenis := t.Jenis
	if jenis == "" {
		jenis = "renstra"
	}
	res, err := tx.ExecContext(ctx,
		`INSERT INTO tb_target_pemda (kode_indikator, target, satuan, tahun, jenis)
		 VALUES (?, ?, ?, ?, ?)`,
		t.KodeIndikator, t.Target, t.Satuan, t.Tahun, jenis,
	)
	if err != nil {
		return domain.TargetPemda{}, err
	}
	newId, _ := res.LastInsertId()
	return r.FindTargetLayerById(ctx, tx, int(newId))
}
func (r *SasaranPemdaRepositoryImpl) UpdateTargetLayerById(
	ctx context.Context, tx *sql.Tx, id int, target, satuan string,
) (domain.TargetPemda, error) {
	if _, err := tx.ExecContext(ctx,
		`UPDATE tb_target_pemda SET target=?, satuan=? WHERE id=?`, target, satuan, id,
	); err != nil {
		return domain.TargetPemda{}, err
	}
	return r.FindTargetLayerById(ctx, tx, id)
}

// ═══════════════════════════════════════════════════════════════════
// UPSERT TARGET PEMDA — rankhir / penetapan
// key: kode_indikator + tahun + jenis
// ═══════════════════════════════════════════════════════════════════
func (r *SasaranPemdaRepositoryImpl) UpsertTargetPemda(
	ctx context.Context, tx *sql.Tx, t domain.TargetPemda,
) (domain.TargetPemda, error) {
	var existingId int
	err := tx.QueryRowContext(ctx,
		`SELECT id FROM tb_target_pemda
		 WHERE kode_indikator=? AND tahun=? AND jenis=? LIMIT 1`,
		t.KodeIndikator, t.Tahun, t.Jenis,
	).Scan(&existingId)
	if err == sql.ErrNoRows {
		res, err := tx.ExecContext(ctx,
			`INSERT INTO tb_target_pemda (kode_indikator, target, satuan, tahun, jenis)
			 VALUES (?, ?, ?, ?, ?)`,
			t.KodeIndikator, t.Target, t.Satuan, t.Tahun, t.Jenis,
		)
		if err != nil {
			return domain.TargetPemda{}, err
		}
		newId, _ := res.LastInsertId()
		return r.FindTargetLayerById(ctx, tx, int(newId))
	}
	if err != nil {
		return domain.TargetPemda{}, err
	}
	if _, err := tx.ExecContext(ctx,
		`UPDATE tb_target_pemda SET target=?, satuan=? WHERE id=?`,
		t.Target, t.Satuan, existingId,
	); err != nil {
		return domain.TargetPemda{}, err
	}
	return r.FindTargetLayerById(ctx, tx, existingId)
}

func (r *SasaranPemdaRepositoryImpl) FindRanwalByTahun(
	ctx context.Context, tx *sql.Tx, tahun, jenisPeriode string,
) ([]domain.SasaranPemda, error) {
	query := `
		SELECT
			sp.id, sp.tujuan_pemda_id, sp.subtema_id, sp.sasaran_pemda, sp.periode_id,
			COALESCE(p.tahun_awal,''), COALESCE(p.tahun_akhir,''), COALESCE(p.jenis_periode,''),
			COALESCE(pk.nama_pohon,'')  AS nama_subtema,
			COALESCE(tp.tujuan_pemda,'') AS tujuan_text,
			COALESCE(i.id, 0),
			COALESCE(i.kode_indikator,''),
			COALESCE(i.indikator,''),
			COALESCE(i.rumus_perhitungan,''),
			COALESCE(i.sumber_data,''),
			COALESCE(i.definisi_operasional,''),
			COALESCE(tr.id, tren.id, 0)                             AS target_id,
			COALESCE(tr.target, tren.target, '')                    AS target_value,
			COALESCE(tr.satuan, tren.satuan, '')                    AS target_satuan,
			CASE WHEN tr.id IS NOT NULL THEN 'ranwal' ELSE COALESCE(tren.jenis,'renstra') END AS target_jenis
		FROM tb_sasaran_pemda sp
		INNER JOIN tb_periode p ON sp.periode_id=p.id
		LEFT JOIN tb_pohon_kinerja pk ON sp.subtema_id=pk.id
		LEFT JOIN tb_tujuan_pemda tp ON sp.tujuan_pemda_id=tp.id
		LEFT JOIN tb_indikator_matrix_pemda i
			ON sp.id=i.sasaran_pemda_id
			AND (i.jenis='renstra' OR i.jenis='' OR i.jenis IS NULL)
		LEFT JOIN tb_target_pemda tren
			ON tren.kode_indikator=i.kode_indikator
			AND tren.tahun=?
			AND (tren.jenis='renstra' OR tren.jenis='' OR tren.jenis IS NULL)
		LEFT JOIN tb_target_pemda tr
			ON tr.kode_indikator=i.kode_indikator
			AND tr.tahun=?
			AND tr.jenis='ranwal'
		WHERE CAST(? AS SIGNED) BETWEEN CAST(p.tahun_awal AS SIGNED) AND CAST(p.tahun_akhir AS SIGNED)
		  AND p.jenis_periode=?
		ORDER BY sp.id, i.id`
	rows, err := tx.QueryContext(ctx, query, tahun, tahun, tahun, jenisPeriode)
	if err != nil {
		return nil, fmt.Errorf("FindRanwalByTahun: %w", err)
	}
	defer rows.Close()
	sasaranMap := make(map[int]*domain.SasaranPemda)
	indMap := make(map[string]*domain.IndikatorPemda)
	for rows.Next() {
		var (
			spId, tujuanPemdaId, subtemaId, periodeId           int
			sasaranText, tahunAwal, tahunAkhir, jenisPeriodeCol string
			namaSubtema, tujuanText                             string
			indDbId                                             int
			kodeIndikator, indText, rumus, sumber, definisi     string
			targetDbId                                          int
			targetValue, targetSatuan, targetJenis              string
		)
		if err := rows.Scan(
			&spId, &tujuanPemdaId, &subtemaId, &sasaranText,
			&periodeId, &tahunAwal, &tahunAkhir, &jenisPeriodeCol,
			&namaSubtema, &tujuanText,
			&indDbId, &kodeIndikator, &indText, &rumus, &sumber, &definisi,
			&targetDbId, &targetValue, &targetSatuan, &targetJenis,
		); err != nil {
			return nil, err
		}
		if _, exists := sasaranMap[spId]; !exists {
			sasaranMap[spId] = &domain.SasaranPemda{
				Id: spId, TujuanPemdaId: tujuanPemdaId,
				SubtemaId: subtemaId, SasaranPemda: sasaranText,
				NamaSubtema: namaSubtema, TujuanPemdaText: tujuanText,
				PeriodeId: periodeId,
				Periode:   domain.Periode{TahunAwal: tahunAwal, TahunAkhir: tahunAkhir, JenisPeriode: jenisPeriodeCol},
				Indikator: []domain.IndikatorPemda{},
			}
		}
		sp := sasaranMap[spId]
		if indDbId == 0 || kodeIndikator == "" {
			continue
		}
		mapKey := fmt.Sprintf("%d:%d", spId, indDbId)
		if _, exists := indMap[mapKey]; !exists {
			ind := domain.IndikatorPemda{
				Id: indDbId, SasaranPemdaId: spId,
				KodeIndikator:       kodeIndikator,
				Indikator:           sql.NullString{String: indText, Valid: true},
				RumusPerhitungan:    sql.NullString{String: rumus, Valid: rumus != ""},
				SumberData:          sql.NullString{String: sumber, Valid: sumber != ""},
				DefinisiOperasional: sql.NullString{String: definisi, Valid: definisi != ""},
				Target:              []domain.TargetPemda{},
			}
			sp.Indikator = append(sp.Indikator, ind)
			indMap[mapKey] = &sp.Indikator[len(sp.Indikator)-1]
		}
		if targetDbId > 0 {
			cur := indMap[mapKey]
			cur.Target = []domain.TargetPemda{{
				Id: targetDbId, KodeIndikator: kodeIndikator,
				Target: targetValue, Satuan: targetSatuan,
				Tahun: tahun, Jenis: targetJenis,
			}}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	result := make([]domain.SasaranPemda, 0, len(sasaranMap))
	for _, sp := range sasaranMap {
		result = append(result, *sp)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Id < result[j].Id })
	return result, nil
}
