package service

import (
	"context"
	"testing"

	"github.com/seyuta/ecommerce-microservices-example/s-catalog/pkg/model"
	"github.com/seyuta/ecommerce-microservices-example/s-catalog/pkg/pb"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type ProductRepositoryMock struct {
	mock.Mock
}

func (r ProductRepositoryMock) Create(context.Context, *model.Product) (model.Product, error) {
	args := r.Called()
	products := model.Product{
		Name: "kopi",
	}
	return products, args.Error(1)
}

func (r ProductRepositoryMock) FindByID(string) (model.Product, error) {
	args := r.Called()
	products := model.Product{
		Name: "kopi",
	}
	return products, args.Error(1)
}

func (r ProductRepositoryMock) FindAll() ([]model.Product, error) {
	args := r.Called()
	products := []model.Product{
		{Name: "kopi"},
		{Name: "teh"},
	}
	return products, args.Error(1)
}

func TestService_Create(t *testing.T) {
	repository := ProductRepositoryMock{}
	repository.On("Create").Return(model.Product{}, nil)

	ctx := context.Background()

	service := ProductSvc{repository, logrus.StandardLogger()}
	req := &pb.ProductDto{
		Name:  "kopi",
		Price: 123,
		Stock: 2,
	}
	product, _ := service.Create(ctx, req)
	assert.Equal(t, product.Name, "kopi")
}

func TestService_GetProductByID(t *testing.T) {
	repository := ProductRepositoryMock{}
	repository.On("FindByID").Return(model.Product{}, nil)

	ctx := context.Background()

	service := ProductSvc{repository, logrus.StandardLogger()}
	product, _ := service.GetProductByID(ctx, &wrapperspb.StringValue{})
	assert.Equal(t, product.Name, "kopi")
}

func TestService_ListProduct(t *testing.T) {
	repository := ProductRepositoryMock{}
	repository.On("FindAll").Return(model.Product{}, nil)

	ctx := context.Background()

	service := ProductSvc{repository, logrus.StandardLogger()}
	product, _ := service.ListProduct(ctx, &emptypb.Empty{})
	assert.NotEmpty(t, product)
}
