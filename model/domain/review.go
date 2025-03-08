package domain

import "time"

type Review struct {
	Id             int
	IdPohonKinerja int
	Review         string
	Keterangan     string
	CreatedBy      string
	Jenis_pokin    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
