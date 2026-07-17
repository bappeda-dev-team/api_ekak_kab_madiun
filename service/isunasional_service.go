package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/isunasional"
)

type IsuNasionalService interface {
	Create(ctx context.Context, request isunasional.IsuNasionalRequest) (isunasional.IsuNasionalResponse, error)
	Update(ctx context.Context, request isunasional.IsuNasionalUpdateRequest) (isunasional.IsuNasionalResponse, error)
	Delete(ctx context.Context, id int) error
}