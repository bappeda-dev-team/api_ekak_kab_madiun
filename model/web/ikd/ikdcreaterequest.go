package ikd

type ProgramOpdTerpilihCreateRequest struct {
	PohonKinerjaId int `json:"pohon_kinerja_id" validate:"required"`
	ProgramOpdId   int `json:"program_opd_id" validate:"required"`
}