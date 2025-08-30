package domain

import "time"

type TaggingPokin struct {
	Id                int
	IdPokin           int
	NamaTagging       string
	KeteranganTagging *string
	CloneFrom         int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
