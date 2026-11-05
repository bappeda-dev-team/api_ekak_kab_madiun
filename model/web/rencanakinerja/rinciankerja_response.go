package rencanakinerja

import (
	"ekak_kabupaten_madiun/model/web/dasarhukum"
	"ekak_kabupaten_madiun/model/web/gambaranumum"
	"ekak_kabupaten_madiun/model/web/inovasi"
	"ekak_kabupaten_madiun/model/web/inovasirekin"
	"ekak_kabupaten_madiun/model/web/permasalahan"
	"ekak_kabupaten_madiun/model/web/rencanaaksi"
	"ekak_kabupaten_madiun/model/web/subkegiatan"
)

type DataRincianKerja struct {
	RencanaKinerja RencanaKinerjaResponse                   `json:"rencana_kinerja"`
	RencanaAksi    rencanaaksi.RencanaAksiTableResponse     `json:"rencana_aksis"`
	Usulan         []UsulanGabunganResponse                 `json:"usulan"`
	SubKegiatan    []subkegiatan.SubKegiatanResponse        `json:"subkegiatan"`
	Permasalahan   []permasalahan.PermasalahanRekinResponse `json:"permasalahan"`
	DasarHukum     []dasarhukum.DasarHukumResponse          `json:"dasar_hukum"`
	GambaranUmum   []gambaranumum.GambaranUmumResponse      `json:"gambaran_umum"`
	Inovasi        []inovasi.InovasiResponse                `json:"inovasi"`
	InovasiRekin   []inovasirekin.InovasiRekinResponse      `json:"inovasi_rekin"`
}

type UsulanGabunganResponse struct {
	Id          string `json:"id"`
	JenisUsulan string `json:"jenis_usulan"` // "musrebang", "mandatori", "pokok_pikiran", "inisiatif"
	Usulan      string `json:"usulan"`
	Uraian      string `json:"uraian"`
	Tahun       string `json:"tahun"`
	RekinId     string `json:"rencana_kinerja_id,omitempty"`
	PegawaiId   string `json:"pegawai_id,omitempty"`
	KodeOpd     string `json:"kode_opd"`
	NamaOpd     string `json:"nama_opd,omitempty"`
	IsActive    bool   `json:"is_active,omitempty"`
	Status      string `json:"status"`
	CreatedAt   string `json:"dibuat_pada,omitempty" time_format:"2006-01-02 15:04:05"`
	// Field khusus per jenis usulan
	Alamat           string `json:"alamat,omitempty"`            // untuk musrebang & pokok pikiran
	Manfaat          string `json:"manfaat,omitempty"`           // untuk inisiatif
	PeraturanTerkait string `json:"peraturan_terkait,omitempty"` // untuk mandatori
}
