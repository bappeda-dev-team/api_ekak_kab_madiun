package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type PermasalahanRekinRepository interface {
	Create(ctx context.Context, tx *sql.Tx, permasalahan domain.PermasalahanRekin) (domain.PermasalahanRekin, error)
	Update(ctx context.Context, tx *sql.Tx, permasalahan domain.PermasalahanRekin) (domain.PermasalahanRekin, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindAll(ctx context.Context, tx *sql.Tx, rekinId *string) ([]domain.PermasalahanRekin, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.PermasalahanRekin, error)
	FindByRekinIds(ctx context.Context, tx *sql.Tx, rekinIds []string) ([]domain.PermasalahanRekin, error)
	BatchCreate(ctx context.Context, tx *sql.Tx, permasalahans []domain.PermasalahanRekin) error
}
