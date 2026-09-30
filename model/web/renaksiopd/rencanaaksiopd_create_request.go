package renaksiopd

import "ekak_kabupaten_madiun/model/web"

type RencanaAksiOpdCreateRequest struct {
	SasaranOpdId web.StringID `json:"sasaranopd_id" validate:"required" swaggertype:"string"`
	RekinId      string       `json:"rekin_id" validate:"required"`
	TahunRenaksi string       `json:"tahun" validate:"required"`
	Keterangan   string       `json:"keterangan"`
}
