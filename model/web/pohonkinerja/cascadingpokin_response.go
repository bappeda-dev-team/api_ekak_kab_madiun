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
