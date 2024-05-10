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

func TestCreateHotel(t *testing.T) {
	// Connect to database
	cfg := config.New()

	db, err := postgres.New(cfg)
	if err != nil {
		return
	}

	// Test  Method Create
	repo := NewHotelRepo(db)

	hotel_id := uuid.New().String()

	hotel := &entity.Hotel{
		HotelId:       hotel_id,
		OwnerId:       uuid.New().String(),
		HotelName:     "test hotel name",
		Description:   "Test description",
		Rating:        4.9,
		ContactNumber: "+9989123456789",
		LicenceUrl:    "test licence url",
		WebsiteUrl:    "test website url",
		Images: []*entity.Image{
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: hotel_id,
				ImageUrl:        "Test image url 1",
			},
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: hotel_id,
				ImageUrl:        "Test image url 2",
			},
		},
		Location: entity.Location{
			LocationId:      uuid.New().String(),
			EstablishmentId: hotel_id,
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

	createdHotel, err := repo.CreateHotel(ctx, hotel)

	assert.NoError(t, err)
	assert.Equal(t, hotel.HotelId, createdHotel.HotelId)
	assert.Equal(t, hotel.OwnerId, createdHotel.OwnerId)
	assert.Equal(t, hotel.HotelName, createdHotel.HotelName)
	assert.Equal(t, hotel.Description, createdHotel.Description)
	assert.Equal(t, hotel.Rating, createdHotel.Rating)
	assert.Equal(t, hotel.ContactNumber, createdHotel.ContactNumber)
	assert.Equal(t, hotel.LicenceUrl, createdHotel.LicenceUrl)
	assert.Equal(t, hotel.WebsiteUrl, createdHotel.WebsiteUrl)
	assert.NotNil(t, createdHotel.Images)
	assert.NotNil(t, createdHotel.Location)
}

