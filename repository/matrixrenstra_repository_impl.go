package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"strings"
)

type MatrixRenstraRepositoryImpl struct{}

func NewMatrixRenstraRepositoryImpl() *MatrixRenstraRepositoryImpl {
	return &MatrixRenstraRepositoryImpl{}
}

//lama
// func (repository *MatrixRenstraRepositoryImpl) GetByKodeSubKegiatan(ctx context.Context, tx *sql.Tx, kodeOpd string, tahunAwal string, tahunAkhir string) ([]domain.SubKegiatanQuery, error) {

// 	checkQuery := `
//     SELECT COUNT(*)
//     FROM tb_subkegiatan_terpilih st
//     JOIN tb_rencana_kinerja rk ON st.rekin_id = rk.id
//     WHERE rk.kode_opd = ?
//     AND rk.tahun BETWEEN ? AND ?
//     `
// 	var count int
// 	err := tx.QueryRowContext(ctx, checkQuery, kodeOpd, tahunAwal, tahunAkhir).Scan(&count)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if count == 0 {
// 		return nil, fmt.Errorf("subkegiatan belum dipilih pada periode tahun %s sampai %s", tahunAwal, tahunAkhir)
// 	}

// 	// Perubahan utama vs versi lama:
// 	// 1. Tambah LEFT JOIN tb_pegawai → nama_pegawai (hilangkan N+1 di service)
// 	// 2. Tambah LEFT JOIN tb_pagu    → pagu_subkegiatan (sumber pagu baru, jenis='renstra')
// 	// 3. Hapus i.pagu_anggaran dari SELECT (pagu tidak lagi dari tb_indikator)
// 	query := `
//     WITH RECURSIVE hierarchy AS (
//         SELECT DISTINCT
//             u.kode_urusan,
//             u.nama_urusan,
//             bu.kode_bidang_urusan,
//             bu.nama_bidang_urusan,
//             p.kode_program,
//             p.nama_program,
//             k.kode_kegiatan,
//             k.nama_kegiatan,
//             s.kode_subkegiatan,
//             s.nama_subkegiatan,
//             rk.tahun AS tahun_subkegiatan,
//             rk.pegawai_id
//         FROM tb_subkegiatan_terpilih st
//         JOIN tb_rencana_kinerja rk ON st.rekin_id = rk.id
//         JOIN tb_subkegiatan s ON st.kode_subkegiatan = s.kode_subkegiatan
//         JOIN tb_master_kegiatan k
//             ON LEFT(s.kode_subkegiatan, LENGTH(k.kode_kegiatan)) = k.kode_kegiatan
//         JOIN tb_master_program p
//             ON LEFT(k.kode_kegiatan, LENGTH(p.kode_program)) = p.kode_program
//         JOIN tb_bidang_urusan bu
//             ON LEFT(p.kode_program, LENGTH(bu.kode_bidang_urusan)) = bu.kode_bidang_urusan
//         JOIN tb_urusan u
//             ON LEFT(bu.kode_bidang_urusan, LENGTH(u.kode_urusan)) = u.kode_urusan
//         WHERE rk.kode_opd = ?
//         AND rk.tahun BETWEEN ? AND ?
//     )
//     SELECT
//         h.kode_urusan,
//         h.nama_urusan,
//         h.kode_bidang_urusan,
//         h.nama_bidang_urusan,
//         h.kode_program,
//         h.nama_program,
//         h.kode_kegiatan,
//         h.nama_kegiatan,
//         h.kode_subkegiatan,
//         h.nama_subkegiatan,
//         h.tahun_subkegiatan,
//         h.pegawai_id,
//         COALESCE(pg.nama, '')    AS nama_pegawai,
//         COALESCE(tp.pagu, 0)     AS pagu_subkegiatan,
// 		im.kode_indikator        AS indikator_id,
// 		im.kode                  AS indikator_kode,
// 		im.indikator,
// 		im.tahun                 AS indikator_tahun,
// 		im.kode_opd              AS indikator_kode_opd,
//       	COALESCE(t.id, '')     AS target_id,
// 		COALESCE(t.target, '') AS target,
// 		COALESCE(t.satuan, '') AS satuan
//     FROM hierarchy h
//     LEFT JOIN tb_pegawai pg
//         ON pg.nip = h.pegawai_id
//     LEFT JOIN tb_pagu tp
//         ON  tp.kode_subkegiatan = h.kode_subkegiatan
//         AND tp.kode_opd         = ?
//         AND tp.jenis            = 'renstra'
//         AND tp.tahun            = h.tahun_subkegiatan
// 	LEFT JOIN tb_indikator_matrix im ON (
// 		im.kode = h.kode_urusan         OR
// 		im.kode = h.kode_bidang_urusan  OR
// 		im.kode = h.kode_program        OR
// 		im.kode = h.kode_kegiatan       OR
// 		im.kode = h.kode_subkegiatan
// 	)
// 		AND im.kode_opd = ?
// 		AND im.jenis    = 'renstra'
// 		AND im.tahun BETWEEN ? AND ?
// 	LEFT JOIN tb_target t ON t.indikator_id = im.kode_indikator
//     ORDER BY
//         h.kode_urusan,
//         h.kode_bidang_urusan,
//         h.kode_program,
//         h.kode_kegiatan,
//         h.kode_subkegiatan,
//         im.tahun
//     `

