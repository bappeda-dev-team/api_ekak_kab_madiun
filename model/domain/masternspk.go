package domain

import "time"

type MasterNSPK struct {
	ID               int
	KodeOpd          string
	NamaOpd          string
	NSPK             string
	Tahun            int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}