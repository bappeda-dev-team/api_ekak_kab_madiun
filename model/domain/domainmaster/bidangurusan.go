package domainmaster

import "time"

type BidangUrusan struct {
	Id               string
	KodeBidangUrusan string
	NamaBidangUrusan string
	KodeUrusan       string
	NamaUrusan       string
	Tahun            string
	KodeSubKegiatan string
	CreatedAt        time.Time
}
