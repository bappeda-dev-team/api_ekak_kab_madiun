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

type WarningType string

type CloneWarning struct {
	Type    WarningType // "indikator", "manual_ik", "renaksi", dll
	Count   int
	Message string
}

type CloneResult struct {
	Warnings []CloneWarning
}
