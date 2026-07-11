package sasaranpemda

type SasaranPemdaResponse struct {
	Id            int                 `json:"id"`
	SubtemaId     int                 `json:"subtema_id,omitempty"`
	NamaSubtema   string              `json:"nama_subtema,omitempty"`
	TujuanPemdaId int                 `json:"tujuan_pemda_id,omitempty"`
	TujuanPemda   string              `json:"tujuan_pemda,omitempty"`
	SasaranPemda  string              `json:"sasaran_pemda"`
	JenisPohon    string              `json:"jenis_pohon,omitempty"`
	Periode       PeriodeResponse     `json:"periode"`
	Indikator     []IndikatorResponse `json:"indikator"`
}
type IndikatorResponse struct {
	Id               int              `json:"id"`             // id DB auto-increment
	KodeIndikator    string           `json:"kode_indikator"` // IND-SAS-PMD-YYYY-xxxxx
	Indikator        string           `json:"indikator"`
	RumusPerhitungan string           `json:"rumus_perhitungan"`
	SumberData       string           `json:"sumber_data"`
	Target           []TargetResponse `json:"target"`
}
type TargetResponse struct {
	Id     int           `json:"id"` // id DB auto-increment (0 = slot kosong)
	Target TargetDisplay `json:"target"`
	Satuan string        `json:"satuan"`
	Tahun  string        `json:"tahun"`
	Jenis  string        `json:"jenis,omitempty"`
}
type PeriodeResponse struct {
	Id           int    `json:"id"`
	TahunAwal    string `json:"tahun_awal"`
	TahunAkhir   string `json:"tahun_akhir"`
	JenisPeriode string `json:"jenis_periode"`
}

// ── FindAllWithPokin ─────────────────────────────────────────────
type TematikResponse struct {
	TematikId   int                  `json:"tematik_id"`
	NamaTematik string               `json:"nama_tematik"`
	Tahun       string               `json:"tahun"`
	IsLock      bool                 `json:"is_lock"`
	Subtematik  []SubtematikResponse `json:"subtematik"`
}
type SubtematikResponse struct {
	SubtematikId   int                             `json:"subtematik_id"`
	NamaSubtematik string                          `json:"nama_subtematik"`
	JenisPohon     string                          `json:"jenis_pohon"`
	LevelPohon     int                             `json:"level_pohon"`
	Tahun          string                          `json:"tahun"`
	IsActive       bool                            `json:"is_active"`
	SasaranPemda   []SasaranPemdaWithPokinResponse `json:"sasaran_pemda"`
}
type SasaranPemdaWithPokinResponse struct {
	IdSasaranPemda int                           `json:"id_sasaran_pemda"`
	SasaranPemda   string                        `json:"sasaran_pemda"`
	Periode        PeriodeResponse               `json:"periode"`
	Indikator      []IndikatorSubtematikResponse `json:"indikator"`
}
type IndikatorSubtematikResponse struct {
	Id               int              `json:"id"`
	KodeIndikator    string           `json:"kode_indikator"`
	Indikator        string           `json:"indikator"`
	RumusPerhitungan string           `json:"rumus_perhitungan"`
	SumberData       string           `json:"sumber_data"`
	Target           []TargetResponse `json:"target"`
}

// ── Dual Response Rankhir ────────────────────────────────────────
type IndikatorRankhirDualResponse struct {
	Id               int              `json:"id"`
	KodeIndikator    string           `json:"kode_indikator"`
	Indikator        string           `json:"indikator"`
	RumusPerhitungan string           `json:"rumus_perhitungan"`
	SumberData       string           `json:"sumber_data"`
	TargetRanwal     []TargetResponse `json:"target_ranwal"`
	TargetRankhir    []TargetResponse `json:"target_rankhir"`
}
type SasaranPemdaRankhirDualResponse struct {
	Id           int                            `json:"id"`
	SasaranPemda string                         `json:"sasaran_pemda"`
	Periode      PeriodeResponse                `json:"periode"`
	Indikator    []IndikatorRankhirDualResponse `json:"indikator"`
}

// ── Dual Response Penetapan ──────────────────────────────────────
type IndikatorPenetapanDualResponse struct {
	Id               int              `json:"id"`
	KodeIndikator    string           `json:"kode_indikator"`
	Indikator        string           `json:"indikator"`
	RumusPerhitungan string           `json:"rumus_perhitungan"`
	SumberData       string           `json:"sumber_data"`
	TargetRankhir    []TargetResponse `json:"target_rankhir"`
	TargetPenetapan  []TargetResponse `json:"target_penetapan"`
}
type SasaranPemdaPenetapanDualResponse struct {
	Id           int                              `json:"id"`
	SasaranPemda string                           `json:"sasaran_pemda"`
	Periode      PeriodeResponse                  `json:"periode"`
	Indikator    []IndikatorPenetapanDualResponse `json:"indikator"`
}
