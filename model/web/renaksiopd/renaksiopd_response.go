package renaksiopd

type RencanaAksiOpdResponse struct {
	SasaranOpdId   int                      `json:"sasaran_opd_id"`
	NamaSasaranOpd string                   `json:"nama_sasaran_opd"`
	TahunRenaksi   string                   `json:"tahun_renaksi"`
	RencanaKinerja []RencanaKinerjaResponse `json:"rencana_kinerja"`
}

type RencanaKinerjaResponse struct {
	Id                 int                   `json:"id_renaksiopd"`
	RekinId            string                `json:"rekin_id"`
	NamaRencanaKinerja string                `json:"nama_rencana_kinerja"`
	NipPegawai         string                `json:"nip_pegawai"`
	NamaPegawai        string                `json:"nama_pegawai"`
	KodeOpd            string                `json:"kode_opd"`
	TotalAnggaran      int64                 `json:"total_anggaran"`
	Tw1                int                   `json:"tw1"`
	Tw2                int                   `json:"tw2"`
	Tw3                int                   `json:"tw3"`
	Tw4                int                   `json:"tw4"`
	Keterangan         *string               `json:"keterangan"`
	SubKegiatan        []SubKegiatanResponse `json:"subkegiatan"`
}

type SubKegiatanResponse struct {
	KodeSubKegiatan string              `json:"kode_subkegiatan"`
	NamaSubKegiatan string              `json:"nama_subkegiatan"`
	Indikator       []IndikatorResponse `json:"indikator"`
}

type IndikatorResponse struct {
	Id        string `json:"id"`
	Indikator string `json:"indikator"`
	Target    string `json:"target"`
	Satuan    string `json:"satuan"`
}

type RencanaAksiOpdRequestResponse struct {
	SasaranOpdId int     `json:"sasaran_opd_id"`
	RekinId      string  `json:"rekin_id"`
	TahunRenaksi string  `json:"tahun_renaksi"`
	Tw1          int     `json:"tw1"`
	Tw2          int     `json:"tw2"`
	Tw3          int     `json:"tw3"`
	Tw4          int     `json:"tw4"`
	Keterangan   *string `json:"keterangan"`
}

type RencanaAksiOpdByIdResponse struct {
	Id                 int                      `json:"id_renaksiopd"`
	RekinId            string                   `json:"rekin_id"`
	TahunRenaksi       string                   `json:"tahun_renaksi"`
	Keterangan         *string                  `json:"keterangan"`
	NamaRencanaKinerja string                   `json:"nama_rencana_kinerja"`
	SasaranOpd         SasaranOpdDetailResponse `json:"sasaran_opd"`
}

type SasaranOpdDetailResponse struct {
	Id             int                           `json:"id"`
	NamaSasaranOpd string                        `json:"nama_sasaran_opd"`
	TahunAwal      string                        `json:"tahun_awal"`
	TahunAkhir     string                        `json:"tahun_akhir"`
	JenisPeriode   string                        `json:"jenis_periode"`
	Indikator      []IndikatorSasaranOpdResponse `json:"indikator"`
}

type IndikatorSasaranOpdResponse struct {
	Id               string         `json:"id"`
	Indikator        string         `json:"indikator"`
	RumusPerhitungan string         `json:"rumus_perhitungan"`
	SumberData       string         `json:"sumber_data"`
	Target           TargetResponse `json:"target"`
}

type TargetResponse struct {
	Id          string `json:"id"`
	IndikatorId string `json:"indikator_id"`
	Tahun       string `json:"tahun"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}
