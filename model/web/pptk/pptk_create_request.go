package pptk

type PptkCreateRequest struct {
	Nip             string  `json:"nip" validate:"required"`
	KodeOpd         string  `json:"kode_opd" validate:"required"`
	Tahun           int     `json:"tahun" validate:"required"`
	KodeSubKegiatan string  `json:"kode_sub_kegiatan" validate:"required"`
	NipAtasan       *string `json:"nip_atasan"`
	AktifAt         string  `json:"aktif_at" validate:"required"`
}
