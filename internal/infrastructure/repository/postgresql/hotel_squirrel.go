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
	hotelTableName = "hotel_table"
	// locationTableName        = "location_table"
	// imageTableName           = "image_table"
	hotelServiceName    = "hotelService"
	hotelSpanRepoPrefix = "hotelRepo"
)

type hotelRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewHotelRepo(db *postgres.PostgresDB) *hotelRepo {
	return &hotelRepo{
		tableName: hotelTableName,
		db:        db,
	}
}

func (p *hotelRepo) HotelSelectQueryPrefix() squirrel.SelectBuilder {
	return p.db.Sq.Builder.Select(
		"hotel_id",
		"owner_id",
		"hotel_name",
		"description",
		"rating",
		"contact_number",
		"licence_url",
		"website_url",
		"created_at",
		"updated_at",
	).From(p.tableName)
}

// create a new hotel
func (p hotelRepo) CreateHotel(ctx context.Context, hotel *entity.Hotel) (*entity.Hotel, error) {

	ctx, span := otlp.Start(ctx, hotelServiceName, hotelSpanRepoPrefix+"Create")
	defer span.End()

	// insert location info to location_table
	dataL := map[string]interface{}{
		"location_id":      hotel.Location.LocationId,
		"establishment_id": hotel.Location.EstablishmentId,
		"address":          hotel.Location.Address,
		"latitude":         hotel.Location.Latitude,
		"longitude":        hotel.Location.Longitude,
		"country":          hotel.Location.Country,
		"city":             hotel.Location.City,
		"state_province":   hotel.Location.StateProvince,
		"category":         hotel.Location.Category,
		"created_at":       hotel.Location.CreatedAt,
		"updated_at":       hotel.Location.UpdatedAt,
	}

	query, args, err := p.db.Sq.Builder.Insert(locationTableName).SetMap(dataL).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for creating hotel's location part: %v", err)
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query for creating hotel's location part: %v", err)
	}

	// insert images to image_table
	for _, image := range hotel.Images {
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
		"hotel_id":       hotel.HotelId,
		"owner_id":       hotel.OwnerId,
		"hotel_name":     hotel.HotelName,
		"description":    hotel.Description,
		"rating":         hotel.Rating,
		"contact_number": hotel.ContactNumber,
		"licence_url":    hotel.LicenceUrl,
		"website_url":    hotel.WebsiteUrl,
		"created_at":     hotel.CreatedAt,
		"updated_at":     hotel.UpdatedAt,
	}
	query, args, err = p.db.Sq.Builder.Insert(p.tableName).SetMap(data).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for creating hotel: %v", err)
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query for creating hotel: %v", err)
	}

	return hotel, nil
}

// get a restaurant
func (p hotelRepo) GetHotel(ctx context.Context, hotel_id string) (*entity.Hotel, error) {

	ctx, span := otlp.Start(ctx, hotelServiceName, hotelSpanRepoPrefix+"Get")
	defer span.End()

	var hotel entity.Hotel

	// Build the query to select attraction details
	queryBuilder := p.HotelSelectQueryPrefix().Where(p.db.Sq.Equal("hotel_id", hotel_id)).Where(p.db.Sq.Equal("deleted_at", nil))

	// Get the SQL query and arguments from the query builder
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for getting hotel: %v", err)
	}

	// Execute the query to fetch hotel details
	if err := p.db.QueryRow(ctx, query, args...).Scan(
		&hotel.HotelId,
		&hotel.OwnerId,
		&hotel.HotelName,
		&hotel.Description,
		&hotel.Rating,
		&hotel.ContactNumber,
		&hotel.LicenceUrl,
		&hotel.WebsiteUrl,
		&hotel.CreatedAt,
		&hotel.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get hotel: %v", err)
	}

	// Fetch location information
	locationQuery := fmt.Sprintf("SELECT location_id, establishment_id, address, latitude, longitude, country, city, state_province, created_at, updated_at FROM %s WHERE establishment_id = $1", locationTableName)
	if err := p.db.QueryRow(ctx, locationQuery, hotel.HotelId).Scan(
		&hotel.Location.LocationId,
		&hotel.Location.EstablishmentId,
		&hotel.Location.Address,
		&hotel.Location.Latitude,
		&hotel.Location.Longitude,
		&hotel.Location.Country,
		&hotel.Location.City,
		&hotel.Location.StateProvince,
		&hotel.Location.CreatedAt,
		&hotel.Location.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get location for hotel: %v", err)
	}

	// Fetch images information
	imagesQuery := fmt.Sprintf("SELECT image_id, establishment_id, image_url, created_at, updated_at FROM %s WHERE establishment_id = $1", imageTableName)
	rows, err := p.db.Query(ctx, imagesQuery, hotel_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get images for hotel: %v", err)
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
		hotel.Images = append(hotel.Images, &image)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered while iterating over image rows: %v", err)
	}

	return &hotel, nil
}

