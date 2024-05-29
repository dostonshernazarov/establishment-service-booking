package postgresql

import (
	"Booking/establishment-service-booking/internal/entity"
	"Booking/establishment-service-booking/internal/pkg/otlp"
	"Booking/establishment-service-booking/internal/pkg/postgres"
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
)

const (
	restaurantTableName = "restaurant_table"
	// locationTableName        = "location_table"
	// imageTableName           = "image_table"
	restaurantServiceName    = "restaurantService"
	restaurantSpanRepoPrefix = "attractionRepo"
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
	).From(p.tableName)
}

// create a new restaurant
func (p restaurantRepo) CreateRestaurant(ctx context.Context, restaurant *entity.Restaurant) (*entity.Restaurant, error) {

	ctx, span := otlp.Start(ctx, restaurantServiceName, restaurantSpanRepoPrefix+"Create")
	defer span.End()

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
		"category":         restaurant.Location.Category,
		"created_at":       restaurant.Location.CreatedAt,
		"updated_at":       restaurant.Location.UpdatedAt,
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
			"category":         image.Category,
			"created_at":       image.CreatedAt,
			"updated_at":       image.UpdatedAt,
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
	}
	query, args, err = p.db.Sq.Builder.Insert(p.tableName).SetMap(data).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for creating restaurant: %v", err)
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query for creating restaurant: %v", err)
	}

	return restaurant, nil
}

// get a restaurant
func (p restaurantRepo) GetRestaurant(ctx context.Context, restaurant_id string) (*entity.Restaurant, error) {

	ctx, span := otlp.Start(ctx, restaurantServiceName, restaurantSpanRepoPrefix+"Get")
	defer span.End()

	var restaurant entity.Restaurant

	// Build the query to select attraction details
	queryBuilder := p.RestaurantSelectQueryPrefix().Where(p.db.Sq.Equal("restaurant_id", restaurant_id)).Where(p.db.Sq.Equal("deleted_at", nil))

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
	); err != nil {
		return nil, fmt.Errorf("failed to get restaurant: %v", err)
	}

	// Fetch location information
	locationQuery := fmt.Sprintf("SELECT location_id, establishment_id, address, latitude, longitude, country, city, state_province, created_at, updated_at FROM %s WHERE establishment_id = $1", locationTableName)
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
	); err != nil {
		return nil, fmt.Errorf("failed to get location for location: %v", err)
	}

	// Fetch images information
	imagesQuery := fmt.Sprintf("SELECT image_id, establishment_id, image_url, created_at, updated_at FROM %s WHERE establishment_id = $1", imageTableName)
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
func (p restaurantRepo) ListRestaurants(ctx context.Context, offset, limit int64) ([]*entity.Restaurant, uint64, error) {

	ctx, span := otlp.Start(ctx, restaurantServiceName, restaurantSpanRepoPrefix+"List")
	defer span.End()

	var restaurants []*entity.Restaurant

	queryBuilder := p.RestaurantSelectQueryPrefix()

	if limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(limit)).Offset(uint64(offset)).Where(p.db.Sq.Equal("deleted_at", nil))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
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
		); err != nil {
			return nil, 0, err
		}

		// Fetch location information for the restaurant
		locationQuery := fmt.Sprintf("SELECT location_id, establishment_id, address, latitude, longitude, country, city, state_province, created_at, updated_at FROM %s WHERE establishment_id = $1", locationTableName)
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
		); err != nil {
			return nil, 0, err
		}

		// Fetch images information for the attraction
		imagesQuery := fmt.Sprintf("SELECT image_id, establishment_id, image_url, created_at, updated_at FROM %s WHERE establishment_id = $1", imageTableName)
		imageRows, err := p.db.Query(ctx, imagesQuery, restaurant.RestaurantId)
		if err != nil {
			return nil, 0, err
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
			); err != nil {
				return nil, 0, err
			}
			restaurant.Images = append(restaurant.Images, &image)
		}
		if err := imageRows.Err(); err != nil {
			return nil, 0, err
		}

		// Append the attraction to the restaurants slice
		restaurants = append(restaurants, &restaurant)
	}

	var overall uint64

	queryC := `SELECT COUNT(*) FROM restaurant_table WHERE deleted_at IS NULL`

	if err := p.db.QueryRow(ctx, queryC).Scan(&overall); err != nil {
		return nil, 0, err
	}

	return restaurants, overall, nil
}

