package domain

type RencanaAksi struct {
	Id               string
	RencanaKinerjaId string
	KodeOpd          string
	Urutan           int
	NamaRencanaAksi  string
	Pelaksanaan      []PelaksanaanRencanaAksi
}
