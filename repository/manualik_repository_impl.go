package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"strconv"
)

type ManualIKRepositoryImpl struct {
}

func NewManualIKRepositoryImpl() *ManualIKRepositoryImpl {
	return &ManualIKRepositoryImpl{}
}

func (repository *ManualIKRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, manualik domain.ManualIK) (domain.ManualIK, error) {
	script := `INSERT INTO tb_manual_ik (
        id, indikator_id, perspektif, tujuan_rekin, definisi, key_activities, 
        formula, jenis_indikator, kinerja, penduduk, spasial,
        unit_penanggung_jawab, unit_penyedia_data, sumber_data,
        jangka_waktu_awal, jangka_waktu_akhir, periode_pelaporan
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?)`

	_, err := tx.ExecContext(ctx, script,
		manualik.Id,
		manualik.IndikatorId,
		manualik.Perspektif,
		manualik.TujuanRekin,
		manualik.Definisi,
		manualik.KeyActivities,
		manualik.Formula,
		manualik.JenisIndikator,
		manualik.Kinerja,
		manualik.Penduduk,
		manualik.Spatial,
		manualik.UnitPenanggungJawab,
		manualik.UnitPenyediaData,
		manualik.SumberData,
		manualik.JangkaWaktuAwal,
		manualik.JangkaWaktuAkhir,
		manualik.PeriodePelaporan,
	)
	if err != nil {
		return manualik, err
	}

	return manualik, nil
}

func (repository *ManualIKRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, manualik domain.ManualIK) (domain.ManualIK, error) {
	script := `UPDATE tb_manual_ik SET 
        perspektif = ?, 
        tujuan_rekin = ?,
        definisi = ?,
        key_activities = ?,
        formula = ?,
        jenis_indikator = ?,
        kinerja = ?,
        penduduk = ?,
        spasial = ?,
        unit_penanggung_jawab = ?,
        unit_penyedia_data = ?,
        sumber_data = ?,
        jangka_waktu_awal = ?,
        jangka_waktu_akhir = ?,
        periode_pelaporan = ?
    WHERE indikator_id = ?`

	_, err := tx.ExecContext(ctx, script,
		manualik.Perspektif,
		manualik.TujuanRekin,
		manualik.Definisi,
		manualik.KeyActivities,
		manualik.Formula,
		manualik.JenisIndikator,
		manualik.Kinerja,
		manualik.Penduduk,
		manualik.Spatial,
		manualik.UnitPenanggungJawab,
		manualik.UnitPenyediaData,
		manualik.SumberData,
		manualik.JangkaWaktuAwal,
		manualik.JangkaWaktuAkhir,
		manualik.PeriodePelaporan,
		manualik.IndikatorId,
	)
	if err != nil {
		return manualik, err
	}

	// Ambil data yang baru diupdate menggunakan SELECT
	script = `SELECT id, indikator_id, perspektif, tujuan_rekin, definisi, key_activities, 
        formula, jenis_indikator, kinerja, penduduk, spasial,
        unit_penanggung_jawab, unit_penyedia_data, sumber_data,
        jangka_waktu_awal, jangka_waktu_akhir, periode_pelaporan 
        FROM tb_manual_ik WHERE indikator_id = ?`

	var result domain.ManualIK
	err = tx.QueryRowContext(ctx, script, manualik.IndikatorId).Scan(
		&result.Id,
		&result.IndikatorId,
		&result.Perspektif,
		&result.TujuanRekin,
		&result.Definisi,
		&result.KeyActivities,
		&result.Formula,
		&result.JenisIndikator,
		&result.Kinerja,
		&result.Penduduk,
		&result.Spatial,
		&result.UnitPenanggungJawab,
		&result.UnitPenyediaData,
		&result.SumberData,
		&result.JangkaWaktuAwal,
		&result.JangkaWaktuAkhir,
		&result.PeriodePelaporan,
	)
	if err != nil {
		return manualik, err
	}

	return result, nil
}

