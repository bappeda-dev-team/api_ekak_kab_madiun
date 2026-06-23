package programunggulan

type ProgramUnggulanCreateRequest struct {
	KodeProgramUnggulan       string `json:"kode_program_unggulan"`
	NamaTagging               string `json:"nama_program_unggulan"`
	KeteranganProgramUnggulan string `json:"rencana_implementasi"`
	Keterangan                string `json:"keterangan"`
	TahunAwal                 string `json:"tahun_awal" validate:"required"`
	TahunAkhir                string `json:"tahun_akhir" validate:"required"`
}

type FindByIdTerkaitRequest struct {
	Ids []int `json:"id_prorgramunggulan" validate:"required,min=1"`
}

type CreateOpdProgramUnggulanRequest struct {
	KodeProgramUnggulan string   `json:"kode_program_unggulan" validate:"required"`
	KodeOpd             []string `json:"kode_opd" validate:"required,min=1,dive,required"`
}
type DeleteOpdProgramUnggulanRequest struct {
	KodeProgramUnggulan string `json:"kode_program_unggulan" validate:"required"`
	KodeOpd             string `json:"kode_opd" validate:"required"`
}
