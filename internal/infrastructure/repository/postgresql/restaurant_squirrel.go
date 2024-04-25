package postgresql

import (
	"Booking/establishment-service-booking/internal/entity"
	"Booking/establishment-service-booking/internal/pkg/postgres"
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

const (
	restaurantTableName = "restaurant_table" //table for storing general info of attraction
	// locationTableName        = "location_table"   // table for storing location info
	// imageTableName           = "image_table"      // table for storing multiple images of establishment
	// restaurantServiceName    = "restaurantService"
	// restaurantSpanRepoPrefix = "attractionRepo"
)

type restaurantRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewRestaurantRepo(db *postgres.PostgresDB) *restaurantRepo {
	return &restaurantRepo{
		tableName: restaurantTableName,
		db:        db,
	}
}

func (p *restaurantRepo) RestaurantSelectQueryPrefix() squirrel.SelectBuilder {
	return p.db.Sq.Builder.Select(
		"restaurant_id",
		"owner_id",
		"restaurant_name",
		"description",
		"rating",
		"opening_hours",
		"contact_number",
		"licence_url",
		"website_url",
		"created_at",
		"updated_at",
		"deleted_at",
	).From(p.tableName)
}

// create a new restaurant
func (p restaurantRepo) CreateRestaurant(ctx context.Context, restaurant *entity.Restaurant) (*entity.Restaurant, error) {

	// insert location info to location_table
	dataL := map[string]interface{}{
		"location_id":      restaurant.Location.LocationId,
		"establishment_id": restaurant.Location.EstablishmentId,
		"address":          restaurant.Location.Address,
		"latitude":         restaurant.Location.Latitude,
		"longitude":        restaurant.Location.Longitude,
		"country":          restaurant.Location.Country,
		"city":             restaurant.Location.City,
		"state_province":   restaurant.Location.StateProvince,
		"created_at":       restaurant.Location.CreatedAt,
		"updated_at":       restaurant.Location.UpdatedAt,
		"deleted_at":       restaurant.Location.DeletedAt,
	}

	query, args, err := p.db.Sq.Builder.Insert(locationTableName).SetMap(dataL).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for creating restaurant's location part: %v", err)
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query for creating restaurant's location part: %v", err)
	}

	// insert images to image_table
	for _, image := range restaurant.Images {
		dataI := map[string]interface{}{
			"image_id":         image.ImageId,
			"establishment_id": image.EstablishmentId,
			"image_url":        image.ImageUrl,
			"created_at":       image.CreatedAt,
			"updated_at":       image.UpdatedAt,
			"deleted_at":       image.DeletedAt,
		}

		query, args, err := p.db.Sq.Builder.Insert(imageTableName).SetMap(dataI).ToSql()
		if err != nil {
			return nil, fmt.Errorf("failed to build SQL query for creating image: %v", err)
		}

		_, err = p.db.Exec(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to execute SQL query for creating image: %v", err)
		}
	}

	// insert general info of attraction
	data := map[string]interface{}{
		"restaurant_id":   restaurant.RestaurantId,
		"owner_id":        restaurant.OwnerId,
		"restaurant_name": restaurant.RestaurantName,
		"description":     restaurant.Description,
		"rating":          restaurant.Rating,
		"opening_hours":   restaurant.OpeningHours,
		"contact_number":  restaurant.ContactNumber,
		"licence_url":     restaurant.LicenceUrl,
		"website_url":     restaurant.WebsiteUrl,
		"created_at":      restaurant.CreatedAt,
		"updated_at":      restaurant.UpdatedAt,
		"deleted_at":      restaurant.DeletedAt,
	}
	query, args, err = p.db.Sq.Builder.Insert(p.tableName).SetMap(data).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for creating restaurant: %v", err)
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query for creating restaurant: %v", err)
	}

	return nil, nil
}

// get a restaurant
func (p restaurantRepo) GetRestaurant(ctx context.Context, restaurant_id string) (*entity.Restaurant, error) {
	var restaurant entity.Restaurant

	// Build the query to select attraction details
	queryBuilder := p.RestaurantSelectQueryPrefix().Where(p.db.Sq.Equal("restaurant_id", restaurant_id))

	// Get the SQL query and arguments from the query builder
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for getting restaurant: %v", err)
	}

	// Execute the query to fetch restaurant details
	if err := p.db.QueryRow(ctx, query, args...).Scan(
		&restaurant.RestaurantId,
		&restaurant.OwnerId,
		&restaurant.RestaurantName,
		&restaurant.Description,
		&restaurant.Rating,
		&restaurant.OpeningHours,
		&restaurant.ContactNumber,
		&restaurant.LicenceUrl,
		&restaurant.WebsiteUrl,
		&restaurant.CreatedAt,
		&restaurant.UpdatedAt,
		&restaurant.DeletedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get restaurant: %v", err)
	}

	// Fetch location information
	locationQuery := fmt.Sprintf("SELECT * FROM %s WHERE establishment_id = $1", locationTableName)
	if err := p.db.QueryRow(ctx, locationQuery, restaurant.RestaurantId).Scan(
		&restaurant.Location.LocationId,
		&restaurant.Location.EstablishmentId,
		&restaurant.Location.Address,
		&restaurant.Location.Latitude,
		&restaurant.Location.Longitude,
		&restaurant.Location.Country,
		&restaurant.Location.City,
		&restaurant.Location.StateProvince,
		&restaurant.Location.CreatedAt,
		&restaurant.Location.UpdatedAt,
		&restaurant.Location.DeletedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get location for location: %v", err)
	}

	// Fetch images information
	imagesQuery := fmt.Sprintf("SELECT * FROM %s WHERE establishment_id = $1", imageTableName)
	rows, err := p.db.Query(ctx, imagesQuery, restaurant_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get images for restaurant: %v", err)
	}
	defer rows.Close()

	// Iterate over the rows and populate the Images slice
	for rows.Next() {
		var image entity.Image
		if err := rows.Scan(
			&image.ImageId,
			&image.EstablishmentId,
			&image.ImageUrl,
			&image.CreatedAt,
			&image.UpdatedAt,
			&image.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan image row: %v", err)
		}
		restaurant.Images = append(restaurant.Images, &image)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered while iterating over image rows: %v", err)
	}

	return &restaurant, nil
}

