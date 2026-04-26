package strategicarahkebijakan

type StrategicArahKebijakanPemdaAllResponse struct {
	KodePemda                   string                               `json:"kode_pemda"`
	NamaPemda                   string                               `json:"nama_pemda"`
	Tahun                       string                               `json:"tahun"`
	IsuStrategisPemda           []IsuStrategiPemdaResponse           `json:"isu_strategis_pemdas"`
	TujuanPemda                 []TujuanPemdaResponse                `json:"tujuan_pemda"`
	StrategiArahKebijakanPemdas []StrategiArahKebijakanPemdaResponse `json:"strategi_arah_kebijakan_pemdas"`
}

type IsuStrategiPemdaResponse struct {
	NamaIsu string `json:"nama_isu_strategis"`
}

type TujuanPemdaResponse struct {
	Id        int    `json:"id"`
	KodePemda string `json:"kode_pemda"`
	Tujuan    string `json:"tujuan"`
}

type StrategiArahKebijakanPemdaResponse struct {
	TujuanPemda   string                 `json:"tujuan_pemda"`
	SasaranPemdas []SasaranPemdaResponse `json:"sasaran_pemdas"`
}

type SasaranPemdaResponse struct {
	SasaranPemda        string                       `json:"sasaran_pemda"`
	StrategiPemda       string                       `json:"strategi_pemda"`
	ArahKebijakanPemdas []ArahKebijakanPemdaResponse `json:"arah_kebijakan_pemdas"`
}

type ArahKebijakanPemdaResponse struct {
	ArahKebijakanPemda string `json:"arah_kebijakan_pemda"`
}