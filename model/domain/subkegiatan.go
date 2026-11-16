package domain

import (
	"database/sql"
	"time"
)

type SubKegiatan struct {
	Id              string
	KodeSubKegiatan string
	NamaSubKegiatan string
	// KodeOpd              string
	NamaOpd string
	// Tahun                string
	RekinId string
	// Status               string
	CreatedAt            time.Time
	Indikator            []Indikator
	IndikatorSubKegiatan []IndikatorSubKegiatan
	PaguSubKegiatan      []PaguSubKegiatan
}

type IndikatorSubKegiatan struct {
	Id            string
	SubKegiatanId string
	NamaIndikator string
}

type PaguSubKegiatan struct {
	Id            string
	SubKegiatanId string
	JenisPagu     string
	PaguAnggaran  int
	Tahun         string
}

type SubKegiatanQuery struct {
	KodeUrusan               string
	NamaUrusan               string
	KodeBidangUrusan         string
	NamaBidangUrusan         string
	KodeProgram              string
	NamaProgram              string
	KodeKegiatan             string
	NamaKegiatan             string
	KodeSubKegiatan          string
	NamaSubKegiatan          string
	TahunSubKegiatan         string
	PegawaiId                string
	NamaPegawai              string
	PaguSubKegiatan          int64
	IndikatorId              string
	IndikatorKode            string
	Indikator                string
	IndikatorTahun           string
	IndikatorKodeOpd         string
	Target                   string
	Satuan                   string
	TargetId                 string
	PaguAnggaran             sql.NullInt64
	TotalAnggaranSubKegiatan int64
	JenisPagu                string
	// penetapan: pagu dari tb_pagu jenis='penetapan' (0 jika belum ada)
	PaguPenetapan  int64
	// penetapan: jenis layer indikator ("rankhir" atau "penetapan")
	JenisIndikator string

	// flat pagu per layer (ranwal & rankhir) — digunakan oleh GET ranwal/rankhir/penetapan v2
	PaguRanwal  int64
	PaguRankhir int64

	// dual-target untuk rankhir view (slot ranwal + slot rankhir per indikator)
	TargetRanwalId  string
	TargetRanwal    string
	SatuanRanwal    string
	TargetRankhirId string
	TargetRankhir   string
	SatuanRankhir   string

	// slot penetapan (untuk penetapan view)
	TargetPenetapanId string
	TargetPenetapan   string
	SatuanPenetapan   string

	// status subkegiatan (untuk rankhir view — baris sintetis dari renstra yg tidak ada rekin tahun ini)
	StatusSubkegiatan string
}

type SubKegiatanKAKQuery struct {
	KodeOpd              string
	NamaOpd              string
	KodeProgram          string
	NamaProgram          string
	IndikatorProgram     string
	TargetProgram        string
	SatuanProgram        string
	KodeKegiatan         string
	NamaKegiatan         string
	IndikatorKegiatan    string
	TargetKegiatan       string
	SatuanKegiatan       string
	KodeSubKegiatan      string
	NamaSubKegiatan      string
	IndikatorSubKegiatan string
	TargetSubKegiatan    string
	SatuanSubKegiatan    string
	PaguAnggaran         int64
	TahunIndikator       string
}
