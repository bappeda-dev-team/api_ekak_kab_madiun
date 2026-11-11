package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/rencanaaksi"
)

type RencanaAksiService interface {
	Create(ctx context.Context, request rencanaaksi.RencanaAksiCreateRequest) (rencanaaksi.RencanaAksiResponse, error)
	Update(ctx context.Context, request rencanaaksi.RencanaAksiUpdateRequest) (rencanaaksi.RencanaAksiResponse, error)
	FindAll(ctx context.Context, rencanaKinerjaId string) ([]rencanaaksi.RencanaAksiResponse, error)
	FindById(ctx context.Context, id string) (rencanaaksi.RencanaAksiResponse, error)
	Delete(ctx context.Context, id string) error
	FindByPegawai(ctx context.Context, nip string, kodeOpd string, tahun string) ([]rencanaaksi.RencanaAksiPegawaiResponse, error)
	FindByRekin(ctx context.Context, rekinId string) ([]rencanaaksi.RencanaAksiResponse, error)
}
