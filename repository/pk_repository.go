package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type PkRepository interface {
	FindByKodeOpdTahun(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun int) (map[int][]domain.PkOpd, error)
	HubungkanRekin(ctx context.Context, tx *sql.Tx, pkTerhubung domain.PkOpd) error
	FindSubkegiatanByKodeOpdTahunRekinIds(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun int, rekinIds []string) (map[string]domain.AllItemPk, error)
	FindTotalPaguAnggaranByRekinIds(ctx context.Context, tx *sql.Tx, rekinIds []string) (map[string]int, error)
	FindPaguPkByKodeSubkegiatans(ctx context.Context, tx *sql.Tx, kodeSubkegiatans []string) (map[string]int64, error)
	FindSasaranPemdaByTahun(ctx context.Context, tx *sql.Tx, tahun int) ([]domain.AllSasaranPemdaPk, error)
	FindSasaranPemdaById(ctx context.Context, tx *sql.Tx, sasaranPemdaId int) (domain.AllSasaranPemdaPk, error)
	// GROUPED BY KODE SUB - PAGU
	PaguPkByKodeOpdTahun(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun int) (map[string]int64, error)
	KunciPK(ctx context.Context, tx *sql.Tx, model domain.KunciPK) (int64, error)
	FindPkTerkunciByKodeOpdTahun(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun int) (map[string]bool, error)
	FindPkPegawaiPenetapan(ctx context.Context, tx *sql.Tx, idPegawai, kodeOpd string, tahun int) ([]domain.PkOpd, error)
	IndikatorTargetPkByIdRekins(ctx context.Context, tx *sql.Tx, idRekins []string) (map[string][]domain.Indikator, error)
}
