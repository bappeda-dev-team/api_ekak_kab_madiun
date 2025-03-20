package programkegiatan

type ProgramKegiatanResponse struct {
	Id          string              `json:"id"`
	KodeProgram string              `json:"kode_program"`
	NamaProgram string              `json:"nama_program"`
	Tahun       string              `json:"tahun"`
	IsActive    bool                `json:"is_active"`
	Indikator   []IndikatorResponse `json:"indikator"`
}

type IndikatorResponse struct {
	Id        string           `json:"id"`
	Kode      string           `json:"kode"`
	KodeOpd   string           `json:"kode_opd"`
	ProgramId string           `json:"program_id"`
	Indikator string           `json:"indikator"`
	Tahun     string           `json:"tahun"`
	Target    []TargetResponse `json:"target"`
}

type TargetResponse struct {
	Id          string `json:"id"`
	IndikatorId string `json:"indikator_id"`
	Tahun       string `json:"tahun,omitempty"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}

type UrusanDetailResponse struct {
	KodeOpd    string         `json:"kode_opd"`
	TahunAwal  string         `json:"tahun_awal"`
	TahunAkhir string         `json:"tahun_akhir"`
	Urusan     UrusanResponse `json:"urusan"`
}

type UrusanResponse struct {
	Kode         string                 `json:"kode"`
	Nama         string                 `json:"nama"`
	Indikator    []IndikatorResponse    `json:"indikator"`
	BidangUrusan []BidangUrusanResponse `json:"bidang_urusan"`
}

type BidangUrusanResponse struct {
	Kode      string              `json:"kode"`
	Nama      string              `json:"nama"`
	Indikator []IndikatorResponse `json:"indikator"`
	Program   []ProgramResponse   `json:"program"`
}

type ProgramResponse struct {
	Kode      string              `json:"kode"`
	Nama      string              `json:"nama"`
	Indikator []IndikatorResponse `json:"indikator"`
	Kegiatan  []KegiatanResponse  `json:"kegiatan"`
}

type KegiatanResponse struct {
	Kode        string                `json:"kode"`
	Nama        string                `json:"nama"`
	Indikator   []IndikatorResponse   `json:"indikator"`
	SubKegiatan []SubKegiatanResponse `json:"subkegiatan"`
}

type SubKegiatanResponse struct {
	Kode      string              `json:"kode"`
	Nama      string              `json:"nama"`
	Tahun     string              `json:"tahun"`
	Indikator []IndikatorResponse `json:"indikator"`
}
