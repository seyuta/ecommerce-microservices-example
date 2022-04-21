package repository

import (
	"context"
	"time"

	"github.com/spf13/viper"
)

// Err
const (
	ErrMgoInvalidOID = "invalid ObjectId: %s"
	ErrMgoOpsFail    = "fail MongoOperation : %s"
)

// TimeOutContext is a global context for mongodb connection access
func TimeOutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(viper.GetInt64("app.mongodb.timeoutinsecond"))*time.Second)
}

// TimeOutContextWithParent is a global context for mongodb connection access with parent Context
func TimeOutContextWithParent(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, time.Duration(viper.GetInt64("app.mongodb.timeoutinsecond"))*time.Second)
}
