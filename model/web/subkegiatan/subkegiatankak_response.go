package subkegiatan

type SubKegiatanKAKResponse struct {
	KodeOpd      string                       `json:"kode_opd"`
	NamaOpd      string                       `json:"nama_opd"`
	Program      ProgramKAKResponse           `json:"program"`
	Kegiatan     KegiatanKAKResponse          `json:"kegiatan"`
	SubKegiatan  SubKegiatanDetailKAKResponse `json:"sub_kegiatan"`
	PaguAnggaran string                       `json:"pagu_anggaran"`
}

type ProgramKAKResponse struct {
	Kode                    string                      `json:"kode"`
	Nama                    string                      `json:"nama"`
	IndikatorKinerjaProgram IndikatorKinerjaKAKResponse `json:"indikator_kinerja_program"`
}

type KegiatanKAKResponse struct {
	Kode                     string                      `json:"kode"`
	Nama                     string                      `json:"nama"`
	IndikatorKinerjaKegiatan IndikatorKinerjaKAKResponse `json:"indikator_kinerja_kegiatan"`
}

type SubKegiatanDetailKAKResponse struct {
	Subkegiatan                 string                      `json:"subkegiatan"`
	Nama                        string                      `json:"nama"`
	IndikatorKinerjaSubKegiatan IndikatorKinerjaKAKResponse `json:"indikator_kinerja_sub_kegiatan"`
}

type IndikatorKinerjaKAKResponse struct {
	Nama   string `json:"nama"`
	Target string `json:"target"`
	Satuan string `json:"satuan"`
}
