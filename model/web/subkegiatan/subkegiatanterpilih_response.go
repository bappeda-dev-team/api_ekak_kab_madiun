package subkegiatan

type SubKegiatanTerpilihResponse struct {
	Id              string              `json:"id,omitempty"`
	KodeSubKegiatan SubKegiatanResponse `json:"kode_subkegiatan"`
}

type SubKegiatanOpdMultipleResponse struct {
	SuccessCount   int                      `json:"success_count"`
	TotalRequested int                      `json:"total_requested"`
	SkippedCount   int                      `json:"skipped_count"`
	SuccessItems   []SubKegiatanOpdResponse `json:"success_items"`
	SkippedItems   []SubKegiatanOpdResponse `json:"skipped_items"`
	Message        string                   `json:"message"`
}

type SubKegiatanOpdResponse struct {
	Id              int    `json:"id"`
	KodeSubkegiatan string `json:"kode_subkegiatan"`
	NamaSubkegiatan string `json:"nama_subkegiatan"`
	KodeOpd         string `json:"kode_opd"`
	NamaOpd         string `json:"nama_opd"`
	Tahun           string `json:"tahun"`
	Status          string `json:"status"`
}
