package subkegiatan

type SubKegiatanUpdateRequest struct {
	Id              string                   `json:"id"`
	NamaSubKegiatan string                   `json:"nama_subkegiatan"`
	KodeOpd         string                   `json:"kode_opd"`
	Tahun           string                   `json:"tahun"`
	Indikator       []IndikatorUpdateRequest `json:"indikator"`
}

type IndikatorUpdateRequest struct {
	Id               string                `json:"id_indikator"`
	RencanaKinerjaId string                `json:"rencana_kinerja_id"`
	NamaIndikator    string                `json:"nama_indikator"`
	Target           []TargetUpdateRequest `json:"targets"`
}

type TargetUpdateRequest struct {
	Id              string `json:"id_target"`
	IndikatorId     string `json:"indikator_id"`
	TargetIndikator string `json:"target"`
	SatuanIndikator string `json:"satuan"`
}

type SubKegiatanOpdUpdateRequest struct {
	Id              int    `json:"id"`
	KodeSubkegiatan string `json:"kode_subkegiatan" validate:"required"`
	KodeOpd         string `json:"kode_opd" validate:"required"`
	Tahun           string `json:"tahun" validate:"required"`
}
