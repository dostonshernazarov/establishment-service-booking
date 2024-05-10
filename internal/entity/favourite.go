package entity

import "time"

type Favourite struct {
	FavouriteId     string
	EstablishmentId string
	UserId          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
}
