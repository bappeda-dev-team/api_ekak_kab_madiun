package rencanakinerja

type ManualIKCreateRequest struct {
	Id                  int                     `json:"id"`
	IdIndikator         string                  `json:"id_indikator"`
	Perspektif          string                  `validate:"required" json:"perspektif"`
	TujuanRekin         string                  `json:"tujuan_rekin"`
	Definisi            string                  `json:"definisi"`
	KeyActivities       string                  `json:"key_activities"`
	Formula             string                  `json:"formula"`
	JenisIndikator      string                  `json:"jenis_indikator"`
	OutputData          OutputDataCreateRequest `json:"output_data"`
	UnitPenanggungJawab string                  `json:"unit_penanggung_jawab"`
	UnitPenyediaData    string                  `json:"unit_penyedia_data"`
	SumberData          string                  `json:"sumber_data"`
	JangkaWaktuAwal     string                  `json:"jangka_waktu_awal"`
	JangkaWaktuAkhir    string                  `json:"jangka_waktu_akhir"`
	PeriodePelaporan    string                  `json:"periode_pelaporan"`
}

type OutputDataCreateRequest struct {
	Kinerja  bool `json:"kinerja"`
	Penduduk bool `json:"penduduk"`
	Spatial  bool `json:"spatial"`
}

type RekinByOpdCloneRequest struct {
	KodeOpd     string `json:"kode_opd" validate:"required"`
	TahunSumber string `json:"tahun_sumber" validate:"required,min=4,max=4"`
	TahunTujuan string `json:"tahun_tujuan" validate:"required,min=4,max=4"`
	UpdatedBy   string `json:"updated_by" validate:"required"`
}
