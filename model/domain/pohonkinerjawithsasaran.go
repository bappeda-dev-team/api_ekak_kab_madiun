package domain

import "database/sql"

type PohonKinerjaWithSasaran struct {
	TematikId   int
	NamaTematik string
	Tahun       string
	Subtematik  []SubtematikWithSasaran
}

type SubtematikWithSasaran struct {
	SubtematikId     int
	NamaSubtematik   string
	JenisPohon       string
	LevelPohon       int
	Keterangan       string
	IsActive         bool
	SasaranPemdaList []SasaranPemdaDetail
}

type SasaranPemdaDetail struct {
	Id           int
	SasaranPemda string
	PeriodeId    int
	Periode      Periode
	Indikator    []IndikatorDetail
}

type IndikatorDetail struct {
	Id               string
	Indikator        string
	RumusPerhitungan sql.NullString
	SumberData       sql.NullString
	Target           []TargetDetail
}

type TargetDetail struct {
	Id          string
	IndikatorId string
	Target      string
	Satuan      string
	Tahun       string
}
