package domain

import "time"

type PohonKinerja struct {
	Id            int
	Parent        int
	NamaPohon     string
	KodeOpd       string
	NamaOpd       string
	Keterangan    string
	Tahun         string
	JenisPohon    string
	LevelPohon    int
	CreatedAt     time.Time
	Indikator     []Indikator
	Pelaksana     []PelaksanaPokin
	Status        string
	CloneFrom     int
	Crosscutting  []Crosscutting
	PegawaiAction interface{}
}

type PegawaiAction struct {
	ApproveBy *string
	RejectBy  *string
	ApproveAt *time.Time
	RejectAt  *time.Time
}
