package service

import (
	"context"
	"testing"

	"github.com/seyuta/ecommerce-microservices-example/s-transactions/pkg/model"
	"github.com/seyuta/ecommerce-microservices-example/s-transactions/pkg/pb"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type OrderRepositoryMock struct {
	mock.Mock
}

func (r *OrderRepositoryMock) Create(ctx context.Context, req *model.Order) (model.Order, error) {
	args := r.Called(ctx, req)
	return args.Get(0).(model.Order), args.Error(1)
}

func (r *OrderRepositoryMock) FindByID(id string) (model.Order, error) {
	args := r.Called(id)
	return args.Get(0).(model.Order), args.Error(1)
}

func (r *OrderRepositoryMock) FindAll() ([]model.Order, error) {
	args := r.Called()
	return args.Get(0).([]model.Order), args.Error(1)
}

var (
	// context
	ctx = context.Background()

	// mock
	orderRepoMock = OrderRepositoryMock{}

	// service
	orderService = NewOrderSvc(&orderRepoMock, logrus.StandardLogger())
)

func TestService_Create(t *testing.T) {
	t.Run("Scenario Normal", func(t *testing.T) {
		orderRepoMock.On("Create", ctx, mock.Anything).Once().Return(model.Order{}, nil)

		req := &pb.OrderReqDto{
			Order: []*pb.OrderDetailDto{
				{
					ProductId: "1",
					Qty:       1,
				},
				{
					ProductId: "2",
					Qty:       1,
				},
			},
		}
		_, err := orderService.Create(ctx, req)
		assert.Empty(t, err)
	})
}

func TestService_GetOrderByID(t *testing.T) {
	t.Run("Scenario Normal", func(t *testing.T) {
		orderRepoMock.On("FindByID", mock.AnythingOfType("string")).Once().Return(model.Order{}, nil)
		_, err := orderService.GetOrderByID(ctx, &wrapperspb.StringValue{})
		assert.Empty(t, err)
	})

	t.Run("Scenario Order Not Found", func(t *testing.T) {
		orderRepoMock.On("FindByID", mock.AnythingOfType("string")).Once().Return(model.Order{}, mongo.ErrNoDocuments)
		_, err := orderService.GetOrderByID(ctx, &wrapperspb.StringValue{})
		assert.NotEmpty(t, err)
	})
}

func TestService_ListOrder(t *testing.T) {
	t.Run("Scenario Normal", func(t *testing.T) {
		orderRepoMock.On("FindAll").Once().Return([]model.Order{}, nil)
		_, err := orderService.ListOrderByUserID(ctx, &wrapperspb.StringValue{})
		assert.Empty(t, err)
	})

	t.Run("Scenario Order Not Found", func(t *testing.T) {
		orderRepoMock.On("FindAll").Once().Return([]model.Order{}, mongo.ErrNoDocuments)
		_, err := orderService.ListOrderByUserID(ctx, &wrapperspb.StringValue{})
		assert.NotEmpty(t, err)
	})
}
