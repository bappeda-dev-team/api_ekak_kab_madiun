package tujuanopd

type TujuanOpdCreateRequest struct {
	KodeOpd          string                   `json:"kode_opd"`
	KodeBidangUrusan string                   `json:"kode_bidang_urusan"`
	Tujuan           string                   `json:"tujuan"`
	PeriodeId        int                      `json:"periode_id"`
	TahunAwal        string                   `json:"tahun_awal"`
	TahunAkhir       string                   `json:"tahun_akhir"`
	JenisPeriode     string                   `json:"jenis_periode"`
	Indikator        []IndikatorCreateRequest `json:"indikator"`
}

type IndikatorCreateRequest struct {
	IdTujuanOpd         string                `json:"id_tujuan_opd"`
	Indikator           string                `json:"indikator"`
	RumusPerhitungan    string                `json:"rumus_perhitungan"`
	SumberData          string                `json:"sumber_data"`
	Jenis               string                `json:"jenis"`
	DefinisiOperasional string                `json:"definisi_operasional"`
	Target              []TargetCreateRequest `json:"target"`
}

type TargetCreateRequest struct {
	IndikatorId string `json:"indikator_id"`
	Target      string `json:"target"`
	Tahun       string `json:"tahun"`
	Satuan      string `json:"satuan"`
}

type IndikatorUpsertRequest struct {
	KodeIndikator       string                `json:"kode_indikator"` // kosong = CREATE baru
	Indikator           string                `json:"indikator"`
	DefinisiOperasional string                `json:"definisi_operasional"`
	RumusPerhitungan    string                `json:"rumus_perhitungan"`
	SumberData          string                `json:"sumber_data"`
	Target              []TargetUpsertRequest `json:"target"`
}
type TargetUpsertRequest struct {
	Id     string `json:"id"` // kosong = CREATE baru
	Tahun  string `json:"tahun"`
	Target string `json:"target"`
	Satuan string `json:"satuan"`
}
