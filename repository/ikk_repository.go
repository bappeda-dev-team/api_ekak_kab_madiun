package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type IkkRepository interface {
	Create(ctx context.Context, tx *sql.Tx, ikk domain.Ikk) (domain.Ikk, error)
	Update(ctx context.Context, tx *sql.Tx, ikk domain.Ikk) (domain.Ikk, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.Ikk, error)
	FindByKodeOpd(ctx context.Context, tx *sql.Tx, jenis string, kodeOpd string) ([]domain.Ikk, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.Ikk, error)
	FindAllByIdPokin(ctx context.Context, tx *sql.Tx, pokinId int) ([]domain.Ikk, error)
	FindAllById(ctx context.Context, tx *sql.Tx, id int) (domain.Ikk, error)
	FindAllByJenisAndKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string, jenis string) ([]domain.Ikk, error)
	FindSelection(ctx context.Context, tx *sql.Tx) ([]domain.BidangUrusanSelection, error)
	FindSelectionByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.BidangUrusanSelection, error)
	PilihIkk(ctx context.Context, tx *sql.Tx, ikd domain.IkkTerpilih) (domain.IkkTerpilih, error)
	DeletePilihanIkk(ctx context.Context, tx *sql.Tx, id int) error
	FindTerpilihById(ctx context.Context, tx *sql.Tx, id int) (domain.IkkTerpilih, error)
	FindAllTerpilihByPokinId(ctx context.Context, tx *sql.Tx, id int) ([]domain.IkkTerpilih, error)
	FindTerpilihPokinIkkById(ctx context.Context, tx *sql.Tx, id int) (domain.IkkTerpilihDetail, error)
}