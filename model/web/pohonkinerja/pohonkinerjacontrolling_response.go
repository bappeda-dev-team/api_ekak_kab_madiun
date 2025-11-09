package pohonkinerja

type ControlPokinOpdResponse struct {
	Data  []ControlPokinOpdData `json:"data"`
	Total ControlPokinOpdTotal  `json:"total"`
}

type ControlPokinOpdData struct {
	LevelPohon                int    `json:"level_pohon"`
	NamaLevel                 string `json:"nama_level"`
	JumlahPokin               int    `json:"jumlah_pokin"`
	JumlahPelaksana           int    `json:"jumlah_pelaksana"`
	JumlahPokinTanpaPelaksana int    `json:"jumlah_pokin_tanpa_pelaksana"`
	Persentase                string `json:"persentase"`
}

type ControlPokinOpdTotal struct {
	TotalPokin               int    `json:"total_pokin"`
	TotalPelaksana           int    `json:"total_pelaksana"`
	TotalPokinTanpaPelaksana int    `json:"total_pokin_tanpa_pelaksana"`
	Persentase               string `json:"persentase"`
}
