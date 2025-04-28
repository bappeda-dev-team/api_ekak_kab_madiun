package sasaranopd

type SasaranOpdResponse struct {
	Id             int                    `json:"id"`
	IdPohon        int                    `json:"id_pohon"`
	NamaPohon      string                 `json:"nama_pohon"`
	JenisPohon     string                 `json:"jenis_pohon"`
	TahunPohon     string                 `json:"tahun_pohon"`
	LevelPohon     int                    `json:"level_pohon"`
	RencanaKinerja []RencanaKinerjaOpd    `json:"rencana_kinerja"`
	Pelaksana      []PelaksanaOpdResponse `json:"pelaksana"`
}

type RencanaKinerjaOpd struct {
	Id                 string              `json:"id"`
	TahunAwal          string              `json:"tahun_awal"`
	TahunAkhir         string              `json:"tahun_akhir"`
	JenisPeriode       string              `json:"jenis_periode"`
	NamaRencanaKinerja string              `json:"nama_rencana_kinerja"`
	Nip                string              `json:"nip"`
	NamaPegawai        string              `json:"nama_pegawai"`
	Indikator          []IndikatorResponse `json:"indikator"`
}

type PelaksanaOpdResponse struct {
	Id          string `json:"id"`
	PegawaiId   string `json:"pegawai_id"`
	Nip         string `json:"nip"`
	NamaPegawai string `json:"nama_pegawai"`
}

type IndikatorResponse struct {
	Id         string            `json:"id"`
	Indikator  string            `json:"indikator"`
	ManualIK   *ManualIKResponse `json:"manual_ik"`
	Formula    string            `json:"rumus_perhitungan"`
	SumberData string            `json:"sumber_data"`
	Target     []TargetResponse  `json:"target"`
}

type ManualIKResponse struct {
	Formula    string `json:"formula"`
	SumberData string `json:"sumber_data"`
}

type TargetResponse struct {
	Tahun  string `json:"tahun"`
	Target string `json:"target"`
	Satuan string `json:"satuan"`
}
