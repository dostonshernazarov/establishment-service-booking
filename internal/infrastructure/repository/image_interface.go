package repository

import (
	"Booking/establishment-service-booking/internal/entity"
	"context"
)

type Image interface{
	CreateImage(ctx context.Context, image *entity.Image) error
}