// get a list of hotels
func (p hotelRepo) ListHotels(ctx context.Context, offset, limit int64) ([]*entity.Hotel, uint64, error) {

	ctx, span := otlp.Start(ctx, hotelServiceName, hotelSpanRepoPrefix+"List")
	defer span.End()

	var hotels []*entity.Hotel

	queryBuilder := p.HotelSelectQueryPrefix()

	if limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(limit)).Offset(uint64(offset)).Where(p.db.Sq.Equal("deleted_at", nil)).OrderBy("rating DESC")
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

	// Iterate over the rows to fetch each hotel's details
	for rows.Next() {
		var hotel entity.Hotel
		if err := rows.Scan(
			&hotel.HotelId,
			&hotel.OwnerId,
			&hotel.HotelName,
			&hotel.Description,
			&hotel.Rating,
			&hotel.ContactNumber,
			&hotel.LicenceUrl,
			&hotel.WebsiteUrl,
			&hotel.CreatedAt,
			&hotel.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		// Fetch location information for the hotel
		locationQuery := fmt.Sprintf("SELECT location_id, establishment_id, address, latitude, longitude, country, city, state_province, created_at, updated_at FROM %s WHERE establishment_id = $1", locationTableName)

		println()
		println(hotel.HotelId)
		println()

		if err := p.db.QueryRow(ctx, locationQuery, hotel.HotelId).Scan(
			&hotel.Location.LocationId,
			&hotel.Location.EstablishmentId,
			&hotel.Location.Address,
			&hotel.Location.Latitude,
			&hotel.Location.Longitude,
			&hotel.Location.Country,
			&hotel.Location.City,
			&hotel.Location.StateProvince,
			&hotel.Location.CreatedAt,
			&hotel.Location.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		// Fetch images information for the attraction
		imagesQuery := fmt.Sprintf("SELECT image_id, establishment_id, image_url, created_at, updated_at FROM %s WHERE establishment_id = $1", imageTableName)
		imageRows, err := p.db.Query(ctx, imagesQuery, hotel.HotelId)
		if err != nil {
			return nil, 0, err
		}

		// Iterate over the image rows and populate the Images slice for the hotel
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
			hotel.Images = append(hotel.Images, &image)
		}
		if err := imageRows.Err(); err != nil {
			return nil, 0, err
		}

		// Append the attraction to the hotels slice
		hotels = append(hotels, &hotel)
	}

	var overall uint64

	queryC := `SELECT COUNT(*) FROM hotel_table WHERE deleted_at IS NULL`

	if err := p.db.QueryRow(ctx, queryC).Scan(&overall); err != nil {
		return nil, 0, err
	}

	return hotels, overall, nil
}

