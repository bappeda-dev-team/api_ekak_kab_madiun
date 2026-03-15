package domain

type TujuanOpd struct {
	Id               int
	KodeOpd          string
	NamaOpd          string
	KodeBidangUrusan string
	NamaBidangUrusan string
	Tujuan           string
	PeriodeId        Periode
	TahunAwal        string
	TahunAkhir       string
	JenisPeriode     string
	Indikator        []Indikator
}
