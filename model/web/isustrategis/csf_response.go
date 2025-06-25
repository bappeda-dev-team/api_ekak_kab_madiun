package isustrategis

import "time"

type CSFResponse struct {
	ID                         int       `json:"id"`
	PohonID                    int       `json:"pohon_id"`
	PernyataanKondisiStrategis string    `json:"pernyataan_kondisi_strategis"`
	AlasanKondisiStrategis     string    `json:"alasan_sebagai_kondisi_strategis"`
	DataTerukur                string    `json:"data_terukur"`
	KondisiTerukur             string    `json:"kondisi_terukur"`
	KondisiWujud               string    `json:"kondisi_wujud"`
	Tahun                      int       `json:"tahun"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}
