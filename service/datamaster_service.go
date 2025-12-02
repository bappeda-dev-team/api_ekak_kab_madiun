package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/datamaster"
)

type DataMasterService interface {
	DataRBByTahun(ctx context.Context, tahunBase int) ([]datamaster.RBResponse, error)
	SaveRB(ctx context.Context, req datamaster.RBRequest, userId int) (datamaster.RBResponse, error)
	UpdateRB(ctx context.Context, req datamaster.RBRequest, userId int, rbId int) (datamaster.RBResponse, error)
}
