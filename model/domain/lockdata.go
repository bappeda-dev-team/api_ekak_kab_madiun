package domain

import "time"

type LockData struct {
	Id        int
	JenisData string
	KodeOpd   string
	Tahun     string
}

type LockDataPemda struct {
	Id        int
	Jenis     string
	Tahun     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
