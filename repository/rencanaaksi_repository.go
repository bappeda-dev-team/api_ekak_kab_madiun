package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type RencanaAksiRepository interface {
	Create(ctx context.Context, tx *sql.Tx, rencanaAksi domain.RencanaAksi) (domain.RencanaAksi, error)
	Update(ctx context.Context, tx *sql.Tx, rencanaAksi domain.RencanaAksi) (domain.RencanaAksi, error)
	Delete(ctx context.Context, tx *sql.Tx, id string) error
	FindById(ctx context.Context, tx *sql.Tx, id string) (domain.RencanaAksi, error)
	FindAll(ctx context.Context, tx *sql.Tx, rencanaKinerjaId string) ([]domain.RencanaAksi, error)
	IsUrutanExistsForRencanaKinerja(ctx context.Context, tx *sql.Tx, rencanaKinerjaId string, urutan int) (bool, error)
	IsUrutanExistsForRencanaKinerjaExcludingId(ctx context.Context, tx *sql.Tx, rencanaKinerjaId string, urutan int, excludeId string) (bool, error)
	GetTotalBobotForRencanaKinerja(ctx context.Context, tx *sql.Tx, rencanaKinerjaId string) (int, error)
	FindRenaksiByRekinIds(ctx context.Context, tx *sql.Tx, rekinIds []string) ([]domain.RencanaAksi, error)
	BatchCreate(ctx context.Context, tx *sql.Tx, renaksis []domain.RencanaAksi) error
}
