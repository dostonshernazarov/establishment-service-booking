package usecase

import (
	"Booking/establishment-service-booking/internal/entity"
	"Booking/establishment-service-booking/internal/infrastructure/repository"
	"context"
	"time"
)

// const (
// 	serviceNameuser = "userService"
// 	spanNameuser    = "userUsecase"
// )

type Hotel interface {
	CreateHotel(ctx context.Context, hotel *entity.Hotel) (*entity.Hotel, error)
	GetHotel(ctx context.Context, hotel_id string) (*entity.Hotel, error)
	ListHotels(ctx context.Context, page, limit int64) ([]*entity.Hotel, error)
	UpdateHotel(ctx context.Context, hotel *entity.Hotel) (*entity.Hotel, error)
	DeleteHotel(ctx context.Context, hotel_id string) error
}

type HotelService struct {
	BaseUseCase
	repo       repository.Hotel
	ctxTimeout time.Duration
}

func NewHotelService(ctxTimeout time.Duration, repo repository.Hotel) HotelService {
	return HotelService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (h HotelService) CreateHotel(ctx context.Context, hotel *entity.Hotel) (*entity.Hotel, error) {
	ctx, cancel := context.WithTimeout(ctx, h.ctxTimeout)
	defer cancel()

	h.beforeRequest(&hotel.HotelId, &hotel.CreatedAt, &hotel.UpdatedAt, nil)

	return h.repo.CreateHotel(ctx, hotel)
}

func (h HotelService) GetHotel(ctx context.Context, hotel_id string) (*entity.Hotel, error) {
	ctx, cancel := context.WithTimeout(ctx, h.ctxTimeout)
	defer cancel()

	return h.repo.GetHotel(ctx, hotel_id)
}

func (h HotelService) ListHotels(ctx context.Context, page, limit int64) ([]*entity.Hotel, error) {
	ctx, cancel := context.WithTimeout(ctx, h.ctxTimeout)
	defer cancel()

	return h.repo.ListHotels(ctx, page, limit)
}

func (h HotelService) UpdateHotel(ctx context.Context, hotel *entity.Hotel) (*entity.Hotel, error) {
	ctx, cancel := context.WithTimeout(ctx, h.ctxTimeout)
	defer cancel()

	h.beforeRequest(nil, nil, &hotel.UpdatedAt, nil)

	return h.repo.UpdateHotel(ctx, hotel)
}

func (h HotelService) DeleteHotel(ctx context.Context, hotel_id string) error {
	ctx, cancel := context.WithTimeout(ctx, h.ctxTimeout)
	defer cancel()

	return h.repo.DeleteHotel(ctx, hotel_id)
}
