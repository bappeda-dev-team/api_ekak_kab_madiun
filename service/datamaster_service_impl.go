package service

import (
	"context"
	"fmt"

	"database/sql"
	"ekak_kabupaten_madiun/model/web/datamaster"
	"ekak_kabupaten_madiun/repository"
)

type DataMasterServiceImpl struct {
	DataMasterRepository repository.DataMasterRepository
	DB                   *sql.DB
}

func NewDataMasterServiceImpl(dataMasterRepository repository.DataMasterRepository, DB *sql.DB) *DataMasterServiceImpl {
	return &DataMasterServiceImpl{
		DataMasterRepository: dataMasterRepository,
		DB:                   DB,
	}
}

func (service *DataMasterServiceImpl) DataRBByTahun(ctx context.Context, tahunBase int) (datamaster.RBResponse, error) {
	return datamaster.RBResponse{}, fmt.Errorf("Belum implmentasi :x")
}
