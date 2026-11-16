package domainmaster

import "ekak_kabupaten_madiun/model/domain"

// MasterIndikatorProgramKegiatan — baris di tb_master_indikator_program_kegiatan (tanpa target di API master).
type MasterIndikatorProgramKegiatan struct {
	Id            int
	KodeIndikator string
	Kode          string
	Indikator     string
}

type ProgramKegiatan struct {
	Id          string
	KodeProgram string
	NamaProgram string
	KodeOPD     string
	IsActive    bool
	Tahun       string
	Indikator   []domain.Indikator
}
