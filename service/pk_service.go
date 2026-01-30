package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/pkopd"
)

type PkService interface {
	FindByKodeOpdTahun(ctx context.Context, kodeOpd string, tahun int) (pkopd.PkOpdResponse, error)
	HubungkanRekin(ctx context.Context, request pkopd.PkOpdRequest) (pkopd.PkOpdResponse, error)
	HubungkanAtasan(ctx context.Context, request pkopd.HubungkanAtasanRequest) (pkopd.PkOpdResponse, error)
}
