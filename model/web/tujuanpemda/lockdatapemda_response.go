package tujuanpemda

type LockDataPemdaResponse struct {
	Id     int    `json:"id"`
	Jenis  string `json:"jenis"`
	Tahun  string `json:"tahun"`
	Locked bool   `json:"locked"`
}
