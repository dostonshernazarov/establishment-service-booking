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
	attractionTableName      = "attraction_table" //table for storing general info of attraction
	locationTableName        = "location_table"   // table for storing location info
	imageTableName           = "image_table"      // table for storing multiple images of establishment
	attractionServiceName    = "attractionService"
	attractionSpanRepoPrefix = "attractionRepo"
)

type attractionRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewAttractionRepo(db *postgres.PostgresDB) *attractionRepo {
	return &attractionRepo{
		tableName: attractionTableName,
		db:        db,
	}
}

func (p *attractionRepo) AttractionSelectQueryPrefix() squirrel.SelectBuilder {
	return p.db.Sq.Builder.Select(
		"attraction_id",
		"attraction_name",
		"owner_id",
		"description",
		"rating",
		"contact_number",
		"licence_url",
		"website_url",
		"created_at",
		"updated_at",
	).From(p.tableName)
}

// create a new attraction
func (p attractionRepo) CreateAttraction(ctx context.Context, attraction *entity.Attraction) (*entity.Attraction, error) {

	ctx, span := otlp.Start(ctx, attractionServiceName, attractionSpanRepoPrefix+"Create")
	defer span.End()

	// insert location info to location_table
	dataL := map[string]interface{}{
		"location_id":      attraction.Location.LocationId,
		"establishment_id": attraction.Location.EstablishmentId,
		"address":          attraction.Location.Address,
		"latitude":         attraction.Location.Latitude,
		"longitude":        attraction.Location.Longitude,
		"country":          attraction.Location.Country,
		"city":             attraction.Location.City,
		"state_province":   attraction.Location.StateProvince,
		"category":         attraction.Location.Category,
		"created_at":       attraction.Location.CreatedAt,
		"updated_at":       attraction.Location.UpdatedAt,
	}

	query, args, err := p.db.Sq.Builder.Insert(locationTableName).SetMap(dataL).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for creating attraction' location part: %v", err)
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query for creating attraction's location: %v", err)
	}

	// insert images to image_table
	for _, image := range attraction.Images {
		dataI := map[string]interface{}{
			"image_id":         image.ImageId,
			"establishment_id": attraction.AttractionId,
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
		"attraction_id":   attraction.AttractionId,
		"attraction_name": attraction.AttractionName,
		"owner_id":        attraction.OwnerId,
		"description":     attraction.Description,
		"rating":          attraction.Rating,
		"contact_number":  attraction.ContactNumber,
		"licence_url":     attraction.LicenceUrl,
		"website_url":     attraction.WebsiteUrl,
		"created_at":      attraction.CreatedAt,
		"updated_at":      attraction.UpdatedAt,
	}
	query, args, err = p.db.Sq.Builder.Insert(p.tableName).SetMap(data).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for creating attraction: %v", err)
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query for creating attraction: %v", err)
	}

	return attraction, nil
}

// get an attraction
func (p attractionRepo) GetAttraction(ctx context.Context, attraction_id string) (*entity.Attraction, error) {

	ctx, span := otlp.Start(ctx, attractionServiceName, attractionSpanRepoPrefix+"Get")
	defer span.End()

	var attraction entity.Attraction

	// Build the query to select attraction details
	queryBuilder := p.AttractionSelectQueryPrefix().Where(p.db.Sq.Equal("attraction_id", attraction_id)).Where(p.db.Sq.Equal("deleted_at", nil))

	// Get the SQL query and arguments from the query builder
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for getting attraction: %v", err)
	}

	// Execute the query to fetch attraction details
	if err := p.db.QueryRow(ctx, query, args...).Scan(
		&attraction.AttractionId,
		&attraction.AttractionName,
		&attraction.OwnerId,
		&attraction.Description,
		&attraction.Rating,
		&attraction.ContactNumber,
		&attraction.LicenceUrl,
		&attraction.WebsiteUrl,
		&attraction.CreatedAt,
		&attraction.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get attraction: %v", err)
	}

	// Fetch location information
	locationQuery := fmt.Sprintf("SELECT location_id, establishment_id, address, latitude, longitude, country, city, state_province, created_at, updated_at FROM %s WHERE establishment_id = $1", locationTableName)
	if err := p.db.QueryRow(ctx, locationQuery, attraction.AttractionId).Scan(
		&attraction.Location.LocationId,
		&attraction.Location.EstablishmentId,
		&attraction.Location.Address,
		&attraction.Location.Latitude,
		&attraction.Location.Longitude,
		&attraction.Location.Country,
		&attraction.Location.City,
		&attraction.Location.StateProvince,
		&attraction.Location.CreatedAt,
		&attraction.Location.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get location for attraction: %v", err)
	}

	// Fetch images information
	imagesQuery := fmt.Sprintf("SELECT image_id, establishment_id, image_url, created_at, updated_at FROM %s WHERE establishment_id = $1", imageTableName)
	rows, err := p.db.Query(ctx, imagesQuery, attraction_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get images for attraction: %v", err)
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
		attraction.Images = append(attraction.Images, &image)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered while iterating over image rows: %v", err)
	}

	return &attraction, nil
}

