package domain

type SasaranPemda struct {
	Id            int
	TujuanPemdaId int
	SubtemaId     int
	IsActive      bool
	NamaSubtema   string
	SasaranPemda  string
	JenisPohon    string
	PeriodeId     int
	TahunAwal     string
	TahunAkhir    string
	JenisPeriode  string
	Periode       Periode
	Indikator     []Indikator
}

type SasaranPemdaWithPokin struct {
	SubtematikId        int
	JenisPohon          string
	LevelPohon          int
	IsActive            bool
	TematikId           int
	NamaTematik         string
	NamaSubtematik      string
	IdsasaranPemda      int
	SasaranPemda        string
	Keterangan          string
	IndikatorSubtematik []Indikator
	SasaranPemdaList    []SasaranPemda
}
