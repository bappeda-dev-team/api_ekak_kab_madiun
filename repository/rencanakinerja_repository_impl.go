package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type RencanaKinerjaRepositoryImpl struct {
}

func NewRencanaKinerjaRepositoryImpl() *RencanaKinerjaRepositoryImpl {
	return &RencanaKinerjaRepositoryImpl{}
}

func (repository *RencanaKinerjaRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, rencanaKinerja domain.RencanaKinerja) (domain.RencanaKinerja, error) {
	script := "INSERT INTO tb_rencana_kinerja (id, id_pohon, nama_rencana_kinerja, tahun, status_rencana_kinerja, catatan, kode_opd, pegawai_id, kode_subkegiatan, tahun_awal, tahun_akhir, jenis_periode, periode_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, rencanaKinerja.Id, rencanaKinerja.IdPohon, rencanaKinerja.NamaRencanaKinerja, rencanaKinerja.Tahun, rencanaKinerja.StatusRencanaKinerja, rencanaKinerja.Catatan, rencanaKinerja.KodeOpd, rencanaKinerja.PegawaiId, rencanaKinerja.KodeSubKegiatan, rencanaKinerja.TahunAwal, rencanaKinerja.TahunAkhir, rencanaKinerja.JenisPeriode, rencanaKinerja.PeriodeId)
	if err != nil {
		return domain.RencanaKinerja{}, fmt.Errorf("error saat menyimpan rencana kinerja: %v", err)
	}

	for _, indikator := range rencanaKinerja.Indikator {
		queryIndikator := "INSERT INTO tb_indikator (id, rencana_kinerja_id, indikator, tahun) VALUES (?, ?, ?, ?)"
		_, err := tx.ExecContext(ctx, queryIndikator, indikator.Id, rencanaKinerja.Id, indikator.Indikator, indikator.Tahun)
		if err != nil {
			return domain.RencanaKinerja{}, fmt.Errorf("error saat menyimpan indikator: %v", err)
		}

		for _, target := range indikator.Target {
			queryTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
			_, err := tx.ExecContext(ctx, queryTarget, target.Id, indikator.Id, target.Target, target.Satuan, target.Tahun)
			if err != nil {
				return domain.RencanaKinerja{}, fmt.Errorf("error saat menyimpan target: %v", err)
			}
		}
	}

	return rencanaKinerja, nil
}

func (repository *RencanaKinerjaRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, rencanaKinerja domain.RencanaKinerja) (domain.RencanaKinerja, error) {
	script := "UPDATE tb_rencana_kinerja SET id_pohon = ?, nama_rencana_kinerja = ?, tahun = ?, status_rencana_kinerja = ?, catatan = ?, kode_opd = ?, pegawai_id = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, rencanaKinerja.IdPohon, rencanaKinerja.NamaRencanaKinerja, rencanaKinerja.Tahun, rencanaKinerja.StatusRencanaKinerja, rencanaKinerja.Catatan, rencanaKinerja.KodeOpd, rencanaKinerja.PegawaiId, rencanaKinerja.Id)
	if err != nil {
		return domain.RencanaKinerja{}, err
	}

	scriptDeleteTarget := "DELETE FROM tb_target WHERE indikator_id IN (SELECT id FROM tb_indikator WHERE rencana_kinerja_id = ?)"
	_, err = tx.ExecContext(ctx, scriptDeleteTarget, rencanaKinerja.Id)
	if err != nil {
		return domain.RencanaKinerja{}, err
	}

	queryDeleteIndikator := "DELETE FROM tb_indikator WHERE rencana_kinerja_id = ?"
	_, err = tx.ExecContext(ctx, queryDeleteIndikator, rencanaKinerja.Id)
	if err != nil {
		return domain.RencanaKinerja{}, err
	}
	for _, indikator := range rencanaKinerja.Indikator {
		scriptIndikator := "INSERT INTO tb_indikator (id, rencana_kinerja_id, indikator, tahun) VALUES (?, ?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptIndikator, indikator.Id, rencanaKinerja.Id, indikator.Indikator, indikator.Tahun)
		if err != nil {
			return domain.RencanaKinerja{}, err
		}

		for _, target := range indikator.Target {
			queryTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
			_, err := tx.ExecContext(ctx, queryTarget, target.Id, indikator.Id, target.Target, target.Satuan, target.Tahun)
			if err != nil {
				return domain.RencanaKinerja{}, err
			}
		}
	}

	return rencanaKinerja, nil
}

func (repository *RencanaKinerjaRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, pegawaiId string, kodeOPD string, tahun string) ([]domain.RencanaKinerja, error) {
	script := "SELECT id, id_pohon, nama_rencana_kinerja, tahun, status_rencana_kinerja, catatan, kode_opd, pegawai_id, created_at FROM tb_rencana_kinerja WHERE 1=1"
	params := []interface{}{}

	if pegawaiId != "" {
		script += " AND pegawai_id = ?"
		params = append(params, pegawaiId)
	}
	if kodeOPD != "" {
		script += " AND kode_opd = ?"
		params = append(params, kodeOPD)
	}
	if tahun != "" {
		script += " AND tahun = ?"
		params = append(params, tahun)
	}

	script += " ORDER BY created_at ASC"

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rencanaKinerjas []domain.RencanaKinerja

	for rows.Next() {
		var rencanaKinerja domain.RencanaKinerja
		err := rows.Scan(&rencanaKinerja.Id, &rencanaKinerja.IdPohon, &rencanaKinerja.NamaRencanaKinerja, &rencanaKinerja.Tahun, &rencanaKinerja.StatusRencanaKinerja, &rencanaKinerja.Catatan, &rencanaKinerja.KodeOpd, &rencanaKinerja.PegawaiId, &rencanaKinerja.CreatedAt)
		if err != nil {
			return nil, err
		}
		rencanaKinerjas = append(rencanaKinerjas, rencanaKinerja)
	}

	return rencanaKinerjas, nil
}

func (repository *RencanaKinerjaRepositoryImpl) FindIndikatorbyRekinId(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.Indikator, error) {
	script := "SELECT id, rencana_kinerja_id, indikator, tahun FROM tb_indikator WHERE rencana_kinerja_id = ?"
	params := []interface{}{rekinId}

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indikators []domain.Indikator

	for rows.Next() {
		var indikator domain.Indikator
		err := rows.Scan(&indikator.Id, &indikator.RencanaKinerjaId, &indikator.Indikator, &indikator.Tahun)
		if err != nil {
			return nil, err
		}
		indikators = append(indikators, indikator)
	}

	return indikators, nil
}

