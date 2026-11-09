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
	Tahun         string `json:"tahun"`
	// single-slot (ranwal view)
	Target       string `json:"target"`
	Satuan       string `json:"satuan"`
	StatusTarget bool   `json:"status_target_renja,omitempty"`
	// dual-slot rankhir view: target_ranwal + target_rankhir
	TargetRanwal    string `json:"target_ranwal,omitempty"`
	SatuanRanwal    string `json:"satuan_ranwal,omitempty"`
	TargetRankhir   string `json:"target_rankhir,omitempty"`
	SatuanRankhir   string `json:"satuan_rankhir,omitempty"`
	// penetapan view
	TargetPenetapan string `json:"target_penetapan,omitempty"`
	SatuanPenetapan string `json:"satuan_penetapan,omitempty"`
}

type TargetResponse struct {
	Id          string `json:"id"`
	IndikatorId string `json:"indikator_id"`
	Tahun       string `json:"tahun,omitempty"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
	Jenis       string `json:"jenis,omitempty"`
}

// IndikatorBaselineResponse: 1 indikator tahun awal (indikator pertama jika ada lebih dari satu),
// beserta target existing di target_baseline.
type IndikatorBaselineResponse struct {
	KodeIndikator  string           `json:"kode_indikator"`
	Kode           string           `json:"kode,omitempty"`
	KodeOpd        string           `json:"kode_opd,omitempty"`
	Indikator      string           `json:"indikator"`
	Tahun          string           `json:"tahun"`
	TargetBaseline []TargetResponse `json:"target_baseline"`
}

type UrusanDetailResponse struct {
	KodeOpd           string                      `json:"kode_opd"`
	TahunAwal         string                      `json:"tahun_awal"`
	TahunAkhir        string                      `json:"tahun_akhir"`
	PaguAnggaranTotal []PaguAnggaranTotalResponse `json:"pagu_total"`
	Urusan            []UrusanResponse            `json:"urusan"`
	Tahun             string                      `json:"tahun,omitempty"`
	// total pagu per layer (level root/detail)
	PaguRanwal      int64 `json:"pagu_total_ranwal,omitempty"`
	PaguRankhir     int64 `json:"pagu_total_rankhir,omitempty"`
	PaguPenetapan   int64 `json:"pagu_total_penetapan,omitempty"`
}

type UrusanResponse struct {
	Kode               string                      `json:"kode"`
	Nama               string                      `json:"nama"`
	Jenis              string                      `json:"jenis"`
	Anggaran           []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	PaguRanwal         int64                       `json:"pagu_ranwal,omitempty"`
	PaguRankhir        int64                       `json:"pagu_rankhir,omitempty"`
	PaguPenetapan      int64                       `json:"pagu_penetapan,omitempty"`
	Indikator          []IndikatorMatrixResponse   `json:"indikator"`
	IndikatorBaseline  []IndikatorBaselineResponse `json:"indikator_baseline"`
	// Diisi untuk rankhir view: indikator per layer terpisah
	IndikatorRanwal    []IndikatorMatrixResponse `json:"indikator_ranwal,omitempty"`
	IndikatorRankhir   []IndikatorMatrixResponse `json:"indikator_rankhir,omitempty"`
	IndikatorPenetapan []IndikatorMatrixResponse `json:"indikator_penetapan,omitempty"`
	BidangUrusan       []BidangUrusanResponse    `json:"bidang_urusan"`
}

type BidangUrusanResponse struct {
	Kode               string                      `json:"kode"`
	Nama               string                      `json:"nama"`
	Jenis              string                      `json:"jenis"`
	Anggaran           []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	PaguRanwal         int64                       `json:"pagu_ranwal,omitempty"`
	PaguRankhir        int64                       `json:"pagu_rankhir,omitempty"`
	PaguPenetapan      int64                       `json:"pagu_penetapan,omitempty"`
	Indikator          []IndikatorMatrixResponse   `json:"indikator"`
	IndikatorBaseline  []IndikatorBaselineResponse `json:"indikator_baseline"`
	// Diisi untuk rankhir view: indikator per layer terpisah
	IndikatorRanwal    []IndikatorMatrixResponse `json:"indikator_ranwal,omitempty"`
	IndikatorRankhir   []IndikatorMatrixResponse `json:"indikator_rankhir,omitempty"`
	IndikatorPenetapan []IndikatorMatrixResponse `json:"indikator_penetapan,omitempty"`
	Program            []ProgramResponse         `json:"program"`
}

type ProgramResponse struct {
	Kode               string                      `json:"kode"`
	Nama               string                      `json:"nama"`
	Jenis              string                      `json:"jenis"`
	Anggaran           []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	PaguRanwal         int64                       `json:"pagu_ranwal,omitempty"`
	PaguRankhir        int64                       `json:"pagu_rankhir,omitempty"`
	PaguPenetapan      int64                       `json:"pagu_penetapan,omitempty"`
	Indikator          []IndikatorMatrixResponse   `json:"indikator"`
	IndikatorBaseline  []IndikatorBaselineResponse `json:"indikator_baseline"`
	// Diisi untuk rankhir view: indikator per layer terpisah
	IndikatorRanwal    []IndikatorMatrixResponse `json:"indikator_ranwal,omitempty"`
	IndikatorRankhir   []IndikatorMatrixResponse `json:"indikator_rankhir,omitempty"`
	IndikatorPenetapan []IndikatorMatrixResponse `json:"indikator_penetapan,omitempty"`
	Kegiatan           []KegiatanResponse        `json:"kegiatan"`
}

type KegiatanResponse struct {
	Kode               string                      `json:"kode"`
	Nama               string                      `json:"nama"`
	Jenis              string                      `json:"jenis"`
	Anggaran           []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	PaguRanwal         int64                       `json:"pagu_ranwal,omitempty"`
	PaguRankhir        int64                       `json:"pagu_rankhir,omitempty"`
	PaguPenetapan      int64                       `json:"pagu_penetapan,omitempty"`
	Indikator          []IndikatorMatrixResponse   `json:"indikator"`
	IndikatorBaseline  []IndikatorBaselineResponse `json:"indikator_baseline"`
	// Diisi untuk rankhir view: indikator per layer terpisah
	IndikatorRanwal    []IndikatorMatrixResponse `json:"indikator_ranwal,omitempty"`
	IndikatorRankhir   []IndikatorMatrixResponse `json:"indikator_rankhir,omitempty"`
	IndikatorPenetapan []IndikatorMatrixResponse `json:"indikator_penetapan,omitempty"`
	SubKegiatan        []SubKegiatanResponse     `json:"subkegiatan"`
}

type SubKegiatanResponse struct {
	Kode               string                      `json:"kode"`
	Nama               string                      `json:"nama"`
	Jenis              string                      `json:"jenis"`
	Tahun              string                      `json:"tahun,omitempty"`
	PegawaiId          string                      `json:"pegawai_id"`
	NamaPegawai        string                      `json:"nama_pegawai"`
	// Status hanya terisi untuk rankhir view — baris sintetis dari renstra yg tidak ada rekin
	Status             string                      `json:"status,omitempty"`
	Anggaran           []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	TotalAnggaran      int64                       `json:"total_anggaran,omitempty"`
	// flat pagu per layer
	PaguRanwal         int64                     `json:"pagu_ranwal,omitempty"`
	PaguRankhir        int64                     `json:"pagu_rankhir,omitempty"`
	PaguPenetapan      int64                     `json:"pagu_penetapan,omitempty"`
	Indikator          []IndikatorMatrixResponse `json:"indikator"`
	IndikatorBaseline  []IndikatorBaselineResponse `json:"indikator_baseline"`
	// Diisi untuk rankhir view: indikator per layer terpisah
	IndikatorRanwal    []IndikatorMatrixResponse `json:"indikator_ranwal,omitempty"`
	IndikatorRankhir   []IndikatorMatrixResponse `json:"indikator_rankhir,omitempty"`
	IndikatorPenetapan []IndikatorMatrixResponse `json:"indikator_penetapan,omitempty"`
}

type PaguAnggaranTotalResponse struct {
	Tahun        string `json:"tahun"`
	PaguAnggaran int64  `json:"pagu_indikatif"`
	JenisPagu    string `json:"jenis_pagu"`
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

type IndikatorPeriodResponse struct {
	KodeIndikator string           `json:"kode_indikator"`
	Kode          string           `json:"kode,omitempty"`
	KodeOpd       string           `json:"kode_opd,omitempty"`
	Indikator     string           `json:"indikator"`
	Tahun         string           `json:"tahun,omitempty"`
	Target        []TargetResponse `json:"target"`
}

type UrusanDetailV2Response struct {
	KodeOpd           string                      `json:"kode_opd"`
	TahunAwal         string                      `json:"tahun_awal"`
	TahunAkhir        string                      `json:"tahun_akhir"`
	PaguAnggaranTotal []PaguAnggaranTotalResponse `json:"pagu_total"`
	Urusan            []UrusanV2Response          `json:"urusan"`
}

type UrusanV2Response struct {
	Kode         string                      `json:"kode"`
	Nama         string                      `json:"nama"`
	Jenis        string                      `json:"jenis"`
	Anggaran     []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	Indikator    []IndikatorPeriodResponse   `json:"indikator"`
	Outcome      []OutcomeMatrixItemResponse `json:"outcome"`
	BidangUrusan []BidangUrusanV2Response    `json:"bidang_urusan"`
}

type BidangUrusanV2Response struct {
	Kode      string                      `json:"kode"`
	Nama      string                      `json:"nama"`
	Jenis     string                      `json:"jenis"`
	Anggaran  []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	Indikator []IndikatorPeriodResponse   `json:"indikator"`
	Outcome   []OutcomeMatrixItemResponse `json:"outcome"`
	Program   []ProgramV2Response         `json:"program"`
}

type ProgramV2Response struct {
	Kode      string                      `json:"kode"`
	Nama      string                      `json:"nama"`
	Jenis     string                      `json:"jenis"`
	Anggaran  []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	Indikator []IndikatorPeriodResponse   `json:"indikator"`
	Outcome   []OutcomeMatrixItemResponse `json:"outcome"`
	Kegiatan  []KegiatanV2Response        `json:"kegiatan"`
}

type KegiatanV2Response struct {
	Kode        string                      `json:"kode"`
	Nama        string                      `json:"nama"`
	Jenis       string                      `json:"jenis"`
	Anggaran    []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	Indikator   []IndikatorPeriodResponse   `json:"indikator"`
	Outcome     []OutcomeMatrixItemResponse `json:"outcome"`
	SubKegiatan []SubKegiatanV2Response     `json:"subkegiatan"`
}

type SubKegiatanV2Response struct {
	Kode        string                      `json:"kode"`
	Nama        string                      `json:"nama"`
	Jenis       string                      `json:"jenis"`
	PegawaiId   string                      `json:"pegawai_id"`
	NamaPegawai string                      `json:"nama_pegawai"`
	Anggaran    []PaguAnggaranTotalResponse `json:"anggaran,omitempty"`
	Indikator   []IndikatorPeriodResponse   `json:"indikator"`
	Outcome     []OutcomeMatrixItemResponse `json:"outcome"`
}

type OutcomeMatrixItemResponse struct {
	Id      int    `json:"id"`
	KodeOpd string `json:"kode_opd"`
	Kode    string `json:"kode"`
	Jenis   string `json:"jenis"`
	Outcome string `json:"outcome"`
}

type IndikatorRenstraUpdateResponse struct {
	KodeIndikator string `json:"kode_indikator"`
	Kode          string `json:"kode"`
	KodeOpd       string `json:"kode_opd"`
	Indikator     string `json:"indikator"`
	Tahun         string `json:"tahun,omitempty"`
	Jenis         string `json:"jenis,omitempty"`
}

type IndikatorV2UpsertResponse struct {
	KodeIndikator string           `json:"kode_indikator"`
	Kode          string           `json:"kode"`
	KodeOpd       string           `json:"kode_opd"`
	Indikator     string           `json:"indikator"`
	Jenis         string           `json:"jenis,omitempty"`
	Target        []TargetResponse `json:"target"`
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

// TargetRenjaLayerResponse — deprecated: gunakan renja.TargetLayerResponse untuk matrix renja.
type TargetRenjaLayerResponse struct {
	KodeIndikator string `json:"kode_indikator"`
	Tahun         string `json:"tahun"`
	Target        string `json:"target"`
	Satuan        string `json:"satuan"`
	Jenis         string `json:"jenis"`
}
