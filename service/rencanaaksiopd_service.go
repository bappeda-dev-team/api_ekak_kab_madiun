package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/renaksiopd"
)

type RencanaAksiOpdService interface {
	FindBySasaranOpdAndTahun(ctx context.Context, sasaranOpdId int, tahun string) ([]renaksiopd.RencanaAksiOpdResponse, error)
	SyncJadwalPelaksanaan(ctx context.Context, rekinId string) error
	Create(ctx context.Context, request renaksiopd.RencanaAksiOpdCreateRequest) (renaksiopd.RencanaAksiOpdRequestResponse, error)
	Update(ctx context.Context, request renaksiopd.RencanaAksiOpdUpdateRequest) (renaksiopd.RencanaAksiOpdRequestResponse, error)
	Delete(ctx context.Context, id int) error
	FindById(ctx context.Context, id int) (renaksiopd.RencanaAksiOpdByIdResponse, error)
	FindAllSasaranByTahun(ctx context.Context, kodeOpd string, tahun string) ([]renaksiopd.SasaranOpdDetailResponse, error)
}
