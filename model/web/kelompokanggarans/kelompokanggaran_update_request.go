package kelompokanggarans

type KelompokAnggaranUpdateRequest struct {
	Id           int    `json:"id" validate:"required"`
	Tahun        string `json:"tahun" validate:"required"`
	Kelompok     string `json:"kelompok" validate:"required"`
	KodeKelompok string `json:"kode_kelompok"`
}
