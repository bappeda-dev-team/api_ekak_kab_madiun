package pkopd

type PkOpdResponse struct {
	KodeOpd      string         `json:"kode_opd"`
	NamaOpd      string         `json:"nama_opd"`
	KepalaOpd    string         `json:"nama_kepala_opd"`
	NipKepalaOpd string         `json:"nip_kepala_opd"`
	Tahun        int            `json:"tahun"`
	PkItem       []PkOpdByLevel `json:"pk_item"`
}

type PkOpdByLevel struct {
	LevelPk  int         `json:"level_pk"`
	Pegawais []PkPegawai `json:"pegawais"`
}

type PkPegawai struct {
	JenisItem      string          `json:"jenis_item"`
	NipAtasan      string          `json:"nip_atasan"`
	NamaAtasan     string          `json:"nama_atasan"`
	JabatanAtasan  string          `json:"jabatan_atasan"`
	Nama           string          `json:"nama_pegawai"`
	Nip            string          `json:"nip"`
	JabatanPegawai string          `json:"jabatan_pegawai"`
	Pks            []PkAsn         `json:"pks"`
	Subkegiatan    []SubkegiatanPk `json:"subkegiatan"`
}

type SubkegiatanPk struct {
	IdRekin         string `json:"id_rekin"`
	KodeProgram     string `json:"kode_program"`
	NamaProgram     string `json:"nama_program"`
	KodeKegiatan    string `json:"kode_kegiatan"`
	NamaKegiatan    string `json:"nama_kegiatan"`
	KodeSubkegiatan string `json:"kode_subkegiatan"`
	NamaSubkegiatan string `json:"nama_subkegiatan"`
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
	PaguAnggaran     int           `json:"pagu_anggaran"`
}

type IndikatorPk struct {
	Indikator string        `json:"indikator"`
	Targets   []TargetIndPk `json:"targets"`
}

type TargetIndPk struct {
	Target string `json:"target"`
	Satuan string `json:"satuan"`
}
