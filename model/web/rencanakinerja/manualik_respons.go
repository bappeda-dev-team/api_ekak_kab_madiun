package rencanakinerja

type ManualIKResponse struct {
	Id                  int           `json:"id,omitempty"`
	IndikatorId         string        `json:"indikator_id,omitempty"`
	DataIndikator       RekinResponse `json:"rencana_kinerja,omitempty"`
	ParentPohon         int           `json:"parent,omitempty"`
	NamaParent          string        `json:"nama_parent,omitempty"`
	Perspektif          string        `json:"perspektif,omitempty"`
	TujuanRekin         string        `json:"tujuan_rekin,omitempty"`
	Definisi            string        `json:"definisi,omitempty"`
	KeyActivities       string        `json:"key_activities,omitempty"`
	Formula             string        `json:"formula,omitempty"`
	JenisIndikator      string        `json:"jenis_indikator,omitempty"`
	OutputData          OutputData    `json:"output_data,omitempty"`
	UnitPenanggungJawab string        `json:"unit_penanggung_jawab,omitempty"`
	UnitPenyediaData    string        `json:"unit_penyedia_data,omitempty"`
	SumberData          string        `json:"sumber_data,omitempty"`
	JangkaWaktuAwal     string        `json:"jangka_waktu_awal,omitempty"`
	JangkaWaktuAkhir    string        `json:"jangka_waktu_akhir,omitempty"`
	PeriodePelaporan    string        `json:"periode_pelaporan,omitempty"`
}

type OutputData struct {
	Kinerja  bool `json:"kinerja,omitempty"`
	Penduduk bool `json:"penduduk,omitempty"`
	Spatial  bool `json:"spatial,omitempty"`
}

type RekinResponse struct {
	RencanaKinerja string              `json:"Nama_rencana_kinerja,omitempty"`
	Indikator      []IndikatorResponse `json:"indikator,omitempty"`
}

type DataOutput struct {
	OutputData []string `json:"output_data,omitempty"`
}
