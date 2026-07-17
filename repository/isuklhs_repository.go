package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type IsuKlhsRepository interface {
	Create(ctx context.Context, tx *sql.Tx, ikk domain.IsuKlhs) (domain.IsuKlhs, error)
	Update(ctx context.Context, tx *sql.Tx, ikk domain.IsuKlhs) (domain.IsuKlhs, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.IsuKlhs, error)
	FindAllById(ctx context.Context, tx *sql.Tx, id int) (domain.IsuKlhs, error)
}