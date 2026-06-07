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
	IsActive                  bool
	OpdList                   []OpdProgramUnggulan
	TahunTerpakai             []string
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

type OpdProgramUnggulan struct {
	Id                  int
	KodeProgramUnggulan string
	KodeOpd             string
	NamaOpd             string
}
