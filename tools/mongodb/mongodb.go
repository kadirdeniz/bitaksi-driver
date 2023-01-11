package mongodb

import (
	"context"
	"driver/internal"
	"driver/pkg"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Database

const DriverCollection = "drivers"

var CTX = context.TODO()

type MongoDB struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
}

type MongoDBInterface interface {
	FindNearestDriver(location internal.Coordinates) error
	BulkCreateDrivers(locations []internal.Location) error
}

func NewMongoDB(config pkg.MongoDBConfig) MongoDBInterface {
	return &MongoDB{
		Username: config.Username,
		Password: config.Password,
		Host:     config.Host,
		Port:     config.Port,
		DBName:   config.Database,
	}
}

func (m *MongoDB) FindNearestDriver(location internal.Coordinates) error {
	return nil
}

func (m *MongoDB) BulkCreateDrivers(locations []internal.Location) error {
	return nil
}
