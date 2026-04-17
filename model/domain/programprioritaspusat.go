package domain

import "time"

type ProgramPrioritasPusat struct {
	Id                        		int
	NamaTagging               		string
	KodeProgramPrioritasPusat       string
	KeteranganProgramPrioritasPusat *string
	Keterangan                		*string
	TahunAwal                 		string
	TahunAkhir                		string
	IsActive                  		bool
	CreatedAt                 		time.Time
	UpdatedAt                 		time.Time
}