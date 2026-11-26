package programkegiatan

// TargetProgramKegiatanResponse — response untuk operasi CRUD tb_target_program_kegiatan.
type TargetProgramKegiatanResponse struct {
	Id            int    `json:"id"`
	KodeIndikator string `json:"kode_indikator"`
	KodeOpd       string `json:"kode_opd"`
	Jenis         string `json:"jenis"`
	Target        string `json:"target"`
	Satuan        string `json:"satuan"`
	Tahun         int    `json:"tahun"`
}

// MasterIndikatorWithTargetResponse — response detail satu master indikator beserta semua target OPD.
type MasterIndikatorWithTargetResponse struct {
	KodeIndikator string                          `json:"kode_indikator"`
	Kode          string                          `json:"kode"`
	Indikator     string                          `json:"indikator"`
	Target        []TargetProgramKegiatanResponse `json:"target"`
}
