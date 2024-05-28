package postgresql

import (
	"Booking/establishment-service-booking/internal/entity"
	"Booking/establishment-service-booking/internal/pkg/otlp"
	"Booking/establishment-service-booking/internal/pkg/postgres"
	"context"
	"fmt"
)


type imageRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewImageRepo(db *postgres.PostgresDB) *imageRepo {
	return &imageRepo{
		tableName: imageTableName,
		db:        db,
	}
}

func (p imageRepo) CreateImage(ctx context.Context, image *entity.Image) error {
	ctx, span := otlp.Start(ctx, hotelServiceName, hotelSpanRepoPrefix+"CreateImage")
	defer span.End()

	dataI := map[string]interface{}{
		"image_id":      image.ImageId,
		"establishment_id": image.EstablishmentId,
		"image_url": image.ImageUrl,
		"category": image.Category,
		"created_at": image.CreatedAt,
		"updated_at": image.UpdatedAt,

	}

	query, args, err := p.db.Sq.Builder.Insert(imageTableName).SetMap(dataI).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query for creating establishment's image: %v", err)
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute SQL query for creating establishment's image: %v", err)
	}

	return nil
}