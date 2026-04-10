package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
)

type LembagaRepositoryImpl struct {
}

func NewLembagaRepositoryImpl() *LembagaRepositoryImpl {
	return &LembagaRepositoryImpl{}
}

func (repository *LembagaRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, lembaga domainmaster.Lembaga) domainmaster.Lembaga {
	script := "INSERT INTO tb_lembaga (id, kode_lembaga, nama_lembaga, nama_kepala_pemda, nip_kepala_pemda, jabatan_kepala_pemda) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, lembaga.Id, lembaga.KodeLembaga, lembaga.NamaLembaga, lembaga.NamaKepalaPemda, lembaga.NipKepalaPemda)
	if err != nil {
		return lembaga
	}
	return lembaga
}

func (repository *LembagaRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, lembaga domainmaster.Lembaga) domainmaster.Lembaga {
	script := "UPDATE tb_lembaga SET kode_lembaga = ?, nama_lembaga = ?, nama_kepala_pemda = ?, nip_kepala_pemda = ?, is_active = ?, jabatan_kepala_pemda = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, lembaga.KodeLembaga, lembaga.NamaLembaga, lembaga.NamaKepalaPemda, lembaga.NipKepalaPemda, lembaga.IsActive, lembaga.JabatanKepalaPemda, lembaga.Id)
	if err != nil {
		return lembaga
	}
	return lembaga
}

func (repository *LembagaRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	script := "DELETE FROM tb_lembaga WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *LembagaRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domainmaster.Lembaga, error) {
	script := "SELECT id, kode_lembaga, nama_lembaga, nama_kepala_pemda, nip_kepala_pemda, jabatan_kepala_pemda, is_active FROM tb_lembaga WHERE id = ?"
	var lembaga domainmaster.Lembaga
	var namaKepalaPemdaNs,
		nipKepalaPemdaNs,
		jabatanKepalaPemdaNs sql.NullString
	err := tx.QueryRowContext(ctx, script, id).Scan(&lembaga.Id,
		&lembaga.KodeLembaga,
		&lembaga.NamaLembaga,
		&namaKepalaPemdaNs,
		&nipKepalaPemdaNs,
		&jabatanKepalaPemdaNs,
		&lembaga.IsActive)

	if namaKepalaPemdaNs.Valid {
		lembaga.NamaKepalaPemda = namaKepalaPemdaNs.String
	}
	if nipKepalaPemdaNs.Valid {
		lembaga.NipKepalaPemda = nipKepalaPemdaNs.String
	}
	if jabatanKepalaPemdaNs.Valid {
		lembaga.JabatanKepalaPemda = jabatanKepalaPemdaNs.String
	}
	if err != nil {
		return domainmaster.Lembaga{}, err
	}
	return lembaga, nil
}

func (repository *LembagaRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domainmaster.Lembaga, error) {
	script := "SELECT id, kode_lembaga, nama_lembaga, nama_kepala_pemda, nip_kepala_pemda, jabatan_kepala_pemda, is_active FROM tb_lembaga"
	rows, err := tx.QueryContext(ctx, script)
	if err != nil {
		return []domainmaster.Lembaga{}, err
	}
	defer rows.Close()
	var lembagas []domainmaster.Lembaga
	for rows.Next() {
		lembaga := domainmaster.Lembaga{}
		var namaKepalaPemdaNs,
			nipKepalaPemdaNs,
			jabatanKepalaPemdaNs sql.NullString
		err := rows.Scan(
			&lembaga.Id,
			&lembaga.KodeLembaga,
			&lembaga.NamaLembaga,
			&namaKepalaPemdaNs,
			&nipKepalaPemdaNs,
			&jabatanKepalaPemdaNs,
			&lembaga.IsActive)
		if err != nil {
			return []domainmaster.Lembaga{}, err
		}

		if namaKepalaPemdaNs.Valid {
			lembaga.NamaKepalaPemda = namaKepalaPemdaNs.String
		}
		if nipKepalaPemdaNs.Valid {
			lembaga.NipKepalaPemda = nipKepalaPemdaNs.String
		}
		if jabatanKepalaPemdaNs.Valid {
			lembaga.JabatanKepalaPemda = jabatanKepalaPemdaNs.String
		}
		lembagas = append(lembagas, lembaga)
	}
	return lembagas, nil
}
