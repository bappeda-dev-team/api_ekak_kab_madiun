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
	FindByNipAndOpd(ctx context.Context, tx *sql.Tx, nip, kodeOpd, tahun string) ([]domain.SasaranOpd, error)
	FindSasaranByPeriod(ctx context.Context, tx *sql.Tx, kodeOpd, tahunAwal, tahunAkhir, jenisPeriode, jenisIndikator string) ([]domain.SasaranOpd, error)
	FindSasaranByTahun(ctx context.Context, tx *sql.Tx, kodeOpd, tahun, jenisPeriode, jenisIndikator string) ([]domain.SasaranOpd, error)
	CreateRenjaIndikator(ctx context.Context, tx *sql.Tx, sasaranOpdId int, indikators []domain.Indikator) error
	UpdateRenjaIndikator(ctx context.Context, tx *sql.Tx, indikators []domain.Indikator) error
	DeleteIndikatorTargetRenja(ctx context.Context, tx *sql.Tx, indikatorId string) error
	FindIndikatorByKodeIndikator(ctx context.Context, tx *sql.Tx, kodeIndikator string) (domain.Indikator, error)
	// FindSasaranTujuanByPokinIdsBatch mengambil ringkasan sasaran OPD (nama saja),
	// tujuan OPD, dan bidang urusan secara batch berdasarkan pokin_id (strategic level 4).
	FindSasaranTujuanByPokinIdsBatch(ctx context.Context, tx *sql.Tx, pokinIds []int) (map[int][]domain.SasaranPokinInfo, error)
	// GetIsHideByPokinIds mengambil nilai is_hide dari tb_sasaran_opd_view secara batch.
	// Key map adalah id_pokin; nilai true berarti is_hide=1, false berarti tidak ada entri atau is_hide=0.
	GetIsHideByPokinIds(ctx context.Context, tx *sql.Tx, pokinIds []int) (map[int]bool, error)
	HideSasaranOpdView(ctx context.Context, tx *sql.Tx, idPokin int) error
	UnhideSasaranOpdView(ctx context.Context, tx *sql.Tx, idPokin int) error
}
