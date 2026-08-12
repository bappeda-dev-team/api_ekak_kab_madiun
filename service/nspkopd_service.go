package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/nspkopd"
)

type NspkOpdService interface {
	Create(ctx context.Context, request nspkopd.NspkRequest) (nspkopd.NspkResponse, error)
	Update(ctx context.Context, request nspkopd.NspkUpdateRequest) (nspkopd.NspkResponse, error)
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context, kodeOpd string) ([]nspkopd.NspkFullResponse, error)
}