package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/isustrategis"
)

type CSFRepository interface {
	FindByTahun(ctx context.Context, tx *sql.Tx, tahun string) ([]isustrategis.CSF, error)
}
