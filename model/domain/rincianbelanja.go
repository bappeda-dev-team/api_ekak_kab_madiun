package domain

type RincianBelanjaAsn struct {
	PegawaiId       string
	NamaPegawai     string
	KodeSubkegiatan string
	NamaSubkegiatan string
	TotalAnggaran   int
	RencanaKinerja  []RencanaKinerjaAsn
}

type RencanaKinerjaAsn struct {
	RencanaKinerja string
	RencanaAksi    []RincianBelanja
}

type RincianBelanja struct {
	Id        int
	RenaksiId string
	Renaksi   string
	Anggaran  int64
}
