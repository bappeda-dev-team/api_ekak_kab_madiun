package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type InovasiRekinRepositoryImpl struct {
}

func NewInovasiRekinRepositoryImpl() *InovasiRekinRepositoryImpl {
	return &InovasiRekinRepositoryImpl{}
}

func (repository *InovasiRekinRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, inovasiRekin domain.InovasiRekin) (domain.InovasiRekin, error) {
	query := `INSERT INTO tb_inovasi_rekin 
	(id, rekin_id, kode_opd, nama_inovasi, jenis_inovasi_id, waktu_implementasi, instansi, inovator) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := tx.ExecContext(ctx, query, inovasiRekin.Id, inovasiRekin.RekinId, inovasiRekin.KodeOpd, inovasiRekin.NamaInovasi, inovasiRekin.JenisInovasiId, 
	inovasiRekin.WaktuImplementasi, inovasiRekin.Instansi, inovasiRekin.Inovator)
	if err != nil {
		return domain.InovasiRekin{}, err
	}
	return inovasiRekin, nil
}

func (repository *InovasiRekinRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, inovasiRekin domain.InovasiRekin) (domain.InovasiRekin, error) {
	query := `UPDATE tb_inovasi_rekin SET 
								nama_inovasi = ?, 
								jenis_inovasi_id = ?, 
								waktu_implementasi = ?, 
								instansi = ?, 
								inovator = ? 
								WHERE id = ?`
	_, err := tx.ExecContext(ctx, query, inovasiRekin.NamaInovasi, inovasiRekin.JenisInovasiId, inovasiRekin.WaktuImplementasi, inovasiRekin.Instansi, inovasiRekin.Inovator, inovasiRekin.Id)
	if err != nil {
		return domain.InovasiRekin{}, err
	}
	return inovasiRekin, nil
}

func (repository *InovasiRekinRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	query := "DELETE FROM tb_inovasi_rekin WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *InovasiRekinRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domain.InovasiRekin, error) {
	query := `SELECT 
	id, rekin_id, kode_opd, nama_inovasi, jenis_inovasi_id, waktu_implementasi, instansi, inovator 
	FROM tb_inovasi_rekin WHERE id = ?`
	row := tx.QueryRowContext(ctx, query, id)
	var inovasiRekin domain.InovasiRekin
	err := row.Scan(&inovasiRekin.Id, &inovasiRekin.RekinId, &inovasiRekin.KodeOpd, &inovasiRekin.NamaInovasi, &inovasiRekin.JenisInovasiId, &inovasiRekin.WaktuImplementasi, &inovasiRekin.Instansi, &inovasiRekin.Inovator)
	if err != nil {
		return domain.InovasiRekin{}, err
	}
	return inovasiRekin, nil
}

func (repository *InovasiRekinRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.InovasiRekin, error) {
	query := `SELECT 
	id, rekin_id, kode_opd, nama_inovasi, jenis_inovasi_id, waktu_implementasi, instansi, inovator 
	FROM tb_inovasi_rekin 
	WHERE rekin_id = ?`
	rows, err := tx.QueryContext(ctx, query, rekinId)
	if err != nil {
		return []domain.InovasiRekin{}, err
	}
	defer rows.Close()

	var inovasiRekinList []domain.InovasiRekin
	for rows.Next() {
		var inovasiRekin domain.InovasiRekin
		err := rows.Scan(&inovasiRekin.Id, &inovasiRekin.RekinId, &inovasiRekin.KodeOpd, &inovasiRekin.NamaInovasi, &inovasiRekin.JenisInovasiId, &inovasiRekin.WaktuImplementasi, &inovasiRekin.Instansi, &inovasiRekin.Inovator)
		if err != nil {
			return []domain.InovasiRekin{}, err
		}

		inovasiRekinList = append(inovasiRekinList, inovasiRekin)
	}

	err = rows.Err()
	if err != nil {
		return []domain.InovasiRekin{}, err
	}

	return inovasiRekinList, nil
}

