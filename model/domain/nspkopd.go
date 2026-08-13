package domain

import "time"

type NspkOpd struct {
	ID               int
	KodeOpd          string
	NamaOpd          string
	IdNspk           int
	NSPK             string
	IdTujuanOpd      int
	TujuanOpd        string
	IdSasaranOpd     string
	SasaranOpd       string
	Tahun            int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}