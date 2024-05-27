package repository

import (
	"Booking/establishment-service-booking/internal/entity"
	"context"
)

type Hotel interface {
	CreateHotel(ctx context.Context, Hotel *entity.Hotel) (*entity.Hotel, error)
	GetHotel(ctx context.Context, hotel_id string) (*entity.Hotel, error)
	ListHotels(ctx context.Context, offset, limit int64) ([]*entity.Hotel, uint64, error)
	UpdateHotel(ctx context.Context, Hotel *entity.Hotel) (*entity.Hotel, error)
	DeleteHotel(ctx context.Context, hotel_id string) error
	ListHotelsByLocation(ctx context.Context, offset, limit uint64, country, city, state_province string) ([]*entity.Hotel, int64, error)
	FindHotelsByName(ctx context.Context, name string) ([]*entity.Hotel, uint64, error)
}
