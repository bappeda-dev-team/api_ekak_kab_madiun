package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type IkkRepository interface {
	Create(ctx context.Context, tx *sql.Tx, ikk domain.IndikatorIkk) (domain.IndikatorIkk, error)
	Update(ctx context.Context, tx *sql.Tx, ikk domain.IndikatorIkk) (domain.IndikatorIkk, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.IndikatorIkk, error)
	FindByKodeOpd(ctx context.Context, tx *sql.Tx, jenis string, kodeOpd string) ([]domain.IndikatorIkk, error)
}