// update a restaurant
func (p restaurantRepo) UpdateRestaurant(ctx context.Context, request *entity.Restaurant) (*entity.Restaurant, error) {

	ctx, span := otlp.Start(ctx, restaurantServiceName, restaurantSpanRepoPrefix+"Update")
	defer span.End()

	clauses := map[string]interface{}{
		"restaurant_name": request.RestaurantName,
		"description":     request.Description,
		"rating":          request.Rating,
		"opening_hours":   request.OpeningHours,
		"contact_number":  request.ContactNumber,
		"licence_url":     request.LicenceUrl,
		"website_url":     request.WebsiteUrl,
		"updated_at":      time.Now().Local(),
	}

	sqlStr, args, err := p.db.Sq.Builder.Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("restaurant_id", request.RestaurantId), p.db.Sq.Equal("deleted_at", nil)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for updating restaurant: %v", err)
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query for updating restaurant: %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		return nil, fmt.Errorf("no rows affected while updating restaurant")
	}

	clausesL := map[string]interface{}{
		"address":        request.Location.Address,
		"latitude":       request.Location.Latitude,
		"longitude":      request.Location.Longitude,
		"country":        request.Location.Country,
		"city":           request.Location.City,
		"state_province": request.Location.StateProvince,
		"updated_at":     time.Now().Local(),
	}

	sqlStrL, args, err := p.db.Sq.Builder.Update("location_table").
		SetMap(clausesL).
		Where(p.db.Sq.Equal("establishment_id", request.RestaurantId), p.db.Sq.Equal("deleted_at", nil)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for updating location of Restaurant: %v", err)
	}

	commandTagL, err := p.db.Exec(ctx, sqlStrL, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query for updating location of Restaurant: %v", err)
	}

	if commandTagL.RowsAffected() == 0 {
		return nil, fmt.Errorf("no rows affected while updating restaurant")
	}

	var restaurant entity.Restaurant

	// Build the query to select restaurant details
	queryBuilder := p.RestaurantSelectQueryPrefix().Where(p.db.Sq.Equal("restaurant_id", request.RestaurantId))

	// Get the SQL query and arguments from the query builder
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for getting restaurant: %v", err)
	}

	// Execute the query to fetch restaurant details
	if err := p.db.QueryRow(ctx, query, args...).Scan(
		&restaurant.RestaurantId,
		&restaurant.RestaurantName,
		&restaurant.OwnerId,
		&restaurant.Description,
		&restaurant.Rating,
		&restaurant.OpeningHours,
		&restaurant.ContactNumber,
		&restaurant.LicenceUrl,
		&restaurant.WebsiteUrl,
		&restaurant.CreatedAt,
		&restaurant.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get restaurant: %v", err)
	}

	// Fetch location information
	locationQuery := fmt.Sprintf("SELECT location_id, establishment_id, address, latitude, longitude, country, city, state_province, created_at, updated_at FROM %s WHERE establishment_id = $1", locationTableName)
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
	); err != nil {
		return nil, fmt.Errorf("failed to get location for restaurant: %v", err)
	}

	// Fetch images information
	imagesQuery := fmt.Sprintf("SELECT image_id, establishment_id, image_url, created_at, updated_at FROM %s WHERE establishment_id = $1", imageTableName)
	rows, err := p.db.Query(ctx, imagesQuery, request.RestaurantId)
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

