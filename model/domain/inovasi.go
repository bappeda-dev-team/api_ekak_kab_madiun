package domain

import "time"

type Inovasi struct {
	Id                    string
	RekinId               string
	JudulInovasi          string
	JenisInovasi          string
	GambaranNilaiKebaruan string
	PegawaiId             string
	KodeOpd               string
	CreatedAt             time.Time
}
