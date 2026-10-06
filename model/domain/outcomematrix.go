package domain

import "time"

type OutcomeMatrix struct {
	Id              int
	KodeSubkegiatan string
	Kode            string
	Outcome         string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
