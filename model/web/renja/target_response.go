package renja

// TargetLayerResponse — hasil upsert/GET target rankhir atau penetapan (tanpa id baris target).
type TargetLayerResponse struct {
	KodeIndikator string `json:"kode_indikator"`
	Tahun         string `json:"tahun"`
	Target        string `json:"target"`
	Satuan        string `json:"satuan"`
	Jenis         string `json:"jenis"`
}
