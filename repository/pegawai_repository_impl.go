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
	script := `SELECT
            peg.id,
            peg.nama,
            peg.nip,
            opd.kode_opd,
            opd.nama_opd,
            jab.nama_jabatan
            FROM tb_pegawai peg
            LEFT JOIN tb_operasional_daerah opd ON peg.kode_opd = opd.kode_opd
			LEFT JOIN tb_jabatan jab
				ON jab.id = (
					SELECT jp.id_jabatan
					FROM tb_jabatan_pegawai jp
					WHERE jp.id_pegawai = peg.nip
					ORDER BY jp.tahun DESC, jp.bulan DESC
					LIMIT 1
				)
            WHERE 1=1 `
	var params []any

	if kodeOpd != "" {
		script += " AND peg.kode_opd = ?"
		params = append(params, kodeOpd)
	}

	script += " ORDER BY peg.nama ASC"
	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return []domainmaster.Pegawai{}, err
	}
	defer rows.Close()
	var pegawais []domainmaster.Pegawai
	for rows.Next() {
		pegawai := domainmaster.Pegawai{}

		var kodeOpd, namaOpd sql.NullString
		var namaJabatan sql.NullString

		err := rows.Scan(
			&pegawai.Id,
			&pegawai.NamaPegawai,
			&pegawai.Nip,
			&kodeOpd,
			&namaOpd,
			&namaJabatan,
		)
		if err != nil {
			return []domainmaster.Pegawai{}, err
		}

		if kodeOpd.Valid {
			pegawai.KodeOpd = kodeOpd.String
		}
		if namaOpd.Valid {
			pegawai.NamaOpd = namaOpd.String
		}
		if namaJabatan.Valid {
			pegawai.NamaJabatan = namaJabatan.String
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

func (repository *PegawaiRepositoryImpl) FindByNipWithJabatan(ctx context.Context, tx *sql.Tx, nip string) (domainmaster.Pegawai, error) {
	script := `SELECT
            peg.id,
            peg.nama,
            peg.nip,
            opd.kode_opd,
            opd.nama_opd,
            jab.nama_jabatan
            FROM tb_pegawai peg
            LEFT JOIN tb_operasional_daerah opd ON peg.kode_opd = opd.kode_opd
            LEFT JOIN tb_jabatan_pegawai jp ON jp.id_pegawai = peg.nip
            LEFT JOIN tb_jabatan jab ON jab.id = jp.id_jabatan
            WHERE nip = ? `
	var pegawai domainmaster.Pegawai
	var kodeOpd, namaOpd sql.NullString
	var namaJabatan sql.NullString
	err := tx.QueryRowContext(ctx, script, nip).Scan(
		&pegawai.Id,
		&pegawai.NamaPegawai,
		&pegawai.Nip,
		&kodeOpd,
		&namaOpd,
		&namaJabatan,
	)
	if kodeOpd.Valid {
		pegawai.KodeOpd = kodeOpd.String
	}
	if namaOpd.Valid {
		pegawai.NamaOpd = namaOpd.String
	}
	if namaJabatan.Valid {
		pegawai.NamaJabatan = namaJabatan.String
	}
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
	args := make([]any, len(nips))
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
