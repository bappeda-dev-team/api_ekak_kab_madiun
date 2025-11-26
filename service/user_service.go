package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/user"
)

type UserService interface {
	Create(ctx context.Context, request user.UserCreateRequest) (user.UserResponse, error)
	Update(ctx context.Context, request user.UserUpdateRequest) (user.UserResponse, error)
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context, kodeOpd string) ([]user.UserResponse, error)
	FindById(ctx context.Context, id int) (user.UserResponse, error)
	Login(ctx context.Context, request user.UserLoginRequest) (user.UserLoginResponse, error)
	FindByKodeOpdAndRole(ctx context.Context, kodeOpd string, roleName string) ([]user.UserResponse, error)
	FindByNip(ctx context.Context, nip string) (user.UserResponse, error)
	CekAdminOpd(ctx context.Context) ([]user.CekAdminOpdResponse, error)
}
