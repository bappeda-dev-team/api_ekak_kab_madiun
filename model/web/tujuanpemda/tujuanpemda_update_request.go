package tujuanpemda

type TujuanPemdaUpdateRequest struct {
	Id          int                      `json:"id"`
	TujuanPemda string                   `json:"tujuan_pemda"`
	TematikId   int                      `json:"tema_id"`
	PeriodeId   int                      `json:"periode_id"`
	Indikator   []IndikatorUpdateRequest `json:"indikator"`
}

type IndikatorUpdateRequest struct {
	Id               string                `json:"id"`
	TujuanPemdaId    string                `json:"tujuan_pemda_id"`
	Indikator        string                `json:"indikator"`
	RumusPerhitungan string                `json:"rumus_perhitungan"`
	SumberData       string                `json:"sumber_data"`
	Target           []TargetUpdateRequest `json:"target"`
}

type TargetUpdateRequest struct {
	Id     string `json:"id"`
	Target string `json:"target"`
	Satuan string `json:"satuan"`
	Tahun  string `json:"tahun"`
}
