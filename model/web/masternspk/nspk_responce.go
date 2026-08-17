package masternspk

import "time"

type NspkResponse struct {
	ID               int       `json:"id"`
	KodeOpd          string    `json:"kode_opd"`
	Nspk             string    `json:"nspk"`
	Tahun            int       `json:"tahun"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type NspkFullResponse struct {
	ID               int       `json:"id"`
	KodeOpd          string    `json:"kode_opd"`
	NamaOpd          string    `json:"nama_opd"`
	Nspk             string    `json:"nspk"`
	Tahun            int       `json:"tahun"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}