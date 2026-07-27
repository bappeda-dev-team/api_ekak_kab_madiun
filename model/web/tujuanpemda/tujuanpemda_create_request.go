package tujuanpemda

type TujuanPemdaCreateRequest struct {
	IdVisi            int                      `json:"id_visi"`
	IdMisi            int                      `json:"id_misi"`
	TujuanPemda       string                   `json:"tujuan_pemda" validate:"required"`
	TematikId         int                      `json:"tema_id" validate:"required"`
	PeriodeId         int                      `json:"periode_id" validate:"required"`
	TahunAwalPeriode  string                   `json:"tahun_awal_periode" validate:"required"`
	TahunAkhirPeriode string                   `json:"tahun_akhir_periode" validate:"required"`
	JenisPeriode      string                   `json:"jenis_periode" validate:"required"`
	Indikator         []IndikatorCreateRequest `json:"indikator"`
}

type IndikatorCreateRequest struct {
	Indikator           string                `json:"indikator"`
	RumusPerhitungan    string                `json:"rumus_perhitungan"`
	SumberData          string                `json:"sumber_data"`
	DefinisiOperasional string                `json:"definisi_operasional"`
	Jenis               string                `json:"jenis"`
	Target              []TargetCreateRequest `json:"target"`
}

type TargetCreateRequest struct {
	Target TargetInput `json:"target"`
	Satuan string      `json:"satuan"`
	Tahun  string      `json:"tahun"`
	Jenis  string      `json:"jenis"`
}
