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
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (r UserRepositoryMock) Create(context.Context, *model.UserAuth) (model.UserAuth, error) {
	args := r.Called()
	users := model.UserAuth{
		Username: "test2",
		Email:    "test2@test.com",
	}
	return users, args.Error(1)
}

func (r UserRepositoryMock) CreateToken(context.Context, *model.UserToken) (model.UserToken, error) {
	args := r.Called()
	users := model.UserToken{
		Username: "test",
	}
	return users, args.Error(1)
}

func (r UserRepositoryMock) DeleteTokenByUserAndDevice(context.Context, string, string) (bool, error) {
	args := r.Called()
	return args.Bool(0), args.Error(1)
}

func (r UserRepositoryMock) FindTokenByUserAndDevice(context.Context, string, string) (model.UserToken, error) {
	args := r.Called()
	ex := time.Now().Add(time.Hour * time.Duration(10))
	users := model.UserToken{
		Token:     "token",
		ExpiredDt: &ex,
	}
	return users, args.Error(1)
}

func (r UserRepositoryMock) FindByPhone(context.Context, string) (model.UserAuth, error) {
	args := r.Called()
	users := model.UserAuth{
		Username: "test",
	}
	return users, args.Error(1)
}

func (r UserRepositoryMock) FindByEmail(context.Context, string) (model.UserAuth, error) {
	args := r.Called()
	users := model.UserAuth{
		Email: "test@test.com",
	}
	return users, args.Error(1)
}

func (r UserRepositoryMock) FindByUsername(context.Context, string) (model.UserAuth, error) {
	args := r.Called()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	users := model.UserAuth{
		Username: "test",
		Password: string(hashedPassword),
		Status:   model.UserStatusActive,
	}
	return users, args.Error(1)
}

func (r UserRepositoryMock) FindByID(context.Context, string) (model.UserAuth, error) {
	args := r.Called()
	users := model.UserAuth{
		Username: "test",
	}
	return users, args.Error(1)
}

func TestService_Login(t *testing.T) {
	repository := UserRepositoryMock{}
	repository.On("FindByUsername").Return(model.UserAuth{}, nil)
	repository.On("FindTokenByUserAndDevice").Return(model.UserAuth{}, nil)

	ctx := context.Background()

	service := AuthSvc{repository, logrus.StandardLogger()}
	req := &pb.LoginDto{
		Username: "test",
		Password: "123",
		DeviceId: "123",
	}
	users, _ := service.Login(ctx, req)
	assert.Equal(t, users.Username, "test")
}

func TestService_Register(t *testing.T) {
	repository := UserRepositoryMock{}
	repository.On("FindByUsername").Return(model.UserAuth{}, nil)
	repository.On("FindByEmail").Return(model.UserAuth{}, nil)
	repository.On("Create").Return(model.UserAuth{}, nil)

	ctx := context.Background()

	service := AuthSvc{repository, logrus.StandardLogger()}
	req := &pb.RegisterDto{
		Username: "test2",
		Email:    "test2@test.com",
		Phone:    "123",
		Password: "123",
	}
	users, _ := service.Register(ctx, req)
	assert.Equal(t, users.Username, "test2")
	assert.Equal(t, users.Email, "test2@test.com")
}
