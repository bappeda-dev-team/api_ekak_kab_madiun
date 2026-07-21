package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type PpdRepository interface {
	Create(ctx context.Context, tx *sql.Tx, ppd domain.PotensiPerangkatDaerah) (domain.PotensiPerangkatDaerah, error)
	Update(ctx context.Context, tx *sql.Tx, ppd domain.PotensiPerangkatDaerah) (domain.PotensiPerangkatDaerah, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.PotensiPerangkatDaerah, error)
	FindAllById(ctx context.Context, tx *sql.Tx, id int) (domain.PotensiPerangkatDaerah, error)
	FindSelectionByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.BidangUrusanSelection, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.PotensiPerangkatDaerah, error)
}