package domain

import (
	"time"
)

type PkOpd struct {
	Id               string
	KodeOpd          string
	NamaOpd          string
	LevelPk          int
	NipAtasan        string
	NamaAtasan       string
	IdRekinAtasan    string
	RekinAtasan      string
	NipPemilikPk     string
	NamaPemilikPk    string
	IdRekinPemilikPk string
	RekinPemilikPk   string
	Tahun            int
	Keterangan       string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	SasaranOpdId     int64
}

type AllItemPk struct {
	RekinId         string
	KodeProgram     string
	NamaProgram     string
	KodeKegiatan    string
	NamaKegiatan    string
	KodeSubkegiatan string
	NamaSubkegiatan string
	PaguSubkegiatan int64
}

type AllSasaranPemdaPk struct {
	JabatanKepalaPemda string
	NamaKepalaPemda    string
	NipKepalaPemda     string
	SasaranPemdaId     int
	SasaranPemda       string
}

type KunciPK struct {
	IdPegawai   string
	KodeOpd     string
	Tahun       int
	DikunciOleh string
	DikunciPada time.Time
	StatusPk    string // terkunci, terbuka, revisi
	PkTerkunci  bool   // true, false, false
}
