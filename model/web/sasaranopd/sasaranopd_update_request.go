package sasaranopd

type SasaranOpdUpdateRequest struct {
	IdSasaranOpd int                      `json:"id"`
	NamaSasaran  string                   `json:"nama_sasaran" validate:"required"`
	TahunAwal    string                   `json:"tahun_awal" validate:"required"`
	TahunAkhir   string                   `json:"tahun_akhir" validate:"required"`
	JenisPeriode string                   `json:"jenis_periode" validate:"required"`
	Indikator    []IndikatorUpdateRequest `json:"indikator"`
}

type IndikatorUpdateRequest struct {
	Id               string                `json:"id"`
	Indikator        string                `json:"indikator"`
	RumusPerhitungan string                `json:"rumus_perhitungan"`
	SumberData       string                `json:"sumber_data"`
	Target           []TargetUpdateRequest `json:"target"`
}

type TargetUpdateRequest struct {
	Id     string `json:"id"`
	Tahun  string `json:"tahun"`
	Target string `json:"target"`
	Satuan string `json:"satuan"`
}
