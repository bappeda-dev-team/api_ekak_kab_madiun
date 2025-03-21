package programkegiatan

type ProgramKegiatanCreateRequest struct {
	Id          string                   `json:"id"`
	KodeProgram string                   `json:"kode_program"`
	NamaProgram string                   `json:"nama_program"`
	KodeOPD     string                   `json:"kode_opd"`
	Tahun       string                   `json:"tahun"`
	IsActive    bool                     `json:"is_active"`
	Indikator   []IndikatorCreateRequest `json:"indikator"`
}

type IndikatorCreateRequest struct {
	Id        string                `json:"id"`
	ProgramId string                `json:"program_id"`
	Indikator string                `json:"indikator"`
	Tahun     string                `json:"tahun"`
	Target    []TargetCreateRequest `json:"target"`
}

type TargetCreateRequest struct {
	Id          string `json:"id"`
	IndikatorId string `json:"indikator_id"`
	Tahun       string `json:"tahun"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}

type IndikatorRenstraCreateRequest struct {
	Kode         string `json:"kode"`
	KodeOpd      string `json:"kode_opd"`
	Indikator    string `json:"indikator"`
	Tahun        string `json:"tahun"`
	PaguAnggaran int64  `json:"pagu_anggaran"`
	Target       string `json:"target"`
	Satuan       string `json:"satuan"`
}
