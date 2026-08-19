package strategic

type StrategicArahKebijakanOpdAllResponse struct {
	KodeOpd                   string                             `json:"kode_opd"`
	NamaOpd                   string                             `json:"nama_opd"`
	Tahun                     string                             `json:"tahun"`
	PermasalahanOpd           []PermasalahanOpdResponse          `json:"permasalahan_opds"`
	IsuStrategisOpd           []IsuStrategiOpdResponse           `json:"isu_strategis_opds"`
	StrategiArahKebijakanOpds []StrategiArahKebijakanOpdResponse `json:"strategi_arah_kebijakan_opds"`
}

type TujuanOpdResponse struct {
	KodeOpd string `json:"kode_opd"`
	Tujuan  string `json:"tujuan"`
}

type IsuStrategiOpdResponse struct {
	NamaIsu string `json:"nama_isu_strategis"`
}

type PermasalahanOpdResponse struct {
	NamaPermasalahan string `json:"permasalahan"`
}

type StrategiArahKebijakanOpdResponse struct {
	TujuanOpd   string               `json:"tujuan_opd"`
	SasaranOpds []SasaranOpdResponse `json:"sasaran_opds"`
}

// type SasaranOpdResponse struct {
// 	SasaranOpd        string                     `json:"sasaran_opd"`
// 	StrategiOpd       string                     `json:"strategi_opd"`
// 	ArahKebijakanOpds []ArahKebijakanOpdResponse `json:"arah_kebijakan_opds"`
// }

// type ArahKebijakanOpdResponse struct {
// 	ArahKebijakanOpd string `json:"arah_kebijakan_opd"`
// }

type SasaranOpdResponse struct {
	SasaranOpd   string                `json:"sasaran_opd"`
	StrategiOpds []StrategiOpdResponse `json:"strategi_opds"`
}

type StrategiOpdResponse struct {
	StrategiOpd  string                `json:"strategi_opd"`
	TacticalOpds []TacticalOpdResponse `json:"tactical_opds"`
}

type TacticalOpdResponse struct {
	TacticalOpd     string                   `json:"tactical_opd"`
	OperasionalOpds []OperasionalOpdResponse `json:"operasional_opds"`
}

type OperasionalOpdResponse struct {
	OperasionalOpd string `json:"operasional_opd"`
}
