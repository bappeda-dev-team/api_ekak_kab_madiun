package domain

import "time"

type IsuGlobal struct {
	ID               int
	KodeOpd          string
	NamaOpd          string
	KodeBidangUrusan string
	NamaBidangUrusan string
	Isu              string
	Tahun            int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type BidangUrusanGlobalSelection struct {
	KodeBidangUrusan string `json:"kode_bidang_urusan"`
	NamaBidangUrusan string `json:"nama_bidang_urusan"`
	KodeOpd          string `json:"kode_opd"`
	NamaOpd          string `json:"nama_opd"`
}