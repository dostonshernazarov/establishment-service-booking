package app

import (
	pb "Booking/establishment-service-booking/genproto/establishment-proto"
	grpc_server "Booking/establishment-service-booking/internal/delivery/grpc/server"
	invest_grpc "Booking/establishment-service-booking/internal/delivery/grpc/services"
	"Booking/establishment-service-booking/internal/infrastructure/grpc_service_clients"
	"Booking/establishment-service-booking/internal/infrastructure/kafka"
	repo "Booking/establishment-service-booking/internal/infrastructure/repository/postgresql"
	"Booking/establishment-service-booking/internal/pkg/config"
	"Booking/establishment-service-booking/internal/pkg/logger"
	"Booking/establishment-service-booking/internal/pkg/postgres"
	"Booking/establishment-service-booking/internal/usecase"
	"Booking/establishment-service-booking/internal/usecase/event"
	"fmt"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	Config            *config.Config
	Logger            *zap.Logger
	DB                *postgres.PostgresDB
	GrpcServer        *grpc.Server
	AttractionUsecase usecase.Attraction
	RestaurantUsecase usecase.Restaurant
	FavouriteUsecase  usecase.Favourite
	Review            usecase.Review
	Image             usecase.Image
	ServiceClients    grpc_service_clients.ServiceClients
	BrokerProducer    event.BrokerProducer
}

func NewApp(cfg *config.Config) (*App, error) {
	// init logger
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}

	kafkaProducer := kafka.NewProducer(cfg, logger)

	// init db
	db, err := postgres.New(cfg)
	if err != nil {
		return nil, err
	}

	// grpc server init
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_server.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_zap.UnaryServerInterceptor(logger),
				grpc_recovery.UnaryServerInterceptor(),
			),
			grpc_server.UnaryInterceptorData(logger),
		)),
	)

	return &App{
		Config:         cfg,
		Logger:         logger,
		DB:             db,
		GrpcServer:     grpcServer,
		BrokerProducer: kafkaProducer,
	}, nil
}

func (a *App) Run() error {
	var (
		contextTimeout time.Duration
	)

	// context timeout initialization
	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error during parse duration for context timeout : %w", err)
	}
	// Initialize Service Clients
	serviceClients, err := grpc_service_clients.New(a.Config)
	if err != nil {
		return fmt.Errorf("error during initialize service clients: %w", err)
	}
	a.ServiceClients = serviceClients

	// repositories initialization
	attractionRepo := repo.NewAttractionRepo(a.DB)
	restaurantRepo := repo.NewRestaurantRepo(a.DB)
	hotelRepo := repo.NewHotelRepo(a.DB)
	favouriteRepo := repo.NewFavouriteRepo(a.DB)
	reviewRepo := repo.NewReviewRepo(a.DB)
	imageRepo := repo.NewImageRepo(a.DB)

	// usecase initialization
	attracationUsecase := usecase.NewAttractionService(contextTimeout, attractionRepo)
	restaurantUsecase := usecase.NewRestaurantService(contextTimeout, restaurantRepo)
	hotelUsecase := usecase.NewHotelService(contextTimeout, hotelRepo)
	favouriteUsecase := usecase.NewFavouriteService(contextTimeout, favouriteRepo)
	reviewUsecase := usecase.NewReviewService(contextTimeout, reviewRepo)
	imageUsecase := usecase.NewImageService(contextTimeout, imageRepo)

	pb.RegisterEstablishmentServiceServer(a.GrpcServer, invest_grpc.NewRPC(a.Logger, attracationUsecase, restaurantUsecase, hotelUsecase, favouriteUsecase,imageUsecase, reviewUsecase, a.BrokerProducer))
	a.Logger.Info("gRPC Server Listening", zap.String("url", a.Config.RPCPort))
	if err := grpc_server.Run(a.Config, a.GrpcServer); err != nil {
		return fmt.Errorf("gRPC fatal to serve grpc server over %s %w", a.Config.RPCPort, err)
	}
	return nil
}

func (a *App) Stop() {
	// close broker producer
	a.BrokerProducer.Close()

	// closing client service connections
	a.ServiceClients.Close()

	// stop gRPC server
	a.GrpcServer.Stop()

	// database connection
	a.DB.Close()

	// zap logger sync
	a.Logger.Sync()
}
