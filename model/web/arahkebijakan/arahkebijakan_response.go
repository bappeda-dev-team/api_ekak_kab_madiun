package arahkebijakan

import "time"

type ArahKebijakanResponse struct {
	ID           int       `json:"id"`
	PokinId      int       `json:"pokin_id"`
	Arah         string    `json:"arah"`
	KodeOpd      string    `json:"kode_opd"`
	Tahun        int       `json:"tahun"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}