package taggingpokin

type TaggingPokinUpdateRequest struct {
	Id                int     `json:"id"`
	NamaTagging       string  `json:"nama_tagging"`
	KeteranganTagging *string `json:"keterangan_tagging"`
}
