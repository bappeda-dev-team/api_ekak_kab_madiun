package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/datamaster"
)

type DataMasterService interface {
	DataRBByTahun(ctx context.Context, tahunBase int) (datamaster.RBResponse, error)
}
