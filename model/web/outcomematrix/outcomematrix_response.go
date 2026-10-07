package outcomematrix

import "time"

type OutcomeMatrixResponse struct {
	Id        int       `json:"id"`
	KodeOpd   string    `json:"kode_opd"`
	Kode      string    `json:"kode"`
	Jenis     string    `json:"jenis"`
	Outcome   string    `json:"outcome"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
