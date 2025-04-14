package pohonkinerja

type IndikatorResponse struct {
	Id            string           `json:"id_indikator"`
	IdPokin       string           `json:"id_pokin,omitempty"`
	IdRekin       string           `json:"id_rekin,omitempty"`
	Kode          string           `json:"kode,omitempty"`
	NamaIndikator string           `json:"nama_indikator"`
	Target        []TargetResponse `json:"targets"`
}

type TargetResponse struct {
	Id              string `json:"id_target"`
	IndikatorId     string `json:"indikator_id"`
	TargetIndikator string `json:"target"`
	SatuanIndikator string `json:"satuan"`
	TahunSasaran    string `json:"tahun,omitempty"`
}
