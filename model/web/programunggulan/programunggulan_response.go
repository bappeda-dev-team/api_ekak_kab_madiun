package programunggulan

type ProgramUnggulanResponse struct {
	Id                        int                          `json:"id"`
	KodeProgramUnggulan       string                       `json:"kode_program_unggulan"`
	NamaTagging               string                       `json:"nama_program_unggulan"`
	KeteranganProgramUnggulan *string                      `json:"rencana_implementasi"`
	Keterangan                *string                      `json:"keterangan"`
	TahunAwal                 string                       `json:"tahun_awal"`
	TahunAkhir                string                       `json:"tahun_akhir"`
	IsActive                  bool                         `json:"is_active"`
	OpdList                   []OpdProgramUnggulanResponse `json:"opd_list"`
	TahunTerpakai             []string                     `json:"tahun_terpakai"`
}

type OpdProgramUnggulanResponse struct {
	Id      int    `json:"id"`
	KodeOpd string `json:"kode_opd"`
	NamaOpd string `json:"nama_opd"`
}

type CreateOpdProgramUnggulanResponse struct {
	KodeProgramUnggulan string                       `json:"kode_program_unggulan"`
	OpdList             []OpdProgramUnggulanResponse `json:"opd_list"`
}
