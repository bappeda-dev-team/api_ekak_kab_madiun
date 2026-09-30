package jenisinovasi

type JenisInovasiRequest struct {
	Jenis string `json:"jenis" validate:"required"`
}

type JenisInovasiUpdateRequest struct {
	ID    int    `json:"id"`
	Jenis string `json:"jenis" validate:"required"`
}