// delete a restaurant softly
func (p restaurantRepo) DeleteRestaurant(ctx context.Context, restaurant_id string) error {

	ctx, span := otlp.Start(ctx, restaurantServiceName, restaurantSpanRepoPrefix+"Delete")
	defer span.End()

	// Build the SQL query
	sqlStr, args, err := p.db.Sq.Builder.Update(p.tableName).
		Set("deleted_at", time.Now().Local()).
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

// list restaurants by location
func (p restaurantRepo) ListRestaurantsByLocation(ctx context.Context, offset, limit uint64, country, city, state_province string) ([]*entity.Restaurant, int64, error) {

	ctx, span := otlp.Start(ctx, restaurantServiceName, restaurantSpanRepoPrefix+"ListL")
	defer span.End()

	countryStr := "%"+country+"%"
	cityStr := "%"+city+"%"
	stateStr := "%"+state_province+"%"

	queryL := fmt.Sprintf("SELECT establishment_id FROM location_table WHERE country LIKE '%s' and city LIKE '%s' and state_province LIKE '%s' and category = 'restaurant' and deleted_at IS NULL LIMIT $1 OFFSET $2", countryStr, cityStr, stateStr)
	rows, err := p.db.Query(ctx, queryL, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var restaurants []*entity.Restaurant

	for rows.Next() {

		var establishment_id string

		if err := rows.Scan(&establishment_id); err != nil {
			return nil, 0, err
		}

		var restaurant entity.Restaurant

		queryA := `SELECT restaurant_id, owner_id, restaurant_name, description, rating, opening_hours, contact_number, licence_url, website_url, created_at, updated_at FROM restaurant_table WHERE restaurant_id = $1`

		if err := p.db.QueryRow(ctx, queryA, establishment_id).Scan(
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
		); err != nil {
			return nil, 0, err
		}

		var location entity.Location

		queryLA := `SELECT location_id, establishment_id, address, latitude, longitude, country, city, state_province, created_at, updated_at FROM location_table WHERE establishment_id = $1`

		if err := p.db.QueryRow(ctx, queryLA, restaurant.RestaurantId).Scan(
			&location.LocationId,
			&location.EstablishmentId,
			&location.Address,
			&location.Latitude,
			&location.Longitude,
			&location.Country,
			&location.City,
			&location.StateProvince,
			&location.CreatedAt,
			&location.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		restaurant.Location = location

		queryI := `SELECT image_id, establishment_id, image_url, created_at, updated_at FROM image_table WHERE establishment_id = $1`

		rowsI, err := p.db.Query(ctx, queryI, restaurant.RestaurantId)
		if err != nil {
			return nil, 0, err
		}
		defer rowsI.Close()

		var images []*entity.Image

		for rowsI.Next() {
			var image entity.Image

			if err := rowsI.Scan(
				&image.ImageId,
				&image.EstablishmentId,
				&image.ImageUrl,
				&image.CreatedAt,
				&image.UpdatedAt,
			); err != nil {
				return nil, 0, err
			}

			images = append(images, &image)
		}

		restaurant.Images = images

		restaurants = append(restaurants, &restaurant)
	}

	var count int64

	queryC := fmt.Sprintf("SELECT COUNT(*) establishment_id FROM location_table WHERE country LIKE '%s' and city LIKE '%s' and state_province LIKE '%s' and category = 'restaurant' and deleted_at IS NULL", countryStr, cityStr, stateStr)

	if err := p.db.QueryRow(ctx, queryC).Scan(&count); err != nil {
		return restaurants, 0, err
	}

	return restaurants, count, nil
}

// find restaurants by name
func (p restaurantRepo) FindRestaurantsByName(ctx context.Context, name string) ([]*entity.Restaurant, uint64, error) {

	ctx, span := otlp.Start(ctx, restaurantServiceName, restaurantSpanRepoPrefix+"Find")
	defer span.End()

	var restaurants []*entity.Restaurant

	query := `SELECT
  restaurant_id,
  owner_id,
  restaurant_name,
  description,
  rating,
	opening_hours,
  contact_number,
  licence_url,
  website_url,
  created_at,
  updated_at
  FROM restaurant_table
  WHERE deleted_at IS NULL
  AND restaurant_name ILIKE '%' || $1 || '%'
  ORDER BY rating DESC`

	rows, err := p.db.Query(ctx, query, name)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Iterate over the rows to fetch each restaurant's details
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
		); err != nil {
			return nil, 0, err
		}

		// Fetch location information for the restaurant
		locationQuery := fmt.Sprintf("SELECT location_id, establishment_id, address, latitude, longitude, country, city, state_province, created_at, updated_at FROM %s WHERE establishment_id = $1", locationTableName)
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
		); err != nil {
			return nil, 0, err
		}

		// Fetch images information for the restaurant
		imagesQuery := fmt.Sprintf("SELECT image_id, establishment_id, image_url, created_at, updated_at FROM %s WHERE establishment_id = $1", imageTableName)
		imageRows, err := p.db.Query(ctx, imagesQuery, restaurant.RestaurantId)
		if err != nil {
			return nil, 0, err
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
			); err != nil {
				return nil, 0, err
			}
			restaurant.Images = append(restaurant.Images, &image)
		}
		if err := imageRows.Err(); err != nil {
			return nil, 0, err
		}

		// Append the restaurant to the restaurants slice
		restaurants = append(restaurants, &restaurant)
	}

	var overall uint64

	queryC := `SELECT COUNT(*) FROM restaurant_table WHERE restaurant_name ILIKE '%' || $1 || '%' and deleted_at IS NULL`

	if err := p.db.QueryRow(ctx, queryC, name).Scan(&overall); err != nil {
		return nil, 0, err
	}

	return restaurants, overall, nil
}
