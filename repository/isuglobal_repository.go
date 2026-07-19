package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type IsuGlobalRepository interface {
	Create(ctx context.Context, tx *sql.Tx, isu domain.IsuGlobal) (domain.IsuGlobal, error)
	Update(ctx context.Context, tx *sql.Tx, isu domain.IsuGlobal) (domain.IsuGlobal, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.IsuGlobal, error)
	FindAllById(ctx context.Context, tx *sql.Tx, id int) (domain.IsuGlobal, error)
	FindSelectionByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.BidangUrusanSelection, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.IsuGlobal, error)
}