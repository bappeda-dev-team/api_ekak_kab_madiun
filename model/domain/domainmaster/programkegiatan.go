package domainmaster

import "ekak_kabupaten_madiun/model/domain"

type ProgramKegiatan struct {
	Id              string
	KodeProgram     string
	NamaProgram     string
	KodeSubKegiatan string
	KodeOPD         string
	IsActive        bool
	Tahun           string
	Indikator       []domain.Indikator
}
