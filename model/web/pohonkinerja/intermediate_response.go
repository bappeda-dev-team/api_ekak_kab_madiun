package pohonkinerja

type IntermediateResponse struct {
	Tahun        string                           `json:"tahun"`
	Intermediate []IntermediateSubtematikResponse `json:"intermediate"`
}

type IntermediateSubtematikResponse struct {
	Id         int                 `json:"id"`
	Parent     int                 `json:"parent"`
	Tema       string              `json:"tema"`
	JenisPohon string              `json:"jenis_pohon"`
	LevelPohon int                 `json:"level_pohon"`
	Indikators []IndikatorResponse `json:"indikator"`
	Child      []interface{}       `json:"childs,omitempty"`
}

type IntermediateSubSubtematikResponse struct {
	Id         int                 `json:"id"`
	Parent     int                 `json:"parent"`
	Tema       string              `json:"tema"`
	JenisPohon string              `json:"jenis_pohon"`
	LevelPohon int                 `json:"level_pohon"`
	Indikators []IndikatorResponse `json:"indikator"`
	Child      []interface{}       `json:"childs,omitempty"`
}

type IntermediateStrategicPemdaResponse struct {
	Id         int                 `json:"id"`
	Parent     int                 `json:"parent"`
	Tema       string              `json:"tema"`
	JenisPohon string              `json:"jenis_pohon"`
	LevelPohon int                 `json:"level_pohon"`
	Indikators []IndikatorResponse `json:"indikator"`
	Child      []interface{}       `json:"childs,omitempty"`
}
