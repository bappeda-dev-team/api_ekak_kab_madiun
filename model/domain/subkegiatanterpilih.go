package domain

type SubKegiatanTerpilih struct {
	Id              string
	KodeSubKegiatan string
	RekinId         string
	SubkegiatanId   string
}

type SubKegiatanOpd struct {
	Id              int
	KodeSubKegiatan string
	KodeOpd         string
	Tahun           string
}
