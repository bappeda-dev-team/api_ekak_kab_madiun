package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"fmt"
	"strings"
)

type PegawaiRepositoryImpl struct {
}

func NewPegawaiRepositoryImpl() *PegawaiRepositoryImpl {
	return &PegawaiRepositoryImpl{}
}

func (repository *PegawaiRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, pegawai domainmaster.Pegawai) (domainmaster.Pegawai, error) {
	script := "INSERT INTO tb_pegawai (id, nama, nip, kode_opd) VALUES (?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, pegawai.Id, pegawai.NamaPegawai, pegawai.Nip, pegawai.KodeOpd)
	if err != nil {
		return domainmaster.Pegawai{}, err
	}
	return pegawai, nil
}

func (repository *PegawaiRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, pegawai domainmaster.Pegawai) domainmaster.Pegawai {
	script := "UPDATE tb_pegawai SET  nama = ?, nip = ?, kode_opd = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, pegawai.NamaPegawai, pegawai.Nip, pegawai.KodeOpd, pegawai.Id)
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
	script := "SELECT id, nama, nip, kode_opd FROM tb_pegawai WHERE id = ?"
	var pegawai domainmaster.Pegawai
	err := tx.QueryRowContext(ctx, script, id).Scan(&pegawai.Id, &pegawai.NamaPegawai, &pegawai.Nip, &pegawai.KodeOpd)
	if err != nil {
		return domainmaster.Pegawai{}, err
	}
	return pegawai, nil
}

func (repository *PegawaiRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domainmaster.Pegawai, error) {
	script := "SELECT id, nama, nip, kode_opd FROM tb_pegawai where 1=1"
	var params []interface{}

	if kodeOpd != "" {
		script += " AND kode_opd = ?"
		params = append(params, kodeOpd)
	}

	script += " ORDER BY nama ASC"
	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return []domainmaster.Pegawai{}, err
	}
	defer rows.Close()
	var pegawais []domainmaster.Pegawai
	for rows.Next() {
		pegawai := domainmaster.Pegawai{}
		err := rows.Scan(&pegawai.Id, &pegawai.NamaPegawai, &pegawai.Nip, &pegawai.KodeOpd)
		if err != nil {
			return []domainmaster.Pegawai{}, err
		}
		pegawais = append(pegawais, pegawai)
	}
	return pegawais, nil
}

func (repository *PegawaiRepositoryImpl) FindByNip(ctx context.Context, tx *sql.Tx, nip string) (domainmaster.Pegawai, error) {
	script := "SELECT id, nama, nip, kode_opd FROM tb_pegawai WHERE nip = ?"
	var pegawai domainmaster.Pegawai
	err := tx.QueryRowContext(ctx, script, nip).Scan(&pegawai.Id, &pegawai.NamaPegawai, &pegawai.Nip, &pegawai.KodeOpd)
	if err != nil {
		return domainmaster.Pegawai{}, err
	}
	return pegawai, nil
}

func (repository *PegawaiRepositoryImpl) FindPegawaiByNipsBatch(ctx context.Context, tx *sql.Tx, nips []string) (map[string]*domainmaster.Pegawai, error) {
	if len(nips) == 0 {
		return make(map[string]*domainmaster.Pegawai), nil
	}

	placeholders := make([]string, len(nips))
	args := make([]interface{}, len(nips))
	for i, nip := range nips {
		placeholders[i] = "?"
		args[i] = nip
	}

	script := fmt.Sprintf(`
		SELECT id, nip, nama
		FROM tb_pegawai
		WHERE nip IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]*domainmaster.Pegawai)
	for rows.Next() {
		var pegawai domainmaster.Pegawai
		err := rows.Scan(&pegawai.Id, &pegawai.Nip, &pegawai.NamaPegawai)
		if err != nil {
			return nil, err
		}
		result[pegawai.Nip] = &pegawai
	}

	return result, nil
}
