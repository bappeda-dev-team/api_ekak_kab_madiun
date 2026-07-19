package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/isuregional"
)

type IsuRegionalService interface {
	Create(ctx context.Context, request isuregional.IsuRegionalRequest) (isuregional.IsuRegionalResponse, error)
	Update(ctx context.Context, request isuregional.IsuRegionalUpdateRequest) (isuregional.IsuRegionalResponse, error)
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context, kodeOpd string) (isuregional.IsuRegionalMasterResponse, error)
}