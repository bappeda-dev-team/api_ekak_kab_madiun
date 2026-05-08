package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/ikk"
)

type IkkService interface {
	Create(ctx context.Context, request ikk.IkkRequest) (ikk.IkkResponse, error)
	Update(ctx context.Context, request ikk.IkkUpdateRequest) (ikk.IkkResponse, error)
	Delete(ctx context.Context, id int) error
	FindById(ctx context.Context, id int) (ikk.IkkResponse, error)
	FindByKodeOpd(ctx context.Context, levelPohon int, kodeOpd string) ([]ikk.IkkFullResponse, error)
}