package pohonkinerja

type PohonKinerjaCloneHierarchyRequest struct {
	IdPokinSource int    `json:"id_pokin_source" validate:"required"`
	TahunSource   string `json:"tahun_source" validate:"required"`
	TahunTarget   string `json:"tahun_target" validate:"required"`
	ParentId      int    `json:"parent_id,omitempty"` // Optional, default 0
}
