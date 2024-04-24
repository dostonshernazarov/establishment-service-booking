package grpc_service_clients

import (
	"Booking/establishment-service-booking/internal/pkg/config"

	"google.golang.org/grpc"
)

type ServiceClients interface {
	Close()
}

type serviceClients struct {
	services []*grpc.ClientConn
}

func New(config *config.Config) (ServiceClients, error) {
	return &serviceClients{
		services: []*grpc.ClientConn{},
	}, nil
}

func (s *serviceClients) Close() {
	for _, conn := range s.services {
		conn.Close()
	}
}
