package domain

import "time"

type CloneRecord struct {
	Id                   int
	KodeClone            string
	KodeOpd              string
	TahunAsal            string
	TahunTarget          string
	KeteranganTahunClone string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	UpdatedBy            string
	Status               string
	ErrorMessage         string
}
