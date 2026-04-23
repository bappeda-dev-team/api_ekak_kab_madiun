package pohonkinerja

type LeaderboardPokinResponse struct {
	KodeOpd             string                   `json:"kode_opd"`
	NamaOpd             string                   `json:"nama_opd"`
	Tematik             []LeaderboardTematikItem `json:"tematik"`
	PersentaseCascading string                   `json:"persentase_cascading"`
	IsHidden            bool                     `json:"is_hidden"`
}

type LeaderboardTematikItem struct {
	Nama  string                   `json:"nama"`
	Anak  []LeaderboardTematikItem `json:"child"`
	Level int                      `json:"level"`
	Jenis string                      `json:"jenis"`
}
