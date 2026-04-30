package ikk

type IkkRequest struct {
	KodeBidangUrusan string `json:"kode_bidang_urusan" validate:"required"`
	Jenis            string `json:"jenis" validate:"required"`
	NamaIndikator    string `json:"nama_indikator" validate:"required"`
	Target           string `json:"target" validate:"required"`
	Satuan           string `json:"satuan" validate:"required"`
	Keterangan       string `json:"keterangan"`
}