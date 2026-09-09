package domain

import "time"

type JenisInovasi struct {
	ID        int
	Jenis     string
	CreatedAt time.Time
	UpdatedAt time.Time
}