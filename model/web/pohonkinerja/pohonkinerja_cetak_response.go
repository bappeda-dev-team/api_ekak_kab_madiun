package pohonkinerja

import "time"

type CetakResponse[T any] struct {
	Nama    string    `json:"nama"`
	Version string    `json:"version"`
	Time    time.Time `json:"time"`
	Item    T         `json:"item"`
}

type PokinCetak struct {
	Id         int           `json:"id"`
	ParentId   int           `json:"parent_id"`
	LevelPohon int           `json:"level_pohon"`
	JenisPohon string        `json:"jenis_pohon"`
	NamaPohon  string        `json:"nama_pohon"`
	Metadata   PokinMetadata `json:"metadata_pohon"`
}

type PokinOpdCetak struct {
	Tahun      int                 `json:"tahun"`
	KodeOpd    string              `json:"kode_opd"`
	NamaOpd    string              `json:"nama_opd"`
	TujuanOpds []TujuanOpdResponse `json:"tujuan_opds"`
	Pokins     []PokinCetak        `json:"pohon_kinerjas"`
}

type PokinMetadata struct {
	IsCrosscutting     bool                `json:"is_crosscutting"`
	CrosscuttingPokins []CrossCuttingPokin `json:"crosscutting_pokins"`
	// metadata lain...
}

type CrossCuttingPokin struct {
	IsCrossCuttingDiterima bool   `json:"is_crosscutting_diterima"`
	NamaPohonPemberi       string `json:"nama_pohon_pemberi"`
	NamaOpdPemberi         string `json:"nama_opd_pemberi"`
	NamaPohonPenerima      string `json:"nama_pohon_penerima"`
	NamaOpdPenerima        string `json:"nama_opd_penerima"`
	KeteranganCrosscutting string `json:"keterangan_crosscutting"`
	StatusCrosscutting     string `json:"status_crosscutting"`
}
