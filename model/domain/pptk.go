package domain

import "time"

type Pptk struct {
	Id                              int
	Nip                     		string
	NamaPegawai                     string
	KodeOpd       					string
	Tahun                           int
	KodeSubKegiatan 				string
	NipAtasan                       *string
	NamaAtasan                      *string
	AktifAt                      	time.Time
	NonAktifAt                      *time.Time
	CreatedAt                       time.Time
	UpdatedAt                       time.Time
}

type KandidatPptk struct {
	PegawaiId   string
	NamaPegawai string
	Level       string
}