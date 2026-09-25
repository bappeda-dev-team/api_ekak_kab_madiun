package inovasirekin

import (
	"ekak_kabupaten_madiun/model/web"
)

type InovasiRekinResponse struct {
	Id           	   string             `json:"id"`
	RekinId      	   string             `json:"rencana_kinerja_id"`
	KodeOpd      	   string             `json:"kode_opd"`
	NamaOpd      	   string             `json:"nama_opd"`
	NamaInovasi        string			  `json:"nama_inovasi"`
	JenisInovasiId     int				  `json:"jenis_inovasi_id"`
	JenisInovasi       string			  `json:"jenis_inovasi"`
	WaktuImplementasi  string	          `json:"waktu_implementasi"`
	Instansi           string			  `json:"instansi"`
	Inovator           string		      `json:"inovator"`
	Kebaruan           string             `json:"kebaruan"`
	AsalInovasi        string             `json:"asal_inovasi"`
	Tahun              int                `json:"tahun"`
	NipInovator        string             `json:"nip_inovator"`
	NamaNipInovator    string             `json:"nama_nip_inovator"`
	Level              string             `json:"level"`
	Action       	   []web.ActionButton `json:"action,omitempty"`
}
type InovasiLaporanResponse struct {
	Id           	   string             `json:"id"`
	RekinId      	   string             `json:"rencana_kinerja_id"`
	KodeOpd      	   string             `json:"kode_opd"`
	NamaOpd      	   string             `json:"nama_opd"`
	NamaInovasi        string			  `json:"nama_inovasi"`
	JenisInovasiId     int				  `json:"jenis_inovasi_id"`
	JenisInovasi       string			  `json:"jenis_inovasi"`
	WaktuImplementasi  string	          `json:"waktu_implementasi"`
	Instansi           string			  `json:"instansi"`
	Inovator           string		      `json:"inovator"`
	Kebaruan           string             `json:"kebaruan"`
	AsalInovasi        string             `json:"asal_inovasi"`
	Tahun              int                `json:"tahun"`
	NipInovator        string             `json:"nip_inovator"`
	NamaNipInovator    string             `json:"nama_nip_inovator"`
	Level              string             `json:"level"`
	Action       	   []web.ActionButton `json:"action,omitempty"`
}

