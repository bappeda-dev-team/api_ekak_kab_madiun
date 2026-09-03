package service

import (
	"bytes"
	"context"
	"ekak_kabupaten_madiun/model/web/strategicarahkebijakan"
)

type StrategicArahKebijakanPemdaService interface {
	FindAll(ctx context.Context, tahunAwal string, tahunAkhir string) ([]strategicarahkebijakan.StrategiArahKebijakanPemdaResponse, error)
	ExportExcel(ctx context.Context, tahunAwal string, tahunAkhir string) (*bytes.Buffer, error)
}