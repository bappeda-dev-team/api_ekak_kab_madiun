package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/programunggulan"
)

type ProgramUnggulanService interface {
	Create(ctx context.Context, request programunggulan.ProgramUnggulanCreateRequest) (programunggulan.ProgramUnggulanResponse, error)
	Update(ctx context.Context, request programunggulan.ProgramUnggulanUpdateRequest) (programunggulan.ProgramUnggulanResponse, error)
	Delete(ctx context.Context, id int) error
	FindById(ctx context.Context, id int) (programunggulan.ProgramUnggulanResponse, error)
	FindAll(ctx context.Context, tahunAwal string, tahunAkhir string) ([]programunggulan.ProgramUnggulanResponse, error)
	FindByKodeProgramUnggulan(ctx context.Context, kodeProgramUnggulan string) (programunggulan.ProgramUnggulanResponse, error)
	FindByTahun(ctx context.Context, tahun string) ([]programunggulan.ProgramUnggulanResponse, error)
	FindUnusedByTahun(ctx context.Context, tahun string) ([]programunggulan.ProgramUnggulanResponse, error)
	FindByIdTerkait(ctx context.Context, request programunggulan.FindByIdTerkaitRequest) ([]programunggulan.ProgramUnggulanResponse, error)
}
