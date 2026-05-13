package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type IkdRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string, jenisPeriode string) ([]domain.IkdDetail, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.ProgramOpdTerpilih, error)
	Create(ctx context.Context, tx *sql.Tx, ikd domain.ProgramOpdTerpilih) (domain.ProgramOpdTerpilih, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
}