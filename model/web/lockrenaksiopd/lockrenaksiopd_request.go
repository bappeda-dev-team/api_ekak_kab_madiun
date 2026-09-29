package lockrenaksiopd

import "ekak_kabupaten_madiun/model/web"

type LockRenaksiOpdRequest struct {
	SasaranId web.StringID `json:"sasaran_id" validate:"required" swaggertype:"string"`
	RekinId   string       `json:"rekin_id" validate:"required"`
	AksiKegiatan    string `json:"aksi_kegiatan"`
	KodeSubKegiatan string `json:"kode_subkegiatan"`
	NamaSubKegiatan string `json:"nama_subkegiatan"`
	Anggaran        int64  `json:"anggaran"`
	NamaPemilik     string `json:"nama_pemilik"`
	Tw1             int    `json:"tw1"`
	Tw2             int    `json:"tw2"`
	Tw3             int    `json:"tw3"`
	Tw4             int    `json:"tw4"`
}
