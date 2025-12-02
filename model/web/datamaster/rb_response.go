package datamaster

type RBResponse struct {
	IdRB          int           `json:"id"`
	JenisRB       string        `json:"jenis_rb"`
	KegiatanUtama string        `json:"kegiatan_utama"`
	Keterangan    string        `json:"keterangan"`
	Indikator     []IndikatorRB `json:"indikator"`
	TahunBaseline int           `json:"tahun_baseline"`
	TahunNext     int           `json:"tahun_next"`
}

type IndikatorRB struct {
	IdIndikator string     `json:"id"`
	IdRB        int        `json:"id_rb"`
	Indikator   string     `json:"indikator"`
	TargetRB    []TargetRB `json:"target"`
}

type TargetRB struct {
	IdTarget          string  `json:"id"`
	IdIndikator       string  `json:"id_indikator"`
	TahunBaseline     int     `json:"tahun_baseline"`
	TargetBaseline    int     `json:"target_baseline,string"`
	RealisasiBaseline float32 `json:"realisasi_baseline,string"`
	SatuanBaseline    string  `json:"satuan_baseline"`
	TahunNext         int     `json:"tahun_next"`
	TargetNext        int     `json:"target_next,string"`
	SatuanNext        string  `json:"satuan_next,string"`
}
