package outcomematrix

type OutcomeMatrixCreateRequest struct {
	KodeSubkegiatan string `json:"kode_subkegiatan" validate:"required"`
	Kode            string `json:"kode" validate:"required"`
	Outcome         string `json:"outcome" validate:"required"`
}

type OutcomeMatrixUpdateRequest struct {
	Id              int    `json:"id"`
	KodeSubkegiatan string `json:"kode_subkegiatan" validate:"required"`
	Kode            string `json:"kode" validate:"required"`
	Outcome         string `json:"outcome" validate:"required"`
}

type OutcomeMatrixBatchItemRequest struct {
	Id              int    `json:"id"`
	KodeSubkegiatan string `json:"kode_subkegiatan" validate:"required"`
	Kode            string `json:"kode" validate:"required"`
	Outcome         string `json:"outcome"`
}
