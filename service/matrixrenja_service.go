package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
)

type MatrixRenjaService interface {
	GetRenja(ctx context.Context, kodeOpd, tahun, jenisIndikator, jenisPagu string) ([]programkegiatan.UrusanDetailResponse, error)
	GetRenjaRankhir(ctx context.Context, kodeOpd string, tahun string) ([]programkegiatan.UrusanDetailResponse, error)
	UpsertBatchIndikatorRenja(ctx context.Context, request programkegiatan.BatchIndikatorRenjaRequest) (programkegiatan.BatchIndikatorRenjaResponse, error)
	UpsertAnggaran(ctx context.Context, request programkegiatan.AnggaranRenjaRequest) (programkegiatan.AnggaranRenjaResponse, error)
}
