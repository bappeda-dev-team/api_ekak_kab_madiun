package domain

import "time"

type StrategicRow struct {
	KodeOpd          string
	NamaTujuanOpd    string
	NamaSasaranOpd   string
	NamaStrategi     string // level 4
	IdTactical       int    // level 5
	NamaTactical     string // level 5
	NamaOperasional  string // level 6
	TahunOperasional int // level 6
	ArahKebijakan    ArahKebijakanRow
}

type ArahKebijakanRow struct {
	ID           int
	KodeOpd      string
	PokinId      int
	Arah         string
}
type ArahKebijakanOpd struct {
	ID           int
	KodeOpd      string
	NamaOpd      string
	PokinId      int
	Arah         string
	Tahun        int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type StrategicPemdaRow struct {
	NamaTujuanPemda   string
	NamaSasaranPemda  string
	NamaStrategi      string // level 1 & 2
	NamaArahKebijakan string // level 4
}
