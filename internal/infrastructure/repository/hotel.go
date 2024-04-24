package repository

import (
	"Booking/establishment-service-booking/internal/entity"
	"context"
)

type Hotel interface {
	CreateHotel(ctx context.Context, Hotel *entity.Hotel) (*entity.Hotel, error)
	GetHotel(ctx context.Context, hotel_id string) (*entity.Hotel, error)
	ListHotels(ctx context.Context, page, limit int64) ([]*entity.Hotel, error)
	UpdateHotel(ctx context.Context, Hotel *entity.Hotel) (*entity.Hotel, error)
	DeleteHotel(ctx context.Context, hotel_id string) error
}
