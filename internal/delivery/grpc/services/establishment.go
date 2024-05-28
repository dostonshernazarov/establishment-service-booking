package services

import (
	pb "Booking/establishment-service-booking/genproto/establishment-proto"
	"Booking/establishment-service-booking/internal/entity"
	"Booking/establishment-service-booking/internal/pkg/otlp"
	"Booking/establishment-service-booking/internal/usecase"
	"Booking/establishment-service-booking/internal/usecase/event"
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type establishmentRPC struct {
	logger             *zap.Logger
	attracationUsecase usecase.Attraction
	restaurantUsecase  usecase.Restaurant
	hotelUsecase       usecase.Hotel
	favouriteUsecase   usecase.Favourite
	imageUsecase       usecase.Image
	reviewUsecase      usecase.Review
	brokerProducer     event.BrokerProducer
}

func NewRPC(logger *zap.Logger, attracationUsecase usecase.Attraction, restaurantUsecase usecase.Restaurant, hotelUsecase usecase.Hotel, favouriteUsecase usecase.Favourite, imageUsecase usecase.Image, reviewUsecase usecase.Review, brokerProducer event.BrokerProducer) pb.EstablishmentServiceServer {
	return &establishmentRPC{
		logger:             logger,
		attracationUsecase: attracationUsecase,
		restaurantUsecase:  restaurantUsecase,
		hotelUsecase:       hotelUsecase,
		favouriteUsecase:   favouriteUsecase,
		reviewUsecase:      reviewUsecase,
		imageUsecase:       imageUsecase,
		brokerProducer:     brokerProducer,
	}
}

//****************************IMPLEMENTATIONS****************************//

// ATTRACTION
func (s establishmentRPC) CreateAttraction(ctx context.Context, attraction *pb.Attraction) (*pb.Attraction, error) {
	ctx, span := otlp.Start(ctx, "attraction_grpc_delivery", "Create")
	span.SetAttributes(
		attribute.Key("attraction_id").String(attraction.AttractionId),
	)
	defer span.End()

	var images []*entity.Image

	for _, i := range attraction.Images {
		var image entity.Image

		image.ImageId = i.ImageId
		image.EstablishmentId = i.EstablishmentId
		image.ImageUrl = i.ImageUrl
		image.Category = i.Category
		image.CreatedAt = time.Now().Local()
		image.UpdatedAt = time.Now().Local()

		images = append(images, &image)
	}
	response, err := s.attracationUsecase.CreateAttraction(ctx, &entity.Attraction{
		AttractionId:   attraction.AttractionId,
		OwnerId:        attraction.OwnerId,
		AttractionName: attraction.AttractionName,
		Description:    attraction.Description,
		Rating:         attraction.Rating,
		ContactNumber:  attraction.ContactNumber,
		LicenceUrl:     attraction.LicenceUrl,
		WebsiteUrl:     attraction.WebsiteUrl,
		Images:         images,
		Location: entity.Location{
			LocationId:      attraction.Location.LocationId,
			EstablishmentId: attraction.Location.EstablishmentId,
			Address:         attraction.Location.Address,
			Latitude:        attraction.Location.Latitude,
			Longitude:       attraction.Location.Longitude,
			Country:         attraction.Location.Country,
			City:            attraction.Location.City,
			StateProvince:   attraction.Location.StateProvince,
			Category:        attraction.Location.Category,
			CreatedAt:       time.Now().Local(),
			UpdatedAt:       time.Now().Local(),
		},
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
	})
	if err != nil {
		return nil, err
	}

	var respImages []*pb.Image

	for _, respImage := range response.Images {
		image := pb.Image{
			ImageId:         respImage.ImageId,
			EstablishmentId: respImage.EstablishmentId,
			ImageUrl:        respImage.ImageUrl,
			CreatedAt:       respImage.CreatedAt.String(),
			UpdatedAt:       respImage.UpdatedAt.String(),
		}

		respImages = append(respImages, &image)
	}

	return &pb.Attraction{
		AttractionId:   response.AttractionId,
		OwnerId:        response.OwnerId,
		AttractionName: response.AttractionName,
		Description:    response.Description,
		Rating:         response.Rating,
		ContactNumber:  response.ContactNumber,
		LicenceUrl:     response.LicenceUrl,
		WebsiteUrl:     response.WebsiteUrl,
		Images:         respImages,
		Location: &pb.Location{
			LocationId:      response.Location.LocationId,
			EstablishmentId: response.Location.EstablishmentId,
			Address:         response.Location.Address,
			Latitude:        response.Location.Latitude,
			Longitude:       response.Location.Longitude,
			Country:         response.Location.Country,
			City:            response.Location.City,
			StateProvince:   response.Location.StateProvince,
			CreatedAt:       response.Location.CreatedAt.String(),
			UpdatedAt:       response.Location.UpdatedAt.String(),
		},
		CreatedAt: response.CreatedAt.String(),
		UpdatedAt: response.UpdatedAt.String(),
	}, nil
}

func (s establishmentRPC) GetAttraction(ctx context.Context, request *pb.GetAttractionRequest) (*pb.GetAttractionResponse, error) {
	ctx, span := otlp.Start(ctx, "attraction_grpc_delivery", "Get")
	span.SetAttributes(
		attribute.Key("attraction_id").String(request.AttractionId),
	)
	defer span.End()

	attraction, err := s.attracationUsecase.GetAttraction(ctx, request.AttractionId)
	if err != nil {
		return nil, err
	}

	var images []*pb.Image

	for _, i := range attraction.Images {
		var image pb.Image

		image.ImageId = i.ImageId
		image.EstablishmentId = i.EstablishmentId
		image.ImageUrl = i.ImageUrl
		image.CreatedAt = i.CreatedAt.String()
		image.UpdatedAt = i.UpdatedAt.String()

		images = append(images, &image)
	}

	return &pb.GetAttractionResponse{
		Attraction: &pb.Attraction{
			AttractionId:   attraction.AttractionId,
			OwnerId:        attraction.OwnerId,
			AttractionName: attraction.AttractionName,
			Description:    attraction.Description,
			Rating:         attraction.Rating,
			ContactNumber:  attraction.ContactNumber,
			LicenceUrl:     attraction.LicenceUrl,
			WebsiteUrl:     attraction.WebsiteUrl,
			Images:         images,
			Location: &pb.Location{
				LocationId:      attraction.Location.LocationId,
				EstablishmentId: attraction.Location.EstablishmentId,
				Address:         attraction.Location.Address,
				Latitude:        attraction.Location.Latitude,
				Longitude:       attraction.Location.Longitude,
				Country:         attraction.Location.Country,
				City:            attraction.Location.City,
				StateProvince:   attraction.Location.StateProvince,
				CreatedAt:       attraction.CreatedAt.String(),
				UpdatedAt:       attraction.UpdatedAt.String(),
			},
			CreatedAt: attraction.CreatedAt.String(),
			UpdatedAt: attraction.UpdatedAt.String(),
		},
	}, nil
}

func (s establishmentRPC) ListAttractions(ctx context.Context, request *pb.ListAttractionsRequest) (*pb.ListAttractionsResponse, error) {
	ctx, span := otlp.Start(ctx, "attraction_grpc_delivery", "List")
	span.SetAttributes(
		attribute.Key("limit").Int64(request.Limit),
	)
	defer span.End()

	attractions, overall, err := s.attracationUsecase.ListAttractions(ctx, request.Offset, request.Limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch attractions: %v", err)
	}

	var pbAttractions []*pb.Attraction
	for _, attraction := range attractions {
		var images []*pb.Image
		for _, i := range attraction.Images {
			images = append(images, &pb.Image{
				ImageId:         i.ImageId,
				EstablishmentId: i.EstablishmentId,
				ImageUrl:        i.ImageUrl,
				CreatedAt:       i.CreatedAt.String(),
				UpdatedAt:       i.UpdatedAt.String(),
			})
		}

		pbAttractions = append(pbAttractions, &pb.Attraction{
			AttractionId:   attraction.AttractionId,
			OwnerId:        attraction.OwnerId,
			AttractionName: attraction.AttractionName,
			Description:    attraction.Description,
			Rating:         attraction.Rating,
			ContactNumber:  attraction.ContactNumber,
			LicenceUrl:     attraction.LicenceUrl,
			WebsiteUrl:     attraction.WebsiteUrl,
			Images:         images,
			Location: &pb.Location{
				LocationId:      attraction.Location.LocationId,
				EstablishmentId: attraction.Location.EstablishmentId,
				Address:         attraction.Location.Address,
				Latitude:        attraction.Location.Latitude,
				Longitude:       attraction.Location.Longitude,
				Country:         attraction.Location.Country,
				City:            attraction.Location.City,
				StateProvince:   attraction.Location.StateProvince,
				CreatedAt:       attraction.CreatedAt.String(),
				UpdatedAt:       attraction.UpdatedAt.String(),
			},
			CreatedAt: attraction.CreatedAt.String(),
			UpdatedAt: attraction.UpdatedAt.String(),
		})
	}

	return &pb.ListAttractionsResponse{
		Attractions: pbAttractions,
		Overall:     overall,
	}, nil
}

func (s establishmentRPC) UpdateAttraction(ctx context.Context, request *pb.UpdateAttractionRequest) (*pb.UpdateAttractionResponse, error) {

	ctx, span := otlp.Start(ctx, "attraction_grpc_delivery", "Update")
	span.SetAttributes(
		attribute.Key("attraction_id").String(request.Attraction.AttractionId),
	)
	defer span.End()

	// var imagesS []*entity.Image

	// for _, i := range request.Attraction.Images {
	// 	var image entity.Image

	// 	image.ImageId = i.ImageId
	// 	image.EstablishmentId = i.EstablishmentId
	// 	image.ImageUrl = i.ImageUrl

	// 	imagesS = append(imagesS, &image)
	// }
	attraction, err := s.attracationUsecase.UpdateAttraction(ctx, &entity.Attraction{
		AttractionId:   request.Attraction.AttractionId,
		OwnerId:        request.Attraction.OwnerId,
		AttractionName: request.Attraction.AttractionName,
		Description:    request.Attraction.Description,
		Rating:         request.Attraction.Rating,
		ContactNumber:  request.Attraction.ContactNumber,
		LicenceUrl:     request.Attraction.LicenceUrl,
		WebsiteUrl:     request.Attraction.WebsiteUrl,
		// Images:         imagesS,
		Location: entity.Location{
			LocationId:      request.Attraction.Location.LocationId,
			EstablishmentId: request.Attraction.Location.EstablishmentId,
			Address:         request.Attraction.Location.Address,
			Latitude:        request.Attraction.Location.Latitude,
			Longitude:       request.Attraction.Location.Longitude,
			Country:         request.Attraction.Location.Country,
			City:            request.Attraction.Location.City,
			StateProvince:   request.Attraction.Location.StateProvince,
		},
	})
	if err != nil {
		return nil, err
	}

	var images []*pb.Image

	for _, i := range attraction.Images {
		var image pb.Image

		image.ImageId = i.ImageId
		image.EstablishmentId = i.EstablishmentId
		image.ImageUrl = i.ImageUrl
		image.CreatedAt = i.CreatedAt.String()
		image.UpdatedAt = i.UpdatedAt.String()

		images = append(images, &image)
	}

	return &pb.UpdateAttractionResponse{
		Attraction: &pb.Attraction{
			AttractionId:   attraction.AttractionId,
			OwnerId:        attraction.OwnerId,
			AttractionName: attraction.AttractionName,
			Description:    attraction.Description,
			Rating:         attraction.Rating,
			ContactNumber:  attraction.ContactNumber,
			LicenceUrl:     attraction.LicenceUrl,
			WebsiteUrl:     attraction.WebsiteUrl,
			Images:         images,
			Location: &pb.Location{
				LocationId:      attraction.Location.LocationId,
				EstablishmentId: attraction.Location.EstablishmentId,
				Address:         attraction.Location.Address,
				Latitude:        attraction.Location.Latitude,
				Longitude:       attraction.Location.Longitude,
				Country:         attraction.Location.Country,
				City:            attraction.Location.City,
				StateProvince:   attraction.Location.StateProvince,
				CreatedAt:       attraction.CreatedAt.String(),
				UpdatedAt:       attraction.UpdatedAt.String(),
			},
			CreatedAt: attraction.CreatedAt.String(),
			UpdatedAt: attraction.UpdatedAt.String(),
		},
	}, nil
}

func (s establishmentRPC) DeleteAttraction(ctx context.Context, request *pb.DeleteAttractionRequest) (*pb.DeleteAttractionResponse, error) {
	ctx, span := otlp.Start(ctx, "attraction_grpc_delivery", "Delete")
	span.SetAttributes(
		attribute.Key("attraction_id").String(request.AttractionId),
	)
	defer span.End()

	if err := s.attracationUsecase.DeleteAttraction(ctx, request.AttractionId); err != nil {
		return &pb.DeleteAttractionResponse{
			Success: false,
		}, err
	}

	return &pb.DeleteAttractionResponse{
		Success: true,
	}, nil
}

func (s establishmentRPC) ListAttractionsByLocation(ctx context.Context, request *pb.ListAttractionsByLocationRequest) (*pb.ListAttractionsByLocationResponse, error) {
	ctx, span := otlp.Start(ctx, "attraction_grpc_delivery", "List")
	span.SetAttributes(
		attribute.Key("state_province").String(request.StateProvince),
	)
	defer span.End()

	attractions, count, err := s.attracationUsecase.ListAttractionsByLocation(ctx, request.Offset, request.Limit, request.Country, request.City, request.StateProvince)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch attractions: %v", err)
	}

	var pbAttractions []*pb.Attraction
	for _, attraction := range attractions {
		var images []*pb.Image
		for _, i := range attraction.Images {
			images = append(images, &pb.Image{
				ImageId:         i.ImageId,
				EstablishmentId: i.EstablishmentId,
				ImageUrl:        i.ImageUrl,
				CreatedAt:       i.CreatedAt.String(),
				UpdatedAt:       i.UpdatedAt.String(),
			})
		}

		pbAttractions = append(pbAttractions, &pb.Attraction{
			AttractionId:   attraction.AttractionId,
			OwnerId:        attraction.OwnerId,
			AttractionName: attraction.AttractionName,
			Description:    attraction.Description,
			Rating:         attraction.Rating,
			ContactNumber:  attraction.ContactNumber,
			LicenceUrl:     attraction.LicenceUrl,
			WebsiteUrl:     attraction.WebsiteUrl,
			Images:         images,
			Location: &pb.Location{
				LocationId:      attraction.Location.LocationId,
				EstablishmentId: attraction.Location.EstablishmentId,
				Address:         attraction.Location.Address,
				Latitude:        attraction.Location.Latitude,
				Longitude:       attraction.Location.Longitude,
				Country:         attraction.Location.Country,
				City:            attraction.Location.City,
				StateProvince:   attraction.Location.StateProvince,
				CreatedAt:       attraction.CreatedAt.String(),
				UpdatedAt:       attraction.UpdatedAt.String(),
			},
			CreatedAt: attraction.CreatedAt.String(),
			UpdatedAt: attraction.UpdatedAt.String(),
		})
	}

	return &pb.ListAttractionsByLocationResponse{
		Attractions: pbAttractions,
		Count:       count,
	}, nil
}

