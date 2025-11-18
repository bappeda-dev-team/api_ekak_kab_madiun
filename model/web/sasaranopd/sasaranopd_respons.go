package sasaranopd

type SasaranOpdResponse struct {
	IdPohon    int                        `json:"id_pohon"`
	KodeOpd    string                     `json:"kode_opd,omitempty"`
	NamaOpd    string                     `json:"nama_opd,omitempty"`
	NamaPohon  string                     `json:"nama_pohon"`
	JenisPohon string                     `json:"jenis_pohon"`
	TahunPohon string                     `json:"tahun_pohon"`
	LevelPohon int                        `json:"level_pohon"`
	Pelaksana  []PelaksanaOpdResponse     `json:"pelaksana"`
	SasaranOpd []SasaranOpdDetailResponse `json:"sasaran_opd"`
}

type SasaranOpdDetailResponse struct {
	Id             string              `json:"id"`
	NamaSasaranOpd string              `json:"nama_sasaran_opd"`
	IdTujuanOpd    int                 `json:"id_tujuan_opd"`
	NamaTujuanOpd  string              `json:"nama_tujuan_opd"`
	TahunAwal      string              `json:"tahun_awal"`
	TahunAkhir     string              `json:"tahun_akhir"`
	JenisPeriode   string              `json:"jenis_periode"`
	Indikator      []IndikatorResponse `json:"indikator"`
}

type PelaksanaOpdResponse struct {
	Id          string `json:"id"`
	PegawaiId   string `json:"pegawai_id"`
	Nip         string `json:"nip"`
	NamaPegawai string `json:"nama_pegawai"`
}

type IndikatorResponse struct {
	Id               string           `json:"id"`
	Indikator        string           `json:"indikator"`
	RumusPerhitungan string           `json:"rumus_perhitungan"`
	SumberData       string           `json:"sumber_data"`
	Target           []TargetResponse `json:"target"`
}

type ManualIKResponse struct {
	Formula    string `json:"formula"`
	SumberData string `json:"sumber_data"`
}

type TargetResponse struct {
	Id     string `json:"id"`
	Tahun  string `json:"tahun"`
	Target string `json:"target"`
	Satuan string `json:"satuan"`
}

// respons create update

type SasaranOpdCreateResponse struct {
	IdPohon        int               `json:"id_pohon"`
	NamaSasaranOpd string            `json:"nama_sasaran_opd"`
	NamaTujuanOpd  string            `json:"nama_tujuan_opd"`
	TahunAwal      string            `json:"tahun_awal"`
	TahunAkhir     string            `json:"tahun_akhir"`
	JenisPeriode   string            `json:"jenis_periode"`
	Indikator      []IndikatorDetail `json:"indikator"`
}

type IndikatorDetail struct {
	Id               string         `json:"id"`
	Indikator        string         `json:"indikator"`
	RumusPerhitungan string         `json:"rumus_perhitungan"`
	SumberData       string         `json:"sumber_data"`
	Target           []TargetDetail `json:"target"`
}

type TargetDetail struct {
	Id     string `json:"id"`
	Tahun  string `json:"tahun"`
	Target string `json:"target"`
	Satuan string `json:"satuan"`
}
