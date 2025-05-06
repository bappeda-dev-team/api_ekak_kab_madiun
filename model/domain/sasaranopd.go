package domain

type SasaranOpd struct {
	Id                int
	IdPohon           int
	NamaPohon         string
	JenisPohon        string
	LevelPohon        int
	TahunPohon        string
	TahunAwalPeriode  string
	TahunAkhirPeriode string
	JenisPeriode      string
	Pelaksana         []PelaksanaPokin
	SasaranOpd        []SasaranOpdDetail
}

type SasaranOpdDetail struct {
	Id             int
	IdPohon        int
	NamaSasaranOpd string
	TahunAwal      string
	TahunAkhir     string
	JenisPeriode   string
	Indikator      []Indikator
}
