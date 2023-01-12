package mock

import (
	"driver/internal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Location = internal.Location{
	Type:        "Point",
	Coordinates: []float64{40.94289771, 29.0390297},
}

var LocationModel = internal.Model{
	ID:       primitive.NewObjectID(),
	Location: Location,
}

var Coordinates = internal.Coordinates{
	Latitude:  40.94289771,
	Longitude: 29.0390297,
}

func BulkCreateDriversRequest(count int) []internal.Model {
	var locations []internal.Model
	for i := 0; i < count; i++ {
		var locationModel = internal.Model{
			ID:       primitive.NewObjectID(),
			Location: Location,
		}

		locations = append(locations, locationModel)
	}
	return locations
}
