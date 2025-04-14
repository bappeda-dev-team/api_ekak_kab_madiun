package rincianbelanja

type RincianBelanjaAsnResponse struct {
	PegawaiId       string                   `json:"pegawai_id"`
	NamaPegawai     string                   `json:"nama_pegawai"`
	KodeSubkegiatan string                   `json:"kode_subkegiatan"`
	NamaSubkegiatan string                   `json:"nama_subkegiatan"`
	TotalAnggaran   int                      `json:"total_anggaran"`
	RincianBelanja  []RincianBelanjaResponse `json:"rincian_belanja"`
}

type RincianBelanjaResponse struct {
	RencanaKinerja string                `json:"rencana_kinerja"`
	RencanaAksi    []RencanaAksiResponse `json:"rencana_aksi"`
}

type RencanaAksiResponse struct {
	RenaksiId string `json:"renaksi_id"`
	Renaksi   string `json:"renaksi"`
	Anggaran  int    `json:"anggaran"`
}
