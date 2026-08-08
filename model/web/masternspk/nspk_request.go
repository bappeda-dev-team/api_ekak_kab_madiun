package masternspk

type NspkRequest struct {
	KodeOpd string `json:"kode_opd" validate:"required"`
	Nspk    string `json:"nspk" validate:"required"`
	Tahun   int    `json:"tahun" validate:"required"`
}

type NspkUpdateRequest struct {
	ID      int    `json:"id"`
	KodeOpd string `json:"kode_opd" validate:"required"`
	Nspk    string `json:"nspk" validate:"required"`
	Tahun   int    `json:"tahun" validate:"required"`
}