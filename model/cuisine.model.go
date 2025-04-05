package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Cuisine represents the structure for a cuisine in the database
type Cuisine struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Origin      string             `bson:"origin" json:"origin"`
	Category    string             `bson:"category" json:"category"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// CuisineCollection returns the name of the MongoDB collection for cuisines
func (Cuisine) CollectionName() string {
	return "cuisines"
}
