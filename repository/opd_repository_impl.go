package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
)

type OpdRepositoryImpl struct {
}

func NewOpdRepositoryImpl() *OpdRepositoryImpl {
	return &OpdRepositoryImpl{}
}

func (repository *OpdRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, opd domainmaster.Opd) (domainmaster.Opd, error) {
	script := `INSERT INTO tb_operasional_daerah (
		id, kode_opd, nama_opd, singkatan, alamat, telepon, fax, 
		email, website, nama_kepala_opd, nip_kepala_opd, pangkat_kepala, id_lembaga
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := tx.ExecContext(ctx, script,
		opd.Id, opd.KodeOpd, opd.NamaOpd, opd.Singkatan, opd.Alamat,
		opd.Telepon, opd.Fax, opd.Email, opd.Website, opd.NamaKepalaOpd,
		opd.NIPKepalaOpd, opd.PangkatKepala, opd.IdLembaga)
	if err != nil {
		return opd, err
	}
	return opd, nil
}

func (repository *OpdRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, opd domainmaster.Opd) (domainmaster.Opd, error) {
	script := `UPDATE tb_operasional_daerah SET 
		kode_opd = ?, nama_opd = ?, singkatan = ?, alamat = ?, 
		telepon = ?, fax = ?, email = ?, website = ?, 
		nama_kepala_opd = ?, nip_kepala_opd = ?, pangkat_kepala = ?, 
		id_lembaga = ? 
		WHERE id = ?`

	_, err := tx.ExecContext(ctx, script,
		opd.KodeOpd, opd.NamaOpd, opd.Singkatan, opd.Alamat,
		opd.Telepon, opd.Fax, opd.Email, opd.Website,
		opd.NamaKepalaOpd, opd.NIPKepalaOpd, opd.PangkatKepala,
		opd.IdLembaga, opd.Id)
	if err != nil {
		return opd, err
	}
	return opd, nil
}

func (repository *OpdRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, opdId string) error {
	script := "DELETE FROM tb_operasional_daerah WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, opdId)
	if err != nil {
		return err
	}
	return nil
}

func (repository *OpdRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domainmaster.Opd, error) {
	script := `SELECT 
		id, kode_opd, nama_opd, singkatan, alamat, telepon, fax,
		email, website, nama_kepala_opd, nip_kepala_opd, pangkat_kepala,
		id_lembaga 
		FROM tb_operasional_daerah`
	rows, err := tx.QueryContext(ctx, script)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var opds []domainmaster.Opd
	for rows.Next() {
		opd := domainmaster.Opd{}
		err := rows.Scan(
			&opd.Id, &opd.KodeOpd, &opd.NamaOpd, &opd.Singkatan,
			&opd.Alamat, &opd.Telepon, &opd.Fax, &opd.Email,
			&opd.Website, &opd.NamaKepalaOpd, &opd.NIPKepalaOpd,
			&opd.PangkatKepala, &opd.IdLembaga)
		if err != nil {
			return nil, err
		}
		opds = append(opds, opd)
	}
	return opds, nil
}

func (repository *OpdRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, opdId string) (domainmaster.Opd, error) {
	script := "SELECT id, kode_opd, nama_opd, singkatan, alamat, telepon, fax, email, website, nama_kepala_opd, nip_kepala_opd, pangkat_kepala, id_lembaga FROM tb_operasional_daerah WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, opdId)
	if err != nil {
		return domainmaster.Opd{}, err
	}
	defer rows.Close()

	var opd domainmaster.Opd
	if rows.Next() {
		err := rows.Scan(&opd.Id, &opd.KodeOpd, &opd.NamaOpd, &opd.Singkatan, &opd.Alamat, &opd.Telepon, &opd.Fax, &opd.Email, &opd.Website, &opd.NamaKepalaOpd, &opd.NIPKepalaOpd, &opd.PangkatKepala, &opd.IdLembaga)
		helper.PanicIfError(err)
	}
	return opd, nil
}

func (repository *OpdRepositoryImpl) FindByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) (domainmaster.Opd, error) {
	script := "SELECT id, kode_opd, nama_opd, singkatan, alamat, telepon, fax, email, website, nama_kepala_opd, nip_kepala_opd, pangkat_kepala, id_lembaga FROM tb_operasional_daerah WHERE kode_opd = ?"
	rows, err := tx.QueryContext(ctx, script, kodeOpd)
	if err != nil {
		return domainmaster.Opd{}, err
	}
	defer rows.Close()

	var opd domainmaster.Opd
	if rows.Next() {
		err := rows.Scan(&opd.Id, &opd.KodeOpd, &opd.NamaOpd, &opd.Singkatan, &opd.Alamat, &opd.Telepon, &opd.Fax, &opd.Email, &opd.Website, &opd.NamaKepalaOpd, &opd.NIPKepalaOpd, &opd.PangkatKepala, &opd.IdLembaga)
		helper.PanicIfError(err)
	}
	return opd, nil
}

// ... existing code ...

func (repository *OpdRepositoryImpl) FindAllWithLembaga(ctx context.Context, tx *sql.Tx) ([]domainmaster.Opd, map[string]domainmaster.Lembaga, error) {
	script := `SELECT 
		o.id, o.kode_opd, o.nama_opd, o.singkatan, o.alamat, o.telepon, o.fax,
		o.email, o.website, o.nama_kepala_opd, o.nip_kepala_opd, o.pangkat_kepala,
		o.id_lembaga,
		l.id as lembaga_id, l.kode_lembaga, l.nama_lembaga, l.is_active
		FROM tb_operasional_daerah o
		LEFT JOIN tb_lembaga l ON o.id_lembaga = l.id`

	rows, err := tx.QueryContext(ctx, script)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var opds []domainmaster.Opd
	lembagaMap := make(map[string]domainmaster.Lembaga)

	for rows.Next() {
		opd := domainmaster.Opd{}
		var lembagaId, kodeLembaga, namaLembaga sql.NullString
		var isActive sql.NullBool

		err := rows.Scan(
			&opd.Id, &opd.KodeOpd, &opd.NamaOpd, &opd.Singkatan,
			&opd.Alamat, &opd.Telepon, &opd.Fax, &opd.Email,
			&opd.Website, &opd.NamaKepalaOpd, &opd.NIPKepalaOpd,
			&opd.PangkatKepala, &opd.IdLembaga,
			&lembagaId, &kodeLembaga, &namaLembaga, &isActive,
		)
		if err != nil {
			return nil, nil, err
		}

		opds = append(opds, opd)

		// Simpan lembaga ke map jika ada
		if lembagaId.Valid && opd.IdLembaga != "" {
			if _, exists := lembagaMap[opd.IdLembaga]; !exists {
				lembagaMap[opd.IdLembaga] = domainmaster.Lembaga{
					Id:          lembagaId.String,
					KodeLembaga: kodeLembaga.String,
					NamaLembaga: namaLembaga.String,
					IsActive:    isActive.Bool,
				}
			}
		}
	}

	return opds, lembagaMap, nil
}
