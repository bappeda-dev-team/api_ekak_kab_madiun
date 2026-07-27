package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/ppd"
)

type PpdService interface {
	Create(ctx context.Context, request ppd.PpdRequest) (ppd.PpdResponse, error)
	Update(ctx context.Context, request ppd.PpdUpdateRequest) (ppd.PpdResponse, error)
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context, kodeOpd string) (ppd.PpdMasterResponse, error)
}