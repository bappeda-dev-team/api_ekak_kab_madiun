package rencanaaksi

import "ekak_kabupaten_madiun/model/web"

type PelaksanaanRencanaAksiResponse struct {
	Id            string             `json:"id"`
	RencanaAksiId string             `json:"rencana_aksi_id"`
	Bulan         int                `json:"bulan"`
	Bobot         int                `json:"bobot"`
	BobotAvail    int                `json:"bobot_tersedia,omitempty"`
	Action        []web.ActionButton `json:"action,omitempty" swaggerignore:"true"`
}
