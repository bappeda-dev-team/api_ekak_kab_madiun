package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type IsuNasionalRepository interface {
	Create(ctx context.Context, tx *sql.Tx, ikk domain.IsuNasional) (domain.IsuNasional, error)
	Update(ctx context.Context, tx *sql.Tx, ikk domain.IsuNasional) (domain.IsuNasional, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.IsuNasional, error)
	FindAllById(ctx context.Context, tx *sql.Tx, id int) (domain.IsuNasional, error)
}