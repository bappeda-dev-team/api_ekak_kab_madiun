package domain

import "time"

type TaggingPokin struct {
	Id                int
	IdPokin           int
	NamaTagging       string
	KeteranganTagging *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
