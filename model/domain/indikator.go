package domain

import (
	"database/sql"
	"time"
)

type Indikator struct {
	Id               string
	PokinId          string
	ProgramId        string
	RencanaKinerjaId string
	KegiatanId       string
	SubKegiatanId    string
	TujuanOpdId      int
	TujuanPemdaId    int
	SasaranPemdaId   int
	SasaranOpdId     int
	Indikator        string
	Tahun            string
	CloneFrom        string
	Sumber           string
	ParentId         int
	ParentOpdId      string
	AsalIku          string
	ParentName       string
	CreatedAt        time.Time
	TahunAwal        string
	TahunAkhir       string
	JenisPeriode     string
	Target           []Target
	RencanaKinerja   RencanaKinerja
	RumusPerhitungan sql.NullString
	SumberData       sql.NullString
	IsActive         bool
	Kode             string
	KodeOpd          string
	PaguAnggaran     int64
}
