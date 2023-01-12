package mongodb

import (
	"context"
	"driver/internal"
	"driver/pkg"
	"driver/tools/zap"
	"go.mongodb.org/mongo-driver/bson"
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
	FlushLocations() error
	FindLocations() ([]*internal.Location, error)
	CreateIndexForGeoJSON() error
	MigrateWithCSVData() error
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
	filter := bson.M{
		"location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{location.Longitude, location.Latitude},
				},
			},
		},
	}

	if err := m.GetDriverCollection().FindOne(CTX, filter).Decode(&location); err != nil {
		return err
	}

	return nil
}

func (m *MongoDB) BulkCreateDrivers(locations []internal.Location) error {
	var drivers []interface{}
	for _, location := range locations {
		drivers = append(drivers, location)
	}

	if _, err := m.GetDriverCollection().InsertMany(CTX, drivers); err != nil {
		return err
	}

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

func (m *MongoDB) FlushLocations() error {
	if _, err := m.GetDriverCollection().DeleteMany(CTX, bson.M{}); err != nil {
		zap.Logger.Error("Error while flushing locations :" + err.Error())
		return err
	}

	return nil
}

func (m *MongoDB) FindLocations() ([]*internal.Location, error) {
	var locations []*internal.Location

	cursor, err := m.GetDriverCollection().Find(CTX, bson.M{})
	if cursor.RemainingBatchLength() == 0 {
		return locations, pkg.ErrDriverNotFound
	}
	if err != nil {
		return nil, err
	}

	if err = cursor.All(CTX, &locations); err != nil {
		return nil, err
	}

	return locations, nil
}

func (m *MongoDB) CreateIndexForGeoJSON() error {
	index := mongo.IndexModel{
		Keys: bson.M{
			"location": "2dsphere",
		},
	}

	if _, err := m.GetDriverCollection().Indexes().CreateOne(CTX, index); err != nil {
		return err
	}

	return nil
}

func (m *MongoDB) MigrateWithCSVData() error {
	const csvFileCount = 101000

	for i := 1; i < csvFileCount; i = i + 1000 {
		locations, err := pkg.GetLocationsFromCSVWithGivenRange(i, i+1000)
		if err != nil {
			return err
		}

		if err := m.BulkCreateDrivers(locations); err != nil {
			return err
		}

	}

	return nil
}
