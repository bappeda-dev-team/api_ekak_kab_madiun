package domain

type LaporanReviewTujuanOpd struct {
	Id               int
	KodeOpd          string
	NamaOpd          string
	KodeBidangUrusan string
	NamaBidangUrusan string
	Tujuan           string
	Review           []LaporanReviewTujuanDetail
}

type LaporanReviewTujuanDetail struct {
	Id             int
	IdTujuanOpd    int
	IdPohonKinerja int
	Review         string
	Keterangan     string
	Catatan        string
	CreatedBy      string
	NamaPegawai    string
	Tahun          int
}