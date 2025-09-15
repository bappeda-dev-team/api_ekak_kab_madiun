package pohonkinerja

type TaggingResponse struct {
	Id                       int                         `json:"id,omitempty"`
	IdPokin                  int                         `json:"id_pokin,omitempty"`
	NamaTagging              string                      `json:"nama_tagging"`
	KeteranganTaggingProgram []KeteranganTaggingResponse `json:"keterangan_tagging_program"`
	CloneFrom                int                         `json:"clone_from"`
}

type KeteranganTaggingResponse struct {
	Id                  int    `json:"id"`
	IdTagging           int    `json:"id_tagging"`
	KodeProgramUnggulan string `json:"kode_program_unggulan"`
	Tahun               string `json:"tahun,omitempty"`
}
