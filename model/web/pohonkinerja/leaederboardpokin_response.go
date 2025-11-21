package pohonkinerja

type LeaderboardPokinResponse struct {
	KodeOpd             string                   `json:"kode_opd"`
	NamaOpd             string                   `json:"nama_opd"`
	Tematik             []LeaderboardTematikItem `json:"tematik"`
	PersentaseCascading string                   `json:"persentase_cascading"`
}

type LeaderboardTematikItem struct {
	Nama string `json:"nama"`
}
