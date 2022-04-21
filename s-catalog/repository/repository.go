package repository

import (
	"context"

	"github.com/seyuta/ecommerce-microservices-example/s-catalog/pkg/model"
)

type (
	// Repository is base interface for any kind of Repository
	Repository interface {
		ProductRepository() ProductRepository
	}

	// ProductRepository ...
	ProductRepository interface {
		Create(c context.Context, t *model.Product) (model.Product, error)
		FindByID(id string) (model.Product, error)
		FindAll() ([]model.Product, error)
	}
)
