package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/isuglobal"
)

type IsuGlobalService interface {
	Create(ctx context.Context, request isuglobal.IsuGlobalRequest) (isuglobal.IsuGlobalResponse, error)
	Update(ctx context.Context, request isuglobal.IsuGlobalUpdateRequest) (isuglobal.IsuGlobalResponse, error)
	Delete(ctx context.Context, id int) error
}