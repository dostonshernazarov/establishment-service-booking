package postgresql

import (
	"Booking/establishment-service-booking/internal/entity"
	"Booking/establishment-service-booking/internal/pkg/postgres"
	"context"
	"fmt"

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
		"deleted_at",
	).From(p.tableName)
}

// create a new attraction
func (p attractionRepo) CreateAttraction(ctx context.Context, attraction *entity.Attraction) (*entity.Attraction, error) {

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
		"created_at":       attraction.Location.CreatedAt,
		"updated_at":       attraction.Location.UpdatedAt,
		"deleted_at":       attraction.Location.DeletedAt,
	}

	query, args, err := p.db.Sq.Builder.Insert(locationTableName).SetMap(dataL).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for creating attraction' location part: %v", err)
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query for creating attraction: %v", err)
	}

	// insert images to image_table
	for _, image := range attraction.Images {
		dataI := map[string]interface{}{
			"image_id":         image.ImageId,
			"establishment_id": attraction.AttractionId,
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
		"deleted_at":      attraction.DeletedAt,
	}
	query, args, err = p.db.Sq.Builder.Insert(p.tableName).SetMap(data).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for creating attraction: %v", err)
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query for creating attraction: %v", err)
	}

	return nil, nil
}

// get an attraction
func (p attractionRepo) GetAttraction(ctx context.Context, attraction_id string) (*entity.Attraction, error) {
	var attraction entity.Attraction

	// Build the query to select attraction details
	queryBuilder := p.AttractionSelectQueryPrefix().Where(p.db.Sq.Equal("attraction_id", attraction_id))

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
		&attraction.DeletedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get attraction: %v", err)
	}

	// Fetch location information
	locationQuery := fmt.Sprintf("SELECT * FROM %s WHERE establishment_id = $1", locationTableName)
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
		&attraction.Location.DeletedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get location for attraction: %v", err)
	}

	// Fetch images information
	imagesQuery := fmt.Sprintf("SELECT * FROM %s WHERE establishment_id = $1", imageTableName)
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
			&image.DeletedAt,
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
func (p attractionRepo) ListAttractions(ctx context.Context, page, limit int64) ([]*entity.Attraction, error) {
	var attractions []*entity.Attraction

	// calculate offset
	offset := (page - 1) * limit

	queryBuilder := p.AttractionSelectQueryPrefix()

	if limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(limit)).Offset(uint64(offset)).Where(p.db.Sq.Equal("deleted_at", nil))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query for listing attractions: %v", err)
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query for listing attractions: %v", err)
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
			&attraction.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row while listing attractions: %v", err)
		}

		// Fetch location information for the attraction
		locationQuery := fmt.Sprintf("SELECT * FROM %s WHERE establishment_id = $1", locationTableName)
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
			&attraction.Location.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to get location for attraction: %v", err)
		}

		// Fetch images information for the attraction
		imagesQuery := fmt.Sprintf("SELECT * FROM %s WHERE establishment_id = $1", imageTableName)
		imageRows, err := p.db.Query(ctx, imagesQuery, attraction.AttractionId)
		if err != nil {
			return nil, fmt.Errorf("failed to get images for attraction: %v", err)
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
				&image.DeletedAt,
			); err != nil {
				return nil, fmt.Errorf("failed to scan image row: %v", err)
			}
			attraction.Images = append(attraction.Images, &image)
		}
		if err := imageRows.Err(); err != nil {
			return nil, fmt.Errorf("error encountered while iterating over image rows: %v", err)
		}

		// Append the attraction to the attractions slice
		attractions = append(attractions, &attraction)
	}

	return attractions, nil
}

// update an attraction
func (p attractionRepo) UpdateAttraction(ctx context.Context, attraction *entity.Attraction) (*entity.Attraction, error) {

	// println("\n\n ", attraction)
	clauses := map[string]interface{}{
		"attraction_name": attraction.AttractionName,
		"description":     attraction.Description,
		"contact_number":  attraction.ContactNumber,
		"licence_url":     attraction.LicenceUrl,
		"website_url":     attraction.WebsiteUrl,
	}

	sqlStr, args, err := p.db.Sq.Builder.Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("attraction_id", attraction.AttractionId), p.db.Sq.Equal("deleted_at", nil)).
		ToSql()
	if err != nil {
		return attraction, fmt.Errorf("failed to build SQL query for updating attracation: %v", err)
	}

	

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return attraction, fmt.Errorf("failed to execute SQL query for updating attraction: %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		return attraction, fmt.Errorf("no rows affected while updating attraction")
	}

	

	return attraction, nil
}

// delete an attraction completely
func (p attractionRepo) DeleteAttraction(ctx context.Context, attractionID string) error {
	// Build the SQL query
	sqlStr, args, err := p.db.Sq.Builder.Delete(p.tableName).
		Where(p.db.Sq.Equal("attraction_id", attractionID)).
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
