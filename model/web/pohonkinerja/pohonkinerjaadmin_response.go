package pohonkinerja

type PohonKinerjaAdminResponse struct {
	Tahun   string            `json:"tahun"`
	Tematik []TematikResponse `json:"tematiks"`
}

type PohonKinerjaAdminResponseData struct {
	Id          int                  `json:"id"`
	Parent      int                  `json:"parent"`
	NamaPohon   string               `json:"nama_pohon"`
	KodeOpd     string               `json:"kode_opd"`
	Keterangan  string               `json:"keterangan"`
	Tahun       string               `json:"tahun"`
	JenisPohon  string               `json:"jenis_pohon"`
	LevelPohon  int                  `json:"level_pohon"`
	Indikators  []IndikatorResponse  `json:"indikators"`
	SubTematiks []SubtematikResponse `json:"sub_tematiks,omitempty"`
}

type TematikResponse struct {
	Id          int                       `json:"id"`
	Parent      *int                      `json:"parent"`
	Tema        string                    `json:"tema"`
	Keterangan  string                    `json:"keterangan"`
	Indikators  []IndikatorSimpleResponse `json:"indikators"`
	SubTematiks []SubtematikResponse      `json:"sub_tematiks,omitempty"`
}

type SubtematikResponse struct {
	Id             int                       `json:"id"`
	Parent         int                       `json:"parent"`
	Tema           string                    `json:"tema"`
	Keterangan     string                    `json:"keterangan"`
	Indikators     []IndikatorSimpleResponse `json:"indikators"`
	SubSubTematiks []SubSubTematikResponse   `json:"sub_sub_tematiks,omitempty"`
	Strategics     []StrategicResponse       `json:"strategics,omitempty"`
}

type SubSubTematikResponse struct {
	Id               int                       `json:"id"`
	Parent           int                       `json:"parent"`
	Tema             string                    `json:"tema"`
	Keterangan       string                    `json:"keterangan"`
	Indikators       []IndikatorSimpleResponse `json:"indikators"`
	SuperSubTematiks []SuperSubTematikResponse `json:"super_sub_tematiks,omitempty"`
	Strategics       []StrategicResponse       `json:"strategics,omitempty"`
}

type SuperSubTematikResponse struct {
	Id         int                       `json:"id"`
	Parent     int                       `json:"parent"`
	Tema       string                    `json:"tema"`
	Keterangan string                    `json:"keterangan"`
	Indikators []IndikatorSimpleResponse `json:"indikators"`
	Strategics []StrategicResponse       `json:"strategics,omitempty"`
}

type IndikatorResponse struct {
	Id               string           `json:"id_indikator"`
	RencanaKinerjaId string           `json:"rencana_kinerja_id"`
	NamaIndikator    string           `json:"nama_indikator"`
	Target           []TargetResponse `json:"targets"`
}

type TargetResponse struct {
	Id              string `json:"id_target"`
	IndikatorId     string `json:"indikator_id"`
	TargetIndikator string `json:"target"`
	SatuanIndikator string `json:"satuan"`
}

type StrategicResponse struct {
	Id              int                       `json:"id"`
	Parent          int                       `json:"parent"`
	Strategi        string                    `json:"strategi"`
	Keterangan      string                    `json:"keterangan"`
	KodeOpd         string                    `json:"kode_perangkat_daerah"`
	PerangkatDaerah string                    `json:"perangkat_daerah"`
	Indikators      []IndikatorSimpleResponse `json:"indikators"`
	Tacticals       []TacticalResponse        `json:"tacticals,omitempty"`
}

type TacticalResponse struct {
	Id              int                       `json:"id"`
	Parent          int                       `json:"parent"`
	Strategi        string                    `json:"strategi"`
	Keterangan      *string                   `json:"keterangan"`
	KodeOpd         string                    `json:"kode_perangkat_daerah"`
	PerangkatDaerah string                    `json:"perangkat_daerah"`
	Indikators      []IndikatorSimpleResponse `json:"indikators"`
	Operationals    []OperationalResponse     `json:"operationals"`
}

type OperationalResponse struct {
	Id              int                       `json:"id"`
	Parent          int                       `json:"parent"`
	Strategi        string                    `json:"strategi"`
	Keterangan      *string                   `json:"keterangan"`
	KodeOpd         string                    `json:"kode_perangkat_daerah"`
	PerangkatDaerah string                    `json:"perangkat_daerah"`
	Indikators      []IndikatorSimpleResponse `json:"indikators"`
}

type IndikatorSimpleResponse struct {
	Indikator string `json:"indikator"`
	Target    string `json:"target"`
	Satuan    string `json:"satuan"`
}