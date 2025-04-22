package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/kelompokanggarans"
)

type KelompokAnggaranService interface {
	Create(ctx context.Context, request kelompokanggarans.KelompokAnggaranCreateRequest) (kelompokanggarans.KelompokAnggaranResponse, error)
	Update(ctx context.Context, request kelompokanggarans.KelompokAnggaranUpdateRequest) (kelompokanggarans.KelompokAnggaranResponse, error)
	Delete(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (kelompokanggarans.KelompokAnggaranResponse, error)
	FindAll(ctx context.Context) ([]kelompokanggarans.KelompokAnggaranResponse, error)
}
