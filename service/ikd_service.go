package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/ikd"
)

type IkdService interface {
	FindAll(ctx context.Context, kodeOpd string, tahun string, jenisPeriode string) ([]ikd.IkdResponse, error)
}