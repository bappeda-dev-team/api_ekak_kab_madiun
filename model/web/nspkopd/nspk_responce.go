package nspkopd

import "time"

type NspkResponse struct {
	ID        		int       `json:"id"`
	KodeOpd   		string    `json:"kode_opd"`
	IdNspk      	int    	  `json:"id_nspk"`
	IdTujuanOpd     int    	  `json:"id_tujuan_opd"`
	IdSasaranOpd    string    `json:"id_sasaran_opd"`
	Tahun     		int       `json:"tahun"`
	CreatedAt 		time.Time `json:"created_at"`
	UpdatedAt 		time.Time `json:"updated_at"`
}

type NspkFullResponse struct {
	ID        	   int       `json:"id"`
	KodeOpd   	   string    `json:"kode_opd"`
	NamaOpd   	   string    `json:"nama_opd"`
	IdNspk    	   int       `json:"id_nspk"`
	Nspk      	   string    `json:"nspk"`
	IdTujuanOpd    int       `json:"id_tujuan_opd"`
	TujuanOpd      string    `json:"tujuan_opd"`
	IdSasaranOpd   string    `json:"id_sasaran_opd"`
	SasaranOpd     string    `json:"sasaran_opd"`
	Tahun     	   int       `json:"tahun"`
	CreatedAt 	   time.Time `json:"created_at"`
	UpdatedAt 	   time.Time `json:"updated_at"`
}