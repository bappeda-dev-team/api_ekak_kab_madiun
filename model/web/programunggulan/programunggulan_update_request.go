package programunggulan

type ProgramUnggulanUpdateRequest struct {
	Id                        int    `json:"id"`
	NamaTagging               string `json:"nama_program_unggulan"`
	KeteranganProgramUnggulan string `json:"rencana_implementasi"`
	Keterangan                string `json:"keterangan"`
	TahunAwal                 string `json:"tahun_awal" validate:"required"`
	TahunAkhir                string `json:"tahun_akhir" validate:"required"`
}
