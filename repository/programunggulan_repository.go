package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type ProgramUnggulanRepository interface {
	Create(ctx context.Context, tx *sql.Tx, programUnggulan domain.ProgramUnggulan) (domain.ProgramUnggulan, error)
	Update(ctx context.Context, tx *sql.Tx, programUnggulan domain.ProgramUnggulan) (domain.ProgramUnggulan, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.ProgramUnggulan, error)
	FindAll(ctx context.Context, tx *sql.Tx, tahunAwal string, tahunAkhir string) ([]domain.ProgramUnggulan, error)
	FindByKodeProgramUnggulan(ctx context.Context, tx *sql.Tx, kodeProgramUnggulan string) (domain.ProgramUnggulan, error)
	FindByTahun(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.ProgramUnggulan, error)
	FindUnusedByTahun(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.ProgramUnggulan, error)
	FindByIdTerkait(ctx context.Context, tx *sql.Tx, ids []int) ([]domain.ProgramUnggulan, error)
	FindProgramUnggulanByKodesBatch(ctx context.Context, tx *sql.Tx, kodes []string) (map[string]*domain.ProgramUnggulan, error)
	// Method baru untuk OPD
	CreateOpdProgramUnggulan(ctx context.Context, tx *sql.Tx, kodeProgramUnggulan string, kodeOpd []string) error
	FindOpdByKodeProgramUnggulanAndKodeOpds(ctx context.Context, tx *sql.Tx, kodeProgramUnggulan string, kodeOpds []string) ([]domain.OpdProgramUnggulan, error)
	DeleteOpdProgramUnggulan(ctx context.Context, tx *sql.Tx, id int) error
	FindOpdByKodeProgramUnggulan(ctx context.Context, tx *sql.Tx, kodeProgramUnggulan string) ([]domain.OpdProgramUnggulan, error)
	FindTahunTerpakai(ctx context.Context, tx *sql.Tx, kodeProgramUnggulan string) ([]string, error)
	// Modifikasi FindByTahun untuk support filter kode_opd
	FindByTahunAndKodeOpd(ctx context.Context, tx *sql.Tx, tahun string, kodeOpd string) ([]domain.ProgramUnggulan, error)
	FindOpdProgramUnggulanById(ctx context.Context, tx *sql.Tx, id int) ([]domain.OpdProgramUnggulan, error)
}
