package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type ReviewRepository interface {
	Create(ctx context.Context, tx *sql.Tx, review domain.Review) (domain.Review, error)
	Update(ctx context.Context, tx *sql.Tx, review domain.Review) (domain.Review, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.Review, error)
	FindByPohonKinerja(ctx context.Context, tx *sql.Tx, idPohonKinerja int) ([]domain.Review, error)
	CountReviewByPohonKinerja(ctx context.Context, tx *sql.Tx, idPohonKinerja int) (int, error)
	FindAllReviewByTematik(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.ReviewTematik, error)
}
