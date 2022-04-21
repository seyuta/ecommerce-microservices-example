package repository

import (
	"context"
	"time"

	"github.com/seyuta/ecommerce-microservices-example/constant"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepository implementation of `service.Repository`
type MongoRepository struct {
	DB          *mongo.Database
	productRepo *ProductRepo
}

func BuildMongoRepository(db *mongo.Database) *MongoRepository {
	logger := logrus.StandardLogger()
	return &MongoRepository{
		DB:          db,
		productRepo: NewProductRepo(db, logger),
	}
}

func (m *MongoRepository) ProductRepository() ProductRepository {
	return m.productRepo
}

func InitMongoDB(logger *logrus.Logger) *mongo.Database {
	logger.Println("Creating connection to database...")

	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString(constant.CfgMongoURI)))
	if err != nil {
		logger.Fatalf("%v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db := client.Database(viper.GetString(constant.CfgMongoDatabase))
	err = client.Connect(ctx)
	if err != nil {
		logger.Fatalf("%v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatalf("%v", err)
	}

	logger.Println("Connected to database...")

	return db
}
