package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
)

type MatrixRenjaRepositoryImpl struct {
}

func NewMatrixRenjaRepositoryImpl() *MatrixRenjaRepositoryImpl {
	return &MatrixRenjaRepositoryImpl{}
}

func (repository *MatrixRenjaRepositoryImpl) GetRenjaRanwal(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.SubKegiatanQuery, error) {
	checkQuery := `
    SELECT COUNT(*) 
    FROM tb_subkegiatan_terpilih st
    JOIN tb_rencana_kinerja rk ON st.rekin_id = rk.id
    WHERE rk.kode_opd = ? 
    AND rk.tahun = ?`

	var count int
	err := tx.QueryRowContext(ctx, checkQuery, kodeOpd, tahun).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("subkegiatan belum dipilih pada tahun %s", tahun)
	}

	query := `
    WITH RECURSIVE hierarchy AS (
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
            rk.tahun    AS tahun_subkegiatan,
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
        AND rk.tahun = ?
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
        COALESCE(tp.pagu, 0) AS pagu_subkegiatan,
        i.id                  AS indikator_id,
        i.kode                AS indikator_kode,
        i.indikator,
        i.tahun               AS indikator_tahun,
        i.kode_opd            AS indikator_kode_opd,
        t.id                  AS target_id,
        t.target,
        t.satuan
    FROM hierarchy h
    LEFT JOIN tb_pegawai pg
        ON pg.nip = h.pegawai_id
    LEFT JOIN tb_pagu tp
        ON  tp.kode_subkegiatan = h.kode_subkegiatan
        AND tp.kode_opd         = ?
        AND tp.jenis            = 'renstra'
        AND tp.tahun            = h.tahun_subkegiatan
    LEFT JOIN tb_indikator i ON (
            i.kode = h.kode_urusan        OR
            i.kode = h.kode_bidang_urusan OR
            i.kode = h.kode_program       OR
            i.kode = h.kode_kegiatan      OR
            i.kode = h.kode_subkegiatan
        )
        AND i.kode_opd = ?
        AND i.tahun    = ?
    LEFT JOIN tb_target t
        ON  t.indikator_id = i.id
        AND t.jenis        = 'ranwal'
    ORDER BY
        h.kode_urusan,
        h.kode_bidang_urusan,
        h.kode_program,
        h.kode_kegiatan,
        h.kode_subkegiatan,
        i.tahun
    `

	rows, err := tx.QueryContext(ctx, query,
		kodeOpd, tahun, // hierarchy
		kodeOpd,        //tb_pagu
		kodeOpd, tahun, // indikator

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

		err := rows.Scan(
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
			&data.NamaPegawai, // dari JOIN tb_pegawai
			&data.PaguSubKegiatan,
			&indikatorId,
			&indikatorKode,
			&indikator,
			&indikatorTahun,
			&indikatorKodeOpd,
			&targetId,
			&target,
			&satuan,
		)
		if err != nil {
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

func (repository *MatrixRenjaRepositoryImpl) GetRenjaRankhir(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.SubKegiatanQuery, error) {

	// Cek apakah ada subkegiatan aktif dari cascading OPD (dengan filter pelaksana)
	checkQuery := `
    SELECT COUNT(DISTINCT sub.kode_subkegiatan)
    FROM tb_rencana_kinerja rk
    JOIN tb_pohon_kinerja pk ON rk.id_pohon = pk.id
    JOIN tb_pelaksana_pokin pp ON pp.pohon_kinerja_id = pk.id
    JOIN tb_pegawai peg ON pp.pegawai_id = peg.id AND peg.nip = rk.pegawai_id
    JOIN tb_subkegiatan_terpilih sub ON sub.rekin_id = rk.id
    WHERE pk.kode_opd = ?
    AND pk.tahun = ?
    AND pk.status NOT IN (
        'menunggu_disetujui', 'tarik pokin opd', 'disetujui',
        'ditolak', 'crosscutting_menunggu', 'crosscutting_ditolak'
    )
    `
	var count int
	err := tx.QueryRowContext(ctx, checkQuery, kodeOpd, tahun).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("tidak ada subkegiatan aktif pada cascading OPD tahun %s", tahun)
	}

	// CTE:
	// 1. valid_rekin   → rencana kinerja yang lolos filter pelaksana cascading OPD
	// 2. subkeg_pagu   → subkegiatan unik beserta total anggaran (SUM jika subkegiatan sama)
	// Lalu hierarchy ke atas: subkegiatan → kegiatan → program → bidang_urusan → urusan
	query := `
    WITH valid_rekin AS (
        SELECT DISTINCT rk.id AS rekin_id
        FROM tb_rencana_kinerja rk
        JOIN tb_pohon_kinerja pk ON rk.id_pohon = pk.id
        JOIN tb_pelaksana_pokin pp ON pp.pohon_kinerja_id = pk.id
        JOIN tb_pegawai peg ON pp.pegawai_id = peg.id AND peg.nip = rk.pegawai_id
        WHERE pk.kode_opd = ?
        AND pk.tahun = ?
        AND pk.status NOT IN (
            'menunggu_disetujui', 'tarik pokin opd', 'disetujui',
            'ditolak', 'crosscutting_menunggu', 'crosscutting_ditolak'
        )
    ),
    subkeg_pagu AS (
        SELECT
            sub.kode_subkegiatan,
            COALESCE(SUM(rb.anggaran), 0) AS total_pagu
        FROM valid_rekin vr
        JOIN tb_subkegiatan_terpilih sub ON sub.rekin_id = vr.rekin_id
        LEFT JOIN tb_rencana_aksi ra ON ra.rencana_kinerja_id = vr.rekin_id
        LEFT JOIN tb_rincian_belanja rb ON rb.renaksi_id = ra.id
        GROUP BY sub.kode_subkegiatan
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
        ? AS tahun_subkegiatan,
        '' AS pegawai_id,
        sp.total_pagu,
        i.id AS indikator_id,
        i.kode AS indikator_kode,
        i.indikator,
        i.tahun AS indikator_tahun,
        i.kode_opd AS indikator_kode_opd,
        i.pagu_anggaran,
        t.id AS target_id,
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
    LEFT JOIN tb_indikator i ON (
            i.kode = u.kode_urusan OR
            i.kode = bu.kode_bidang_urusan OR
            i.kode = p.kode_program OR
            i.kode = k.kode_kegiatan OR
            i.kode = s.kode_subkegiatan
        )
        AND i.kode_opd = ?
        AND i.tahun = ?
 LEFT JOIN tb_target t ON t.indikator_id = i.id
    AND t.jenis = 'rankhir'
    ORDER BY
        u.kode_urusan,
        bu.kode_bidang_urusan,
        p.kode_program,
        k.kode_kegiatan,
        s.kode_subkegiatan,
        i.tahun
    `

	rows, err := tx.QueryContext(ctx, query,
		kodeOpd, tahun, // valid_rekin filter
		tahun,          // tahun_subkegiatan di SELECT
		kodeOpd, tahun, // filter tb_indikator
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
			paguAnggaran                          sql.NullInt64
			totalAnggaran                         int64
		)

		err := rows.Scan(
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
			&totalAnggaran,
			&indikatorId,
			&indikatorKode,
			&indikator,
			&indikatorTahun,
			&indikatorKodeOpd,
			&paguAnggaran,
			&targetId,
			&target,
			&satuan,
		)
		if err != nil {
			return nil, err
		}

		data.TotalAnggaranSubKegiatan = totalAnggaran

		if indikatorId.Valid {
			data.IndikatorId = indikatorId.String
			data.IndikatorKode = indikatorKode.String
			data.Indikator = indikator.String
			data.IndikatorTahun = indikatorTahun.String
			data.IndikatorKodeOpd = indikatorKodeOpd.String
			data.PaguAnggaran = paguAnggaran
		}
		if target.Valid {
			data.TargetId = targetId.String
			data.Target = target.String
			data.Satuan = satuan.String
		}

		result = append(result, data)
	}

	return result, nil
}

func (repository *MatrixRenjaRepositoryImpl) SaveTargetRenja(ctx context.Context, tx *sql.Tx, target domain.Target) error {
	// Tidak ada kolom tahun — tahun mengikuti indikator
	query := `INSERT INTO tb_target (id, indikator_id, target, satuan, jenis) VALUES (?, ?, ?, ?, ?)`
	_, err := tx.ExecContext(ctx, query,
		target.Id, target.IndikatorId, target.Target, target.Satuan, target.Jenis,
	)
	return err
}

func (repository *MatrixRenjaRepositoryImpl) UpdateTargetRenja(ctx context.Context, tx *sql.Tx, target domain.Target) error {
	// Filter by id AND jenis (tanpa tahun)
	query := `UPDATE tb_target SET target = ?, satuan = ? WHERE id = ? AND jenis = ?`
	_, err := tx.ExecContext(ctx, query,
		target.Target, target.Satuan, target.Id, target.Jenis,
	)
	return err
}
