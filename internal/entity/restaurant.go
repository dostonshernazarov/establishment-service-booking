package entity

import "time"

type Restaurant struct {
	RestaurantId   string
	OwnerId        string
	RestaurantName string
	Description    string
	Rating         float32
	OpeningHours   string
	ContactNumber  string
	LicenceUrl     string
	WebsiteUrl     string
	Images         []*Image
	Location       Location
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}
