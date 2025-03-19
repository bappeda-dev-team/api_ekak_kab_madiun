package domain

import "time"

type SubKegiatan struct {
	Id              string
	KodeSubKegiatan string
	NamaSubKegiatan string
	// KodeOpd              string
	NamaOpd string
	// Tahun                string
	RekinId string
	// Status               string
	CreatedAt            time.Time
	Indikator            []Indikator
	IndikatorSubKegiatan []IndikatorSubKegiatan
	PaguSubKegiatan      []PaguSubKegiatan
}

type IndikatorSubKegiatan struct {
	Id            string
	SubKegiatanId string
	NamaIndikator string
}

type PaguSubKegiatan struct {
	Id            string
	SubKegiatanId string
	JenisPagu     string
	PaguAnggaran  int
	Tahun         string
}