func (repository *RencanaKinerjaRepositoryImpl) FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error) {
	script := "SELECT id, indikator_id, target, satuan, tahun FROM tb_target WHERE indikator_id = ?"
	params := []interface{}{indikatorId}

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []domain.Target

	for rows.Next() {
		var target domain.Target
		err := rows.Scan(&target.Id, &target.IndikatorId, &target.Target, &target.Satuan, &target.Tahun)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}

	return targets, nil
}

func (repository *RencanaKinerjaRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string, kodeOPD string, tahun string) (domain.RencanaKinerja, error) {
	script := "SELECT id, id_pohon, nama_rencana_kinerja, tahun, status_rencana_kinerja, catatan, kode_opd, pegawai_id FROM tb_rencana_kinerja WHERE id = ?"
	params := []interface{}{id}

	if kodeOPD != "" {
		script += " AND kode_opd = ?"
		params = append(params, kodeOPD)
	}
	if tahun != "" {
		script += " AND tahun = ?"
		params = append(params, tahun)
	}

	row := tx.QueryRowContext(ctx, script, params...)
	var rencanaKinerja domain.RencanaKinerja
	err := row.Scan(&rencanaKinerja.Id, &rencanaKinerja.IdPohon, &rencanaKinerja.NamaRencanaKinerja, &rencanaKinerja.Tahun, &rencanaKinerja.StatusRencanaKinerja, &rencanaKinerja.Catatan, &rencanaKinerja.KodeOpd, &rencanaKinerja.PegawaiId)
	if err != nil {
		return domain.RencanaKinerja{}, err
	}

	return rencanaKinerja, nil
}

func (repository *RencanaKinerjaRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	script := []string{
		"DELETE FROM tb_manual_ik WHERE indikator_id IN (SELECT id FROM tb_indikator WHERE rencana_kinerja_id = ?)",
		"DELETE FROM tb_target WHERE indikator_id IN (SELECT id FROM tb_indikator WHERE rencana_kinerja_id = ?)",
		"DELETE FROM tb_indikator WHERE rencana_kinerja_id = ?",
		"DELETE FROM tb_rencana_kinerja WHERE id = ?",
	}

	for _, script := range script {
		_, err := tx.ExecContext(ctx, script, id)
		if err != nil {
			return fmt.Errorf("gagal menghapus data: %v", err)
		}
	}

	return nil
}

func (repository *RencanaKinerjaRepositoryImpl) FindAllRincianKak(ctx context.Context, tx *sql.Tx, rencanaKinerjaId string, pegawaiId string) ([]domain.RencanaKinerja, error) {
	log.Printf("Mencari rencana kinerja dengan ID: %s dan PegawaiID: %s", rencanaKinerjaId, pegawaiId)

	script := `
		SELECT 
			id, 
			id_pohon, 
			nama_rencana_kinerja, 
			tahun, 
			status_rencana_kinerja, 
			catatan, 
			kode_opd, 
			pegawai_id, 
			kode_subkegiatan, 
			created_at 
		FROM tb_rencana_kinerja 
		WHERE 1=1
	`
	var params []interface{}

	if rencanaKinerjaId != "" {
		script += " AND id = ?"
		params = append(params, rencanaKinerjaId)
	}

	if pegawaiId != "" {
		script += " AND pegawai_id = ?"
		params = append(params, pegawaiId)
	}

	script += " ORDER BY created_at ASC"

	log.Printf("Executing query: %s with params: %v", script, params)

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, fmt.Errorf("error querying rencana kinerja: %v", err)
	}
	defer rows.Close()

	var rencanaKinerjas []domain.RencanaKinerja

	for rows.Next() {
		var rencanaKinerja domain.RencanaKinerja
		err := rows.Scan(
			&rencanaKinerja.Id,
			&rencanaKinerja.IdPohon,
			&rencanaKinerja.NamaRencanaKinerja,
			&rencanaKinerja.Tahun,
			&rencanaKinerja.StatusRencanaKinerja,
			&rencanaKinerja.Catatan,
			&rencanaKinerja.KodeOpd,
			&rencanaKinerja.PegawaiId,
			&rencanaKinerja.KodeSubKegiatan,
			&rencanaKinerja.CreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("error scanning rencana kinerja: %v", err)
		}
		rencanaKinerjas = append(rencanaKinerjas, rencanaKinerja)
	}

	log.Printf("Found %d rencana kinerja records", len(rencanaKinerjas))
	return rencanaKinerjas, nil
}

