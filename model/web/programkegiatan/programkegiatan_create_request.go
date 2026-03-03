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

type BatchIndikatorRenstraCreateRequest struct {
	Indikator []IndikatorRenstraCreateRequest `json:"indikator" validate:"required,min=1"`
}

type IndikatorRenstraCreateRequest struct {
	Kode      string `json:"kode"`
	KodeOpd   string `json:"kode_opd"`
	Indikator string `json:"indikator"`
	Tahun     string `json:"tahun"`
	Target    string `json:"target"`
	Satuan    string `json:"satuan"`
}

// Fungsi khusus anggaran (upsert)
type AnggaranRenstraRequest struct {
	KodeSubKegiatan string `json:"kode_subkegiatan" validate:"required"`
	KodeOpd         string `json:"kode_opd"         validate:"required"`
	Tahun           string `json:"tahun"            validate:"required"`
	Pagu            int64  `json:"pagu_indikatif" validate:"required"`
}
