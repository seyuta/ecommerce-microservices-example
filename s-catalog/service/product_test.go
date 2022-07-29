package service

import (
	"context"
	"testing"

	"github.com/seyuta/ecommerce-microservices-example/s-catalog/pkg/model"
	"github.com/seyuta/ecommerce-microservices-example/s-catalog/pkg/pb"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type ProductRepositoryMock struct {
	mock.Mock
}

func (r *ProductRepositoryMock) Create(ctx context.Context, req *model.Product) (model.Product, error) {
	args := r.Called(ctx, req)
	return args.Get(0).(model.Product), args.Error(1)
}

func (r *ProductRepositoryMock) FindByID(id string) (model.Product, error) {
	args := r.Called(id)
	return args.Get(0).(model.Product), args.Error(1)
}

func (r *ProductRepositoryMock) FindAll() ([]model.Product, error) {
	args := r.Called()
	return args.Get(0).([]model.Product), args.Error(1)
}

var (
	// context
	ctx = context.Background()

	// mock
	productRepoMock = ProductRepositoryMock{}

	// service
	productService = NewProductSvc(&productRepoMock, logrus.StandardLogger())
)

func TestService_Create(t *testing.T) {
	t.Run("Scenario Normal", func(t *testing.T) {
		productRepoMock.On("Create", ctx, mock.Anything).Once().Return(model.Product{}, nil)

		req := &pb.ProductDto{
			Name:  "coffee",
			Price: 123,
			Stock: 2,
		}
		_, err := productService.Create(ctx, req)
		assert.Empty(t, err)
	})

	t.Run("Scenario Blocker Name Empty", func(t *testing.T) {
		req := &pb.ProductDto{Name: ""}
		_, err := productService.Create(ctx, req)
		assert.NotEmpty(t, err)
	})
}

func TestService_GetProductByID(t *testing.T) {
	t.Run("Scenario Normal", func(t *testing.T) {
		productRepoMock.On("FindByID", mock.AnythingOfType("string")).Once().Return(model.Product{}, nil)
		_, err := productService.GetProductByID(ctx, &wrapperspb.StringValue{})
		assert.Empty(t, err)
	})

	t.Run("Scenario Product Not Found", func(t *testing.T) {
		productRepoMock.On("FindByID", mock.AnythingOfType("string")).Once().Return(model.Product{}, mongo.ErrNoDocuments)
		_, err := productService.GetProductByID(ctx, &wrapperspb.StringValue{})
		assert.NotEmpty(t, err)
	})
}

func TestService_ListProduct(t *testing.T) {
	t.Run("Scenario Normal", func(t *testing.T) {
		productRepoMock.On("FindAll").Once().Return([]model.Product{}, nil)
		_, err := productService.ListProduct(ctx, &emptypb.Empty{})
		assert.Empty(t, err)
	})

	t.Run("Scenario Product Not Found", func(t *testing.T) {
		productRepoMock.On("FindAll").Once().Return([]model.Product{}, mongo.ErrNoDocuments)
		_, err := productService.ListProduct(ctx, &emptypb.Empty{})
		assert.NotEmpty(t, err)
	})
}
