package postgresql

import (
	"context"
	"testing"
	"time"

	"Booking/establishment-service-booking/internal/entity"
	"Booking/establishment-service-booking/internal/pkg/config"
	"Booking/establishment-service-booking/internal/pkg/postgres"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateAttraction(t *testing.T) {
	// Connect to database
	cfg := config.New()

	db, err := postgres.New(cfg)
	if err != nil {
		return
	}

	// Test  Method Create
	repo := NewAttractionRepo(db)

	attraction_id := uuid.New().String()

	attraction := &entity.Attraction{
		AttractionId:   attraction_id,
		OwnerId:        uuid.New().String(),
		AttractionName: "test attraction name",
		Description:    "Test description",
		Rating:         4.9,
		ContactNumber:  "+9989123456789",
		LicenceUrl:     "test licence url",
		WebsiteUrl:     "test website url",
		Images: []*entity.Image{
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: attraction_id,
				ImageUrl:        "Test image url 1",
			},
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: attraction_id,
				ImageUrl:        "Test image url 2",
			},
		},
		Location: entity.Location{
			LocationId:      uuid.New().String(),
			EstablishmentId: attraction_id,
			Address:         "test address",
			Latitude:        1.1,
			Longitude:       2.2,
			Country:         "Test country",
			City:            "Test city",
			StateProvince:   "Test state province",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	createdAttraction, err := repo.CreateAttraction(ctx, attraction)

	assert.NoError(t, err)
	assert.Equal(t, attraction.AttractionId, createdAttraction.AttractionId)
	assert.Equal(t, attraction.OwnerId, createdAttraction.OwnerId)
	assert.Equal(t, attraction.AttractionName, createdAttraction.AttractionName)
	assert.Equal(t, attraction.Description, createdAttraction.Description)
	assert.Equal(t, attraction.Rating, createdAttraction.Rating)
	assert.Equal(t, attraction.ContactNumber, createdAttraction.ContactNumber)
	assert.Equal(t, attraction.LicenceUrl, createdAttraction.LicenceUrl)
	assert.Equal(t, attraction.WebsiteUrl, createdAttraction.WebsiteUrl)
	assert.NotNil(t, createdAttraction.Images)
	assert.NotNil(t, createdAttraction.Location)
}

func TestGetAttraction(t *testing.T) {
	// Connect to database
	cfg := config.New()

	db, err := postgres.New(cfg)
	if err != nil {
		return
	}

	repo := NewAttractionRepo(db)

	attractionID := uuid.New().String()

	attraction := &entity.Attraction{
		AttractionId:   attractionID,
		OwnerId:        uuid.New().String(),
		AttractionName: "test attraction name",
		Description:    "Test description",
		Rating:         4.9,
		ContactNumber:  "+9989123456789",
		LicenceUrl:     "test licence url",
		WebsiteUrl:     "test website url",
		Images: []*entity.Image{
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: attractionID,
				ImageUrl:        "Test image url 1",
			},
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: attractionID,
				ImageUrl:        "Test image url 2",
			},
		},
		Location: entity.Location{
			LocationId:      uuid.New().String(),
			EstablishmentId: attractionID,
			Address:         "test address",
			Latitude:        1.1,
			Longitude:       2.2,
			Country:         "Test country",
			City:            "Test city",
			StateProvince:   "Test state province",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	_, err = repo.CreateAttraction(ctx, attraction)
	if err != nil {
		t.Fatalf("failed to insert attraction for testing: %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	gotAttraction, err := repo.GetAttraction(ctx, attractionID)

	assert.NoError(t, err)
	assert.NotNil(t, gotAttraction)
	assert.Equal(t, attraction.AttractionId, gotAttraction.AttractionId)
	assert.Equal(t, attraction.OwnerId, gotAttraction.OwnerId)
	assert.Equal(t, attraction.AttractionName, gotAttraction.AttractionName)
	assert.Equal(t, attraction.Description, gotAttraction.Description)
	assert.Equal(t, attraction.Rating, gotAttraction.Rating)
	assert.Equal(t, attraction.ContactNumber, gotAttraction.ContactNumber)
	assert.Equal(t, attraction.LicenceUrl, gotAttraction.LicenceUrl)
	assert.Equal(t, attraction.WebsiteUrl, gotAttraction.WebsiteUrl)

	assert.NotNil(t, gotAttraction.Location)
	assert.Equal(t, attraction.Location.LocationId, gotAttraction.Location.LocationId)
	assert.Equal(t, attraction.Location.EstablishmentId, gotAttraction.Location.EstablishmentId)
	assert.Equal(t, attraction.Location.Address, gotAttraction.Location.Address)
	assert.Equal(t, attraction.Location.Latitude, gotAttraction.Location.Latitude)
	assert.Equal(t, attraction.Location.Longitude, gotAttraction.Location.Longitude)
	assert.Equal(t, attraction.Location.Country, gotAttraction.Location.Country)
	assert.Equal(t, attraction.Location.City, gotAttraction.Location.City)
	assert.Equal(t, attraction.Location.StateProvince, gotAttraction.Location.StateProvince)

	assert.NotNil(t, gotAttraction.Images)
	assert.Len(t, gotAttraction.Images, len(attraction.Images))
	for i, expectedImage := range attraction.Images {
		assert.Equal(t, expectedImage.ImageId, gotAttraction.Images[i].ImageId)
		assert.Equal(t, expectedImage.EstablishmentId, gotAttraction.Images[i].EstablishmentId)
		assert.Equal(t, expectedImage.ImageUrl, gotAttraction.Images[i].ImageUrl)
	}
}

// func TestListAttractions(t *testing.T) {
// 	// Connect to database
// 	cfg := config.New()

// 	db, err := postgres.New(cfg)
// 	if err != nil {
// 		return
// 	}

// 	// Test Method ListAttractions
// 	repo := NewAttractionRepo(db)

// 	// Create sample attractions for testing
// 	var attractions []*entity.Attraction
// 	numAttractions := 5

// 	for i := 0; i < numAttractions; i++ {
// 		attractionID := uuid.New().String()
// 		attraction := &entity.Attraction{
// 			AttractionId:   attractionID,
// 			OwnerId:        uuid.New().String(),
// 			AttractionName: "test attraction name",
// 			Description:    "Test description",
// 			Rating:         4.9,
// 			ContactNumber:  "+9989123456789",
// 			LicenceUrl:     "test licence url",
// 			WebsiteUrl:     "test website url",
// 			Images: []*entity.Image{
// 				{
// 					ImageId:         uuid.New().String(),
// 					EstablishmentId: attractionID,
// 					ImageUrl:        "Test image url 1",
// 				},
// 				{
// 					ImageId:         uuid.New().String(),
// 					EstablishmentId: attractionID,
// 					ImageUrl:        "Test image url 2",
// 				},
// 			},
// 			Location: entity.Location{
// 				LocationId:      uuid.New().String(),
// 				EstablishmentId: attractionID,
// 				Address:         "test address",
// 				Latitude:        1.1,
// 				Longitude:       2.2,
// 				Country:         "Test country",
// 				City:            "Test city",
// 				StateProvince:   "Test state province",
// 			},
// 		}
// 		attractions = append(attractions, attraction)

// 		// Insert attraction into the database
// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
// 		defer cancel()

// 		_, err := repo.CreateAttraction(ctx, attraction)
// 		if err != nil {
// 			t.Fatalf("failed to insert attraction for testing: %v", err)
// 		}
// 	}

// 	// Test listing attractions
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
// 	defer cancel()

// 	page := int64(1)
// 	limit := int64(10)

// 	listedAttractions, err := repo.ListAttractions(ctx, page, limit)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, listedAttractions)
// 	assert.Len(t, listedAttractions, numAttractions)

// 	for i, expectedAttraction := range attractions {
// 		assert.Equal(t, expectedAttraction.AttractionId, listedAttractions[i].AttractionId)
// 		assert.Equal(t, expectedAttraction.OwnerId, listedAttractions[i].OwnerId)
// 		assert.Equal(t, expectedAttraction.AttractionName, listedAttractions[i].AttractionName)
// 		assert.Equal(t, expectedAttraction.Description, listedAttractions[i].Description)
// 		assert.Equal(t, expectedAttraction.Rating, listedAttractions[i].Rating)
// 		assert.Equal(t, expectedAttraction.ContactNumber, listedAttractions[i].ContactNumber)
// 		assert.Equal(t, expectedAttraction.LicenceUrl, listedAttractions[i].LicenceUrl)
// 		assert.Equal(t, expectedAttraction.WebsiteUrl, listedAttractions[i].WebsiteUrl)

// 		// Ensure location data is populated correctly
// 		assert.NotNil(t, listedAttractions[i].Location)
// 		assert.Equal(t, expectedAttraction.Location.LocationId, listedAttractions[i].Location.LocationId)
// 		assert.Equal(t, expectedAttraction.Location.EstablishmentId, listedAttractions[i].Location.EstablishmentId)
// 		assert.Equal(t, expectedAttraction.Location.Address, listedAttractions[i].Location.Address)
// 		assert.Equal(t, expectedAttraction.Location.Latitude, listedAttractions[i].Location.Latitude)
// 		assert.Equal(t, expectedAttraction.Location.Longitude, listedAttractions[i].Location.Longitude)
// 		assert.Equal(t, expectedAttraction.Location.Country, listedAttractions[i].Location.Country)
// 		assert.Equal(t, expectedAttraction.Location.City, listedAttractions[i].Location.City)
// 		assert.Equal(t, expectedAttraction.Location.StateProvince, listedAttractions[i].Location.StateProvince)

// 		// Ensure images data is populated correctly
// 		assert.NotNil(t, listedAttractions[i].Images)
// 		assert.Len(t, listedAttractions[i].Images, len(expectedAttraction.Images))
// 		for j, expectedImage := range expectedAttraction.Images {
// 			assert.Equal(t, expectedImage.ImageId, listedAttractions[i].Images[j].ImageId)
// 			assert.Equal(t, expectedImage.EstablishmentId, listedAttractions[i].Images[j].EstablishmentId)
// 			assert.Equal(t, expectedImage.ImageUrl, listedAttractions[i].Images[j].ImageUrl)
// 		}
// 	}
// }

func TestUpdateAttraction(t *testing.T) {
	// Connect to database
	cfg := config.New()

	db, err := postgres.New(cfg)
	if err != nil {
		return
	}

	repo := NewAttractionRepo(db)

	attractionID := uuid.New().String()

	attraction := &entity.Attraction{
		AttractionId:   attractionID,
		OwnerId:        uuid.New().String(),
		AttractionName: "test attraction name",
		Description:    "Test description",
		Rating:         4.9,
		ContactNumber:  "+9989123456789",
		LicenceUrl:     "test licence url",
		WebsiteUrl:     "test website url",
		Images: []*entity.Image{
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: attractionID,
				ImageUrl:        "Test image url 1",
			},
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: attractionID,
				ImageUrl:        "Test image url 2",
			},
		},
		Location: entity.Location{
			LocationId:      uuid.New().String(),
			EstablishmentId: attractionID,
			Address:         "test address",
			Latitude:        1.1,
			Longitude:       2.2,
			Country:         "Test country",
			City:            "Test city",
			StateProvince:   "Test state province",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	_, err = repo.CreateAttraction(ctx, attraction)
	if err != nil {
		t.Fatalf("failed to insert attraction for testing: %v", err)
	}

	// Update the sample attraction
	attraction.AttractionName = "updated attraction name"
	attraction.Description = "Updated description"
	attraction.ContactNumber = "+998976543210"
	attraction.LicenceUrl = "updated licence url"
	attraction.WebsiteUrl = "updated website url"

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	updatedAttraction, err := repo.UpdateAttraction(ctx, attraction)

	assert.NoError(t, err)
	assert.NotNil(t, updatedAttraction)
	assert.Equal(t, attraction.AttractionId, updatedAttraction.AttractionId)
	assert.Equal(t, attraction.OwnerId, updatedAttraction.OwnerId)
	assert.Equal(t, attraction.AttractionName, updatedAttraction.AttractionName)
	assert.Equal(t, attraction.Description, updatedAttraction.Description)
	assert.Equal(t, attraction.Rating, updatedAttraction.Rating)
	assert.Equal(t, attraction.ContactNumber, updatedAttraction.ContactNumber)
	assert.Equal(t, attraction.LicenceUrl, updatedAttraction.LicenceUrl)
	assert.Equal(t, attraction.WebsiteUrl, updatedAttraction.WebsiteUrl)

	// Ensure location data is populated correctly
	assert.NotNil(t, updatedAttraction.Location)
	assert.Equal(t, attraction.Location.LocationId, updatedAttraction.Location.LocationId)
	assert.Equal(t, attraction.Location.EstablishmentId, updatedAttraction.Location.EstablishmentId)
	assert.Equal(t, attraction.Location.Address, updatedAttraction.Location.Address)
	assert.Equal(t, attraction.Location.Latitude, updatedAttraction.Location.Latitude)
	assert.Equal(t, attraction.Location.Longitude, updatedAttraction.Location.Longitude)
	assert.Equal(t, attraction.Location.Country, updatedAttraction.Location.Country)
	assert.Equal(t, attraction.Location.City, updatedAttraction.Location.City)
	assert.Equal(t, attraction.Location.StateProvince, updatedAttraction.Location.StateProvince)

	// Ensure images data is populated correctly
	assert.NotNil(t, updatedAttraction.Images)
	assert.Len(t, updatedAttraction.Images, len(attraction.Images))
	for i, expectedImage := range attraction.Images {
		assert.Equal(t, expectedImage.ImageId, updatedAttraction.Images[i].ImageId)
		assert.Equal(t, expectedImage.EstablishmentId, updatedAttraction.Images[i].EstablishmentId)
		assert.Equal(t, expectedImage.ImageUrl, updatedAttraction.Images[i].ImageUrl)
	}
}

func TestDeleteAttraction(t *testing.T) {
	// Connect to database
	cfg := config.New()

	db, err := postgres.New(cfg)
	if err != nil {
		return
	}

	// Test Method DeleteAttraction
	repo := NewAttractionRepo(db)

	// Create a sample attraction for testing
	attractionID := uuid.New().String()

	// Insert sample attraction data into the database
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	_, err = repo.CreateAttraction(ctx, &entity.Attraction{
		AttractionId:   attractionID,
		OwnerId:        uuid.New().String(),
		AttractionName: "test attraction name",
		Description:    "Test description",
		Rating:         4.9,
		ContactNumber:  "+9989123456789",
		LicenceUrl:     "test licence url",
		WebsiteUrl:     "test website url",
		Images: []*entity.Image{
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: attractionID,
				ImageUrl:        "Test image url 1",
			},
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: attractionID,
				ImageUrl:        "Test image url 2",
			},
		},
		Location: entity.Location{
			LocationId:      uuid.New().String(),
			EstablishmentId: attractionID,
			Address:         "test address",
			Latitude:        1.1,
			Longitude:       2.2,
			Country:         "Test country",
			City:            "Test city",
			StateProvince:   "Test state province",
		},
	})
	if err != nil {
		t.Fatalf("failed to insert attraction for testing: %v", err)
	}

	// Test deleting the attraction
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	err = repo.DeleteAttraction(ctx, attractionID)

	assert.NoError(t, err)
}