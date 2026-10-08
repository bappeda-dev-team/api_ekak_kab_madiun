package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/outcomematrix"
)

type OutcomeMatrixService interface {
	Create(ctx context.Context, request outcomematrix.OutcomeMatrixCreateRequest) (outcomematrix.OutcomeMatrixResponse, error)
	Update(ctx context.Context, request outcomematrix.OutcomeMatrixUpdateRequest) (outcomematrix.OutcomeMatrixResponse, error)
	Delete(ctx context.Context, id int) error
	FindById(ctx context.Context, id int) (outcomematrix.OutcomeMatrixResponse, error)
	FindAll(ctx context.Context, kode, kodeOpd, jenis string) ([]outcomematrix.OutcomeMatrixResponse, error)
	UpsertBatch(ctx context.Context, requests []outcomematrix.OutcomeMatrixBatchItemRequest) ([]outcomematrix.OutcomeMatrixResponse, error)
}
