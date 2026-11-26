package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
)

// TargetProgramKegiatanService — CRUD untuk tb_target_program_kegiatan.
// Parameter jenis pada Create menentukan fase: "renstra" / "renja" / "rankhir" / "penetapan".
// Update hanya mengizinkan perubahan pada kolom target dan satuan.
type TargetProgramKegiatanService interface {
	Create(ctx context.Context, jenis string, req programkegiatan.TargetProgramKegiatanCreateRequest) (programkegiatan.TargetProgramKegiatanResponse, error)
	Update(ctx context.Context, id int, req programkegiatan.TargetProgramKegiatanUpdateRequest) (programkegiatan.TargetProgramKegiatanResponse, error)
	Delete(ctx context.Context, id int) error
	FindById(ctx context.Context, id int) (programkegiatan.TargetProgramKegiatanResponse, error)
	FindByKodeIndikatorAndKodeOpd(ctx context.Context, kodeIndikator, kodeOpd string) ([]programkegiatan.TargetProgramKegiatanResponse, error)
}
