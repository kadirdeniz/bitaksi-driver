package mock

import (
	"driver/internal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Location = internal.Location{
	ID:          primitive.NewObjectID(),
	Coordinates: []float64{1.0, 1.0},
	Distance:    1.0,
}

func BulkCreateDriversRequest(count int) []internal.Location {
	var locations []internal.Location
	for i := 0; i < count; i++ {
		locations = append(locations, Location)
	}
	return locations
}
