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
}
