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

func (repository *MatrixRenjaRepositoryImpl) GetRenja(ctx context.Context, tx *sql.Tx, kodeOpd, tahun, jenisIndikator, jenisPagu string) ([]domain.SubKegiatanQuery, error) {
	// Pastikan ada subkegiatan yang dipilih
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
        COALESCE(pg.nama, '')  AS nama_pegawai,
        COALESCE(tp.pagu, 0)   AS pagu_subkegiatan,
        im.kode_indikator      AS indikator_id,
        im.kode                AS indikator_kode,
        im.indikator,
        COALESCE(im.tahun, '') AS indikator_tahun,
        im.kode_opd            AS indikator_kode_opd,
        t.id                   AS target_id,
        t.target,
        t.satuan
    FROM hierarchy h
    LEFT JOIN tb_pegawai pg
        ON pg.nip = h.pegawai_id
    LEFT JOIN tb_pagu tp
        ON  tp.kode_subkegiatan = h.kode_subkegiatan
        AND tp.kode_opd         = ?
        AND tp.jenis            = ?
        AND tp.tahun            = h.tahun_subkegiatan
    LEFT JOIN tb_indikator_matrix im
        ON (
            im.kode = h.kode_urusan        OR
            im.kode = h.kode_bidang_urusan OR
            im.kode = h.kode_program       OR
            im.kode = h.kode_kegiatan      OR
            im.kode = h.kode_subkegiatan
        )
            AND im.kode_opd = ?
            AND im.jenis    = ?
            AND im.tahun    = ?  
    LEFT JOIN tb_target t
        ON t.indikator_id = im.kode_indikator
    ORDER BY
        h.kode_urusan,
        h.kode_bidang_urusan,
        h.kode_program,
        h.kode_kegiatan,
        h.kode_subkegiatan
    `
	rows, err := tx.QueryContext(ctx, query,
		kodeOpd, tahun, // hierarchy WHERE
		kodeOpd, jenisPagu, // tb_pagu
		kodeOpd, jenisIndikator, tahun, // tb_indikator_matrix + filter tahun
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.SubKegiatanQuery
	for rows.Next() {
		var data domain.SubKegiatanQuery
		var (
			indikatorId, indikatorKode, indikator sql.NullString
			indikatorTahun, indikatorKodeOpd      sql.NullString
			targetId, target, satuan              sql.NullString
		)
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
			&data.PaguSubKegiatan,
			&indikatorId,
			&indikatorKode,
			&indikator,
			&indikatorTahun,
			&indikatorKodeOpd,
			&targetId,
			&target,
			&satuan,
		); err != nil {
			return nil, err
		}
		if indikatorId.Valid {
			data.IndikatorId = indikatorId.String
			data.IndikatorKode = indikatorKode.String
			data.Indikator = indikator.String
			data.IndikatorTahun = indikatorTahun.String
			data.IndikatorKodeOpd = indikatorKodeOpd.String
		}
		if targetId.Valid {
			data.TargetId = targetId.String
			data.Target = target.String
			data.Satuan = satuan.String
		}
		result = append(result, data)
	}
	return result, nil
}

func (repository *MatrixRenjaRepositoryImpl) GetRenjaRankhir(
	ctx context.Context,
	tx *sql.Tx,
	kodeOpd, tahun, jenisIndikator string,
) ([]domain.SubKegiatanQuery, error) {
	// Check sama seperti ranwal
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
	query := `
    WITH subkeg_pagu AS (
        -- Sumber subkegiatan SAMA dengan ranwal (tb_subkegiatan_terpilih)
        -- Pagu dari SUM rincian_belanja (beda dari ranwal yang pakai tb_pagu)
        SELECT
            sub.kode_subkegiatan,
            rk.pegawai_id,
            rk.tahun                          AS tahun_rekin,
            COALESCE(SUM(rb.anggaran), 0)     AS total_pagu
        FROM tb_subkegiatan_terpilih sub
        JOIN tb_rencana_kinerja rk ON sub.rekin_id = rk.id
        LEFT JOIN tb_rencana_aksi ra ON ra.rencana_kinerja_id = rk.id
        LEFT JOIN tb_rincian_belanja rb ON rb.renaksi_id = ra.id
        WHERE rk.kode_opd = ?
          AND rk.tahun    = ?
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
        sp.tahun_rekin                 AS tahun_subkegiatan,
        sp.pegawai_id,
        COALESCE(pg.nama, '')          AS nama_pegawai,
        sp.total_pagu,
        im.kode_indikator              AS indikator_id,
        im.kode                        AS indikator_kode,
        im.indikator,
        COALESCE(im.tahun, '')         AS indikator_tahun,
        im.kode_opd                    AS indikator_kode_opd,
        t.id                           AS target_id,
        t.target,
        t.satuan
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
    LEFT JOIN tb_pegawai pg
        ON pg.nip = sp.pegawai_id
    LEFT JOIN tb_indikator_matrix im
        ON (
            im.kode = u.kode_urusan         OR
            im.kode = bu.kode_bidang_urusan OR
            im.kode = p.kode_program        OR
            im.kode = k.kode_kegiatan       OR
            im.kode = s.kode_subkegiatan
        )
        AND im.kode_opd = ?
        AND im.jenis    = ?
        AND im.tahun    = ?
    LEFT JOIN tb_target t
        ON t.indikator_id = im.kode_indikator
    ORDER BY
        u.kode_urusan,
        bu.kode_bidang_urusan,
        p.kode_program,
        k.kode_kegiatan,
        s.kode_subkegiatan
    `
	rows, err := tx.QueryContext(ctx, query,
		kodeOpd, tahun, // subkeg_pagu WHERE
		kodeOpd, jenisIndikator, tahun, // tb_indikator_matrix filter
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.SubKegiatanQuery
	for rows.Next() {
		var data domain.SubKegiatanQuery
		var (
			indikatorId, indikatorKode, indikator sql.NullString
			indikatorTahun, indikatorKodeOpd      sql.NullString
			targetId, target, satuan              sql.NullString
			totalAnggaran                         int64
		)
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
			&totalAnggaran,
			&indikatorId,
			&indikatorKode,
			&indikator,
			&indikatorTahun,
			&indikatorKodeOpd,
			&targetId,
			&target,
			&satuan,
		); err != nil {
			return nil, err
		}
		data.TotalAnggaranSubKegiatan = totalAnggaran
		if indikatorId.Valid {
			data.IndikatorId = indikatorId.String
			data.IndikatorKode = indikatorKode.String
			data.Indikator = indikator.String
			data.IndikatorTahun = indikatorTahun.String
			data.IndikatorKodeOpd = indikatorKodeOpd.String
		}
		if targetId.Valid {
			data.TargetId = targetId.String
			data.Target = target.String
			data.Satuan = satuan.String
		}
		result = append(result, data)
	}
	return result, nil
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

func (r *MatrixRenjaRepositoryImpl) CountKodeIndikatorByPrefix(
	ctx context.Context, tx *sql.Tx, prefix string,
) (int, error) {
	var count int
	err := tx.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM tb_indikator_matrix WHERE kode_indikator LIKE ?",
		prefix+"%",
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
