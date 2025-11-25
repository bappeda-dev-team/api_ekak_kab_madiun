package domain

type DetailRekins struct {
	Id                 string
	IdPohon            int
	LevelPohon         int
	Parent             int
	NamaRencanaKinerja string
	Tahun              string
	PegawaiId          string
	NamaPegawai        string
	KodeOpd            string
	KodeSubKegiatan    string
	Indikator          []Indikator
}
