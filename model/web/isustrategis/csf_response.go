package isustrategis

import "ekak_kabupaten_madiun/model/web/pohonkinerja"

type CSFResponse struct {
	ID                         int                              `json:"id"`
	PohonID                    int                              `json:"pohon_id"`
	PernyataanKondisiStrategis string                           `json:"pernyataan_kondisi_strategis"`
	AlasanKondisiStrategis     string                           `json:"alasan_sebagai_kondisi_strategis"`
	DataTerukur                string                           `json:"data_terukur_pendukung_pernyataan"`
	KondisiTerukur             string                           `json:"kondisi_terukur_yang_diharapkan"`
	KondisiWujud               string                           `json:"kondisi_yang_ingin_diwujudkan"`
	Tahun                      int                              `json:"tahun"`
	JenisPohon                 string                           `json:"jenis_pohon"`
	LevelPohon                 int                              `json:"level_pohon"`
	Strategi                   string                           `json:"tema"`
	NamaPohon                  string                           `json:"nama_pohon"`
	Keterangan                 string                           `json:"keterangan"`
	IsActive                   bool                             `json:"is_active"`
	Indikators                 []pohonkinerja.IndikatorResponse `json:"indikator,omitempty"`
}
