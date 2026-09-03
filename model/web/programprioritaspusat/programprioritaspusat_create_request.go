package programprioritaspusat

type ProgramPrioritasPusatCreateRequest struct {
	KodeProgramPrioritasPusat       string `json:"kode_program_prioritas_pusat"`
	NamaTagging                     string `json:"nama_program_prioritas_pusat"`
	KeteranganProgramPrioritasPusat string `json:"rencana_implementasi"`
	Keterangan                      string `json:"keterangan"`
	TahunAwal                       string `json:"tahun_awal" validate:"required"`
	TahunAkhir                      string `json:"tahun_akhir" validate:"required"`
}

type FindByIdTerkaitRequest struct {
	Ids []int `json:"id_programprioritaspusat" validate:"required,min=1"`
}