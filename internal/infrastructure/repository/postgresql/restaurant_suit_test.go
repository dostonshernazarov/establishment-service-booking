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

func TestCreateRestaurant(t *testing.T) {
	// Connect to database
	cfg := config.New()

	db, err := postgres.New(cfg)
	if err != nil {
		return
	}

	// Test  Method Create
	repo := NewRestaurantRepo(db)

	restaurant_id := uuid.New().String()

	restaurant := &entity.Restaurant{
		RestaurantId:   restaurant_id,
		OwnerId:        uuid.New().String(),
		RestaurantName: "test restaurant name",
		Description:    "Test description",
		Rating:         4.9,
		OpeningHours:   "09:00 - 00:00",
		ContactNumber:  "+9989123456789",
		LicenceUrl:     "test licence url",
		WebsiteUrl:     "test website url",
		Images: []*entity.Image{
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: restaurant_id,
				ImageUrl:        "Test image url 1",
			},
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: restaurant_id,
				ImageUrl:        "Test image url 2",
			},
		},
		Location: entity.Location{
			LocationId:      uuid.New().String(),
			EstablishmentId: restaurant_id,
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

	createdrestaurant, err := repo.CreateRestaurant(ctx, restaurant)

	assert.NoError(t, err)
	assert.Equal(t, restaurant.RestaurantId, createdrestaurant.RestaurantId)
	assert.Equal(t, restaurant.OwnerId, createdrestaurant.OwnerId)
	assert.Equal(t, restaurant.RestaurantName, createdrestaurant.RestaurantName)
	assert.Equal(t, restaurant.Description, createdrestaurant.Description)
	assert.Equal(t, restaurant.Rating, createdrestaurant.Rating)
	assert.Equal(t, restaurant.ContactNumber, createdrestaurant.ContactNumber)
	assert.Equal(t, restaurant.LicenceUrl, createdrestaurant.LicenceUrl)
	assert.Equal(t, restaurant.WebsiteUrl, createdrestaurant.WebsiteUrl)
	assert.NotNil(t, createdrestaurant.Images)
	assert.NotNil(t, createdrestaurant.Location)
}

