package tujuanpemda

type TujuanPemdaUpdateRequest struct {
	Id                int                      `json:"id"`
	IdVisi            int                      `json:"id_visi"`
	IdMisi            int                      `json:"id_misi"`
	TujuanPemda       string                   `json:"tujuan_pemda"`
	TematikId         int                      `json:"tema_id"`
	PeriodeId         int                      `json:"periode_id"`
	TahunAwalPeriode  string                   `json:"tahun_awal_periode"`
	TahunAkhirPeriode string                   `json:"tahun_akhir_periode"`
	JenisPeriode      string                   `json:"jenis_periode"`
	Indikator         []IndikatorUpdateRequest `json:"indikator"`
}

type IndikatorUpdateRequest struct {
	Id                  string                `json:"id"`
	IdIndikator         int                   `json:"id_indikator"`
	TujuanPemdaId       string                `json:"tujuan_pemda_id"`
	KodeIndikator       string                `json:"kode_indikator"`
	Jenis               string                `json:"jenis"`
	DefinisiOperasional string                `json:"definisi_operasional"`
	Indikator           string                `json:"indikator"`
	RumusPerhitungan    string                `json:"rumus_perhitungan"`
	SumberData          string                `json:"sumber_data"`
	Target              []TargetUpdateRequest `json:"target"`
}

type TargetUpdateRequest struct {
	Id     int         `json:"id"`
	Target TargetInput `json:"target"`
	Satuan string      `json:"satuan"`
	Tahun  string      `json:"tahun"`
	Jenis  string      `json:"jenis"`
}
