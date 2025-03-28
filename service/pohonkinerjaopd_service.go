package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
)

type PohonKinerjaOpdService interface {
	Create(ctx context.Context, request pohonkinerja.PohonKinerjaCreateRequest) (pohonkinerja.PohonKinerjaOpdResponse, error)
	Update(ctx context.Context, request pohonkinerja.PohonKinerjaUpdateRequest) (pohonkinerja.PohonKinerjaOpdResponse, error)
	Delete(ctx context.Context, id string) error
	FindById(ctx context.Context, id int) (pohonkinerja.PohonKinerjaOpdResponse, error)
	FindAll(ctx context.Context, kodeOpd, tahun string) (pohonkinerja.PohonKinerjaOpdAllResponse, error)
	FindStrategicNoParent(ctx context.Context, kodeOpd, tahun string) ([]pohonkinerja.StrategicOpdResponse, error)
	DeletePelaksana(ctx context.Context, pelaksanaId string) error
	FindPokinByPelaksana(ctx context.Context, pegawaiId string, tahun string) ([]pohonkinerja.PohonKinerjaOpdResponse, error)
	DeletePokinPemdaInOpd(ctx context.Context, id int) error
	UpdateParent(ctx context.Context, pohonKinerja pohonkinerja.PohonKinerjaUpdateRequest) (pohonkinerja.PohonKinerjaOpdResponse, error)
	FindidPokinWithAllTema(ctx context.Context, id int) (pohonkinerja.PohonKinerjaAdminResponse, error)
}
