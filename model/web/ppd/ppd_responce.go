package ppd

import "time"

type PpdResponse struct {
	ID               int       `json:"id"`
	KodeBidangUrusan string    `json:"kode_bidang_urusan"`
	KodeOpd          string    `json:"kode_opd"`
	Potensi          string    `json:"potensi"`
	Tahun            int       `json:"tahun"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type PpdFullResponse struct {
	ID               int		`json:"id"`
	KodeOpd 		 string		`json:"kode_opd"`
	NamaOpd 		 string		`json:"nama_opd"`
	KodeBidangUrusan string		`json:"kode_bidang_urusan"`
	NamaBidangUrusan string		`json:"nama_bidang_urusan"`
	Potensi          string		`json:"potensi"`
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

type PpdMasterResponse struct {
	BidangUrusanSelections []BidangUrusanSelectionResponse `json:"bidang_urusan_selections"`
	Ppds                   []PpdFullResponse               `json:"ppds"`
}