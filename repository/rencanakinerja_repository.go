package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type RencanaKinerjaRepository interface {
	Create(ctx context.Context, tx *sql.Tx, rencanaKinerja domain.RencanaKinerja) (domain.RencanaKinerja, error)
	FindAll(ctx context.Context, tx *sql.Tx, pegawaiId string, kodeOPD string, tahun string) ([]domain.RencanaKinerja, error)
	FindIndikatorbyRekinId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Indikator, error)
	FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, targetId string) ([]domain.Target, error)
	FindById(ctx context.Context, tx *sql.Tx, id string, kodeOPD string, tahun string) (domain.RencanaKinerja, error)
	Update(ctx context.Context, tx *sql.Tx, rencanaKinerja domain.RencanaKinerja) (domain.RencanaKinerja, error)
	Delete(ctx context.Context, tx *sql.Tx, id string) error
	FindAllRincianKak(ctx context.Context, tx *sql.Tx, rencanakinerjaid, pegawaiId string) ([]domain.RencanaKinerja, error)
	FindRekinLevel3(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.RencanaKinerja, error)
	// FindRekinAtasan(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.RencanaKinerja, error)
	FindParentPokin(ctx context.Context, tx *sql.Tx, pokinId int) (domain.PohonKinerja, error)
	ValidateRekinId(ctx context.Context, tx *sql.Tx, rekinId string) error
	//sasaran opd
	CreateRekinLevel1(ctx context.Context, tx *sql.Tx, rencanaKinerja domain.RencanaKinerja) (domain.RencanaKinerja, error)
	UpdateRekinLevel1(ctx context.Context, tx *sql.Tx, rencanaKinerja domain.RencanaKinerja) (domain.RencanaKinerja, error)
	FindIdRekinLevel1(ctx context.Context, tx *sql.Tx, id string) (domain.RencanaKinerja, error)
	RekinsasaranOpd(ctx context.Context, tx *sql.Tx, pegawaiId string, kodeOPD string, tahun string) ([]domain.RencanaKinerja, error)
	FindIndikatorSasaranbyRekinId(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.Indikator, error)
	FindTargetByIndikatorIdAndTahun(ctx context.Context, tx *sql.Tx, indikatorId string, tahun string) ([]domain.Target, error)
	FindByPokinId(ctx context.Context, tx *sql.Tx, pokinId int) ([]domain.RencanaKinerja, error)

	// Method untuk clone
	CloneRencanaKinerja(ctx context.Context, tx *sql.Tx, rekinId string, tahunBaru string) (domain.RencanaKinerja, error)
	CloneIndikator(ctx context.Context, tx *sql.Tx, rekinIdLama string, rekinIdBaru string) error
	CloneTarget(ctx context.Context, tx *sql.Tx, indikatorIdLama string, indikatorIdBaru string, tahunBaru string) error
	CloneRencanaAksi(ctx context.Context, tx *sql.Tx, rekinIdLama string, rekinIdBaru string) error
	CloneDasarHukum(ctx context.Context, tx *sql.Tx, rekinIdLama string, rekinIdBaru string) error
	CloneGambaranUmum(ctx context.Context, tx *sql.Tx, rekinIdLama string, rekinIdBaru string) error
	CloneInovasi(ctx context.Context, tx *sql.Tx, rekinIdLama string, rekinIdBaru string) error
	ClonePermasalahan(ctx context.Context, tx *sql.Tx, rekinIdLama string, rekinIdBaru string) error
	CreateIndikatorClone(ctx context.Context, tx *sql.Tx, newIndikatorId string, rekinIdBaru string, indikator string, tahunBaru string) error

	FindRekinByFilters(ctx context.Context, tx *sql.Tx, filter domain.FilterParams) ([]domain.RencanaKinerja, error)
}
