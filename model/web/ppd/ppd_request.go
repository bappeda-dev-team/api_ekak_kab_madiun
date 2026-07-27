package ppd

type PpdRequest struct {
	KodeBidangUrusan string `json:"kode_bidang_urusan" validate:"required"`
	KodeOpd          string `json:"kode_opd" validate:"required"`
	Potensi          string `json:"potensi" validate:"required"`
	Tahun            int    `json:"tahun" validate:"required"`
}

type PpdUpdateRequest struct {
	ID               int    `json:"id"`
	KodeBidangUrusan string `json:"kode_bidang_urusan" validate:"required"`
	KodeOpd          string `json:"kode_opd" validate:"required"`
	Potensi          string `json:"potensi" validate:"required"`
	Tahun            int    `json:"tahun" validate:"required"`
}