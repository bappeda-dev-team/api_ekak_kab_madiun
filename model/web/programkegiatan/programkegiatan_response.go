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
	ProgramId string           `json:"program_id"`
	Indikator string           `json:"indikator"`
	Tahun     string           `json:"tahun"`
	Target    []TargetResponse `json:"target"`
}

type TargetResponse struct {
	Id          string `json:"id"`
	IndikatorId string `json:"indikator_id"`
	Tahun       string `json:"tahun"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}
