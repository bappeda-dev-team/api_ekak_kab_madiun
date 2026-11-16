package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/rencanakinerja"
)

type ManualIKService interface {
	Create(ctx context.Context, request rencanakinerja.ManualIKCreateRequest, indikatorId string) (rencanakinerja.ManualIKResponse, error)
	Update(ctx context.Context, request rencanakinerja.ManualIKUpdateRequest, indikatorId string) (rencanakinerja.ManualIKResponse, error)
	FindManualIKByIndikatorId(ctx context.Context, indikatorId string) (rencanakinerja.ManualIKResponse, error)
	FindManualIKSasaranOpdByIndikatorId(ctx context.Context, indikatorId string, tahun string) (rencanakinerja.ManualIKResponse, error)
	FindManualIKByPegawai(ctx context.Context, nip string, kodeOpd string, tahun string) ([]rencanakinerja.ManualIKResponse, error)
	FindManualIKByRekin(ctx context.Context, rekinId string) ([]rencanakinerja.ManualIKResponse, error)
}
