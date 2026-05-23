package domain

type IkdDetail struct {
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

	Pelaksana          []PelaksanaDetail          `json:"pelaksana"`
	SasaranOpd         []SasaranOpdDetail         `json:"sasaran_opd"`
	ProgramOpd         []ProgramOpdDetail         `json:"program_opd"`
	ProgramOpdTerpilih []ProgramOpdTerpilihDetail `json:"program_opd_terpilih"`
}

type PelaksanaDetail struct {
	Id          string `json:"id"`
	PegawaiId   string `json:"pegawai_id"`
	Nip         string `json:"nip"`
	NamaPegawai string `json:"nama_pegawai"`
}

type ProgramOpdDetail struct {
	Id          int    `json:"id"`
	Parent      int    `json:"parent"`
	NamaProgram string `json:"nama_program"`
}

type ProgramOpdTerpilihDetail struct {
	Id          int    `json:"id"`
	TacticalId  int    `json:"tactical_id"`
	Parent      int    `json:"parent"`
	NamaProgram string `json:"nama_program"`
	IsLocked    bool   `json:"is_locked"`
}

type ProgramOpdTerpilih struct {
	Id             int
	PohonKinerjaId int
	ProgramOpdId   int
}

type PokinIkd struct {
	Id        int
	NamaPokin string
}
