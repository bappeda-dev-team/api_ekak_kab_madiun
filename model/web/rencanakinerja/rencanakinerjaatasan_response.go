package rencanakinerja

type RekinAtasanResponse struct {
	KodeSubkegiatan string              `json:"kode_subkegiatan"`
	Subkegiatan     string              `json:"subkegiatan"`
	PaguSubkegiatan string              `json:"pagu_subkegiatan"`
	KodeKegiatan    string              `json:"kode_kegiatan"`
	Kegiatan        string              `json:"kegiatan"`
	PaguKegiatan    string              `json:"pagu_kegiatan"`
	RekinAtasan     []RekinAtasanDetail `json:"rekin_atasan"`
}

type RekinAtasanDetail struct {
	Id                   string `json:"id"`
	NamaRencanaKinerja   string `json:"nama_rencana_kinerja"`
	IdPohon              int    `json:"id_pohon"`
	Tahun                string `json:"tahun"`
	StatusRencanaKinerja string `json:"status_rencana_kinerja"`
	Catatan              string `json:"catatan"`
	KodeOpd              string `json:"kode_opd"`
	PegawaiId            string `json:"nip"`
	NamaPegawai          string `json:"nama_pegawai"`
	KodeProgram          string `json:"kode_program"`
	Program              string `json:"program"`
	PaguProgram          string `json:"pagu_program"`
}
