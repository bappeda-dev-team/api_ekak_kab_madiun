package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type SasaranOpdRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx, KodeOpd string, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]domain.SasaranOpd, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (*domain.SasaranOpd, error)
	FindIdPokinSasaran(ctx context.Context, tx *sql.Tx, id int) (domain.PohonKinerja, error)
	FindByIdSasaran(ctx context.Context, tx *sql.Tx, id int) (*domain.SasaranOpdDetail, error)
	Create(ctx context.Context, tx *sql.Tx, domain domain.SasaranOpdDetail) error
	Update(ctx context.Context, tx *sql.Tx, sasaranOpd domain.SasaranOpdDetail) (domain.SasaranOpdDetail, error)
	Delete(ctx context.Context, tx *sql.Tx, id string) error
	FindByIdPokin(ctx context.Context, tx *sql.Tx, idPokin int, tahun string) (*domain.SasaranOpd, error)
	FindByTahun(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string, jenisPeriode string) ([]domain.SasaranOpd, error)
	FindSasaranByPeriod(ctx context.Context, tx *sql.Tx, kodeOpd, tahunAwal, tahunAkhir, jenisPeriode, jenisIndikator string) ([]domain.SasaranOpd, error)
	FindSasaranByTahun(ctx context.Context, tx *sql.Tx, kodeOpd, tahun, jenisPeriode, jenisIndikator string) ([]domain.SasaranOpd, error)
	CreateRenjaIndikator(ctx context.Context, tx *sql.Tx, sasaranOpdId int, indikators []domain.Indikator) error
	UpdateRenjaIndikator(ctx context.Context, tx *sql.Tx, indikators []domain.Indikator) error
	DeleteIndikatorTargetRenja(ctx context.Context, tx *sql.Tx, indikatorId string) error
	FindIndikatorByKodeIndikator(ctx context.Context, tx *sql.Tx, kodeIndikator string) (domain.Indikator, error)
}
