package sasaranpemda

type LockDataPemdaResponse struct {
	Id     int    `json:"id,omitempty"`
	Jenis  string `json:"jenis"`
	Tahun  string `json:"tahun"`
	Locked bool   `json:"locked"`
}
