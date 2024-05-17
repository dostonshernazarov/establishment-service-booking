package entity

import "time"

type Attraction struct {
	AttractionId   string
	OwnerId        string
	AttractionName string
	Description    string
	Rating         float32
	ContactNumber  string
	LicenceUrl     string
	WebsiteUrl     string
	Images         []*Image
	Location       Location
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}

type Image struct {
	ImageId         string
	EstablishmentId string
	ImageUrl        string
	Category        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
}

type Location struct {
	LocationId      string
	EstablishmentId string
	Address         string
	Latitude        float32
	Longitude       float32
	Country         string
	City            string
	StateProvince   string
	Category        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
}
