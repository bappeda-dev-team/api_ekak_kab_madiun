package kegiatan

type KegiatanCreateRequest struct {
	NamaKegiatan string `json:"nama_kegiatan"`
	KodeKegiatan string `json:"kode_kegiatan"`
	Indikator    []IndikatorCreateRequest
}

type IndikatorCreateRequest struct {
	Indikator string `json:"indikator"`
	Tahun     string `json:"tahun"`
	Target    []TargetCreateRequest
}

type TargetCreateRequest struct {
	IndikatorId string `json:"indikator_id"`
	Tahun       string `json:"tahun"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}
