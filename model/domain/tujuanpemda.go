package domain

type TujuanPemda struct {
	Id          int
	TematikId   int
	NamaTematik string
	TujuanPemda string
	JenisPohon  string
	IsActive    bool
	// TujuanPemdaId    int
	// NamaTujuanPemda  string
	PeriodeId         int
	TahunAwalPeriode  string
	TahunAkhirPeriode string
	JenisPeriode      string
	Periode           Periode
	Indikator         []Indikator
	IdVisi            int
	Visi              string
	IdMisi            int
	Misi              string
}

type TujuanPemdaWithPokin struct {
	PokinId     int    `json:"pokin_id"`
	NamaPohon   string `json:"nama_pohon"`
	JenisPohon  string `json:"jenis_pohon"`
	LevelPohon  int    `json:"level_pohon"`
	IsActive    bool
	KodeOpd     string `json:"kode_opd"`
	Keterangan  string `json:"keterangan"`
	IdVisi      int
	Visi        string
	IdMisi      int
	Misi        string
	TahunPokin  string        `json:"tahun_pokin"`
	TujuanPemda []TujuanPemda `json:"tujuan_pemda,omitempty"`
}
