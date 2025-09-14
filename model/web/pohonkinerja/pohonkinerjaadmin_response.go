package pohonkinerja

import (
	"ekak_kabupaten_madiun/model/web/opdmaster"
)

type PohonKinerjaAdminResponse struct {
	Tahun   string            `json:"tahun,omitempty"`
	Tematik []TematikResponse `json:"tematiks"`
}

type PohonKinerjaAdminResponseData struct {
	Id              int                          `json:"id"`
	Parent          int                          `json:"parent,omitempty"`
	NamaPohon       string                       `json:"nama_pohon"`
	KodeOpd         string                       `json:"kode_opd,omitempty"`
	NamaOpd         string                       `json:"nama_opd,omitempty"`
	PerangkatDaerah *opdmaster.OpdResponseForAll `json:"perangkat_daerah,omitempty"`
	Keterangan      string                       `json:"keterangan,omitempty"`
	Tahun           string                       `json:"tahun"`
	NamaOpdPengaju  string                       `json:"nama_opd_pengaju,omitempty"`
	JenisPohon      string                       `json:"jenis_pohon"`
	LevelPohon      int                          `json:"level_pohon"`
	Status          string                       `json:"status"`
	Tagging         []TaggingResponse            `json:"tagging"`
	IsActive        bool                         `json:"is_active"`
	CountReview     int                          `json:"jumlah_review"`
	Pelaksana       []PelaksanaOpdResponse       `json:"pelaksana,omitempty"`
	Indikators      []IndikatorResponse          `json:"indikator,omitempty"`
	Childs          []interface{}                `json:"childs,omitempty"`
	CSFResponse     `json:",inline"`
	UpdatedBy       string `json:"updated_by"`
	// SubTematiks []SubtematikResponse `json:"sub_tematiks,omitempty"`
}

type CSFResponse struct {
	PernyataanKondisiStrategis string `json:"pernyataan_kondisi_strategis"`
	AlasanKondisiStrategis     string `json:"alasan_sebagai_kondisi_strategis"`
	DataTerukur                string `json:"data_terukur_pendukung_pernyataan"`
	KondisiTerukur             string `json:"kondisi_terukur_yang_diharapkan"`
	KondisiWujud               string `json:"kondisi_yang_ingin_diwujudkan"`
}

type TematikResponse struct {
	// CSF         CSFApiResponse      `json:"csf"`
	Id           int                 `json:"id"`
	Parent       *int                `json:"parent"`
	Tema         string              `json:"tema"`
	JenisPohon   string              `json:"jenis_pohon"`
	LevelPohon   int                 `json:"level_pohon"`
	Keterangan   string              `json:"keterangan"`
	CountReview  int                 `json:"jumlah_review"`
	IsActive     bool                `json:"is_active"`
	TaggingPokin []TaggingResponse   `json:"tagging"`
	Indikators   []IndikatorResponse `json:"indikator"`
	// SubTematiks []SubtematikResponse `json:"childs,omitempty"`
	// Strategics  []StrategicResponse  `json:"strategics,omitempty"`
	Child []interface{} `json:"childs,omitempty"`
}

type SubtematikResponse struct {
	// Outcome     []outcome.OutcomeResponse `json:"outcome"`
	Id           int                 `json:"id"`
	Parent       int                 `json:"parent"`
	Tema         string              `json:"tema"`
	JenisPohon   string              `json:"jenis_pohon"`
	LevelPohon   int                 `json:"level_pohon"`
	Keterangan   string              `json:"keterangan"`
	Indikators   []IndikatorResponse `json:"indikator"`
	CountReview  int                 `json:"jumlah_review"`
	IsActive     bool                `json:"is_active"`
	TaggingPokin []TaggingResponse   `json:"tagging"`
	// SubSubTematiks []SubSubTematikResponse `json:"childs,omitempty"`
	// Strategics     []StrategicResponse     `json:"strategics,omitempty"`
	Child []interface{} `json:"childs,omitempty"`
}

type SubSubTematikResponse struct {
	Id           int                 `json:"id"`
	Parent       int                 `json:"parent"`
	Tema         string              `json:"tema"`
	JenisPohon   string              `json:"jenis_pohon"`
	LevelPohon   int                 `json:"level_pohon"`
	Keterangan   string              `json:"keterangan"`
	CountReview  int                 `json:"jumlah_review"`
	IsActive     bool                `json:"is_active"`
	TaggingPokin []TaggingResponse   `json:"tagging"`
	Indikators   []IndikatorResponse `json:"indikator"`
	// SuperSubTematiks []SuperSubTematikResponse `json:"childs,omitempty"`
	// Strategics       []StrategicResponse       `json:"strategics,omitempty"`
	Child []interface{} `json:"childs,omitempty"`
}

