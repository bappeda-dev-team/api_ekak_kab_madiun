package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/rincianbelanja"
)

type RincianBelanjaService interface {
	Create(ctx context.Context, rincianBelanja rincianbelanja.RincianBelanjaCreateRequest) (rincianbelanja.RencanaAksiResponse, error)
	Update(ctx context.Context, rincianBelanja rincianbelanja.RincianBelanjaUpdateRequest) (rincianbelanja.RencanaAksiResponse, error)
	FindRincianBelanjaAsn(ctx context.Context, pegawaiId string, tahun string) []rincianbelanja.RincianBelanjaAsnResponse
	LaporanRincianBelanjaOpd(ctx context.Context, kodeOpd string, tahun string) ([]rincianbelanja.RincianBelanjaAsnResponse, error)
	LaporanRincianBelanjaPegawai(ctx context.Context, pegawaiId string, tahun string) ([]rincianbelanja.RincianBelanjaAsnResponse, error)
}
