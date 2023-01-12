package mongodb

import (
	"context"
	"driver/internal"
	"driver/pkg"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	GetMongoDBURI() string
	Connect() (*MongoDB, error)
	GetDriverCollection() *mongo.Collection
	GetDatabase() *mongo.Database
	FindNearestDriver(location internal.Coordinates) error
	BulkCreateDrivers(locations []internal.Location) error
	FlushUsers() error
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

func (m *MongoDB) FlushUsers() error {
	return nil
}

func (m *MongoDB) GetMongoDBURI() string {
	return "mongodb://" + m.Username + ":" + m.Password + "@" + m.Host + ":" + m.Port
}

func (m *MongoDB) Connect() (*MongoDB, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(m.GetMongoDBURI()))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	MongoClient = client.Database(m.DBName)

	return m, nil
}

func (m *MongoDB) GetDriverCollection() *mongo.Collection {
	return MongoClient.Collection(DriverCollection)
}

func (m *MongoDB) GetDatabase() *mongo.Database {
	return MongoClient
}
