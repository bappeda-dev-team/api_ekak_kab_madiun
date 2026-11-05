package laporancrosscutting

import (
	"ekak_kabupaten_madiun/model/domain"
)

type LaporanCrosscuttingResponse struct {
	IdPohon                int                            `json:"id_pohon"`
	Parent                 int                            `json:"parent"`
	KeteranganCrosscutting string                         `json:"keterangan_crosscutting"`
	NamaPohon              string                         `json:"nama_pohon"`
	Indikator              string                         `json:"indikator"`
	Satuan                 string                         `json:"satuan"`
	Target                 string                         `json:"target"`
	Penerima               []PenerimaCrosscuttingResponse `json:"penerima"`
}

type PenerimaCrosscuttingResponse struct {
	NamaPelaksana string `json:"nama_pelaksana"`
	Nip           string `json:"nip"`
	SubKegiatan   string `json:"sub_kegiatan"`
	PaguAnggaran  string `json:"pagu_anggaran"`
	NamaOpd       string `json:"nama_opd"`
	NamaPohon     string `json:"nama_pohon"`
	Indikator     string `json:"indikator"`
	Satuan        string `json:"satuan"`
	Target        string `json:"target"`
}

type SubKegiatan struct {
	NamaPelaksana string `json:"nama_pelaksana"`
	Nip           string `json:"nip"`
	SubKegiatan   string `json:"sub_kegiatan"`
	PaguAnggaran  string `json:"pagu_anggaran"`
}

type PohonKinerjaCrosscuttingItemResponse struct {
	PohonKinerja        domain.PohonKinerja           `json:"pohon_kinerja"`
	Crosscutting        []CrosscuttingPokinResponse   `json:"crosscutting"`
	CrosscuttingDikirim []CrosscuttingDikirimResponse `json:"crosscutting_dikirim"`
	Indikator           []IndikatorResponse           `json:"indikator"`
	Pelaksana           []PelaksanaOpdResponse        `json:"pelaksana"`
}

type CrosscuttingPokinResponse struct {
	IdCrosscutting         int    `json:"id_crosscutting"`
	KeteranganCrosscutting string `json:"keterangan_crosscutting"`
	NamaPohonAsal          string `json:"nama_pohon_asal"`
	KodeOpdAsal            string `json:"kode_opd_asal"`
	NamaOpdAsal            string `json:"nama_opd_asal,omitempty"`
	Status                 string `json:"status"`
}

type CrosscuttingDikirimResponse struct {
	IdCrosscutting         int    `json:"id_crosscutting"`
	KeteranganCrosscutting string `json:"keterangan_crosscutting"`
	NamaPohonTujuan        string `json:"nama_pohon_tujuan"`
	KodeOpdTujuan          string `json:"kode_opd_tujuan"`
	NamaOpdTujuan          string `json:"nama_opd_tujuan,omitempty"`
	Status                 string `json:"status"`
}

type IndikatorResponse struct {
	Id            string           `json:"id_indikator"`
	IdPokin       string           `json:"id_pokin,omitempty"`
	IdRekin       string           `json:"id_rekin,omitempty"`
	Kode          string           `json:"kode,omitempty"`
	NamaIndikator string           `json:"nama_indikator"`
	Target        []TargetResponse `json:"targets"`
}

type TargetResponse struct {
	Id              string `json:"id_target"`
	IndikatorId     string `json:"indikator_id"`
	TargetIndikator string `json:"target"`
	SatuanIndikator string `json:"satuan"`
	TahunSasaran    string `json:"tahun,omitempty"`
}


type PelaksanaOpdResponse struct {
	Id             string `json:"id_pelaksana"`
	PohonKinerjaId string `json:"pohon_kinerja_id,omitempty"`
	PegawaiId      string `json:"pegawai_id"`
	Nip            string `json:"nip"`
	NamaPegawai    string `json:"nama_pegawai"`
}

type PohonKinerjaCrosscuttingAllResponse struct {
	KodeOpd      string                                 `json:"kode_opd"`
	NamaOpd      string                                 `json:"nama_opd"`
	Tahun        string                                 `json:"tahun"`
	PohonKinerja []PohonKinerjaCrosscuttingItemResponse `json:"pohon_kinerja"`
}