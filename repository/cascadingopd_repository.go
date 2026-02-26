package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type CascadingOpdRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.PohonKinerja, error)
	FindIndikatorByPokinId(ctx context.Context, tx *sql.Tx, pokinId string) ([]domain.Indikator, error)
	FindByKodeAndOpdAndTahun(ctx context.Context, tx *sql.Tx, kode string, kodeOpd string, tahun string) ([]domain.Indikator, error)
	FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error)
	FindIndikatorTargetByPokinIds(ctx context.Context, tx *sql.Tx, pokinIds []int) (map[int][]domain.Indikator, error)

	//by rekin
	FindPokinByRekinId(ctx context.Context, tx *sql.Tx, rekinId string) (domain.PohonKinerja, error)
	FindPokinById(ctx context.Context, tx *sql.Tx, pokinId int) (domain.PohonKinerja, error)
	FindStrategicByChildPokin(ctx context.Context, tx *sql.Tx, pokinId int) (domain.PohonKinerja, error)
	CalculateTotalAnggaranByPokinId(ctx context.Context, tx *sql.Tx, pokinId int) (int64, error)
	GetAnggaranByRekinId(ctx context.Context, tx *sql.Tx, rekinId string) (int64, error)
	FindOperationalChildrenByTacticalId(ctx context.Context, tx *sql.Tx, tacticalId int) ([]int, error)
	FindTacticalChildrenByStrategicId(ctx context.Context, tx *sql.Tx, strategicId int) ([]int, error)
	GetTotalAnggaranByPokinId(ctx context.Context, tx *sql.Tx, pokinId int) (int64, error)
	FindKodeSubkegiatanByPokinId(ctx context.Context, tx *sql.Tx, pokinId int) ([]string, error)
	FindKodeSubkegiatanFromChildren(ctx context.Context, tx *sql.Tx, pokinId int) ([]string, error)
	FindPokinByNipAndTahun(ctx context.Context, tx *sql.Tx, nip string, tahun string) ([]domain.PohonKinerja, error)
	GetTotalAnggaranByPokinIdWithPelaksana(ctx context.Context, tx *sql.Tx, pokinId int) (int64, error)
	FindTargetByIndikatorIdsBatch(ctx context.Context, tx *sql.Tx, indikatorIds []string) (map[string][]domain.Target, error)
}
