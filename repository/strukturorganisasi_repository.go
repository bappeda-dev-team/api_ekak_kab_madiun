package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type StrukturOrganisasiRepository interface {
	Create(ctx context.Context, tx *sql.Tx, strukturOrganisasi domain.StrukturOrganisasi) error
	AtasanBawahanByKodeOpdTahun(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun int) (map[string]string, error)
}
