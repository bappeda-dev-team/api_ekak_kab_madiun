package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
)

type PegawaiRepositoryImpl struct {
}

func NewPegawaiRepositoryImpl() *PegawaiRepositoryImpl {
	return &PegawaiRepositoryImpl{}
}

func (repository *PegawaiRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, pegawai domainmaster.Pegawai) (domainmaster.Pegawai, error) {
	script := "INSERT INTO tb_pegawai (id, nama, nip ) VALUES (?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, pegawai.Id, pegawai.NamaPegawai, pegawai.Nip)
	if err != nil {
		return domainmaster.Pegawai{}, err
	}
	return pegawai, nil
}

func (repository *PegawaiRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, pegawai domainmaster.Pegawai) domainmaster.Pegawai {
	script := "UPDATE tb_pegawai SET  nama = ?, nip = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, pegawai.NamaPegawai, pegawai.Nip, pegawai.Id)
	if err != nil {
		return pegawai
	}

	return pegawai
}

func (repository *PegawaiRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	script := "DELETE FROM tb_pegawai WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *PegawaiRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domainmaster.Pegawai, error) {
	script := "SELECT id, nama, nip FROM tb_pegawai WHERE id = ?"
	var pegawai domainmaster.Pegawai
	err := tx.QueryRowContext(ctx, script, id).Scan(&pegawai.Id, &pegawai.NamaPegawai, &pegawai.Nip)
	if err != nil {
		return domainmaster.Pegawai{}, err
	}
	return pegawai, nil
}

func (repository *PegawaiRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domainmaster.Pegawai, error) {
	script := "SELECT id, nama, nip FROM tb_pegawai"
	rows, err := tx.QueryContext(ctx, script)
	if err != nil {
		return []domainmaster.Pegawai{}, err
	}
	defer rows.Close()
	var pegawais []domainmaster.Pegawai
	for rows.Next() {
		pegawai := domainmaster.Pegawai{}
		err := rows.Scan(&pegawai.Id, &pegawai.NamaPegawai, &pegawai.Nip)
		if err != nil {
			return []domainmaster.Pegawai{}, err
		}
		pegawais = append(pegawais, pegawai)
	}
	return pegawais, nil
}
