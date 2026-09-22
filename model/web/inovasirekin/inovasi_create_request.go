package inovasirekin

type InovasiRekinCreateRequest struct {
	RekinId           string `json:"rencana_kinerja_id" validate:"required"`
	KodeOpd           string `json:"kode_opd" validate:"required"`
	NamaInovasi       string `json:"nama_inovasi" validate:"required"`
	JenisInovasiId    int    `json:"jenis_inovasi_id" validate:"required"`
	WaktuImplementasi string `json:"waktu_implementasi" validate:"required"`
	Instansi          string `json:"instansi"`
	Inovator          string `json:"inovator"`
	Kebaruan          string `json:"kebaruan" validate:"required"`
	AsalInovasi       string `json:"asal_inovasi" validate:"required"`
	Tahun             int    `json:"tahun" validate:"required"`
	NipInovator       string `json:"nip_inovator"`
}