package repository

import (
	"Booking/establishment-service-booking/internal/entity"
	"context"
)

type Restaurant interface {
	CreateRestaurant(ctx context.Context, restaurant *entity.Restaurant) (*entity.Restaurant, error)
	GetRestaurant(ctx context.Context, restaurant_id string) (*entity.Restaurant, error)
	ListRestaurants(ctx context.Context, offset, limit int64) ([]*entity.Restaurant, error)
	UpdateRestaurant(ctx context.Context, restaurant *entity.Restaurant) (*entity.Restaurant, error)
	DeleteRestaurant(ctx context.Context, restaurant_id string) error
}
