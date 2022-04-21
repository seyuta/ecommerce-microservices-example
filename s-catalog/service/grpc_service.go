package service

import (
	"github.com/seyuta/ecommerce-microservices-example/s-catalog/pkg/pb"
	"github.com/seyuta/ecommerce-microservices-example/s-catalog/repository"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// GRPCService ...
type GRPCService struct {
	server     *grpc.Server
	ProductSvc *ProductSvc
}

// BuildGRPCService register gRPC services implementation
func BuildGRPCService(server *grpc.Server, repository repository.Repository) *GRPCService {
	logger := logrus.StandardLogger()

	productSvc := NewProductSvc(repository.ProductRepository(), logger)
	pb.RegisterProductApiServer(server, productSvc)

	return &GRPCService{
		server:     server,
		ProductSvc: productSvc,
	}
}
