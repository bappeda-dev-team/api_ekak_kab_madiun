package rencanakinerja

type RekinAtasanResponse struct {
	RekinAtasan []RekinAtasanDetail `json:"rekin_atasan"`
}

type RekinAtasanDetail struct {
	Id                   string `json:"id"`
	NamaRencanaKinerja   string `json:"nama_rencana_kinerja"`
	IdPohon              int    `json:"id_pohon"`
	Tahun                string `json:"tahun"`
	StatusRencanaKinerja string `json:"status_rencana_kinerja"`
	Catatan              string `json:"catatan"`
	KodeOpd              string `json:"kode_opd"`
	PegawaiId            string `json:"pegawai_id"`
	NamaPegawai          string `json:"nama_pegawai"`
	NipPegawai           string `json:"nip_pegawai"`
}
