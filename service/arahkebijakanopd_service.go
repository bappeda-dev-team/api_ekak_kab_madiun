package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/arahkebijakan"
)

type ArahKebijakanService interface {
	Create(ctx context.Context, request arahkebijakan.ArahKebijakanRequest) (arahkebijakan.ArahKebijakanResponse, error)
	Update(ctx context.Context, request arahkebijakan.ArahKebijakanUpdateRequest) (arahkebijakan.ArahKebijakanResponse, error)
}