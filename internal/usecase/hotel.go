package usecase

import (
	"Booking/establishment-service-booking/internal/entity"
	"Booking/establishment-service-booking/internal/infrastructure/repository"
	"Booking/establishment-service-booking/internal/pkg/otlp"
	"context"
	"time"
)

const (
	hotelServiceName = "hotelService"
	spanNameHotel    = "hotelUsecase"
)

type Hotel interface {
	CreateHotel(ctx context.Context, hotel *entity.Hotel) (*entity.Hotel, error)
	GetHotel(ctx context.Context, hotel_id string) (*entity.Hotel, error)
	ListHotels(ctx context.Context, page, limit int64) ([]*entity.Hotel, uint64, error)
	UpdateHotel(ctx context.Context, hotel *entity.Hotel) (*entity.Hotel, error)
	DeleteHotel(ctx context.Context, hotel_id string) error
	ListHotelsByLocation(ctx context.Context, offset, limit uint64, country, city, state_province string) ([]*entity.Hotel, int64, error)
	FindHotelsByName(ctx context.Context, name string) ([]*entity.Hotel, uint64, error)
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

	ctx, span := otlp.Start(ctx, hotelServiceName, spanNameHotel+"Create")
	defer span.End()

	return h.repo.CreateHotel(ctx, hotel)
}

func (h HotelService) GetHotel(ctx context.Context, hotel_id string) (*entity.Hotel, error) {
	ctx, cancel := context.WithTimeout(ctx, h.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, hotelServiceName, spanNameHotel+"Get")
	defer span.End()

	return h.repo.GetHotel(ctx, hotel_id)
}

func (h HotelService) ListHotels(ctx context.Context, offset, limit int64) ([]*entity.Hotel, uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, h.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, hotelServiceName, spanNameHotel+"List")
	defer span.End()

	return h.repo.ListHotels(ctx, offset, limit)
}

func (h HotelService) UpdateHotel(ctx context.Context, hotel *entity.Hotel) (*entity.Hotel, error) {
	ctx, cancel := context.WithTimeout(ctx, h.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, hotelServiceName, spanNameHotel+"Update")
	defer span.End()

	return h.repo.UpdateHotel(ctx, hotel)
}

func (h HotelService) DeleteHotel(ctx context.Context, hotel_id string) error {
	ctx, cancel := context.WithTimeout(ctx, h.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, hotelServiceName, spanNameHotel+"Delete")
	defer span.End()

	return h.repo.DeleteHotel(ctx, hotel_id)
}

func (h HotelService) ListHotelsByLocation(ctx context.Context, offset, limit uint64, country, city, state_province string) ([]*entity.Hotel, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, h.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, hotelServiceName, spanNameHotel+"ListL")
	defer span.End()

	return h.repo.ListHotelsByLocation(ctx, offset, limit, country, city, state_province)
}

func (h HotelService) FindHotelsByName(ctx context.Context, name string) ([]*entity.Hotel, uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, h.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, hotelServiceName, spanNameHotel+"List")
	defer span.End()

	return h.repo.FindHotelsByName(ctx, name)
}
