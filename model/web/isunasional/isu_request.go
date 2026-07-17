package isunasional

type IsuNasionalRequest struct {
	KodeBidangUrusan string `json:"kode_bidang_urusan" validate:"required"`
	KodeOpd          string `json:"kode_opd" validate:"required"`
	Isu              string `json:"isu" validate:"required"`
	Tahun            int    `json:"tahun" validate:"required"`
}

type IsuNasionalUpdateRequest struct {
	ID               int    `json:"id"`
	KodeBidangUrusan string `json:"kode_bidang_urusan" validate:"required"`
	KodeOpd          string `json:"kode_opd" validate:"required"`
	Isu              string `json:"isu" validate:"required"`
	Tahun            int    `json:"tahun" validate:"required"`
}