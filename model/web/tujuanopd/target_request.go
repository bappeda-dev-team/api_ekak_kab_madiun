package tujuanopd

type LayerTargetItemRequest struct {
	KodeIndikator string `json:"kode_indikator"`
	Tahun         string `json:"tahun"`
	Target        string `json:"target"`
	Satuan        string `json:"satuan"`
}
type LayerTargetBatchRequest struct {
	Targets []LayerTargetItemRequest `json:"targets"`
}
type LayerTargetUpdateItemRequest struct {
	Id     string `json:"id"`
	Target string `json:"target"`
	Satuan string `json:"satuan"`
}
type LayerTargetUpdateBatchRequest struct {
	Targets []LayerTargetUpdateItemRequest `json:"targets"`
}
