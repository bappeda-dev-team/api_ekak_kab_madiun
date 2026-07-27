package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type SubKegiatanRepository interface {
	Create(ctx context.Context, tx *sql.Tx, subKegiatan domain.SubKegiatan) (domain.SubKegiatan, error)
	CountAll(ctx context.Context, tx *sql.Tx, kodeSubKegiatan, namaSubKegiatan string) (int, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeSubKegiatan, namaSubKegiatan string, limit, offset int) ([]domain.SubKegiatan, error)
	FindIndikatorsBySubKegiatanIds(ctx context.Context, tx *sql.Tx, subKegiatanIds []string) ([]domain.Indikator, error)
	FindTargetsByIndikatorIds(ctx context.Context, tx *sql.Tx, indikatorIds []string) ([]domain.Target, error)
	Update(ctx context.Context, tx *sql.Tx, subKegiatan domain.SubKegiatan) (domain.SubKegiatan, error)
	FindById(ctx context.Context, tx *sql.Tx, subKegiatanId string) (domain.SubKegiatan, error)
	Delete(ctx context.Context, tx *sql.Tx, subKegiatanId string) error
	FindIndikatorBySubKegiatanId(ctx context.Context, tx *sql.Tx, subKegiatanId string) ([]domain.Indikator, error)
	FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error)
	FindByKodeSubKegiatan(ctx context.Context, tx *sql.Tx, kodeSubKegiatan string) (domain.SubKegiatan, error)
	FindSubKegiatanKAK(ctx context.Context, tx *sql.Tx, kodeSubKegiatan string, kode string, tahun string) (domain.SubKegiatanKAKQuery, error)
}
