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
	FindSasaranPemdaByTahun(ctx context.Context, tx *sql.Tx, tahun int) ([]domain.AllSasaranPemdaPk, error)
	FindSasaranPemdaById(ctx context.Context, tx *sql.Tx, sasaranPemdaId int) (domain.AllSasaranPemdaPk, error)
}
