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
	script := "INSERT INTO tb_rencana_kinerja (id, id_pohon, nama_rencana_kinerja, tahun, status_rencana_kinerja, catatan, kode_opd, pegawai_id, kode_subkegiatan) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, rencanaKinerja.Id, rencanaKinerja.IdPohon, rencanaKinerja.NamaRencanaKinerja, rencanaKinerja.Tahun, rencanaKinerja.StatusRencanaKinerja, rencanaKinerja.Catatan, rencanaKinerja.KodeOpd, rencanaKinerja.PegawaiId, rencanaKinerja.KodeSubKegiatan)
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
