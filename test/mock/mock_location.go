package mock

import (
	"driver/internal"
)

var Location = internal.Location{
	Coordinates: []float64{1.0, 1.0},
}

func BulkCreateDriversRequest(count int) []internal.Location {
	var locations []internal.Location
	for i := 0; i < count; i++ {
		locations = append(locations, Location)
	}
	return locations
}
