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
	FindAll(ctx context.Context, tx *sql.Tx, kode, kodeSubkegiatan string) ([]domain.OutcomeMatrix, error)
	FindByKodes(ctx context.Context, tx *sql.Tx, kodes []string, kodeSubkegiatans []string) ([]domain.OutcomeMatrix, error)
}
