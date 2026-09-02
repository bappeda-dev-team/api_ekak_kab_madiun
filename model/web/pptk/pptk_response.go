package pptk

import "time"

type PptkResponse struct {
	Id              int        `json:"id"`
	Nip             string     `json:"nip"`
	NamaPegawai     string     `json:"nama_pegawai"`
	KodeOpd         string     `json:"kode_opd"`
	Tahun           int        `json:"tahun"`
	KodeSubKegiatan string     `json:"kode_sub_kegiatan"`
	NipAtasan       *string    `json:"nip_atasan"`
	NamaAtasan      *string    `json:"nama_atasan"`
	AktifAt         time.Time  `json:"aktif_at"`
	NonAktifAt      *time.Time `json:"nonaktif_at"`
}