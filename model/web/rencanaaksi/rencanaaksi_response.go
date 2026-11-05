package rencanaaksi

import "ekak_kabupaten_madiun/model/web"

type RencanaAksiResponse struct {
	Id                     string                           `json:"id"`
	RencanaKinerjaId       string                           `json:"rekin_id"`
	KodeOpd                string                           `json:"kode_opd,omitempty"`
	Urutan                 int                              `json:"urutan"`
	NamaRencanaAksi        string                           `json:"nama_rencana_aksi"`
	PelaksanaanRencanaAksi []PelaksanaanRencanaAksiResponse `json:"pelaksanaan"`
	JumlahBobot            int                              `json:"jumlah_bobot,omitempty"`
	TotalBobotRencanaAksi  int                              `json:"total_bobot_rencana_aksi,omitempty"`
	Action                 []web.ActionButton               `json:"action,omitempty" swaggerignore:"true"`
}

type BobotBulanan struct {
	Bulan      int `json:"bulan"`
	TotalBobot int `json:"total_bobot"`
}

type RencanaAksiTableResponse struct {
	RencanaAksi      []RencanaAksiResponse `json:"rencana_aksi"`
	TotalPerBulan    []BobotBulanan        `json:"total_per_bulan"`
	TotalKeseluruhan int                   `json:"total_keseluruhan"`
	WaktuDibutuhkan  int                   `json:"waktu_dibutuhkan"`
}

type RencanaAksiPegawaiResponse struct {
	RencanaKinerjaId   string                `json:"rekin_id"`
	NamaRencanaKinerja string                `json:"nama_rencana_kinerja"`
	Nip                string                `json:"nip"`
	KodeOpd            string                `json:"kode_opd"`
	Tahun              string                `json:"tahun"`
	JumlahBobot        int                   `json:"jumlah_bobot"`
	RencanaAksi        []RencanaAksiResponse `json:"rencana_aksi"`
}
