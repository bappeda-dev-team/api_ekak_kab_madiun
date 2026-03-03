package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
)

type MatrixRenjaService interface {
	GetRenjaRanwal(ctx context.Context, kodeOpd string, tahun string) ([]programkegiatan.UrusanDetailResponse, error)
	GetRenjaRankhir(ctx context.Context, kodeOpd string, tahun string) ([]programkegiatan.UrusanDetailResponse, error)
	CreateOrUpdateTarget(ctx context.Context, request programkegiatan.TargetRenjaRequest) (programkegiatan.TargetResponse, error)
}
