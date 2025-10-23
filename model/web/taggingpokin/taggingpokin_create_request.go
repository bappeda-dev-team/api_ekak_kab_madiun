package taggingpokin

type TaggingPokinCreateRequest struct {
	NamaTagging       string  `json:"nama_tagging"`
	KeteranganTagging *string `json:"keterangan_tagging"`
}
