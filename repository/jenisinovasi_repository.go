package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type JenisInovasiRepository interface {
	Create(ctx context.Context, tx *sql.Tx, ji domain.JenisInovasi) (domain.JenisInovasi, error)
	Update(ctx context.Context, tx *sql.Tx, ji domain.JenisInovasi) (domain.JenisInovasi, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.JenisInovasi, error)
	FindByIds(ctx context.Context, tx *sql.Tx, ids []int) ([]domain.JenisInovasi, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]domain.JenisInovasi, error)
}