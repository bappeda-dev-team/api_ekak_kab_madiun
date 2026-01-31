package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/pegawai"
)

type PegawaiService interface {
	Create(ctx context.Context, request pegawai.PegawaiCreateRequest) (pegawai.PegawaiResponse, error)
	Update(ctx context.Context, request pegawai.PegawaiUpdateRequest) (pegawai.PegawaiResponse, error)
	Delete(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (pegawai.PegawaiResponse, error)
	FindAll(ctx context.Context, kodeOpd string) ([]pegawai.PegawaiResponse, error)
	TambahJabatan(ctx context.Context, request pegawai.TambahJabatanRequest) (pegawai.PegawaiResponse, error)
}
