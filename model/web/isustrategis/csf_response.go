package isustrategis

type CSFResponse struct {
	ID                         int       `json:"id"`
	PohonID                    int       `json:"pohon_id"`
	PernyataanKondisiStrategis string    `json:"pernyataan_kondisi_strategis"`
	AlasanKondisiStrategis     string    `json:"alasan_sebagai_kondisi_strategis"`
	DataTerukur                string    `json:"data_terukur"`
	KondisiTerukur             string    `json:"kondisi_terukur"`
	KondisiWujud               string    `json:"kondisi_wujud"`
	Tahun                      int       `json:"tahun"`
	JenisPohon                 string    `json:"jenis_pohon"`
	LevelPohon                 int       `json:"level_pohon"`
	Strategi                   string    `json:"tema"`
	Keterangan                 string    `json:"keterangan"`
	IsActive                   bool      `json:"is_active"`
}
