package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/iku"
)

type IkuService interface {
	FindAll(ctx context.Context, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]iku.IkuResponse, error)
	FindAllIkuOpd(ctx context.Context, kodeOpd string, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]iku.IkuOpdResponse, error)
	UpdateIkuActive(ctx context.Context, id string, request iku.IkuUpdateActiveRequest) error
	UpdateIkuOpdActive(ctx context.Context, id string, request iku.IkuUpdateActiveRequest) error
	FindAllIkuRenja(ctx context.Context, kodeOpd string, tahun string, jenisPeriode string, jenisIndikator string) ([]iku.IkuOpdResponse, error)
}
