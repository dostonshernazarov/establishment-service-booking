package usecase

import (
	"Booking/establishment-service-booking/internal/entity"
	"Booking/establishment-service-booking/internal/infrastructure/repository"
	"Booking/establishment-service-booking/internal/pkg/otlp"
	"context"
	"time"
)

const (
	imageServiceName = "imageService"
	spanNameImage    = "imageUsecase"
)

type Image interface {
	CreateImage(ctx context.Context, image *entity.Image) error
}

type ImageService struct {
	BaseUseCase
	repo       repository.Image
	ctxTimeout time.Duration
}


func NewImageService(ctxTimeout time.Duration, repo repository.Image) ImageService {
	return ImageService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (h ImageService) CreateImage(ctx context.Context, image *entity.Image) error {
	ctx, cancel := context.WithTimeout(ctx, h.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, imageServiceName, spanNameImage+"Create")
	defer span.End()

	return h.repo.CreateImage(ctx, image)
}
