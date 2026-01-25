package bidangurusanresponse

type BidangUrusanUpdateRequest struct {
	Id               string `json:"id"`
	KodeBidangUrusan string `json:"kode_bidang_urusan"`
	NamaBidangUrusan string `json:"nama_bidang_urusan"`
	Tahun            string `json:"tahun"`
}

type BidangUrusanOPDUpdateRequest struct {
	Id               int    `json:"id"`
	KodeBidangUrusan string `json:"kode_bidang_urusan"`
	KodeOpd          string `json:"kode_opd"`
}
