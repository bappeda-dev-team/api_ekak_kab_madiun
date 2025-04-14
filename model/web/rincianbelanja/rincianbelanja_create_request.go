package rincianbelanja

type RincianBelanjaCreateRequest struct {
	RenaksiId string `json:"renaksi_id"`
	Anggaran  int64  `json:"anggaran"`
}
