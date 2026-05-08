package ikk

import "time"

type IkkResponse struct {
	ID               int		`json:"id"`
	KodeBidangUrusan string		`json:"kode_bidang_urusan"`
	Jenis            string		`json:"jenis"`
	NamaIndikator    string		`json:"nama_indikator"`
	Target           string		`json:"target"`
	Satuan           string		`json:"satuan"`
	Keterangan       string		`json:"keterangan"`
	CreatedAt        time.Time	`json:"created_at"`
	UpdatedAt        time.Time	`json:"updated_at"`
}

type IkkFullResponse struct {
	ID               int		`json:"id"`
	KodeBidangUrusan string		`json:"kode_bidang_urusan"`
	NamaOpd 		 string		`json:"nama_opd"`
	Jenis            string		`json:"jenis"`
	NamaIndikator    string		`json:"nama_indikator"`
	Target           string		`json:"target"`
	Satuan           string		`json:"satuan"`
	Keterangan       string		`json:"keterangan"`
	CreatedAt        time.Time	`json:"created_at"`
	UpdatedAt        time.Time	`json:"updated_at"`
}