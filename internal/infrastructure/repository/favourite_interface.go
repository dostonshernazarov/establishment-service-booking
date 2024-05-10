package repository

import (
	"Booking/establishment-service-booking/internal/entity"
	"context"
)

type Favourite interface {
	AddToFavourites(ctx context.Context, favourite *entity.Favourite) (*entity.Favourite, error)
	RemoveFromFavourites(ctx context.Context, favourite_id string) error
	ListFavouritesByUserId(ctx context.Context, user_id string) ([]*entity.Favourite, error)
}
