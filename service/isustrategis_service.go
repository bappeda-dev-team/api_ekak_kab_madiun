package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/isustrategis"
)

type CSFService interface {
	AllCsfsByTahun(ctx context.Context, tahun string) ([]isustrategis.CSFResponse, error)
	FindByTahun(ctx context.Context, tahun string) ([]isustrategis.CSFResponse, error)
	FindById(ctx context.Context, csfID int) (isustrategis.CSFResponse, error)
}
