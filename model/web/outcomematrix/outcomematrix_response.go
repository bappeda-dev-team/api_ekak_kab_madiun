package outcomematrix

import "time"

type OutcomeMatrixResponse struct {
	Id              int       `json:"id"`
	KodeSubkegiatan string    `json:"kode_subkegiatan"`
	Kode            string    `json:"kode"`
	Outcome         string    `json:"outcome"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
