package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type RencanaAksiOpdRepository interface {
	FindBySasaranOpdAndTahun(ctx context.Context, tx *sql.Tx, sasaranOpdId int, tahun string) ([]domain.RencanaAksiOpd, error)
	SyncJadwalPelaksanaan(ctx context.Context, tx *sql.Tx, rekinId string) error
	Create(ctx context.Context, tx *sql.Tx, rencanaAksiOpd domain.RencanaAksiOpd) (domain.RencanaAksiOpd, error)
	Update(ctx context.Context, tx *sql.Tx, rencanaAksiOpd domain.RencanaAksiOpd) domain.RencanaAksiOpd
	Delete(ctx context.Context, tx *sql.Tx, Id int) error
	FindById(ctx context.Context, tx *sql.Tx, Id int) (domain.RencanaAksiOpd, error)
	FindAllSasaranByTahun(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.SasaranOpdDetailRenaksi, error)
}
