package domain

type CascadingOpd struct {
	KodeOpd      string
	NamaOpd      string
	Tahun        string
	BidangUrusan []BidangUrusanCascading
}

type BidangUrusanCascading struct {
	KodeBidangUrusan string
	NamaBidangUrusan string
	Program          []ProgramCascading
}

type ProgramCascading struct {
	KodeProgram string
	NamaProgram string
	Kegiatan    []KegiatanCascading
}

type KegiatanCascading struct {
	KodeKegiatan string
	NamaKegiatan string
	SubKegiatan  []SubKegiatanCascading
}

type SubKegiatanCascading struct {
	KodeSubKegiatan string
	NamaSubKegiatan string
	RencanaKinerja  []RencanaKinerjaCascading
}

type RencanaKinerjaCascading struct {
	Id                   string
	IdPohon              int
	NamaRencanaKinerja   string
	Tahun                string
	StatusRencanaKinerja string
	Catatan              string
	KodeOpd              string
	PegawaiId            string
	NamaPegawai          string
	Indikator            []IndikatorCascading
}

type IndikatorCascading struct {
	Id               string
	RencanaKinerjaId string
	Indikator        string
	Tahun            string
	Target           []TargetCascading
}

type TargetCascading struct {
	Id          string
	IndikatorId string
	Target      string
	Satuan      string
	Tahun       string
}
