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
	KodeIndikator string `json:"kode_indikator"`
	Kode          string `json:"kode"`
	KodeOpd       string `json:"kode_opd"`
	Indikator     string `json:"indikator"`
	Tahun         string `json:"tahun"`
	Target        string `json:"target"`
	Satuan        string `json:"satuan"`
}

// Fungsi khusus anggaran (upsert)
type AnggaranRenstraRequest struct {
	KodeSubKegiatan string `json:"kode_subkegiatan" validate:"required"`
	KodeOpd         string `json:"kode_opd"         validate:"required"`
	Tahun           string `json:"tahun"            validate:"required"`
	Pagu            int64  `json:"pagu_indikatif" validate:"required"`
}

type AnggaranRenjaRequest struct {
	KodeSubKegiatan string `json:"kode_subkegiatan" validate:"required"`
	KodeOpd         string `json:"kode_opd"         validate:"required"`
	Tahun           string `json:"tahun"            validate:"required"`
	Pagu            int64  `json:"pagu_indikatif" validate:"required"`
}

type TargetRenjaRequest struct {
	Id          string `json:"id"`
	IndikatorId string `json:"indikator_id"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}

type BatchIndikatorRenjaRequest struct {
	Kode      string               `json:"kode" validate:"required"`
	KodeOpd   string               `json:"kode_opd" validate:"required"`
	Tahun     string               `json:"tahun" validate:"required"` // ← wajib ada
	Jenis     string               `json:"jenis" validate:"required"` // "ranwal"/"rankhir"
	Indikator []IndikatorRenjaItem `json:"indikator" validate:"required,min=1"`
}

type IndikatorRenjaItem struct {
	KodeIndikator string               `json:"kode_indikator"`
	Indikator     string               `json:"indikator" validate:"required"`
	Target        []TargetRenjaRequest `json:"target"`
}

type IndikatorRenjaCreateRequest struct {
	KodeIndikator string `json:"kode_indikator"`
	Kode          string `json:"kode" validate:"required"`
	KodeOpd       string `json:"kode_opd" validate:"required"`
	Indikator     string `json:"indikator" validate:"required"`
	Tahun         string `json:"tahun" validate:"required"`
	Jenis         string `json:"jenis"`
	Target        string `json:"target" validate:"required"`
	Satuan        string `json:"satuan" validate:"required"`
}
