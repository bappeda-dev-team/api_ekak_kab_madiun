package renaksiopd

type RencanaAksiOpdUpdateRequest struct {
	Id         int    `json:"id" `
	RekinId    string `json:"rekin_id" validate:"required"`
	Keterangan string `json:"keterangan"`
}
