package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Photo represents the structure for a photo in the database
type Photo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	URL         string             `bson:"url" json:"url"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	EntityID    primitive.ObjectID `bson:"entity_id" json:"entity_id"`
	EntityType  string             `bson:"entity_type" json:"entity_type"` // "hotel", "restaurant", etc.
	IsPrimary   bool               `bson:"is_primary" json:"is_primary"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// PhotoCollection returns the name of the MongoDB collection for photos
func (Photo) CollectionName() string {
	return "photos"
}