// update a hotel
func (p hotelRepo) UpdateHotel(ctx context.Context, request *entity.Hotel) (*entity.Hotel, error) {

	ctx, span := otlp.Start(ctx, hotelServiceName, hotelSpanRepoPrefix+"Update")
	defer span.End()

	clauses := map[string]interface{}{
		"hotel_name":     request.HotelName,
		"description":    request.Description,
		"rating":         request.Rating,
		"contact_number": request.ContactNumber,
		"licence_url":    request.LicenceUrl,
		"website_url":    request.WebsiteUrl,
		"updated_at":     time.Now().Local(),
	}

	sqlStr, args, err := p.db.Sq.Builder.Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("hotel_id", request.HotelId), p.db.Sq.Equal("deleted_at", nil)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for updating hotel: %v", err)
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query for updating hotel: %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		return nil, fmt.Errorf("no rows affected while updating hotel")
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
		Where(p.db.Sq.Equal("establishment_id", request.HotelId), p.db.Sq.Equal("deleted_at", nil)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for updating location: %v", err)
	}

	commandTagL, err := p.db.Exec(ctx, sqlStrL, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query for updating location: %v", err)
	}

	if commandTagL.RowsAffected() == 0 {
		return nil, fmt.Errorf("no rows affected while updating hotel")
	}

	var hotel entity.Hotel

	// Build the query to select hotel details
	queryBuilder := p.HotelSelectQueryPrefix().Where(p.db.Sq.Equal("hotel_id", request.HotelId))

	// Get the SQL query and arguments from the query builder
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for getting hotel: %v", err)
	}

	// Execute the query to fetch hotel details
	if err := p.db.QueryRow(ctx, query, args...).Scan(
		&hotel.HotelId,
		&hotel.HotelName,
		&hotel.OwnerId,
		&hotel.Description,
		&hotel.Rating,
		&hotel.ContactNumber,
		&hotel.LicenceUrl,
		&hotel.WebsiteUrl,
		&hotel.CreatedAt,
		&hotel.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get hotel: %v", err)
	}

	// Fetch location information
	locationQuery := fmt.Sprintf("SELECT location_id, establishment_id, address, latitude, longitude, country, city, state_province, created_at, updated_at FROM %s WHERE establishment_id = $1", locationTableName)
	if err := p.db.QueryRow(ctx, locationQuery, hotel.HotelId).Scan(
		&hotel.Location.LocationId,
		&hotel.Location.EstablishmentId,
		&hotel.Location.Address,
		&hotel.Location.Latitude,
		&hotel.Location.Longitude,
		&hotel.Location.Country,
		&hotel.Location.City,
		&hotel.Location.StateProvince,
		&hotel.Location.CreatedAt,
		&hotel.Location.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get location for hotel: %v", err)
	}

	// Fetch images information
	imagesQuery := fmt.Sprintf("SELECT image_id, establishment_id, image_url, created_at, updated_at FROM %s WHERE establishment_id = $1", imageTableName)
	rows, err := p.db.Query(ctx, imagesQuery, request.HotelId)
	if err != nil {
		return nil, fmt.Errorf("failed to get images for hotel: %v", err)
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
		hotel.Images = append(hotel.Images, &image)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered while iterating over image rows: %v", err)
	}

	return &hotel, nil
}

// delete a hotel softly
func (p hotelRepo) DeleteHotel(ctx context.Context, hotel_id string) error {

	ctx, span := otlp.Start(ctx, hotelServiceName, hotelSpanRepoPrefix+"Delete")
	defer span.End()

	// Build the SQL query
	sqlStr, args, err := p.db.Sq.Builder.Update(p.tableName).
		Set("deleted_at", time.Now().Local()).
		Where(p.db.Sq.Equal("hotel_id", hotel_id)).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query for deleting hotel: %v", err)
	}

	// Execute the SQL query
	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("failed to execute SQL query for deleting hotel: %v", err)
	}

	// Check if any rows were affected
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no rows affected while deleting hotel")
	}

	return nil
}

