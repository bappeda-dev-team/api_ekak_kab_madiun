package domain

import "time"

type TaggingPokin struct {
	Id                       int
	IdPokin                  int
	NamaTagging              string
	CloneFrom                int
	CreatedAt                time.Time
	UpdatedAt                time.Time
	KeteranganTaggingProgram []KeteranganTagging
}

type KeteranganTagging struct {
	Id                  int
	IdTagging           int
	KodeProgramUnggulan string
	Tahun               string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
