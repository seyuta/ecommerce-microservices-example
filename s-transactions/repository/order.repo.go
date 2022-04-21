package repository

import (
	"context"

	"github.com/globalsign/mgo/bson"
	"github.com/seyuta/ecommerce-microservices-example/s-transactions/pkg/model"
	util "github.com/seyuta/ecommerce-microservices-example/utils/context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Order repository
type OrderRepo struct {
	db  *mongo.Database
	col *mongo.Collection
	log *logrus.Logger
}

// Instantiate new OrderRepo
func NewOrderRepo(db *mongo.Database, log *logrus.Logger) *OrderRepo {
	return &OrderRepo{
		db:  db,
		col: db.Collection(model.OrderCollection),
		log: log,
	}
}

// Create new Order
func (r *OrderRepo) Create(c context.Context, t *model.Order) (model.Order, error) {
	r.log.Debugf("Create(%s)", t)
	ctx, cancel := TimeOutContextWithParent(c)
	defer cancel()

	var (
		order model.Order
		now   = bson.Now()
	)

	t.Audit.CreatedDt = &now
	t.UserID = util.GetUserID(c)
	res, err := r.col.InsertOne(ctx, t)
	if err != nil {
		r.log.Errorf(ErrMgoOpsFail, err)
		return order, err
	}

	err = r.col.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&order)
	if err != nil {
		r.log.Errorf(ErrMgoOpsFail, err)
		return order, err
	}

	return order, err
}

// FindByID is finding order by its id
func (r *OrderRepo) FindByID(id string) (model.Order, error) {
	r.log.Debugf("FindByID(%s) \n", id)
	ctx, cancel := TimeOutContext()
	defer cancel()

	var order model.Order
	oid, _ := primitive.ObjectIDFromHex(id)
	err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&order)
	if err != nil {
		r.log.Println(err)
		return order, err
	}

	return order, nil
}

// FindAll orders
func (r *OrderRepo) FindAll() ([]model.Order, error) {
	r.log.Debugf("FindAll() \n")
	ctx, cancel := TimeOutContext()
	defer cancel()

	var orders []model.Order
	opts := options.Find()
	opts.SetSort(bson.M{"no": 1})

	rcur, err := r.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		r.log.Println(err)
		return orders, err
	}
	defer rcur.Close(ctx)
	for rcur.Next(ctx) {
		var order model.Order
		err := rcur.Decode(&order)
		if err != nil {
			r.log.Errorf(ErrMgoOpsFail, err)
		}
		orders = append(orders, order)
	}
	if err := rcur.Err(); err != nil {
		r.log.Error(err)
	}

	return orders, nil
}
