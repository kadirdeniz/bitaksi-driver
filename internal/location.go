package internal

import "go.mongodb.org/mongo-driver/bson/primitive"

type Model struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Location `bson:"location" json:"location"`
	Distance float64 `bson:"omitempty" json:"distance"`
}

type Location struct {
	Type        string    `json:"type,omitempty" bson:"type,omitempty"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
