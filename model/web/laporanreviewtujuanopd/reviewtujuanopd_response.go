package laporanreviewtujuanopd

type LaporanReviewTujuanOpd struct {
	Id               int                 `json:"id"`
	KodeOpd          string              `json:"kode_opd"`
	NamaOpd          string              `json:"nama_opd"`
	KodeBidangUrusan string              `json:"kode_bidang_urusan"`
	NamaBidangUrusan string              `json:"nama_bidang_urusan"`
	Tujuan           string              `json:"tujuan"`
	Review           LaporanReviewTujuan `json:"review"`
}

// type LaporanReviewTujuanDetail struct {
// 	Id             int    `json:"id"`
// 	IdTujuanOpd    int    `json:"id_tujuan_opd"`
// 	IdPohonKinerja int    `json:"id_pohon_kinerja"`
// 	Review         string `json:"review"`
// 	Keterangan     string `json:"keterangan"`
// 	Catatan        string `json:"catatan"`
// 	CreatedBy      string `json:"created_by"`
// 	NamaPegawai    string `json:"nama_pegawai"`
// 	Tahun          int    `json:"tahun"`
// }

type LaporanReviewTujuan struct {
	Validasi    *LaporanReviewTujuanDetail `json:"validasi"`
	Kesesuaian  *LaporanReviewTujuanDetail `json:"kesesuaian"`
	Kelengkapan *LaporanReviewTujuanDetail `json:"kelengkapan"`
	Catatan     *LaporanReviewTujuanDetail `json:"catatan"`
}

type LaporanReviewTujuanDetail struct {
	Id             *int    `json:"id"`
	IdTujuanOpd    *int    `json:"id_tujuan_opd"`
	IdPohonKinerja *int    `json:"id_pohon_kinerja"`
	Review         string  `json:"review"`
	Keterangan     string  `json:"keterangan"`
	Catatan        string  `json:"catatan"`
	CreatedBy      *string `json:"created_by"`
	NamaPegawai    *string `json:"nama_pegawai"`
	Tahun          *int    `json:"tahun"`
}