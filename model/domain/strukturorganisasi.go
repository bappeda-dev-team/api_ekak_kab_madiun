package domain

import (
	"time"
)

type StrukturOrganisasi struct {
	Id         int
	NipBawahan string
	NipAtasan  string
	KodeOpd    string
	Tahun      int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
