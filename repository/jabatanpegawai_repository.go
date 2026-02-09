package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
)

type JabatanPegawaiRepository interface {
	TambahJabatanPegawai(ctx context.Context, tx *sql.Tx, jabatanPegawai domainmaster.JabatanPegawai) error
}
