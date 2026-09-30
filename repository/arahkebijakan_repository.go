package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type ArahKebijakanRepository interface {
	Create(ctx context.Context, tx *sql.Tx, ar domain.ArahKebijakanOpd) (domain.ArahKebijakanOpd, error)
	Update(ctx context.Context, tx *sql.Tx, ar domain.ArahKebijakanOpd) (domain.ArahKebijakanOpd, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.ArahKebijakanOpd, error)
}