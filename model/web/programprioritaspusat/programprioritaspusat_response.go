package programprioritaspusat

type ProgramPrioritasPusatResponse struct {
	Id                              int     `json:"id"`
	KodeProgramPrioritasPusat       string  `json:"kode_program_prioritas_pusat"`
	NamaTagging                     string  `json:"nama_program_prioritas_pusat"`
	KeteranganProgramPrioritasPusat *string `json:"rencana_implementasi"`
	Keterangan                      *string `json:"keterangan"`
	TahunAwal                       string  `json:"tahun_awal"`
	TahunAkhir                      string  `json:"tahun_akhir"`
	IsActive                        bool    `json:"is_active"`
}