package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Restaurant represents the structure for a restaurant in the database
type Restaurant struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name          string               `bson:"name" json:"name"`
	Description   string               `bson:"description" json:"description"`
	Address       string               `bson:"address" json:"address"`
	Location      string               `bson:"location" json:"location"`
	Cuisines      []primitive.ObjectID `bson:"cuisines" json:"cuisines"`
	PriceRange    string               `bson:"price_range" json:"price_range"`
	Rating        float64              `bson:"rating" json:"rating"`
	ContactNumber string               `bson:"contact_number" json:"contact_number"`
	Email         string               `bson:"email" json:"email"`
	Website       string               `bson:"website" json:"website"`
	OpeningHours  map[string]string    `bson:"opening_hours" json:"opening_hours"`
	CreatedAt     time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time            `bson:"updated_at" json:"updated_at"`
}

// RestaurantCollection returns the name of the MongoDB collection for restaurants
func (Restaurant) CollectionName() string {
	return "restaurants"
}
