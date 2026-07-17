package sasaranpemda

type SasaranPemdaUpdateRequest struct {
	Id            int    `json:"id"`
	TujuanPemdaId int    `json:"tujuan_pemda_id"`
	SubtemaId     int    `json:"subtema_id"`
	SasaranPemda  string `json:"sasaran_pemda"`
	// PeriodeId     int                      `json:"periode_id"`
	Indikator []IndikatorUpdateRequest `json:"indikator"`
}
type IndikatorUpdateRequest struct {
	IdIndikator      int                   `json:"id_indikator"`
	KodeIndikator    string                `json:"kode_indikator"`
	Indikator        string                `json:"indikator"`
	RumusPerhitungan string                `json:"rumus_perhitungan"`
	SumberData       string                `json:"sumber_data"`
	Target           []TargetUpdateRequest `json:"target"`
}
type TargetUpdateRequest struct {
	Id     int         `json:"id"`
	Target TargetInput `json:"target"`
	Satuan string      `json:"satuan"`
	Tahun  string      `json:"tahun"`
}

// Layer update (rankhir / penetapan)
type LayerTargetUpdateItemRequest struct {
	Id     int         `json:"id"`
	Target TargetInput `json:"target"`
	Satuan string      `json:"satuan"`
}
type LayerTargetUpdateBatchRequest struct {
	Targets []LayerTargetUpdateItemRequest `json:"targets"`
}
