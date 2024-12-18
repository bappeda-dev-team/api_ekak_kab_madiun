package tujuanopd

type TujuanOpdUpdateRequest struct {
	Id               int                      `json:"id"`
	KodeOpd          string                   `json:"kode_opd"`
	Tujuan           string                   `json:"tujuan"`
	RumusPerhitungan string                   `json:"rumus_perhitungan"`
	SumberData       string                   `json:"sumber_data"`
	TahunAwal        string                   `json:"tahun_awal"`
	TahunAkhir       string                   `json:"tahun_akhir"`
	Indikator        []IndikatorUpdateRequest `json:"indikator"`
}

type IndikatorUpdateRequest struct {
	Id          string                `json:"id"`
	IdTujuanOpd string                `json:"id_tujuan_opd"`
	Indikator   string                `json:"indikator"`
	Target      []TargetUpdateRequest `json:"target"`
}

type TargetUpdateRequest struct {
	Id          string `json:"id"`
	IndikatorId string `json:"indikator_id"`
	Tahun       string `json:"tahun"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}
