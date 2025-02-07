package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"log"
)

type SubKegiatanRepositoryImpl struct {
}

func NewSubKegiatanRepositoryImpl() *SubKegiatanRepositoryImpl {
	return &SubKegiatanRepositoryImpl{}
}

func (repository *SubKegiatanRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, subKegiatan domain.SubKegiatan) (domain.SubKegiatan, error) {
	scriptSubKegiatan := `INSERT INTO tb_subkegiatan (id, kode_subkegiatan, nama_subkegiatan, kode_opd, tahun, status) 
                         VALUES (?, ?, ?, ?, ?, ?)`

	_, err := tx.ExecContext(ctx, scriptSubKegiatan,
		subKegiatan.Id,
		subKegiatan.KodeSubKegiatan,
		subKegiatan.NamaSubKegiatan,
		subKegiatan.KodeOpd,
		subKegiatan.Tahun,
		subKegiatan.Status)
	if err != nil {
		return domain.SubKegiatan{}, err
	}

	for _, indikator := range subKegiatan.Indikator {
		scriptIndikator := `INSERT INTO tb_indikator (id, subkegiatan_id, indikator, tahun) 
						   VALUES (?, ?, ?, ?)`

		_, err = tx.ExecContext(ctx, scriptIndikator,
			indikator.Id,
			subKegiatan.Id,
			indikator.Indikator,
			indikator.Tahun)
		if err != nil {
			return domain.SubKegiatan{}, err
		}

		for _, target := range indikator.Target {
			scriptTarget := `INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) 
						   VALUES (?, ?, ?, ?, ?)`

			_, err = tx.ExecContext(ctx, scriptTarget,
				target.Id,
				indikator.Id,
				target.Target,
				target.Satuan,
				target.Tahun)
			if err != nil {
				return domain.SubKegiatan{}, err
			}
		}
	}

	return subKegiatan, nil
}

func (repository *SubKegiatanRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, subKegiatan domain.SubKegiatan) (domain.SubKegiatan, error) {
	// Update SubKegiatan
	scriptSubKegiatan := `UPDATE tb_subkegiatan 
                         SET nama_subkegiatan = ?, 
                             kode_opd = ?, 
                             tahun = ? 
                         WHERE id = ?`

	_, err := tx.ExecContext(ctx, scriptSubKegiatan,
		subKegiatan.NamaSubKegiatan,
		subKegiatan.KodeOpd,
		subKegiatan.Tahun,
		subKegiatan.Id)
	if err != nil {
		log.Printf("Error updating subkegiatan: %v", err)
		return domain.SubKegiatan{}, err
	}

	// Hapus indikator dan target yang lama
	scriptDeleteTarget := `DELETE FROM tb_target 
                          WHERE indikator_id IN (
                              SELECT id FROM tb_indikator 
                              WHERE subkegiatan_id = ?
                          )`
	_, err = tx.ExecContext(ctx, scriptDeleteTarget, subKegiatan.Id)
	if err != nil {
		log.Printf("Error deleting old targets: %v", err)
		return domain.SubKegiatan{}, err
	}

	scriptDeleteIndikator := `DELETE FROM tb_indikator WHERE subkegiatan_id = ?`
	_, err = tx.ExecContext(ctx, scriptDeleteIndikator, subKegiatan.Id)
	if err != nil {
		log.Printf("Error deleting old indicators: %v", err)
		return domain.SubKegiatan{}, err
	}

	// Insert indikator baru
	for _, indikator := range subKegiatan.Indikator {
		scriptIndikator := `INSERT INTO tb_indikator (id, subkegiatan_id, indikator, tahun) 
						   VALUES (?, ?, ?, ?)`

		_, err = tx.ExecContext(ctx, scriptIndikator,
			indikator.Id,
			subKegiatan.Id,
			indikator.Indikator,
			indikator.Tahun)
		if err != nil {
			log.Printf("Error inserting new indicator: %v", err)
			return domain.SubKegiatan{}, err
		}

		// Insert target baru untuk setiap indikator
		for _, target := range indikator.Target {
			scriptTarget := `INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) 
						   VALUES (?, ?, ?, ?, ?)`

			_, err = tx.ExecContext(ctx, scriptTarget,
				target.Id,
				indikator.Id,
				target.Target,
				target.Satuan,
				target.Tahun)
			if err != nil {
				log.Printf("Error inserting new target: %v", err)
				return domain.SubKegiatan{}, err
			}
		}
	}

	return subKegiatan, nil
}

func (repository *SubKegiatanRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, rekinId string, status string) ([]domain.SubKegiatan, error) {
	script := `SELECT id, kode_subkegiatan, nama_subkegiatan, kode_opd, rekin_id, status, tahun, created_at FROM tb_subkegiatan WHERE 1=1`
	var params []interface{}

	if kodeOpd != "" {
		script += ` AND kode_opd = ?`
		params = append(params, kodeOpd)
	}
	if rekinId != "" {
		script += ` AND rekin_id = ?`
		params = append(params, rekinId)
	}
	if status != "" {
		script += ` AND status = ?`
		params = append(params, status)
	}
	script += ` ORDER BY created_at ASC`

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subKegiatans []domain.SubKegiatan
	for rows.Next() {
		subKegiatan := domain.SubKegiatan{}
		err := rows.Scan(&subKegiatan.Id, &subKegiatan.KodeSubKegiatan, &subKegiatan.NamaSubKegiatan, &subKegiatan.KodeOpd, &subKegiatan.RekinId, &subKegiatan.Status, &subKegiatan.Tahun, &subKegiatan.CreatedAt)
		if err != nil {
			return nil, err
		}
		subKegiatans = append(subKegiatans, subKegiatan)
	}

	return subKegiatans, nil
}

