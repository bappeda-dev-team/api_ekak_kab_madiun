package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/datamaster"
)

type DataMasterRepository interface {
	DataRBByTahun(ctx context.Context, tx *sql.Tx, tahun int) ([]datamaster.MasterRB, error)
}
