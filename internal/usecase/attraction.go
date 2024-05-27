package usecase

import (
	"Booking/establishment-service-booking/internal/entity"
	"Booking/establishment-service-booking/internal/infrastructure/repository"
	"Booking/establishment-service-booking/internal/pkg/otlp"
	"context"
	"time"
)

const (
	attractionServiceName = "attractionService"
	spanNameAttraction    = "attractionUsecase"
)

type Attraction interface {
	CreateAttraction(ctx context.Context, attracation *entity.Attraction) (*entity.Attraction, error)
	GetAttraction(ctx context.Context, attraction_id string) (*entity.Attraction, error)
	ListAttractions(ctx context.Context, page, limit int64) ([]*entity.Attraction, uint64, error)
	UpdateAttraction(ctx context.Context, attracation *entity.Attraction) (*entity.Attraction, error)
	DeleteAttraction(ctx context.Context, attraction_id string) error
	ListAttractionsByLocation(ctx context.Context, offset, limit uint64, country, city, state_province string) ([]*entity.Attraction, int64, error)
	FindAttractionsByName(ctx context.Context, name string) ([]*entity.Attraction, uint64, error)
}

type AttractionService struct {
	BaseUseCase
	repo       repository.Attraction
	ctxTimeout time.Duration
}

func NewAttractionService(ctxTimeout time.Duration, repo repository.Attraction) AttractionService {
	return AttractionService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (a AttractionService) CreateAttraction(ctx context.Context, attracation *entity.Attraction) (*entity.Attraction, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, attractionServiceName, spanNameAttraction+"Create")
	defer span.End()

	return a.repo.CreateAttraction(ctx, attracation)
}

func (a AttractionService) GetAttraction(ctx context.Context, attraction_id string) (*entity.Attraction, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, attractionServiceName, spanNameAttraction+"Get")
	defer span.End()

	return a.repo.GetAttraction(ctx, attraction_id)
}

func (a AttractionService) ListAttractions(ctx context.Context, offset, limit int64) ([]*entity.Attraction, uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, attractionServiceName, spanNameAttraction+"List")
	defer span.End()

	return a.repo.ListAttractions(ctx, offset, limit)
}

func (a AttractionService) UpdateAttraction(ctx context.Context, attracation *entity.Attraction) (*entity.Attraction, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, attractionServiceName, spanNameAttraction+"Update")
	defer span.End()

	return a.repo.UpdateAttraction(ctx, attracation)
}

func (a AttractionService) DeleteAttraction(ctx context.Context, attraction_id string) error {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, attractionServiceName, spanNameAttraction+"Delete")
	defer span.End()

	return a.repo.DeleteAttraction(ctx, attraction_id)
}

func (a AttractionService) ListAttractionsByLocation(ctx context.Context, offset, limit uint64, country, city, state_province string) ([]*entity.Attraction, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, attractionServiceName, spanNameAttraction+"ListL")
	defer span.End()

	return a.repo.ListAttractionsByLocation(ctx, offset, limit, country, city, state_province)
}

func (a AttractionService) FindAttractionsByName(ctx context.Context, name string) ([]*entity.Attraction, uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, attractionServiceName, spanNameAttraction+"ListL")
	defer span.End()

	return a.repo.FindAttractionsByName(ctx, name)
}
