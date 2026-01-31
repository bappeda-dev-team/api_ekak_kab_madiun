package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
)

type JabatanPegawaiRepositoryImpl struct {
}

func NewJabatanPegawaiRepositoryImpl() *JabatanPegawaiRepositoryImpl {
	return &JabatanPegawaiRepositoryImpl{}
}

func (repository *JabatanPegawaiRepositoryImpl) TambahJabatanPegawai(
	ctx context.Context,
	tx *sql.Tx,
	jabatanPegawai domainmaster.JabatanPegawai,
) error {
	query := `
		INSERT INTO tb_jabatan_pegawai (
			id,
			id_jabatan,
			id_pegawai,
			status,
			is_active,
			bulan,
			tahun,
            kode_opd
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		jabatanPegawai.Id,
		jabatanPegawai.IdJabatan,
		jabatanPegawai.IdPegawai,
		jabatanPegawai.Status,
		jabatanPegawai.IsActive,
		jabatanPegawai.Bulan,
		jabatanPegawai.Tahun,
		jabatanPegawai.KodeOpd,
	)
	if err != nil {
		return err
	}

	return err
}
