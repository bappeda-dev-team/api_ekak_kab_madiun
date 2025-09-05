package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/taggingpokin"
)

type TaggingPokinService interface {
	Create(ctx context.Context, request taggingpokin.TaggingPokinCreateRequest) (taggingpokin.TaggingPokinResponse, error)
	Update(ctx context.Context, request taggingpokin.TaggingPokinUpdateRequest) (taggingpokin.TaggingPokinResponse, error)
	Delete(ctx context.Context, id int) error
	FindById(ctx context.Context, id int) (taggingpokin.TaggingPokinResponse, error)
	FindAll(ctx context.Context) ([]taggingpokin.TaggingPokinResponse, error)
}
