package ikk

import "time"

type IkkResponse struct {
	ID               int		`json:"id"`
	KodeBidangUrusan string		`json:"kode_bidang_urusan"`
	KodeOpd 		 string		`json:"kode_opd"`
	Jenis            string		`json:"jenis"`
	Tahun            int		`json:"tahun"`
	Keterangan       string		`json:"keterangan"`
	Indikators    	 []IndikatorResponse	`json:"indikators"`
	CreatedAt        time.Time	`json:"created_at"`
	UpdatedAt        time.Time	`json:"updated_at"`
}

type TargetResponse struct {
	ID      int    `json:"id"`
	Target  string `json:"target"`
	Satuan  string `json:"satuan"`
	Tahun   int    `json:"tahun"`
}

type IndikatorResponse struct {
	ID         int              `json:"id"`
	Indikator  string           `json:"indikator"`
	Targets    []TargetResponse `json:"targets"`
}

type IkkFullResponse struct {
	ID               int		`json:"id"`
	KodeOpd 		 string		`json:"kode_opd"`
	NamaOpd 		 string		`json:"nama_opd"`
	KodeBidangUrusan string		`json:"kode_bidang_urusan"`
	NamaBidangUrusan string		`json:"nama_bidang_urusan"`
	Jenis            string		`json:"jenis"`
	Tahun            int		`json:"tahun"`
	Keterangan       string		`json:"keterangan"`
	Indikators    	 []IndikatorResponse	`json:"indikators"`
	CreatedAt        time.Time	`json:"created_at"`
	UpdatedAt        time.Time	`json:"updated_at"`
}

type BidangUrusanSelectionResponse struct {
	KodeBidangUrusan string `json:"kode_bidang_urusan"`
	NamaBidangUrusan string `json:"nama_bidang_urusan"`
	KodeOpd          string `json:"kode_opd"`
	NamaOpd          string `json:"nama_opd"`
}

type IkkMasterResponse struct {
	BidangUrusanSelections []BidangUrusanSelectionResponse `json:"bidang_urusan_selections"`
	Ikks                   []IkkFullResponse               `json:"ikks"`
}

type IkkTerpilihResponse struct {
	Id              int    `json:"id"`
	PokinId         int    `json:"pokin_id"`
	IkkId           int    `json:"ikk_id"`
	NamaPokin       string `json:"nama_pokin"`
	JenisIkk        string `json:"jenis_ikk"`
	KeteranganIkk   string `json:"keterangan_ikk"`
}