package pohonkinerja

type FindByMultipleRekinRequest struct {
	RekinIds []string `json:"rekin_ids" validate:"required,min=1"`
}

type MultiRekinDetailsByOpdAndTahunRequest struct {
	KodeOpd string
	Tahun   string
}
