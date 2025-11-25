package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
)

type CascadingOpdService interface {
	FindAll(ctx context.Context, kodeOpd, tahun string) (pohonkinerja.CascadingOpdResponse, error)
	FindByRekinPegawaiAndId(ctx context.Context, rekinId string) (pohonkinerja.CascadingRekinPegawaiResponse, error)
	FindByIdPokin(ctx context.Context, pokinId int) (pohonkinerja.CascadingRekinPegawaiResponse, error)
	FindByNip(ctx context.Context, nip string, tahun string) ([]pohonkinerja.CascadingRekinPegawaiResponse, error)
	FindByMultipleRekinPegawai(ctx context.Context, request pohonkinerja.FindByMultipleRekinRequest) ([]pohonkinerja.CascadingRekinPegawaiResponse, error)
	MultiRekinDetails(ctx context.Context, request pohonkinerja.FindByMultipleRekinRequest) ([]pohonkinerja.DetailRekinResponse, error)
	MultiRekinDetailsByOpdTahun(ctx context.Context, request pohonkinerja.MultiRekinDetailsByOpdAndTahunRequest) ([]pohonkinerja.DetailRekinResponse, error)
}
