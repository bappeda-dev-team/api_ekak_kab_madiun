package domain

type Users struct {
	Id          int
	Nip         string
	Email       string
	Password    string
	IsActive    bool
	Role        []Roles
	PegawaiId   string
	NamaPegawai string
}