func TestGetRestaurant(t *testing.T) {
	// Connect to database
	cfg := config.New()

	db, err := postgres.New(cfg)
	if err != nil {
		return
	}

	repo := NewRestaurantRepo(db)

	restaurant_id := uuid.New().String()

	restaurant := &entity.Restaurant{
		RestaurantId:   restaurant_id,
		OwnerId:        uuid.New().String(),
		RestaurantName: "test restaurant name",
		Description:    "Test description",
		Rating:         4.9,
		OpeningHours:   "09:00 - 00:00",
		ContactNumber:  "+9989123456789",
		LicenceUrl:     "test licence url",
		WebsiteUrl:     "test website url",
		Images: []*entity.Image{
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: restaurant_id,
				ImageUrl:        "Test image url 1",
			},
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: restaurant_id,
				ImageUrl:        "Test image url 2",
			},
		},
		Location: entity.Location{
			LocationId:      uuid.New().String(),
			EstablishmentId: restaurant_id,
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

	_, err = repo.CreateRestaurant(ctx, restaurant)
	if err != nil {
		t.Fatalf("failed to insert restaurant for testing: %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	gotrestaurant, err := repo.GetRestaurant(ctx, restaurant_id)

	assert.NoError(t, err)
	assert.NotNil(t, gotrestaurant)
	assert.Equal(t, restaurant.RestaurantId, gotrestaurant.RestaurantId)
	assert.Equal(t, restaurant.OwnerId, gotrestaurant.OwnerId)
	assert.Equal(t, restaurant.RestaurantName, gotrestaurant.RestaurantName)
	assert.Equal(t, restaurant.Description, gotrestaurant.Description)
	assert.Equal(t, restaurant.Rating, gotrestaurant.Rating)
	assert.Equal(t, restaurant.ContactNumber, gotrestaurant.ContactNumber)
	assert.Equal(t, restaurant.LicenceUrl, gotrestaurant.LicenceUrl)
	assert.Equal(t, restaurant.WebsiteUrl, gotrestaurant.WebsiteUrl)

	assert.NotNil(t, gotrestaurant.Location)
	assert.Equal(t, restaurant.Location.LocationId, gotrestaurant.Location.LocationId)
	assert.Equal(t, restaurant.Location.EstablishmentId, gotrestaurant.Location.EstablishmentId)
	assert.Equal(t, restaurant.Location.Address, gotrestaurant.Location.Address)
	assert.Equal(t, restaurant.Location.Latitude, gotrestaurant.Location.Latitude)
	assert.Equal(t, restaurant.Location.Longitude, gotrestaurant.Location.Longitude)
	assert.Equal(t, restaurant.Location.Country, gotrestaurant.Location.Country)
	assert.Equal(t, restaurant.Location.City, gotrestaurant.Location.City)
	assert.Equal(t, restaurant.Location.StateProvince, gotrestaurant.Location.StateProvince)

	assert.NotNil(t, gotrestaurant.Images)
	assert.Len(t, gotrestaurant.Images, len(restaurant.Images))
	for i, expectedImage := range restaurant.Images {
		assert.Equal(t, expectedImage.ImageId, gotrestaurant.Images[i].ImageId)
		assert.Equal(t, expectedImage.EstablishmentId, gotrestaurant.Images[i].EstablishmentId)
		assert.Equal(t, expectedImage.ImageUrl, gotrestaurant.Images[i].ImageUrl)
	}
}

// func TestListRestaurants(t *testing.T) {
// 	// Connect to database
// 	cfg := config.New()

// 	db, err := postgres.New(cfg)
// 	if err != nil {
// 		return
// 	}

// 	repo := NewRestaurantRepo(db)

// 	var restaurants []*entity.Restaurant
// 	numrestaurants := 5

// 	for i := 0; i < numrestaurants; i++ {
// 		restaurant_id := uuid.New().String()
// 		restaurant := &entity.Restaurant{
// 			RestaurantId:   restaurant_id,
// 			OwnerId:        uuid.New().String(),
// 			RestaurantName: "test restaurant name",
// 			Description:    "Test description",
// 			Rating:         4.9,
// 			OpeningHours:   "09:00 - 00:00",
// 			ContactNumber:  "+9989123456789",
// 			LicenceUrl:     "test licence url",
// 			WebsiteUrl:     "test website url",
// 			Images: []*entity.Image{
// 				{
// 					ImageId:         uuid.New().String(),
// 					EstablishmentId: restaurant_id,
// 					ImageUrl:        "Test image url 1",
// 				},
// 				{
// 					ImageId:         uuid.New().String(),
// 					EstablishmentId: restaurant_id,
// 					ImageUrl:        "Test image url 2",
// 				},
// 			},
// 			Location: entity.Location{
// 				LocationId:      uuid.New().String(),
// 				EstablishmentId: restaurant_id,
// 				Address:         "test address",
// 				Latitude:        1.1,
// 				Longitude:       2.2,
// 				Country:         "Test country",
// 				City:            "Test city",
// 				StateProvince:   "Test state province",
// 			},
// 		}
// 		restaurants = append(restaurants, restaurant)

// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
// 		defer cancel()

// 		_, err := repo.CreateRestaurant(ctx, restaurant)
// 		if err != nil {
// 			t.Fatalf("failed to insert restaurant for testing: %v", err)
// 		}
// 	}

// 	// Test listing restaurants
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
// 	defer cancel()

// 	offset := int64(1)
// 	limit := int64(10)

// 	listedrestaurants, err := repo.ListRestaurants(ctx, offset, limit)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, listedrestaurants)
// 	assert.Len(t, listedrestaurants, numrestaurants)

// 	for i, expectedrestaurant := range restaurants {
// 		assert.Equal(t, expectedrestaurant.RestaurantId, listedrestaurants[i].RestaurantId)
// 		assert.Equal(t, expectedrestaurant.OwnerId, listedrestaurants[i].OwnerId)
// 		assert.Equal(t, expectedrestaurant.RestaurantName, listedrestaurants[i].RestaurantName)
// 		assert.Equal(t, expectedrestaurant.Description, listedrestaurants[i].Description)
// 		assert.Equal(t, expectedrestaurant.Rating, listedrestaurants[i].Rating)
// 		assert.Equal(t, expectedrestaurant.ContactNumber, listedrestaurants[i].ContactNumber)
// 		assert.Equal(t, expectedrestaurant.LicenceUrl, listedrestaurants[i].LicenceUrl)
// 		assert.Equal(t, expectedrestaurant.WebsiteUrl, listedrestaurants[i].WebsiteUrl)

// 		// Ensure location data is populated correctly
// 		assert.NotNil(t, listedrestaurants[i].Location)
// 		assert.Equal(t, expectedrestaurant.Location.LocationId, listedrestaurants[i].Location.LocationId)
// 		assert.Equal(t, expectedrestaurant.Location.EstablishmentId, listedrestaurants[i].Location.EstablishmentId)
// 		assert.Equal(t, expectedrestaurant.Location.Address, listedrestaurants[i].Location.Address)
// 		assert.Equal(t, expectedrestaurant.Location.Latitude, listedrestaurants[i].Location.Latitude)
// 		assert.Equal(t, expectedrestaurant.Location.Longitude, listedrestaurants[i].Location.Longitude)
// 		assert.Equal(t, expectedrestaurant.Location.Country, listedrestaurants[i].Location.Country)
// 		assert.Equal(t, expectedrestaurant.Location.City, listedrestaurants[i].Location.City)
// 		assert.Equal(t, expectedrestaurant.Location.StateProvince, listedrestaurants[i].Location.StateProvince)

// 		// Ensure images data is populated correctly
// 		assert.NotNil(t, listedrestaurants[i].Images)
// 		assert.Len(t, listedrestaurants[i].Images, len(expectedrestaurant.Images))
// 		for j, expectedImage := range expectedrestaurant.Images {
// 			assert.Equal(t, expectedImage.ImageId, listedrestaurants[i].Images[j].ImageId)
// 			assert.Equal(t, expectedImage.EstablishmentId, listedrestaurants[i].Images[j].EstablishmentId)
// 			assert.Equal(t, expectedImage.ImageUrl, listedrestaurants[i].Images[j].ImageUrl)
// 		}
// 	}
// }

func TestUpdaterestaurant(t *testing.T) {
	// Connect to database
	cfg := config.New()

	db, err := postgres.New(cfg)
	if err != nil {
		return
	}

	repo := NewRestaurantRepo(db)

	restaurant_id := uuid.New().String()

	restaurant := &entity.Restaurant{
		RestaurantId:   restaurant_id,
		OwnerId:        uuid.New().String(),
		RestaurantName: "test restaurant name",
		Description:    "Test description",
		Rating:         4.9,
		OpeningHours:   "09:00 - 00:00",
		ContactNumber:  "+9989123456789",
		LicenceUrl:     "test licence url",
		WebsiteUrl:     "test website url",
		Images: []*entity.Image{
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: restaurant_id,
				ImageUrl:        "Test image url 1",
			},
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: restaurant_id,
				ImageUrl:        "Test image url 2",
			},
		},
		Location: entity.Location{
			LocationId:      uuid.New().String(),
			EstablishmentId: restaurant_id,
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

	_, err = repo.CreateRestaurant(ctx, restaurant)
	if err != nil {
		t.Fatalf("failed to insert restaurant for testing: %v", err)
	}

	// Update the sample restaurant
	restaurant.RestaurantName = "updated restaurant name"
	restaurant.Description = "Updated description"
	restaurant.ContactNumber = "+998976543210"
	restaurant.LicenceUrl = "updated licence url"
	restaurant.WebsiteUrl = "updated website url"

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	updatedrestaurant, err := repo.UpdateRestaurant(ctx, restaurant)

	assert.NoError(t, err)
	assert.NotNil(t, updatedrestaurant)
	assert.Equal(t, restaurant.RestaurantId, updatedrestaurant.RestaurantId)
	assert.Equal(t, restaurant.OwnerId, updatedrestaurant.OwnerId)
	assert.Equal(t, restaurant.RestaurantName, updatedrestaurant.RestaurantName)
	assert.Equal(t, restaurant.Description, updatedrestaurant.Description)
	assert.Equal(t, restaurant.Rating, updatedrestaurant.Rating)
	assert.Equal(t, restaurant.ContactNumber, updatedrestaurant.ContactNumber)
	assert.Equal(t, restaurant.LicenceUrl, updatedrestaurant.LicenceUrl)
	assert.Equal(t, restaurant.WebsiteUrl, updatedrestaurant.WebsiteUrl)

	// Ensure location data is populated correctly
	assert.NotNil(t, updatedrestaurant.Location)
	assert.Equal(t, restaurant.Location.LocationId, updatedrestaurant.Location.LocationId)
	assert.Equal(t, restaurant.Location.EstablishmentId, updatedrestaurant.Location.EstablishmentId)
	assert.Equal(t, restaurant.Location.Address, updatedrestaurant.Location.Address)
	assert.Equal(t, restaurant.Location.Latitude, updatedrestaurant.Location.Latitude)
	assert.Equal(t, restaurant.Location.Longitude, updatedrestaurant.Location.Longitude)
	assert.Equal(t, restaurant.Location.Country, updatedrestaurant.Location.Country)
	assert.Equal(t, restaurant.Location.City, updatedrestaurant.Location.City)
	assert.Equal(t, restaurant.Location.StateProvince, updatedrestaurant.Location.StateProvince)

	// Ensure images data is populated correctly
	assert.NotNil(t, updatedrestaurant.Images)
	assert.Len(t, updatedrestaurant.Images, len(restaurant.Images))
	for i, expectedImage := range restaurant.Images {
		assert.Equal(t, expectedImage.ImageId, updatedrestaurant.Images[i].ImageId)
		assert.Equal(t, expectedImage.EstablishmentId, updatedrestaurant.Images[i].EstablishmentId)
		assert.Equal(t, expectedImage.ImageUrl, updatedrestaurant.Images[i].ImageUrl)
	}
}

func TestDeleterestaurant(t *testing.T) {
	// Connect to database
	cfg := config.New()

	db, err := postgres.New(cfg)
	if err != nil {
		return
	}

	repo := NewRestaurantRepo(db)

	// Create a sample restaurant for testing
	restaurant_id := uuid.New().String()

	// Insert sample restaurant data into the database
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	_, err = repo.CreateRestaurant(ctx, &entity.Restaurant{
		RestaurantId:   restaurant_id,
		OwnerId:        uuid.New().String(),
		RestaurantName: "test restaurant name",
		Description:    "Test description",
		Rating:         4.9,
		OpeningHours:   "09:00",
		ContactNumber:  "+9989123456789",
		LicenceUrl:     "test licence url",
		WebsiteUrl:     "test website url",
		Images: []*entity.Image{
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: restaurant_id,
				ImageUrl:        "Test image url 1",
			},
			{
				ImageId:         uuid.New().String(),
				EstablishmentId: restaurant_id,
				ImageUrl:        "Test image url 2",
			},
		},
		Location: entity.Location{
			LocationId:      uuid.New().String(),
			EstablishmentId: restaurant_id,
			Address:         "test address",
			Latitude:        1.1,
			Longitude:       2.2,
			Country:         "Test country",
			City:            "Test city",
			StateProvince:   "Test state province",
		},
	})
	if err != nil {
		t.Fatalf("failed to insert restaurant for testing: %v", err)
	}

	// Test deleting the restaurant
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	err = repo.DeleteRestaurant(ctx, restaurant_id)

	assert.NoError(t, err)
}
