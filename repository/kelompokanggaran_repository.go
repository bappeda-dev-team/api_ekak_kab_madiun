package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type KelompokAnggaranRepository interface {
	Create(ctx context.Context, tx *sql.Tx, ka domain.KelompokAnggaran) (domain.KelompokAnggaran, error)
	Update(ctx context.Context, tx *sql.Tx, ka domain.KelompokAnggaran) (domain.KelompokAnggaran, error)
	FindAll(ctx context.Context, tx *sql.Tx) []domain.KelompokAnggaran
	FindById(ctx context.Context, tx *sql.Tx, id string) (domain.KelompokAnggaran, error)
	Delete(ctx context.Context, tx *sql.Tx, id string) error
}
