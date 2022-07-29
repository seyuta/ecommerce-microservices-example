package service

import (
	"context"
	"testing"
	"time"

	"github.com/seyuta/ecommerce-microservices-example/s-auth/pkg/model"
	"github.com/seyuta/ecommerce-microservices-example/s-auth/pkg/pb"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (r *UserRepositoryMock) Create(ctx context.Context, req *model.UserAuth) (model.UserAuth, error) {
	args := r.Called(ctx, req)
	return args.Get(0).(model.UserAuth), args.Error(1)
}

func (r *UserRepositoryMock) CreateToken(ctx context.Context, req *model.UserToken) (model.UserToken, error) {
	args := r.Called(ctx, req)
	return args.Get(0).(model.UserToken), args.Error(1)
}

func (r *UserRepositoryMock) DeleteTokenByUserAndDevice(ctx context.Context, username, deviceId string) (bool, error) {
	args := r.Called(ctx, username, deviceId)
	return args.Bool(0), args.Error(1)
}

func (r *UserRepositoryMock) FindTokenByUserAndDevice(ctx context.Context, username, deviceId string) (model.UserToken, error) {
	args := r.Called(ctx, username, deviceId)
	return args.Get(0).(model.UserToken), args.Error(1)
}

func (r *UserRepositoryMock) FindByPhone(ctx context.Context, mobile string) (model.UserAuth, error) {
	args := r.Called(ctx, mobile)
	return args.Get(0).(model.UserAuth), args.Error(1)
}

func (r *UserRepositoryMock) FindByEmail(ctx context.Context, email string) (model.UserAuth, error) {
	args := r.Called(ctx, email)
	return args.Get(0).(model.UserAuth), args.Error(1)
}

func (r *UserRepositoryMock) FindByUsername(ctx context.Context, username string) (model.UserAuth, error) {
	args := r.Called(ctx, username)
	return args.Get(0).(model.UserAuth), args.Error(1)
}

func (r *UserRepositoryMock) FindByID(ctx context.Context, id string) (model.UserAuth, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(model.UserAuth), args.Error(1)
}

var (
	// mock
	userRepoMock = UserRepositoryMock{}

	// service
	userService = NewAuthSvc(&userRepoMock, logrus.StandardLogger())

	// context
	ctx = context.Background()
)

func TestService_Register(t *testing.T) {
	t.Run("Scenario Normal", func(t *testing.T) {
		userRepoMock.On("FindByUsername", ctx, mock.AnythingOfType("string")).Once().Return(model.UserAuth{}, nil)
		userRepoMock.On("FindByEmail", ctx, mock.AnythingOfType("string")).Once().Return(model.UserAuth{}, nil)
		userRepoMock.On("Create", ctx, mock.Anything).Once().Return(model.UserAuth{}, nil)

		req := &pb.RegisterDto{
			Username: "usertest",
			Email:    "test@test.mail",
			Phone:    "1234567890",
			Password: "123",
		}
		_, err := userService.Register(ctx, req)
		assert.Empty(t, err)
	})

	t.Run("Scenario Blocker Username", func(t *testing.T) {
		userRepoMock.On("FindByUsername", ctx, mock.AnythingOfType("string")).Once().Return(model.UserAuth{Username: "usertest"}, nil)

		req := &pb.RegisterDto{
			Username: "usertest",
			Email:    "test@test.mail",
			Phone:    "1234567890",
			Password: "123",
		}
		_, err := userService.Register(ctx, req)
		assert.NotEmpty(t, err)
	})

	t.Run("Scenario Blocker Email", func(t *testing.T) {
		userRepoMock.On("FindByUsername", ctx, mock.AnythingOfType("string")).Once().Return(model.UserAuth{}, nil)
		userRepoMock.On("FindByEmail", ctx, mock.AnythingOfType("string")).Once().Return(model.UserAuth{Email: "test@test.mail"}, nil)

		req := &pb.RegisterDto{
			Username: "usertest",
			Email:    "test@test.mail",
			Phone:    "1234567890",
			Password: "123",
		}
		_, err := userService.Register(ctx, req)
		assert.NotEmpty(t, err)
	})
}

func TestService_Login(t *testing.T) {
	t.Run("Scenario Normal", func(t *testing.T) {
		userRepoMock.On("FindByUsername", ctx, mock.AnythingOfType("string")).Once().Return(
			model.UserAuth{
				Password: "$2a$12$jCaKCFEl723czhZmc7fDEeU8CezXZ.Qmram4G5IGrJf/enmEDjzSC",
				Status:   "active",
			},
			nil,
		)
		userRepoMock.On("FindTokenByUserAndDevice", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(
			model.UserToken{
				ExpiredDt: &time.Time{},
			},
			nil,
		)
		userRepoMock.On("DeleteTokenByUserAndDevice", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(true, nil)
		userRepoMock.On("CreateToken", ctx, mock.Anything).Return(
			model.UserToken{
				ExpiredDt: &time.Time{},
			},
			nil,
		)

		req := &pb.LoginDto{
			Username: "usertest",
			Password: "123",
			DeviceId: "123",
		}
		_, err := userService.Login(ctx, req)
		assert.Empty(t, err)
	})

	t.Run("Scenario Blocker Username Empty", func(t *testing.T) {
		req := &pb.LoginDto{Username: ""}
		_, err := userService.Login(ctx, req)
		assert.NotEmpty(t, err)
	})

	t.Run("Scenario Blocker Invalid Username", func(t *testing.T) {
		req := &pb.LoginDto{Username: "+_)"}
		_, err := userService.Login(ctx, req)
		assert.NotEmpty(t, err)
	})

	t.Run("Scenario Blocker DeviceID Empty", func(t *testing.T) {
		req := &pb.LoginDto{
			Username: "usertest",
			DeviceId: "",
		}
		_, err := userService.Login(ctx, req)
		assert.NotEmpty(t, err)
	})

	t.Run("Scenario Username Not Found", func(t *testing.T) {
		userRepoMock.On("FindByUsername", ctx, mock.AnythingOfType("string")).Once().Return(model.UserAuth{}, mongo.ErrNoDocuments)

		req := &pb.LoginDto{
			Username: "usertest",
			Password: "123",
			DeviceId: "123",
		}
		_, err := userService.Login(ctx, req)
		assert.NotEmpty(t, err)
	})

	t.Run("Scenario Invalid Password", func(t *testing.T) {
		userRepoMock.On("FindByUsername", ctx, mock.AnythingOfType("string")).Once().Return(model.UserAuth{}, nil)

		req := &pb.LoginDto{
			Username: "usertest",
			Password: "123",
			DeviceId: "123",
		}
		_, err := userService.Login(ctx, req)
		assert.NotEmpty(t, err)
	})

	t.Run("Scenario Username Not Active", func(t *testing.T) {
		userRepoMock.On("FindByUsername", ctx, mock.AnythingOfType("string")).Once().Return(
			model.UserAuth{Password: "$2a$12$jCaKCFEl723czhZmc7fDEeU8CezXZ.Qmram4G5IGrJf/enmEDjzSC"},
			nil,
		)

		req := &pb.LoginDto{
			Username: "usertest",
			Password: "123",
			DeviceId: "123",
		}
		_, err := userService.Login(ctx, req)
		assert.NotEmpty(t, err)
	})
}
