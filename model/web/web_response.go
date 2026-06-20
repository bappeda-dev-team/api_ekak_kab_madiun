package web

type WebResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   any    `json:"data"`
	Error  string `json:"error,omitempty"`
}

type WebRencanaKinerjaResponse struct {
	Code   int            `json:"code"`
	Status string         `json:"status"`
	Action []ActionButton `json:"rencana_kinerja_action,omitempty"`
	Data   interface{}    `json:"rencana_kinerja"`
}

type WebRencanaAksiResponse struct {
	Code   int            `json:"code"`
	Status string         `json:"status"`
	Action []ActionButton `json:"rencana_aksi_action,omitempty"`
	Data   interface{}    `json:"renaksi"`
}

type WebPelaksanaanRencanaAksiResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"pelaksanaan_renaksi"`
}

type WebUsulanMusrebangResponse struct {
	Code        int            `json:"code"`
	Status      string         `json:"status"`
	Data        interface{}    `json:"usulan_musrebang,omitempty"`
	Action      []ActionButton `json:"pilihan_action,omitempty"`
	DataPilihan interface{}    `json:"usulan_terpilih_musrebang,omitempty"`
}

type WebUsulanMandatoriResponse struct {
	Code        int            `json:"code"`
	Status      string         `json:"status"`
	Data        interface{}    `json:"usulan_mandatori,omitempty"`
	Action      []ActionButton `json:"pilihan_action,omitempty"`
	DataPilihan interface{}    `json:"usulan_terpilih_mandatori,omitempty"`
}

type WebUsulanPokokPikiranResponse struct {
	Code        int            `json:"code"`
	Status      string         `json:"status"`
	Data        interface{}    `json:"usulan_pokok_pikiran,omitempty"`
	Action      []ActionButton `json:"pilihan_action,omitempty"`
	DataPilihan interface{}    `json:"usulan_terpilih_pokir,omitempty"`
}

type WebUsulanInisiatifResponse struct {
	Code        int            `json:"code"`
	Status      string         `json:"status"`
	Data        interface{}    `json:"usulan_inisiatif,omitempty"`
	Action      []ActionButton `json:"pilihan_action,omitempty"`
	DataPilihan interface{}    `json:"usulan_terpilih_inisiatif,omitempty"`
}

type WebUsulanTerpilihResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"usulan_terpilih"`
}

type WebGambaranUmumResponse struct {
	Code   int            `json:"code"`
	Status string         `json:"status"`
	Action []ActionButton `json:"action,omitempty"`
	Data   interface{}    `json:"gambaran_umum"`
}

type WebDasarHukumResponse struct {
	Code   int            `json:"code"`
	Status string         `json:"status"`
	Action []ActionButton `json:"action,omitempty"`
	Data   interface{}    `json:"dasar_hukum"`
}

type WebInovasiResponse struct {
	Code   int            `json:"code"`
	Status string         `json:"status"`
	Action []ActionButton `json:"action,omitempty"`
	Data   interface{}    `json:"inovasi"`
}

type WebSubKegiatanResponse struct {
	Code   int            `json:"code"`
	Status string         `json:"status"`
	Action []ActionButton `json:"pilihan_subkegiatan_action,omitempty"`
	Data   interface{}    `json:"sub_kegiatan"`
}

type WebSubKegiatanTerpilihResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"sub_kegiatan_terpilih"`
}
