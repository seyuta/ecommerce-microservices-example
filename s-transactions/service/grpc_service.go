package service

import (
	"github.com/seyuta/ecommerce-microservices-example/s-transactions/pkg/pb"
	"github.com/seyuta/ecommerce-microservices-example/s-transactions/repository"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// GRPCService ...
type GRPCService struct {
	server   *grpc.Server
	OrderSvc *OrderSvc
}

// BuildGRPCService register gRPC services implementation
func BuildGRPCService(server *grpc.Server, repository repository.Repository) *GRPCService {
	logger := logrus.StandardLogger()

	orderSvc := NewOrderSvc(repository.OrderRepository(), logger)
	pb.RegisterOrderApiServer(server, orderSvc)

	return &GRPCService{
		server:   server,
		OrderSvc: orderSvc,
	}
}
