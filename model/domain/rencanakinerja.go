package domain

import (
	"time"
)

type RencanaKinerja struct {
	Id                   string
	IdPohon              int
	LevelPohon           int
	ParentPohon          int
	NamaPohon            string
	NamaRencanaKinerja   string
	Tahun                string
	StatusRencanaKinerja string
	Catatan              string
	KodeOpd              string
	NamaOpd              string
	PegawaiId            string
	NamaPegawai          string
	TahunAwal            string
	TahunAkhir           string
	JenisPeriode         string
	PeriodeId            int
	CreatedAt            time.Time
	Indikator            []Indikator
	//tambahan
	Formula            string
	SumberData         string
	KodeSubKegiatan    string
	NamaSubKegiatan    string
	KodeKegiatan       string
	NamaKegiatan       string
	SasaranOpdId       int
	NamaSasaranOpd     string
	PohonKinerjaParent PohonKinerja
	KodeProgram        string
	Program            string
	PaguProgram        string
}
