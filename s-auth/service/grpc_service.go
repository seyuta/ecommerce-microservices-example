package service

import (
	"github.com/seyuta/ecommerce-microservices-example/s-auth/pkg/pb"
	"github.com/seyuta/ecommerce-microservices-example/s-auth/repository"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// GRPCService ...
type GRPCService struct {
	server  *grpc.Server
	AuthSvc *AuthSvc
}

// BuildGRPCService register gRPC services implementation
func BuildGRPCService(server *grpc.Server, repository repository.Repository) *GRPCService {
	logger := logrus.StandardLogger()

	authSvc := NewAuthSvc(repository.UserRepository(), logger)
	pb.RegisterAuthApiServer(server, authSvc)

	return &GRPCService{
		server:  server,
		AuthSvc: authSvc,
	}
}
