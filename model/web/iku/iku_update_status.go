package iku

type IkuUpdateActiveRequest struct {
	IndikatorId string `json:"indikator_id"`
	IsActive    bool   `json:"is_active"`
}
