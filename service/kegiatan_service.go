package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/kegiatan"
)

type KegiatanService interface {
	Create(ctx context.Context, request kegiatan.KegiatanCreateRequest) (kegiatan.KegiatanResponse, error)
	Update(ctx context.Context, request kegiatan.KegiatanUpdateRequest) (kegiatan.KegiatanResponse, error)
	Delete(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (kegiatan.KegiatanResponse, error)
	FindAll(ctx context.Context) ([]kegiatan.KegiatanResponse, error)
	FindAllWithIndikator(ctx context.Context) ([]kegiatan.KegiatanWithIndikatorV2Response, error)
}
