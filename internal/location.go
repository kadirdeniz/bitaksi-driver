package internal

import "go.mongodb.org/mongo-driver/bson/primitive"

type Location struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Type        string             `json:"type,omitempty" bson:"type,omitempty"`
	Coordinates []float64          `json:"coordinates" bson:"coordinates"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
