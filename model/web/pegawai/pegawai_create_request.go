package pegawai

type PegawaiCreateRequest struct {
	NamaPegawai string `json:"nama_pegawai"`
	Nip         string `json:"nip"`
	KodeOpd     string `json:"kode_opd"`
}

type TambahJabatanRequest struct {
	Nip       string `json:"nip" validate:"required"`
	IdJabatan string `json:"id_jabatan" validate:"required"`
	Bulan     int    `json:"bulan" validate:"required"`
	Tahun     int    `json:"tahun" validate:"required"`
	KodeOpd   string `json:"kode_opd" validate:"required"`
}
