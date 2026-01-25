package bidangurusanresponse

type BidangUrusanCreateRequest struct {
	KodeBidangUrusan string `json:"kode_bidang_urusan"`
	NamaBidangUrusan string `json:"nama_bidang_urusan"`
	Tahun            string `json:"tahun"`
}

type BidangUrusanOPDCreateRequest struct {
	KodeBidangUrusan string `json:"kode_bidang_urusan"`
	KodeOpd          string `json:"kode_opd"`
}
