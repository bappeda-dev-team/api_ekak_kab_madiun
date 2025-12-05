package datamaster

type MasterRB struct {
	Id             int
	JenisRB        string
	KegiatanUtama  string
	Keterangan     string
	Indikator      []IndikatorRB
	TahunBaseline  int
	TahunNext      int
	LastUpdatedBy  int
	CurrentVersion int
	RencanaAksis  []RencanaAksiRB
}

type IndikatorRB struct {
	IdRB        int
	IdIndikator string
	Indikator   string
	TargetRB    []TargetRB
}

type TargetRB struct {
	IdTarget          string
	IdIndikator       string
	TahunBaseline     int
	TargetBaseline    int
	RealisasiBaseline float32
	SatuanBaseline    string
	TahunNext         int
	TargetNext        int
	SatuanNext        string
}

type PokinIdRBTagging struct {
	IdTagging int
	KodeRB int
	KegiatanUtama string
	IdPokin int
	NamaTagging string
	NamaPohon string
	KodeOpd string
	JenisPohon string
}

type RencanaAksiRB struct {
	RencanaAksi     string
	IndikatorOutput []IndikatorRencanaAksiRB
	Anggaran        int
	Realisasi       int
	OpdKoordinator  string
	NipPelaksana    string
	NamaPelaksana   string
	OpdCrosscutting []OpdCrosscutting
}

type IndikatorRencanaAksiRB struct {
	Indikator       string
	TargetIndikator []TargetIndikatorRencanaAksiRB
}

type TargetIndikatorRencanaAksiRB struct {
	Target    int
	Realisasi int
	Satuan    string
	Capaian   int
}

type OpdCrosscutting struct {
	KodeOpd       string
	NamaOpd       string
	NipPelaksana  string
	NamaPelaksana string
}
