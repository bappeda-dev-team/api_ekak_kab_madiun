package kelompokanggarans

type KelompokAnggaranResponse struct {
	Id           int    `json:"id"`
	Tahun        string `json:"tahun"`
	Kelompok     string `json:"kelompok"`
	KodeKelompok string `json:"kode_kelompok"`
	TahunView    string `json:"tahun_view,omitempty"`
}
