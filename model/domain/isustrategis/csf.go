package isustrategis

import (
	"ekak_kabupaten_madiun/model/domain"
	"time"
)

type CSF struct {
	ID                         int
	PohonID                    int
	PernyataanKondisiStrategis string
	AlasanKondisiStrategis     string
	DataTerukur                string
	KondisiTerukur             string
	KondisiWujud               string
	Tahun                      int
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}

type CSFPokin struct {
	ID                         int
	PohonID                    int
	PernyataanKondisiStrategis string
	AlasanKondisiStrategis     string
	DataTerukur                string
	KondisiTerukur             string
	KondisiWujud               string
	Tahun                      int
	JenisPohon                 string
	LevelPohon                 int
	Strategi                   string
	Keterangan                 string
	IsActive                   bool
	Indikator                  []domain.Indikator
}
