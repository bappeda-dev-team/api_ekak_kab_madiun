package repository

import (
	"context"
	"database/sql"

	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
)

type ProgramRepository interface {
	Create(ctx context.Context, tx *sql.Tx, program domainmaster.ProgramKegiatan) (domainmaster.ProgramKegiatan, error)
	Update(ctx context.Context, tx *sql.Tx, program domainmaster.ProgramKegiatan) (domainmaster.ProgramKegiatan, error)
	Delete(ctx context.Context, tx *sql.Tx, id string) error
	FindById(ctx context.Context, tx *sql.Tx, id string) (domainmaster.ProgramKegiatan, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]domainmaster.ProgramKegiatan, error)
	FindIndikatorByProgramId(ctx context.Context, tx *sql.Tx, programId string) ([]domain.Indikator, error)
	FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error)
	FindByKodeProgram(ctx context.Context, tx *sql.Tx, kodeProgram string) (domainmaster.ProgramKegiatan, error)
	FindIndikatorTargetByKodeRelasiBatch(ctx context.Context, tx *sql.Tx, kodeRelasis []string) (map[string][]domain.Indikator, error)
	FindMasterIndikatorProgramKegiatanByKodes(ctx context.Context, tx *sql.Tx, kodes []string) (map[string][]domainmaster.MasterIndikatorProgramKegiatan, error)

	// Upsert indikator program (create/update/delete dalam 1 request)
	FindIndikatorIdsByProgramId(ctx context.Context, tx *sql.Tx, programId string) ([]string, error)
	UpsertIndikator(ctx context.Context, tx *sql.Tx, ind domain.Indikator) error
	DeleteIndikatorWithTargets(ctx context.Context, tx *sql.Tx, indikatorId string) error
	ReplaceTargets(ctx context.Context, tx *sql.Tx, indikatorId string, targets []domain.Target) error
	FindIndikatorWithTargetsByProgramId(ctx context.Context, tx *sql.Tx, programId string) ([]domain.Indikator, error)
}
