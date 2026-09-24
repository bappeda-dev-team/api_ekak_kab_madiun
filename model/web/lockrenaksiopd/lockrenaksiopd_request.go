package lockrenaksiopd

type LockRenaksiOpdRequest struct {
	SasaranId    int    `json:"sasaran_id" validate:"required"`
	RekinId      string `json:"rekin_id" validate:"required"`
	AksiKegiatan string `json:"aksi_kegiatan"`
	SubKegiatan  string `json:"sub_kegiatan"`
	Anggaran     int64  `json:"anggaran"`
	NamaPemilik  string `json:"nama_pemilik"`
	Tw1          int    `json:"tw1"`
	Tw2          int    `json:"tw2"`
	Tw3          int    `json:"tw3"`
	Tw4          int    `json:"tw4"`
}
