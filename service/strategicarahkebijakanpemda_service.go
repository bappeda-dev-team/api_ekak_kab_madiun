package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/strategicarahkebijakan"
)

type StrategicArahKebijakanPemdaService interface {
	FindAll(ctx context.Context, tahunAwal, tahunAkhir string) (strategicarahkebijakan.StrategicArahKebijakanPemdaAllResponse, error)
}