package programkegiatan

type ProgramKegiatanUpdateRequest struct {
	Id          string                   `json:"id"`
	KodeProgram string                   `json:"kode_program"`
	NamaProgram string                   `json:"nama_program"`
	KodeOPD     string                   `json:"kode_opd"`
	Tahun       string                   `json:"tahun"`
	IsActive    bool                     `json:"is_active"`
	Indikator   []IndikatorUpdateRequest `json:"indikator"`
}

type IndikatorUpdateRequest struct {
	Id        string                `json:"id"`
	ProgramId string                `json:"program_id"`
	Indikator string                `json:"indikator"`
	Tahun     string                `json:"tahun"`
	Target    []TargetUpdateRequest `json:"target"`
}

type TargetUpdateRequest struct {
	Id          string `json:"id"`
	IndikatorId string `json:"indikator_id"`
	Tahun       string `json:"tahun"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}

type UpdateIndikatorRequest struct {
	Id           string `json:"id" validate:"required"`
	Kode         string `json:"kode" validate:"required"`
	KodeOpd      string `json:"kode_opd" validate:"required"`
	Indikator    string `json:"indikator" validate:"required"`
	Tahun        string `json:"tahun" validate:"required"`
	PaguAnggaran int64  `json:"pagu_anggaran"`
	Target       string `json:"target" validate:"required"`
	Satuan       string `json:"satuan" validate:"required"`
}
