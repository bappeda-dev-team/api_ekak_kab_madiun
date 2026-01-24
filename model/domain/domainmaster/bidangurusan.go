package domainmaster

import "time"

type BidangUrusan struct {
	Id               string
	KodeBidangUrusan string
	NamaBidangUrusan string
	KodeUrusan       string
	NamaUrusan       string
	Tahun            string
	CreatedAt        time.Time
}

type BidangUrusanOpd struct {
	Id               int
	KodeBidangUrusan string
	NamaBidangUrusan string
	KodeOpd          string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
