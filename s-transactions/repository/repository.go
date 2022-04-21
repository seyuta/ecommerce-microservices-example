package repository

import (
	"context"

	"github.com/seyuta/ecommerce-microservices-example/s-transactions/pkg/model"
)

type (
	// Repository is base interface for any kind of Repository
	Repository interface {
		OrderRepository() OrderRepository
	}

	// OrderRepository ...
	OrderRepository interface {
		Create(c context.Context, t *model.Order) (model.Order, error)
		FindByID(id string) (model.Order, error)
		FindAll() ([]model.Order, error)
	}
)
