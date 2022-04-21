package repository

import (
	"context"
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/seyuta/ecommerce-microservices-example/s-auth/pkg/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserAuthRepo repository
type UserAuthRepo struct {
	db  *mongo.Database
	col *mongo.Collection
	log *logrus.Logger
}

// NewUserAuthRepo is new instantiation of UserAuthRepo
func NewUserAuthRepo(db *mongo.Database, log *logrus.Logger) *UserAuthRepo {
	return &UserAuthRepo{
		db:  db,
		col: db.Collection(model.UserAuthCollection),
		log: log,
	}
}

// Create new UserAuth
func (r *UserAuthRepo) Create(c context.Context, t *model.UserAuth) (model.UserAuth, error) {
	r.log.Debugf("Create(%s)", t)
	ctx, cancel := TimeOutContextWithParent(c)
	defer cancel()

	var uauth model.UserAuth
	now := bson.Now()
	t.Audit.CreatedDt = &now
	res, err := r.col.InsertOne(ctx, t)
	if err != nil {
		r.log.Errorf(ErrMgoOpsFail, err)
		return uauth, err
	}

	err = r.col.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&uauth)
	if err != nil {
		r.log.Errorf(ErrMgoOpsFail, err)
		return uauth, err
	}

	r.log.Debugf("Saved UserAuth(%v)", uauth)

	return uauth, nil
}

// CreateToken new UserToken for UserAuth
func (r *UserAuthRepo) CreateToken(c context.Context, t *model.UserToken) (model.UserToken, error) {
	r.log.Debugf("CreateToken(%s)", t)
	ctx, cancel := TimeOutContextWithParent(c)
	defer cancel()

	col := r.db.Collection(model.UserTokenCollection)

	var utoken model.UserToken
	res, err := col.InsertOne(ctx, t)
	if err != nil {
		r.log.Errorf(ErrMgoOpsFail, err)
		return utoken, err
	}

	err = col.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&utoken)
	if err != nil {
		r.log.Errorf(ErrMgoOpsFail, err)
		return utoken, err
	}

	r.log.Debugf("Saved UserToken(%v)", utoken)

	return utoken, err
}

// DeleteTokenByUserAndDevice is removing UserToken by Username & DeviceID
func (r *UserAuthRepo) DeleteTokenByUserAndDevice(c context.Context, uname, dvcid string) (bool, error) {
	r.log.Debugf("DeleteTokenByUserAndDevice(%s,%s)", uname, dvcid)
	ctx, cancel := TimeOutContextWithParent(c)
	defer cancel()

	col := r.db.Collection(model.UserTokenCollection)
	filter := bson.M{
		"uname": uname,
		"dvcid": dvcid,
	}

	res, err := col.DeleteOne(ctx, filter)
	if err != nil {
		r.log.Errorf(ErrMgoOpsFail, err)
		return false, err
	}

	return res.DeletedCount > 0, nil
}

// FindTokenByUserAndDevice is finding UserToken by Username & DeviceID which not expired yet
func (r *UserAuthRepo) FindTokenByUserAndDevice(c context.Context, uname, dvcid string) (model.UserToken, error) {
	r.log.Debugf("FindTokenByUserAndDevice(%s,%s)", uname, dvcid)
	ctx, cancel := TimeOutContextWithParent(c)
	defer cancel()

	col := r.db.Collection(model.UserTokenCollection)
	filter := bson.M{
		"uname": uname,
		"dvcid": dvcid,
		"exdt": bson.M{
			"$gt": bson.Now(),
		},
	}

	var utoken model.UserToken
	err := col.FindOne(ctx, filter).Decode(&utoken)
	if err != nil {
		r.log.Errorf(ErrMgoOpsFail, err)
		return utoken, err
	}

	return utoken, nil
}

// FindByMobile is finding UserAuth by its Mobile
func (r *UserAuthRepo) FindByPhone(c context.Context, mobile string) (model.UserAuth, error) {
	r.log.Debugf("FindByMobile(%s) \n", mobile)
	ctx, cancel := TimeOutContextWithParent(c)
	defer cancel()

	var uauth model.UserAuth
	err := r.col.FindOne(ctx, bson.M{"mobile": mobile}).Decode(&uauth)
	if err != nil {
		r.log.Errorf(ErrMgoOpsFail, err)
		return uauth, err
	}

	return uauth, nil
}

// FindByEmail is finding UserAuth by its Email
func (r *UserAuthRepo) FindByEmail(c context.Context, email string) (model.UserAuth, error) {
	r.log.Debugf("FindByEmail(%s) \n", email)
	ctx, cancel := TimeOutContextWithParent(c)
	defer cancel()

	var uauth model.UserAuth
	err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&uauth)
	if err != nil {
		r.log.Errorf(ErrMgoOpsFail, err)
		return uauth, err
	}

	return uauth, nil
}

// FindByUsername is finding UserAuth by its Username
func (r *UserAuthRepo) FindByUsername(c context.Context, uname string) (model.UserAuth, error) {
	fmt.Println(r)
	r.log.Debugf("FindByUsername(%s) \n", uname)
	ctx, cancel := TimeOutContextWithParent(c)
	defer cancel()

	var uauth model.UserAuth
	err := r.col.FindOne(ctx, bson.M{"uname": uname}).Decode(&uauth)
	if err != nil {
		r.log.Errorf(ErrMgoOpsFail, err)
		return uauth, err
	}

	return uauth, nil
}

// FindByID is finding UserAuth by its ID
func (r *UserAuthRepo) FindByID(c context.Context, id string) (model.UserAuth, error) {
	r.log.Debugf("FindByID(%s) \n", id)
	ctx, cancel := TimeOutContextWithParent(c)
	defer cancel()

	var uauth model.UserAuth
	oid, _ := primitive.ObjectIDFromHex(id)
	err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&uauth)
	if err != nil {
		r.log.Errorf(ErrMgoOpsFail, err)
		return uauth, err
	}

	return uauth, nil
}
