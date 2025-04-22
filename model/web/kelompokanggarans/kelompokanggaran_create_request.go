package kelompokanggarans

type KelompokAnggaranCreateRequest struct {
	Tahun        string `json:"tahun" validate:"required"`
	Kelompok     string `json:"kelompok" validate:"required"`
	KodeKelompok string `json:"kode_kelompok" `
}