func TestGetHotel(t *testing.T) {
	// Connect to database
	cfg := config.New()

	db, err := postgres.New(cfg)
	if err != nil {
		return
	}

	repo := NewHotelRepo(db)

	hotel_id := uuid.New().String()

	hotel := &entity.Hotel{
		HotelId:       hotel_id,
		OwnerId:       uuid.New().String(),
		HotelName:     "test hotel name",
		Description:   "Test description",
		Rating:        4.9,
		ContactNumber: "+9989123456789",
		LicenceUrl:    "test licence url",
		WebsiteUrl:    "test website url",
		Images: []*entity.Image{
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: hotel_id,
				ImageUrl:        "Test image url 1",
			},
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: hotel_id,
				ImageUrl:        "Test image url 2",
			},
		},
		Location: entity.Location{
			LocationId:      uuid.New().String(),
			EstablishmentId: hotel_id,
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

	_, err = repo.CreateHotel(ctx, hotel)
	if err != nil {
		t.Fatalf("failed to insert hotel for testing: %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	gotHotel, err := repo.GetHotel(ctx, hotel_id)

	assert.NoError(t, err)
	assert.NotNil(t, gotHotel)
	assert.Equal(t, hotel.HotelId, gotHotel.HotelId)
	assert.Equal(t, hotel.OwnerId, gotHotel.OwnerId)
	assert.Equal(t, hotel.HotelName, gotHotel.HotelName)
	assert.Equal(t, hotel.Description, gotHotel.Description)
	assert.Equal(t, hotel.Rating, gotHotel.Rating)
	assert.Equal(t, hotel.ContactNumber, gotHotel.ContactNumber)
	assert.Equal(t, hotel.LicenceUrl, gotHotel.LicenceUrl)
	assert.Equal(t, hotel.WebsiteUrl, gotHotel.WebsiteUrl)

	assert.NotNil(t, gotHotel.Location)
	assert.Equal(t, hotel.Location.LocationId, gotHotel.Location.LocationId)
	assert.Equal(t, hotel.Location.EstablishmentId, gotHotel.Location.EstablishmentId)
	assert.Equal(t, hotel.Location.Address, gotHotel.Location.Address)
	assert.Equal(t, hotel.Location.Latitude, gotHotel.Location.Latitude)
	assert.Equal(t, hotel.Location.Longitude, gotHotel.Location.Longitude)
	assert.Equal(t, hotel.Location.Country, gotHotel.Location.Country)
	assert.Equal(t, hotel.Location.City, gotHotel.Location.City)
	assert.Equal(t, hotel.Location.StateProvince, gotHotel.Location.StateProvince)

	assert.NotNil(t, gotHotel.Images)
	assert.Len(t, gotHotel.Images, len(hotel.Images))
	for i, expectedImage := range hotel.Images {
		assert.Equal(t, expectedImage.ImageId, gotHotel.Images[i].ImageId)
		assert.Equal(t, expectedImage.EstablishmentId, gotHotel.Images[i].EstablishmentId)
		assert.Equal(t, expectedImage.ImageUrl, gotHotel.Images[i].ImageUrl)
	}
}

// func TestListHotels(t *testing.T) {
// 	// Connect to database
// 	cfg := config.New()

// 	db, err := postgres.New(cfg)
// 	if err != nil {
// 		return
// 	}

// 	repo := NewHotelRepo(db)

// 	var hotels []*entity.Hotel
// 	numHotels := 5

// 	for i := 0; i < numHotels; i++ {
// 		hotel_id := uuid.New().String()
// 		hotel := &entity.Hotel{
// 			HotelId:       hotel_id,
// 			OwnerId:       uuid.New().String(),
// 			HotelName:     "test hotel name",
// 			Description:   "Test description",
// 			Rating:        4.9,
// 			ContactNumber: "+9989123456789",
// 			LicenceUrl:    "test licence url",
// 			WebsiteUrl:    "test website url",
// 			Images: []*entity.Image{
// 				{
// 					ImageId:         uuid.New().String(),
// 					EstablishmentId: hotel_id,
// 					ImageUrl:        "Test image url 1",
// 				},
// 				{
// 					ImageId:         uuid.New().String(),
// 					EstablishmentId: hotel_id,
// 					ImageUrl:        "Test image url 2",
// 				},
// 			},
// 			Location: entity.Location{
// 				LocationId:      uuid.New().String(),
// 				EstablishmentId: hotel_id,
// 				Address:         "test address",
// 				Latitude:        1.1,
// 				Longitude:       2.2,
// 				Country:         "Test country",
// 				City:            "Test city",
// 				StateProvince:   "Test state province",
// 			},
// 		}
// 		hotels = append(hotels, hotel)

// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
// 		defer cancel()

// 		_, err := repo.CreateHotel(ctx, hotel)
// 		if err != nil {
// 			t.Fatalf("failed to insert hotel for testing: %v", err)
// 		}
// 	}

// 	// Test listing hotels
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
// 	defer cancel()

// 	offset := int64(1)
// 	limit := int64(10)

// 	listedHotels, err := repo.ListHotels(ctx, offset, limit)

// 	assert.NoError(t, err)
// 	// assert.NotNil(t, listedHotels)
// 	// assert.Len(t, listedHotels, numHotels)

// 	for i, expectedHotel := range hotels {
// 		assert.Equal(t, expectedHotel.HotelId, listedHotels[i].HotelId)
// 		assert.Equal(t, expectedHotel.OwnerId, listedHotels[i].OwnerId)
// 		assert.Equal(t, expectedHotel.HotelName, listedHotels[i].HotelName)
// 		assert.Equal(t, expectedHotel.Description, listedHotels[i].Description)
// 		assert.Equal(t, expectedHotel.Rating, listedHotels[i].Rating)
// 		assert.Equal(t, expectedHotel.ContactNumber, listedHotels[i].ContactNumber)
// 		assert.Equal(t, expectedHotel.LicenceUrl, listedHotels[i].LicenceUrl)
// 		assert.Equal(t, expectedHotel.WebsiteUrl, listedHotels[i].WebsiteUrl)

// 		// Ensure location data is populated correctly
// 		assert.NotNil(t, listedHotels[i].Location)
// 		assert.Equal(t, expectedHotel.Location.LocationId, listedHotels[i].Location.LocationId)
// 		assert.Equal(t, expectedHotel.Location.EstablishmentId, listedHotels[i].Location.EstablishmentId)
// 		assert.Equal(t, expectedHotel.Location.Address, listedHotels[i].Location.Address)
// 		assert.Equal(t, expectedHotel.Location.Latitude, listedHotels[i].Location.Latitude)
// 		assert.Equal(t, expectedHotel.Location.Longitude, listedHotels[i].Location.Longitude)
// 		assert.Equal(t, expectedHotel.Location.Country, listedHotels[i].Location.Country)
// 		assert.Equal(t, expectedHotel.Location.City, listedHotels[i].Location.City)
// 		assert.Equal(t, expectedHotel.Location.StateProvince, listedHotels[i].Location.StateProvince)

// 		// Ensure images data is populated correctly
// 		assert.NotNil(t, listedHotels[i].Images)
// 		assert.Len(t, listedHotels[i].Images, len(expectedHotel.Images))
// 		for j, expectedImage := range expectedHotel.Images {
// 			assert.Equal(t, expectedImage.ImageId, listedHotels[i].Images[j].ImageId)
// 			assert.Equal(t, expectedImage.EstablishmentId, listedHotels[i].Images[j].EstablishmentId)
// 			assert.Equal(t, expectedImage.ImageUrl, listedHotels[i].Images[j].ImageUrl)
// 		}
// 	}
// }

func TestUpdateHotel(t *testing.T) {
	// Connect to database
	cfg := config.New()

	db, err := postgres.New(cfg)
	if err != nil {
		return
	}

	repo := NewHotelRepo(db)

	hotel_id := uuid.New().String()

	hotel := &entity.Hotel{
		HotelId:       hotel_id,
		OwnerId:       uuid.New().String(),
		HotelName:     "test hotel name",
		Description:   "Test description",
		Rating:        4.9,
		ContactNumber: "+9989123456789",
		LicenceUrl:    "test licence url",
		WebsiteUrl:    "test website url",
		Images: []*entity.Image{
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: hotel_id,
				ImageUrl:        "Test image url 1",
			},
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: hotel_id,
				ImageUrl:        "Test image url 2",
			},
		},
		Location: entity.Location{
			LocationId:      uuid.New().String(),
			EstablishmentId: hotel_id,
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

	_, err = repo.CreateHotel(ctx, hotel)
	if err != nil {
		t.Fatalf("failed to insert hotel for testing: %v", err)
	}

	// Update the sample hotel
	hotel.HotelName = "updated hotel name"
	hotel.Description = "Updated description"
	hotel.ContactNumber = "+998976543210"
	hotel.LicenceUrl = "updated licence url"
	hotel.WebsiteUrl = "updated website url"

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	updatedHotel, err := repo.UpdateHotel(ctx, hotel)

	assert.NoError(t, err)
	assert.NotNil(t, updatedHotel)
	assert.Equal(t, hotel.HotelId, updatedHotel.HotelId)
	assert.Equal(t, hotel.OwnerId, updatedHotel.OwnerId)
	assert.Equal(t, hotel.HotelName, updatedHotel.HotelName)
	assert.Equal(t, hotel.Description, updatedHotel.Description)
	assert.Equal(t, hotel.Rating, updatedHotel.Rating)
	assert.Equal(t, hotel.ContactNumber, updatedHotel.ContactNumber)
	assert.Equal(t, hotel.LicenceUrl, updatedHotel.LicenceUrl)
	assert.Equal(t, hotel.WebsiteUrl, updatedHotel.WebsiteUrl)

	// Ensure location data is populated correctly
	assert.NotNil(t, updatedHotel.Location)
	assert.Equal(t, hotel.Location.LocationId, updatedHotel.Location.LocationId)
	assert.Equal(t, hotel.Location.EstablishmentId, updatedHotel.Location.EstablishmentId)
	assert.Equal(t, hotel.Location.Address, updatedHotel.Location.Address)
	assert.Equal(t, hotel.Location.Latitude, updatedHotel.Location.Latitude)
	assert.Equal(t, hotel.Location.Longitude, updatedHotel.Location.Longitude)
	assert.Equal(t, hotel.Location.Country, updatedHotel.Location.Country)
	assert.Equal(t, hotel.Location.City, updatedHotel.Location.City)
	assert.Equal(t, hotel.Location.StateProvince, updatedHotel.Location.StateProvince)

	// Ensure images data is populated correctly
	assert.NotNil(t, updatedHotel.Images)
	assert.Len(t, updatedHotel.Images, len(hotel.Images))
	for i, expectedImage := range hotel.Images {
		assert.Equal(t, expectedImage.ImageId, updatedHotel.Images[i].ImageId)
		assert.Equal(t, expectedImage.EstablishmentId, updatedHotel.Images[i].EstablishmentId)
		assert.Equal(t, expectedImage.ImageUrl, updatedHotel.Images[i].ImageUrl)
	}
}

func TestDeleteHotel(t *testing.T) {
	// Connect to database
	cfg := config.New()

	db, err := postgres.New(cfg)
	if err != nil {
		return
	}

	repo := NewHotelRepo(db)

	// Create a sample hotel for testing
	hotel_id := uuid.New().String()

	// Insert sample hotel data into the database
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	_, err = repo.CreateHotel(ctx, &entity.Hotel{
		HotelId:       hotel_id,
		OwnerId:       uuid.New().String(),
		HotelName:     "test hotel name",
		Description:   "Test description",
		Rating:        4.9,
		ContactNumber: "+9989123456789",
		LicenceUrl:    "test licence url",
		WebsiteUrl:    "test website url",
		Images: []*entity.Image{
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: hotel_id,
				ImageUrl:        "Test image url 1",
			},
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: hotel_id,
				ImageUrl:        "Test image url 2",
			},
		},
		Location: entity.Location{
			LocationId:      uuid.New().String(),
			EstablishmentId: hotel_id,
			Address:         "test address",
			Latitude:        1.1,
			Longitude:       2.2,
			Country:         "Test country",
			City:            "Test city",
			StateProvince:   "Test state province",
		},
	})
	if err != nil {
		t.Fatalf("failed to insert hotel for testing: %v", err)
	}

	// Test deleting the hotel
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	err = repo.DeleteHotel(ctx, hotel_id)

	assert.NoError(t, err)
}
