package pohonkinerja

type CSFAPIResponse struct {
	Data []CSFApiResponse `json:"data"`
}

type CSFApiResponse struct {
	Id                         int                     `json:"id"`
	PohonId                    int                     `json:"pohon_id"`
	PernyataanKondisiStrategis string                  `json:"pernyataan_kondisi_strategis"`
	AlasanKondisi              []AlasanKondisiResponse `json:"alasan_kondisi"`
}

type AlasanKondisiResponse struct {
	Id                             int                   `json:"id"`
	CSFid                          int                   `json:"csf_id"`
	AlasanKondisiStrategis         string                `json:"alasan_kondisi_strategis"`
	DataTerukurPendukungPernyataan []DataTerukurResponse `json:"data_terukur"`
}

type DataTerukurResponse struct {
	Id              int    `json:"id"`
	AlasanKondisiId int    `json:"alasan_kondisi_id"`
	DataTerukur     string `json:"data_terukur"`
}
