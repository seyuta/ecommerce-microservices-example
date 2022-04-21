package service

import (
	"context"
	"testing"

	"github.com/seyuta/ecommerce-microservices-example/s-transactions/pkg/model"
	"github.com/seyuta/ecommerce-microservices-example/s-transactions/pkg/pb"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type OrderRepositoryMock struct {
	mock.Mock
}

func (r OrderRepositoryMock) Create(context.Context, *model.Order) (model.Order, error) {
	args := r.Called()
	order := model.Order{
		NoInvoice: "INV/123/123",
	}
	return order, args.Error(1)
}

func (r OrderRepositoryMock) FindByID(string) (model.Order, error) {
	args := r.Called()
	order := model.Order{
		NoInvoice: "INV/123/123",
	}
	return order, args.Error(1)
}

func (r OrderRepositoryMock) FindAll() ([]model.Order, error) {
	args := r.Called()
	order := []model.Order{
		{NoInvoice: "INV/123/123"},
		{NoInvoice: "INV/456/456"},
	}
	return order, args.Error(1)
}

func TestService_Create(t *testing.T) {
	repository := OrderRepositoryMock{}
	repository.On("Create").Return(model.Order{}, nil)

	ctx := context.Background()

	service := OrderSvc{repository, logrus.StandardLogger()}
	req := &pb.OrderReqDto{}
	order, _ := service.Create(ctx, req)
	assert.Equal(t, order.NoInv, "INV/123/123")
}

func TestService_GetOrderByID(t *testing.T) {
	repository := OrderRepositoryMock{}
	repository.On("FindByID").Return(model.Order{}, nil)

	ctx := context.Background()

	service := OrderSvc{repository, logrus.StandardLogger()}
	order, _ := service.GetOrderByID(ctx, &wrapperspb.StringValue{})
	assert.Equal(t, order.NoInv, "INV/123/123")
}

func TestService_ListOrder(t *testing.T) {
	repository := OrderRepositoryMock{}
	repository.On("FindAll").Return(model.Order{}, nil)

	ctx := context.Background()

	service := OrderSvc{repository, logrus.StandardLogger()}
	order, _ := service.ListOrderByUserID(ctx, &wrapperspb.StringValue{})
	assert.NotEmpty(t, order, "get order success")
}
