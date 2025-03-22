package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type PeriodeRepository interface {
	Save(ctx context.Context, tx *sql.Tx, periode domain.Periode) (domain.Periode, error)
	Update(ctx context.Context, tx *sql.Tx, periode domain.Periode) (domain.Periode, error)
	SaveTahunPeriode(ctx context.Context, tx *sql.Tx, tahunPeriode domain.TahunPeriode) error
	FindById(ctx context.Context, tx *sql.Tx, periodeId int) (domain.Periode, error)
	IsIdExists(ctx context.Context, tx *sql.Tx, id int) bool
	FindOverlappingPeriodes(ctx context.Context, tx *sql.Tx, tahunAwal, tahunAkhir, jenisPeriode string) ([]domain.Periode, error)
	DeleteTahunPeriode(ctx context.Context, tx *sql.Tx, periodeId int) error
	FindByTahun(ctx context.Context, tx *sql.Tx, tahun string) (domain.Periode, error)
	FindOverlappingPeriodesExcludeCurrent(ctx context.Context, tx *sql.Tx, currentId int, tahunAwal, tahunAkhir, jenisPeriode string) ([]domain.Periode, error)
	FindAll(ctx context.Context, tx *sql.Tx, jenis_periode string) ([]domain.Periode, error)
	Delete(ctx context.Context, tx *sql.Tx, periodeId int) error
	FindRPJMDByTahun(ctx context.Context, tx *sql.Tx, tahun string) (domain.Periode, error)
}
