package ikd

type IkdResponse struct {
	Id                     int    `json:"id"`
	NamaPohon              string `json:"nama_pohon"`
	Parent                 int    `json:"parent"`
	JenisPohon             string `json:"jenis_pohon"`
	LevelPohon             int    `json:"level_pohon"`
	KodeOpd                string `json:"kode_opd"`
	Keterangan             string `json:"keterangan"`
	KeteranganCrosscutting string `json:"keterangan_crosscutting"`
	Tahun                  string `json:"tahun"`
	Status                 string `json:"status"`
	IsActive               bool   `json:"is_active"`

	Pelaksana          []PelaksanaResponse  `json:"pelaksana"`
	SasaranOpd         []SasaranOpdResponse `json:"sasaran_opd"`
	ProgramOpd         []ProgramOpdResponse `json:"program_opd"`
	ProgramOpdTerpilih []ProgramOpdResponse `json:"program_opd_terpilih"`
}

type PelaksanaResponse struct {
	Id          string `json:"id"`
	PegawaiId   string `json:"pegawai_id"`
	Nip         string `json:"nip"`
	NamaPegawai string `json:"nama_pegawai"`
}

type SasaranOpdResponse struct {
	Id             int    `json:"id"`
	IdPohon        int    `json:"id_pohon"`
	NamaSasaranOpd string `json:"nama_sasaran_opd"`
	IdTujuanOpd    int    `json:"id_tujuan_opd"`
	NamaTujuanOpd  string `json:"nama_tujuan_opd"`
	TahunAwal      string `json:"tahun_awal"`
	TahunAkhir     string `json:"tahun_akhir"`
	JenisPeriode   string `json:"jenis_periode"`

	Indikator []IndikatorResponse `json:"indikator"`
}

type IndikatorResponse struct {
	Id               string `json:"id"`
	Indikator        string `json:"indikator"`
	RumusPerhitungan string `json:"rumus_perhitungan"`
	SumberData       string `json:"sumber_data"`

	Target []TargetResponse `json:"target"`
}

type TargetResponse struct {
	Id          string `json:"id"`
	IndikatorId string `json:"indikator_id"`
	Tahun       string `json:"tahun"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}

type ProgramOpdResponse struct {
	Id          int    `json:"id"`
	Parent      int    `json:"parent"`
	NamaProgram string `json:"nama_program"`
}

type ProgramOpdTerpilihResponse struct {
	Id             int `json:"id"`
	PohonKinerjaId int `json:"pohon_kinerja_id"`
	ProgramOpdId   int `json:"program_opd_id"`
}