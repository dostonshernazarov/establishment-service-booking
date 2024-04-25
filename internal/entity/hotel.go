package entity

import "time"

type Hotel struct {
	HotelId       string
	OwnerId       string
	HotelName     string
	Description   string
	Rating        float32
	ContactNumber string
	LicenceUrl    string
	WebsiteUrl    string
	Images        []*Image
	Location      Location
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}