func (s establishmentRPC) FindAttractionsByName(ctx context.Context, request *pb.FindAttractionsByNameRequest) (*pb.FindAttractionsByNameResponse, error) {
	ctx, span := otlp.Start(ctx, "attraction_grpc_delivery", "Find")
	span.SetAttributes(
		attribute.Key("name").String(request.Name),
	)
	defer span.End()

	attractions, overall, err := s.attracationUsecase.FindAttractionsByName(ctx, request.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch attractions: %v", err)
	}

	var pbAttractions []*pb.Attraction
	for _, attraction := range attractions {
		var images []*pb.Image
		for _, i := range attraction.Images {
			images = append(images, &pb.Image{
				ImageId:         i.ImageId,
				EstablishmentId: i.EstablishmentId,
				ImageUrl:        i.ImageUrl,
				CreatedAt:       i.CreatedAt.String(),
				UpdatedAt:       i.UpdatedAt.String(),
			})
		}

		pbAttractions = append(pbAttractions, &pb.Attraction{
			AttractionId:   attraction.AttractionId,
			OwnerId:        attraction.OwnerId,
			AttractionName: attraction.AttractionName,
			Description:    attraction.Description,
			Rating:         attraction.Rating,
			ContactNumber:  attraction.ContactNumber,
			LicenceUrl:     attraction.LicenceUrl,
			WebsiteUrl:     attraction.WebsiteUrl,
			Images:         images,
			Location: &pb.Location{
				LocationId:      attraction.Location.LocationId,
				EstablishmentId: attraction.Location.EstablishmentId,
				Address:         attraction.Location.Address,
				Latitude:        attraction.Location.Latitude,
				Longitude:       attraction.Location.Longitude,
				Country:         attraction.Location.Country,
				City:            attraction.Location.City,
				StateProvince:   attraction.Location.StateProvince,
				CreatedAt:       attraction.CreatedAt.String(),
				UpdatedAt:       attraction.UpdatedAt.String(),
			},
			CreatedAt: attraction.CreatedAt.String(),
			UpdatedAt: attraction.UpdatedAt.String(),
		})
	}

	return &pb.FindAttractionsByNameResponse{
		Attractions: pbAttractions,
		Count:       overall,
	}, nil
}

