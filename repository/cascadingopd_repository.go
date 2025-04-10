package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type CascadingOpdRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.PohonKinerja, error)
	FindIndikatorByPokinId(ctx context.Context, tx *sql.Tx, pokinId string) ([]domain.Indikator, error)
}
