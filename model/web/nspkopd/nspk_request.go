package nspkopd

type NspkRequest struct {
	KodeOpd      string `json:"kode_opd" validate:"required"`
	IdNspk       int    `json:"id_nspk" validate:"required"`
	IdTujuanOpd  int    `json:"id_tujuan_opd" validate:"required"`
	IdSasaranOpd int    `json:"id_sasaran_opd" validate:"required"`
	Tahun        int    `json:"tahun" validate:"required"`
}

type NspkUpdateRequest struct {
	ID           int    `json:"id"`
	KodeOpd      string `json:"kode_opd" validate:"required"`
	IdNspk       int    `json:"id_nspk" validate:"required"`
	IdTujuanOpd  int    `json:"id_tujuan_opd" validate:"required"`
	IdSasaranOpd int    `json:"id_sasaran_opd" validate:"required"`
	Tahun        int    `json:"tahun" validate:"required"`
}