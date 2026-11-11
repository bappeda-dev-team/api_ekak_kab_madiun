package renja

// IndikatorMatrixResponse — satu indikator matrix renja (ranwal / rankhir / penetapan).
// Hanya kode_indikator; tidak ada field id terpisah.
type IndikatorMatrixResponse struct {
	KodeIndikator   string `json:"kode_indikator"`
	Kode            string `json:"kode,omitempty"`
	KodeOpd         string `json:"kode_opd,omitempty"`
	Indikator       string `json:"indikator"`
	Tahun           string `json:"tahun"`
	TargetId        string `json:"target_id"`
	Target          string `json:"target"`
	Satuan          string `json:"satuan"`
	StatusTarget    bool   `json:"status_target_renja,omitempty"`
	TargetRanwal    string `json:"target_ranwal,omitempty"`
	SatuanRanwal    string `json:"satuan_ranwal,omitempty"`
	TargetRankhir   string `json:"target_rankhir,omitempty"`
	SatuanRankhir   string `json:"satuan_rankhir,omitempty"`
	TargetPenetapan string `json:"target_penetapan,omitempty"`
	SatuanPenetapan string `json:"satuan_penetapan,omitempty"`
}
