package domain

type StrategicRow struct {
	KodeOpd         string
	NamaTujuanOpd   string
	NamaSasaranOpd  string
	NamaStrategi    string // level 4
	NamaTactical    string // level 5
	NamaOperasional string // level 6
}

type StrategicPemdaRow struct {
	NamaTujuanPemda   string
	NamaSasaranPemda  string
	NamaStrategi      string // level 1 & 2
	NamaArahKebijakan string // level 4
}
