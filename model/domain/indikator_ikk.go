package domain

import "time"

type IndikatorIkk struct {
	ID               int   
	KodeBidangUrusan string 
	NamaOpd 		 string 
	Jenis            string 
	NamaIndikator    string 
	Target           string 
	Satuan           string 
	Keterangan       string 
	CreatedAt        time.Time
	UpdatedAt        time.Time
}