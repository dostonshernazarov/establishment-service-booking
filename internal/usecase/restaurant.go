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

type Restaurant interface {
	CreateRestaurant(ctx context.Context, restaurant *entity.Restaurant) (*entity.Restaurant, error)
	GetRestaurant(ctx context.Context, restaurant_id string) (*entity.Restaurant, error)
	ListRestaurants(ctx context.Context, page, limit int64) ([]*entity.Restaurant, uint64, error)
	UpdateRestaurant(ctx context.Context, restaurant *entity.Restaurant) (*entity.Restaurant, error)
	DeleteRestaurant(ctx context.Context, restaurant_id string) error
}

type RestaurantService struct {
	BaseUseCase
	repo       repository.Restaurant
	ctxTimeout time.Duration
}

func NewRestaurantService(ctxTimeout time.Duration, repo repository.Restaurant) RestaurantService {
	return RestaurantService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (r RestaurantService) CreateRestaurant(ctx context.Context, restaurant *entity.Restaurant) (*entity.Restaurant, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	r.beforeRequest(nil, &restaurant.CreatedAt, &restaurant.UpdatedAt, nil)

	return r.repo.CreateRestaurant(ctx, restaurant)
}

func (r RestaurantService) GetRestaurant(ctx context.Context, restaurant_id string) (*entity.Restaurant, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	return r.repo.GetRestaurant(ctx, restaurant_id)
}

func (r RestaurantService) ListRestaurants(ctx context.Context, offset, limit int64) ([]*entity.Restaurant, uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	return r.repo.ListRestaurants(ctx, offset, limit)
}

func (r RestaurantService) UpdateRestaurant(ctx context.Context, restaurant *entity.Restaurant) (*entity.Restaurant, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	r.beforeRequest(nil, nil, &restaurant.UpdatedAt, nil)

	return r.repo.UpdateRestaurant(ctx, restaurant)
}

func (r RestaurantService) DeleteRestaurant(ctx context.Context, restaurant_id string) error {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	return r.repo.DeleteRestaurant(ctx, restaurant_id)
}
