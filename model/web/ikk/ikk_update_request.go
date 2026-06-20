package ikk

type IkkUpdateRequest struct {
	ID               int                `json:"id"`
	KodeBidangUrusan string             `json:"kode_bidang_urusan" validate:"required"`
	KodeOpd          string             `json:"kode_opd" validate:"required"`
	Jenis            string             `json:"jenis" validate:"required"`
	Tahun            int                `json:"tahun" validate:"required"`
	Keterangan       string             `json:"keterangan"`
	Indikators       []IndikatorRequest `json:"indikators" validate:"required,dive"`
}