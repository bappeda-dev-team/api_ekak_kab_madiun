package rencanakinerja

import (
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/opdmaster"
	"ekak_kabupaten_madiun/model/web/subkegiatan"
)

type RencanaKinerjaResponse struct {
	Id                    string                          `json:"id_rencana_kinerja,omitempty"`
	IdPohon               int                             `json:"id_pohon,omitempty"`
	IdParentPohon         int                             `json:"id_parent_pohon,omitempty"`
	PerluUbahPohonKinerja bool                            `json:"perlu_ubah_pokin"`
	NamaPohon             string                          `json:"nama_pohon,omitempty"`
	LevelPohon            int                             `json:"level_pohon,omitempty"`
	NamaRencanaKinerja    string                          `json:"nama_rencana_kinerja,omitempty"`
	TahunAwal             string                          `json:"tahun_awal,omitempty"`
	TahunAkhir            string                          `json:"tahun_akhir,omitempty"`
	JenisPeriode          string                          `json:"jenis_periode,omitempty"`
	Tahun                 string                          `json:"tahun,omitempty"`
	StatusRencanaKinerja  string                          `json:"status_rencana_kinerja,omitempty"`
	Catatan               string                          `json:"catatan,omitempty"`
	KodeOpd               opdmaster.OpdResponseForAll     `json:"operasional_daerah"`
	PegawaiId             string                          `json:"pegawai_id,omitempty"`
	NamaPegawai           string                          `json:"nama_pegawai,omitempty"`
	Indikator             []IndikatorResponse             `json:"indikator,omitempty"`
	SubKegiatan           subkegiatan.SubKegiatanResponse `json:"sub_kegiatan"`
	Action                []web.ActionButton              `json:"action,omitempty"`
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

type RencanaKinerjaLevel1Response struct {
	Id                   string                      `json:"id_rencana_kinerja,omitempty"`
	IdPohon              int                         `json:"id_pohon,omitempty"`
	SasaranOpdId         int                         `json:"sasaran_opd_id"`
	NamaSasaranOpd       string                      `json:"nama_sasaran_opd,omitempty"`
	NamaPohon            string                      `json:"nama_pohon,omitempty"`
	NamaRencanaKinerja   string                      `json:"nama_rencana_kinerja,omitempty"`
	TahunAwal            string                      `json:"tahun_awal"`
	TahunAkhir           string                      `json:"tahun_akhir"`
	JenisPeriode         string                      `json:"jenis_periode"`
	Tahun                string                      `json:"tahun,omitempty"`
	StatusRencanaKinerja string                      `json:"status_rencana_kinerja,omitempty"`
	Catatan              string                      `json:"catatan,omitempty"`
	KodeOpd              opdmaster.OpdResponseForAll `json:"operasional_daerah"`
	PegawaiId            string                      `json:"pegawai_id,omitempty"`
	NamaPegawai          string                      `json:"nama_pegawai,omitempty"`
	Indikator            []IndikatorResponseLevel1   `json:"indikator"`
}

type IndikatorResponseLevel1 struct {
	Id               string           `json:"id_indikator,omitempty"`
	RencanaKinerjaId string           `json:"rencana_kinerja_id,omitempty"`
	NamaIndikator    string           `json:"nama_indikator,omitempty"`
	Formula          string           `json:"rumus_perhitungan"`
	SumberData       string           `json:"sumber_data"`
	Target           []TargetResponse `json:"targets,omitempty"`
	ManualIKExist    bool             `json:"manual_ik_exist"`
}
