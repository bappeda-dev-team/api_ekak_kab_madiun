package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type PptkRepository interface {
	Create(ctx context.Context, tx *sql.Tx, pptk domain.Pptk) (domain.Pptk, error)
	Update(ctx context.Context, tx *sql.Tx, pptk domain.Pptk) (domain.Pptk, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.Pptk, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.Pptk, error)
	FindAllByNip(ctx context.Context, tx *sql.Tx, nip string, tahun string) ([]domain.Pptk, error)
}