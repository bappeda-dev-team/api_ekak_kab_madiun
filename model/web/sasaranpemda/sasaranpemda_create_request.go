package sasaranpemda

type SasaranPemdaCreateRequest struct {
	SubtemaId     int                      `json:"subtema_id"`
	TujuanPemdaId int                      `json:"tujuan_pemda_id"`
	SasaranPemda  string                   `json:"sasaran_pemda"`
	PeriodeId     int                      `json:"periode_id"`
	Indikator     []IndikatorCreateRequest `json:"indikator"`
}
type IndikatorCreateRequest struct {
	Indikator        string                `json:"indikator"`
	RumusPerhitungan string                `json:"rumus_perhitungan"`
	SumberData       string                `json:"sumber_data"`
	Target           []TargetCreateRequest `json:"target"`
}
type TargetCreateRequest struct {
	Target TargetInput `json:"target"`
	Satuan string      `json:"satuan"`
	Tahun  string      `json:"tahun"`
}

// Layer create (rankhir / penetapan) — pakai kode_indikator
type LayerTargetItemRequest struct {
	KodeIndikator string      `json:"kode_indikator"`
	Tahun         string      `json:"tahun"`
	Target        TargetInput `json:"target"`
	Satuan        string      `json:"satuan"`
}
type LayerTargetBatchRequest struct {
	Targets []LayerTargetItemRequest `json:"targets"`
}
