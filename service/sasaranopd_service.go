package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/model/web/sasaranopd"
)

type SasaranOpdService interface {
	FindAll(ctx context.Context, KodeOpd string, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]sasaranopd.SasaranOpdResponse, error)
	FindById(ctx context.Context, id int) (*sasaranopd.SasaranOpdResponse, error)
	Create(ctx context.Context, request sasaranopd.SasaranOpdCreateRequest) (*sasaranopd.SasaranOpdCreateResponse, error)
	Update(ctx context.Context, request sasaranopd.SasaranOpdUpdateRequest) (*sasaranopd.SasaranOpdCreateResponse, error)
	Delete(ctx context.Context, id string) error
	FindByIdPokin(ctx context.Context, idPokin int, tahun string) (*sasaranopd.SasaranOpdResponse, error)
	FindIdPokinSasaran(ctx context.Context, id int) (pohonkinerja.PohonKinerjaOpdResponse, error)
	FindByTahun(ctx context.Context, kodeOpd string, tahun string, jenisPeriode string) ([]sasaranopd.SasaranOpdResponse, error)
	FindSasaranRenstra(ctx context.Context, kodeOpd, tahunAwal, tahunAkhir, jenisPeriode string) ([]sasaranopd.SasaranOpdResponse, error)
	FindSasaranRanwal(ctx context.Context, kodeOpd, tahun, jenisPeriode string) ([]sasaranopd.SasaranOpdResponse, error)
	FindSasaranRankhir(ctx context.Context, kodeOpd, tahun, jenisPeriode string) ([]sasaranopd.SasaranOpdResponse, error)
	CreateRenjaIndikator(ctx context.Context, sasaranOpdId int, jenis string, requests []sasaranopd.IndikatorCreateRequest) ([]sasaranopd.IndikatorResponse, error)
	UpdateRenjaIndikator(ctx context.Context, kodeIndikator string, jenis string, request sasaranopd.IndikatorUpdateRequest) (sasaranopd.IndikatorResponse, error)
	DeleteRenjaIndikator(ctx context.Context, kodeIndikator string) error
}
