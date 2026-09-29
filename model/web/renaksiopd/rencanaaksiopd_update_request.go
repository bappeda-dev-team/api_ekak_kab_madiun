package renaksiopd

type RencanaAksiOpdUpdateRequest struct {
	Id         int    `json:"id" `
	RekinId    string `json:"rekin_id" validate:"required"`
	Urutan     int    `json:"urutan" validate:"required,min=1"`
	Keterangan string `json:"keterangan"`
}
