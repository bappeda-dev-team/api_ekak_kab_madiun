package service

import (
	"context"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/rencanakinerja"
)

type RencanaKinerjaService interface {
	Create(ctx context.Context, request rencanakinerja.RencanaKinerjaCreateRequest) (rencanakinerja.RencanaKinerjaResponse, error)
	Update(ctx context.Context, request rencanakinerja.RencanaKinerjaUpdateRequest) (rencanakinerja.RencanaKinerjaResponse, error)
	FindAll(ctx context.Context, pegawaiId string, kodeOPD string, tahun string) ([]rencanakinerja.RencanaKinerjaResponse, error)
	FindById(ctx context.Context, id string, kodeOPD string, tahun string) (rencanakinerja.RencanaKinerjaResponse, error)
	Delete(ctx context.Context, id string) error
	RekinsasaranOpd(ctx context.Context, pegawaiId string, kodeOPD string, tahun string) ([]rencanakinerja.RencanaKinerjaResponse, error)

	FindAllRincianKak(ctx context.Context, pegawaiId string, rencanaKinerjaId string) ([]rencanakinerja.DataRincianKerja, error)

	//rencana kinerja level 1
	CreateRekinLevel1(ctx context.Context, request rencanakinerja.RencanaKinerjaCreateRequest) (rencanakinerja.RencanaKinerjaResponse, error)
	UpdateRekinLevel1(ctx context.Context, request rencanakinerja.RencanaKinerjaUpdateRequest) (rencanakinerja.RencanaKinerjaResponse, error)
	FindIdRekinLevel1(ctx context.Context, id string) (rencanakinerja.RencanaKinerjaLevel1Response, error)

	//rencana kinerja level 3
	FindRekinLevel3(ctx context.Context, kodeOpd string, tahun string) ([]rencanakinerja.RencanaKinerjaResponse, error)

	//rencana kinerja atasan
	FindRekinAtasan(ctx context.Context, rekinId string) (rencanakinerja.RekinAtasanResponse, error)

	CloneRencanaKinerja(ctx context.Context, rekinId string, tahunBaru string) (rencanakinerja.RencanaKinerjaResponse, error)

	FindByFilter(ctx context.Context, filter domain.FilterParams) ([]rencanakinerja.RencanaKinerjaResponse, error)
}
