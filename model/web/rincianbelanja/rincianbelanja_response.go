package rincianbelanja

import "time"

type RincianBelanjaAsnResponse struct {
	KodeOpd              string                   `json:"kode_opd,omitempty"`
	PegawaiId            string                   `json:"pegawai_id,omitempty"`
	NamaPegawai          string                   `json:"nama_pegawai,omitempty"`
	KodeSubkegiatan      string                   `json:"kode_subkegiatan"`
	NamaSubkegiatan      string                   `json:"nama_subkegiatan"`
	IndikatorSubkegiatan []IndikatorResponse      `json:"indikator_subkegiatan"`
	TotalAnggaran        int                      `json:"total_anggaran"`
	RincianBelanja       []RincianBelanjaResponse `json:"rincian_belanja"`
	Pptk                 []PptkResponse           `json:"pptk"`
}

type PptkResponse struct {
	Id              int        `json:"id"`
	Nip             string     `json:"nip"`
	NamaPegawai     string     `json:"nama_pegawai"`
	KodeOpd         string     `json:"kode_opd"`
	Tahun           int        `json:"tahun"`
	KodeSubKegiatan string     `json:"kode_sub_kegiatan"`
	NipAtasan       *string    `json:"nip_atasan"`
	NamaAtasan      *string    `json:"nama_atasan"`
	AktifAt         time.Time  `json:"aktif_at"`
	NonAktifAt      *time.Time `json:"nonaktif_at"`
}

type RincianBelanjaResponse struct {
	RencanaKinerjaId string                `json:"rencana_kinerja_id"`
	RencanaKinerja   string                `json:"rencana_kinerja"`
	PegawaiId        string                `json:"pegawai_id,omitempty"`
	NamaPegawai      string                `json:"nama_pegawai,omitempty"`
	Indikator        []IndikatorResponse   `json:"indikator"`
	TotalAnggaran    int                   `json:"total_anggaran"`
	RencanaAksi      []RencanaAksiResponse `json:"rencana_aksi"`
}

type RencanaAksiResponse struct {
	RenaksiId string `json:"renaksi_id"`
	Renaksi   string `json:"renaksi"`
	Anggaran  int    `json:"anggaran"`
}

type IndikatorResponse struct {
	Id               string           `json:"id_indikator,omitempty"`
	KodeSubkegiatan  string           `json:"kode_subkegiatan,omitempty"`
	KodeOPD          string           `json:"kode_opd,omitempty"`
	RencanaKinerjaId string           `json:"rencana_kinerja_id,omitempty"`
	NamaIndikator    string           `json:"nama_indikator,omitempty"`
	Target           []TargetResponse `json:"targets,omitempty"`
}

type TargetResponse struct {
	Id          string `json:"id_target,omitempty"`
	IndikatorId string `json:"indikator_id"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}