func (repository *ManualIKRepositoryImpl) GetManualIK(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.ManualIK, error) {
	script := `SELECT 
        id, indikator_id, perspektif, tujuan_rekin, definisi, key_activities, 
        formula, jenis_indikator, kinerja, penduduk, spasial,
        unit_penanggung_jawab, unit_penyedia_data, sumber_data,
        jangka_waktu_awal, jangka_waktu_akhir, periode_pelaporan 
        FROM tb_manual_ik WHERE indikator_id = ?`

	rows, err := tx.QueryContext(ctx, script, indikatorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var manualIKs []domain.ManualIK
	for rows.Next() {
		var manualIK domain.ManualIK
		err := rows.Scan(
			&manualIK.Id,
			&manualIK.IndikatorId,
			&manualIK.Perspektif,
			&manualIK.TujuanRekin,
			&manualIK.Definisi,
			&manualIK.KeyActivities,
			&manualIK.Formula,
			&manualIK.JenisIndikator,
			&manualIK.Kinerja,
			&manualIK.Penduduk,
			&manualIK.Spatial,
			&manualIK.UnitPenanggungJawab,
			&manualIK.UnitPenyediaData,
			&manualIK.SumberData,
			&manualIK.JangkaWaktuAwal,
			&manualIK.JangkaWaktuAkhir,
			&manualIK.PeriodePelaporan,
		)
		if err != nil {
			return nil, err
		}
		manualIKs = append(manualIKs, manualIK)
	}

	return manualIKs, nil
}

// GetRencanaKinerja mengambil data rencana kinerja dan indikator
func (repository *ManualIKRepositoryImpl) GetRencanaKinerjaWithTarget(ctx context.Context, tx *sql.Tx, indikatorId string) (domain.Indikator, domain.RencanaKinerja, []domain.Target, error) {
	// Query untuk mendapatkan indikator dan rencana kinerja
	scriptIndikator := `
        SELECT 
            i.id, 
            i.rencana_kinerja_id, 
            i.indikator, 
            i.tahun,
            rk.id_pohon,
            rk.nama_rencana_kinerja,
            rk.tahun,
            rk.status_rencana_kinerja,
            rk.catatan,
            rk.kode_opd,
            rk.pegawai_id
        FROM tb_indikator i
        JOIN tb_rencana_kinerja rk ON i.rencana_kinerja_id = rk.id
        WHERE i.id = ?`

	var indikator domain.Indikator
	var rencanaKinerja domain.RencanaKinerja

	err := tx.QueryRowContext(ctx, scriptIndikator, indikatorId).Scan(
		&indikator.Id,
		&indikator.RencanaKinerjaId,
		&indikator.Indikator,
		&indikator.Tahun,
		&rencanaKinerja.IdPohon,
		&rencanaKinerja.NamaRencanaKinerja,
		&rencanaKinerja.Tahun,
		&rencanaKinerja.StatusRencanaKinerja,
		&rencanaKinerja.Catatan,
		&rencanaKinerja.KodeOpd,
		&rencanaKinerja.PegawaiId,
	)
	if err != nil && err != sql.ErrNoRows {
		return domain.Indikator{}, domain.RencanaKinerja{}, nil, err
	}

	// Query untuk mendapatkan target
	scriptTarget := `
        SELECT 
            id, 
            indikator_id, 
            target, 
            satuan,
            tahun
        FROM tb_target 
        WHERE indikator_id = ?`

	rows, err := tx.QueryContext(ctx, scriptTarget, indikatorId)
	if err != nil {
		return domain.Indikator{}, domain.RencanaKinerja{}, nil, err
	}
	defer rows.Close()

	var targets []domain.Target
	for rows.Next() {
		var target domain.Target
		err := rows.Scan(
			&target.Id,
			&target.IndikatorId,
			&target.Target,
			&target.Satuan,
			&target.Tahun,
		)
		if err != nil {
			return domain.Indikator{}, domain.RencanaKinerja{}, nil, err
		}
		targets = append(targets, target)
	}

	return indikator, rencanaKinerja, targets, nil
}

func (repository *ManualIKRepositoryImpl) FindByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) (domain.ManualIK, error) {
	script := `SELECT 
        id, perspektif, tujuan_rekin, definisi, key_activities, 
        formula, jenis_indikator, kinerja, penduduk, spasial,
        unit_penanggung_jawab, unit_penyedia_data, sumber_data,
        jangka_waktu_awal, jangka_waktu_akhir, periode_pelaporan 
        FROM tb_manual_ik WHERE indikator_id = ?`

	var manualIK domain.ManualIK
	err := tx.QueryRowContext(ctx, script, indikatorId).Scan(
		&manualIK.Id,
		&manualIK.Perspektif,
		&manualIK.TujuanRekin,
		&manualIK.Definisi,
		&manualIK.KeyActivities,
		&manualIK.Formula,
		&manualIK.JenisIndikator,
		&manualIK.Kinerja,
		&manualIK.Penduduk,
		&manualIK.Spatial,
		&manualIK.UnitPenanggungJawab,
		&manualIK.UnitPenyediaData,
		&manualIK.SumberData,
		&manualIK.JangkaWaktuAwal,
		&manualIK.JangkaWaktuAkhir,
		&manualIK.PeriodePelaporan,
	)

	// Jika tidak ada data, kembalikan manual IK kosong
	if err == sql.ErrNoRows {
		return manualIK, nil
	}
	if err != nil {
		return manualIK, err
	}

	manualIK.IndikatorId = indikatorId
	return manualIK, nil
}

