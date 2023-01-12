package repository

import (
	"driver/internal"
	"driver/pkg"
	"driver/tools/mongodb"
	"driver/tools/zap"
)

type Repository struct {
	MongoDBInterface mongodb.MongoDBInterface
}

//go:generate mockgen -source=repository.go -destination=./../../test/mock/mock_repository.go -package=mock
type IRepository interface {
	FindNearestDriver(location internal.Coordinates) (internal.Model, error)
	BulkCreateDrivers(locations []internal.Model) error
	Migration() error
}

func NewRepository(mongoDBInterface mongodb.MongoDBInterface) IRepository {
	return &Repository{
		MongoDBInterface: mongoDBInterface,
	}
}

func (r *Repository) FindNearestDriver(location internal.Coordinates) (internal.Model, error) {
	var result internal.Model
	response, err := r.MongoDBInterface.FindNearestDriver(location)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return result, pkg.ErrDriverNotFound
		}

		zap.Logger.Error("Error while finding nearest driver" + err.Error())
		return result, err
	}

	return response, nil
}

func (r *Repository) BulkCreateDrivers(locations []internal.Model) error {
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
