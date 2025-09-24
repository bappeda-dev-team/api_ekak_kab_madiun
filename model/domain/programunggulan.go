package domain

import "time"

type ProgramUnggulan struct {
	Id                        int
	NamaTagging               string
	KodeProgramUnggulan       string
	KeteranganProgramUnggulan *string
	Keterangan                *string
	TahunAwal                 string
	TahunAkhir                string
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}
