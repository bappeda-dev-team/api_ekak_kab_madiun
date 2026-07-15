package pohonkinerja

import "time"

type CetakResponse[T any] struct {
	Nama    string    `json:"nama"`
	Version string    `json:"version"`
	Time    time.Time `json:"time"`
	Item    T         `json:"item"`
}

type PokinCetak struct {
	Id         int    `json:"id"`
	ParentId   int    `json:"parent_id"`
	LevelPohon int    `json:"level_pohon"`
	JenisPohon string `json:"jenis_pohon"`
	NamaPohon  string `json:"nama_pohon"`
}

type PokinOpdCetak struct {
	Tahun      int                 `json:"tahun"`
	KodeOpd    string              `json:"kode_opd"`
	NamaOpd    string              `json:"nama_opd"`
	TujuanOpds []TujuanOpdResponse `json:"tujuan_opds"`
	Pokins     []PokinCetak        `json:"pohon_kinerjas"`
}
