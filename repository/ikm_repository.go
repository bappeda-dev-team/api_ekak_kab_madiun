package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type IkmRepository interface {
	FindAllByPeriode(ctx context.Context, tx *sql.Tx, tahunAwal, tahunAkhir string) ([]domain.IndikatorIkm, error)
	FindById(ctx context.Context, tx *sql.Tx, ikmId string) (domain.IndikatorIkm, error)
	ExistsById(ctx context.Context, tx *sql.Tx, ikmId string) (bool, error)
	Create(ctx context.Context, tx *sql.Tx, request domain.IndikatorIkm) (domain.IndikatorIkm, error)
	Update(ctx context.Context, tx *sql.Tx, request domain.IndikatorIkm, ikmId string) (domain.IndikatorIkm, error)
	Delete(ctx context.Context, tx *sql.Tx, ikmId string) error
}
