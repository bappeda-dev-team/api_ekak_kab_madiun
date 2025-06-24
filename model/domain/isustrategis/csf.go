package isustrategis

import "time"

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
