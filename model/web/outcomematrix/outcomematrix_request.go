package outcomematrix

type OutcomeMatrixCreateRequest struct {
	KodeOpd string `json:"kode_opd" validate:"required"`
	Kode    string `json:"kode" validate:"required"`
	Jenis   string `json:"jenis"`
	Outcome string `json:"outcome" validate:"required"`
}

type OutcomeMatrixUpdateRequest struct {
	Id      int    `json:"id"`
	KodeOpd string `json:"kode_opd" validate:"required"`
	Kode    string `json:"kode" validate:"required"`
	Jenis   string `json:"jenis"`
	Outcome string `json:"outcome" validate:"required"`
}

type OutcomeMatrixBatchItemRequest struct {
	Id      int    `json:"id"`
	KodeOpd string `json:"kode_opd" validate:"required"`
	Kode    string `json:"kode" validate:"required"`
	Jenis   string `json:"jenis"`
	Outcome string `json:"outcome"`
}
