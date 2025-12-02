package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/datamaster"
)

type DataMasterRepository interface {
	DataRBByTahun(ctx context.Context, tx *sql.Tx) ([]datamaster.MasterRB, error)
}
