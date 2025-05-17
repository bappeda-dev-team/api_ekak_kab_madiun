package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"log"
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

func (repository *RencanaKinerjaRepositoryImpl) FindByPokinId(ctx context.Context, tx *sql.Tx, pokinId int) ([]domain.RencanaKinerja, error) {
	SQL := `
    SELECT 
        rk.id,
        rk.nama_rencana_kinerja,
        rk.pegawai_id,
        p.nama as nama_pegawai,
        rk.id_pohon,
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
		var namaPegawai, kodeSubkegiatan, namaSubkegiatan, kodeKegiatan, namaKegiatan sql.NullString
		err := rows.Scan(
			&rk.Id,
			&rk.NamaRencanaKinerja,
			&rk.PegawaiId,
			&namaPegawai,
			&rk.IdPohon,
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
