package subkegiatan

type SubKegiatanTerpilihResponse struct {
	Id              string              `json:"id,omitempty"`
	KodeSubKegiatan SubKegiatanResponse `json:"kode_subkegiatan"`
}

type SubKegiatanOpdResponse struct {
	Id              int    `json:"id"`
	KodeSubkegiatan string `json:"kode_subkegiatan"`
	NamaSubkegiatan string `json:"nama_subkegiatan"`
	KodeOpd         string `json:"kode_opd"`
	NamaOpd         string `json:"nama_opd"`
	Tahun           string `json:"tahun"`
}
