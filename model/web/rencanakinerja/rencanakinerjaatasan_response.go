package rencanakinerja

// type RekinAtasanResponse struct {
// 	KodeSubkegiatan string              `json:"kode_subkegiatan"`
// 	Subkegiatan     string              `json:"subkegiatan"`
// 	PaguSubkegiatan string              `json:"pagu_subkegiatan"`
// 	KodeKegiatan    string              `json:"kode_kegiatan"`
// 	Kegiatan        string              `json:"kegiatan"`
// 	PaguKegiatan    string              `json:"pagu_kegiatan"`
// 	RekinAtasan     []RekinAtasanDetail `json:"rekin_atasan"`
// }

// type RekinAtasanDetail struct {
// 	Id                   string `json:"id"`
// 	NamaRencanaKinerja   string `json:"nama_rencana_kinerja"`
// 	IdPohon              int    `json:"id_pohon"`
// 	Tahun                string `json:"tahun"`
// 	StatusRencanaKinerja string `json:"status_rencana_kinerja"`
// 	Catatan              string `json:"catatan"`
// 	KodeOpd              string `json:"kode_opd"`
// 	PegawaiId            string `json:"nip"`
// 	NamaPegawai          string `json:"nama_pegawai"`
// 	KodeProgram          string `json:"kode_program"`
// 	Program              string `json:"program"`
// 	PaguProgram          string `json:"pagu_program"`
// }

type RekinAtasanResponse struct {
	PokinParent       PokinParentInfo             `json:"pokin_parent"`
	RekinAtasan       []RekinAtasanDetail         `json:"rekin_atasan"`
	Program           []ProgramAtasanResponse     `json:"program,omitempty"`
	Kegiatan          []KegiatanAtasanResponse    `json:"kegiatan,omitempty"`
	Subkegiatan       []SubKegiatanAtasanResponse `json:"sub_kegiatan,omitempty"`
	PaguAnggaranTotal int64                       `json:"pagu_anggaran_total"`
}

type PokinParentInfo struct {
	Id         int    `json:"id"`
	NamaPohon  string `json:"nama_pohon"`
	LevelPohon int    `json:"level_pohon"`
	KodeOpd    string `json:"kode_opd"`
	NamaOpd    string `json:"nama_opd"`
}

type RekinAtasanDetail struct {
	Id                   string `json:"id_rencana_kinerja"`
	NamaRencanaKinerja   string `json:"nama_rencana_kinerja"`
	IdPohon              int    `json:"id_pohon"`
	Tahun                string `json:"tahun"`
	StatusRencanaKinerja string `json:"status_rencana_kinerja"`
	Catatan              string `json:"catatan"`
	KodeOpd              string `json:"kode_opd"`
	PegawaiId            string `json:"pegawai_id"`
	NamaPegawai          string `json:"nama_pegawai"`
}

type ProgramAtasanResponse struct {
	KodeProgram string `json:"kode_program"`
	NamaProgram string `json:"nama_program"`
	PaguProgram int64  `json:"pagu_program"`
}

type KegiatanAtasanResponse struct {
	KodeKegiatan string `json:"kode_kegiatan"`
	NamaKegiatan string `json:"nama_kegiatan"`
	PaguKegiatan int64  `json:"pagu_kegiatan"`
}

type SubKegiatanAtasanResponse struct {
	KodeSubkegiatan string `json:"kode_subkegiatan"`
	NamaSubkegiatan string `json:"nama_subkegiatan"`
	PaguSubkegiatan int64  `json:"pagu_subkegiatan"`
}
