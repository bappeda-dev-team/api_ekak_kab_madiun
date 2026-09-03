package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type NspkRepository interface {
	Create(ctx context.Context, tx *sql.Tx, nspk domain.MasterNSPK) (domain.MasterNSPK, error)
	Update(ctx context.Context, tx *sql.Tx, nspk domain.MasterNSPK) (domain.MasterNSPK, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.MasterNSPK, error)
	FindByIds(ctx context.Context, tx *sql.Tx, ids []int) ([]domain.MasterNSPK, error)
	FindAllById(ctx context.Context, tx *sql.Tx, id int) (domain.MasterNSPK, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.MasterNSPK, error)
}