func (repository *ManualIKRepositoryImpl) FindManualIKSasaranOpdByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string, tahun string) (domain.ManualIK, error) {
	var manualIK domain.ManualIK
	var indikator domain.Indikator

	// Query untuk data indikator dan rencana kinerja terlebih dahulu
	scriptIndikator := `
    SELECT 
        i.id as indikator_id,
        i.indikator,
        rk.nama_rencana_kinerja,
        rk.tahun_awal,
        rk.tahun_akhir,
        rk.jenis_periode
    FROM tb_indikator i
    JOIN tb_rencana_kinerja rk ON i.rencana_kinerja_id = rk.id
    WHERE i.id = ?`

	err := tx.QueryRowContext(ctx, scriptIndikator, indikatorId).Scan(
		&indikator.Id,
		&indikator.Indikator,
		&indikator.RencanaKinerja.NamaRencanaKinerja,
		&indikator.TahunAwal,
		&indikator.TahunAkhir,
		&indikator.JenisPeriode,
	)

	if err == sql.ErrNoRows {
		return manualIK, fmt.Errorf("indikator tidak ditemukan")
	}
	if err != nil {
		return manualIK, err
	}

	// Validasi tahun
	tahunInt, _ := strconv.Atoi(tahun)
	tahunAwalInt, _ := strconv.Atoi(indikator.TahunAwal)
	tahunAkhirInt, _ := strconv.Atoi(indikator.TahunAkhir)

	if tahunInt < tahunAwalInt || tahunInt > tahunAkhirInt {
		return manualIK, fmt.Errorf("tahun %s diluar range (%s-%s)",
			tahun, indikator.TahunAwal, indikator.TahunAkhir)
	}

	// Query untuk manual IK jika ada
	scriptManualIK := `
    SELECT 
        m.id,
        m.perspektif,
        m.tujuan_rekin,
        m.definisi,
        m.key_activities,
        m.formula,
        m.jenis_indikator,
        m.kinerja,
        m.penduduk,
        m.spasial,
        m.unit_penanggung_jawab,
        m.unit_penyedia_data,
        m.sumber_data,
        m.jangka_waktu_awal,
        m.jangka_waktu_akhir,
        m.periode_pelaporan
    FROM tb_manual_ik m
    WHERE m.indikator_id = ?`

	err = tx.QueryRowContext(ctx, scriptManualIK, indikatorId).Scan(
		&manualIK.Id,
		&manualIK.Perspektif,
		&manualIK.TujuanRekin,
		&manualIK.Definisi,
		&manualIK.KeyActivities,
		&manualIK.Formula,
		&manualIK.JenisIndikator,
		&manualIK.Kinerja,
		&manualIK.Penduduk,
		&manualIK.Spatial,
		&manualIK.UnitPenanggungJawab,
		&manualIK.UnitPenyediaData,
		&manualIK.SumberData,
		&manualIK.JangkaWaktuAwal,
		&manualIK.JangkaWaktuAkhir,
		&manualIK.PeriodePelaporan,
	)

	// Jika manual IK tidak ditemukan, tetap lanjutkan dengan nilai default
	if err == sql.ErrNoRows {
		// Set nilai default untuk manual IK
		manualIK = domain.ManualIK{
			Id:                  0,
			IndikatorId:         indikatorId,
			Perspektif:          "",
			TujuanRekin:         "",
			Definisi:            "",
			KeyActivities:       "",
			Formula:             "",
			JenisIndikator:      "",
			Kinerja:             false,
			Penduduk:            false,
			Spatial:             false,
			UnitPenanggungJawab: "",
			UnitPenyediaData:    "",
			SumberData:          "",
			JangkaWaktuAwal:     "",
			JangkaWaktuAkhir:    "",
			PeriodePelaporan:    "",
		}
	} else if err != nil {
		return manualIK, err
	}

	// Query untuk target
	scriptTarget := `
    SELECT 
        target,
        satuan,
        tahun
    FROM tb_target 
    WHERE indikator_id = ? AND tahun = ?`

	var target domain.Target
	err = tx.QueryRowContext(ctx, scriptTarget, indikatorId, tahun).Scan(
		&target.Target,
		&target.Satuan,
		&target.Tahun,
	)

	// Jika target tidak ditemukan, buat target kosong
	if err == sql.ErrNoRows {
		target = domain.Target{
			IndikatorId: indikatorId,
			Target:      "",
			Satuan:      "",
			Tahun:       tahun,
		}
	} else if err != nil {
		return manualIK, err
	}

	// Set data indikator dan target
	indikator.Target = []domain.Target{target}
	manualIK.IndikatorId = indikatorId
	manualIK.DataIndikator = indikator

	return manualIK, nil
}

func (repository *ManualIKRepositoryImpl) DeleteByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) error {
	script := `DELETE FROM tb_manual_ik WHERE indikator_id = ?`
	_, err := tx.ExecContext(ctx, script, indikatorId)
	return err
}

func (repository *ManualIKRepositoryImpl) IsIndikatorExist(ctx context.Context, tx *sql.Tx, indikatorId string) (bool, error) {
	// Query untuk memeriksa keberadaan indikator di tabel manual_ik
	script := `SELECT COUNT(*) FROM tb_manual_ik WHERE indikator_id = ?`

	var count int
	err := tx.QueryRowContext(ctx, script, indikatorId).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("gagal memeriksa keberadaan indikator: %v", err)
	}

	// Mengembalikan true jika count > 0 (indikator ditemukan)
	return count > 0, nil
}