// RESTAURANT
func (s establishmentRPC) CreateRestaurant(ctx context.Context, restaurant *pb.Restaurant) (*pb.Restaurant, error) {
	ctx, span := otlp.Start(ctx, "restaurant_grpc_delivery", "Create")
	span.SetAttributes(
		attribute.Key("restaurant_id").String(restaurant.RestaurantId),
	)
	defer span.End()

	var images []*entity.Image

	for _, i := range restaurant.Images {
		var image entity.Image

		image.ImageId = i.ImageId
		image.EstablishmentId = i.EstablishmentId
		image.ImageUrl = i.ImageUrl
		image.Category = i.Category
		image.CreatedAt = time.Now().Local()
		image.UpdatedAt = time.Now().Local()

		images = append(images, &image)
	}
	response, err := s.restaurantUsecase.CreateRestaurant(ctx, &entity.Restaurant{
		RestaurantId:   restaurant.RestaurantId,
		OwnerId:        restaurant.OwnerId,
		RestaurantName: restaurant.RestaurantName,
		Description:    restaurant.Description,
		Rating:         restaurant.Rating,
		OpeningHours:   restaurant.OpeningHours,
		ContactNumber:  restaurant.ContactNumber,
		LicenceUrl:     restaurant.LicenceUrl,
		WebsiteUrl:     restaurant.WebsiteUrl,
		Images:         images,
		Location: entity.Location{
			LocationId:      restaurant.Location.LocationId,
			EstablishmentId: restaurant.Location.EstablishmentId,
			Address:         restaurant.Location.Address,
			Latitude:        restaurant.Location.Latitude,
			Longitude:       restaurant.Location.Longitude,
			Country:         restaurant.Location.Country,
			City:            restaurant.Location.City,
			StateProvince:   restaurant.Location.StateProvince,
			Category:        restaurant.Location.Category,
			CreatedAt:       time.Now().Local(),
			UpdatedAt:       time.Now().Local(),
		},
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
	})
	if err != nil {
		return nil, err
	}

	var respImages []*pb.Image

	for _, respImage := range response.Images {
		image := pb.Image{
			ImageId:         respImage.ImageId,
			EstablishmentId: respImage.EstablishmentId,
			ImageUrl:        respImage.ImageUrl,
			CreatedAt:       respImage.CreatedAt.String(),
			UpdatedAt:       respImage.UpdatedAt.String(),
		}

		respImages = append(respImages, &image)
	}

	return &pb.Restaurant{
		RestaurantId:   response.RestaurantId,
		OwnerId:        response.OwnerId,
		RestaurantName: response.RestaurantName,
		Description:    response.Description,
		Rating:         response.Rating,
		OpeningHours:   response.OpeningHours,
		ContactNumber:  response.ContactNumber,
		LicenceUrl:     response.LicenceUrl,
		WebsiteUrl:     response.WebsiteUrl,
		Images:         respImages,
		Location: &pb.Location{
			LocationId:      response.Location.LocationId,
			EstablishmentId: response.Location.EstablishmentId,
			Address:         response.Location.Address,
			Latitude:        response.Location.Latitude,
			Longitude:       response.Location.Longitude,
			Country:         response.Location.Country,
			City:            response.Location.City,
			StateProvince:   response.Location.StateProvince,
			CreatedAt:       response.Location.CreatedAt.String(),
			UpdatedAt:       response.Location.UpdatedAt.String(),
		},
		CreatedAt: response.CreatedAt.String(),
		UpdatedAt: response.UpdatedAt.String(),
	}, nil
}

