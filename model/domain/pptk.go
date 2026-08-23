package domain

import "time"

type Pptk struct {
	Id                              int
	Nip                     		string
	KodeOpd       					string
	Tahun                           int
	KodeSubKegiatan 				string
	NipAtasan                       *string
	AktifAt                      	time.Time
	NonAktifAt                      *time.Time
	CreatedAt                       time.Time
	UpdatedAt                       time.Time
}