package programunggulan

type ProgramUnggulanCreateRequest struct {
	KodeProgramUnggulan       string `json:"kode_program_unggulan"`
	NamaTagging               string `json:"nama_tagging"`
	KeteranganProgramUnggulan string `json:"keterangan_program_unggulan"`
	Keterangan                string `json:"keterangan"`
	TahunAwal                 string `json:"tahun_awal" validate:"required"`
	TahunAkhir                string `json:"tahun_akhir" validate:"required"`
}
