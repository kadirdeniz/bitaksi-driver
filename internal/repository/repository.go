package repository

import (
	"driver/internal"
	"driver/tools/mongodb"
)

type Repository struct {
	MongoDBInterface mongodb.MongoDBInterface
}

//go:generate mockgen -source=repository.go -destination=./../../test/mock/mock_repository.go -package=mock
type IRepository interface {
	FindNearestDriver(location internal.Coordinates) (internal.Location, error)
	BulkCreateDrivers(locations []internal.Location) error
	Migration() error
}

func NewRepository(mongoDBInterface mongodb.MongoDBInterface) IRepository {
	return &Repository{
		MongoDBInterface: mongoDBInterface,
	}
}

func (r *Repository) FindNearestDriver(location internal.Coordinates) (internal.Location, error) {
	var result internal.Location
	err := r.MongoDBInterface.FindNearestDriver(location)
	return result, err
}

func (r *Repository) BulkCreateDrivers(locations []internal.Location) error {
	return r.MongoDBInterface.BulkCreateDrivers(locations)
}

func (r *Repository) Migration() error {
	err := r.MongoDBInterface.FlushLocations()
	if err != nil {
		return err
	}
	err = r.MongoDBInterface.CreateIndexForGeoJSON()
	if err != nil {
		return err
	}
	err = r.MongoDBInterface.MigrateWithCSVData()
	if err != nil {
		return err
	}

	return nil
}
