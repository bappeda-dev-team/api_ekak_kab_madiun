package lockrenaksiopd

type LockRenaksiOpdResponse struct {
	Id           int    `json:"id"`
	KodeOpd      string `json:"kode_opd"`
	Tahun        string `json:"tahun"`
	SasaranId    int    `json:"sasaran_id"`
	RekinId      string `json:"rekin_id"`
	AksiKegiatan string `json:"aksi_kegiatan"`
	SubKegiatan  string `json:"sub_kegiatan"`
	Anggaran     int64  `json:"anggaran"`
	NamaPemilik  string `json:"nama_pemilik"`
	Tw1          int    `json:"tw1"`
	Tw2          int    `json:"tw2"`
	Tw3          int    `json:"tw3"`
	Tw4          int    `json:"tw4"`
	Locked       bool   `json:"locked"`
}
