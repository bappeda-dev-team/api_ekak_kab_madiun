package pohonkinerja

type ReviewUpdateRequest struct {
	Id             int    `json:"id"`
	IdPohonKinerja int    `json:"id_pohon_kinerja"`
	Review         string `json:"review"`
	Keterangan     string `json:"keterangan"`
	Catatan        string `json:"catatan"`
}

type ReviewTujuanOpdUpdateRequest struct {
	Id         int    `json:"id"`
	Review     string `json:"review"`
	Keterangan string `json:"keterangan"`
	Catatan    string `json:"catatan"`
}
