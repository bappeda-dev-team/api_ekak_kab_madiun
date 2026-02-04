package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
)

type OpdRepository interface {
	Create(ctx context.Context, tx *sql.Tx, opd domainmaster.Opd) (domainmaster.Opd, error)
	Update(ctx context.Context, tx *sql.Tx, opd domainmaster.Opd) (domainmaster.Opd, error)
	Delete(ctx context.Context, tx *sql.Tx, opdId string) error
	FindAll(ctx context.Context, tx *sql.Tx) ([]domainmaster.Opd, error)
	FindById(ctx context.Context, tx *sql.Tx, opdId string) (domainmaster.Opd, error)
	FindByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) (domainmaster.Opd, error)
	FindAllWithLembaga(ctx context.Context, tx *sql.Tx) ([]domainmaster.Opd, map[string]domainmaster.Lembaga, error)
}
