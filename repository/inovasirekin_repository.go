package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type InovasiRekinRepository interface {
	Create(ctx context.Context, tx *sql.Tx, inovasiRekin domain.InovasiRekin) (domain.InovasiRekin, error)
	Update(ctx context.Context, tx *sql.Tx, inovasiRekin domain.InovasiRekin) (domain.InovasiRekin, error)
	Delete(ctx context.Context, tx *sql.Tx, id string) error
	FindById(ctx context.Context, tx *sql.Tx, id string) (domain.InovasiRekin, error)
	FindAll(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.InovasiRekin, error)
}