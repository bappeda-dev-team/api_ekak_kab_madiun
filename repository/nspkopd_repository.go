package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type NspkOpdRepository interface {
	Create(ctx context.Context, tx *sql.Tx, nspk domain.NspkOpd) (domain.NspkOpd, error)
	Update(ctx context.Context, tx *sql.Tx, nspk domain.NspkOpd) (domain.NspkOpd, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.NspkOpd, error)
	FindAllById(ctx context.Context, tx *sql.Tx, id int) (domain.NspkOpd, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.NspkOpd, error)
}