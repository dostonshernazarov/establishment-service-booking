package entity

import "time"

type Review struct {
	ReviewId        string
	EstablishmentId string
	UserId          string
	Rating          float64
	Comment         string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
}