// 	rows, err := tx.QueryContext(ctx, query,
// 		kodeOpd, tahunAwal, tahunAkhir, // params CTE hierarchy
// 		kodeOpd,                        // params JOIN tb_pagu (kode_opd)
// 		kodeOpd, tahunAwal, tahunAkhir, // params JOIN tb_indikator
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var result []domain.SubKegiatanQuery
// 	for rows.Next() {
// 		var data domain.SubKegiatanQuery
// 		var (
// 			indikatorId, indikatorKode, indikator sql.NullString
// 			indikatorTahun, indikatorKodeOpd      sql.NullString
// 			targetId, target, satuan              sql.NullString
// 		)

// 		err := rows.Scan(
// 			&data.KodeUrusan,
// 			&data.NamaUrusan,
// 			&data.KodeBidangUrusan,
// 			&data.NamaBidangUrusan,
// 			&data.KodeProgram,
// 			&data.NamaProgram,
// 			&data.KodeKegiatan,
// 			&data.NamaKegiatan,
// 			&data.KodeSubKegiatan,
// 			&data.NamaSubKegiatan,
// 			&data.TahunSubKegiatan,
// 			&data.PegawaiId,
// 			&data.NamaPegawai,     // dari JOIN tb_pegawai
// 			&data.PaguSubKegiatan, // dari JOIN tb_pagu jenis='renstra'
// 			&indikatorId,
// 			&indikatorKode,
// 			&indikator,
// 			&indikatorTahun,
// 			&indikatorKodeOpd,
// 			&targetId,
// 			&target,
// 			&satuan,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if indikatorId.Valid {
// 			data.IndikatorId = indikatorId.String
// 			data.IndikatorKode = indikatorKode.String
// 			data.Indikator = indikator.String
// 			data.IndikatorTahun = indikatorTahun.String
// 			data.IndikatorKodeOpd = indikatorKodeOpd.String
// 		}
// 		if target.Valid {
// 			data.TargetId = targetId.String
// 			data.Target = target.String
// 			data.Satuan = satuan.String
// 		}

// 		result = append(result, data)
// 	}

// 	return result, nil
// }

