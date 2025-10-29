package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/subkegiatan"
)

type SubKegiatanTerpilihService interface {
	Update(ctx context.Context, request subkegiatan.SubKegiatanTerpilihUpdateRequest) (subkegiatan.SubKegiatanTerpilihResponse, error)
	FindByKodeSubKegiatan(ctx context.Context, kodeSubKegiatan string) (subkegiatan.SubKegiatanTerpilihResponse, error)
	Delete(ctx context.Context, id string, kodeSubKegiatan string) error
	CreateRekin(ctx context.Context, request subkegiatan.SubKegiatanCreateRekinRequest) ([]subkegiatan.SubKegiatanResponse, error)
	DeleteSubKegiatanTerpilih(ctx context.Context, idSubKegiatan string) error
	// CreateOpd(ctx context.Context, request subkegiatan.SubKegiatanOpdCreateRequest) (subkegiatan.SubKegiatanOpdResponse, error)
	UpdateOpd(ctx context.Context, request subkegiatan.SubKegiatanOpdUpdateRequest) (subkegiatan.SubKegiatanOpdResponse, error)
	FindAllOpd(ctx context.Context, kodeOpd, tahun *string) ([]subkegiatan.SubKegiatanOpdResponse, error)
	FindById(ctx context.Context, id int) (subkegiatan.SubKegiatanOpdResponse, error)
	DeleteOpd(ctx context.Context, id int) error
	FindAllSubkegiatanByBidangUrusanOpd(ctx context.Context, kodeOpd string) ([]subkegiatan.SubKegiatanResponse, error)
	CreateOpdMultiple(ctx context.Context, request subkegiatan.SubKegiatanOpdMultipleCreateRequest) (subkegiatan.SubKegiatanOpdMultipleResponse, error)
}
