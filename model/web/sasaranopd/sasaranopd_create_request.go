package sasaranopd

type SasaranOpdCreateRequest struct {
	IdPohon      int                      `json:"id_pohon"`
	NamaSasaran  string                   `json:"nama_sasaran" validate:"required"`
	IdTujuanOpd  int                      `json:"id_tujuan_opd" validate:"required"`
	TahunAwal    string                   `json:"tahun_awal" validate:"required"`
	TahunAkhir   string                   `json:"tahun_akhir" validate:"required"`
	JenisPeriode string                   `json:"jenis_periode" validate:"required"`
	Indikator    []IndikatorCreateRequest `json:"indikator" `
}

type IndikatorCreateRequest struct {
	Id               string                `json:"id"`
	Indikator        string                `json:"indikator"`
	RumusPerhitungan string                `json:"rumus_perhitungan"`
	SumberData       string                `json:"sumber_data"`
	Target           []TargetCreateRequest `json:"target"`
}

type TargetCreateRequest struct {
	Tahun  string `json:"tahun"`
	Target string `json:"target"`
	Satuan string `json:"satuan"`
}
