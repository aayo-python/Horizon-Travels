package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Thumbnail represents the structure for a thumbnail in the database
type Thumbnail struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	URL        string             `bson:"url" json:"url"`
	Width      int                `bson:"width" json:"width"`
	Height     int                `bson:"height" json:"height"`
	EntityID   primitive.ObjectID `bson:"entity_id" json:"entity_id"`
	EntityType string             `bson:"entity_type" json:"entity_type"` // "hotel", "restaurant", "photo", etc.
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

// ThumbnailCollection returns the name of the MongoDB collection for thumbnails
func (Thumbnail) CollectionName() string {
	return "thumbnails"
}
