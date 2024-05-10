package usecase

import (
	"Booking/establishment-service-booking/internal/entity"
	"Booking/establishment-service-booking/internal/infrastructure/repository"
	"context"
	"time"
)

type Review interface {
	CreateReview(ctx context.Context, review *entity.Review) (*entity.Review, error)
	ListReviews(ctx context.Context, establishment_id string) ([]*entity.Review, uint64, error)
	DeleteReview(ctx context.Context, review_id string) error
}

type ReviewService struct {
	BaseUseCase
	repo       repository.Review
	ctxTimeout time.Duration
}

func NewReviewService(ctxTimeout time.Duration, repo repository.Review) ReviewService {
	return ReviewService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (r ReviewService) CreateReview(ctx context.Context, review *entity.Review) (*entity.Review, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	return r.repo.CreateReview(ctx, review)
}

func (r ReviewService) ListReviews(ctx context.Context, establishment_id string) ([]*entity.Review, uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	return r.repo.ListReviews(ctx, establishment_id)
}

func (r ReviewService) DeleteReview(ctx context.Context, establishment_id string) error {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	return r.repo.DeleteReview(ctx, establishment_id)
}
