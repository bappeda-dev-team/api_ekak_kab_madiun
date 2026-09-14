package inovasirekin

import "time"

type InovasiRekinCreateRequest struct {
	RekinId           string 	  `json:"rencana_kinerja_id"`
	KodeOpd           string 	  `json:"kode_opd"`
	NamaInovasi       string 	  `json:"nama_inovasi"`
	JenisInovasiId    int   	  `json:"jenis_inovasi_id"`
	WaktuImplementasi time.Time   `json:"waktu_implementasi"`
	Instansi          string 	  `json:"instansi"`
	Inovator          string 	  `json:"inovator"`
}