func (s establishmentRPC) GetRestaurant(ctx context.Context, request *pb.GetRestaurantRequest) (*pb.GetRestaurantResponse, error) {
	ctx, span := otlp.Start(ctx, "restaurant_grpc_delivery", "Get")
	span.SetAttributes(
		attribute.Key("restaurant_id").String(request.RestaurantId),
	)
	defer span.End()

	restaurant, err := s.restaurantUsecase.GetRestaurant(ctx, request.RestaurantId)
	if err != nil {
		return nil, err
	}

	var images []*pb.Image

	for _, i := range restaurant.Images {
		var image pb.Image

		image.ImageId = i.ImageId
		image.EstablishmentId = i.EstablishmentId
		image.ImageUrl = i.ImageUrl
		image.CreatedAt = i.CreatedAt.String()
		image.UpdatedAt = i.UpdatedAt.String()

		images = append(images, &image)
	}

	return &pb.GetRestaurantResponse{
		Restaurant: &pb.Restaurant{
			RestaurantId:   restaurant.RestaurantId,
			OwnerId:        restaurant.OwnerId,
			RestaurantName: restaurant.RestaurantName,
			Description:    restaurant.Description,
			Rating:         restaurant.Rating,
			OpeningHours:   restaurant.OpeningHours,
			ContactNumber:  restaurant.ContactNumber,
			LicenceUrl:     restaurant.LicenceUrl,
			WebsiteUrl:     restaurant.WebsiteUrl,
			Images:         images,
			Location: &pb.Location{
				LocationId:      restaurant.Location.LocationId,
				EstablishmentId: restaurant.Location.EstablishmentId,
				Address:         restaurant.Location.Address,
				Latitude:        restaurant.Location.Latitude,
				Longitude:       restaurant.Location.Longitude,
				Country:         restaurant.Location.Country,
				City:            restaurant.Location.City,
				StateProvince:   restaurant.Location.StateProvince,
				CreatedAt:       restaurant.CreatedAt.String(),
				UpdatedAt:       restaurant.UpdatedAt.String(),
			},
			CreatedAt: restaurant.CreatedAt.String(),
			UpdatedAt: restaurant.UpdatedAt.String(),
		},
	}, nil
}

func (s establishmentRPC) ListRestaurants(ctx context.Context, request *pb.ListRestaurantsRequest) (*pb.ListRestaurantsResponse, error) {
	ctx, span := otlp.Start(ctx, "restaurant_grpc_delivery", "List")
	span.SetAttributes(
		attribute.Key("limit").Int64(request.Limit),
	)
	defer span.End()

	restaurants, overall, err := s.restaurantUsecase.ListRestaurants(ctx, request.Offset, request.Limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch restaurants: %v", err)
	}

	var pbRestaurants []*pb.Restaurant
	for _, restaurant := range restaurants {
		var images []*pb.Image
		for _, i := range restaurant.Images {
			images = append(images, &pb.Image{
				ImageId:         i.ImageId,
				EstablishmentId: i.EstablishmentId,
				ImageUrl:        i.ImageUrl,
				CreatedAt:       i.CreatedAt.String(),
				UpdatedAt:       i.UpdatedAt.String(),
			})
		}

		pbRestaurants = append(pbRestaurants, &pb.Restaurant{
			RestaurantId:   restaurant.RestaurantId,
			OwnerId:        restaurant.OwnerId,
			RestaurantName: restaurant.RestaurantName,
			Description:    restaurant.Description,
			Rating:         restaurant.Rating,
			OpeningHours:   restaurant.OpeningHours,
			ContactNumber:  restaurant.ContactNumber,
			LicenceUrl:     restaurant.LicenceUrl,
			WebsiteUrl:     restaurant.WebsiteUrl,
			Images:         images,
			Location: &pb.Location{
				LocationId:      restaurant.Location.LocationId,
				EstablishmentId: restaurant.Location.EstablishmentId,
				Address:         restaurant.Location.Address,
				Latitude:        restaurant.Location.Latitude,
				Longitude:       restaurant.Location.Longitude,
				Country:         restaurant.Location.Country,
				City:            restaurant.Location.City,
				StateProvince:   restaurant.Location.StateProvince,
				CreatedAt:       restaurant.CreatedAt.String(),
				UpdatedAt:       restaurant.UpdatedAt.String(),
			},
			CreatedAt: restaurant.CreatedAt.String(),
			UpdatedAt: restaurant.UpdatedAt.String(),
		})
	}

	return &pb.ListRestaurantsResponse{
		Restaurants: pbRestaurants,
		Overall:     overall,
	}, nil
}

func (s establishmentRPC) UpdateRestaurant(ctx context.Context, request *pb.UpdateRestaurantRequest) (*pb.UpdateRestaurantResponse, error) {
	ctx, span := otlp.Start(ctx, "restaurant_grpc_delivery", "Update")
	span.SetAttributes(
		attribute.Key("restaurant_id").String(request.Restaurant.RestaurantId),
	)
	defer span.End()

	var imagesS []*entity.Image

	for _, i := range request.Restaurant.Images {
		var image entity.Image

		image.ImageId = i.ImageId
		image.EstablishmentId = i.EstablishmentId
		image.ImageUrl = i.ImageUrl

		imagesS = append(imagesS, &image)
	}

	restaurant, err := s.restaurantUsecase.UpdateRestaurant(ctx, &entity.Restaurant{
		RestaurantId:   request.Restaurant.RestaurantId,
		OwnerId:        request.Restaurant.OwnerId,
		RestaurantName: request.Restaurant.RestaurantName,
		Description:    request.Restaurant.Description,
		Rating:         request.Restaurant.Rating,
		OpeningHours:   request.Restaurant.OpeningHours,
		ContactNumber:  request.Restaurant.ContactNumber,
		LicenceUrl:     request.Restaurant.LicenceUrl,
		WebsiteUrl:     request.Restaurant.WebsiteUrl,
		Images:         imagesS,
		Location: entity.Location{
			LocationId:      request.Restaurant.Location.LocationId,
			EstablishmentId: request.Restaurant.Location.EstablishmentId,
			Address:         request.Restaurant.Location.Address,
			Latitude:        request.Restaurant.Location.Latitude,
			Longitude:       request.Restaurant.Location.Longitude,
			Country:         request.Restaurant.Location.Country,
			City:            request.Restaurant.Location.City,
			StateProvince:   request.Restaurant.Location.StateProvince,
		},
	})
	if err != nil {
		return nil, err
	}

	var images []*pb.Image

	for _, i := range restaurant.Images {
		var image pb.Image

		image.ImageId = i.ImageId
		image.EstablishmentId = i.EstablishmentId
		image.ImageUrl = i.ImageUrl
		image.CreatedAt = i.CreatedAt.String()
		image.UpdatedAt = i.UpdatedAt.String()

		images = append(images, &image)
	}

	return &pb.UpdateRestaurantResponse{
		Restaurant: &pb.Restaurant{
			RestaurantId:   restaurant.RestaurantId,
			OwnerId:        restaurant.OwnerId,
			RestaurantName: restaurant.RestaurantName,
			Description:    restaurant.Description,
			Rating:         restaurant.Rating,
			ContactNumber:  restaurant.ContactNumber,
			LicenceUrl:     restaurant.LicenceUrl,
			WebsiteUrl:     restaurant.WebsiteUrl,
			Images:         images,
			Location: &pb.Location{
				LocationId:      restaurant.Location.LocationId,
				EstablishmentId: restaurant.Location.EstablishmentId,
				Address:         restaurant.Location.Address,
				Latitude:        restaurant.Location.Latitude,
				Longitude:       restaurant.Location.Longitude,
				Country:         restaurant.Location.Country,
				City:            restaurant.Location.City,
				StateProvince:   restaurant.Location.StateProvince,
				CreatedAt:       restaurant.CreatedAt.String(),
				UpdatedAt:       restaurant.UpdatedAt.String(),
			},
			CreatedAt: restaurant.CreatedAt.String(),
			UpdatedAt: restaurant.UpdatedAt.String(),
		},
	}, nil
}

