package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
)

type PohonKinerjaAdminService interface {
	Create(ctx context.Context, request pohonkinerja.PohonKinerjaAdminCreateRequest) (pohonkinerja.PohonKinerjaAdminResponseData, error)
	Update(ctx context.Context, request pohonkinerja.PohonKinerjaAdminUpdateRequest) (pohonkinerja.PohonKinerjaAdminResponseData, error)
	Delete(ctx context.Context, id int) error
	FindById(ctx context.Context, id int) (pohonkinerja.PohonKinerjaAdminResponseData, error)
	FindAll(ctx context.Context, tahun string) (pohonkinerja.PohonKinerjaAdminResponse, error)
	FindSubTematik(ctx context.Context, tahun string) (pohonkinerja.PohonKinerjaAdminResponse, error)
	FindPokinAdminByIdHierarki(ctx context.Context, idPokin int) (pohonkinerja.TematikResponse, error)
	CreateStrategicAdmin(ctx context.Context, request pohonkinerja.PohonKinerjaAdminStrategicCreateRequest) (pohonkinerja.PohonKinerjaAdminResponseData, error)
	CloneStrategiFromPemda(ctx context.Context, request pohonkinerja.PohonKinerjaAdminStrategicCreateRequest) (pohonkinerja.PohonKinerjaAdminResponseData, error)
	TolakPokin(ctx context.Context, request pohonkinerja.PohonKinerjaAdminTolakRequest) error
	CrosscuttingOpd(ctx context.Context, request pohonkinerja.PohonKinerjaAdminStrategicCreateRequest) (pohonkinerja.PohonKinerjaAdminResponseData, error)
	TolakCrosscutting(ctx context.Context, request pohonkinerja.PohonKinerjaAdminTolakRequest) error
	SetujuiCrosscutting(ctx context.Context, request pohonkinerja.PohonKinerjaAdminTolakRequest) error
	AktiforNonAktifTematik(ctx context.Context, request pohonkinerja.TematikStatusRequest) (string, error)
	FindListOpdAllTematik(ctx context.Context, tahun string) ([]pohonkinerja.TematikListOpdResponse, error)

	//find pokin for dropdown
	FindPokinByTematik(ctx context.Context, tahun string) ([]pohonkinerja.PohonKinerjaAdminResponseData, error)
	FindPokinByStrategic(ctx context.Context, kodeOpd string, tahun string) ([]pohonkinerja.PohonKinerjaAdminResponseData, error)
	FindPokinByTactical(ctx context.Context, kodeOpd string, tahun string) ([]pohonkinerja.PohonKinerjaAdminResponseData, error)
	FindPokinByOperational(ctx context.Context, kodeOpd string, tahun string) ([]pohonkinerja.PohonKinerjaAdminResponseData, error)
	FindPokinByStatus(ctx context.Context, kodeOpd string, tahun string) ([]pohonkinerja.PohonKinerjaAdminResponseData, error)
	FindPokinByCrosscuttingStatus(ctx context.Context, kodeOpd string, tahun string) ([]pohonkinerja.PohonKinerjaAdminResponseData, error)
	FindPokinFromPemda(ctx context.Context, kodeOpd string, tahun string) ([]pohonkinerja.PohonKinerjaAdminResponseData, error)
	FindPokinFromOpd(ctx context.Context, kodeOpd string, tahun string, levelPohon int) ([]pohonkinerja.PohonKinerjaAdminResponseData, error)
}