func (repository *SubKegiatanRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, subKegiatanId string) (domain.SubKegiatan, error) {
	script := `SELECT id, kode_subkegiatan, rekin_id, nama_subkegiatan, kode_opd, status, tahun FROM tb_subkegiatan WHERE id = ?`

	rows, err := tx.QueryContext(ctx, script, subKegiatanId)
	if err != nil {
		return domain.SubKegiatan{}, err
	}
	defer rows.Close()

	subKegiatan := domain.SubKegiatan{}
	if rows.Next() {
		err := rows.Scan(&subKegiatan.Id, &subKegiatan.KodeSubKegiatan, &subKegiatan.RekinId, &subKegiatan.NamaSubKegiatan, &subKegiatan.KodeOpd, &subKegiatan.Status, &subKegiatan.Tahun)
		if err != nil {
			return domain.SubKegiatan{}, err
		}
		return subKegiatan, nil
	}

	return domain.SubKegiatan{}, fmt.Errorf("subkegiatan dengan id %s tidak ditemukan", subKegiatanId)
}

func (repository *SubKegiatanRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, subKegiatanId string) error {
	// Urutan query untuk menghapus data secara berurutan
	deleteQueries := []string{
		`DELETE FROM tb_target 
		 WHERE indikator_id IN (
			 SELECT id FROM tb_indikator 
			 WHERE subkegiatan_id = ?
		 )`,
		`DELETE FROM tb_indikator WHERE subkegiatan_id = ?`,
		`DELETE FROM tb_subkegiatan WHERE id = ?`,
	}

	// Eksekusi setiap query secara berurutan
	for _, query := range deleteQueries {
		_, err := tx.ExecContext(ctx, query, subKegiatanId)
		if err != nil {
			return fmt.Errorf("gagal menghapus data: %v", err)
		}
	}

	return nil
}

func (repository *SubKegiatanRepositoryImpl) FindIndikatorBySubKegiatanId(ctx context.Context, tx *sql.Tx, subKegiatanId string) ([]domain.Indikator, error) {
	script := "SELECT id, subkegiatan_id, indikator FROM tb_indikator WHERE subkegiatan_id = ?"
	params := []interface{}{subKegiatanId}

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indikators []domain.Indikator
	for rows.Next() {
		var indikator domain.Indikator
		err := rows.Scan(&indikator.Id, &indikator.SubKegiatanId, &indikator.Indikator)
		if err != nil {
			return nil, err
		}
		indikators = append(indikators, indikator)
	}

	return indikators, nil
}

func (repository *SubKegiatanRepositoryImpl) FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error) {
	script := "SELECT id, indikator_id, target, satuan FROM tb_target WHERE indikator_id = ?"
	params := []interface{}{indikatorId}

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []domain.Target
	for rows.Next() {
		var target domain.Target
		err := rows.Scan(&target.Id, &target.IndikatorId, &target.Target, &target.Satuan)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}

	return targets, nil
}

func (repository *SubKegiatanRepositoryImpl) FindByKodeSubKegiatan(ctx context.Context, tx *sql.Tx, kodeSubKegiatan string) (domain.SubKegiatan, error) {
	script := "SELECT id, kode_subkegiatan, nama_subkegiatan FROM tb_subkegiatan WHERE kode_subkegiatan = ?"
	rows, err := tx.QueryContext(ctx, script, kodeSubKegiatan)
	if err != nil {
		return domain.SubKegiatan{}, err
	}
	defer rows.Close()

	subKegiatan := domain.SubKegiatan{}
	if rows.Next() {
		err := rows.Scan(&subKegiatan.Id, &subKegiatan.KodeSubKegiatan, &subKegiatan.NamaSubKegiatan)
		if err != nil {
			return domain.SubKegiatan{}, err
		}
		return subKegiatan, nil
	}

	return domain.SubKegiatan{}, fmt.Errorf("subkegiatan dengan kode %s tidak ditemukan", kodeSubKegiatan)
}

func (repository *SubKegiatanRepositoryImpl) CreateRekin(ctx context.Context, tx *sql.Tx, idSubKegiatan string, rekinId string) error {
	script := "UPDATE tb_subkegiatan SET rekin_id = ?, status = 'subkegiatan_diambil' WHERE id = ?"
	result, err := tx.ExecContext(ctx, script, rekinId, idSubKegiatan)
	if err != nil {
		return fmt.Errorf("error saat mengupdate rekin subkegiatan: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error saat memeriksa rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subkegiatan dengan id %s tidak ditemukan", idSubKegiatan)
	}

	return nil
}

func (repository *SubKegiatanRepositoryImpl) DeleteSubKegiatanTerpilih(ctx context.Context, tx *sql.Tx, idSubKegiatan string) error {
	script := "UPDATE tb_subkegiatan SET rekin_id = '', status = 'belum_diambil' WHERE id = ?"
	result, err := tx.ExecContext(ctx, script, idSubKegiatan)
	if err != nil {
		return fmt.Errorf("error saat menghapus subkegiatan terpilih: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error saat memeriksa rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subkegiatan dengan id %s tidak ditemukan", idSubKegiatan)
	}

	return nil
}
