package subkegiatan

type SubKegiatanTerpilihUpdateRequest struct {
	Id              string `json:"id"`
	KodeSubKegiatan string `json:"kode_subkegiatan"`
}

type SubKegiatanOpdMultipleCreateRequest struct {
	KodeSubkegiatan []string `json:"kode_subkegiatan" validate:"required,min=1"`
	KodeOpd         string   `json:"kode_opd" validate:"required"`
	Tahun           string   `json:"tahun" validate:"required"`
}

// Tetap pertahankan yang lama untuk backward compatibility
type SubKegiatanOpdCreateRequest struct {
	KodeSubkegiatan string `json:"kode_subkegiatan" validate:"required"`
	KodeOpd         string `json:"kode_opd" validate:"required"`
	Tahun           string `json:"tahun" validate:"required"`
}
