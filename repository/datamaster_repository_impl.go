package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/datamaster"
	"fmt"
)

type DataMasterRepositoryImpl struct {
}

func NewDataMasterRepositoryImpl() *DataMasterRepositoryImpl {
	return &DataMasterRepositoryImpl{}
}

func (repository *DataMasterRepositoryImpl) DataRBByTahun(ctx context.Context, tx *sql.Tx) ([]datamaster.MasterRB, error) {
	return nil, fmt.Errorf("No data yet")
}
