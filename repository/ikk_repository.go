package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type IkkRepository interface {
	Create(ctx context.Context, tx *sql.Tx, ikk domain.Ikk) (domain.Ikk, error)
	Update(ctx context.Context, tx *sql.Tx, ikk domain.Ikk) (domain.Ikk, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.Ikk, error)
	FindByKodeOpd(ctx context.Context, tx *sql.Tx, jenis string, kodeOpd string) ([]domain.Ikk, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.Ikk, error)
	FindAllByJenisAndKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string, jenis string) ([]domain.Ikk, error)
	FindSelection(ctx context.Context, tx *sql.Tx) ([]domain.BidangUrusanSelection, error)
	FindSelectionByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.BidangUrusanSelection, error)
}