func (s establishmentRPC) DeleteRestaurant(ctx context.Context, request *pb.DeleteRestaurantRequest) (*pb.DeleteRestaurantResponse, error) {
	ctx, span := otlp.Start(ctx, "restaurant_grpc_delivery", "Delete")
	span.SetAttributes(
		attribute.Key("restaurant_id").String(request.RestaurantId),
	)
	defer span.End()

	if err := s.restaurantUsecase.DeleteRestaurant(ctx, request.RestaurantId); err != nil {
		return &pb.DeleteRestaurantResponse{
			Success: false,
		}, err
	}

	return &pb.DeleteRestaurantResponse{
		Success: true,
	}, nil
}

func (s establishmentRPC) ListRestaurantsByLocation(ctx context.Context, request *pb.ListRestaurantsByLocationRequest) (*pb.ListRestaurantsByLocationResponse, error) {
	ctx, span := otlp.Start(ctx, "restaurant_grpc_delivery", "List")
	span.SetAttributes(
		attribute.Key("state_province").String(request.StateProvince),
	)
	defer span.End()

	restaurants, count, err := s.restaurantUsecase.ListRestaurantsByLocation(ctx, request.Offset, request.Limit, request.Country, request.City, request.StateProvince)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch restaurants: %v", err)
	}

	var pbRestaurants []*pb.Restaurant
	for _, restaurant := range restaurants {
		var images []*pb.Image
		for _, i := range restaurant.Images {
			images = append(images, &pb.Image{
				ImageId:         i.ImageId,
				EstablishmentId: i.EstablishmentId,
				ImageUrl:        i.ImageUrl,
				CreatedAt:       i.CreatedAt.String(),
				UpdatedAt:       i.UpdatedAt.String(),
			})
		}

		pbRestaurants = append(pbRestaurants, &pb.Restaurant{
			RestaurantId:   restaurant.RestaurantId,
			OwnerId:        restaurant.OwnerId,
			RestaurantName: restaurant.RestaurantName,
			Description:    restaurant.Description,
			Rating:         restaurant.Rating,
			OpeningHours:   restaurant.OpeningHours,
			ContactNumber:  restaurant.ContactNumber,
			LicenceUrl:     restaurant.LicenceUrl,
			WebsiteUrl:     restaurant.WebsiteUrl,
			Images:         images,
			Location: &pb.Location{
				LocationId:      restaurant.Location.LocationId,
				EstablishmentId: restaurant.Location.EstablishmentId,
				Address:         restaurant.Location.Address,
				Latitude:        restaurant.Location.Latitude,
				Longitude:       restaurant.Location.Longitude,
				Country:         restaurant.Location.Country,
				City:            restaurant.Location.City,
				StateProvince:   restaurant.Location.StateProvince,
				CreatedAt:       restaurant.CreatedAt.String(),
				UpdatedAt:       restaurant.UpdatedAt.String(),
			},
			CreatedAt: restaurant.CreatedAt.String(),
			UpdatedAt: restaurant.UpdatedAt.String(),
		})
	}

	return &pb.ListRestaurantsByLocationResponse{
		Restaurants: pbRestaurants,
		Count:       count,
	}, nil
}

func (s establishmentRPC) FindRestaurantsByName(ctx context.Context, request *pb.FindRestaurantsByNameRequest) (*pb.FindRestaurantsByNameResponse, error) {
	ctx, span := otlp.Start(ctx, "restaurant_grpc_delivery", "Find")
	span.SetAttributes(
		attribute.Key("name").String(request.Name),
	)
	defer span.End()

	restaurants, overall, err := s.restaurantUsecase.FindRestaurantsByName(ctx, request.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch attractions: %v", err)
	}

	var pbRestaurants []*pb.Restaurant
	for _, restaurant := range restaurants {
		var images []*pb.Image
		for _, i := range restaurant.Images {
			images = append(images, &pb.Image{
				ImageId:         i.ImageId,
				EstablishmentId: i.EstablishmentId,
				ImageUrl:        i.ImageUrl,
				CreatedAt:       i.CreatedAt.String(),
				UpdatedAt:       i.UpdatedAt.String(),
			})
		}

		pbRestaurants = append(pbRestaurants, &pb.Restaurant{
			RestaurantId:   restaurant.RestaurantId,
			OwnerId:        restaurant.OwnerId,
			RestaurantName: restaurant.RestaurantName,
			Description:    restaurant.Description,
			Rating:         restaurant.Rating,
			OpeningHours:   restaurant.OpeningHours,
			ContactNumber:  restaurant.ContactNumber,
			LicenceUrl:     restaurant.LicenceUrl,
			WebsiteUrl:     restaurant.WebsiteUrl,
			Images:         images,
			Location: &pb.Location{
				LocationId:      restaurant.Location.LocationId,
				EstablishmentId: restaurant.Location.EstablishmentId,
				Address:         restaurant.Location.Address,
				Latitude:        restaurant.Location.Latitude,
				Longitude:       restaurant.Location.Longitude,
				Country:         restaurant.Location.Country,
				City:            restaurant.Location.City,
				StateProvince:   restaurant.Location.StateProvince,
				CreatedAt:       restaurant.CreatedAt.String(),
				UpdatedAt:       restaurant.UpdatedAt.String(),
			},
			CreatedAt: restaurant.CreatedAt.String(),
			UpdatedAt: restaurant.UpdatedAt.String(),
		})
	}

	return &pb.FindRestaurantsByNameResponse{
		Restaurants: pbRestaurants,
		Count:       overall,
	}, nil
}

