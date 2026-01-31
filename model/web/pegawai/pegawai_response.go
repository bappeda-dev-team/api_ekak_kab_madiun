package pegawai

type PegawaiResponse struct {
	Id          string `json:"id"`
	NamaPegawai string `json:"nama_pegawai"`
	Nip         string `json:"nip"`
	KodeOpd     string `json:"kode_opd"`
	NamaOpd     string `json:"nama_opd"`
	NamaJabatan string `json:"nama_jabatan"`
}
