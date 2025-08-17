package web

type JWTClaim struct {
	Issuer    string   `json:"iss"`
	Subject   string   `json:"sub"`
	UserId    int      `json:"user_id"`
	PegawaiId string   `json:"pegawai_id"`
	KodeOpd   string   `json:"kode_opd"`
	NamaOpd   string   `json:"nama_opd"`
	Email     string   `json:"email"`
	Nip       string   `json:"nip"`
	Roles     []string `json:"roles"`
	Iat       int64    `json:"iat"`
	Exp       int64    `json:"exp"`
}
