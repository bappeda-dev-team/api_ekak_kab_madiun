package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type InovasiRepository interface {
	Create(ctx context.Context, tx *sql.Tx, inovasi domain.Inovasi) (domain.Inovasi, error)
	Update(ctx context.Context, tx *sql.Tx, inovasi domain.Inovasi) (domain.Inovasi, error)
	FindAll(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.Inovasi, error)
	FindById(ctx context.Context, tx *sql.Tx, id string) (domain.Inovasi, error)
	Delete(ctx context.Context, tx *sql.Tx, id string) error
	FindByRekinIds(ctx context.Context, tx *sql.Tx, rekinIds []string) ([]domain.Inovasi, error)
	BatchCreate(ctx context.Context, tx *sql.Tx, inovasis []domain.Inovasi) error
}