// HOTEL
func (s establishmentRPC) CreateHotel(ctx context.Context, hotel *pb.Hotel) (*pb.Hotel, error) {
	ctx, span := otlp.Start(ctx, "hotel_grpc_delivery", "Create")
	span.SetAttributes(
		attribute.Key("hotel_id").String(hotel.HotelId),
	)
	defer span.End()

	var images []*entity.Image

	for _, i := range hotel.Images {
		var image entity.Image

		image.ImageId = i.ImageId
		image.EstablishmentId = i.EstablishmentId
		image.ImageUrl = i.ImageUrl
		image.Category = i.Category
		image.CreatedAt = time.Now().Local()
		image.UpdatedAt = time.Now().Local()

		images = append(images, &image)
	}
	response, err := s.hotelUsecase.CreateHotel(ctx, &entity.Hotel{
		HotelId:       hotel.HotelId,
		OwnerId:       hotel.OwnerId,
		HotelName:     hotel.HotelName,
		Description:   hotel.Description,
		Rating:        hotel.Rating,
		ContactNumber: hotel.ContactNumber,
		LicenceUrl:    hotel.LicenceUrl,
		WebsiteUrl:    hotel.WebsiteUrl,
		Images:        images,
		Location: entity.Location{
			LocationId:      hotel.Location.LocationId,
			EstablishmentId: hotel.Location.EstablishmentId,
			Address:         hotel.Location.Address,
			Latitude:        hotel.Location.Latitude,
			Longitude:       hotel.Location.Longitude,
			Country:         hotel.Location.Country,
			City:            hotel.Location.City,
			StateProvince:   hotel.Location.StateProvince,
			Category:        hotel.Location.Category,
			CreatedAt:       time.Now().Local(),
			UpdatedAt:       time.Now().Local(),
		},
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
	})
	if err != nil {
		return nil, err
	}

	var respImages []*pb.Image

	for _, respImage := range response.Images {
		image := pb.Image{
			ImageId:         respImage.ImageId,
			EstablishmentId: respImage.EstablishmentId,
			ImageUrl:        respImage.ImageUrl,
			CreatedAt:       respImage.CreatedAt.String(),
			UpdatedAt:       respImage.UpdatedAt.String(),
		}

		respImages = append(respImages, &image)
	}

	return &pb.Hotel{
		HotelId:       response.HotelId,
		OwnerId:       response.OwnerId,
		HotelName:     response.HotelName,
		Description:   response.Description,
		Rating:        response.Rating,
		ContactNumber: response.ContactNumber,
		LicenceUrl:    response.LicenceUrl,
		WebsiteUrl:    response.WebsiteUrl,
		Images:        respImages,
		Location: &pb.Location{
			LocationId:      response.Location.LocationId,
			EstablishmentId: response.Location.EstablishmentId,
			Address:         response.Location.Address,
			Latitude:        response.Location.Latitude,
			Longitude:       response.Location.Longitude,
			Country:         response.Location.Country,
			City:            response.Location.City,
			StateProvince:   response.Location.StateProvince,
			CreatedAt:       response.Location.CreatedAt.String(),
			UpdatedAt:       response.Location.UpdatedAt.String(),
		},
		CreatedAt: response.CreatedAt.String(),
		UpdatedAt: response.UpdatedAt.String(),
	}, nil
}

func (s establishmentRPC) GetHotel(ctx context.Context, request *pb.GetHotelRequest) (*pb.GetHotelResponse, error) {
	ctx, span := otlp.Start(ctx, "hotel_grpc_delivery", "Get")
	span.SetAttributes(
		attribute.Key("hotel_id").String(request.HotelId),
	)
	defer span.End()

	hotel, err := s.hotelUsecase.GetHotel(ctx, request.HotelId)
	if err != nil {
		return nil, err
	}

	var images []*pb.Image

	for _, i := range hotel.Images {
		var image pb.Image

		image.ImageId = i.ImageId
		image.EstablishmentId = i.EstablishmentId
		image.ImageUrl = i.ImageUrl
		image.CreatedAt = i.CreatedAt.String()
		image.UpdatedAt = i.UpdatedAt.String()

		images = append(images, &image)
	}

	return &pb.GetHotelResponse{
		Hotel: &pb.Hotel{
			HotelId:       hotel.HotelId,
			OwnerId:       hotel.OwnerId,
			HotelName:     hotel.HotelName,
			Description:   hotel.Description,
			Rating:        hotel.Rating,
			ContactNumber: hotel.ContactNumber,
			LicenceUrl:    hotel.LicenceUrl,
			WebsiteUrl:    hotel.WebsiteUrl,
			Images:        images,
			Location: &pb.Location{
				LocationId:      hotel.Location.LocationId,
				EstablishmentId: hotel.Location.EstablishmentId,
				Address:         hotel.Location.Address,
				Latitude:        hotel.Location.Latitude,
				Longitude:       hotel.Location.Longitude,
				Country:         hotel.Location.Country,
				City:            hotel.Location.City,
				StateProvince:   hotel.Location.StateProvince,
				CreatedAt:       hotel.CreatedAt.String(),
				UpdatedAt:       hotel.UpdatedAt.String(),
			},
			CreatedAt: hotel.CreatedAt.String(),
			UpdatedAt: hotel.UpdatedAt.String(),
		},
	}, nil
}

func (s establishmentRPC) ListHotels(ctx context.Context, request *pb.ListHotelsRequest) (*pb.ListHotelsResponse, error) {
	ctx, span := otlp.Start(ctx, "hotel_grpc_delivery", "Create")
	span.SetAttributes(
		attribute.Key("limit").Int64(request.Limit),
	)
	defer span.End()

	hotels, overall, err := s.hotelUsecase.ListHotels(ctx, request.Offset, request.Limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch hotels: %v", err)
	}

	// Convert []*entity.Hotel to []*pb.Hotel
	var pbHotels []*pb.Hotel
	for _, hotel := range hotels {
		var images []*pb.Image
		for _, i := range hotel.Images {
			images = append(images, &pb.Image{
				ImageId:         i.ImageId,
				EstablishmentId: i.EstablishmentId,
				ImageUrl:        i.ImageUrl,
				CreatedAt:       i.CreatedAt.String(),
				UpdatedAt:       i.UpdatedAt.String(),
			})
		}

		pbHotels = append(pbHotels, &pb.Hotel{
			HotelId:       hotel.HotelId,
			OwnerId:       hotel.OwnerId,
			HotelName:     hotel.HotelName,
			Description:   hotel.Description,
			Rating:        hotel.Rating,
			ContactNumber: hotel.ContactNumber,
			LicenceUrl:    hotel.LicenceUrl,
			WebsiteUrl:    hotel.WebsiteUrl,
			Images:        images,
			Location: &pb.Location{
				LocationId:      hotel.Location.LocationId,
				EstablishmentId: hotel.Location.EstablishmentId,
				Address:         hotel.Location.Address,
				Latitude:        hotel.Location.Latitude,
				Longitude:       hotel.Location.Longitude,
				Country:         hotel.Location.Country,
				City:            hotel.Location.City,
				StateProvince:   hotel.Location.StateProvince,
				CreatedAt:       hotel.CreatedAt.String(),
				UpdatedAt:       hotel.UpdatedAt.String(),
			},
			CreatedAt: hotel.CreatedAt.String(),
			UpdatedAt: hotel.UpdatedAt.String(),
		})
	}

	return &pb.ListHotelsResponse{
		Hotels:  pbHotels,
		Overall: overall,
	}, nil
}

