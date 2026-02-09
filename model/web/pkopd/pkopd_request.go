package pkopd

type PkOpdRequest struct {
	IdPohon          int    `json:"id_pohon"`
	KodeOpd          string `json:"kode_opd"`
	LevelPk          int    `json:"level_pk"`
	NipAtasan        string `json:"nip_atasan"`
	IdRekinAtasan    string `json:"id_rekin_atasan"`
	NipPemilikPk     string `json:"nip_pemilik_pk"`
	IdRekinPemilikPk string `json:"id_rekin_pemilik_pk"`
	Tahun            int    `json:"tahun"`
	Keterangan       string `json:"keterangan"`
}

type HubungkanAtasanRequest struct {
	KodeOpd    string `json:"kode_opd" validate:"required"`
	Tahun      int    `json:"tahun" validate:"required"`
	NipBawahan string `json:"nip_bawahan" validate:"required"`
	NipAtasan  string `json:"nip_atasan" validate:"required"`
}
