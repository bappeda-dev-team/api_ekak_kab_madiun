package pohonkinerja

type PohonKinerjaUpdateRequest struct {
	Id           int                      `json:"id"`
	Parent       int                      `json:"parent"`
	NamaPohon    string                   `json:"nama_pohon"`
	JenisPohon   string                   `json:"jenis_pohon"`
	LevelPohon   int                      `json:"level_pohon"`
	KodeOpd      string                   `json:"kode_opd"`
	Keterangan   string                   `json:"keterangan"`
	Tahun        string                   `json:"tahun"`
	Status       string                   `json:"status"`
	PelaksanaId  []PelaksanaUpdateRequest `json:"pelaksana"`
	Indikator    []IndikatorUpdateRequest `json:"indikator"`
	TaggingPokin []TaggingUpdateRequest   `json:"tagging"`
}

type PelaksanaUpdateRequest struct {
	PegawaiId string `json:"pegawai_id"`
}

type PohonKinerjaAdminUpdateRequest struct {
	Id           int                      `json:"id"`
	Parent       int                      `json:"parent"`
	NamaPohon    string                   `json:"nama_pohon"`
	JenisPohon   string                   `json:"jenis_pohon"`
	KodeOpd      string                   `json:"kode_opd,omitempty"`
	LevelPohon   int                      `json:"level_pohon"`
	Keterangan   string                   `json:"keterangan"`
	Tahun        string                   `json:"tahun"`
	Status       string                   `json:"status"`
	TaggingPokin []TaggingUpdateRequest   `json:"tagging"`
	Pelaksana    []PelaksanaUpdateRequest `json:"pelaksana"`
	Indikator    []IndikatorUpdateRequest `json:"indikator"`
	CSFRequest   `json:",inline"`
}

type IndikatorUpdateRequest struct {
	Id             string                `json:"id"`
	PohonKinerjaId int                   `json:"pohon_id"`
	NamaIndikator  string                `json:"indikator"`
	Target         []TargetUpdateRequest `json:"target"`
}

type TargetUpdateRequest struct {
	Id          string `json:"id"`
	IndikatorId string `json:"indikator_id"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}

type PohonKinerjaAdminTolakRequest struct {
	Id int `json:"id" validate:"required"`
}

type TaggingUpdateRequest struct {
	Id                int    `json:"id"`
	NamaTagging       string `json:"nama_tagging"`
	KeteranganTagging string `json:"keterangan_tagging"`
}
