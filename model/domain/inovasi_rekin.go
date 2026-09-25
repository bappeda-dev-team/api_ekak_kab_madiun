package domain

import "time"

type InovasiRekin struct {
	Id                 string
	RekinId            string
	KodeOpd            string
	NamaOpd            string
	NamaInovasi        string
	JenisInovasiId     int
	JenisInovasi       string
	WaktuImplementasi  string
	Instansi           string
	Inovator           string
	Kebaruan           string
	AsalInovasi        string
	Tahun              int
	NipInovator        string
	NamaNipInovator    string
	Level              string
	CreatedAt          time.Time
}

type InovasiLaporan struct {
	Id                 string
	RekinId            string
	KodeOpd            string
	NamaOpd            string
	NamaInovasi        string
	JenisInovasiId     int
	JenisInovasi       string
	WaktuImplementasi  string
	Instansi           string
	Inovator           string
	Kebaruan           string
	AsalInovasi        string
	Tahun              int
	NipInovator        string
	NamaNipInovator    string
	Level              string
	CreatedAt          time.Time
}