// get a list of attractions
func (p attractionRepo) ListAttractions(ctx context.Context, offset, limit int64) ([]*entity.Attraction, uint64, error) {

	ctx, span := otlp.Start(ctx, attractionServiceName, attractionSpanRepoPrefix+"List")
	defer span.End()

	var attractions []*entity.Attraction

	queryBuilder := p.AttractionSelectQueryPrefix()

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

	// Iterate over the rows to fetch each attraction's details
	for rows.Next() {
		var attraction entity.Attraction
		if err := rows.Scan(
			&attraction.AttractionId,
			&attraction.AttractionName,
			&attraction.OwnerId,
			&attraction.Description,
			&attraction.Rating,
			&attraction.ContactNumber,
			&attraction.LicenceUrl,
			&attraction.WebsiteUrl,
			&attraction.CreatedAt,
			&attraction.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		// Fetch location information for the attraction
		locationQuery := fmt.Sprintf("SELECT location_id, establishment_id, address, latitude, longitude, country, city, state_province, created_at, updated_at FROM %s WHERE establishment_id = $1", locationTableName)
		if err := p.db.QueryRow(ctx, locationQuery, attraction.AttractionId).Scan(
			&attraction.Location.LocationId,
			&attraction.Location.EstablishmentId,
			&attraction.Location.Address,
			&attraction.Location.Latitude,
			&attraction.Location.Longitude,
			&attraction.Location.Country,
			&attraction.Location.City,
			&attraction.Location.StateProvince,
			&attraction.Location.CreatedAt,
			&attraction.Location.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		// Fetch images information for the attraction
		imagesQuery := fmt.Sprintf("SELECT image_id, establishment_id, image_url, created_at, updated_at FROM %s WHERE establishment_id = $1", imageTableName)
		imageRows, err := p.db.Query(ctx, imagesQuery, attraction.AttractionId)
		if err != nil {
			return nil, 0, err
		}

		// Iterate over the image rows and populate the Images slice for the attraction
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
			attraction.Images = append(attraction.Images, &image)
		}
		if err := imageRows.Err(); err != nil {
			return nil, 0, err
		}

		// Append the attraction to the attractions slice
		attractions = append(attractions, &attraction)
	}

	var overall uint64

	queryC := `SELECT COUNT(*) FROM attraction_table WHERE deleted_at IS NULL`

	if err := p.db.QueryRow(ctx, queryC).Scan(&overall); err != nil {
		return nil, 0, err
	}

	return attractions, overall, nil
}

// update an attraction
func (p attractionRepo) UpdateAttraction(ctx context.Context, request *entity.Attraction) (*entity.Attraction, error) {

	ctx, span := otlp.Start(ctx, attractionServiceName, attractionSpanRepoPrefix+"Update")
	defer span.End()

	clauses := map[string]interface{}{
		"attraction_name": request.AttractionName,
		"description":     request.Description,
		"rating":          request.Rating,
		"contact_number":  request.ContactNumber,
		"licence_url":     request.LicenceUrl,
		"website_url":     request.WebsiteUrl,
		"updated_at":      time.Now().Local(),
	}

	sqlStr, args, err := p.db.Sq.Builder.Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("attraction_id", request.AttractionId), p.db.Sq.Equal("deleted_at", nil)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for updating attracation: %v", err)
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query for updating attraction: %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		return nil, fmt.Errorf("no rows affected while updating attraction")
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
		Where(p.db.Sq.Equal("establishment_id", request.AttractionId), p.db.Sq.Equal("deleted_at", nil)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for updating location: %v", err)
	}

	commandTagL, err := p.db.Exec(ctx, sqlStrL, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query for updating attraction: %v", err)
	}

	if commandTagL.RowsAffected() == 0 {
		return nil, fmt.Errorf("no rows affected while updating attraction")
	}

	var attraction entity.Attraction

	// Build the query to select attraction details
	queryBuilder := p.AttractionSelectQueryPrefix().Where(p.db.Sq.Equal("attraction_id", request.AttractionId))

	// Get the SQL query and arguments from the query builder
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for getting attraction: %v", err)
	}

	// Execute the query to fetch attraction details
	if err := p.db.QueryRow(ctx, query, args...).Scan(
		&attraction.AttractionId,
		&attraction.AttractionName,
		&attraction.OwnerId,
		&attraction.Description,
		&attraction.Rating,
		&attraction.ContactNumber,
		&attraction.LicenceUrl,
		&attraction.WebsiteUrl,
		&attraction.CreatedAt,
		&attraction.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get attraction: %v", err)
	}

	// Fetch location information
	locationQuery := fmt.Sprintf("SELECT location_id, establishment_id, address, latitude, longitude, country, city, state_province, created_at, updated_at FROM %s WHERE establishment_id = $1", locationTableName)
	if err := p.db.QueryRow(ctx, locationQuery, attraction.AttractionId).Scan(
		&attraction.Location.LocationId,
		&attraction.Location.EstablishmentId,
		&attraction.Location.Address,
		&attraction.Location.Latitude,
		&attraction.Location.Longitude,
		&attraction.Location.Country,
		&attraction.Location.City,
		&attraction.Location.StateProvince,
		&attraction.Location.CreatedAt,
		&attraction.Location.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get location for attraction: %v", err)
	}

	// Fetch images information
	imagesQuery := fmt.Sprintf("SELECT image_id, establishment_id, image_url, created_at, updated_at FROM %s WHERE establishment_id = $1", imageTableName)
	rows, err := p.db.Query(ctx, imagesQuery, request.AttractionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get images for attraction: %v", err)
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
		attraction.Images = append(attraction.Images, &image)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered while iterating over image rows: %v", err)
	}

	return &attraction, nil
}

