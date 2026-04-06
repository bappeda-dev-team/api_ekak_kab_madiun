package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type DasarHukumRepository interface {
	Create(ctx context.Context, tx *sql.Tx, dasarHukum domain.DasarHukum) (domain.DasarHukum, error)
	Update(ctx context.Context, tx *sql.Tx, dasarHukum domain.DasarHukum) (domain.DasarHukum, error)
	FindAll(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.DasarHukum, error)
	FindById(ctx context.Context, tx *sql.Tx, id string) (domain.DasarHukum, error)
	Delete(ctx context.Context, tx *sql.Tx, id string) error
	GetLastUrutan(ctx context.Context, tx *sql.Tx) (int, error)
	GetLastUrutanByRekinId(ctx context.Context, tx *sql.Tx, rekinId string) (int, error)
	FindByRekinIds(ctx context.Context, tx *sql.Tx, rekinIds []string) ([]domain.DasarHukum, error)
	BatchCreate(ctx context.Context, tx *sql.Tx, dasarHukums []domain.DasarHukum) error
}
