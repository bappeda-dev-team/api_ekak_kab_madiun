package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type TaggingPokinRepository interface {
	Create(ctx context.Context, tx *sql.Tx, tagging domain.TaggingPokin) (domain.TaggingPokin, error)
	Update(ctx context.Context, tx *sql.Tx, tagging domain.TaggingPokin) (domain.TaggingPokin, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.TaggingPokin, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]domain.TaggingPokin, error)
}
