package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/iku"
)

type IkuService interface {
	FindAll(ctx context.Context, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]iku.IkuResponse, error)
	FindAllIkuOpd(ctx context.Context, kodeOpd string, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]iku.IkuOpdResponse, error)
	UpdateIkuActive(ctx context.Context, id string, request iku.IkuUpdateActiveRequest) error
	UpdateIkuOpdActive(ctx context.Context, id string, request iku.IkuUpdateActiveRequest) error
	FindAllIkuRenja(ctx context.Context, kodeOpd string, tahun string, jenisPeriode string, jenisIndikator string) ([]iku.IkuOpdResponse, error)

	// v2 — IKU pemda filter berdasarkan tematik.tahun
	// ranwal: target renstra sebagai base
	FindIkuPemdaRanwalV2(ctx context.Context, tahun, jenisPeriode string) ([]iku.IkuResponse, error)
	// rankhir: dual target (target_ranwal=renstra base + target_rankhir=override)
	FindIkuPemdaRankhirDualV2(ctx context.Context, tahun, jenisPeriode string) ([]iku.IkuPemdaRankhirDualResponse, error)
	// penetapan: dual target (target_rankhir=base + target_penetapan=override)
	FindIkuPemdaPenetapanDualV2(ctx context.Context, tahun, jenisPeriode string) ([]iku.IkuPemdaPenetapanDualResponse, error)
}
