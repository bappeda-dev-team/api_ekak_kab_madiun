package domain

import "time"

type InovasiRekin struct {
	Id                 string
	RekinId            string
	KodeOpd            string
	NamaInovasi        string
	JenisInovasiId     int
	JenisInovasi       string
	WaktuImplementasi  time.Time
	Instansi           string
	Inovator           string
	CreatedAt          time.Time
}
