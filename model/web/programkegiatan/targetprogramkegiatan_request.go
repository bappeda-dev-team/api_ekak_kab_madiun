package programkegiatan

// TargetProgramKegiatanCreateRequest — body untuk POST create per jenis.
// Kolom jenis TIDAK perlu diisi karena sudah ditentukan oleh endpoint yang dipanggil.
type TargetProgramKegiatanCreateRequest struct {
	KodeIndikator string `json:"kode_indikator" validate:"required"`
	KodeOpd       string `json:"kode_opd"       validate:"required"`
	Target        string `json:"target"         validate:"required"`
	Satuan        string `json:"satuan"         validate:"required"`
	Tahun         int    `json:"tahun"          validate:"required,min=2000"`
}

// TargetProgramKegiatanUpdateRequest — body untuk PUT update per jenis.
// Hanya target dan satuan yang dapat diubah.
type TargetProgramKegiatanUpdateRequest struct {
	Target string `json:"target" validate:"required"`
	Satuan string `json:"satuan" validate:"required"`
}
