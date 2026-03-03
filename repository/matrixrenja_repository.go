package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type MatrixRenjaRepository interface {
	GetRenjaRanwal(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.SubKegiatanQuery, error)
	GetRenjaRankhir(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.SubKegiatanQuery, error)
	SaveTargetRenja(ctx context.Context, tx *sql.Tx, target domain.Target) error
	UpdateTargetRenja(ctx context.Context, tx *sql.Tx, target domain.Target) error
}
