package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type OutcomeMatrixRepository interface {
	Create(ctx context.Context, tx *sql.Tx, data domain.OutcomeMatrix) (domain.OutcomeMatrix, error)
	Update(ctx context.Context, tx *sql.Tx, data domain.OutcomeMatrix) (domain.OutcomeMatrix, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.OutcomeMatrix, error)
	FindAll(ctx context.Context, tx *sql.Tx, kode, kodeOpd, jenis string) ([]domain.OutcomeMatrix, error)
	FindByKodeAndKodeOpd(ctx context.Context, tx *sql.Tx, kodes []string, kodeOpd string) ([]domain.OutcomeMatrix, error)
}
