package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type TujuanPemdaRepository interface {
	Create(ctx context.Context, tx *sql.Tx, tujuanPemda domain.TujuanPemda) (domain.TujuanPemda, error)
	CreateIndikator(ctx context.Context, tx *sql.Tx, indikator domain.IndikatorPemda) (domain.IndikatorPemda, error)
	CreateTarget(ctx context.Context, tx *sql.Tx, target domain.TargetPemda) (domain.TargetPemda, error)
	Update(ctx context.Context, tx *sql.Tx, tujuanPemda domain.TujuanPemda) (domain.TujuanPemda, error)
	Delete(ctx context.Context, tx *sql.Tx, tujuanPemdaId int) error
	FindById(ctx context.Context, tx *sql.Tx, tujuanPemdaId int) (domain.TujuanPemda, error)
	FindAll(ctx context.Context, tx *sql.Tx, tahun string, jenisPeriode string) ([]domain.TujuanPemda, error)
	FindAllBetweenTahun(ctx context.Context, tx *sql.Tx, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]domain.TujuanPemda, error)
	DeleteIndikator(ctx context.Context, tx *sql.Tx, tujuanPemdaId int) error
	IsIdExists(ctx context.Context, tx *sql.Tx, id int) bool
	UpdatePeriode(ctx context.Context, tx *sql.Tx, tujuanPemda domain.TujuanPemda) (domain.TujuanPemda, error)
	FindAllWithPokin(ctx context.Context, tx *sql.Tx, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]domain.TujuanPemdaWithPokin, error)
	FindAllWithPokinRenstra(ctx context.Context, tx *sql.Tx, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]domain.TujuanPemdaWithPokin, error)
	IsPokinIdExists(ctx context.Context, tx *sql.Tx, pokinId int) (bool, error)
	TargetPemdaExistsByKey(ctx context.Context, tx *sql.Tx, kodeIndikator, tahun, jenis string) (bool, error)
	FindTargetPemdaById(ctx context.Context, tx *sql.Tx, id int) (domain.TargetPemda, error)
	UpdateTargetPemda(ctx context.Context, tx *sql.Tx, id int, target, satuan string) (domain.TargetPemda, error)

	//for rkpd
	FindAllWithPokinByTargetJenis(ctx context.Context, tx *sql.Tx, tahunAwal, tahunAkhir, jenisPeriode, targetJenis string) ([]domain.TujuanPemdaWithPokin, error)
	FindIndikatorPemdaByKode(ctx context.Context, tx *sql.Tx, kodeIndikator string) (domain.IndikatorPemda, error)
	UpsertTargetPemda(ctx context.Context, tx *sql.Tx, target domain.TargetPemda) (domain.TargetPemda, error)
	FindAllByTahun(ctx context.Context, tx *sql.Tx, tahun, jenisPeriode, targetJenis string) ([]domain.TujuanPemda, error)
}
