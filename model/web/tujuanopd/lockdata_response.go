package tujuanopd

// LockDataOpdResponse digunakan sebagai respons lock/unlock tujuan OPD.
type LockDataOpdResponse struct {
	Id      int    `json:"id,omitempty"`
	Jenis   string `json:"jenis"`
	KodeOpd string `json:"kode_opd"`
	Tahun   string `json:"tahun"`
	Locked  bool   `json:"locked"`
}
