package programkegiatan

type PengampuSubKegiatanRequest struct {
	KodeSubKegiatan string `json:"kode_subkegiatan" validate:"required"`
	KodeOpd         string `json:"kode_opd" validate:"required"`
	TahunAwal       string `json:"tahun_awal" validate:"required"`
	TahunAkhir      string `json:"tahun_akhir" validate:"required"`
}

type PengampuSubKegiatanResponse struct {
	KodeSubKegiatan string             `json:"kode_subkegiatan"`
	NamaSubKegiatan string             `json:"nama_subkegiatan"`
	KodeOpd         string             `json:"kode_opd"`
	TahunAwal       string             `json:"tahun_awal"`
	TahunAkhir      string             `json:"tahun_akhir"`
	Pengampu        []PengampuResponse `json:"pengampu"`
}

type PengampuResponse struct {
	Nip              string `json:"nip"`
	NamaPegawai      string `json:"nama_pegawai"`
	RencanaKinerjaId string `json:"rencana_kinerja_id"`
	RencanaKinerja   string `json:"rencana_kinerja"`
	Tahun            string `json:"tahun"`
}
