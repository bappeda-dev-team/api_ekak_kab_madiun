package tujuanpemda

type TujuanPemdaResponse struct {
	Id          int                 `json:"id"`
	IdVisi      int                 `json:"id_visi,omitempty"`
	Visi        string              `json:"visi,omitempty"`
	IdMisi      int                 `json:"id_misi,omitempty"`
	Misi        string              `json:"misi,omitempty"`
	TujuanPemda string              `json:"tujuan_pemda"`
	TematikId   int                 `json:"tematik_id,omitempty"`
	NamaTematik string              `json:"nama_tematik,omitempty"`
	JenisPohon  string              `json:"jenis_pohon,omitempty"`
	PeriodeId   int                 `json:"periode_id,omitempty"`
	Periode     PeriodeResponse     `json:"periode"`
	Indikator   []IndikatorResponse `json:"indikator"`
}
type IndikatorResponse struct {
	Id               string           `json:"id"`
	Indikator        string           `json:"indikator"`
	RumusPerhitungan string           `json:"rumus_perhitungan"`
	SumberData       string           `json:"sumber_data"`
	Target           []TargetResponse `json:"target"`
}

type TargetResponse struct {
	Id     string `json:"id"`
	Target string `json:"target"`
	Satuan string `json:"satuan"`
	Tahun  string `json:"tahun"`
}

type PeriodeResponse struct {
	TahunAwal    string `json:"tahun_awal"`
	TahunAkhir   string `json:"tahun_akhir"`
	JenisPeriode string `json:"jenis_periode"`
}

type TujuanPemdaWithPokinResponse struct {
	PokinId     int                   `json:"pokin_id"`
	NamaPohon   string                `json:"nama_tematik"`
	JenisPohon  string                `json:"jenis_pohon"`
	LevelPohon  int                   `json:"level_pohon"`
	IsActive    bool                  `json:"is_active"`
	Keterangan  string                `json:"keterangan"`
	TahunPokin  string                `json:"tahun_pokin"`
	TujuanPemda []TujuanPemdaResponse `json:"tujuan_pemda"`
}

// pokin with periode
type PokinWithPeriodeResponse struct {
	Id         int                      `json:"id"`
	NamaPohon  string                   `json:"nama_pohon"`
	Parent     int                      `json:"parent,omitempty"`
	JenisPohon string                   `json:"jenis_pohon,omitempty"`
	LevelPohon int                      `json:"level_pohon,omitempty"`
	KodeOpd    string                   `json:"kode_opd,omitempty"`
	Keterangan string                   `json:"keterangan,omitempty"`
	Tahun      string                   `json:"tahun,omitempty"`
	Status     string                   `json:"status,omitempty"`
	Periode    PokinPeriodeResponse     `json:"periode"`
	Indikator  []PokinIndikatorResponse `json:"indikator"`
}

type PokinPeriodeResponse struct {
	Id         int    `json:"id"`
	TahunAwal  string `json:"tahun_awal"`
	TahunAkhir string `json:"tahun_akhir"`
}

type PokinIndikatorResponse struct {
	Id               string                `json:"id"`
	Indikator        string                `json:"indikator"`
	RumusPerhitungan string                `json:"rumus_perhitungan"`
	SumberData       string                `json:"sumber_data"`
	Target           []PokinTargetResponse `json:"target"`
}

type PokinTargetResponse struct {
	Id     string `json:"id"`
	Target string `json:"target"`
	Satuan string `json:"satuan"`
	Tahun  string `json:"tahun"`
}
