package domainmaster

type JabatanPegawai struct {
	Id        string
	IdJabatan string
	IdPegawai string
	Status    string
	IsActive  bool //default true
	Bulan     string
	Tahun     string
	Pangkat   string
	Golongan  string
	KodeOpd   string
}
