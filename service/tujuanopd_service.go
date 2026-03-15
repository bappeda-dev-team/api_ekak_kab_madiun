package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/tujuanopd"
)

type TujuanOpdService interface {
	Create(ctx context.Context, request tujuanopd.TujuanOpdCreateRequest) (tujuanopd.TujuanOpdResponse, error)
	Update(ctx context.Context, request tujuanopd.TujuanOpdUpdateRequest) (tujuanopd.TujuanOpdResponse, error)
	Delete(ctx context.Context, tujuanOpdId int) error
	FindById(ctx context.Context, tujuanOpdId int) (tujuanopd.TujuanOpdResponse, error)
	FindAll(ctx context.Context, kodeOpd string, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]tujuanopd.TujuanOpdwithBidangUrusanResponse, error)
	FindTujuanOpdOnlyName(ctx context.Context, kodeOpd string, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]tujuanopd.TujuanOpdResponse, error)
	FindTujuanOpdByTahun(ctx context.Context, kodeOpd string, tahun string, jenisPeriode string) ([]tujuanopd.TujuanOpdwithBidangUrusanResponse, error)
	FindTujuanRenstra(ctx context.Context, kodeOpd, tahunAwal, tahunAkhir, jenisPeriode string) ([]tujuanopd.TujuanOpdwithBidangUrusanResponse, error)
	FindTujuanRanwal(ctx context.Context, kodeOpd, tahun, jenisPeriode string) ([]tujuanopd.TujuanOpdwithBidangUrusanResponse, error)
	FindTujuanRankhir(ctx context.Context, kodeOpd, tahun, jenisPeriode string) ([]tujuanopd.TujuanOpdwithBidangUrusanResponse, error)
	CreateTujuanRenjaIndikator(ctx context.Context, tujuanOpdId int, jenis string, requests []tujuanopd.IndikatorCreateRequest) ([]tujuanopd.IndikatorResponse, error)
	UpdateTujuanRenjaIndikator(ctx context.Context, kodeIndikator string, jenis string, request tujuanopd.IndikatorUpdateRequest) (tujuanopd.IndikatorResponse, error)
	DeleteTujuanRenjaIndikator(ctx context.Context, kodeIndikator string) error
}
