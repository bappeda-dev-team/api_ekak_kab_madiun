package isuregional

import "time"

type IsuRegionalResponse struct {
	ID               int       `json:"id"`
	KodeBidangUrusan string    `json:"kode_bidang_urusan"`
	KodeOpd          string    `json:"kode_opd"`
	Isu              string    `json:"isu"`
	Tahun            int       `json:"tahun"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type IsuRegionalFullResponse struct {
	ID               int		`json:"id"`
	KodeOpd 		 string		`json:"kode_opd"`
	NamaOpd 		 string		`json:"nama_opd"`
	KodeBidangUrusan string		`json:"kode_bidang_urusan"`
	NamaBidangUrusan string		`json:"nama_bidang_urusan"`
	Isu              string		`json:"isu"`
	Tahun            int		`json:"tahun"`
	CreatedAt        time.Time	`json:"created_at"`
	UpdatedAt        time.Time	`json:"updated_at"`
}

type BidangUrusanSelectionResponse struct {
	KodeBidangUrusan string `json:"kode_bidang_urusan"`
	NamaBidangUrusan string `json:"nama_bidang_urusan"`
	KodeOpd          string `json:"kode_opd"`
	NamaOpd          string `json:"nama_opd"`
}

type IsuRegionalMasterResponse struct {
	BidangUrusanSelections []BidangUrusanSelectionResponse `json:"bidang_urusan_selections"`
	Isus                   []IsuRegionalFullResponse               `json:"isu_regionals"`
}