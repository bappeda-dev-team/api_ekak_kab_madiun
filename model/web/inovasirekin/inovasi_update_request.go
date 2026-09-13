package inovasirekin

import "time"

type InovasiRekinUpdateRequest struct {
	Id                string    `json:"id"`
	NamaInovasi       string    `json:"nama_inovasi"`
	JenisInovasiId    int       `json:"jenis_inovasi_id"`
	WaktuImplementasi time.Time `json:"waktu_implementasi"`
	Instansi          string    `json:"instansi"`
	Inovator          string    `json:"inovator"`
}
