package ikk

type IkkRequest struct {
	KodeBidangUrusan string             `json:"kode_bidang_urusan" validate:"required"`
	KodeOpd          string             `json:"kode_opd" validate:"required"`
	Jenis            string             `json:"jenis" validate:"required"`
	Indikators       []IndikatorRequest `json:"indikators" validate:"required,dive"`
	Tahun            int                `json:"tahun" validate:"required"`
	Keterangan       string             `json:"keterangan"`
}

type IndikatorRequest struct {
	ID        int             `json:"id"`
	Indikator string          `json:"indikator" validate:"required"`
	Targets   []TargetRequest `json:"targets" validate:"dive"`
}

type TargetRequest struct {
	ID     int    `json:"id"`
	Target string `json:"target"`
	Satuan string `json:"satuan"`
}

type IkkTerpilihCreateRequest struct {
	PohonKinerjaId int `json:"pohon_kinerja_id"`
	IkkId          int `json:"ikk_id"`
}