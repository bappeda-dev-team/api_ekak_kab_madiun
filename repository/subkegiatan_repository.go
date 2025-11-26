package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type SubKegiatanRepository interface {
	Create(ctx context.Context, tx *sql.Tx, subKegiatan domain.SubKegiatan) (domain.SubKegiatan, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]domain.SubKegiatan, error)
	Update(ctx context.Context, tx *sql.Tx, subKegiatan domain.SubKegiatan) (domain.SubKegiatan, error)
	FindById(ctx context.Context, tx *sql.Tx, subKegiatanId string) (domain.SubKegiatan, error)
	Delete(ctx context.Context, tx *sql.Tx, subKegiatanId string) error
	FindIndikatorBySubKegiatanId(ctx context.Context, tx *sql.Tx, subKegiatanId string) ([]domain.Indikator, error)
	FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error)
	FindByKodeSubKegiatan(ctx context.Context, tx *sql.Tx, kodeSubKegiatan string) (domain.SubKegiatan, error)
	FindSubKegiatanKAK(ctx context.Context, tx *sql.Tx, kodeSubKegiatan string, kode string, tahun string) (domain.SubKegiatanKAKQuery, error)
	FindByKodeSubs(ctx context.Context, tx *sql.Tx, kodeSubKegiatans []string) ([]domain.SubKegiatan, error)
}
