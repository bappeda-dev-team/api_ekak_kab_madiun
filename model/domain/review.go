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

type ReviewTematik struct {
	IdTematik  int
	NamaPohon  string
	LevelPohon int
	Review     []ReviewDetail
}

type ReviewDetail struct {
	IdPohon    int
	Parent     int
	NamaPohon  string
	LevelPohon int
	JenisPohon string
	Review     string
	Keterangan string
	CreatedBy  string
	JenisPokin string
	CreatedAt  string
	UpdatedAt  string
}

type ReviewOpd struct {
	IdPohon    int
	Parent     int
	NamaPohon  string
	LevelPohon int
	JenisPohon string
	Review     string
	Keterangan string
	CreatedBy  string
	CreatedAt  string
	UpdatedAt  string
}
