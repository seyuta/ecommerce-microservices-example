package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

const (
	// UserAuthCollection is mongo collection name for UserAuth
	UserAuthCollection = "userauth"
	// UserTokenCollection is mongo collection name for UserToken
	UserTokenCollection = "usertoken"
)

type UserStatus string

const (
	UserStatusActive  UserStatus = "active"
	UserStatusSuspend UserStatus = "inactive"
)

// UserAuth ...
type UserAuth struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Username string        `bson:"uname"`
	Email    string        `bson:"email"`
	Phone    string        `bson:"phone"`
	Password string        `bson:"passwd"`
	Status   UserStatus    `bson:"status"`
	Audit
}

// UserToken ...
type UserToken struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	UserID    bson.ObjectId `bson:"uid"`
	Username  string        `bson:"uname"`
	DeviceID  string        `bson:"dvcid"`
	LoginDt   *time.Time    `bson:"lgdt"`
	ExpiredDt *time.Time    `bson:"exdt"`
	Token     string        `bson:"token"`
}
