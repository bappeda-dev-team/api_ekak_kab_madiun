package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type RincianBelanjaRepository interface {
	Create(ctx context.Context, tx *sql.Tx, rincianBelanja domain.RincianBelanja) (domain.RincianBelanja, error)
	Update(ctx context.Context, tx *sql.Tx, rincianBelanja domain.RincianBelanja) (domain.RincianBelanja, error)
	FindByRenaksiId(ctx context.Context, tx *sql.Tx, renaksiId string) (domain.RincianBelanja, error)
	FindRincianBelanjaAsn(ctx context.Context, tx *sql.Tx, pegawaiId string, tahun string) ([]domain.RincianBelanjaAsn, error)
	FindAnggaranByRenaksiId(ctx context.Context, tx *sql.Tx, renaksiId string) (domain.RincianBelanja, error)
}