func (repository *RencanaKinerjaRepositoryImpl) RekinsasaranOpd(ctx context.Context, tx *sql.Tx, pegawaiId string, kodeOPD string, tahun string) ([]domain.RencanaKinerja, error) {
	script := `
              SELECT DISTINCT 
            rk.id, 
            rk.id_pohon, 
            rk.nama_rencana_kinerja,
            rk.tahun_awal,
            rk.tahun_akhir, 
            rk.status_rencana_kinerja, 
            COALESCE(rk.catatan, ''), 
            rk.kode_opd, 
            rk.pegawai_id,
            rk.created_at
        FROM tb_rencana_kinerja rk
        INNER JOIN tb_pegawai p ON rk.pegawai_id = p.nip
        INNER JOIN tb_pohon_kinerja pk ON rk.id_pohon = pk.id
        INNER JOIN tb_pelaksana_pokin pl ON pk.id = pl.pohon_kinerja_id
        INNER JOIN tb_pegawai pp ON pl.pegawai_id = pp.id
        INNER JOIN tb_indikator i ON rk.id = i.rencana_kinerja_id
        WHERE 1=1
        AND ? BETWEEN rk.tahun_awal AND rk.tahun_akhir
    `
	params := []interface{}{tahun}

	if pegawaiId != "" {
		script += " AND pp.nip = ?"
		params = append(params, pegawaiId)
	}
	if kodeOPD != "" {
		script += " AND rk.kode_opd = ?"
		params = append(params, kodeOPD)
	}

	script += " ORDER BY rk.created_at ASC"

	// Hapus join dan filter dengan tb_target di query utama

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rencanaKinerjas []domain.RencanaKinerja
	seenIds := make(map[string]bool)

	for rows.Next() {
		var rencanaKinerja domain.RencanaKinerja
		err := rows.Scan(
			&rencanaKinerja.Id,
			&rencanaKinerja.IdPohon,
			&rencanaKinerja.NamaRencanaKinerja,
			&rencanaKinerja.TahunAwal,
			&rencanaKinerja.TahunAkhir,
			&rencanaKinerja.StatusRencanaKinerja,
			&rencanaKinerja.Catatan,
			&rencanaKinerja.KodeOpd,
			&rencanaKinerja.PegawaiId,
			&rencanaKinerja.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if !seenIds[rencanaKinerja.Id] {
			seenIds[rencanaKinerja.Id] = true
			rencanaKinerjas = append(rencanaKinerjas, rencanaKinerja)
		}
	}

	return rencanaKinerjas, nil
}

func (repository *RencanaKinerjaRepositoryImpl) FindIndikatorSasaranbyRekinId(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.Indikator, error) {
	script := `
        SELECT 
            id,
            rencana_kinerja_id,
            indikator,
            COALESCE(tahun, ''),
            created_at
        FROM tb_indikator 
        WHERE rencana_kinerja_id = ?
    `

	rows, err := tx.QueryContext(ctx, script, rekinId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indikators []domain.Indikator
	for rows.Next() {
		var indikator domain.Indikator
		err := rows.Scan(
			&indikator.Id,
			&indikator.RencanaKinerjaId,
			&indikator.Indikator,
			&indikator.Tahun,
			&indikator.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		indikators = append(indikators, indikator)
	}

	return indikators, nil
}

func (repository *RencanaKinerjaRepositoryImpl) FindTargetByIndikatorIdAndTahun(ctx context.Context, tx *sql.Tx, indikatorId string, tahun string) ([]domain.Target, error) {
	script := `
        SELECT 
            id,
            indikator_id,
            COALESCE(target, ''),
            COALESCE(satuan, ''),
            COALESCE(tahun, '')
        FROM tb_target 
        WHERE indikator_id = ?
        AND tahun = ?
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
			&target.IndikatorId,
			&target.Target,
			&target.Satuan,
			&target.Tahun,
		)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}

	// Jika tidak ada target untuk tahun tersebut, kembalikan target kosong
	if len(targets) == 0 {
		targets = append(targets, domain.Target{
			Id:          "",
			IndikatorId: indikatorId,
			Target:      "",
			Satuan:      "",
			Tahun:       tahun,
		})
	}

	return targets, nil
}

func (repository *RencanaKinerjaRepositoryImpl) CreateRekinLevel1(ctx context.Context, tx *sql.Tx, rencanaKinerja domain.RencanaKinerja) (domain.RencanaKinerja, error) {
	script := "INSERT INTO tb_rencana_kinerja (id, id_pohon, sasaranopd_id, nama_rencana_kinerja, tahun, status_rencana_kinerja, catatan, kode_opd, pegawai_id, periode_id, tahun_awal, tahun_akhir, jenis_periode, kode_subkegiatan) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, rencanaKinerja.Id, rencanaKinerja.IdPohon, rencanaKinerja.SasaranOpdId, rencanaKinerja.NamaRencanaKinerja, rencanaKinerja.Tahun, rencanaKinerja.StatusRencanaKinerja, rencanaKinerja.Catatan, rencanaKinerja.KodeOpd, rencanaKinerja.PegawaiId, rencanaKinerja.PeriodeId, rencanaKinerja.TahunAwal, rencanaKinerja.TahunAkhir, rencanaKinerja.JenisPeriode, rencanaKinerja.KodeSubKegiatan)
	if err != nil {
		return domain.RencanaKinerja{}, fmt.Errorf("error saat menyimpan rencana kinerja: %v", err)
	}

	for _, indikator := range rencanaKinerja.Indikator {
		queryIndikator := "INSERT INTO tb_indikator (id, rencana_kinerja_id, indikator, tahun) VALUES (?, ?, ?, ?)"
		_, err := tx.ExecContext(ctx, queryIndikator, indikator.Id, rencanaKinerja.Id, indikator.Indikator, indikator.Tahun)
		if err != nil {
			return domain.RencanaKinerja{}, fmt.Errorf("error saat menyimpan indikator: %v", err)
		}

		for _, target := range indikator.Target {
			queryTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
			_, err := tx.ExecContext(ctx, queryTarget, target.Id, indikator.Id, target.Target, target.Satuan, target.Tahun)
			if err != nil {
				return domain.RencanaKinerja{}, fmt.Errorf("error saat menyimpan target: %v", err)
			}
		}
	}

	return rencanaKinerja, nil
}

func (repository *RencanaKinerjaRepositoryImpl) UpdateRekinLevel1(ctx context.Context, tx *sql.Tx, rencanaKinerja domain.RencanaKinerja) (domain.RencanaKinerja, error) {
	script := "UPDATE tb_rencana_kinerja SET id_pohon = ?, sasaranopd_id = ?, nama_rencana_kinerja = ?, tahun = ?, status_rencana_kinerja = ?, catatan = ?, kode_opd = ?, pegawai_id = ?, periode_id = ?, tahun_awal = ?, tahun_akhir = ?, jenis_periode = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, rencanaKinerja.IdPohon, rencanaKinerja.SasaranOpdId, rencanaKinerja.NamaRencanaKinerja, rencanaKinerja.Tahun, rencanaKinerja.StatusRencanaKinerja, rencanaKinerja.Catatan, rencanaKinerja.KodeOpd, rencanaKinerja.PegawaiId, rencanaKinerja.PeriodeId, rencanaKinerja.TahunAwal, rencanaKinerja.TahunAkhir, rencanaKinerja.JenisPeriode, rencanaKinerja.Id)
	if err != nil {
		return domain.RencanaKinerja{}, err
	}

	// Hapus target yang terkait dengan indikator yang akan dihapus
	scriptDeleteTarget := "DELETE FROM tb_target WHERE indikator_id IN (SELECT id FROM tb_indikator WHERE rencana_kinerja_id = ?)"
	_, err = tx.ExecContext(ctx, scriptDeleteTarget, rencanaKinerja.Id)
	if err != nil {
		return domain.RencanaKinerja{}, err
	}

	// Hapus indikator yang akan dihapus
	queryDeleteIndikator := "DELETE FROM tb_indikator WHERE rencana_kinerja_id = ?"
	_, err = tx.ExecContext(ctx, queryDeleteIndikator, rencanaKinerja.Id)
	if err != nil {
		return domain.RencanaKinerja{}, err
	}

	// Insert indikator baru
	for _, indikator := range rencanaKinerja.Indikator {
		scriptIndikator := "INSERT INTO tb_indikator (id, rencana_kinerja_id, indikator, tahun) VALUES (?, ?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptIndikator, indikator.Id, rencanaKinerja.Id, indikator.Indikator, indikator.Tahun)
		if err != nil {
			return domain.RencanaKinerja{}, err
		}

		// Insert target untuk indikator
		for _, target := range indikator.Target {
			queryTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
			_, err := tx.ExecContext(ctx, queryTarget, target.Id, indikator.Id, target.Target, target.Satuan, target.Tahun)
			if err != nil {
				return domain.RencanaKinerja{}, err
			}
		}
	}

	return rencanaKinerja, nil
}

func (repository *RencanaKinerjaRepositoryImpl) FindIdRekinLevel1(ctx context.Context, tx *sql.Tx, id string) (domain.RencanaKinerja, error) {
	script := `
        SELECT 
            rk.id,
            rk.id_pohon,
            rk.sasaranopd_id,
            rk.nama_rencana_kinerja,
            rk.tahun,
            rk.status_rencana_kinerja,
            rk.catatan,
            rk.kode_opd,
            rk.pegawai_id,
            i.id as indikator_id,
            i.indikator,
            i.tahun as indikator_tahun,
            t.id as target_id,
            t.target,
            t.satuan,
            t.tahun as target_tahun,
            m.formula,
            m.sumber_data
        FROM tb_rencana_kinerja rk
        LEFT JOIN tb_indikator i ON rk.id = i.rencana_kinerja_id
        LEFT JOIN tb_target t ON i.id = t.indikator_id
        LEFT JOIN tb_manual_ik m ON i.id = m.indikator_id
        WHERE rk.id = ?`

	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return domain.RencanaKinerja{}, err
	}
	defer rows.Close()

	var rencanaKinerja domain.RencanaKinerja
	rencanaKinerja.Indikator = []domain.Indikator{}
	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var indikator domain.Indikator
		var target domain.Target
		var formula, sumberData sql.NullString
		var indikatorId, indikatorNama, indikatorTahun sql.NullString
		var targetId, targetNama, targetSatuan, targetTahun sql.NullString

		err := rows.Scan(
			&rencanaKinerja.Id,
			&rencanaKinerja.IdPohon,
			&rencanaKinerja.SasaranOpdId,
			&rencanaKinerja.NamaRencanaKinerja,
			&rencanaKinerja.Tahun,
			&rencanaKinerja.StatusRencanaKinerja,
			&rencanaKinerja.Catatan,
			&rencanaKinerja.KodeOpd,
			&rencanaKinerja.PegawaiId,
			&indikatorId,
			&indikatorNama,
			&indikatorTahun,
			&targetId,
			&targetNama,
			&targetSatuan,
			&targetTahun,
			&formula,
			&sumberData,
		)
		if err != nil {
			return domain.RencanaKinerja{}, err
		}

		// Jika tidak ada indikator, lanjutkan ke baris berikutnya
		if !indikatorId.Valid {
			continue
		}

		// Set nilai indikator
		indikator.Id = indikatorId.String
		indikator.Indikator = indikatorNama.String
		indikator.Tahun = indikatorTahun.String
		indikator.RumusPerhitungan = formula
		indikator.SumberData = sumberData

		// Cek apakah indikator sudah ada di map
		if existingIndikator, exists := indikatorMap[indikator.Id]; exists {
			// Tambahkan target ke indikator yang sudah ada jika ada target
			if targetId.Valid {
				target = domain.Target{
					Id:          targetId.String,
					Target:      targetNama.String,
					Satuan:      targetSatuan.String,
					Tahun:       targetTahun.String,
					IndikatorId: indikator.Id,
				}
				existingIndikator.Target = append(existingIndikator.Target, target)
			}
		} else {
			// Buat indikator baru
			indikator.Target = []domain.Target{}
			if targetId.Valid {
				target = domain.Target{
					Id:          targetId.String,
					Target:      targetNama.String,
					Satuan:      targetSatuan.String,
					Tahun:       targetTahun.String,
					IndikatorId: indikator.Id,
				}
				indikator.Target = append(indikator.Target, target)
			}
			indikatorMap[indikator.Id] = &indikator
		}
	}

	// Konversi map ke slice
	for _, indikator := range indikatorMap {
		rencanaKinerja.Indikator = append(rencanaKinerja.Indikator, *indikator)
	}

	return rencanaKinerja, nil
}

// func (repository *RencanaKinerjaRepositoryImpl) FindByPokinId(ctx context.Context, tx *sql.Tx, pokinId int) ([]domain.RencanaKinerja, error) {
// 	SQL := `
//     SELECT
//         rk.id,
//         rk.nama_rencana_kinerja,
//         rk.pegawai_id,
//         p.nama as nama_pegawai,
//         rk.id_pohon,
//         st.kode_subkegiatan,
//         sk.nama_subkegiatan,
//         -- Menggunakan SUBSTRING_INDEX untuk mengambil kode kegiatan dari kode_subkegiatan
//         SUBSTRING_INDEX(st.kode_subkegiatan, '.', 5) as kode_kegiatan,
//         k.nama_kegiatan
//     FROM tb_rencana_kinerja rk
//     LEFT JOIN tb_pegawai p ON p.nip = rk.pegawai_id
//     LEFT JOIN tb_subkegiatan_terpilih st ON st.rekin_id = rk.id
//     LEFT JOIN tb_subkegiatan sk ON sk.kode_subkegiatan = st.kode_subkegiatan
//     LEFT JOIN tb_master_kegiatan k ON k.kode_kegiatan = SUBSTRING_INDEX(st.kode_subkegiatan, '.', 5)
//     WHERE rk.id_pohon = ?
//     ORDER BY rk.id ASC
//     `

// 	rows, err := tx.QueryContext(ctx, SQL, pokinId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var rencanaKinerjas []domain.RencanaKinerja
// 	for rows.Next() {
// 		var rk domain.RencanaKinerja
// 		var namaPegawai, kodeSubkegiatan, namaSubkegiatan, kodeKegiatan, namaKegiatan sql.NullString
// 		err := rows.Scan(
// 			&rk.Id,
// 			&rk.NamaRencanaKinerja,
// 			&rk.PegawaiId,
// 			&namaPegawai,
// 			&rk.IdPohon,
// 			&kodeSubkegiatan,
// 			&namaSubkegiatan,
// 			&kodeKegiatan,
// 			&namaKegiatan,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Handle null values
// 		if namaPegawai.Valid {
// 			rk.NamaPegawai = namaPegawai.String
// 		}
// 		if kodeSubkegiatan.Valid {
// 			rk.KodeSubKegiatan = kodeSubkegiatan.String
// 		}
// 		if namaSubkegiatan.Valid {
// 			rk.NamaSubKegiatan = namaSubkegiatan.String
// 		}
// 		if kodeKegiatan.Valid {
// 			rk.KodeKegiatan = kodeKegiatan.String
// 		}
// 		if namaKegiatan.Valid {
// 			rk.NamaKegiatan = namaKegiatan.String
// 		}

// 		rencanaKinerjas = append(rencanaKinerjas, rk)
// 	}

// 	return rencanaKinerjas, nil
// }

func (repository *RencanaKinerjaRepositoryImpl) FindByPokinId(ctx context.Context, tx *sql.Tx, pokinId int) ([]domain.RencanaKinerja, error) {
	SQL := `
    SELECT 
        rk.id,
        rk.nama_rencana_kinerja,
        rk.pegawai_id,
        p.nama as nama_pegawai,
        rk.id_pohon,
        rk.tahun,
        COALESCE(rk.status_rencana_kinerja, '') as status_rencana_kinerja,
        COALESCE(rk.catatan, '') as catatan,
        rk.kode_opd,
        st.kode_subkegiatan,
        sk.nama_subkegiatan,
        -- Menggunakan SUBSTRING_INDEX untuk mengambil kode kegiatan dari kode_subkegiatan
        SUBSTRING_INDEX(st.kode_subkegiatan, '.', 5) as kode_kegiatan,
        k.nama_kegiatan
    FROM tb_rencana_kinerja rk
    LEFT JOIN tb_pegawai p ON p.nip = rk.pegawai_id
    LEFT JOIN tb_subkegiatan_terpilih st ON st.rekin_id = rk.id
    LEFT JOIN tb_subkegiatan sk ON sk.kode_subkegiatan = st.kode_subkegiatan
    LEFT JOIN tb_master_kegiatan k ON k.kode_kegiatan = SUBSTRING_INDEX(st.kode_subkegiatan, '.', 5)
    WHERE rk.id_pohon = ?
    ORDER BY rk.id ASC
    `

	rows, err := tx.QueryContext(ctx, SQL, pokinId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rencanaKinerjas []domain.RencanaKinerja
	for rows.Next() {
		var rk domain.RencanaKinerja
		var namaPegawai, tahun, statusRekin, catatan, kodeOpd, kodeSubkegiatan, namaSubkegiatan, kodeKegiatan, namaKegiatan sql.NullString
		err := rows.Scan(
			&rk.Id,
			&rk.NamaRencanaKinerja,
			&rk.PegawaiId,
			&namaPegawai,
			&rk.IdPohon,
			&tahun,
			&statusRekin,
			&catatan,
			&kodeOpd,
			&kodeSubkegiatan,
			&namaSubkegiatan,
			&kodeKegiatan,
			&namaKegiatan,
		)
		if err != nil {
			return nil, err
		}

		// Handle null values
		if namaPegawai.Valid {
			rk.NamaPegawai = namaPegawai.String
		}
		if tahun.Valid {
			rk.Tahun = tahun.String
		}
		if statusRekin.Valid {
			rk.StatusRencanaKinerja = statusRekin.String
		}
		if catatan.Valid {
			rk.Catatan = catatan.String
		}
		if kodeOpd.Valid {
			rk.KodeOpd = kodeOpd.String
		}
		if kodeSubkegiatan.Valid {
			rk.KodeSubKegiatan = kodeSubkegiatan.String
		}
		if namaSubkegiatan.Valid {
			rk.NamaSubKegiatan = namaSubkegiatan.String
		}
		if kodeKegiatan.Valid {
			rk.KodeKegiatan = kodeKegiatan.String
		}
		if namaKegiatan.Valid {
			rk.NamaKegiatan = namaKegiatan.String
		}

		rencanaKinerjas = append(rencanaKinerjas, rk)
	}

	return rencanaKinerjas, nil
}

func (repository *RencanaKinerjaRepositoryImpl) FindRekinLevel3(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.RencanaKinerja, error) {
	script := `
        SELECT DISTINCT 
            rk.id,
            rk.id_pohon,
            rk.nama_rencana_kinerja,
            rk.tahun,
            rk.status_rencana_kinerja,
            COALESCE(rk.catatan, ''),
            rk.kode_opd,
            rk.pegawai_id,
            rk.created_at,
            rk.kode_subkegiatan
        FROM tb_rencana_kinerja rk
        INNER JOIN tb_subkegiatan_terpilih st ON rk.id = st.rekin_id
        INNER JOIN tb_users u ON rk.pegawai_id = u.nip
        INNER JOIN tb_user_role ur ON u.id = ur.user_id
        INNER JOIN tb_role r ON ur.role_id = r.id
        WHERE r.role = 'level_3'
        AND rk.kode_opd = ?
        AND rk.tahun = ?
        ORDER BY rk.created_at ASC
    `

	rows, err := tx.QueryContext(ctx, script, kodeOpd, tahun)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data rencana kinerja level 3: %v", err)
	}
	defer rows.Close()

	var rencanaKinerjas []domain.RencanaKinerja
	for rows.Next() {
		var rk domain.RencanaKinerja
		err := rows.Scan(
			&rk.Id,
			&rk.IdPohon,
			&rk.NamaRencanaKinerja,
			&rk.Tahun,
			&rk.StatusRencanaKinerja,
			&rk.Catatan,
			&rk.KodeOpd,
			&rk.PegawaiId,
			&rk.CreatedAt,
			&rk.KodeSubKegiatan,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan data rencana kinerja: %v", err)
		}
		rencanaKinerjas = append(rencanaKinerjas, rk)
	}

	return rencanaKinerjas, nil
}

// func (repository *RencanaKinerjaRepositoryImpl) FindRekinAtasan(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.RencanaKinerja, error) {
// 	// Query untuk mendapatkan id_pohon dari rencana kinerja
// 	scriptGetPokin := `
//         SELECT id_pohon
//         FROM tb_rencana_kinerja
//         WHERE id = ?`

// 	var idPohon int
// 	err := tx.QueryRowContext(ctx, scriptGetPokin, rekinId).Scan(&idPohon)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, fmt.Errorf("rencana kinerja tidak ditemukan")
// 		}
// 		return nil, err
// 	}

// 	// Query untuk mendapatkan parent dari pohon kinerja
// 	scriptGetParent := `
//         SELECT parent
//         FROM tb_pohon_kinerja
//         WHERE id = ?`

// 	var parentId sql.NullInt64
// 	err = tx.QueryRowContext(ctx, scriptGetParent, idPohon).Scan(&parentId)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return []domain.RencanaKinerja{}, nil // Return empty slice instead of error
// 		}
// 		return nil, err
// 	}

// 	if !parentId.Valid {
// 		return []domain.RencanaKinerja{}, nil // Return empty slice if no parent
// 	}

// 	// Query untuk mendapatkan semua rencana kinerja atasan dan data pegawai
// 	scriptRekinAtasan := `
//         SELECT
//             rk.id,
//             rk.nama_rencana_kinerja,
//             rk.id_pohon,
//             rk.tahun,
//             rk.status_rencana_kinerja,
//             rk.catatan,
//             rk.kode_opd,
//             rk.pegawai_id,
//             p.nama as nama_pegawai,
//             p.nip as nip_pegawai
//         FROM tb_rencana_kinerja rk
//         INNER JOIN tb_pegawai p ON rk.pegawai_id = p.nip
//         WHERE rk.id_pohon = ?`

// 	rows, err := tx.QueryContext(ctx, scriptRekinAtasan, parentId.Int64)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var rekins []domain.RencanaKinerja
// 	for rows.Next() {
// 		var rekin domain.RencanaKinerja
// 		var pegawai domainmaster.Pegawai

// 		err := rows.Scan(
// 			&rekin.Id,
// 			&rekin.NamaRencanaKinerja,
// 			&rekin.IdPohon,
// 			&rekin.Tahun,
// 			&rekin.StatusRencanaKinerja,
// 			&rekin.Catatan,
// 			&rekin.KodeOpd,
// 			&rekin.PegawaiId,
// 			&pegawai.NamaPegawai,
// 			&pegawai.Nip,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		rekin.NamaPegawai = pegawai.NamaPegawai
// 		rekins = append(rekins, rekin)
// 	}

// 	if len(rekins) == 0 {
// 		return []domain.RencanaKinerja{}, nil
// 	}

// 	return rekins, nil
// }

func (repository *RencanaKinerjaRepositoryImpl) FindParentPokin(ctx context.Context, tx *sql.Tx, pokinId int) (domain.PohonKinerja, error) {
	script := `
		SELECT 
			parent_pk.id,
			COALESCE(parent_pk.nama_pohon, '') as nama_pohon,
			COALESCE(parent_pk.parent, 0) as parent,
			COALESCE(parent_pk.jenis_pohon, '') as jenis_pohon,
			COALESCE(parent_pk.level_pohon, 0) as level_pohon,
			COALESCE(parent_pk.kode_opd, '') as kode_opd,
			COALESCE(parent_pk.keterangan, '') as keterangan,
			COALESCE(parent_pk.keterangan_crosscutting, '') as keterangan_crosscutting,
			COALESCE(parent_pk.tahun, '') as tahun,
			COALESCE(parent_pk.status, '') as status,
			COALESCE(parent_pk.is_active) as is_active
		FROM tb_pohon_kinerja pk
		INNER JOIN tb_pohon_kinerja parent_pk ON pk.parent = parent_pk.id
		WHERE pk.id = ?
	`

	var pokin domain.PohonKinerja
	err := tx.QueryRowContext(ctx, script, pokinId).Scan(
		&pokin.Id,
		&pokin.NamaPohon,
		&pokin.Parent,
		&pokin.JenisPohon,
		&pokin.LevelPohon,
		&pokin.KodeOpd,
		&pokin.Keterangan,
		&pokin.KeteranganCrosscutting,
		&pokin.Tahun,
		&pokin.Status,
		&pokin.IsActive,
	)
	if err != nil {
		return domain.PohonKinerja{}, err
	}

	return pokin, nil
}
func (repository *RencanaKinerjaRepositoryImpl) ValidateRekinId(ctx context.Context, tx *sql.Tx, rekinId string) error {
	script := "SELECT id FROM tb_rencana_kinerja WHERE id = ?"

	var existingId string
	err := tx.QueryRowContext(ctx, script, rekinId).Scan(&existingId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("rencana kinerja dengan id %s tidak ditemukan", rekinId)
		}
		return err
	}

	return nil
}
func (repository *RencanaKinerjaRepositoryImpl) CloneRencanaKinerja(ctx context.Context, tx *sql.Tx, rekinId string, tahunBaru string) (domain.RencanaKinerja, error) {

	randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
	year := time.Now().Year()
	newRekinId := fmt.Sprintf("REKIN-PEG-%v-%v", year, randomDigits)

	script := `
		INSERT INTO tb_rencana_kinerja (
			id, id_pohon, nama_rencana_kinerja, tahun, 
			status_rencana_kinerja, catatan, kode_opd, pegawai_id, kode_subkegiatan, tahun_awal, tahun_akhir, jenis_periode
		)
		SELECT 
			?,
			0,
			nama_rencana_kinerja,
			?,
			'',
			'',
			kode_opd,
			pegawai_id,
			'',
			'',
			'',
			''
		FROM tb_rencana_kinerja
		WHERE id = ?
	`

	_, err := tx.ExecContext(ctx, script, newRekinId, tahunBaru, rekinId)
	if err != nil {
		return domain.RencanaKinerja{}, fmt.Errorf("gagal clone rencana kinerja: %v", err)
	}

	// Retrieve the cloned record menggunakan ID yang baru di-generate
	var newRekin domain.RencanaKinerja
	querySelect := `
		SELECT id, id_pohon, nama_rencana_kinerja, tahun, 
		       status_rencana_kinerja, catatan, kode_opd, pegawai_id, kode_subkegiatan, tahun_awal, tahun_akhir, jenis_periode
		FROM tb_rencana_kinerja
		WHERE id = ?
	`

	err = tx.QueryRowContext(ctx, querySelect, newRekinId).Scan(
		&newRekin.Id,
		&newRekin.IdPohon,
		&newRekin.NamaRencanaKinerja,
		&newRekin.Tahun,
		&newRekin.StatusRencanaKinerja,
		&newRekin.Catatan,
		&newRekin.KodeOpd,
		&newRekin.PegawaiId,
		&newRekin.KodeSubKegiatan,
		&newRekin.TahunAwal,
		&newRekin.TahunAkhir,
		&newRekin.JenisPeriode,
	)

	if err != nil {
		return domain.RencanaKinerja{}, fmt.Errorf("gagal mengambil data clone: %v", err)
	}

	return newRekin, nil
}

// CloneIndikator - Clone indikator dengan mapping indikator lama ke baru
func (repository *RencanaKinerjaRepositoryImpl) CloneIndikator(ctx context.Context, tx *sql.Tx, rekinIdLama string, rekinIdBaru string) error {
	script := `
		INSERT INTO tb_indikator (
			id, rencana_kinerja_id, indikator, tahun
		)
		SELECT 
			REPLACE(UUID(), '-', ''),
			?,
			indikator,
			(SELECT tahun FROM tb_rencana_kinerja WHERE id = ?)
		FROM tb_indikator
		WHERE rencana_kinerja_id = ?
	`

	_, err := tx.ExecContext(ctx, script, rekinIdBaru, rekinIdBaru, rekinIdLama)
	if err != nil {
		return fmt.Errorf("gagal clone indikator: %v", err)
	}

	return nil
}

// CloneTarget - Clone target untuk indikator baru
func (repository *RencanaKinerjaRepositoryImpl) CloneTarget(ctx context.Context, tx *sql.Tx, indikatorIdLama string, indikatorIdBaru string, tahunBaru string) error {
	script := `
		INSERT INTO tb_target (
			id, indikator_id, target, satuan, tahun
		)
		SELECT 
			REPLACE(UUID(), '-', ''),
			?,
			target,
			satuan,
			?
		FROM tb_target
		WHERE indikator_id = ?
	`

	_, err := tx.ExecContext(ctx, script, indikatorIdBaru, tahunBaru, indikatorIdLama)
	if err != nil {
		return fmt.Errorf("gagal clone target: %v", err)
	}

	return nil
}

// CloneRencanaAksi - Clone rencana aksi tanpa pelaksanaan
func (repository *RencanaKinerjaRepositoryImpl) CloneRencanaAksi(ctx context.Context, tx *sql.Tx, rekinIdLama string, rekinIdBaru string) error {
	script := `
		INSERT INTO tb_rencana_aksi (
			id, rencana_kinerja_id, kode_opd, urutan, nama_rencana_aksi
		)
		SELECT 
			REPLACE(UUID(), '-', ''),
			?,
			kode_opd,
			urutan,
			nama_rencana_aksi
		FROM tb_rencana_aksi
		WHERE rencana_kinerja_id = ?
	`

	_, err := tx.ExecContext(ctx, script, rekinIdBaru, rekinIdLama)
	if err != nil {
		return fmt.Errorf("gagal clone rencana aksi: %v", err)
	}

	return nil
}

// CloneDasarHukum - Clone dasar hukum
func (repository *RencanaKinerjaRepositoryImpl) CloneDasarHukum(ctx context.Context, tx *sql.Tx, rekinIdLama string, rekinIdBaru string) error {
	script := `
		INSERT INTO tb_dasar_hukum (
			id, rekin_id, urutan, peraturan_terkait, uraian
		)
		SELECT 
			REPLACE(UUID(), '-', ''),
			?,
			urutan,
			peraturan_terkait,
			uraian
		FROM tb_dasar_hukum
		WHERE rekin_id = ?
	`

	_, err := tx.ExecContext(ctx, script, rekinIdBaru, rekinIdLama)
	if err != nil {
		return fmt.Errorf("gagal clone dasar hukum: %v", err)
	}

	return nil
}

// CloneGambaranUmum - Clone gambaran umum
func (repository *RencanaKinerjaRepositoryImpl) CloneGambaranUmum(ctx context.Context, tx *sql.Tx, rekinIdLama string, rekinIdBaru string) error {
	script := `
		INSERT INTO tb_gambaran_umum (
			id, rekin_id, kode_opd, urutan, gambaran_umum
		)
		SELECT 
			REPLACE(UUID(), '-', ''),
			?,
			kode_opd,
			urutan,
			gambaran_umum
		FROM tb_gambaran_umum
		WHERE rekin_id = ?
	`

	_, err := tx.ExecContext(ctx, script, rekinIdBaru, rekinIdLama)
	if err != nil {
		return fmt.Errorf("gagal clone gambaran umum: %v", err)
	}

	return nil
}

// CloneInovasi - Clone inovasi
func (repository *RencanaKinerjaRepositoryImpl) CloneInovasi(ctx context.Context, tx *sql.Tx, rekinIdLama string, rekinIdBaru string) error {
	script := `
		INSERT INTO tb_inovasi (
			id, rekin_id, judul_inovasi, jenis_inovasi, gambaran_nilai_kebaruan
		)
		SELECT 
			REPLACE(UUID(), '-', ''),
			?,
			judul_inovasi,
			jenis_inovasi,
			gambaran_nilai_kebaruan
		FROM tb_inovasi
		WHERE rekin_id = ?
	`

	_, err := tx.ExecContext(ctx, script, rekinIdBaru, rekinIdLama)
	if err != nil {
		return fmt.Errorf("gagal clone inovasi: %v", err)
	}

	return nil
}

// ClonePermasalahan - Clone permasalahan
func (repository *RencanaKinerjaRepositoryImpl) ClonePermasalahan(ctx context.Context, tx *sql.Tx, rekinIdLama string, rekinIdBaru string) error {
	script := `
		INSERT INTO tb_permasalahan (
			rekin_id, permasalahan, penyebab_internal, penyebab_eksternal, jenis_permasalahan
		)
		SELECT 
			?,
			permasalahan,
			penyebab_internal,
			penyebab_eksternal,
			jenis_permasalahan
		FROM tb_permasalahan
		WHERE rekin_id = ?
	`

	_, err := tx.ExecContext(ctx, script, rekinIdBaru, rekinIdLama)
	if err != nil {
		return fmt.Errorf("gagal clone permasalahan: %v", err)
	}

	return nil
}

func (repository *RencanaKinerjaRepositoryImpl) CreateIndikatorClone(ctx context.Context, tx *sql.Tx, newIndikatorId string, rekinIdBaru string, indikator string, tahunBaru string) error {
	script := `
		INSERT INTO tb_indikator (
			id, rencana_kinerja_id, indikator, tahun
		) VALUES (?, ?, ?, ?)
	`

	_, err := tx.ExecContext(ctx, script, newIndikatorId, rekinIdBaru, indikator, tahunBaru)
	if err != nil {
		return fmt.Errorf("gagal insert indikator clone: %v", err)
	}

	return nil
}

func (repository *RencanaKinerjaRepositoryImpl) FindRekinByFilters(ctx context.Context, tx *sql.Tx, filter domain.FilterParams) ([]domain.RencanaKinerja, error) {
	if len(filter) == 0 {
		return nil, errors.New("filter tidak boleh kosong")
	}
	baseQuery := `
        SELECT
            rk.id,
            rk.id_pohon,
            pk.parent as id_pohon_atasan,
            rk.nama_rencana_kinerja,
            rk.tahun,
            rk.status_rencana_kinerja,
            rk.catatan,
            rk.kode_opd,
            rk.pegawai_id,
            i.id as indikator_id,
            i.indikator,
            i.tahun as indikator_tahun,
            t.id as target_id,
            t.target,
            t.satuan,
            t.tahun as target_tahun,
            m.formula,
            m.sumber_data,
            pk.nama_pohon,
            pk.level_pohon,
            pg.nama,
            opd.nama_opd
        FROM tb_rencana_kinerja rk
        LEFT JOIN tb_pegawai pg ON rk.pegawai_id = pg.nip
        LEFT JOIN tb_pohon_kinerja pk ON pk.id = rk.id_pohon
        INNER JOIN tb_pelaksana_pokin plp ON plp.pohon_kinerja_id = rk.id_pohon
        LEFT JOIN tb_indikator i ON rk.id = i.rencana_kinerja_id
        LEFT JOIN tb_target t ON i.id = t.indikator_id
        LEFT JOIN tb_manual_ik m ON i.id = m.indikator_id
        LEFT JOIN tb_operasional_daerah opd ON opd.kode_opd = rk.kode_opd
        WHERE 1=1`

	args := []any{}

	// Build WHERE dari filter
	if v, ok := filter["kode_opd"]; ok {
		baseQuery += " AND rk.kode_opd = ?"
		args = append(args, v)
	}

	if v, ok := filter["tahun"]; ok {
		baseQuery += " AND rk.tahun = ?"
		args = append(args, v)
	}
	// TAMBAHKAN FILTER LAIN JIKA PERLU
	if v, ok := filter["pegawai_id"]; ok {
		baseQuery += " AND rk.pegawai_id = ?"
		args = append(args, v)
	}

	if v, ok := filter["rekin_id"]; ok {
		baseQuery += " AND rk.id = ?"
		args = append(args, v)
	}

	rows, err := tx.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rencanaMap := make(map[string]*domain.RencanaKinerja)

	for rows.Next() {
		var (
			rk        domain.RencanaKinerja
			indikator domain.Indikator
			target    domain.Target

			indikatorId, indikatorNama, indikatorTahun      sql.NullString
			targetId, targetNama, targetSatuan, targetTahun sql.NullString
			formula, sumberData                             sql.NullString
			namaPohon, namaPegawai, namaOpd                 sql.NullString
			levelPohon, parentPohon                         sql.NullInt64
		)

		err := rows.Scan(
			&rk.Id,
			&rk.IdPohon,
			&parentPohon,
			&rk.NamaRencanaKinerja,
			&rk.Tahun,
			&rk.StatusRencanaKinerja,
			&rk.Catatan,
			&rk.KodeOpd,
			&rk.PegawaiId,
			&indikatorId,
			&indikatorNama,
			&indikatorTahun,
			&targetId,
			&targetNama,
			&targetSatuan,
			&targetTahun,
			&formula,
			&sumberData,
			&namaPohon,
			&levelPohon,
			&namaPegawai,
			&namaOpd,
		)
		if err != nil {
			return nil, err
		}

		existingRk, ok := rencanaMap[rk.Id]
		if !ok {
			rk.Indikator = []domain.Indikator{}
			if namaPohon.Valid {
				rk.NamaPohon = namaPohon.String
			}
			if levelPohon.Valid {
				rk.LevelPohon = int(levelPohon.Int64)
			}
			if parentPohon.Valid {
				rk.ParentPohon = int(parentPohon.Int64)
			}
			if namaPegawai.Valid {
				rk.NamaPegawai = namaPegawai.String
			}

			if namaOpd.Valid {
				rk.NamaOpd = namaOpd.String
			}

			rencanaMap[rk.Id] = &rk
			existingRk = &rk
		}

		if !indikatorId.Valid {
			continue
		}

		indikatorIndex, exists := -1, false

		for i := range existingRk.Indikator {
			if existingRk.Indikator[i].Id == indikatorId.String {
				indikatorIndex = i
				exists = true
				break
			}
		}

		if !exists {
			indikator = domain.Indikator{
				Id:               indikatorId.String,
				Indikator:        indikatorNama.String,
				Tahun:            indikatorTahun.String,
				RumusPerhitungan: formula,
				SumberData:       sumberData,
				Target:           []domain.Target{},
			}
			existingRk.Indikator = append(existingRk.Indikator, indikator)
			indikatorIndex = len(existingRk.Indikator) - 1
		}

		if targetId.Valid {
			target = domain.Target{
				Id:          targetId.String,
				Target:      targetNama.String,
				Satuan:      targetSatuan.String,
				Tahun:       targetTahun.String,
				IndikatorId: indikatorId.String,
			}
			existingRk.Indikator[indikatorIndex].Target =
				append(existingRk.Indikator[indikatorIndex].Target, target)
		}
	}

	results := make([]domain.RencanaKinerja, 0, len(rencanaMap))
	for _, rk := range rencanaMap {
		results = append(results, *rk)
	}

	return results, nil
}
