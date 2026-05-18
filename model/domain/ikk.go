package domain

import "time"

type Ikk struct {
	ID               int   
	KodeOpd 		 string 
	NamaOpd 		 string 
	KodeBidangUrusan string 
	NamaBidangUrusan string 
	Jenis            string 
	NamaIndikator    string 
	Indikators    	 []IndikatorIkk
	Target           string 
	Satuan           string 
	Tahun            int 
	Keterangan       string 
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type IndikatorIkk struct {
	ID               int 
	IDIkk            int
	KodeOpd 		 string 
	KodeBidangUrusan string 
	Indikator   	 string
	Targets   	     []TargetIkk
	Tahun            int 
}

type TargetIkk struct {
	ID               int 
	IDIndikator      int
	Target 		     string 
	Satuan           string 
	Tahun            int 
}

type BidangUrusanSelection struct {
	KodeBidangUrusan string `json:"kode_bidang_urusan"`
	NamaBidangUrusan string `json:"nama_bidang_urusan"`
	KodeOpd          string `json:"kode_opd"`
	NamaOpd          string `json:"nama_opd"`
}

type IkkTerpilih struct {
	Id             int
	PohonKinerjaId int
	IkkId   	   int
}

type IkkTerpilihDetail struct {
	Id             int
	PohonKinerjaId int
	IkkId   	   int
	NamaPokin      string
	JenisIkk       string
	KeteranganIkk  string
}