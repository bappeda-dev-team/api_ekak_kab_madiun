package tujuanopd

type LayerTargetItemRequest struct {
	KodeIndikator string `json:"kode_indikator" validate:"required"`
	Tahun         string `json:"tahun" validate:"required"`
	Target        string `json:"target" validate:"required"`
	Satuan        string `json:"satuan" validate:"required"`
}
type LayerTargetBatchRequest struct {
	Targets []LayerTargetItemRequest `json:"targets"`
}
type LayerTargetUpdateItemRequest struct {
	Id     string `json:"id" validate:"required"`
	Target string `json:"target" validate:"required"`
	Satuan string `json:"satuan" validate:"required"`
}
type LayerTargetUpdateBatchRequest struct {
	Targets []LayerTargetUpdateItemRequest `json:"targets"`
}
