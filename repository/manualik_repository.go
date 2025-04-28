package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type ManualIKRepository interface {
	Create(ctx context.Context, tx *sql.Tx, manualik domain.ManualIK) (domain.ManualIK, error)
	Update(ctx context.Context, tx *sql.Tx, manualik domain.ManualIK) (domain.ManualIK, error)
	// Delete(ctx context.Context, tx *sql.Tx, manualikId int) error
	// FindBy(ctx context.Context, tx *sql.Tx, manualikId int) ([]domain.ManualIK, error)
	// FindManualIKByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.ManualIK, error)
	GetManualIK(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.ManualIK, error)
	GetRencanaKinerjaWithTarget(ctx context.Context, tx *sql.Tx, indikatorId string) (domain.Indikator, domain.RencanaKinerja, []domain.Target, error)
	FindByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) (domain.ManualIK, error)
	FindManualIKSasaranOpdByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string, tahun string) (domain.ManualIK, error)
	DeleteByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) error
	IsIndikatorExist(ctx context.Context, tx *sql.Tx, indikatorId string) (bool, error)
}
