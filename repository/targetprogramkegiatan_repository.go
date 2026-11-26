package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

// TargetProgramKegiatanRepository — CRUD untuk tb_target_program_kegiatan.
type TargetProgramKegiatanRepository interface {
	Create(ctx context.Context, tx *sql.Tx, t domain.TargetProgramKegiatan) (domain.TargetProgramKegiatan, error)
	Update(ctx context.Context, tx *sql.Tx, t domain.TargetProgramKegiatan) (domain.TargetProgramKegiatan, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.TargetProgramKegiatan, error)
	FindByKodeIndikatorAndKodeOpd(ctx context.Context, tx *sql.Tx, kodeIndikator, kodeOpd string) ([]domain.TargetProgramKegiatan, error)
}
