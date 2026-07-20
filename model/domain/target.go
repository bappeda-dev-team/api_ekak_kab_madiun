package domain

type Target struct {
	Id          string
	IndikatorId string
	Target      string
	Satuan      string
	Tahun       string
	CloneFrom   string
	Jenis       string
}

type TargetPemda struct {
	Id            int
	KodeIndikator string
	Target        string
	Satuan        string
	Tahun         string
	Jenis         string
}
