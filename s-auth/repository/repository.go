package repository

import (
	"context"

	"github.com/seyuta/ecommerce-microservices-example/s-auth/pkg/model"
)

type (
	// Repository is base interface for any kind of Repository
	Repository interface {
		UserRepository() UserRepository
	}

	// UserRepository ...
	UserRepository interface {
		Create(c context.Context, t *model.UserAuth) (model.UserAuth, error)
		CreateToken(c context.Context, t *model.UserToken) (model.UserToken, error)
		DeleteTokenByUserAndDevice(c context.Context, uname, dvcid string) (bool, error)
		FindTokenByUserAndDevice(c context.Context, uname, dvcid string) (model.UserToken, error)
		FindByPhone(c context.Context, mobile string) (model.UserAuth, error)
		FindByEmail(c context.Context, email string) (model.UserAuth, error)
		FindByUsername(c context.Context, uname string) (model.UserAuth, error)
		FindByID(c context.Context, id string) (model.UserAuth, error)
	}
)
