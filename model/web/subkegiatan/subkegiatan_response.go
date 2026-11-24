package subkegiatan

import (
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
)

type SubKegiatanWithIndikatorV2Response struct {
	Id              string                                                   `json:"id"`
	KodeSubKegiatan string                                                   `json:"kode_subkegiatan"`
	NamaSubKegiatan string                                                   `json:"nama_sub_kegiatan"`
	Indikator       []programkegiatan.MasterIndikatorProgramKegiatanResponse `json:"indikator"`
}

type SubKegiatanResponse struct {
	SubKegiatanTerpilihId string                         `json:"subkegiatanterpilih_id,omitempty"`
	Id                    string                         `json:"id,omitempty"`
	RekinId               string                         `json:"rekin_id,omitempty"`
	Status                string                         `json:"status,omitempty"`
	KodeSubKegiatan       string                         `json:"kode_subkegiatan,omitempty"`
	NamaSubKegiatan       string                         `json:"nama_sub_kegiatan,omitempty"`
	KodeOpd               string                         `json:"kode_opd,omitempty"`
	NamaOpd               string                         `json:"nama_opd,omitempty"`
	Tahun                 string                         `json:"tahun,omitempty"`
	Indikator             []IndikatorResponse            `json:"indikator,omitempty"`
	IndikatorSubkegiatan  []IndikatorSubKegiatanResponse `json:"indikator_subkegiatan,omitempty"`
	PaguSubKegiatan       []PaguSubKegiatanResponse      `json:"pagu,omitempty"`
	Action                []web.ActionButton             `json:"action,omitempty"`
}

type IndikatorResponse struct {
	IdIndikator   int              `json:"id,omitempty"`
	Id            string           `json:"id_indikator"`
	KodeIndikator string           `json:"kode_indikator,omitempty"`
	Kode          string           `json:"kode,omitempty"`
	NamaIndikator string           `json:"nama_indikator"`
	Indikator     string           `json:"indikator,omitempty"`
	Target        []TargetResponse `json:"targets"`
}

type TargetResponse struct {
	Id              string `json:"id_target"`
	IndikatorId     string `json:"indikator_id"`
	TargetIndikator string `json:"target"`
	SatuanIndikator string `json:"satuan"`
}

type IndikatorSubKegiatanResponse struct {
	Id            string `json:"id"`
	SubKegiatanId string `json:"sub_kegiatan_id"`
	NamaIndikator string `json:"indikator"`
}

type PaguSubKegiatanResponse struct {
	Id            string `json:"id"`
	SubKegiatanId string `json:"sub_kegiatan_id"`
	JenisPagu     string `json:"jenis"`
	PaguAnggaran  int    `json:"pagu_anggaran"`
	Tahun         string `json:"tahun"`
}
