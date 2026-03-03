package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
)

type MatrixRenstraRepositoryImpl struct{}

func NewMatrixRenstraRepositoryImpl() *MatrixRenstraRepositoryImpl {
	return &MatrixRenstraRepositoryImpl{}
}

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

	// Perubahan utama vs versi lama:
	// 1. Tambah LEFT JOIN tb_pegawai → nama_pegawai (hilangkan N+1 di service)
	// 2. Tambah LEFT JOIN tb_pagu    → pagu_subkegiatan (sumber pagu baru, jenis='renstra')
	// 3. Hapus i.pagu_anggaran dari SELECT (pagu tidak lagi dari tb_indikator)
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
        COALESCE(pg.nama, '')    AS nama_pegawai,
        COALESCE(tp.pagu, 0)     AS pagu_subkegiatan,
        i.id                     AS indikator_id,
        i.kode                   AS indikator_kode,
        i.indikator,
        i.tahun                  AS indikator_tahun,
        i.kode_opd               AS indikator_kode_opd,
        t.id                     AS target_id,
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
            i.kode = h.kode_urusan         OR
            i.kode = h.kode_bidang_urusan  OR
            i.kode = h.kode_program        OR
            i.kode = h.kode_kegiatan       OR
            i.kode = h.kode_subkegiatan
        )
        AND i.kode_opd = ?
        AND i.tahun BETWEEN ? AND ?
    LEFT JOIN tb_target t ON t.indikator_id = i.id
    	AND t.jenis = 'renstra'
    ORDER BY
        h.kode_urusan,
        h.kode_bidang_urusan,
        h.kode_program,
        h.kode_kegiatan,
        h.kode_subkegiatan,
        i.tahun
    `

	rows, err := tx.QueryContext(ctx, query,
		kodeOpd, tahunAwal, tahunAkhir, // params CTE hierarchy
		kodeOpd,                        // params JOIN tb_pagu (kode_opd)
		kodeOpd, tahunAwal, tahunAkhir, // params JOIN tb_indikator
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
			&data.NamaPegawai,     // dari JOIN tb_pegawai
			&data.PaguSubKegiatan, // dari JOIN tb_pagu jenis='renstra'
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
		if target.Valid {
			data.TargetId = targetId.String
			data.Target = target.String
			data.Satuan = satuan.String
		}

		result = append(result, data)
	}

	return result, nil
}

func (repository *MatrixRenstraRepositoryImpl) SaveIndikator(ctx context.Context, tx *sql.Tx, indikator domain.Indikator) error {
	query := `INSERT INTO tb_indikator (id, kode, kode_opd, indikator, tahun, pagu_anggaran) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := tx.ExecContext(ctx, query, indikator.Id, indikator.Kode, indikator.KodeOpd, indikator.Indikator, indikator.Tahun, indikator.PaguAnggaran)
	return err
}

func (repository *MatrixRenstraRepositoryImpl) SaveTarget(ctx context.Context, tx *sql.Tx, target domain.Target) error {
	query := `INSERT INTO tb_target (id, indikator_id, target, satuan, jenis) VALUES (?, ?, ?, ?, 'renstra')`
	_, err := tx.ExecContext(ctx, query, target.Id, target.IndikatorId, target.Target, target.Satuan)
	return err
}

func (repository *MatrixRenstraRepositoryImpl) FindIndikatorById(ctx context.Context, tx *sql.Tx, indikatorId string) (domain.Indikator, error) {
	query := `
        SELECT 
        i.id, i.kode, i.kode_opd, i.indikator, i.tahun, i.pagu_anggaran,
        t.id as target_id, t.target, t.satuan 
    FROM tb_indikator i 
    LEFT JOIN tb_target t ON t.indikator_id = i.id AND t.jenis = 'renstra'
    WHERE i.id = ?
    `
	var indikator domain.Indikator
	var target domain.Target
	// Gunakan NullString dan NullInt64 untuk handle nilai NULL
	var targetId, targetValue, targetSatuan sql.NullString

	err := tx.QueryRowContext(ctx, query, indikatorId).Scan(
		&indikator.Id,
		&indikator.Kode,
		&indikator.KodeOpd,
		&indikator.Indikator,
		&indikator.Tahun,
		&indikator.PaguAnggaran,
		&targetId,
		&targetValue,
		&targetSatuan,
	)
	if err != nil {
		return domain.Indikator{}, err
	}

	// Set target jika ada nilainya
	if targetId.Valid {
		target = domain.Target{
			Id:     targetId.String,
			Target: targetValue.String,
			Satuan: targetSatuan.String,
		}
		indikator.Target = []domain.Target{target}
	} else {
		indikator.Target = []domain.Target{} // Set empty slice jika tidak ada target
	}

	return indikator, nil
}
func (repository *MatrixRenstraRepositoryImpl) UpdateIndikator(ctx context.Context, tx *sql.Tx, indikator domain.Indikator) error {
	query := `UPDATE tb_indikator SET kode = ?, kode_opd = ?, indikator = ?, tahun = ?, pagu_anggaran = ? WHERE id = ?`
	_, err := tx.ExecContext(ctx, query, indikator.Kode, indikator.KodeOpd, indikator.Indikator, indikator.Tahun, indikator.PaguAnggaran, indikator.Id)
	return err
}

func (repository *MatrixRenstraRepositoryImpl) UpdateTarget(ctx context.Context, tx *sql.Tx, target domain.Target) error {
	query := `UPDATE tb_target SET target = ?, satuan = ? WHERE id = ? AND jenis = 'renstra'`
	_, err := tx.ExecContext(ctx, query, target.Target, target.Satuan, target.Id)
	return err
}

func (repository *MatrixRenstraRepositoryImpl) DeleteIndikator(ctx context.Context, tx *sql.Tx, indikatorId string) error {
	query := `DELETE FROM tb_indikator WHERE id = ?`
	_, err := tx.ExecContext(ctx, query, indikatorId)
	return err
}

func (repository *MatrixRenstraRepositoryImpl) DeleteTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) error {
	query := `DELETE FROM tb_target WHERE indikator_id = ? AND jenis = 'renstra'`
	_, err := tx.ExecContext(ctx, query, indikatorId)
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
