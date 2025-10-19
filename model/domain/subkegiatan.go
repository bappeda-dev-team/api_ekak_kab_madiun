package domain

import (
	"database/sql"
	"time"
)

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

type SubKegiatanQuery struct {
	KodeUrusan       string
	NamaUrusan       string
	KodeBidangUrusan string
	NamaBidangUrusan string
	KodeProgram      string
	NamaProgram      string
	KodeKegiatan     string
	NamaKegiatan     string
	KodeSubKegiatan  string
	NamaSubKegiatan  string
	TahunSubKegiatan string
	PegawaiId        string
	IndikatorId      string
	IndikatorKode    string
	Indikator        string
	IndikatorTahun   string
	IndikatorKodeOpd string
	Target           string
	Satuan           string
	TargetId         string
	PaguAnggaran     sql.NullInt64
}

type SubKegiatanKAKQuery struct {
	KodeOpd              string
	NamaOpd              string
	KodeProgram          string
	NamaProgram          string
	IndikatorProgram     string
	TargetProgram        string
	SatuanProgram        string
	KodeKegiatan         string
	NamaKegiatan         string
	IndikatorKegiatan    string
	TargetKegiatan       string
	SatuanKegiatan       string
	KodeSubKegiatan      string
	NamaSubKegiatan      string
	IndikatorSubKegiatan string
	TargetSubKegiatan    string
	SatuanSubKegiatan    string
	PaguAnggaran         int64
	TahunIndikator       string
}
