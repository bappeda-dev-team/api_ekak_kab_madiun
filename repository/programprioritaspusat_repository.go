package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type ProgramPrioritasPusatRepository interface {
	Create(ctx context.Context, tx *sql.Tx, programPrioritasPusat domain.ProgramPrioritasPusat) (domain.ProgramPrioritasPusat, error)
	Update(ctx context.Context, tx *sql.Tx, programPrioritasPusat domain.ProgramPrioritasPusat) (domain.ProgramPrioritasPusat, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.ProgramPrioritasPusat, error)
	FindAll(ctx context.Context, tx *sql.Tx, tahunAwal string, tahunAkhir string) ([]domain.ProgramPrioritasPusat, error)
	FindByKodeProgramPrioritasPusat(ctx context.Context, tx *sql.Tx, kodeProgramPrioritasPusat string) (domain.ProgramPrioritasPusat, error)
	FindByTahun(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.ProgramPrioritasPusat, error)
	FindUnusedByTahun(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.ProgramPrioritasPusat, error)
	FindByIdTerkait(ctx context.Context, tx *sql.Tx, ids []int) ([]domain.ProgramPrioritasPusat, error)
	FindProgramPrioritasPusatByKodesBatch(ctx context.Context, tx *sql.Tx, kodes []string) (map[string]*domain.ProgramPrioritasPusat, error)
}
