package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/usulan"
)

type UsulanPokokPikiranService interface {
	Create(ctx context.Context, request usulan.UsulanPokokPikiranCreateRequest) (usulan.UsulanPokokPikiranResponse, error)
	Update(ctx context.Context, request usulan.UsulanPokokPikiranUpdateRequest) (usulan.UsulanPokokPikiranResponse, error)
	FindById(ctx context.Context, idUsulan string) (usulan.UsulanPokokPikiranResponse, error)
	FindAll(ctx context.Context, pegawaiId *string, isActive *bool, rekinId *string) ([]usulan.UsulanPokokPikiranResponse, error)
	Delete(ctx context.Context, idUsulan string) error
}