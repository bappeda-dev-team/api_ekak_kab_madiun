package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"strings"
)

type MatrixRenjaRepositoryImpl struct {
}

func NewMatrixRenjaRepositoryImpl() *MatrixRenjaRepositoryImpl {
	return &MatrixRenjaRepositoryImpl{}
}

// GetRenja — ranwal view: indikator dari renstra, target dari ranwal (fallback: renstra), pagu dari tb_pagu jenis='renstra'.
// Opsi B: tampilkan semua subkegiatan dari periode renstra yg memiliki pagu untuk tahun ini,
// walau tidak ada rekin di tahun ini secara spesifik.
func (repository *MatrixRenjaRepositoryImpl) GetRenja(ctx context.Context, tx *sql.Tx, kodeOpd, tahun, _ string) ([]domain.SubKegiatanQuery, error) {
	// Cek apakah OPD sudah pernah memilih subkegiatan (any year)
	var count int
	if err := tx.QueryRowContext(ctx, `
        SELECT COUNT(*)
        FROM tb_subkegiatan_terpilih st
        JOIN tb_rencana_kinerja rk ON st.rekin_id = rk.id
        WHERE rk.kode_opd = ?`,
		kodeOpd,
	).Scan(&count); err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("subkegiatan belum dipilih untuk opd ini")
	}

	// ── Main query: hierarki + pagu (tanpa inline indikator JOIN) ──
	query := `
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
            rk.tahun     AS tahun_subkegiatan,
            rk.pegawai_id
        FROM tb_subkegiatan_terpilih st
        JOIN tb_rencana_kinerja rk  ON st.rekin_id        = rk.id
        JOIN tb_subkegiatan s       ON st.kode_subkegiatan = s.kode_subkegiatan
        JOIN tb_master_kegiatan k
            ON LEFT(s.kode_subkegiatan, LENGTH(k.kode_kegiatan)) = k.kode_kegiatan
        JOIN tb_master_program p
            ON LEFT(k.kode_kegiatan, LENGTH(p.kode_program)) = p.kode_program
        JOIN tb_bidang_urusan bu
            ON LEFT(p.kode_program, LENGTH(bu.kode_bidang_urusan)) = bu.kode_bidang_urusan
        JOIN tb_urusan u
            ON LEFT(bu.kode_bidang_urusan, LENGTH(u.kode_urusan)) = u.kode_urusan
        WHERE rk.kode_opd = ?
          AND rk.tahun    = ?
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
        COALESCE(pg.nama, '')    AS nama_pegawai,
        COALESCE(tp_rw.pagu, 0)  AS pagu_ranwal
    FROM hierarchy h
    LEFT JOIN tb_pegawai pg
        ON pg.nip = h.pegawai_id
    LEFT JOIN tb_pagu tp_rw
        ON  tp_rw.kode_subkegiatan = h.kode_subkegiatan
        AND tp_rw.kode_opd         = ?
        AND tp_rw.jenis            = 'renstra'
        AND tp_rw.tahun            = h.tahun_subkegiatan
    ORDER BY
        h.kode_urusan,
        h.kode_bidang_urusan,
        h.kode_program,
        h.kode_kegiatan,
        h.kode_subkegiatan
    `
	rows, err := tx.QueryContext(ctx, query, kodeOpd, tahun, kodeOpd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var baseResult []domain.SubKegiatanQuery
	for rows.Next() {
		var data domain.SubKegiatanQuery
		if err := rows.Scan(
			&data.KodeUrusan,
			&data.NamaUrusan,
			&data.KodeBidangUrusan,
			&data.NamaBidangUrusan,
			&data.KodeProgram,
			&data.NamaProgram,
			&data.KodeKegiatan,
			&data.NamaKegiatan,
			&data.KodeSubKegiatan,
			&data.NamaSubKegiatan,
			&data.TahunSubKegiatan,
			&data.PegawaiId,
			&data.NamaPegawai,
			&data.PaguRanwal,
		); err != nil {
			return nil, err
		}
		baseResult = append(baseResult, data)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// ── Separate indikator query (seperti pola GetByKodeSubKegiatanVersiKedua) ──
	indByKode, err := repository.fetchIndikatorRenjaPeriod(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	applyRanwalInd := func(row *domain.SubKegiatanQuery, ind indRowRenja) {
		row.TargetRanwalId = ind.targetRanwalId
		row.TargetRanwal = ind.targetRanwal
		row.SatuanRanwal = ind.satuanRanwal
		row.TargetId = ind.targetRanwalId
		row.Target = ind.targetRanwal
		row.Satuan = ind.satuanRanwal
	}
	// Merge indikator ke baris hierarki (satu baris per indikator yang cocok)
	result := mergeIndikator(baseResult, indByKode, applyRanwalInd)

	// ── Opsi B: Extra synthesis — subkegiatan renstra yg tidak ada rekin tahun ini ──
	seenSubkeg := make(map[string]bool)
	for _, row := range result {
		seenSubkeg[row.KodeSubKegiatan] = true
	}
	templates, err := repository.fetchAllSubkegTemplates(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	extraPagu, err := repository.fetchExtraPagu(ctx, tx, kodeOpd, tahun, "renstra")
	if err != nil {
		return nil, err
	}
	for kodeSub, pagu := range extraPagu {
		if seenSubkeg[kodeSub] {
			continue
		}
		tmpl, ok := templates[kodeSub]
		if !ok {
			continue
		}
		synth := tmpl
		synth.TahunSubKegiatan = tahun
		synth.PaguRanwal = pagu
		// Indikator renstra tetap ditampilkan walau tidak ada rekin tahun ini
		result = append(result, mergeIndikator([]domain.SubKegiatanQuery{synth}, indByKode, applyRanwalInd)...)
	}

	return result, nil
}

// GetRenjaRankhir — rankhir view: indikator dari renstra, dual-slot target (ranwal + rankhir), pagu_ranwal + pagu_rankhir flat.
// Opsi B: tampilkan semua subkegiatan dari periode renstra yg memiliki pagu untuk tahun ini.
func (repository *MatrixRenjaRepositoryImpl) GetRenjaRankhir(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string) ([]domain.SubKegiatanQuery, error) {
	// Cek apakah OPD sudah pernah memilih subkegiatan (any year)
	var count int
	if err := tx.QueryRowContext(ctx, `
        SELECT COUNT(*)
        FROM tb_subkegiatan_terpilih st
        JOIN tb_rencana_kinerja rk ON st.rekin_id = rk.id
        WHERE rk.kode_opd = ?`,
		kodeOpd,
	).Scan(&count); err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("subkegiatan belum dipilih untuk opd ini")
	}

	// ── Main query: hierarki + dual pagu (tanpa inline indikator JOIN) ──
	query := `
    WITH subkeg_pagu AS (
        SELECT
            sub.kode_subkegiatan,
            rk.pegawai_id,
            rk.tahun AS tahun_rekin,
            COALESCE(SUM(rb.anggaran), 0) AS pagu_rankhir
        FROM tb_subkegiatan_terpilih sub
        JOIN tb_rencana_kinerja rk ON sub.rekin_id = rk.id
        LEFT JOIN tb_rencana_aksi ra ON ra.rencana_kinerja_id = rk.id
        LEFT JOIN tb_rincian_belanja rb ON rb.renaksi_id = ra.id
        WHERE rk.kode_opd = ?
          AND rk.tahun = ?
        GROUP BY sub.kode_subkegiatan, rk.pegawai_id, rk.tahun
    )
    SELECT
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
        sp.tahun_rekin      AS tahun_subkegiatan,
        sp.pegawai_id,
        COALESCE(pg.nama,'') AS nama_pegawai,
        COALESCE(tp_rw.pagu, 0) AS pagu_ranwal,
        sp.pagu_rankhir
    FROM subkeg_pagu sp
    JOIN tb_subkegiatan s ON sp.kode_subkegiatan = s.kode_subkegiatan
    JOIN tb_master_kegiatan k
        ON LEFT(s.kode_subkegiatan, LENGTH(k.kode_kegiatan)) = k.kode_kegiatan
    JOIN tb_master_program p
        ON LEFT(k.kode_kegiatan, LENGTH(p.kode_program)) = p.kode_program
    JOIN tb_bidang_urusan bu
        ON LEFT(p.kode_program, LENGTH(bu.kode_bidang_urusan)) = bu.kode_bidang_urusan
    JOIN tb_urusan u
        ON LEFT(bu.kode_bidang_urusan, LENGTH(u.kode_urusan)) = u.kode_urusan
    LEFT JOIN tb_pegawai pg ON pg.nip = sp.pegawai_id
    LEFT JOIN tb_pagu tp_rw
        ON  tp_rw.kode_subkegiatan = sp.kode_subkegiatan
        AND tp_rw.kode_opd         = ?
        AND tp_rw.jenis            = 'renstra'
        AND tp_rw.tahun            = sp.tahun_rekin
    ORDER BY
        u.kode_urusan,
        bu.kode_bidang_urusan,
        p.kode_program,
        k.kode_kegiatan,
        s.kode_subkegiatan
	`
	rows, err := tx.QueryContext(ctx, query, kodeOpd, tahun, kodeOpd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var baseResult []domain.SubKegiatanQuery
	for rows.Next() {
		var row domain.SubKegiatanQuery
		if err := rows.Scan(
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
			&row.PaguRanwal,
			&row.PaguRankhir,
		); err != nil {
			return nil, err
		}
		baseResult = append(baseResult, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// ── Separate indikator query (dual-slot: ranwal + rankhir) ──
	indByKode, err := repository.fetchIndikatorRenjaPeriod(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	applyRankhirInd := func(row *domain.SubKegiatanQuery, ind indRowRenja) {
		row.TargetRanwalId = ind.targetRanwalId
		row.TargetRanwal = ind.targetRanwal
		row.SatuanRanwal = ind.satuanRanwal
		row.TargetRankhirId = ind.targetRankhirId
		row.TargetRankhir = ind.targetRankhir
		row.SatuanRankhir = ind.satuanRankhir
	}
	result := mergeIndikator(baseResult, indByKode, applyRankhirInd)

	// ── Opsi B: Extra synthesis (rankhir) ──
	seenSubkegRk := make(map[string]bool)
	for _, row := range result {
		seenSubkegRk[row.KodeSubKegiatan] = true
	}
	templatesRk, err := repository.fetchAllSubkegTemplates(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	extraPaguRk, err := repository.fetchExtraPagu(ctx, tx, kodeOpd, tahun, "renstra")
	if err != nil {
		return nil, err
	}
	for kodeSub, pagu := range extraPaguRk {
		if seenSubkegRk[kodeSub] {
			continue
		}
		tmpl, ok := templatesRk[kodeSub]
		if !ok {
			continue
		}
		synth := tmpl
		synth.TahunSubKegiatan = tahun
		synth.PaguRanwal = pagu
		synth.PaguRankhir = 0
		synth.StatusSubkegiatan = "subkegiatan tidak ada di tahun ini"
		result = append(result, mergeIndikator([]domain.SubKegiatanQuery{synth}, indByKode, applyRankhirInd)...)
	}

	if len(result) == 0 {
		return []domain.SubKegiatanQuery{}, nil
	}
	return result, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Helper: Opsi B — extra pagu synthesis
// ─────────────────────────────────────────────────────────────────────────────

// fetchAllSubkegTemplates mengembalikan map[kode_subkegiatan → SubKegiatanQuery]
// berisi info hierarki untuk SEMUA subkegiatan yang pernah dipilih OPD ini (semua tahun).
// preferredTahun diprioritaskan di ORDER BY sehingga data pegawai dari tahun itu diutamakan.
func (repository *MatrixRenjaRepositoryImpl) fetchAllSubkegTemplates(
	ctx context.Context, tx *sql.Tx, kodeOpd, preferredTahun string,
) (map[string]domain.SubKegiatanQuery, error) {
	rows, err := tx.QueryContext(ctx, `
        SELECT
            u.kode_urusan, u.nama_urusan,
            bu.kode_bidang_urusan, bu.nama_bidang_urusan,
            p.kode_program, p.nama_program,
            k.kode_kegiatan, k.nama_kegiatan,
            s.kode_subkegiatan, s.nama_subkegiatan,
            rk.pegawai_id,
            COALESCE(pg.nama, '') AS nama_pegawai
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
        LEFT JOIN tb_pegawai pg ON pg.nip = rk.pegawai_id
        WHERE rk.kode_opd = ?
        ORDER BY (rk.tahun = ?) DESC, rk.tahun DESC
    `, kodeOpd, preferredTahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]domain.SubKegiatanQuery)
	for rows.Next() {
		var row domain.SubKegiatanQuery
		if err := rows.Scan(
			&row.KodeUrusan, &row.NamaUrusan,
			&row.KodeBidangUrusan, &row.NamaBidangUrusan,
			&row.KodeProgram, &row.NamaProgram,
			&row.KodeKegiatan, &row.NamaKegiatan,
			&row.KodeSubKegiatan, &row.NamaSubKegiatan,
			&row.PegawaiId, &row.NamaPegawai,
		); err != nil {
			return nil, err
		}
		// Hanya ambil baris pertama per kode_subkegiatan
		// (ORDER BY sudah mengutamakan preferredTahun)
		if _, exists := result[row.KodeSubKegiatan]; !exists {
			result[row.KodeSubKegiatan] = row
		}
	}
	return result, rows.Err()
}

// indRowRenja menyimpan satu indikator (renstra atau renja-spesifik) beserta
// semua slot target-nya (ranwal, rankhir, penetapan, renstra-fallback).
type indRowRenja struct {
	id, kode, kodeOpd, indikator, tahun string
	// Target ranwal (dari tb_target jenis='ranwal')
	targetRanwalId, targetRanwal, satuanRanwal string
	// Target rankhir (dari tb_target jenis='rankhir')
	targetRankhirId, targetRankhir, satuanRankhir string
	// Target penetapan (dari tb_target jenis='penetapan')
	targetPenetapanId, targetPenetapan, satuanPenetapan string
	// Fallback renstra (dari tb_target jenis='renstra' atau jenis='')
	targetRenstra, satuanRenstra string
}

// fetchIndikatorRenjaPeriod mengembalikan map[im.kode → []indRowRenja] berisi
// semua indikator yang relevan untuk tampilan renja (kode_opd, tahun tertentu).
//
// Mencakup dua sumber indikator:
//   - jenis='renstra', im.tahun BETWEEN tahunAwal AND tahunAkhir  (renstra reference)
//   - jenis IN ('ranwal','rankhir','penetapan'), im.tahun = tahun  (renja-spesifik)
//
// Target untuk masing-masing jenis diambil dari tb_target dengan filter tahun renja.
func (repository *MatrixRenjaRepositoryImpl) fetchIndikatorRenjaPeriod(
	ctx context.Context, tx *sql.Tx, kodeOpd, tahun string,
) (map[string][]indRowRenja, error) {
	tahunAwal, tahunAkhir := repository.getRenstraPeriod(ctx, tx, tahun)

	query := `
        SELECT
            im.kode_indikator,
            im.kode,
            im.kode_opd,
            im.indikator,
            COALESCE(im.tahun, ''),
            COALESCE(t_rw.id,''),     COALESCE(t_rw.target,''),     COALESCE(t_rw.satuan,''),
            COALESCE(t_rk.id,''),     COALESCE(t_rk.target,''),     COALESCE(t_rk.satuan,''),
            COALESCE(t_pt.id,''),     COALESCE(t_pt.target,''),     COALESCE(t_pt.satuan,''),
            COALESCE(t_rs.target,''), COALESCE(t_rs.satuan,'')
        FROM tb_indikator_matrix im
        LEFT JOIN tb_target t_rw
            ON t_rw.indikator_id = im.kode_indikator AND t_rw.jenis = 'ranwal'    AND t_rw.tahun = ?
        LEFT JOIN tb_target t_rk
            ON t_rk.indikator_id = im.kode_indikator AND t_rk.jenis = 'rankhir'   AND t_rk.tahun = ?
        LEFT JOIN tb_target t_pt
            ON t_pt.indikator_id = im.kode_indikator AND t_pt.jenis = 'penetapan' AND t_pt.tahun = ?
        LEFT JOIN tb_target t_rs
            ON t_rs.indikator_id = im.kode_indikator
           AND t_rs.tahun = ?
        WHERE im.kode_opd = ?
          AND (
              (im.jenis = 'renstra' AND (
                  EXISTS (
                      SELECT 1 FROM tb_target tx
                      WHERE tx.indikator_id = im.kode_indikator AND tx.tahun = ?
                  )
                  OR im.tahun BETWEEN ? AND ?
              ))
           OR (im.jenis IN ('ranwal','rankhir','penetapan') AND im.tahun = ?)
          )
        ORDER BY im.kode, im.kode_indikator
    `
	rows, err := tx.QueryContext(ctx, query,
		tahun, tahun, tahun, tahun, // t_rw, t_rk, t_pt, t_rs tahun filter
		kodeOpd,               // im.kode_opd
		tahun,                 // EXISTS: punya target tahun renja
		tahunAwal, tahunAkhir, // renstra im.tahun BETWEEN
		tahun, // renja-spesifik = tahun
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]indRowRenja)
	for rows.Next() {
		var r indRowRenja
		if err := rows.Scan(
			&r.id, &r.kode, &r.kodeOpd, &r.indikator, &r.tahun,
			&r.targetRanwalId, &r.targetRanwal, &r.satuanRanwal,
			&r.targetRankhirId, &r.targetRankhir, &r.satuanRankhir,
			&r.targetPenetapanId, &r.targetPenetapan, &r.satuanPenetapan,
			&r.targetRenstra, &r.satuanRenstra,
		); err != nil {
			return nil, err
		}
		// Ranwal: boleh fallback ke target renstra tahun ini (baseline ranwal).
		// Rankhir & penetapan: hanya dari tb_target jenis masing-masing — tanpa fallback renstra.
		if r.targetRanwal == "" && r.targetRenstra != "" {
			r.targetRanwal = r.targetRenstra
			r.satuanRanwal = r.satuanRenstra
		}
		result[r.kode] = append(result[r.kode], r)
	}
	return result, rows.Err()
}

// mergeIndikator menggabungkan baris hierarki dengan indikator dari indByKode.
// Untuk setiap baris hierarki, menghasilkan satu baris per indikator yang cocok
// (cocok = im.kode sama dengan salah satu dari 5 level kode hierarki).
// Jika tidak ada indikator yang cocok, baris hierarki tetap dimasukkan tanpa indikator.
func mergeIndikator(
	rows []domain.SubKegiatanQuery,
	indByKode map[string][]indRowRenja,
	applyFn func(*domain.SubKegiatanQuery, indRowRenja),
) []domain.SubKegiatanQuery {
	if len(indByKode) == 0 {
		return rows
	}
	result := make([]domain.SubKegiatanQuery, 0, len(rows))
	for _, h := range rows {
		kodeLevels := [5]string{h.KodeUrusan, h.KodeBidangUrusan, h.KodeProgram, h.KodeKegiatan, h.KodeSubKegiatan}
		seenInd := make(map[string]struct{})
		var matchedInds []indRowRenja
		for _, kode := range kodeLevels {
			if kode == "" {
				continue
			}
			for _, ind := range indByKode[kode] {
				if _, seen := seenInd[ind.id]; seen {
					continue
				}
				seenInd[ind.id] = struct{}{}
				matchedInds = append(matchedInds, ind)
			}
		}
		if len(matchedInds) == 0 {
			result = append(result, h)
		} else {
			for _, ind := range matchedInds {
				row := h
				row.IndikatorId = ind.id
				row.IndikatorKode = ind.kode
				row.Indikator = ind.indikator
				row.IndikatorTahun = ind.tahun
				row.IndikatorKodeOpd = ind.kodeOpd
				applyFn(&row, ind)
				result = append(result, row)
			}
		}
	}
	return result
}

// getRenstraPeriod mengembalikan tahun_awal dan tahun_akhir dari periode renstra
// yang mencakup tahun renja. Fallback ke tahun itu sendiri jika periode tidak ditemukan.
func (repository *MatrixRenjaRepositoryImpl) getRenstraPeriod(
	ctx context.Context, tx *sql.Tx, tahun string,
) (tahunAwal, tahunAkhir string) {
	err := tx.QueryRowContext(ctx, `
        SELECT p.tahun_awal, p.tahun_akhir
        FROM tb_periode p
        JOIN tb_tahun_periode tp ON p.id = tp.id_periode
        WHERE tp.tahun = ?
        LIMIT 1`, tahun).Scan(&tahunAwal, &tahunAkhir)
	if err != nil {
		return tahun, tahun // fallback: pakai tahun renja sebagai keduanya
	}
	return tahunAwal, tahunAkhir
}

// fetchExtraPagu mengembalikan map[kode_subkegiatan → pagu] dari tb_pagu
// untuk jenis dan tahun tertentu.
func (repository *MatrixRenjaRepositoryImpl) fetchExtraPagu(
	ctx context.Context, tx *sql.Tx, kodeOpd, tahun, jenis string,
) (map[string]int64, error) {
	rows, err := tx.QueryContext(ctx, `
        SELECT kode_subkegiatan, pagu
        FROM tb_pagu
        WHERE kode_opd = ? AND jenis = ? AND tahun = ?
    `, kodeOpd, jenis, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var kode string
		var pagu int64
		if err := rows.Scan(&kode, &pagu); err != nil {
			return nil, err
		}
		result[kode] = pagu
	}
	return result, rows.Err()
}

func (r *MatrixRenjaRepositoryImpl) UpsertBatchIndikatorRenja(
	ctx context.Context,
	tx *sql.Tx,
	items []domain.Indikator,
) error {
	for _, indikator := range items {
		var exists int
		if err := tx.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM tb_indikator_matrix WHERE kode_indikator = ?",
			indikator.KodeIndikator,
		).Scan(&exists); err != nil {
			return err
		}
		if exists > 0 {
			// ── UPDATE indikator ────────────────────────────────────────
			if _, err := tx.ExecContext(ctx, `
                UPDATE tb_indikator_matrix
                SET indikator            = ?,
                    rumus_perhitungan    = ?,
                    sumber_data          = ?,
                    definisi_operasional = ?,
                    tahun                = ?
                WHERE kode_indikator = ? AND jenis = ?`,
				indikator.Indikator,
				indikator.RumusPerhitungan.String,
				indikator.SumberData.String,
				indikator.DefinisiOperasional.String,
				indikator.Tahun, // ← update tahun juga
				indikator.KodeIndikator,
				indikator.Jenis,
			); err != nil {
				return err
			}
			// Upsert target: INSERT ... ON DUPLICATE KEY UPDATE
			// Aman untuk kasus: target sudah ada (update) atau belum ada (insert)
			if len(indikator.Target) > 0 {
				t := indikator.Target[0]
				if _, err := tx.ExecContext(ctx, `
                    INSERT INTO tb_target (id, indikator_id, tahun, target, satuan, jenis)
                    VALUES (?, ?, ?, ?, ?, ?)
                    ON DUPLICATE KEY UPDATE
                        target = VALUES(target),
                        satuan = VALUES(satuan),
                        tahun  = VALUES(tahun)`,
					t.Id,
					indikator.KodeIndikator,
					t.Tahun,
					t.Target,
					t.Satuan,
					indikator.Jenis,
				); err != nil {
					return err
				}
			}
		} else {
			// ── INSERT indikator baru ──────────────────────────────────
			if _, err := tx.ExecContext(ctx, `
                INSERT INTO tb_indikator_matrix
                    (kode_indikator, kode, kode_opd, tahun, indikator,
                     rumus_perhitungan, sumber_data, definisi_operasional, jenis)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				indikator.KodeIndikator,
				indikator.Kode,
				indikator.KodeOpd,
				indikator.Tahun, // ← bug fix: tahun ikut di-insert
				indikator.Indikator,
				indikator.RumusPerhitungan.String,
				indikator.SumberData.String,
				indikator.DefinisiOperasional.String,
				indikator.Jenis,
			); err != nil {
				return err
			}
			// INSERT target
			if len(indikator.Target) > 0 {
				t := indikator.Target[0]
				if _, err := tx.ExecContext(ctx,
					"INSERT INTO tb_target (id, indikator_id, tahun, target, satuan, jenis) VALUES (?, ?, ?, ?, ?, ?)",
					t.Id,
					indikator.KodeIndikator,
					t.Tahun,
					t.Target,
					t.Satuan,
					indikator.Jenis,
				); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *MatrixRenjaRepositoryImpl) CountIndikatorMatrixByScope(ctx context.Context, tx *sql.Tx, kode, kodeOpd, tahun, jenis, prefix string) (int, error) {
	var count int
	err := tx.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM tb_indikator_matrix
		WHERE kode = ? AND kode_opd = ? AND tahun = ? AND jenis = ?
		  AND kode_indikator LIKE ?`,
		kode, kodeOpd, tahun, jenis, prefix+"%",
	).Scan(&count)
	return count, err
}

func (r *MatrixRenjaRepositoryImpl) FindIndikatorRenjaByKode(
	ctx context.Context, tx *sql.Tx, kodeIndikator string,
) (domain.Indikator, error) {
	var ind domain.Indikator
	err := tx.QueryRowContext(ctx, `
        SELECT kode_indikator,
               COALESCE(kode, ''),
               COALESCE(kode_opd, ''),
               COALESCE(indikator, ''),
               COALESCE(rumus_perhitungan, ''),
               COALESCE(sumber_data, ''),
               COALESCE(definisi_operasional, ''),
               COALESCE(jenis, '')
        FROM tb_indikator_matrix
        WHERE kode_indikator = ?`, kodeIndikator).
		Scan(
			&ind.KodeIndikator,
			&ind.Kode,
			&ind.KodeOpd,
			&ind.Indikator,
			&ind.RumusPerhitungan.String,
			&ind.SumberData.String,
			&ind.DefinisiOperasional.String,
			&ind.Jenis,
		)
	if err != nil {
		return domain.Indikator{}, err
	}
	ind.RumusPerhitungan.Valid = true
	ind.SumberData.Valid = true
	ind.DefinisiOperasional.Valid = true
	// Ambil target yang sudah ada
	rows, err := tx.QueryContext(ctx,
		`SELECT id, indikator_id, COALESCE(tahun,''), target, satuan, COALESCE(jenis,'')
         FROM tb_target WHERE indikator_id = ?`, kodeIndikator)
	if err != nil {
		return ind, err
	}
	defer rows.Close()
	for rows.Next() {
		var t domain.Target
		if err := rows.Scan(&t.Id, &t.IndikatorId, &t.Tahun, &t.Target, &t.Satuan, &t.Jenis); err != nil {
			return ind, err
		}
		ind.Target = append(ind.Target, t)
	}
	return ind, nil
}

func (repository *MatrixRenjaRepositoryImpl) UpsertAnggaran(ctx context.Context, tx *sql.Tx, kodeSubkegiatan, kodeOpd, tahun string, pagu int64) error {
	query := `
        INSERT INTO tb_pagu (kode_subkegiatan, kode_opd, tahun, jenis, pagu)
        VALUES (?, ?, ?, 'penetapan', ?)
        ON DUPLICATE KEY UPDATE pagu = VALUES(pagu)
    `
	_, err := tx.ExecContext(ctx, query, kodeSubkegiatan, kodeOpd, tahun, pagu)
	return err
}

func (r *MatrixRenjaRepositoryImpl) DeleteIndicatorsExcept(
	ctx context.Context, tx *sql.Tx,
	kode, kodeOpd, tahun, jenis string, keepList []string,
) error {
	if len(keepList) == 0 {
		// Hapus semua target dalam scope
		delTarget := `
            DELETE FROM tb_target
            WHERE indikator_id IN (
                SELECT kode_indikator FROM tb_indikator_matrix
                WHERE kode = ? AND kode_opd = ? AND tahun = ? AND jenis = ?
            )
        `
		if _, err := tx.ExecContext(ctx, delTarget, kode, kodeOpd, tahun, jenis); err != nil {
			return err
		}
		// Hapus semua indikator dalam scope
		delInd := `DELETE FROM tb_indikator_matrix WHERE kode = ? AND kode_opd = ? AND tahun = ? AND jenis = ?`
		_, err := tx.ExecContext(ctx, delInd, kode, kodeOpd, tahun, jenis)
		return err
	}
	placeholders := make([]string, len(keepList))
	args := make([]interface{}, 0, 4+len(keepList))
	args = append(args, kode, kodeOpd, tahun, jenis)
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
            WHERE kode = ? AND kode_opd = ? AND tahun = ? AND jenis = ?
              AND kode_indikator NOT IN (%s)
        )
    `, inClause)
	if _, err := tx.ExecContext(ctx, delTarget, args...); err != nil {
		return err
	}
	// Hapus indikator yang TIDAK ada di keepList
	delInd := fmt.Sprintf(`
        DELETE FROM tb_indikator_matrix
        WHERE kode = ? AND kode_opd = ? AND tahun = ? AND jenis = ?
          AND kode_indikator NOT IN (%s)
    `, inClause)
	_, err := tx.ExecContext(ctx, delInd, args...)
	return err
}

// GetRenjaPenetapan: indikator dari renstra, slot rankhir + penetapan, pagu_rankhir + pagu_penetapan flat.
// JenisIndikator di setiap baris: "rankhir" (slot rankhir) atau "penetapan" (slot penetapan).
// Hanya menampilkan subkegiatan yang ada di rekin tahun ini (tidak ada extra synthesis).
func (repository *MatrixRenjaRepositoryImpl) GetRenjaPenetapan(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string, _ string) ([]domain.SubKegiatanQuery, error) {
	// Cek apakah ada subkegiatan di rekin tahun ini
	var count int
	if err := tx.QueryRowContext(ctx, `
        SELECT COUNT(*)
        FROM tb_subkegiatan_terpilih st
        JOIN tb_rencana_kinerja rk ON st.rekin_id = rk.id
        WHERE rk.kode_opd = ? AND rk.tahun = ?`,
		kodeOpd, tahun,
	).Scan(&count); err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("subkegiatan belum dipilih pada tahun %s", tahun)
	}

	// ── Main query: hierarki + dual pagu (tanpa inline indikator JOIN) ──
	query := `
    WITH subkeg_pagu AS (
        SELECT
            sub.kode_subkegiatan,
            rk.pegawai_id,
            rk.tahun AS tahun_rekin,
            COALESCE(SUM(rb.anggaran), 0) AS pagu_rankhir
        FROM tb_subkegiatan_terpilih sub
        JOIN tb_rencana_kinerja rk ON sub.rekin_id = rk.id
        LEFT JOIN tb_rencana_aksi ra ON ra.rencana_kinerja_id = rk.id
        LEFT JOIN tb_rincian_belanja rb ON rb.renaksi_id = ra.id
        WHERE rk.kode_opd = ?
          AND rk.tahun = ?
        GROUP BY sub.kode_subkegiatan, rk.pegawai_id, rk.tahun
    )
    SELECT
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
        sp.tahun_rekin      AS tahun_subkegiatan,
        sp.pegawai_id,
        COALESCE(pg.nama,'') AS nama_pegawai,
        sp.pagu_rankhir,
        COALESCE(tp_pt.pagu,0) AS pagu_penetapan
    FROM subkeg_pagu sp
    JOIN tb_subkegiatan s ON sp.kode_subkegiatan = s.kode_subkegiatan
    JOIN tb_master_kegiatan k
        ON LEFT(s.kode_subkegiatan, LENGTH(k.kode_kegiatan)) = k.kode_kegiatan
    JOIN tb_master_program p
        ON LEFT(k.kode_kegiatan, LENGTH(p.kode_program)) = p.kode_program
    JOIN tb_bidang_urusan bu
        ON LEFT(p.kode_program, LENGTH(bu.kode_bidang_urusan)) = bu.kode_bidang_urusan
    JOIN tb_urusan u
        ON LEFT(bu.kode_bidang_urusan, LENGTH(u.kode_urusan)) = u.kode_urusan
    LEFT JOIN tb_pegawai pg ON pg.nip = sp.pegawai_id
    LEFT JOIN tb_pagu tp_pt
        ON  tp_pt.kode_subkegiatan = sp.kode_subkegiatan
        AND tp_pt.kode_opd         = ?
        AND tp_pt.jenis            = 'penetapan'
        AND tp_pt.tahun            = sp.tahun_rekin
    ORDER BY
        u.kode_urusan,
        bu.kode_bidang_urusan,
        p.kode_program,
        k.kode_kegiatan,
        s.kode_subkegiatan
	`
	rows, err := tx.QueryContext(ctx, query, kodeOpd, tahun, kodeOpd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hierResult []domain.SubKegiatanQuery
	for rows.Next() {
		var row domain.SubKegiatanQuery
		if err := rows.Scan(
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
			&row.PaguRankhir,
			&row.PaguPenetapan,
		); err != nil {
			return nil, err
		}
		hierResult = append(hierResult, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// ── Separate indikator query (dual-slot: rankhir + penetapan) ──
	indByKode, err := repository.fetchIndikatorRenjaPeriod(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	baseResult := mergeIndikator(hierResult, indByKode, func(row *domain.SubKegiatanQuery, ind indRowRenja) {
		row.TargetRankhirId = ind.targetRankhirId
		row.TargetRankhir = ind.targetRankhir
		row.SatuanRankhir = ind.satuanRankhir
		row.TargetPenetapanId = ind.targetPenetapanId
		row.TargetPenetapan = ind.targetPenetapan
		row.SatuanPenetapan = ind.satuanPenetapan
	})

	// Penetapan hanya menampilkan subkegiatan yg ada di rekin tahun ini.
	// Tidak ada extra synthesis — subkegiatan yg tidak ada di rekin tidak ditampilkan.

	// Emit 2 baris per indikator: satu ber-JenisIndikator="rankhir", satu="penetapan"
	result := make([]domain.SubKegiatanQuery, 0, len(baseResult)*2)
	seenSubkeg := make(map[string]bool)
	for _, row := range baseResult {
		if row.IndikatorId != "" {
			// Baris rankhir
			rk := row
			rk.JenisIndikator = "rankhir"
			rk.Target = row.TargetRankhir
			rk.Satuan = row.SatuanRankhir
			rk.TargetId = row.TargetRankhirId
			result = append(result, rk)
			// Baris penetapan
			pt := row
			pt.JenisIndikator = "penetapan"
			pt.Target = row.TargetPenetapan
			pt.Satuan = row.SatuanPenetapan
			pt.TargetId = row.TargetPenetapanId
			result = append(result, pt)
		} else {
			// Baris tanpa indikator (subkegiatan masih ditampilkan)
			if !seenSubkeg[row.KodeSubKegiatan] {
				result = append(result, row)
				seenSubkeg[row.KodeSubKegiatan] = true
			}
		}
	}
	if len(result) == 0 {
		return []domain.SubKegiatanQuery{}, nil
	}
	return result, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// CRUD: target per layer (rankhir / penetapan)
// ─────────────────────────────────────────────────────────────────────────────

// UpsertTargetRenjaLayer — upsert tb_target by (indikator_id, tahun, jenis).
// Kolom tb_target.indikator_id diisi = kode_indikator dari request (sama dengan im.kode_indikator).
func (r *MatrixRenjaRepositoryImpl) UpsertTargetRenjaLayer(
	ctx context.Context, tx *sql.Tx,
	kodeIndikator, tahun, target, satuan, jenis string,
) error {
	var exists int
	if err := tx.QueryRowContext(ctx, `
        SELECT COUNT(*) FROM tb_target
        WHERE indikator_id = ? AND tahun = ? AND jenis = ?`,
		kodeIndikator, tahun, jenis,
	).Scan(&exists); err != nil {
		return err
	}
	if exists > 0 {
		_, err := tx.ExecContext(ctx, `
            UPDATE tb_target
            SET target = ?, satuan = ?
            WHERE indikator_id = ? AND tahun = ? AND jenis = ?`,
			target, satuan, kodeIndikator, tahun, jenis,
		)
		return err
	}
	_, err := tx.ExecContext(ctx, `
        INSERT INTO tb_target (id, indikator_id, tahun, target, satuan, jenis)
        VALUES (UUID(), ?, ?, ?, ?, ?)`,
		kodeIndikator, tahun, target, satuan, jenis,
	)
	return err
}

// DeleteTargetRenjaLayer — hapus satu baris target berdasarkan indikator + tahun + jenis.
func (r *MatrixRenjaRepositoryImpl) DeleteTargetRenjaLayer(
	ctx context.Context, tx *sql.Tx,
	kodeIndikator, tahun, jenis string,
) error {
	_, err := tx.ExecContext(ctx,
		`DELETE FROM tb_target WHERE indikator_id = ? AND tahun = ? AND jenis = ?`,
		kodeIndikator, tahun, jenis,
	)
	return err
}

// FindTargetRenjaLayer — ambil satu baris target berdasarkan indikator + tahun + jenis.
func (r *MatrixRenjaRepositoryImpl) FindTargetRenjaLayer(
	ctx context.Context, tx *sql.Tx,
	kodeIndikator, tahun, jenis string,
) (domain.Target, error) {
	var t domain.Target
	err := tx.QueryRowContext(ctx,
		`SELECT id, indikator_id, COALESCE(tahun,''), target, satuan, COALESCE(jenis,'')
         FROM tb_target WHERE indikator_id = ? AND tahun = ? AND jenis = ?`,
		kodeIndikator, tahun, jenis,
	).Scan(&t.Id, &t.IndikatorId, &t.Tahun, &t.Target, &t.Satuan, &t.Jenis)
	if err != nil {
		return domain.Target{}, err
	}
	return t, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// CRUD: pagu per jenis (ranwal / rankhir / penetapan)
// ─────────────────────────────────────────────────────────────────────────────

// UpsertPaguRenjaByJenis — INSERT atau UPDATE satu baris pagu untuk jenis tertentu.
func (r *MatrixRenjaRepositoryImpl) UpsertPaguRenjaByJenis(
	ctx context.Context, tx *sql.Tx,
	kodeSubkegiatan, kodeOpd, tahun, jenis string,
	pagu int64,
) error {
	_, err := tx.ExecContext(ctx, `
        INSERT INTO tb_pagu (kode_subkegiatan, kode_opd, tahun, jenis, pagu)
        VALUES (?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE pagu = VALUES(pagu)`,
		kodeSubkegiatan, kodeOpd, tahun, jenis, pagu,
	)
	return err
}
