package lembaga

type LembagaCreateRequest struct {
	KodeLembaga        string `json:"kode_lembaga" validate:"required"`
	NamaLembaga        string `json:"nama_lembaga" validate:"required"`
	NamaKepalaPemda    string `json:"nama_kepala_pemda"`
	JabatanKepalaPemda string `json:"jabatan_kepala_pemda"`
	NipKepalaPemda     string `json:"nip_kepala_pemda"`
}
