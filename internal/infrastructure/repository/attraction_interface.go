package repository

import (
	"Booking/establishment-service-booking/internal/entity"
	"context"
)

type Attraction interface {
	CreateAttraction(ctx context.Context, attraction *entity.Attraction) (*entity.Attraction, error)
	GetAttraction(ctx context.Context, attraction_id string) (*entity.Attraction, error)
	ListAttractions(ctx context.Context, page, limit int64) ([]*entity.Attraction, error)
	UpdateAttraction(ctx context.Context, attraction *entity.Attraction) (*entity.Attraction, error)
	DeleteAttraction(ctx context.Context, attraction_id string) error
}
