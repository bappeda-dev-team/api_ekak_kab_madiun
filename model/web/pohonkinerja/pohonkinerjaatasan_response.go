package pohonkinerja

type PokinAtasanResponse struct {
	Id        int               `json:"id"`
	NamaPohon string            `json:"nama_pohon"`
	Pegawai   []PegawaiResponse `json:"pegawai"`
}

type PegawaiResponse struct {
	IdPegawai   string `json:"id_pegawai"`
	NipPegawai  string `json:"nip_pegawai"`
	NamaPegawai string `json:"nama_pegawai"`
}
