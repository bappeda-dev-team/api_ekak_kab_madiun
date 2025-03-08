package pohonkinerja

type ReviewResponse struct {
	Id             int    `json:"id"`
	IdPohonKinerja int    `json:"id_pohon_kinerja"`
	Review         string `json:"review"`
	Keterangan     string `json:"keterangan"`
	CreatedBy      string `json:"created_by,omitempty"`
	NamaPegawai    string `json:"nama_pegawai,omitempty"`
	JenisPokin     string `json:"jenis_pokin"`
}
