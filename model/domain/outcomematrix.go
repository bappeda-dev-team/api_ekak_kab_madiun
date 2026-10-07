package domain

import "time"

type OutcomeMatrix struct {
	Id        int
	KodeOpd   string
	Kode      string
	Jenis     string
	Outcome   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
