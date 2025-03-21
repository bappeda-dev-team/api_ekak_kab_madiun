package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
)

type MatrixRenstraService interface {
	GetByKodeSubKegiatan(ctx context.Context, kodeOpd string, tahunAwal string, tahunAkhir string) ([]programkegiatan.UrusanDetailResponse, error)
	CreateIndikator(ctx context.Context, request programkegiatan.IndikatorRenstraCreateRequest) (programkegiatan.IndikatorResponse, error)
	UpdateIndikator(ctx context.Context, request programkegiatan.UpdateIndikatorRequest) (programkegiatan.IndikatorResponse, error)
	DeleteIndikator(ctx context.Context, indikatorId string) error
	FindIndikatorById(ctx context.Context, indikatorId string) (programkegiatan.IndikatorResponse, error)
}
