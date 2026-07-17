package isuglobal

import "time"

type IsuGlobalResponse struct {
	ID               int                 `json:"id"`
	KodeBidangUrusan string              `json:"kode_bidang_urusan"`
	KodeOpd          string              `json:"kode_opd"`
	Isu              string              `json:"isu"`
	Tahun            int                 `json:"tahun"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
}