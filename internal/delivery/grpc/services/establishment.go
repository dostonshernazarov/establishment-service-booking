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
	brokerProducer     event.BrokerProducer
}

func NewRPC(logger *zap.Logger, attracationUsecase usecase.Attraction, restaurantUsecase usecase.Restaurant, hotelUsecase usecase.Hotel, brokerProducer event.BrokerProducer) pb.EstablishmentServiceServer {
	return &establishmentRPC{
		logger:             logger,
		attracationUsecase: attracationUsecase,
		restaurantUsecase:  restaurantUsecase,
		hotelUsecase:       hotelUsecase,
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

func (s establishmentRPC) GetAttraction(ctx context.Context, req *pb.GetAttractionRequest) (*pb.GetAttractionResponse, error) {
	attraction, err := s.attracationUsecase.GetAttraction(ctx, req.AttractionId)

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

func (s establishmentRPC) ListAttractions(ctx context.Context, req *pb.ListAttractionsRequest) (*pb.ListAttractionsResponse, error) {
	attractions, err := s.attracationUsecase.ListAttractions(ctx, req.Offset, req.Limit)
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
	}, nil
}

func (s establishmentRPC) UpdateAttraction(ctx context.Context, request *pb.UpdateAttractionRequest) (*pb.UpdateAttractionResponse, error) {
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

func (s establishmentRPC) DeleteAttraction(ctx context.Context, req *pb.DeleteAttractionRequest) (*pb.DeleteAttractionResponse, error) {
	if err := s.attracationUsecase.DeleteAttraction(ctx, req.AttractionId); err != nil {
		return &pb.DeleteAttractionResponse{
			Success: false,
		}, err
	}

	return &pb.DeleteAttractionResponse{
		Success: true,
	}, nil
}

// RESTAURANT
func (s establishmentRPC) CreateRestaurant(ctx context.Context, restaurant *pb.Restaurant) (*pb.Restaurant, error) {
	var images []*entity.Image

	for _, i := range restaurant.Images {
		var image entity.Image

		image.ImageId = i.ImageId
		image.EstablishmentId = i.EstablishmentId
		image.ImageUrl = i.ImageUrl

		images = append(images, &image)
	}
	_, err := s.restaurantUsecase.CreateRestaurant(ctx, &entity.Restaurant{
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
		},
	})
	if err != nil {
		return nil, err
	}

	return restaurant, nil
}

func (s establishmentRPC) GetRestaurant(ctx context.Context, request *pb.GetRestaurantRequest) (*pb.GetRestaurantResponse, error) {
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

func (s establishmentRPC) ListRestaurants(ctx context.Context, req *pb.ListRestaurantsRequest) (*pb.ListRestaurantsResponse, error) {
	restaurants, err := s.restaurantUsecase.ListRestaurants(ctx, req.Offset, req.Limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch restaurants: %v", err)
	}

	// Convert []*entity.Restaurant to []*pb.Restaurant
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
	}, nil
}

func (s establishmentRPC) UpdateRestaurant(ctx context.Context, request *pb.UpdateRestaurantRequest) (*pb.UpdateRestaurantResponse, error) {
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

func (s establishmentRPC) DeleteRestaurant(ctx context.Context, req *pb.DeleteRestaurantRequest) (*pb.DeleteRestaurantResponse, error) {
	if err := s.restaurantUsecase.DeleteRestaurant(ctx, req.RestaurantId); err != nil {
		return &pb.DeleteRestaurantResponse{
			Success: false,
		}, err
	}

	return &pb.DeleteRestaurantResponse{
		Success: true,
	}, nil
}

// HOTEL
func (s establishmentRPC) CreateHotel(ctx context.Context, hotel *pb.Hotel) (*pb.Hotel, error) {
	var images []*entity.Image

	for _, i := range hotel.Images {
		var image entity.Image

		image.ImageId = i.ImageId
		image.EstablishmentId = i.EstablishmentId
		image.ImageUrl = i.ImageUrl
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

func (s establishmentRPC) ListHotels(ctx context.Context, req *pb.ListHotelsRequest) (*pb.ListHotelsResponse, error) {
	hotels, err := s.hotelUsecase.ListHotels(ctx, req.Offset, req.Limit)
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
		Hotels: pbHotels,
	}, nil
}

func (s establishmentRPC) UpdateHotel(ctx context.Context, request *pb.UpdateHotelRequest) (*pb.UpdateHotelResponse, error) {
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
	if err := s.hotelUsecase.DeleteHotel(ctx, request.HotelId); err != nil {
		return &pb.DeleteHotelResponse{
			Success: false,
		}, err
	}

	return &pb.DeleteHotelResponse{
		Success: true,
	}, nil
}
