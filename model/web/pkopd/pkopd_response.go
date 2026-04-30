package pkopd

type PkOpdResponse struct {
	KodeOpd       string           `json:"kode_opd"`
	NamaOpd       string           `json:"nama_opd"`
	KepalaOpd     string           `json:"nama_kepala_opd"`
	NipKepalaOpd  string           `json:"nip_kepala_opd"`
	Tahun         int              `json:"tahun"`
	PkItem        []PkOpdByLevel   `json:"pk_item"`
	SasaranPemdas []SasaranPemdaPk `json:"sasaran_pemdas"`
}

type SasaranPemdaPk struct {
	NamaKepalaPemda  string `json:"nama_kepala_pemda"`
	NipKepalaPemda   string `json:"nip_kepala_pemda"`
	IdSasaranPemda   int    `json:"id_sasaran_pemda"`
	NamaSasaranPemda string `json:"sasaran_pemda"`
}

type PkOpdByLevel struct {
	LevelPk  int         `json:"level_pk"`
	Pegawais []PkPegawai `json:"pegawais"`
}

type PkPegawai struct {
	NipAtasan      string   `json:"nip_atasan"`
	NamaAtasan     string   `json:"nama_atasan"`
	JabatanAtasan  string   `json:"jabatan_atasan"`
	Nama           string   `json:"nama_pegawai"`
	Nip            string   `json:"nip"`
	JabatanPegawai string   `json:"jabatan_pegawai"`
	Pks            []PkAsn  `json:"pks"`
	LevelPk        int      `json:"level_pk"`
	JenisItem      string   `json:"jenis_item"`
	Item           []ItemPk `json:"item_pk"`
	TotalPagu      int64    `json:"total_pagu"`
	Roles          []string `json:"roles"`
}

type ItemPk struct {
	RekinId  string `json:"id_rekin"`
	KodeItem string `json:"kode_item"`
	NamaItem string `json:"nama_item"`
	PaguItem int64  `json:"pagu_item"`
}

type PkAsn struct {
	Id               string        `json:"id"`
	IdPohon          int           `json:"id_pohon"`
	IdParentPohon    int           `json:"id_parent_pohon"`
	KodeOpd          string        `json:"kode_opd"`
	NamaOpd          string        `json:"nama_opd"`
	LevelPk          int           `json:"level_pk"`
	NipAtasan        string        `json:"nip_atasan"`
	NamaAtasan       string        `json:"nama_atasan"`
	IdRekinAtasan    string        `json:"id_rekin_atasan"`
	RekinAtasan      string        `json:"rekin_atasan"`
	NipPemilikPk     string        `json:"nip_pemilik_pk"`
	NamaPemilikPk    string        `json:"nama_pemilik_pk"`
	IdRekinPemilikPk string        `json:"id_rekin_pemilik_pk"`
	RekinPemilikPk   string        `json:"rekin_pemilik_pk"`
	Tahun            int           `json:"tahun"`
	Keterangan       string        `json:"keterangan"`
	Indikators       []IndikatorPk `json:"indikators"`
	// daftar atasan untuk menghubungkan rekin pegawai
	AtasanCandidates []AtasanCandidate `json:"atasan_candidates"`
}

type AtasanCandidate struct {
	IdPegawai           string `json:"id_pegawai"`
	NamaPegawai         string `json:"nama_pegawai"`
	LevelPegawai        int    `json:"level_pegawai"`
	KodeOpd             string `json:"kode_opd"`
	NamaOpd             string `json:"nama_opd"`
	IdPohonAtasan       int    `json:"id_pohon_atasan"`
	IdParentPohonAtasan int    `json:"id_parent_pohon_atasan"`
}

type IndikatorPk struct {
	IdRekin     string        `json:"id_rekin"`
	IdIndikator string        `json:"id_indikator"`
	Indikator   string        `json:"indikator"`
	Targets     []TargetIndPk `json:"targets"`
}

type TargetIndPk struct {
	IdIndikator string `json:"id_indikator"`
	IdTarget    string `json:"id_target"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}
