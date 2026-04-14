package strategic

type StrategicArahKebijakanOpdAllResponse struct {
	KodeOpd                   string                          `json:"kode_opd"`
	NamaOpd                   string                          `json:"nama_opd"`
	Tahun                     string                          `json:"tahun"`
	IsuStrategisOpd           []IsuStrategiOpdResponse        `json:"isu_strategis_opds"`
	TujuanOpd                 []TujuanOpdResponse             `json:"tujuan_opd"`
	StrategiArahKebijakanOpds []StrategiArahKebijakanOpdResponse `json:"strategi_arah_kebijakan_opds"`
}

type TujuanOpdResponse struct {
	Id      int    `json:"id"`
	KodeOpd string `json:"kode_opd"`
	Tujuan  string `json:"tujuan"`
}

type IsuStrategiOpdResponse struct {
	NamaIsu string `json:"nama_isu_strategis"`
}

type StrategiArahKebijakanOpdResponse struct {
	TujuanOpd string `json:"tujuan_opd"`
	SasaranOpds []SasaranOpdResponse `json:"sasaran_opds"`
}

type SasaranOpdResponse struct {
	SasaranOpd string `json:"sasaran_opd"`
	StrategiOpd string `json:"strategi_opd"`
	ArahKebijakanOpds []ArahKebijakanOpdResponse `json:"arah_kebijakan_opds"`
}

type ArahKebijakanOpdResponse struct {
	ArahKebijakanOpd string `json:"arah_kebijakan_opd"`
}

