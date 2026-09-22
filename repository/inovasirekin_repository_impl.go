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
	(id, rekin_id, kode_opd, nama_inovasi, jenis_inovasi_id, waktu_implementasi, instansi, inovator, kebaruan, asal_inovasi, tahun, nip_inovator) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := tx.ExecContext(ctx, query, inovasiRekin.Id, inovasiRekin.RekinId, inovasiRekin.KodeOpd, inovasiRekin.NamaInovasi, inovasiRekin.JenisInovasiId, 
	inovasiRekin.WaktuImplementasi, inovasiRekin.Instansi, inovasiRekin.Inovator, inovasiRekin.Kebaruan, inovasiRekin.AsalInovasi, inovasiRekin.Tahun, inovasiRekin.NipInovator)
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
								inovator = ?,
								kebaruan = ?,
								asal_inovasi = ?,
								nip_inovator = ?
								WHERE id = ?`
	_, err := tx.ExecContext(ctx, query, inovasiRekin.NamaInovasi, inovasiRekin.JenisInovasiId, inovasiRekin.WaktuImplementasi, inovasiRekin.Instansi, inovasiRekin.Inovator, inovasiRekin.Kebaruan, inovasiRekin.AsalInovasi, inovasiRekin.NipInovator, inovasiRekin.Id)
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
	tir.id, tir.rekin_id, tir.kode_opd, tir.nama_inovasi, tir.jenis_inovasi_id, ji.jenis, tir.waktu_implementasi, tir.instansi, tir.inovator,
	tir.kebaruan, tir.asal_inovasi, tir.tahun, tir.nip_inovator, COALESCE(tp.nama, '') AS nama_nip_inovator, COALESCE(tro.role, '') AS level
	FROM tb_inovasi_rekin tir
	LEFT JOIN tb_jenis_inovasi ji
		ON ji.id = tir.jenis_inovasi_id 
	LEFT JOIN tb_pegawai tp
		ON tp.nip = tir.nip_inovator 
	LEFT JOIN tb_users tu
		ON tu.nip = tir.nip_inovator
	LEFT JOIN tb_user_role tur
		ON tur.user_id = tu.id
	LEFT JOIN tb_role tro
		ON tro.id = tur.role_id
	WHERE tir.id = ?`
	row := tx.QueryRowContext(ctx, query, id)
	var inovasiRekin domain.InovasiRekin
	err := row.Scan(&inovasiRekin.Id, &inovasiRekin.RekinId, &inovasiRekin.KodeOpd, &inovasiRekin.NamaInovasi, &inovasiRekin.JenisInovasiId, &inovasiRekin.JenisInovasi,  &inovasiRekin.WaktuImplementasi, &inovasiRekin.Instansi, &inovasiRekin.Inovator, &inovasiRekin.Kebaruan, &inovasiRekin.AsalInovasi, &inovasiRekin.Tahun, &inovasiRekin.NipInovator, &inovasiRekin.NamaNipInovator, &inovasiRekin.Level)
	if err != nil {
		return domain.InovasiRekin{}, err
	}
	return inovasiRekin, nil
}

func (repository *InovasiRekinRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.InovasiRekin, error) {
	query := `SELECT 
	tir.id, tir.rekin_id, tir.kode_opd, od.nama_opd, tir.nama_inovasi, tir.jenis_inovasi_id, ji.jenis, tir.waktu_implementasi, tir.instansi, tir.inovator,
	tir.kebaruan, tir.asal_inovasi, tir.tahun, tir.nip_inovator, COALESCE(tp.nama, '') AS nama_nip_inovator, COALESCE(tro.role, '') AS level
	FROM tb_inovasi_rekin tir
	LEFT JOIN tb_jenis_inovasi ji
		ON ji.id = tir.jenis_inovasi_id 
	LEFT JOIN tb_operasional_daerah od
		ON od.kode_opd = tir.kode_opd
	LEFT JOIN tb_pegawai tp
		ON tp.nip = tir.nip_inovator
	LEFT JOIN tb_users tu
		ON tu.nip = tir.nip_inovator
	LEFT JOIN tb_user_role tur
		ON tur.user_id = tu.id
	LEFT JOIN tb_role tro 
		ON tro.id = tur.role_id
	WHERE tir.rekin_id = ?`
	rows, err := tx.QueryContext(ctx, query, rekinId)
	if err != nil {
		return []domain.InovasiRekin{}, err
	}
	defer rows.Close()

	var inovasiRekinList []domain.InovasiRekin
	for rows.Next() {
		var inovasiRekin domain.InovasiRekin
		err := rows.Scan(&inovasiRekin.Id, &inovasiRekin.RekinId, &inovasiRekin.KodeOpd, &inovasiRekin.NamaOpd, &inovasiRekin.NamaInovasi, &inovasiRekin.JenisInovasiId, &inovasiRekin.JenisInovasi, &inovasiRekin.WaktuImplementasi, &inovasiRekin.Instansi, &inovasiRekin.Inovator, &inovasiRekin.Kebaruan, &inovasiRekin.AsalInovasi, &inovasiRekin.Tahun, &inovasiRekin.NipInovator, &inovasiRekin.NamaNipInovator, &inovasiRekin.Level)
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

