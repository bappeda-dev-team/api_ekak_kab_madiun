package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type MisiPemdaRepository interface {
	Create(ctx context.Context, tx *sql.Tx, visiPemda domain.MisiPemda) (domain.MisiPemda, error)
	Update(ctx context.Context, tx *sql.Tx, visiPemda domain.MisiPemda) (domain.MisiPemda, error)
	Delete(ctx context.Context, tx *sql.Tx, visiPemdaId int) error
	FindById(ctx context.Context, tx *sql.Tx, visiPemdaId int) (domain.MisiPemda, error)
	FindAll(ctx context.Context, tx *sql.Tx, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]domain.MisiPemda, error)
	FindByIdWithDefault(ctx context.Context, tx *sql.Tx, visiPemdaId int) (domain.MisiPemda, error)
	CheckUrutanExists(ctx context.Context, tx *sql.Tx, idVisi int, urutan int) (bool, error)
	CheckUrutanExistsExcept(ctx context.Context, tx *sql.Tx, idVisi int, urutan int, id int) (bool, error)
	FindByIdVisi(ctx context.Context, tx *sql.Tx, idVisi int) ([]domain.MisiPemda, error)
	FindByTematikId(ctx context.Context, tx *sql.Tx, tematikId int) ([]domain.MisiPemda, error)
}
