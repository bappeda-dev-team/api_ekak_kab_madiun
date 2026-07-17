package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/isuklhs"
)

type IsuKlhsService interface {
	Create(ctx context.Context, request isuklhs.IsuKlhsRequest) (isuklhs.IsuKlhsResponse, error)
	Update(ctx context.Context, request isuklhs.IsuKlhsUpdateRequest) (isuklhs.IsuKlhsResponse, error)
	Delete(ctx context.Context, id int) error
}