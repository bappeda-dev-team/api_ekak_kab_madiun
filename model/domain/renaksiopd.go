package domain

type RencanaAksiOpd struct {
	Id                 int
	RekinId            string
	SasaranOpdId       int
	NamaSasaranOpd     string
	TahunRenaksi       string
	Tw1                int
	Tw2                int
	Tw3                int
	Tw4                int
	Keterangan         *string
	RencanaKinerja     []RencanaKinerjaOpd
	NamaRencanaKinerja string
	SasaranOpd         SasaranOpdDetailRenaksi
}

type RencanaKinerjaOpd struct {
	Id                 int
	RekinId            string
	NamaRencanaKinerja string
	NipPegawai         string
	NamaPegawai        string
	KodeOpd            string
	Tw1                int
	Tw2                int
	Tw3                int
	Tw4                int
	Keterangan         *string
	TotalAnggaran      int64
	SubKegiatan        []SubKegiatanOpdRenaksi
}

type SubKegiatanOpdRenaksi struct {
	KodeSubKegiatan string
	NamaSubKegiatan string
	Indikator       []IndikatorSubKegiatanOpd
}

type IndikatorSubKegiatanOpd struct {
	Id           string
	Indikator    string
	Target       string
	Satuan       string
	PaguAnggaran int64
}

type RencanaAksiOpdDomainRequest struct {
	SasaranOpdId int
	RekinId      string
	TahunRenaksi string
	Tw1          int
	Tw2          int
	Tw3          int
	Tw4          int
	Keterangan   *string
}

type SasaranOpdDetailRenaksi struct {
	Id             int
	NamaSasaranOpd string
	TahunAwal      string
	TahunAkhir     string
	JenisPeriode   string
	Indikator      []IndikatorSasaranOpd
}

type IndikatorSasaranOpd struct {
	Id               string
	SasaranOpdId     int
	Indikator        string
	RumusPerhitungan string
	SumberData       string
	Target           TargetRenaksiOpd
}

type TargetRenaksiOpd struct {
	Id          string
	IndikatorId string
	Tahun       string
	Target      string
	Satuan      string
}