// get a list of restaurants
func (p restaurantRepo) ListRestaurants(ctx context.Context, offset, limit int64) ([]*entity.Restaurant, error) {
	var restaurants []*entity.Restaurant

	queryBuilder := p.RestaurantSelectQueryPrefix()

	if limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(limit)).Offset(uint64(offset)).Where(p.db.Sq.Equal("deleted_at", nil))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for listing restaurants: %v", err)
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query for listing restaurants: %v", err)
	}
	defer rows.Close()

	// Iterate over the rows to fetch each restaurants's details
	for rows.Next() {
		var restaurant entity.Restaurant
		if err := rows.Scan(
			&restaurant.RestaurantId,
			&restaurant.OwnerId,
			&restaurant.RestaurantName,
			&restaurant.Description,
			&restaurant.Rating,
			&restaurant.OpeningHours,
			&restaurant.ContactNumber,
			&restaurant.LicenceUrl,
			&restaurant.WebsiteUrl,
			&restaurant.CreatedAt,
			&restaurant.UpdatedAt,
			&restaurant.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row while listing restaurants: %v", err)
		}

		// Fetch location information for the restaurant
		locationQuery := fmt.Sprintf("SELECT * FROM %s WHERE establishment_id = $1", locationTableName)
		if err := p.db.QueryRow(ctx, locationQuery, restaurant.RestaurantId).Scan(
			&restaurant.Location.LocationId,
			&restaurant.Location.EstablishmentId,
			&restaurant.Location.Address,
			&restaurant.Location.Latitude,
			&restaurant.Location.Longitude,
			&restaurant.Location.Country,
			&restaurant.Location.City,
			&restaurant.Location.StateProvince,
			&restaurant.Location.CreatedAt,
			&restaurant.Location.UpdatedAt,
			&restaurant.Location.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to get location for restaurant: %v", err)
		}

		// Fetch images information for the attraction
		imagesQuery := fmt.Sprintf("SELECT * FROM %s WHERE establishment_id = $1", imageTableName)
		imageRows, err := p.db.Query(ctx, imagesQuery, restaurant.RestaurantId)
		if err != nil {
			return nil, fmt.Errorf("failed to get images for restaurant: %v", err)
		}

		// Iterate over the image rows and populate the Images slice for the restaurant
		defer imageRows.Close()
		for imageRows.Next() {
			var image entity.Image
			if err := imageRows.Scan(
				&image.ImageId,
				&image.EstablishmentId,
				&image.ImageUrl,
				&image.CreatedAt,
				&image.UpdatedAt,
				&image.DeletedAt,
			); err != nil {
				return nil, fmt.Errorf("failed to scan image row: %v", err)
			}
			restaurant.Images = append(restaurant.Images, &image)
		}
		if err := imageRows.Err(); err != nil {
			return nil, fmt.Errorf("error encountered while iterating over image rows: %v", err)
		}

		// Append the attraction to the restaurants slice
		restaurants = append(restaurants, &restaurant)
	}

	return restaurants, nil
}

// update a restaurant
func (p restaurantRepo) UpdateRestaurant(ctx context.Context, restaurant *entity.Restaurant) (*entity.Restaurant, error) {

	clauses := map[string]interface{}{
		"restaurant_name": restaurant.RestaurantName,
		"description":     restaurant.Description,
		"contact_number":  restaurant.ContactNumber,
		"opening_hours":   restaurant.OpeningHours,
		"licence_url":     restaurant.LicenceUrl,
		"website_url":     restaurant.WebsiteUrl,
	}

	sqlStr, args, err := p.db.Sq.Builder.Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("restaurant_id", restaurant.RestaurantId), p.db.Sq.Equal("deleted_at", nil)).
		ToSql()
	if err != nil {
		return restaurant, fmt.Errorf("failed to build SQL query for updating restaurant: %v", err)
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return restaurant, fmt.Errorf("failed to execute SQL query for updating restaurant: %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		return restaurant, fmt.Errorf("no rows affected while updating restaurant")
	}

	return restaurant, nil
}

// delete a restaurant
func (p restaurantRepo) DeleteRestaurant(ctx context.Context, restaurant_id string) error {
	// Build the SQL query
	sqlStr, args, err := p.db.Sq.Builder.Delete(p.tableName).
		Where(p.db.Sq.Equal("restaurant_id", restaurant_id)).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query for deleting restaurant: %v", err)
	}

	// Execute the SQL query
	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("failed to execute SQL query for deleting restaurant: %v", err)
	}

	// Check if any rows were affected
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no rows affected while deleting restaurant")
	}

	return nil
}
