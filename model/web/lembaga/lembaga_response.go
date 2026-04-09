package lembaga

type LembagaResponse struct {
	Id                 string `json:"id"`
	KodeLembaga        string `json:"kode_lembaga"`
	NamaLembaga        string `json:"nama_lembaga"`
	NamaKepalaPemda    string `json:"nama_kepala_pemda"`
	JabatanKepalaPemda string `json:"jabatan_kepala_pemda"`
	NipKepalaPemda     string `json:"nip_kepala_pemda"`
	IsActive           bool   `json:"is_active"`
}
