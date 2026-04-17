package pohonkinerja

type LeaderboardHiddenUpsertRequest struct {
	KodeOpd  string `json:"kode_opd" validate:"required"`
	Tahun    string `json:"tahun" validate:"required,len=4"`
	IsHidden bool   `json:"is_hidden"`
}
