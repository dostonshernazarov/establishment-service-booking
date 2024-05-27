package repository

import (
	"Booking/establishment-service-booking/internal/entity"
	"context"
)

type Attraction interface {
	CreateAttraction(ctx context.Context, attraction *entity.Attraction) (*entity.Attraction, error)
	GetAttraction(ctx context.Context, attraction_id string) (*entity.Attraction, error)
	ListAttractions(ctx context.Context, offset, limit int64) ([]*entity.Attraction, uint64, error)
	UpdateAttraction(ctx context.Context, attraction *entity.Attraction) (*entity.Attraction, error)
	DeleteAttraction(ctx context.Context, attraction_id string) error
	ListAttractionsByLocation(ctx context.Context, offset, limit uint64, country, city, state_province string) ([]*entity.Attraction, int64, error)
	FindAttractionsByName(ctx context.Context, name string) ([]*entity.Attraction, uint64, error)
}
