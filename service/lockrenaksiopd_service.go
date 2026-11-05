package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/lockrenaksiopd"
)

type LockRenaksiOpdService interface {
	Lock(ctx context.Context, kodeOpd, tahun string, request lockrenaksiopd.LockRenaksiOpdRequest) (lockrenaksiopd.LockRenaksiOpdResponse, error)
	Unlock(ctx context.Context, kodeOpd, tahun string, id int) error
	FindById(ctx context.Context, kodeOpd, tahun string, id int) (lockrenaksiopd.LockRenaksiOpdResponse, error)
	FindAll(ctx context.Context, kodeOpd, tahun string) ([]lockrenaksiopd.LockRenaksiOpdResponse, error)
}
