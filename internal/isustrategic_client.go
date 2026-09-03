package internal

import "context"

type IsustrategicClient interface {
	GetDataIsuStrategic(ctx context.Context, kodeOpd string, tahun string) ([]IsuStrategisResponse, error)
	GetDataPermasalahan(ctx context.Context, kodeOpd string, tahun string) ([]PermasalahanResp, error)
}
