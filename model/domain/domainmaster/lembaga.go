package domainmaster

type Lembaga struct {
	Id                 string
	KodeLembaga        string
	NamaLembaga        string
	NamaKepalaPemda    string
	NipKepalaPemda     string
	JabatanKepalaPemda string
	IsActive           bool //default true
}
