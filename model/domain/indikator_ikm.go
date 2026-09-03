package domain

import (
	"time"
)

type IndikatorIkm struct {
	Id                  string
	Indikator           string
	KodeBidangUrusan    string
	NamaBidangUrusan    string
	IsActive            bool
	DefinisiOperasional string
	RumusPerhitungan    string
	SumberData          string
	Jenis               string
	TahunAwal           string
	TahunAkhir          string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
