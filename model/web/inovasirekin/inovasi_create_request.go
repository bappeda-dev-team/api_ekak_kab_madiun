package inovasirekin

type InovasiRekinCreateRequest struct {
	RekinId           string `json:"rencana_kinerja_id" validate:"required"`
	KodeOpd           string `json:"kode_opd" validate:"required"`
	NamaInovasi       string `json:"nama_inovasi" validate:"required"`
	JenisInovasiId    int    `json:"jenis_inovasi_id" validate:"required"`
	WaktuImplementasi string `json:"waktu_implementasi" validate:"required"`
	Instansi          string `json:"instansi"`
	Inovator          string `json:"inovator" validate:"required"`
}