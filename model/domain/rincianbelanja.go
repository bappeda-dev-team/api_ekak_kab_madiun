package domain

type RincianBelanjaAsn struct {
	PegawaiId       string
	KodeOpd         string
	NamaPegawai     string
	KodeSubkegiatan string
	NamaSubkegiatan string
	Indikator       []Indikator
	TotalAnggaran   int
	RencanaKinerja  []RencanaKinerjaAsn
}

type RencanaKinerjaAsn struct {
	RencanaKinerjaId string
	RencanaKinerja   string
	PegawaiId        string
	NamaPegawai      string
	Indikator        []Indikator
	RencanaAksi      []RincianBelanja
}

type RincianBelanja struct {
	Id        int
	RenaksiId string
	Renaksi   string
	Anggaran  int64
}
