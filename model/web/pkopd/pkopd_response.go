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
	NipAtasan      string `json:"nip_atasan"`
	NamaAtasan     string `json:"nama_atasan"`
	JabatanAtasan  string `json:"jabatan_atasan"`
	Nama           string `json:"nama_pegawai"`
	Nip            string `json:"nip"`
	JabatanPegawai string `json:"jabatan_pegawai"`
	LevelPk        int    `json:"level_pk"`
	JenisItem      string `json:"jenis_item"`
	// program kegiatan subkegiatan
	Item      []ItemPk `json:"item_pk"`
	TotalPagu int64    `json:"total_pagu"`
	Roles     []string `json:"roles"`
	// daftar atasan untuk menghubungkan rekin pegawai
	AtasanCandidates []AtasanCandidate `json:"atasan_candidates"`
	// rekin pegawai
	PkTerkunci bool    `json:"pk_terkunci"` // true, false, false
	Pks        []PkAsn `json:"pks"`
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
	SasaranOpdId     int64         `json:"id_sasaran_opd"`
	RekinPemilikPk   string        `json:"rekin_pemilik_pk"`
	AnggaranPk       int           `json:"anggaran_pk"`
	Tahun            int           `json:"tahun"`
	Keterangan       string        `json:"keterangan"`
	Indikators       []IndikatorPk `json:"indikators"`
	Renaksis         []RenaksiItem `json:"renaksi"`
}

type RenaksiItem struct {
	Id               string         `json:"id_renaksi"`
	RencanaKinerjaId string         `json:"rekin_id"`
	KodeOpd          string         `json:"kode_opd,omitempty"`
	Urutan           int            `json:"urutan"`
	Anggaran         int            `json:"anggaran"`
	NamaRencanaAksi  string         `json:"nama_rencana_aksi"`
	Pelaksanaan      []BobotBulanan `json:"pelaksanaan"`
}

type BobotBulanan struct {
	Id    string `json:"id_pelaksanaan"`
	Bulan int    `json:"bulan"`
	Bobot int    `json:"bobot"`
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

type KunciPKResponse struct {
	IdKunci    int64  `json:"id_pk"`
	IdPegawai  string `json:"id_pegawai"`
	StatusPk   string `json:"status_pk"`   // terkunci, terbuka, revisi
	PkTerkunci bool   `json:"pk_terkunci"` // true, false, false
}
