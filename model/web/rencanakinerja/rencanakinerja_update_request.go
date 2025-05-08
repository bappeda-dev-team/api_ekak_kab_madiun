package rencanakinerja

type RencanaKinerjaUpdateRequest struct {
	Id                   string                   `json:"id"`
	SasaranOpdId         int                      `json:"sasaranopd_id" validate:"required"`
	IdPohon              int                      `json:"id_pohon" validate:"required"`
	NamaRencanaKinerja   string                   `json:"nama_rencana_kinerja" validate:"required"`
	Tahun                string                   `json:"tahun" validate:"required"`
	StatusRencanaKinerja string                   `json:"status_rencana_kinerja" validate:"required"`
	Catatan              string                   `json:"catatan"`
	KodeOpd              string                   `json:"kode_opd" validate:"required"`
	PegawaiId            string                   `json:"pegawai_id" validate:"required"`
	PeriodeId            int                      `json:"periode_id"`
	TahunAwal            string                   `json:"tahun_awal"`
	TahunAkhir           string                   `json:"tahun_akhir"`
	JenisPeriode         string                   `json:"jenis_periode"`
	Indikator            []IndikatorUpdateRequest `json:"indikator"`
}

type IndikatorUpdateRequest struct {
	Id               string                `json:"id_indikator"`
	RencanaKinerjaId string                `json:"rencana_kinerja_id"`
	Formula          string                `json:"rumus_perhitungan,omitempty"`
	SumberData       string                `json:"sumber_data,omitempty"`
	Indikator        string                `json:"nama_indikator"`
	Tahun            string                `json:"tahun"`
	Target           []TargetUpdateRequest `json:"target"`
}

type TargetUpdateRequest struct {
	Id              string `json:"id_target"`
	IndikatorId     string `json:"indikator_id"`
	Tahun           string `json:"tahun"`
	Target          string `json:"target"`
	SatuanIndikator string `json:"satuan"`
}
