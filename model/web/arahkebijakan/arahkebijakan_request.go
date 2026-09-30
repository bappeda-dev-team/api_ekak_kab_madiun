package arahkebijakan

type ArahKebijakanRequest struct {
	PokinId int    `json:"pokin_id" validate:"required"`
	Arah    string `json:"arah" validate:"required"`
	KodeOpd string `json:"kode_opd" validate:"required"`
	Tahun   int    `json:"tahun" validate:"required"`
}

type ArahKebijakanUpdateRequest struct {
	ID      int    `json:"id"`
	PokinId int    `json:"pokin_id" validate:"required"`
	Arah    string `json:"arah" validate:"required"`
	KodeOpd string `json:"kode_opd" validate:"required"`
	Tahun   int    `json:"tahun" validate:"required"`
}