// baru
func (repository *MatrixRenstraRepositoryImpl) GetByKodeSubKegiatan(ctx context.Context, tx *sql.Tx, kodeOpd string, tahunAwal string, tahunAkhir string) ([]domain.SubKegiatanQuery, error) {
	checkQuery := `
    SELECT COUNT(*)
    FROM tb_subkegiatan_terpilih st
    JOIN tb_rencana_kinerja rk ON st.rekin_id = rk.id
    WHERE rk.kode_opd = ?
    AND rk.tahun BETWEEN ? AND ?
    `
	var count int
	err := tx.QueryRowContext(ctx, checkQuery, kodeOpd, tahunAwal, tahunAkhir).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("subkegiatan belum dipilih pada periode tahun %s sampai %s", tahunAwal, tahunAkhir)
	}
	hierarchyQuery := `
    WITH hierarchy AS (
        SELECT DISTINCT
            u.kode_urusan,
            u.nama_urusan,
            bu.kode_bidang_urusan,
            bu.nama_bidang_urusan,
            p.kode_program,
            p.nama_program,
            k.kode_kegiatan,
            k.nama_kegiatan,
            s.kode_subkegiatan,
            s.nama_subkegiatan,
            rk.tahun AS tahun_subkegiatan,
            rk.pegawai_id
        FROM tb_subkegiatan_terpilih st
        JOIN tb_rencana_kinerja rk ON st.rekin_id = rk.id
        JOIN tb_subkegiatan s ON st.kode_subkegiatan = s.kode_subkegiatan
        JOIN tb_master_kegiatan k
            ON LEFT(s.kode_subkegiatan, LENGTH(k.kode_kegiatan)) = k.kode_kegiatan
        JOIN tb_master_program p
            ON LEFT(k.kode_kegiatan, LENGTH(p.kode_program)) = p.kode_program
        JOIN tb_bidang_urusan bu
            ON LEFT(p.kode_program, LENGTH(bu.kode_bidang_urusan)) = bu.kode_bidang_urusan
        JOIN tb_urusan u
            ON LEFT(bu.kode_bidang_urusan, LENGTH(u.kode_urusan)) = u.kode_urusan
        WHERE rk.kode_opd = ?
        AND rk.tahun BETWEEN ? AND ?
    )
    SELECT
        h.kode_urusan,
        h.nama_urusan,
        h.kode_bidang_urusan,
        h.nama_bidang_urusan,
        h.kode_program,
        h.nama_program,
        h.kode_kegiatan,
        h.nama_kegiatan,
        h.kode_subkegiatan,
        h.nama_subkegiatan,
        h.tahun_subkegiatan,
        h.pegawai_id,
        COALESCE(pg.nama, '') AS nama_pegawai,
        COALESCE(tp.pagu, 0)   AS pagu_subkegiatan
    FROM hierarchy h
    LEFT JOIN tb_pegawai pg ON pg.nip = h.pegawai_id
    LEFT JOIN tb_pagu tp
        ON tp.kode_subkegiatan = h.kode_subkegiatan
        AND tp.kode_opd = ?
        AND tp.jenis = 'renstra'
        AND tp.tahun = h.tahun_subkegiatan
    ORDER BY
        h.kode_urusan,
        h.kode_bidang_urusan,
        h.kode_program,
        h.kode_kegiatan,
        h.kode_subkegiatan,
        h.tahun_subkegiatan
    `
	hRows, err := tx.QueryContext(ctx, hierarchyQuery, kodeOpd, tahunAwal, tahunAkhir, kodeOpd)
	if err != nil {
		return nil, err
	}
	defer hRows.Close()
	var baseRows []domain.SubKegiatanQuery
	for hRows.Next() {
		var row domain.SubKegiatanQuery
		if err := hRows.Scan(
			&row.KodeUrusan,
			&row.NamaUrusan,
			&row.KodeBidangUrusan,
			&row.NamaBidangUrusan,
			&row.KodeProgram,
			&row.NamaProgram,
			&row.KodeKegiatan,
			&row.NamaKegiatan,
			&row.KodeSubKegiatan,
			&row.NamaSubKegiatan,
			&row.TahunSubKegiatan,
			&row.PegawaiId,
			&row.NamaPegawai,
			&row.PaguSubKegiatan,
		); err != nil {
			return nil, err
		}
		baseRows = append(baseRows, row)
	}
	if err := hRows.Err(); err != nil {
		return nil, err
	}
	if len(baseRows) == 0 {
		return []domain.SubKegiatanQuery{}, nil
	}
	type indRow struct {
		KodeIndikator string
		Kode          string
		KodeOpd       string
		Indikator     string
		Tahun         string
		TargetId      string
		Target        string
		Satuan        string
	}
	indQuery := `
		SELECT
			im.kode_indikator,
			im.kode,
			im.kode_opd,
			im.indikator,
			COALESCE(im.tahun, '') AS indikator_tahun,
			COALESCE(t.id, '')     AS target_id,
			COALESCE(t.target, '') AS target,
			COALESCE(t.satuan, '') AS satuan
		FROM tb_indikator_matrix im
		LEFT JOIN tb_target t ON t.indikator_id = im.kode_indikator
		WHERE im.kode_opd = ?
		  AND im.jenis = 'renstra'
		  AND im.tahun BETWEEN ? AND ?
		ORDER BY im.kode, im.tahun, im.kode_indikator
	`
	rows, err := tx.QueryContext(ctx, indQuery, kodeOpd, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	renstraByKode := make(map[string][]indRow)
	for rows.Next() {
		var r indRow
		if err := rows.Scan(
			&r.KodeIndikator, &r.Kode, &r.KodeOpd, &r.Indikator, &r.Tahun,
			&r.TargetId, &r.Target, &r.Satuan,
		); err != nil {
			return nil, err
		}
		renstraByKode[r.Kode] = append(renstraByKode[r.Kode], r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	result := make([]domain.SubKegiatanQuery, 0, len(baseRows)*4)
	for _, base := range baseRows {
		kodes := []string{
			base.KodeUrusan,
			base.KodeBidangUrusan,
			base.KodeProgram,
			base.KodeKegiatan,
			base.KodeSubKegiatan,
		}
		totalInd := 0
		for _, k := range kodes {
			if k == "" {
				continue
			}
			for _, ind := range renstraByKode[k] {
				totalInd++
				row := base
				row.IndikatorId = ind.KodeIndikator
				row.IndikatorKode = ind.Kode
				row.Indikator = ind.Indikator
				row.IndikatorTahun = ind.Tahun
				row.IndikatorKodeOpd = ind.KodeOpd
				if ind.TargetId != "" {
					row.TargetId = ind.TargetId
					row.Target = ind.Target
					row.Satuan = ind.Satuan
				}
				result = append(result, row)
			}
		}
		if totalInd == 0 {
			result = append(result, base)
		}
	}
	// Pagu per tahun dari tb_pagu (bukan hanya tahun yang punya baris rekin di hierarchy).
	kodeSubSet := make(map[string]struct{})
	for _, b := range baseRows {
		if b.KodeSubKegiatan != "" {
			kodeSubSet[b.KodeSubKegiatan] = struct{}{}
		}
	}
	if len(kodeSubSet) > 0 {
		placeholders := make([]string, 0, len(kodeSubSet))
		args := make([]interface{}, 0, 2+len(kodeSubSet))
		args = append(args, kodeOpd, tahunAwal, tahunAkhir)
		for k := range kodeSubSet {
			placeholders = append(placeholders, "?")
			args = append(args, k)
		}
		paguQuery := fmt.Sprintf(`
			SELECT kode_subkegiatan, tahun, pagu
			FROM tb_pagu
			WHERE kode_opd = ?
			  AND jenis = 'renstra'
			  AND tahun BETWEEN ? AND ?
			  AND kode_subkegiatan IN (%s)
		`, strings.Join(placeholders, ","))
		pRows, err := tx.QueryContext(ctx, paguQuery, args...)
		if err != nil {
			return nil, err
		}
		type paguRow struct {
			kodeSub string
			tahun   string
			pagu    int64
		}
		var extraPagu []paguRow
		for pRows.Next() {
			var pr paguRow
			if err := pRows.Scan(&pr.kodeSub, &pr.tahun, &pr.pagu); err != nil {
				pRows.Close()
				return nil, err
			}
			extraPagu = append(extraPagu, pr)
		}
		if err := pRows.Close(); err != nil {
			return nil, err
		}
		seenPair := make(map[string]struct{}, len(result))
		templateBySub := make(map[string]domain.SubKegiatanQuery)
		for _, row := range result {
			if row.KodeSubKegiatan == "" {
				continue
			}
			seenPair[row.KodeSubKegiatan+"\x00"+row.TahunSubKegiatan] = struct{}{}
			if _, ok := templateBySub[row.KodeSubKegiatan]; !ok {
				templateBySub[row.KodeSubKegiatan] = row
			}
		}
		for _, pr := range extraPagu {
			key := pr.kodeSub + "\x00" + pr.tahun
			if _, ok := seenPair[key]; ok {
				continue
			}
			tmpl, ok := templateBySub[pr.kodeSub]
			if !ok {
				continue
			}
			synth := tmpl
			synth.TahunSubKegiatan = pr.tahun
			synth.PaguSubKegiatan = pr.pagu
			synth.IndikatorId = ""
			synth.IndikatorKode = ""
			synth.Indikator = ""
			synth.IndikatorTahun = ""
			synth.IndikatorKodeOpd = ""
			synth.TargetId = ""
			synth.Target = ""
			synth.Satuan = ""
			result = append(result, synth)
			seenPair[key] = struct{}{}
		}
	}
	return result, nil
}

func (r *MatrixRenstraRepositoryImpl) DeleteIndikator(ctx context.Context, tx *sql.Tx, kodeIndikator string) error {
	_, err := tx.ExecContext(ctx, `DELETE FROM tb_indikator_matrix WHERE kode_indikator = ?`, kodeIndikator)
	return err
}

func (r *MatrixRenstraRepositoryImpl) DeleteTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) error {
	// indikatorId = kode_indikator dari tb_indikator_matrix
	_, err := tx.ExecContext(ctx, `DELETE FROM tb_target WHERE indikator_id = ?`, indikatorId)
	return err
}
func (repository *MatrixRenstraRepositoryImpl) UpsertAnggaran(ctx context.Context, tx *sql.Tx, kodeSubkegiatan, kodeOpd, tahun string, pagu int64) error {
	query := `
        INSERT INTO tb_pagu (kode_subkegiatan, kode_opd, tahun, jenis, pagu)
        VALUES (?, ?, ?, 'renstra', ?)
        ON DUPLICATE KEY UPDATE pagu = VALUES(pagu)
    `
	_, err := tx.ExecContext(ctx, query, kodeSubkegiatan, kodeOpd, tahun, pagu)
	return err
}

func (r *MatrixRenstraRepositoryImpl) UpsertIndikator(ctx context.Context, tx *sql.Tx, ind domain.Indikator) error {
	query := `
        INSERT INTO tb_indikator_matrix
            (kode_indikator, kode, kode_opd, indikator, tahun, jenis)
        VALUES (?, ?, ?, ?, ?, 'renstra')
        ON DUPLICATE KEY UPDATE
            indikator = VALUES(indikator),
            tahun     = VALUES(tahun)
    `
	_, err := tx.ExecContext(ctx, query,
		ind.KodeIndikator,
		ind.Kode,
		ind.KodeOpd,
		ind.Indikator,
		ind.Tahun,
	)
	return err
}
func (r *MatrixRenstraRepositoryImpl) UpsertTarget(ctx context.Context, tx *sql.Tx, t domain.Target) error {
	query := `
        INSERT INTO tb_target (id, indikator_id, target, satuan, tahun)
        VALUES (?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            target  = VALUES(target),
            satuan  = VALUES(satuan),
            tahun   = VALUES(tahun)
    `
	_, err := tx.ExecContext(ctx, query,
		t.Id,
		t.IndikatorId, // = kode_indikator dari indikator
		t.Target,
		t.Satuan,
		t.Tahun,
	)
	return err
}

func (r *MatrixRenstraRepositoryImpl) FindIndikatorByKodeIndikator(ctx context.Context, tx *sql.Tx, kodeIndikator string) (domain.Indikator, error) {
	query := `
        SELECT
            im.kode_indikator,
            im.kode,
            im.kode_opd,
            im.indikator,
            im.tahun,
            COALESCE(im.rumus_perhitungan, '')    AS rumus_perhitungan,
            COALESCE(im.sumber_data, '')          AS sumber_data,
            COALESCE(im.definisi_operasional, '') AS definisi_operasional,
            COALESCE(t.id, '')                    AS target_id,
            COALESCE(t.target, '')                AS target,
            COALESCE(t.satuan, '')                AS satuan,
            COALESCE(t.tahun, '')                 AS target_tahun
        FROM tb_indikator_matrix im
        LEFT JOIN tb_target t ON t.indikator_id = im.kode_indikator
        WHERE im.kode_indikator = ?
          AND im.jenis = 'renstra'
    `
	rows, err := tx.QueryContext(ctx, query, kodeIndikator)
	if err != nil {
		return domain.Indikator{}, err
	}
	defer rows.Close()
	var ind domain.Indikator
	found := false
	for rows.Next() {
		var (
			targetId, targetVal, targetSatuan, targetTahun string
			rumus, sumber, definisi                        string
		)
		err := rows.Scan(
			&ind.KodeIndikator,
			&ind.Kode,
			&ind.KodeOpd,
			&ind.Indikator,
			&ind.Tahun,
			&rumus, &sumber, &definisi,
			&targetId, &targetVal, &targetSatuan, &targetTahun,
		)
		if err != nil {
			return domain.Indikator{}, err
		}
		ind.RumusPerhitungan = sql.NullString{String: rumus, Valid: rumus != ""}
		ind.SumberData = sql.NullString{String: sumber, Valid: sumber != ""}
		ind.DefinisiOperasional = sql.NullString{String: definisi, Valid: definisi != ""}
		found = true
		if targetId != "" {
			ind.Target = append(ind.Target, domain.Target{
				Id:          targetId,
				IndikatorId: ind.KodeIndikator,
				Target:      targetVal,
				Satuan:      targetSatuan,
				Tahun:       targetTahun,
			})
		}
	}
	if !found {
		return domain.Indikator{}, sql.ErrNoRows
	}
	return ind, nil
}

func (r *MatrixRenstraRepositoryImpl) CountKodeIndikatorByPrefix(ctx context.Context, tx *sql.Tx, prefix string) (int, error) {
	query := `SELECT COUNT(*) FROM tb_indikator_matrix WHERE kode_indikator LIKE ? AND jenis = 'renstra'`
	var count int
	err := tx.QueryRowContext(ctx, query, prefix+"%").Scan(&count)
	return count, err
}

func (r *MatrixRenstraRepositoryImpl) DeleteIndicatorsExcept(
	ctx context.Context, tx *sql.Tx,
	kode, kodeOpd, tahun string, keepList []string,
) error {
	if len(keepList) == 0 {
		// Hapus semua target dalam scope ini
		delTarget := `
            DELETE FROM tb_target
            WHERE indikator_id IN (
                SELECT kode_indikator FROM tb_indikator_matrix
                WHERE kode = ? AND kode_opd = ? AND tahun = ? AND jenis = 'renstra'
            )
        `
		if _, err := tx.ExecContext(ctx, delTarget, kode, kodeOpd, tahun); err != nil {
			return err
		}
		// Hapus semua indikator dalam scope ini
		delInd := `
            DELETE FROM tb_indikator_matrix
            WHERE kode = ? AND kode_opd = ? AND tahun = ? AND jenis = 'renstra'
        `
		_, err := tx.ExecContext(ctx, delInd, kode, kodeOpd, tahun)
		return err
	}
	// Bangun placeholder IN (?,?,...)
	placeholders := make([]string, len(keepList))
	args := make([]interface{}, 0, 3+len(keepList))
	args = append(args, kode, kodeOpd, tahun)
	for i, k := range keepList {
		placeholders[i] = "?"
		args = append(args, k)
	}
	inClause := strings.Join(placeholders, ",")
	// Hapus target dari indikator yang TIDAK ada di keepList
	delTarget := fmt.Sprintf(`
        DELETE FROM tb_target
        WHERE indikator_id IN (
            SELECT kode_indikator FROM tb_indikator_matrix
            WHERE kode = ? AND kode_opd = ? AND tahun = ? AND jenis = 'renstra'
              AND kode_indikator NOT IN (%s)
        )
    `, inClause)
	if _, err := tx.ExecContext(ctx, delTarget, args...); err != nil {
		return err
	}
	// Hapus indikator yang TIDAK ada di keepList
	delInd := fmt.Sprintf(`
        DELETE FROM tb_indikator_matrix
        WHERE kode = ? AND kode_opd = ? AND tahun = ? AND jenis = 'renstra'
          AND kode_indikator NOT IN (%s)
    `, inClause)
	_, err := tx.ExecContext(ctx, delInd, args...)
	return err
}

func (repo *MatrixRenstraRepositoryImpl) FindIndikatorRenstra(ctx context.Context, tx *sql.Tx, kodeOpd string, tahunAwal string, tahunAkhir string) ([]domain.Indikator, error) {
	query := `
        SELECT
            im.kode_indikator,
            im.kode,
            im.kode_opd,
            im.indikator,
            im.tahun,
            COALESCE(t.id, '')                    AS target_id,
            COALESCE(t.target, '')                AS target,
            COALESCE(t.satuan, '')                AS satuan,
            COALESCE(t.tahun, '')                 AS target_tahun
        FROM tb_indikator_matrix im
        LEFT JOIN tb_target t ON t.indikator_id = im.kode_indikator
        WHERE im.kode_indikator != ''
          AND im.jenis = 'renstra'
          AND im.kode_opd = ?
	  AND im.tahun BETWEEN ? AND ?
    `
	rows, err := tx.QueryContext(ctx, query, kodeOpd, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.Indikator

	for rows.Next() {
		var r domain.Indikator
		var (
			targetId, targetVal, targetSatuan, targetTahun string
		)
		err := rows.Scan(
			&r.KodeIndikator,
			&r.Kode,
			&r.KodeOpd,
			&r.Indikator,
			&r.Tahun,
			&targetId, &targetVal, &targetSatuan, &targetTahun,
		)
		if err != nil {
			return nil, err
		}
		if targetId != "" {
			r.Target = append(r.Target, domain.Target{
				Id:          targetId,
				IndikatorId: r.KodeIndikator,
				Target:      targetVal,
				Satuan:      targetSatuan,
				Tahun:       targetTahun,
			})
		}
		result = append(result, r)
	}

	return result, nil
}

func (repo *MatrixRenstraRepositoryImpl) FindIndikatorLama(ctx context.Context, tx *sql.Tx, kodeOpd string, tahunAwal string, tahunAkhir string) ([]domain.Indikator, error) {
	query := `
        SELECT
            im.id,
            im.kode,
            im.kode_opd,
            im.indikator,
            im.tahun,
            COALESCE(t.id, '')                    AS target_id,
            COALESCE(t.target, '')                AS target,
            COALESCE(t.satuan, '')                AS satuan,
            COALESCE(t.tahun, '')                 AS target_tahun
        FROM tb_indikator im
        LEFT JOIN tb_target t ON t.indikator_id = im.id
        WHERE im.kode != ''
          AND im.kode_opd = ?
	  AND im.tahun BETWEEN ? AND ?
    `
	rows, err := tx.QueryContext(ctx, query, kodeOpd, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.Indikator

	for rows.Next() {
		var r domain.Indikator
		var (
			targetId, targetVal, targetSatuan, targetTahun string
		)
		err := rows.Scan(
			&r.KodeIndikator,
			&r.Kode,
			&r.KodeOpd,
			&r.Indikator,
			&r.Tahun,
			&targetId, &targetVal, &targetSatuan, &targetTahun,
		)
		if err != nil {
			return nil, err
		}
		if targetId != "" {
			r.Target = append(r.Target, domain.Target{
				Id:          targetId,
				IndikatorId: r.KodeIndikator,
				Target:      targetVal,
				Satuan:      targetSatuan,
				Tahun:       targetTahun,
			})
		}
		result = append(result, r)
	}

	return result, nil
}
