package pohonkinerja

type FindByMultipleRekinRequest struct {
	RekinIds []string `json:"rekin_ids" validate:"required,min=1"`
}
