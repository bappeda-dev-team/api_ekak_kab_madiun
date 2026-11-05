package domain

type RencanaAksi struct {
	Id                 string
	RencanaKinerjaId   string
	NamaRencanaKinerja string
	PegawaiId          string
	Tahun              string
	KodeOpd            string
	Urutan             int
	NamaRencanaAksi    string
	Pelaksanaan        []PelaksanaanRencanaAksi
}