// delete an attraction softly
func (p attractionRepo) DeleteAttraction(ctx context.Context, attraction_id string) error {

	ctx, span := otlp.Start(ctx, attractionServiceName, attractionSpanRepoPrefix+"Delete")
	defer span.End()

	// Build the SQL query
	sqlStr, args, err := p.db.Sq.Builder.Update(p.tableName).
		Set("deleted_at", time.Now().Local()).
		Where(p.db.Sq.Equal("attraction_id", attraction_id)).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query for deleting attraction: %v", err)
	}

	// Execute the SQL query
	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("failed to execute SQL query for deleting attraction: %v", err)
	}

	// Check if any rows were affected
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no rows affected while deleting attraction")
	}

	return nil
}

// list attractions by location
func (p attractionRepo) ListAttractionsByLocation(ctx context.Context, offset, limit uint64, country, city, state_province string) ([]*entity.Attraction, int64, error) {

	ctx, span := otlp.Start(ctx, attractionServiceName, attractionSpanRepoPrefix+"ListL")
	defer span.End()

	countryStr := "%"+country+"%"
	cityStr := "%"+city+"%"
	stateStr := "%"+state_province+"%"

	queryL := fmt.Sprintf("SELECT establishment_id FROM location_table WHERE country LIKE '%s' and city LIKE '%s' and state_province LIKE '%s' and category = 'attraction' and deleted_at IS NULL LIMIT $1 OFFSET $2",countryStr, cityStr, stateStr)
	rows, err := p.db.Query(ctx, queryL, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var attractions []*entity.Attraction

	for rows.Next() {

		var establishment_id string

		if err := rows.Scan(&establishment_id); err != nil {
			return nil, 0, err
		}

		var attraction entity.Attraction

		queryA := `SELECT attraction_id, owner_id, attraction_name, description, rating, contact_number, licence_url, website_url, created_at, updated_at FROM attraction_table WHERE attraction_id = $1`

		if err := p.db.QueryRow(ctx, queryA, establishment_id).Scan(
			&attraction.AttractionId,
			&attraction.OwnerId,
			&attraction.AttractionName,
			&attraction.Description,
			&attraction.Rating,
			&attraction.ContactNumber,
			&attraction.LicenceUrl,
			&attraction.WebsiteUrl,
			&attraction.CreatedAt,
			&attraction.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		var location entity.Location

		queryLA := `SELECT location_id, establishment_id, address, latitude, longitude, country, city, state_province, created_at, updated_at FROM location_table WHERE establishment_id = $1`

		if err := p.db.QueryRow(ctx, queryLA, attraction.AttractionId).Scan(
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

		attraction.Location = location

		queryI := `SELECT image_id, establishment_id, image_url, created_at, updated_at FROM image_table WHERE establishment_id = $1`

		rowsI, err := p.db.Query(ctx, queryI, attraction.AttractionId)
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

		attraction.Images = images

		attractions = append(attractions, &attraction)
	}

	var count int64

	queryC := `SELECT COUNT(*) establishment_id FROM location_table where category = 'attraction'`

	if err := p.db.QueryRow(ctx, queryC).Scan(&count); err != nil {
		return attractions, 0, err
	}

	return attractions, count, nil
}

// find attractions by name
func (p attractionRepo) FindAttractionsByName(ctx context.Context, name string) ([]*entity.Attraction, uint64, error) {

	ctx, span := otlp.Start(ctx, attractionServiceName, attractionSpanRepoPrefix+"Find")
	defer span.End()

	var attractions []*entity.Attraction

	query := `SELECT
  attraction_id,
  attraction_name,
  owner_id,
  description,
  rating,
  contact_number,
  licence_url,
  website_url,
  created_at,
  updated_at
  FROM attraction_table
  WHERE deleted_at IS NULL
  AND attraction_name ILIKE '%' || $1 || '%'
  ORDER BY rating DESC`

	rows, err := p.db.Query(ctx, query, name)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Iterate over the rows to fetch each attraction's details
	for rows.Next() {
		var attraction entity.Attraction
		if err := rows.Scan(
			&attraction.AttractionId,
			&attraction.AttractionName,
			&attraction.OwnerId,
			&attraction.Description,
			&attraction.Rating,
			&attraction.ContactNumber,
			&attraction.LicenceUrl,
			&attraction.WebsiteUrl,
			&attraction.CreatedAt,
			&attraction.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		// Fetch location information for the attraction
		locationQuery := fmt.Sprintf("SELECT location_id, establishment_id, address, latitude, longitude, country, city, state_province, created_at, updated_at FROM %s WHERE establishment_id = $1", locationTableName)
		if err := p.db.QueryRow(ctx, locationQuery, attraction.AttractionId).Scan(
			&attraction.Location.LocationId,
			&attraction.Location.EstablishmentId,
			&attraction.Location.Address,
			&attraction.Location.Latitude,
			&attraction.Location.Longitude,
			&attraction.Location.Country,
			&attraction.Location.City,
			&attraction.Location.StateProvince,
			&attraction.Location.CreatedAt,
			&attraction.Location.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		// Fetch images information for the attraction
		imagesQuery := fmt.Sprintf("SELECT image_id, establishment_id, image_url, created_at, updated_at FROM %s WHERE establishment_id = $1", imageTableName)
		imageRows, err := p.db.Query(ctx, imagesQuery, attraction.AttractionId)
		if err != nil {
			return nil, 0, err
		}

		// Iterate over the image rows and populate the Images slice for the attraction
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
			attraction.Images = append(attraction.Images, &image)
		}
		if err := imageRows.Err(); err != nil {
			return nil, 0, err
		}

		// Append the attraction to the attractions slice
		attractions = append(attractions, &attraction)
	}

	var overall uint64

	queryC := `SELECT COUNT(*) FROM attraction_table WHERE attraction_name ILIKE '%' || $1 || '%' and deleted_at IS NULL`

	if err := p.db.QueryRow(ctx, queryC, name).Scan(&overall); err != nil {
		return nil, 0, err
	}

	return attractions, overall, nil
}
