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
	Id            string           `json:"id,omitempty"`
	KodeIndikator string           `json:"kode_indikator"`
	Kode          string           `json:"kode,omitempty"`
	KodeOpd       string           `json:"kode_opd,omitempty"`
	ProgramId     string           `json:"program_id,omitempty"`
	Indikator     string           `json:"indikator"`
	PaguAnggaran  int64            `json:"pagu_anggaran,omitempty"`
	Tahun         string           `json:"tahun"`
	Target        []TargetResponse `json:"target"`
	StatusTarget  bool             `json:"status_target_renja,omitempty"`
}

type IndikatorMatrixResponse struct {
	Id            string `json:"id,omitempty"`
	KodeIndikator string `json:"kode_indikator"`
	Kode          string `json:"kode,omitempty"`
	KodeOpd       string `json:"kode_opd,omitempty"`
	ProgramId     string `json:"program_id,omitempty"`
	Indikator     string `json:"indikator"`
	// PaguAnggaran  int64  `json:"pagu_anggaran,omitempty"`
	Tahun        string `json:"tahun"`
	Target       string `json:"target"`
	Satuan       string `json:"satuan"`
	StatusTarget bool   `json:"status_target_renja,omitempty"`
}

type TargetResponse struct {
	Id          string `json:"id"`
	IndikatorId string `json:"indikator_id"`
	Tahun       string `json:"tahun,omitempty"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}

type UrusanDetailResponse struct {
	KodeOpd           string                      `json:"kode_opd"`
	TahunAwal         string                      `json:"tahun_awal"`
	TahunAkhir        string                      `json:"tahun_akhir"`
	PaguAnggaranTotal []PaguAnggaranTotalResponse `json:"pagu_total"`
	Urusan            []UrusanResponse            `json:"urusan"`
	Tahun             string                      `json:"tahun,omitempty"`
}

type UrusanResponse struct {
	Kode         string                      `json:"kode"`
	Nama         string                      `json:"nama"`
	Jenis        string                      `json:"jenis"`
	Anggaran     []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	Indikator    []IndikatorMatrixResponse   `json:"indikator"`
	BidangUrusan []BidangUrusanResponse      `json:"bidang_urusan"`
}

type BidangUrusanResponse struct {
	Kode      string                      `json:"kode"`
	Nama      string                      `json:"nama"`
	Jenis     string                      `json:"jenis"`
	Anggaran  []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	Indikator []IndikatorMatrixResponse   `json:"indikator"`
	Program   []ProgramResponse           `json:"program"`
}

type ProgramResponse struct {
	Kode      string                      `json:"kode"`
	Nama      string                      `json:"nama"`
	Jenis     string                      `json:"jenis"`
	Anggaran  []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	Indikator []IndikatorMatrixResponse   `json:"indikator"`
	Kegiatan  []KegiatanResponse          `json:"kegiatan"`
}

type KegiatanResponse struct {
	Kode        string                      `json:"kode"`
	Nama        string                      `json:"nama"`
	Jenis       string                      `json:"jenis"`
	Anggaran    []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	Indikator   []IndikatorMatrixResponse   `json:"indikator"`
	SubKegiatan []SubKegiatanResponse       `json:"subkegiatan"`
}

type SubKegiatanResponse struct {
	Kode          string                      `json:"kode"`
	Nama          string                      `json:"nama"`
	Jenis         string                      `json:"jenis"`
	Tahun         string                      `json:"tahun,omitempty"`
	PegawaiId     string                      `json:"pegawai_id"`
	NamaPegawai   string                      `json:"nama_pegawai"`
	Anggaran      []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	TotalAnggaran int64                       `json:"total_anggaran,omitempty"`
	Indikator     []IndikatorMatrixResponse   `json:"indikator"`
}

type PaguAnggaranTotalResponse struct {
	Tahun        string `json:"tahun"`
	PaguAnggaran int64  `json:"pagu_indikatif"`
}

type AnggaranRenstraResponse struct {
	KodeSubKegiatan string `json:"kode_subkegiatan" validate:"required"`
	KodeOpd         string `json:"kode_opd"         validate:"required"`
	Tahun           string `json:"tahun"            validate:"required"`
	Pagu            int64  `json:"pagu_indikatif" validate:"required"`
}

type AnggaranRenjaResponse struct {
	KodeSubKegiatan string `json:"kode_subkegiatan" validate:"required"`
	KodeOpd         string `json:"kode_opd"         validate:"required"`
	Tahun           string `json:"tahun"            validate:"required"`
	Pagu            int64  `json:"pagu_indikatif" validate:"required"`
}

type BatchIndikatorRenjaResponse struct {
	Kode      string                         `json:"kode"`
	KodeOpd   string                         `json:"kode_opd"`
	Tahun     string                         `json:"tahun"`
	Jenis     string                         `json:"jenis"`
	Indikator []IndikatorRenjaUpsertResponse `json:"indikator"`
}
type IndikatorRenjaUpsertResponse struct {
	KodeIndikator string         `json:"kode_indikator"`
	Indikator     string         `json:"indikator"`
	Jenis         string         `json:"jenis"`
	Target        TargetResponse `json:"target"`
}

type IndikatorUpsertResponse struct {
	KodeIndikator string `json:"kode_indikator"`
	Kode          string `json:"kode"`
	KodeOpd       string `json:"kode_opd"`
	Indikator     string `json:"indikator"`
	Tahun         string `json:"tahun"`
	Jenis         string `json:"jenis,omitempty"`
	// Target        TargetResponse `json:"target"`
	Id          string `json:"id"`
	IndikatorId string `json:"indikator_id"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}
