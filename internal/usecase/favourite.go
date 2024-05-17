package usecase

import (
	"Booking/establishment-service-booking/internal/entity"
	"Booking/establishment-service-booking/internal/infrastructure/repository"
	"Booking/establishment-service-booking/internal/pkg/otlp"
	"context"
	"time"
)

const (
	favouriteServiceName = "favouriteService"
	spanNameFavourite    = "favouriteUsecase"
)

type Favourite interface {
	AddToFavourites(ctx context.Context, favourite *entity.Favourite) (*entity.Favourite, error)
	RemoveFromFavourites(ctx context.Context, favourite_id string) error
	ListFavouritesByUserId(ctx context.Context, user_id string) ([]*entity.Favourite, error)
}

type FavouriteService struct {
	BaseUseCase
	repo       repository.Favourite
	ctxTimeout time.Duration
}

func NewFavouriteService(ctxTimeout time.Duration, repo repository.Favourite) FavouriteService {
	return FavouriteService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (f FavouriteService) AddToFavourites(ctx context.Context, favourite *entity.Favourite) (*entity.Favourite, error) {
	ctx, cancel := context.WithTimeout(ctx, f.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, favouriteServiceName, spanNameFavourite+"Create")
	defer span.End()

	return f.repo.AddToFavourites(ctx, favourite)
}

func (f FavouriteService) RemoveFromFavourites(ctx context.Context, favourite_id string) error {
	ctx, cancel := context.WithTimeout(ctx, f.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, favouriteServiceName, spanNameFavourite+"Delete")
	defer span.End()

	return f.repo.RemoveFromFavourites(ctx, favourite_id)
}

func (f FavouriteService) ListFavouritesByUserId(ctx context.Context, user_id string) ([]*entity.Favourite, error) {
	ctx, cancel := context.WithTimeout(ctx, f.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, favouriteServiceName, spanNameFavourite+"List")
	defer span.End()

	return f.repo.ListFavouritesByUserId(ctx, user_id)
}
