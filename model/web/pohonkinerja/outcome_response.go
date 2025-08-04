package pohonkinerja

type OutcomeResponse struct {
	Tahun   string                   `json:"tahun,omitempty"`
	Tematik []OutcomeTematikResponse `json:"tematiks"`
}

type OutcomeTematikResponse struct {
	Id         int                 `json:"id"`
	Parent     *int                `json:"parent"`
	Tema       string              `json:"tema"`
	JenisPohon string              `json:"jenis_pohon"`
	LevelPohon int                 `json:"level_pohon"`
	Indikators []IndikatorResponse `json:"indikator"`
	Child      []interface{}       `json:"childs,omitempty"`
}

type OutcomeSubtematikResponse struct {
	// Outcome    []outcome.OutcomeResponse `json:"outcome"`
	Id         int                 `json:"id"`
	Parent     int                 `json:"parent"`
	Tema       string              `json:"tema"`
	JenisPohon string              `json:"jenis_pohon"`
	LevelPohon int                 `json:"level_pohon"`
	Indikators []IndikatorResponse `json:"indikator"`
	Child      []interface{}       `json:"childs,omitempty"`
}
