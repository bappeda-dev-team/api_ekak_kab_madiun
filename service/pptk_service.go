package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/pptk"
)

type PptkService interface {
	Create(ctx context.Context, request pptk.PptkCreateRequest) (pptk.PptkResponse, error)
	Update(ctx context.Context, request pptk.PptkUpdateRequest) (pptk.PptkResponse, error)
	Delete(ctx context.Context, id int) error
	FindById(ctx context.Context, id int) (pptk.PptkResponse, error)
	FindAll(ctx context.Context, kodeSubkegiatan string, kodeOpd string, tahun string) ([]pptk.PptkResponse, error)
	// FindAllByNip(ctx context.Context, nip string, tahun string) ([]pptk.PptkResponse, error)
}