package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type MatrixRenjaRepository interface {
	GetByKodeOpdAndTahun(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.SubKegiatanQuery, error)
}
