package renja

// UrusanDetailResponse — root response GET matrix renja (ranwal, rankhir, penetapan).
type UrusanDetailResponse struct {
	KodeOpd         string           `json:"kode_opd"`
	Tahun           string           `json:"tahun,omitempty"`
	PaguRanwal      int64            `json:"pagu_total_ranwal,omitempty"`
	PaguRankhir     int64            `json:"pagu_total_rankhir,omitempty"`
	PaguPenetapan   int64            `json:"pagu_total_penetapan,omitempty"`
	Urusan          []UrusanResponse `json:"urusan"`
}

type UrusanResponse struct {
	Kode               string                    `json:"kode"`
	Nama               string                    `json:"nama"`
	Jenis              string                    `json:"jenis"`
	PaguRanwal         int64                     `json:"pagu_ranwal,omitempty"`
	PaguRankhir        int64                     `json:"pagu_rankhir,omitempty"`
	PaguPenetapan      int64                     `json:"pagu_penetapan,omitempty"`
	Indikator          []IndikatorMatrixResponse `json:"indikator"`
	IndikatorRanwal    []IndikatorMatrixResponse `json:"indikator_ranwal,omitempty"`
	IndikatorRankhir   []IndikatorMatrixResponse `json:"indikator_rankhir,omitempty"`
	IndikatorPenetapan []IndikatorMatrixResponse `json:"indikator_penetapan,omitempty"`
	BidangUrusan       []BidangUrusanResponse    `json:"bidang_urusan"`
}

type BidangUrusanResponse struct {
	Kode               string                    `json:"kode"`
	Nama               string                    `json:"nama"`
	Jenis              string                    `json:"jenis"`
	PaguRanwal         int64                     `json:"pagu_ranwal,omitempty"`
	PaguRankhir        int64                     `json:"pagu_rankhir,omitempty"`
	PaguPenetapan      int64                     `json:"pagu_penetapan,omitempty"`
	Indikator          []IndikatorMatrixResponse `json:"indikator"`
	IndikatorRanwal    []IndikatorMatrixResponse `json:"indikator_ranwal,omitempty"`
	IndikatorRankhir   []IndikatorMatrixResponse `json:"indikator_rankhir,omitempty"`
	IndikatorPenetapan []IndikatorMatrixResponse `json:"indikator_penetapan,omitempty"`
	Program            []ProgramResponse         `json:"program"`
}

type ProgramResponse struct {
	Kode               string                    `json:"kode"`
	Nama               string                    `json:"nama"`
	Jenis              string                    `json:"jenis"`
	PaguRanwal         int64                     `json:"pagu_ranwal,omitempty"`
	PaguRankhir        int64                     `json:"pagu_rankhir,omitempty"`
	PaguPenetapan      int64                     `json:"pagu_penetapan,omitempty"`
	Indikator          []IndikatorMatrixResponse `json:"indikator"`
	IndikatorRanwal    []IndikatorMatrixResponse `json:"indikator_ranwal,omitempty"`
	IndikatorRankhir   []IndikatorMatrixResponse `json:"indikator_rankhir,omitempty"`
	IndikatorPenetapan []IndikatorMatrixResponse `json:"indikator_penetapan,omitempty"`
	Kegiatan           []KegiatanResponse        `json:"kegiatan"`
}

type KegiatanResponse struct {
	Kode               string                    `json:"kode"`
	Nama               string                    `json:"nama"`
	Jenis              string                    `json:"jenis"`
	PaguRanwal         int64                     `json:"pagu_ranwal,omitempty"`
	PaguRankhir        int64                     `json:"pagu_rankhir,omitempty"`
	PaguPenetapan      int64                     `json:"pagu_penetapan,omitempty"`
	Indikator          []IndikatorMatrixResponse `json:"indikator"`
	IndikatorRanwal    []IndikatorMatrixResponse `json:"indikator_ranwal,omitempty"`
	IndikatorRankhir   []IndikatorMatrixResponse `json:"indikator_rankhir,omitempty"`
	IndikatorPenetapan []IndikatorMatrixResponse `json:"indikator_penetapan,omitempty"`
	SubKegiatan        []SubKegiatanResponse     `json:"subkegiatan"`
}

type SubKegiatanResponse struct {
	Kode               string                    `json:"kode"`
	Nama               string                    `json:"nama"`
	Jenis              string                    `json:"jenis"`
	Tahun              string                    `json:"tahun,omitempty"`
	PegawaiId          string                    `json:"pegawai_id"`
	NamaPegawai        string                    `json:"nama_pegawai"`
	Status             string                    `json:"status,omitempty"`
	PaguRanwal         int64                     `json:"pagu_ranwal,omitempty"`
	PaguRankhir        int64                     `json:"pagu_rankhir,omitempty"`
	PaguPenetapan      int64                     `json:"pagu_penetapan,omitempty"`
	Indikator          []IndikatorMatrixResponse `json:"indikator"`
	IndikatorRanwal    []IndikatorMatrixResponse `json:"indikator_ranwal,omitempty"`
	IndikatorRankhir   []IndikatorMatrixResponse `json:"indikator_rankhir,omitempty"`
	IndikatorPenetapan []IndikatorMatrixResponse `json:"indikator_penetapan,omitempty"`
}
