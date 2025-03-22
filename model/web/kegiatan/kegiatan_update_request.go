package kegiatan

type KegiatanUpdateRequest struct {
	Id           string `json:"id"`
	NamaKegiatan string `json:"nama_kegiatan"`
	KodeKegiatan string `json:"kode_kegiatan"`
	Indikator    []IndikatorUpdateRequest
}

type IndikatorUpdateRequest struct {
	Id         string `json:"id"`
	KegiatanId string `json:"kegiatan_id"`
	Indikator  string `json:"indikator"`
	Tahun      string `json:"tahun"`
	Target     []TargetUpdateRequest
}

type TargetUpdateRequest struct {
	Id          string `json:"id"`
	IndikatorId string `json:"indikator_id"`
	Tahun       string `json:"tahun"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}
