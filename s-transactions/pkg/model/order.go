package model

import (
	"github.com/globalsign/mgo/bson"
)

const (
	OrderCollection       = "order"
	OrderDetailCollection = "orderdetail"
)

const (
	OrderPaid       int = 1
	OrderOnProgress int = 2
	OrderComplete   int = 3
	OrderCancel     int = 4
	OrderRefund     int = 5
)

// Order ...
type Order struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	UserID      string        `bson:"userid"`
	NoInvoice   string        `bson:"noinv"`
	Status      bool          `bson:"status"`
	OrderDetail []OrderDetail
	Audit
}

// OrderDetail ...
type OrderDetail struct {
	ProductID string `bson:"productid"`
	Qty       int64  `bson:"qty"`
}
