package rencanakinerja

import (
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/opdmaster"
)

type RencanaKinerjaResponse struct {
	Id                   string                      `json:"id_rencana_kinerja,omitempty"`
	IdPohon              int                         `json:"id_pohon,omitempty"`
	NamaPohon            string                      `json:"nama_pohon,omitempty"`
	NamaRencanaKinerja   string                      `json:"nama_rencana_kinerja,omitempty"`
	TahunAwal            string                      `json:"tahun_awal,omitempty"`
	TahunAkhir           string                      `json:"tahun_akhir,omitempty"`
	JenisPeriode         string                      `json:"jenis_periode,omitempty"`
	Tahun                string                      `json:"tahun,omitempty"`
	StatusRencanaKinerja string                      `json:"status_rencana_kinerja,omitempty"`
	Catatan              string                      `json:"catatan,omitempty"`
	KodeOpd              opdmaster.OpdResponseForAll `json:"operasional_daerah,omitempty"`
	PegawaiId            string                      `json:"pegawai_id,omitempty"`
	NamaPegawai          string                      `json:"nama_pegawai,omitempty"`
	Indikator            []IndikatorResponse         `json:"indikator,omitempty"`
	// SubKegiatan          subkegiatan.SubKegiatanResponse `json:"sub_kegiatan,omitempty"`
	Action []web.ActionButton `json:"action,omitempty"`
}

type IndikatorResponse struct {
	Id               string           `json:"id_indikator,omitempty"`
	RencanaKinerjaId string           `json:"rencana_kinerja_id,omitempty"`
	NamaIndikator    string           `json:"nama_indikator,omitempty"`
	Target           []TargetResponse `json:"targets,omitempty"`
	ManualIK         *DataOutput      `json:"data_output,omitempty"`
	ManualIKExist    bool             `json:"manual_ik_exist"`
}

type TargetResponse struct {
	Id              string `json:"id_target,omitempty"`
	IndikatorId     string `json:"indikator_id"`
	TargetIndikator string `json:"target"`
	SatuanIndikator string `json:"satuan"`
	Tahun           string `json:"tahun,omitempty"`
}
