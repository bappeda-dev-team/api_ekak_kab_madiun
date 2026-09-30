package inovasirekin

type InovasiRekinUpdateRequest struct {
	Id                string `json:"id"`
	NamaInovasi       string `json:"nama_inovasi" validate:"required"`
	JenisInovasiId    int    `json:"jenis_inovasi_id" validate:"required"`
	WaktuImplementasi string `json:"waktu_implementasi" validate:"required"`
	Instansi          string `json:"instansi"`
	Inovator          string `json:"inovator"`
	Kebaruan          string `json:"kebaruan" validate:"required"`
	AsalInovasi       string `json:"asal_inovasi" validate:"required"`
	NipInovator       string `json:"nip_inovator"`
}
