package pohonkinerja

type TaggingResponse struct {
	Id                int     `json:"id,omitempty"`
	IdPokin           int     `json:"id_pokin,omitempty"`
	NamaTagging       string  `json:"nama_tagging"`
	KeteranganTagging *string `json:"keterangan_tagging"`
}
