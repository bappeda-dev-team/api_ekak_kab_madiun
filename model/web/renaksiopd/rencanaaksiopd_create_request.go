package renaksiopd

type RencanaAksiOpdCreateRequest struct {
	SasaranOpdId int    `json:"sasaranopd_id" validate:"required"`
	RekinId      string `json:"rekin_id" validate:"required"`
	TahunRenaksi string `json:"tahun" validate:"required"`
	Keterangan   string `json:"keterangan"`
}
