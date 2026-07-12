package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/sasaranpemda"
)

type SasaranPemdaService interface {
	Create(ctx context.Context, request sasaranpemda.SasaranPemdaCreateRequest) (sasaranpemda.SasaranPemdaResponse, error)
	Update(ctx context.Context, request sasaranpemda.SasaranPemdaUpdateRequest) (sasaranpemda.SasaranPemdaResponse, error)
	Delete(ctx context.Context, id int) error
	FindById(ctx context.Context, sasaranPemdaId int) (sasaranpemda.SasaranPemdaResponse, error)
	FindAll(ctx context.Context, tahun string) ([]sasaranpemda.SasaranPemdaResponse, error)
	FindAllWithPokin(ctx context.Context, tahunAwal, tahunAkhir, jenisPeriode string) ([]sasaranpemda.TematikResponse, error)
	FindSasaranPemdaRanwal(ctx context.Context, tahun, jenisPeriode string) ([]sasaranpemda.SasaranPemdaResponse, error)
	// Dual target
	FindSasaranPemdaRankhirDual(ctx context.Context, tahun, jenisPeriode string) ([]sasaranpemda.SasaranPemdaRankhirDualResponse, error)
	FindSasaranPemdaPenetapanDual(ctx context.Context, tahun, jenisPeriode string) ([]sasaranpemda.SasaranPemdaPenetapanDualResponse, error)
	// Target layer
	CreateTargetSasaranLayer(ctx context.Context, jenis string, req sasaranpemda.LayerTargetBatchRequest) ([]sasaranpemda.TargetResponse, error)
	UpdateTargetSasaranLayer(ctx context.Context, jenis string, req sasaranpemda.LayerTargetUpdateBatchRequest) ([]sasaranpemda.TargetResponse, error)
	LockSasaranPemda(ctx context.Context, tahun string) (sasaranpemda.LockDataPemdaResponse, error)
	UnlockSasaranPemda(ctx context.Context, tahun string) (sasaranpemda.LockDataPemdaResponse, error)
	IsSasaranPemdaLocked(ctx context.Context, tahun string) (sasaranpemda.LockDataPemdaResponse, error)
	FindAllLockSasaranPemda(ctx context.Context) ([]sasaranpemda.LockDataPemdaResponse, error)
}
