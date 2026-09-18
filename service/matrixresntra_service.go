package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
)

type MatrixRenstraService interface {
	GetByKodeSubKegiatan(ctx context.Context, kodeOpd, tahunAwal, tahunAkhir string) ([]programkegiatan.UrusanDetailResponse, error)
	GetByKodeSubKegiatanVersiKedua(ctx context.Context, kodeOpd, tahunAwal, tahunAkhir string) ([]programkegiatan.UrusanDetailV2Response, error)
	CreateIndikatorV2(ctx context.Context, requests []programkegiatan.IndikatorRenstraV2CreateRequest) ([]programkegiatan.IndikatorV2UpsertResponse, error)
	UpsertTarget(ctx context.Context, request programkegiatan.TargetRenstraUpsertRequest) (programkegiatan.TargetResponse, error)
	UpsertBatchIndikator(ctx context.Context, requests []programkegiatan.IndikatorRenstraCreateRequest) ([]programkegiatan.IndikatorUpsertResponse, error)
	DeleteIndikator(ctx context.Context, kodeIndikator string) error
	FindIndikatorByKodeIndikator(ctx context.Context, kodeIndikator string) (programkegiatan.IndikatorResponse, error)
	UpsertAnggaran(ctx context.Context, request programkegiatan.AnggaranRenstraRequest) (programkegiatan.AnggaranRenstraResponse, error)
}