// list hotels by location
func (p hotelRepo) ListHotelsByLocation(ctx context.Context, offset, limit uint64, country, city, state_province string) ([]*entity.Hotel, int64, error) {

	ctx, span := otlp.Start(ctx, hotelServiceName, hotelSpanRepoPrefix+"ListL")
	defer span.End()

	countryStr := "%"+country+"%"
	cityStr := "%"+city+"%"
	stateStr := "%"+state_province+"%"

	// println("\n\n chech \n")
	queryL := fmt.Sprintf("SELECT establishment_id FROM location_table WHERE country LIKE '%s' and city LIKE '%s' and state_province LIKE '%s' and category = 'hotel' and deleted_at IS NULL LIMIT $1 OFFSET $2", countryStr, cityStr, stateStr)
	rows, err := p.db.Query(ctx, queryL, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var hotels []*entity.Hotel

	for rows.Next() {

		var establishment_id string

		if err := rows.Scan(&establishment_id); err != nil {
			return nil, 0, err
		}

		var hotel entity.Hotel

		queryA := `SELECT hotel_id, owner_id, hotel_name, description, rating, contact_number, licence_url, website_url, created_at, updated_at FROM hotel_table WHERE hotel_id = $1`

		if err := p.db.QueryRow(ctx, queryA, establishment_id).Scan(
			&hotel.HotelId,
			&hotel.OwnerId,
			&hotel.HotelName,
			&hotel.Description,
			&hotel.Rating,
			&hotel.ContactNumber,
			&hotel.LicenceUrl,
			&hotel.WebsiteUrl,
			&hotel.CreatedAt,
			&hotel.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		var location entity.Location

		queryLA := `SELECT location_id, establishment_id, address, latitude, longitude, country, city, state_province, created_at, updated_at FROM location_table WHERE establishment_id = $1`

		if err := p.db.QueryRow(ctx, queryLA, hotel.HotelId).Scan(
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

		hotel.Location = location

		queryI := `SELECT image_id, establishment_id, image_url, created_at, updated_at FROM image_table WHERE establishment_id = $1`

		rowsI, err := p.db.Query(ctx, queryI, hotel.HotelId)
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

		hotel.Images = images

		hotels = append(hotels, &hotel)
	}

	var count int64

	queryC := fmt.Sprintf("SELECT COUNT(*) establishment_id FROM location_table WHERE country LIKE '%s' and city LIKE '%s' and state_province LIKE '%s' and category = 'hotel' and deleted_at IS NULL", countryStr, cityStr, stateStr)

	if err := p.db.QueryRow(ctx, queryC).Scan(&count); err != nil {
		return hotels, 0, err
	}

	return hotels, count, nil
}

// find hotels by name
func (p hotelRepo) FindHotelsByName(ctx context.Context, name string) ([]*entity.Hotel, uint64, error) {

	ctx, span := otlp.Start(ctx, hotelServiceName, hotelSpanRepoPrefix+"Find")
	defer span.End()

	var hotels []*entity.Hotel

	query := `SELECT
  hotel_id,
  owner_id,
  hotel_name,
  description,
  rating,
  contact_number,
  licence_url,
  website_url,
  created_at,
  updated_at
  FROM hotel_table
  WHERE deleted_at IS NULL
  AND hotel_name ILIKE '%' || $1 || '%'
  ORDER BY rating DESC`

	rows, err := p.db.Query(ctx, query, name)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Iterate over the rows to fetch each hotel's details
	for rows.Next() {
		var hotel entity.Hotel
		if err := rows.Scan(
			&hotel.HotelId,
			&hotel.OwnerId,
			&hotel.HotelName,
			&hotel.Description,
			&hotel.Rating,
			&hotel.ContactNumber,
			&hotel.LicenceUrl,
			&hotel.WebsiteUrl,
			&hotel.CreatedAt,
			&hotel.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		// Fetch location information for the hotel
		locationQuery := fmt.Sprintf("SELECT location_id, establishment_id, address, latitude, longitude, country, city, state_province, created_at, updated_at FROM %s WHERE establishment_id = $1", locationTableName)
		if err := p.db.QueryRow(ctx, locationQuery, hotel.HotelId).Scan(
			&hotel.Location.LocationId,
			&hotel.Location.EstablishmentId,
			&hotel.Location.Address,
			&hotel.Location.Latitude,
			&hotel.Location.Longitude,
			&hotel.Location.Country,
			&hotel.Location.City,
			&hotel.Location.StateProvince,
			&hotel.Location.CreatedAt,
			&hotel.Location.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		// Fetch images information for the hotel
		imagesQuery := fmt.Sprintf("SELECT image_id, establishment_id, image_url, created_at, updated_at FROM %s WHERE establishment_id = $1", imageTableName)
		imageRows, err := p.db.Query(ctx, imagesQuery, hotel.HotelId)
		if err != nil {
			return nil, 0, err
		}

		// Iterate over the image rows and populate the Images slice for the hotel
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
			hotel.Images = append(hotel.Images, &image)
		}
		if err := imageRows.Err(); err != nil {
			return nil, 0, err
		}

		// Append the hotel to the hotels slice
		hotels = append(hotels, &hotel)
	}

	var overall uint64

	queryC := `SELECT COUNT(*) FROM hotel_table WHERE hotel_name ILIKE '%' || $1 || '%' and deleted_at IS NULL`

	if err := p.db.QueryRow(ctx, queryC, name).Scan(&overall); err != nil {
		return nil, 0, err
	}

	return hotels, overall, nil
}
