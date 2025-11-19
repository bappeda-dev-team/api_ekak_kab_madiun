package pohonkinerja

import "ekak_kabupaten_madiun/model/web/opdmaster"

type CascadingRekinPegawaiResponse struct {
	Id              int                                 `json:"id"`
	Parent          *int                                `json:"parent"`
	NamaPohon       string                              `json:"nama_pohon"`
	JenisPohon      string                              `json:"jenis_pohon"`
	LevelPohon      int                                 `json:"level_pohon"`
	Keterangan      string                              `json:"keterangan"`
	Status          string                              `json:"status"`
	PerangkatDaerah opdmaster.OpdResponseForAll         `json:"perangkat_daerah"`
	IsActive        bool                                `json:"is_active"`
	RencanaKinerja  []RencanaKinerjaResponse            `json:"rencana_kinerja"`
	Indikator       []IndikatorResponse                 `json:"indikator,omitempty"`
	Program         []ProgramCascadingRekinResponse     `json:"program"`      // Untuk level 4 & 5
	Kegiatan        []KegiatanCascadingRekinResponse    `json:"kegiatan"`     // Untuk level 6
	SubKegiatan     []SubKegiatanCascadingRekinResponse `json:"sub_kegiatan"` // Untuk level 6
	PaguAnggaran    int64                               `json:"pagu_anggaran_total"`
}

type DetailRekinResponse struct {
	Id                 string                               `json:"id_rencana_kinerja,omitempty"`
	IdPohon            int                                  `json:"id_pohon,omitempty"`
	NamaRencanaKinerja string                               `json:"nama_rencana_kinerja,omitempty"`
	Tahun              string                               `json:"tahun,omitempty"`
	PegawaiId          string                               `json:"pegawai_id,omitempty"`
	LevelPohon         int                                  `json:"level_pohon"`
	Urusan             []UrusanCascadingRekinResponse       `json:"urusan,omitempty"`        // Untuk level 4
	BidangUrusan       []BidangUrusanCascadingRekinResponse `json:"bidang_urusan,omitempty"` // Untuk level 4
	Program            []ProgramRekinResponse               `json:"program,omitempty"`       // Untuk level 4 & 5
	Kegiatan           []KegiatanRekinResponse              `json:"kegiatan,omitempty"`      // Untuk level 6
	SubKegiatan        []SubKegiatanRekinResponse           `json:"sub_kegiatan,omitempty"`  // Untuk level 6
}

// Urusan response untuk level 4
type UrusanCascadingRekinResponse struct {
	KodeUrusan string `json:"kode_urusan"`
	NamaUrusan string `json:"nama_urusan"`
}

// Bidang Urusan response untuk level 4
type BidangUrusanCascadingRekinResponse struct {
	KodeBidangUrusan string `json:"kode_bidang_urusan"`
	NamaBidangUrusan string `json:"nama_bidang_urusan"`
}

// Program response untuk level 4 & 5
type ProgramCascadingRekinResponse struct {
	KodeProgram string              `json:"kode_program"`
	NamaProgram string              `json:"nama_program"`
	Indikator   []IndikatorResponse `json:"indikator"`
}

// Kegiatan response untuk level 6
type KegiatanCascadingRekinResponse struct {
	KodeKegiatan string              `json:"kode_kegiatan"`
	NamaKegiatan string              `json:"nama_kegiatan"`
	Indikator    []IndikatorResponse `json:"indikator"`
}

// SubKegiatan response untuk level 6
type SubKegiatanCascadingRekinResponse struct {
	KodeSubkegiatan string              `json:"kode_subkegiatan"`
	NamaSubkegiatan string              `json:"nama_subkegiatan"`
	Indikator       []IndikatorResponse `json:"indikator"`
}

// Program response untuk level 4 & 5
type ProgramRekinResponse struct {
	KodeProgram string `json:"kode_program"`
	NamaProgram string `json:"nama_program"`
}

// Kegiatan response untuk level 6
type KegiatanRekinResponse struct {
	KodeKegiatan string `json:"kode_kegiatan"`
	NamaKegiatan string `json:"nama_kegiatan"`
}

// SubKegiatan response untuk level 6
type SubKegiatanRekinResponse struct {
	KodeSubkegiatan string `json:"kode_subkegiatan"`
	NamaSubkegiatan string `json:"nama_subkegiatan"`
}
