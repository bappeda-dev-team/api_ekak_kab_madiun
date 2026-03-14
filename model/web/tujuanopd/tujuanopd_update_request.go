package tujuanopd

type TujuanOpdUpdateRequest struct {
	Id               int                      `json:"id"`
	KodeOpd          string                   `json:"kode_opd"`
	KodeBidangUrusan string                   `json:"kode_bidang_urusan"`
	Tujuan           string                   `json:"tujuan"`
	PeriodeId        int                      `json:"periode_id"`
	TahunAwal        string                   `json:"tahun_awal"`
	TahunAkhir       string                   `json:"tahun_akhir"`
	JenisPeriode     string                   `json:"jenis_periode"`
	Indikator        []IndikatorUpdateRequest `json:"indikator"`
}

type IndikatorUpdateRequest struct {
	Id                  string                `json:"id"`
	IdTujuanOpd         string                `json:"id_tujuan_opd"`
	KodeIndikator       string                `json:"kode_indikator"`
	Jenis               string                `json:"jenis"`
	DefinisiOperasional string                `json:"definisi_operasional"`
	Indikator           string                `json:"indikator"`
	RumusPerhitungan    string                `json:"rumus_perhitungan"`
	SumberData          string                `json:"sumber_data"`
	Target              []TargetUpdateRequest `json:"target"`
}

type TargetUpdateRequest struct {
	Id          string `json:"id"`
	IndikatorId string `json:"indikator_id"`
	Tahun       string `json:"tahun"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}
