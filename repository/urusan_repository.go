package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
)

type UrusanRepository interface {
	Create(ctx context.Context, tx *sql.Tx, urusan domainmaster.Urusan) (domainmaster.Urusan, error)
	Update(ctx context.Context, tx *sql.Tx, urusan domainmaster.Urusan) (domainmaster.Urusan, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]domainmaster.Urusan, error)
	FindById(ctx context.Context, tx *sql.Tx, id string) (domainmaster.Urusan, error)
	Delete(ctx context.Context, tx *sql.Tx, id string) error
	FindByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domainmaster.Urusan, error)
	FindUrusanAndBidangByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domainmaster.Urusan, error)
	FindByKodeUrusan(ctx context.Context, tx *sql.Tx, kodeUrusan string) (domainmaster.Urusan, error)
	FindByKodeSubKegiatan(ctx context.Context, tx *sql.Tx, kodeSubKegiatan string) (domainmaster.Urusan, error)
}
