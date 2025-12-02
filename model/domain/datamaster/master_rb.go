package datamaster

type MasterRB struct {
	Id            int
	JenisRB       string
	KegiatanUtama string
	Keterangan    string
	Indikator     []IndikatorRB
}

type IndikatorRB struct {
	IdRB        int
	IdIndikator string
	Indikator   string
	TargetRB    []TargetRB
}

type TargetRB struct {
	IdIndikator       string
	IdTarget          string
	TahunBaseline     int
	TargetBaseline    int
	RealisasiBaseline float32
	SatuanBaseline    string
	TahunNext         int
	TargetNext        int
	SatuanNext        string
}
