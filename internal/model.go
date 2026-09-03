package internal

import "time"

type IsuStrategisResponse struct {
	Id               int                    `json:"id"`
	KodeOpd          string                 `json:"kode_opd"`
	TahunAwal        string                 `json:"tahun_awal"`
	TahunAkhir       string                 `json:"tahun_akhir"`
	IsuStrategis     string                 `json:"isu_strategis"`
	CreatedAt        time.Time              `json:"created_at"`
}
type PermasalahanResp struct {
	PermasalahanOpd  []PermasalahanResponse `json:"permasalahan_opd"`
}

type PermasalahanResponse struct {
	Id           int                  `json:"id"`
	Permasalahan string               `json:"masalah"`
	LevelPohon   int                  `json:"level_pohon"`
	JenisMasalah string               `json:"jenis_masalah"`
}

type IsuStrategicWrapper struct {
	Code           int                `json:"code"`
	Status         string             `json:"status"`
	Data 		   []IsuStrategisResponse `json:"data"`
}

type PermasalahanWrapper struct {
	Code           int                `json:"code"`
	Status         string             `json:"status"`
	Data 		   []PermasalahanResp `json:"data"`
}

type FindByKodeOpdsTahunsRequest struct {
	KodeOpd []string `json:"kode_opd" validate:"required,min=1"`
	Tahun []string `json:"tahun" validate:"required,min=1"`
}