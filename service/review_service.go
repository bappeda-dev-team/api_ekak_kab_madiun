package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
)

type ReviewService interface {
	Create(ctx context.Context, request pohonkinerja.ReviewCreateRequest) (pohonkinerja.ReviewResponse, error)
	Update(ctx context.Context, request pohonkinerja.ReviewUpdateRequest) (pohonkinerja.ReviewResponse, error)
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context, idPohonKinerja int) ([]pohonkinerja.ReviewResponse, error)
	FindById(ctx context.Context, id int) (pohonkinerja.ReviewResponse, error)
	FindAllReviewByTematik(ctx context.Context, tahun string) ([]pohonkinerja.ReviewTematikResponse, error)
	FindAllReviewOpd(ctx context.Context, kodeOpd, tahun string) ([]pohonkinerja.ReviewOpdResponse, error)
}
