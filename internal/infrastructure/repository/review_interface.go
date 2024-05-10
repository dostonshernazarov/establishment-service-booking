package repository

import (
	"Booking/establishment-service-booking/internal/entity"
	"context"
)

type Review interface {
	CreateReview(ctx context.Context, review *entity.Review) (*entity.Review, error)
	ListReviews(ctx context.Context, establishment_id string) ([]*entity.Review, uint64, error)
	DeleteReview(ctx context.Context, review_id string) error
}
