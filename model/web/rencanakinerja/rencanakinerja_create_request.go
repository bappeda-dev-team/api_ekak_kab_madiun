package rencanakinerja

type RencanaKinerjaCreateRequest struct {
	IdPohon              int                      `json:"id_pohon"`
	NamaRencanaKinerja   string                   `json:"nama_rencana_kinerja"`
	Tahun                string                   `json:"tahun"`
	StatusRencanaKinerja string                   `json:"status_rencana_kinerja"`
	Catatan              string                   `json:"catatan"`
	KodeOpd              string                   `json:"kode_opd"`
	PegawaiId            string                   `json:"pegawai_id"`
	PeriodeId            int                      `json:"periode_id"`
	TahunAwal            string                   `json:"tahun_awal"`
	TahunAkhir           string                   `json:"tahun_akhir"`
	JenisPeriode         string                   `json:"jenis_periode"`
	Indikator            []IndikatorCreateRequest `json:"indikator"`
}

type IndikatorCreateRequest struct {
	NamaIndikator string                `json:"nama_indikator"`
	Formula       string                `json:"rumus_perhitungan"`
	SumberData    string                `json:"sumber_data"`
	Tahun         string                `json:"tahun"`
	Target        []TargetCreateRequest `json:"target"`
}

type TargetCreateRequest struct {
	Tahun           string `json:"tahun"`
	Target          string `json:"target"`
	SatuanIndikator string `json:"satuan"`
}
