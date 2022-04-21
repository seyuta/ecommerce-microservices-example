package model

import (
	"github.com/globalsign/mgo/bson"
)

const (
	ProductCollection = "products"
)

const (
	ProductActive   bool = true
	ProductInactive bool = false
)

// Product ...
type Product struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Description string        `bson:"desc"`
	Price       int64         `bson:"price"`
	Stock       int64         `bson:"stock"`
	Status      bool          `bson:"status"`
	Audit
}
