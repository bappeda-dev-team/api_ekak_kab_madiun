package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
)

type MatrixRenjaService interface {
	GetByKodeOpdAndTahun(ctx context.Context, kodeOpd string, tahun string) ([]programkegiatan.UrusanDetailResponse, error)
}
