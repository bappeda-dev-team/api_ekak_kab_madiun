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

type ReviewTematikResponse struct {
	IdTematik  int                    `json:"id_tematik"`
	NamaPohon  string                 `json:"nama_pohon"`
	LevelPohon int                    `json:"level_pohon"`
	Review     []ReviewDetailResponse `json:"review"`
}

type ReviewDetailResponse struct {
	IdPohon     int    `json:"id_pohon"`
	Parent      int    `json:"parent"`
	NamaPohon   string `json:"nama_pohon"`
	LevelPohon  int    `json:"level_pohon"`
	JenisPohon  string `json:"jenis_pohon"`
	Review      string `json:"review"`
	Keterangan  string `json:"keterangan"`
	NamaPegawai string `json:"created_by"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
