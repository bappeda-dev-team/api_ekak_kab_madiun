package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/programprioritaspusat"
)

type ProgramPrioritasPusatService interface {
	Create(ctx context.Context, request programprioritaspusat.ProgramPrioritasPusatCreateRequest) (programprioritaspusat.ProgramPrioritasPusatResponse, error)
	Update(ctx context.Context, request programprioritaspusat.ProgramPrioritasPusatUpdateRequest) (programprioritaspusat.ProgramPrioritasPusatResponse, error)
	Delete(ctx context.Context, id int) error
	FindById(ctx context.Context, id int) (programprioritaspusat.ProgramPrioritasPusatResponse, error)
	FindAll(ctx context.Context, tahunAwal string, tahunAkhir string) ([]programprioritaspusat.ProgramPrioritasPusatResponse, error)
	FindByKodeProgramPrioritasPusat(ctx context.Context, kodeProgramPrioritasPusat string) (programprioritaspusat.ProgramPrioritasPusatResponse, error)
	FindByTahun(ctx context.Context, tahun string) ([]programprioritaspusat.ProgramPrioritasPusatResponse, error)
	FindUnusedByTahun(ctx context.Context, tahun string) ([]programprioritaspusat.ProgramPrioritasPusatResponse, error)
	FindByIdTerkait(ctx context.Context, request programprioritaspusat.FindByIdTerkaitRequest) ([]programprioritaspusat.ProgramPrioritasPusatResponse, error)
}
