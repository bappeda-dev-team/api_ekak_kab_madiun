package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
)

type CascadingOpdService interface {
	FindAll(ctx context.Context, kodeOpd, tahun string) (pohonkinerja.CascadingOpdResponse, error)
}
