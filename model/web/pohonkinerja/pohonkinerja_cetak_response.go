package pohonkinerja

import "time"

type CetakResponse[T any] struct {
	Nama    string    `json:"nama"`
	Version string    `json:"version"`
	Time    time.Time `json:"time"`
	Item    T         `json:"item"`
}

type PokinCetak struct {
	Id         int          `json:"id"`
	LevelPohon int          `json:"level_pohon"`
	JenisPohon string       `json:"jenis_pohon"`
	NamaPohon  string       `json:"nama_pohon"`
	Childs     []PokinCetak `json:"childs"`
}
