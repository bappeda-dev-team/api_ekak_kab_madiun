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
	JumlahPokinAdaPelaksana   int    `json:"jumlah_pokin_ada_pelaksana"`
	JumlahPokinTanpaPelaksana int    `json:"jumlah_pokin_tanpa_pelaksana"`
	JumlahRencanaKinerja      int    `json:"jumlah_rencana_kinerja"`   // ← BARU
	JumlahPokinAdaRekin       int    `json:"jumlah_pokin_ada_rekin"`   // ← BARU
	JumlahPokinTanpaRekin     int    `json:"jumlah_pokin_tanpa_rekin"` // ← BARU
	Persentase                string `json:"persentase"`               // pelaksana
	PersentaseCascading       string `json:"persentase_cascading"`     // ← BARU: cascading
}

type ControlPokinOpdTotal struct {
	TotalPokin               int    `json:"total_pokin"`
	TotalPelaksana           int    `json:"total_pelaksana"`
	TotalPokinAdaPelaksana   int    `json:"total_pokin_ada_pelaksana"`
	TotalPokinTanpaPelaksana int    `json:"total_pokin_tanpa_pelaksana"`
	TotalRencanaKinerja      int    `json:"total_rencana_kinerja"`   // ← BARU
	TotalPokinAdaRekin       int    `json:"total_pokin_ada_rekin"`   // ← BARU
	TotalPokinTanpaRekin     int    `json:"total_pokin_tanpa_rekin"` // ← BARU
	Persentase               string `json:"persentase"`              // pelaksana
	PersentaseCascading      string `json:"persentase_cascading"`    // ← BARU: cascading
}
