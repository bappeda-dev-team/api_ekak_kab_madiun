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
	FindIndikatorByRekinId(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.Indikator, error)
	FindIndikatorSubkegiatanByKodeAndOpd(ctx context.Context, tx *sql.Tx, kodeSubkegiatan string, kodeOpd string, tahun string) ([]domain.Indikator, error)
	LaporanRincianBelanjaOpd(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.RincianBelanjaAsn, error)
	LaporanRincianBelanjaPegawai(ctx context.Context, tx *sql.Tx, pegawaiId string, tahun string) ([]domain.RincianBelanjaAsn, error)
	Upsert(ctx context.Context, tx *sql.Tx, rincianBelanja domain.RincianBelanja) (domain.RincianBelanja, error)
	TotalAnggaranByIdRekins(ctx context.Context, tx *sql.Tx, rekinIds []string) (map[string]int64, error)
}
