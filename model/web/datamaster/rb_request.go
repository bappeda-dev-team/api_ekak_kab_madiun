package datamaster

type RBRequest struct {
    JenisRB       string             `json:"jenis_rb" validate:"required"`
    KegiatanUtama string             `json:"kegiatan_utama" validate:"required"`
    Keterangan    string             `json:"keterangan"`
    TahunBaseline int                `json:"tahun_baseline" validate:"required"`
    TahunNext     int                `json:"tahun_next" validate:"required"`
    Indikator     []IndikatorRequest `json:"indikator" validate:"required,dive"`
}

type IndikatorRequest struct {
    Indikator string               `json:"indikator" validate:"required"`
    Target    []TargetRBRequest    `json:"target" validate:"required,dive"`
}

type TargetRBRequest struct {
    // Baseline
    TahunBaseline     *int     `json:"tahun_baseline"`      // boleh null
    TargetBaseline    *int     `json:"target_baseline,string"`
    RealisasiBaseline *float32 `json:"realisasi_baseline,string"`
    SatuanBaseline    *string  `json:"satuan_baseline"`

    // Next
    TahunNext  *int    `json:"tahun_next"`
    TargetNext *int    `json:"target_next,string"`
    SatuanNext *string `json:"satuan_next"`
}
