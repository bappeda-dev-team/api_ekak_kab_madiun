package rencanakinerja

type FindByIdRekinsRequest struct {
	Ids []string `json:"id_rekins" validate:"required,min=1"`
	Bulan int `json:"bulan"`
	Tahun int `json:"tahun"`
}
