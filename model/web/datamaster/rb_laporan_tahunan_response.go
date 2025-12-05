package datamaster

type RbLaporanTahunanResponse struct {
	IdRB          int             `json:"id"`
	JenisRB       string          `json:"jenis_rb"`
	KegiatanUtama string          `json:"kegiatan_utama"`
	Keterangan    string          `json:"keterangan"`
	TahunBaseline int             `json:"tahun_baseline"`
	TahunNext     int             `json:"tahun_next"`
	Indikator     []IndikatorRB   `json:"indikator"`
	RencanaAksis  []RencanaAksiRB `json:"rencana_aksis"`
}

type RencanaAksiRB struct {
	RencanaAksi     string                   `json:"rencana_aksi"` // rencana kinerja
	IndikatorOutput []IndikatorRencanaAksiRB `json:"indikator_rencana_aksis"`
	Anggaran        int                      `json:"anggaran,string"`
	Realisasi       int                      `json:"realisasi_anggaran,string"`
	OpdKoordinator  string                   `json:"opd_koordinator"`
	NipPelaksana    string                   `json:"nip_pelaksana"`
	NamaPelaksana   string                   `json:"nama_pelaksana"`
	OpdCrosscutting []OpdCrosscutting        `json:"opd_crosscuttings"`
}

type IndikatorRencanaAksiRB struct {
	Indikator       string                         `json:"indikator"`
	TargetIndikator []TargetIndikatorRencanaAksiRB `json:"targets"`
}

type TargetIndikatorRencanaAksiRB struct {
	Target    int    `json:"target,string"`
	Realisasi int    `json:"realisasi,string"`
	Satuan    string `json:"satuan"`
	Capaian   int    `json:"capaian,string"`
}

type OpdCrosscutting struct {
	KodeOpd       string `json:"kode_opd"`
	NamaOpd       string `json:"nama_opd"`
	NipPelaksana  string `json:"nip_pelaksana"`
	NamaPelaksana string `json:"nama_pelaksana"`
}
