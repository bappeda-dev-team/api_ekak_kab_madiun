package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/indikator"
)

type IkmService interface {
	FindAllByPeriode(ctx context.Context, tahunAwal, tahunAkhir string) ([]indikator.IkmResponse, error)
	FindById(ctx context.Context, ikmId string) (indikator.IkmResponse, error)
	Create(ctx context.Context, request indikator.IkmRequest) (indikator.IkmResponse, error)
	Update(ctx context.Context, request indikator.IkmRequest, ikmId string) (indikator.IkmResponse, error)
	Delete(ctx context.Context, ikmId string) error
}
