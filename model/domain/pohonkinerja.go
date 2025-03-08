package domain

import "time"

type PohonKinerja struct {
	IdCrosscutting         int
	Id                     int
	Parent                 int
	NamaPohon              string
	KodeOpd                string
	NamaOpd                string
	Keterangan             string
	KeteranganCrosscutting *string
	Tahun                  string
	JenisPohon             string
	LevelPohon             int
	CreatedAt              time.Time
	Indikator              []Indikator
	Pelaksana              []PelaksanaPokin
	Status                 string
	CloneFrom              int
	Crosscutting           []Crosscutting
	PegawaiAction          interface{}
	CrosscuttingTo         int
	CountReview            int
	IsActive               bool
}

type PegawaiAction struct {
	ApproveBy *string
	RejectBy  *string
	ApproveAt *time.Time
	RejectAt  *time.Time
}
