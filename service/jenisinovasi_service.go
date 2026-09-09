package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/jenisinovasi"
)

type JenisInovasiService interface {
	Create(ctx context.Context, request jenisinovasi.JenisInovasiRequest) (jenisinovasi.JenisInovasiResponse, error)
	Update(ctx context.Context, request jenisinovasi.JenisInovasiUpdateRequest) (jenisinovasi.JenisInovasiResponse, error)
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context) ([]jenisinovasi.JenisInovasiResponse, error)
}