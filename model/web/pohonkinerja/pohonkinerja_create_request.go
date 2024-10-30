package pohonkinerja

type PohonKinerjaCreateRequest struct {
	Parent     int    `json:"parent"`
	NamaPohon  string `json:"nama_pohon"`
	JenisPohon string `json:"jenis_pohon"`
	LevelPohon int    `json:"level_pohon"`
	KodeOpd    string `json:"kode_opd"`
	Keterangan string `json:"keterangan"`
	Tahun      string `json:"tahun"`
}