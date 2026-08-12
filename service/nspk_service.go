package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/masternspk"
)

type NspkService interface {
	Create(ctx context.Context, request masternspk.NspkRequest) (masternspk.NspkResponse, error)
	Update(ctx context.Context, request masternspk.NspkUpdateRequest) (masternspk.NspkResponse, error)
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context, kodeOpd string) ([]masternspk.NspkFullResponse, error)
}