package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/inovasirekin"
)

type InovasiRekinService interface {
	Create(ctx context.Context, request inovasirekin.InovasiRekinCreateRequest) (inovasirekin.InovasiRekinResponse, error)
	Update(ctx context.Context, request inovasirekin.InovasiRekinUpdateRequest) (inovasirekin.InovasiRekinResponse, error)
	FindById(ctx context.Context, id string) (inovasirekin.InovasiRekinResponse, error)
	FindAll(ctx context.Context, rekinId string) ([]inovasirekin.InovasiRekinResponse, error)
	Delete(ctx context.Context, id string) error
}