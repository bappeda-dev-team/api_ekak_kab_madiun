package subkegiatan

type SubKegiatanFindAllFilter struct {
	KodeSubKegiatan string
	NamaSubKegiatan string
	Page            int
	Limit           int
}

type SubKegiatanPaginatedResponse struct {
	Items        []SubKegiatanResponse `json:"items"`
	Page         int                   `json:"page"          example:"1"`
	Limit        int                   `json:"limit"         example:"10"`
	Total        int                   `json:"total"         example:"250"`
	TotalPages   int                   `json:"total_pages"   example:"25"`
	HasNext      bool                  `json:"has_next"      example:"true"`
	HasPrevious  bool                  `json:"has_previous"  example:"false"`
	NextPage     int                   `json:"next_page"     example:"2"`
	PreviousPage int                   `json:"previous_page" example:"0"`
}