func (s establishmentRPC) UpdateHotel(ctx context.Context, request *pb.UpdateHotelRequest) (*pb.UpdateHotelResponse, error) {
	ctx, span := otlp.Start(ctx, "hotel_grpc_delivery", "Update")
	span.SetAttributes(
		attribute.Key("hotel_id").String(request.Hotel.HotelId),
	)
	defer span.End()

	// var imagesS []*entity.Image

	// for _, i := range request.Hotel.Images {
	// 	var image entity.Image

	// 	image.ImageId = i.ImageId
	// 	image.EstablishmentId = i.EstablishmentId
	// 	image.ImageUrl = i.ImageUrl

	// 	imagesS = append(imagesS, &image)
	// }

	hotel, err := s.hotelUsecase.UpdateHotel(ctx, &entity.Hotel{
		HotelId:       request.Hotel.HotelId,
		OwnerId:       request.Hotel.OwnerId,
		HotelName:     request.Hotel.HotelName,
		Description:   request.Hotel.Description,
		Rating:        request.Hotel.Rating,
		ContactNumber: request.Hotel.ContactNumber,
		LicenceUrl:    request.Hotel.LicenceUrl,
		WebsiteUrl:    request.Hotel.WebsiteUrl,
		// Images:        imagesS,
		Location: entity.Location{
			LocationId:      request.Hotel.Location.LocationId,
			EstablishmentId: request.Hotel.Location.EstablishmentId,
			Address:         request.Hotel.Location.Address,
			Latitude:        request.Hotel.Location.Latitude,
			Longitude:       request.Hotel.Location.Longitude,
			Country:         request.Hotel.Location.Country,
			City:            request.Hotel.Location.City,
			StateProvince:   request.Hotel.Location.StateProvince,
		},
	})
	if err != nil {
		return nil, err
	}

	var images []*pb.Image

	for _, i := range hotel.Images {
		var image pb.Image

		image.ImageId = i.ImageId
		image.EstablishmentId = i.EstablishmentId
		image.ImageUrl = i.ImageUrl
		image.CreatedAt = i.CreatedAt.String()
		image.UpdatedAt = i.UpdatedAt.String()

		images = append(images, &image)
	}

	return &pb.UpdateHotelResponse{
		Hotel: &pb.Hotel{
			HotelId:       hotel.HotelId,
			OwnerId:       hotel.OwnerId,
			HotelName:     hotel.HotelName,
			Description:   hotel.Description,
			Rating:        hotel.Rating,
			ContactNumber: hotel.ContactNumber,
			LicenceUrl:    hotel.LicenceUrl,
			WebsiteUrl:    hotel.WebsiteUrl,
			Images:        images,
			Location: &pb.Location{
				LocationId:      hotel.Location.LocationId,
				EstablishmentId: hotel.Location.EstablishmentId,
				Address:         hotel.Location.Address,
				Latitude:        hotel.Location.Latitude,
				Longitude:       hotel.Location.Longitude,
				Country:         hotel.Location.Country,
				City:            hotel.Location.City,
				StateProvince:   hotel.Location.StateProvince,
				CreatedAt:       hotel.CreatedAt.String(),
				UpdatedAt:       hotel.UpdatedAt.String(),
			},
			CreatedAt: hotel.CreatedAt.String(),
			UpdatedAt: hotel.UpdatedAt.String(),
		},
	}, nil
}

func (s establishmentRPC) DeleteHotel(ctx context.Context, request *pb.DeleteHotelRequest) (*pb.DeleteHotelResponse, error) {
	ctx, span := otlp.Start(ctx, "hotel_grpc_delivery", "Delete")
	span.SetAttributes(
		attribute.Key("hotel_id").String(request.HotelId),
	)
	defer span.End()

	if err := s.hotelUsecase.DeleteHotel(ctx, request.HotelId); err != nil {
		return &pb.DeleteHotelResponse{
			Success: false,
		}, err
	}

	return &pb.DeleteHotelResponse{
		Success: true,
	}, nil
}

func (s establishmentRPC) ListHotelsByLocation(ctx context.Context, request *pb.ListHotelsByLocationRequest) (*pb.ListHotelsByLocationResponse, error) {
	ctx, span := otlp.Start(ctx, "hotel_grpc_delivery", "List")
	span.SetAttributes(
		attribute.Key("state_province").String(request.StateProvince),
	)
	defer span.End()

	hotels, count, err := s.hotelUsecase.ListHotelsByLocation(ctx, request.Offset, request.Limit, request.Country, request.City, request.StateProvince)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch hotels: %v", err)
	}

	var pbHotels []*pb.Hotel
	for _, hotel := range hotels {
		var images []*pb.Image
		for _, i := range hotel.Images {
			images = append(images, &pb.Image{
				ImageId:         i.ImageId,
				EstablishmentId: i.EstablishmentId,
				ImageUrl:        i.ImageUrl,
				CreatedAt:       i.CreatedAt.String(),
				UpdatedAt:       i.UpdatedAt.String(),
			})
		}

		pbHotels = append(pbHotels, &pb.Hotel{
			HotelId:       hotel.HotelId,
			OwnerId:       hotel.OwnerId,
			HotelName:     hotel.HotelName,
			Description:   hotel.Description,
			Rating:        hotel.Rating,
			ContactNumber: hotel.ContactNumber,
			LicenceUrl:    hotel.LicenceUrl,
			WebsiteUrl:    hotel.WebsiteUrl,
			Images:        images,
			Location: &pb.Location{
				LocationId:      hotel.Location.LocationId,
				EstablishmentId: hotel.Location.EstablishmentId,
				Address:         hotel.Location.Address,
				Latitude:        hotel.Location.Latitude,
				Longitude:       hotel.Location.Longitude,
				Country:         hotel.Location.Country,
				City:            hotel.Location.City,
				StateProvince:   hotel.Location.StateProvince,
				CreatedAt:       hotel.CreatedAt.String(),
				UpdatedAt:       hotel.UpdatedAt.String(),
			},
			CreatedAt: hotel.CreatedAt.String(),
			UpdatedAt: hotel.UpdatedAt.String(),
		})
	}

	return &pb.ListHotelsByLocationResponse{
		Hotels: pbHotels,
		Count:  uint64(count),
	}, nil
}

func (s establishmentRPC) FindHotelsByName(ctx context.Context, request *pb.FindHotelsByNameRequest) (*pb.FindHotelsByNameResponse, error) {
	ctx, span := otlp.Start(ctx, "hotel_grpc_delivery", "Find")
	span.SetAttributes(
		attribute.Key("name").String(request.Name),
	)
	defer span.End()

	hotels, overall, err := s.hotelUsecase.FindHotelsByName(ctx, request.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch hotels: %v", err)
	}

	var pbHotels []*pb.Hotel
	for _, hotel := range hotels {
		var images []*pb.Image
		for _, i := range hotel.Images {
			images = append(images, &pb.Image{
				ImageId:         i.ImageId,
				EstablishmentId: i.EstablishmentId,
				ImageUrl:        i.ImageUrl,
				CreatedAt:       i.CreatedAt.String(),
				UpdatedAt:       i.UpdatedAt.String(),
			})
		}

		pbHotels = append(pbHotels, &pb.Hotel{
			HotelId:       hotel.HotelId,
			OwnerId:       hotel.OwnerId,
			HotelName:     hotel.HotelName,
			Description:   hotel.Description,
			Rating:        hotel.Rating,
			ContactNumber: hotel.ContactNumber,
			LicenceUrl:    hotel.LicenceUrl,
			WebsiteUrl:    hotel.WebsiteUrl,
			Images:        images,
			Location: &pb.Location{
				LocationId:      hotel.Location.LocationId,
				EstablishmentId: hotel.Location.EstablishmentId,
				Address:         hotel.Location.Address,
				Latitude:        hotel.Location.Latitude,
				Longitude:       hotel.Location.Longitude,
				Country:         hotel.Location.Country,
				City:            hotel.Location.City,
				StateProvince:   hotel.Location.StateProvince,
				CreatedAt:       hotel.CreatedAt.String(),
				UpdatedAt:       hotel.UpdatedAt.String(),
			},
			CreatedAt: hotel.CreatedAt.String(),
			UpdatedAt: hotel.UpdatedAt.String(),
		})
	}

	return &pb.FindHotelsByNameResponse{
		Hotels: pbHotels,
		Count:  overall,
	}, nil
}

