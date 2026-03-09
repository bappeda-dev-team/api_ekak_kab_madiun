package user

type UserResponse struct {
	Id          int            `json:"id,omitempty"`
	Nip         string         `json:"nip"`
	Email       string         `json:"email,omitempty"`
	NamaPegawai string         `json:"nama_pegawai"`
	IdJabatan   string         `json:"id_jabatan"`
	NamaJabatan string         `json:"nama_jabatan"`
	IsActive    bool           `json:"is_active"`
	PegawaiId   string         `json:"pegawai_id,omitempty"`
	Role        []RoleResponse `json:"role"`
}

type CekAdminOpdResponse struct {
	KodeOpd    string               `json:"kode_opd"`
	NamaOpd    string               `json:"nama_opd"`
	AdminUsers []AdminOpdUserDetail `json:"admin_users"` // bisa kosong array
}

type AdminOpdUserDetail struct {
	UserId      int    `json:"user_id"`
	Nip         string `json:"nip"`
	NamaPegawai string `json:"nama_pegawai"`
	Email       string `json:"email"`
	IsActive    bool   `json:"is_active"`
}
