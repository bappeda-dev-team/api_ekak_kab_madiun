package iku

import "time"

type IkuResponse struct {
	IndikatorId         string           `json:"indikator_id"`
	Sumber              string           `json:"asal_iku"`
	IkuActive           bool             `json:"iku_active"`
	IsActive            bool             `json:"is_active,omitempty"`
	Indikator           string           `json:"indikator"`
	DefinisiOperasional string           `json:"definisi_operasional"`
	RumusPerhitungan    string           `json:"rumus_perhitungan"`
	SumberData          string           `json:"sumber_data"`
	CreatedAt           time.Time        `json:"created_at"`
	TahunAwal           string           `json:"tahun_awal"`
	TahunAkhir          string           `json:"tahun_akhir"`
	JenisPeriode        string           `json:"jenis_periode"`
	Target              []TargetResponse `json:"target"`
}

type TargetResponse struct {
	Target string `json:"target"`
	Satuan string `json:"satuan"`
	Tahun  string `json:"tahun"`
}

type IkuOpdResponse struct {
	IndikatorId         string              `json:"indikator_id"`
	AsalIku             string              `json:"asal_iku"`
	Indikator           string              `json:"indikator"`
	IkuActive           bool                `json:"iku_active"`
	IsHide              bool                `json:"is_hide"`
	CreatedAt           time.Time           `json:"created_at"`
	DefinisiOperasional string              `json:"definisi_operasional"`
	RumusPerhitungan    string              `json:"rumus_perhitungan"`
	SumberData          string              `json:"sumber_data"`
	Jenis               string              `json:"jenis"`
	TahunAwal           string              `json:"tahun_awal"`
	TahunAkhir          string              `json:"tahun_akhir"`
	JenisPeriode        string              `json:"jenis_periode"`
	Target              []TargetOpdResponse `json:"target"`
}

type TargetOpdResponse struct {
	Target string `json:"target"`
	Satuan string `json:"satuan"`
	Tahun  string `json:"tahun"`
}

// ─── IKU Pemda V2 — Dual Target Response ────────────────────────

// Digunakan untuk endpoint /iku_pemda/v2/rankhir/:tahun
// Menampilkan target_ranwal (base renstra) dan target_rankhir (override)
type IkuPemdaRankhirDualResponse struct {
	IndikatorId         string           `json:"indikator_id"`
	Sumber              string           `json:"asal_iku"`
	IkuActive           bool             `json:"iku_active"`
	Indikator           string           `json:"indikator"`
	DefinisiOperasional string           `json:"definisi_operasional"`
	RumusPerhitungan    string           `json:"rumus_perhitungan"`
	SumberData          string           `json:"sumber_data"`
	CreatedAt           time.Time        `json:"created_at"`
	TahunAwal           string           `json:"tahun_awal"`
	TahunAkhir          string           `json:"tahun_akhir"`
	JenisPeriode        string           `json:"jenis_periode"`
	TargetRanwal        []TargetResponse `json:"target_ranwal"`
	TargetRankhir       []TargetResponse `json:"target_rankhir"`
}

// Digunakan untuk endpoint /iku_pemda/v2/penetapan/:tahun
// Menampilkan target_rankhir (base) dan target_penetapan (override)
type IkuPemdaPenetapanDualResponse struct {
	IndikatorId         string           `json:"indikator_id"`
	Sumber              string           `json:"asal_iku"`
	IkuActive           bool             `json:"iku_active"`
	Indikator           string           `json:"indikator"`
	DefinisiOperasional string           `json:"definisi_operasional"`
	RumusPerhitungan    string           `json:"rumus_perhitungan"`
	SumberData          string           `json:"sumber_data"`
	CreatedAt           time.Time        `json:"created_at"`
	TahunAwal           string           `json:"tahun_awal"`
	TahunAkhir          string           `json:"tahun_akhir"`
	JenisPeriode        string           `json:"jenis_periode"`
	TargetRankhir       []TargetResponse `json:"target_rankhir"`
	TargetPenetapan     []TargetResponse `json:"target_penetapan"`
}
