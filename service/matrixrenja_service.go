package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
)

type MatrixRenjaService interface {
	GetRenja(ctx context.Context, kodeOpd, tahun, jenisIndikator, jenisPagu string) ([]programkegiatan.UrusanDetailResponse, error)
	GetRenjaRankhir(ctx context.Context, kodeOpd string, tahun string) ([]programkegiatan.UrusanDetailResponse, error)
	UpsertBatchIndikatorRenja(ctx context.Context, requests []programkegiatan.IndikatorRenjaCreateRequest) ([]programkegiatan.IndikatorUpsertResponse, error)
	UpsertAnggaran(ctx context.Context, request programkegiatan.AnggaranRenjaRequest) (programkegiatan.AnggaranRenjaResponse, error)
	GetRenjaPenetapan(ctx context.Context, kodeOpd, tahun, jenisPagu string) ([]programkegiatan.UrusanDetailResponse, error)
}
