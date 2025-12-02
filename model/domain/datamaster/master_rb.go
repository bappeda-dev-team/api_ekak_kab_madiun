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
