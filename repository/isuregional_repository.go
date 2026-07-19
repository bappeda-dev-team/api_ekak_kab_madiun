package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type IsuRegionalRepository interface {
	Create(ctx context.Context, tx *sql.Tx, ikk domain.IsuRegional) (domain.IsuRegional, error)
	Update(ctx context.Context, tx *sql.Tx, ikk domain.IsuRegional) (domain.IsuRegional, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.IsuRegional, error)
	FindAllById(ctx context.Context, tx *sql.Tx, id int) (domain.IsuRegional, error)
	FindSelectionByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.BidangUrusanSelection, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.IsuRegional, error)
}