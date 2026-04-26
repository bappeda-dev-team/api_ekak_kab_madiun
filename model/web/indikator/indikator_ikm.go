package indikator

import "ekak_kabupaten_madiun/model/web/iku"

type IkmRequest struct {
	KodeBidangUrusan    string               `json:"kode_bidang_urusan"`
	NamaBidangUrusan    string               `json:"nama_bidang_urusan"`
	Indikator           string               `json:"indikator"`
	IsActive            bool                 `json:"is_active"`
	DefinisiOperasional string               `json:"definisi_operasional"`
	RumusPerhitungan    string               `json:"rumus_perhitungan"`
	SumberData          string               `json:"sumber_data"`
	Jenis               string               `json:"jenis"`
	TahunAwal           string               `json:"tahun_awal"`
	TahunAkhir          string               `json:"tahun_akhir"`
	Target              []iku.TargetResponse `json:"target"`
}

type IkmResponse struct {
	Id                  string               `json:"id"`
	KodeBidangUrusan    string               `json:"kode_bidang_urusan"`
	NamaBidangUrusan    string               `json:"nama_bidang_urusan"`
	Indikator           string               `json:"indikator"`
	IsActive            bool                 `json:"is_active"`
	DefinisiOperasional string               `json:"definisi_operasional"`
	RumusPerhitungan    string               `json:"rumus_perhitungan"`
	SumberData          string               `json:"sumber_data"`
	Jenis               string               `json:"jenis"`
	TahunAwal           string               `json:"tahun_awal"`
	TahunAkhir          string               `json:"tahun_akhir"`
	Target              []iku.TargetResponse `json:"target"`
}