// FAVOURITE
func (s establishmentRPC) AddToFavourites(ctx context.Context, request *pb.AddToFavouritesRequest) (*pb.AddToFavouritesResponse, error) {
	ctx, span := otlp.Start(ctx, "favourite_grpc_delivery", "Create")
	span.SetAttributes(
		attribute.Key("favourite_id").String(request.Favourite.FavouriteId),
	)
	defer span.End()

	response, err := s.favouriteUsecase.AddToFavourites(ctx, &entity.Favourite{
		FavouriteId:     request.Favourite.FavouriteId,
		EstablishmentId: request.Favourite.EstablishmentId,
		UserId:          request.Favourite.UserId,
		CreatedAt:       time.Now().Local(),
		UpdatedAt:       time.Now().Local(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.AddToFavouritesResponse{
		Favourite: &pb.Favourite{
			FavouriteId:     response.FavouriteId,
			EstablishmentId: response.EstablishmentId,
			UserId:          response.UserId,
			CreatedAt:       response.CreatedAt.String(),
			UpdatedAt:       response.UpdatedAt.String(),
		},
	}, nil
}

func (s establishmentRPC) RemoveFromFavourites(ctx context.Context, request *pb.RemoveFromFavouritesRequest) (*pb.RemoveFromFavouritesResponse, error) {
	ctx, span := otlp.Start(ctx, "favourite_grpc_delivery", "Delete")
	span.SetAttributes(
		attribute.Key("favourite_id").String(request.FavouriteId),
	)
	defer span.End()

	if err := s.favouriteUsecase.RemoveFromFavourites(ctx, request.FavouriteId); err != nil {
		return &pb.RemoveFromFavouritesResponse{
			Success: false,
		}, nil
	}

	return &pb.RemoveFromFavouritesResponse{
		Success: true,
	}, nil
}

func (s establishmentRPC) ListFavouritesByUserId(ctx context.Context, request *pb.ListFavouritesByUserIdRequest) (*pb.ListFavouritesByUserIdResponse, error) {
	ctx, span := otlp.Start(ctx, "favourite_grpc_delivery", "List")
	span.SetAttributes(
		attribute.Key("user_id").String(request.UserId),
	)
	defer span.End()

	response, err := s.favouriteUsecase.ListFavouritesByUserId(ctx, request.UserId)
	if err != nil {
		return nil, err
	}

	var favourites []*pb.Favourite

	for _, respFavourite := range response {
		favourite := pb.Favourite{
			FavouriteId:     respFavourite.FavouriteId,
			EstablishmentId: respFavourite.EstablishmentId,
			UserId:          respFavourite.UserId,
			CreatedAt:       respFavourite.CreatedAt.String(),
			UpdatedAt:       respFavourite.UpdatedAt.String(),
		}

		favourites = append(favourites, &favourite)
	}

	return &pb.ListFavouritesByUserIdResponse{
		Favourites: favourites,
	}, nil
}

// REVIEW
func (s establishmentRPC) CreateReview(ctx context.Context, request *pb.CreateReviewRequest) (*pb.CreateReviewResponse, error) {
	ctx, span := otlp.Start(ctx, "review_grpc_delivery", "Create")
	span.SetAttributes(
		attribute.Key("review_id").String(request.Review.ReviewId),
	)
	defer span.End()

	response, err := s.reviewUsecase.CreateReview(ctx, &entity.Review{
		ReviewId:        request.Review.ReviewId,
		EstablishmentId: request.Review.EstablishmentId,
		UserId:          request.Review.UserId,
		Rating:          float64(request.Review.Rating),
		Comment:         request.Review.Comment,
		CreatedAt:       time.Now().Local(),
		UpdatedAt:       time.Now().Local(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateReviewResponse{
		Review: &pb.Review{
			ReviewId:        response.ReviewId,
			EstablishmentId: response.EstablishmentId,
			UserId:          response.UserId,
			Rating:          float32(response.Rating),
			Comment:         response.Comment,
			CreatedAt:       response.CreatedAt.String(),
			UpdatedAt:       response.UpdatedAt.String(),
		},
	}, nil
}

func (s establishmentRPC) ListReviews(ctx context.Context, request *pb.ListReviewsRequest) (*pb.ListReviewsResponse, error) {
	ctx, span := otlp.Start(ctx, "review_grpc_delivery", "List")
	span.SetAttributes(
		attribute.Key("establishment_id").String(request.EstablishmentId),
	)
	defer span.End()

	response, count, err := s.reviewUsecase.ListReviews(ctx, request.EstablishmentId)
	if err != nil {
		return nil, err
	}

	var reviews []*pb.Review

	for _, respReview := range response {
		review := pb.Review{
			ReviewId:        respReview.ReviewId,
			EstablishmentId: respReview.EstablishmentId,
			UserId:          respReview.UserId,
			Rating:          float32(respReview.Rating),
			Comment:         respReview.Comment,
			CreatedAt:       respReview.CreatedAt.String(),
			UpdatedAt:       respReview.UpdatedAt.String(),
		}

		reviews = append(reviews, &review)
	}

	return &pb.ListReviewsResponse{
		Reviews: reviews,
		Count:   count,
	}, nil
}

func (s establishmentRPC) DeleteReview(ctx context.Context, request *pb.DeleteReviewRequest) (*pb.DeleteReviewResponse, error) {
	ctx, span := otlp.Start(ctx, "review_grpc_delivery", "Delete")
	span.SetAttributes(
		attribute.Key("review_id").String(request.ReviewId),
	)
	defer span.End()

	if err := s.reviewUsecase.DeleteReview(ctx, request.ReviewId); err != nil {
		return &pb.DeleteReviewResponse{
			Success: false,
		}, err
	}

	return &pb.DeleteReviewResponse{
		Success: true,
	}, nil
}

// MEDIA
func (s establishmentRPC) CreateMedia(ctx context.Context, image *pb.Image) (*pb.CreateImageRes, error) {
	ctx, span := otlp.Start(ctx, "media_grpc_delivery", "Create")
	span.SetAttributes(
		attribute.Key("establishment_id").String(image.EstablishmentId),
	)
	defer span.End()

	err := s.imageUsecase.CreateImage(ctx, &entity.Image{
		ImageId:         image.ImageId,
		EstablishmentId: image.EstablishmentId,
		ImageUrl:        image.ImageUrl,
		Category:        image.Category,
		CreatedAt:       time.Now().Local(),
		UpdatedAt:       time.Now().Local(),
		DeletedAt:       time.Now().Local(),
	})

	if err != nil {
		return nil, err
	}

	return &pb.CreateImageRes{
		Result: "Image has been created",
	}, nil
}
