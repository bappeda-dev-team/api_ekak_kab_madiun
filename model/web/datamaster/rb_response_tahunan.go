package datamaster

type RbResponseTahunan struct {
	IdRB          int           `json:"id"`
	JenisRB       string        `json:"jenis_rb"`
	KegiatanUtama string        `json:"kegiatan_utama"`
	Keterangan    string        `json:"keterangan"`
	TahunBaseline int           `json:"tahun_baseline"`
	TahunNext     int           `json:"tahun_next"`
}
