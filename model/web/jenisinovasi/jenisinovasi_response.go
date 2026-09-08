package jenisinovasi

import "time"

type JenisInovasiResponse struct {
	ID        int       `json:"id"`
	Jenis     string    `json:"jenis"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