type SuperSubTematikResponse struct {
	Id           int                 `json:"id"`
	Parent       int                 `json:"parent"`
	Tema         string              `json:"tema"`
	JenisPohon   string              `json:"jenis_pohon"`
	LevelPohon   int                 `json:"level_pohon"`
	Keterangan   string              `json:"keterangan"`
	CountReview  int                 `json:"jumlah_review"`
	IsActive     bool                `json:"is_active"`
	TaggingPokin []TaggingResponse   `json:"tagging"`
	Indikators   []IndikatorResponse `json:"indikator"`
	Childs       []interface{}       `json:"childs,omitempty"`
}

type StrategicResponse struct {
	// Intermediate []intermediate.IntermediateResponse `json:"intermediate"`
	Id           int                          `json:"id"`
	Parent       int                          `json:"parent"`
	Strategi     string                       `json:"tema"`
	JenisPohon   string                       `json:"jenis_pohon"`
	LevelPohon   int                          `json:"level_pohon"`
	Keterangan   string                       `json:"keterangan"`
	Status       string                       `json:"status"`
	CountReview  int                          `json:"jumlah_review"`
	IsActive     bool                         `json:"is_active"`
	TaggingPokin []TaggingResponse            `json:"tagging"`
	KodeOpd      *opdmaster.OpdResponseForAll `json:"perangkat_daerah,omitempty"`
	Pelaksana    []PelaksanaOpdResponse       `json:"pelaksana,omitempty"`
	Indikators   []IndikatorResponse          `json:"indikator"`
	Childs       []interface{}                `json:"childs,omitempty"`
}

type TacticalResponse struct {
	Id           int                          `json:"id"`
	Parent       int                          `json:"parent"`
	Strategi     string                       `json:"tema"`
	JenisPohon   string                       `json:"jenis_pohon"`
	LevelPohon   int                          `json:"level_pohon"`
	Keterangan   *string                      `json:"keterangan"`
	Status       string                       `json:"status"`
	CountReview  int                          `json:"jumlah_review"`
	IsActive     bool                         `json:"is_active"`
	TaggingPokin []TaggingResponse            `json:"tagging"`
	KodeOpd      *opdmaster.OpdResponseForAll `json:"perangkat_daerah,omitempty"`
	Pelaksana    []PelaksanaOpdResponse       `json:"pelaksana,omitempty"`
	Indikators   []IndikatorResponse          `json:"indikator"`
	Childs       []interface{}                `json:"childs,omitempty"`
}

type OperationalResponse struct {
	Id           int                          `json:"id"`
	Parent       int                          `json:"parent"`
	Strategi     string                       `json:"tema"`
	JenisPohon   string                       `json:"jenis_pohon"`
	LevelPohon   int                          `json:"level_pohon"`
	Keterangan   *string                      `json:"keterangan"`
	Status       string                       `json:"status"`
	CountReview  int                          `json:"jumlah_review"`
	IsActive     bool                         `json:"is_active"`
	TaggingPokin []TaggingResponse            `json:"tagging"`
	KodeOpd      *opdmaster.OpdResponseForAll `json:"perangkat_daerah,omitempty"`
	Pelaksana    []PelaksanaOpdResponse       `json:"pelaksana,omitempty"`
	Indikators   []IndikatorResponse          `json:"indikator"`
	Childs       []interface{}                `json:"childs,omitempty"`
}

type OperationalNResponse struct {
	Id           int                          `json:"id"`
	Parent       int                          `json:"parent"`
	Strategi     string                       `json:"tema"`
	JenisPohon   string                       `json:"jenis_pohon"`
	LevelPohon   int                          `json:"level_pohon"`
	Keterangan   *string                      `json:"keterangan"`
	Status       string                       `json:"status"`
	CountReview  int                          `json:"jumlah_review"`
	IsActive     bool                         `json:"is_active"`
	TaggingPokin []TaggingResponse            `json:"tagging"`
	KodeOpd      *opdmaster.OpdResponseForAll `json:"perangkat_daerah,omitempty"`
	Pelaksana    []PelaksanaOpdResponse       `json:"pelaksana,omitempty"`
	Indikators   []IndikatorResponse          `json:"indikator"`
	Childs       []OperationalNResponse       `json:"childs,omitempty"`
}

type TematikListOpdResponse struct {
	Tematik    string            `json:"tematik"`
	LevelPohon int               `json:"level_pohon"`
	Tahun      string            `json:"tahun"`
	IsActive   bool              `json:"is_active"`
	ListOpd    []OpdListResponse `json:"list_opd"`
}

type OpdListResponse struct {
	KodeOpd         string `json:"kode_opd"`
	PerangkatDaerah string `json:"perangkat_daerah"`
}
