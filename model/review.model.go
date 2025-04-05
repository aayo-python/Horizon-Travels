package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Review represents the structure for a review in the database
type Review struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	EntityID   primitive.ObjectID `bson:"entity_id" json:"entity_id"`
	EntityType string             `bson:"entity_type" json:"entity_type"` // "hotel", "restaurant", etc.
	Title      string             `bson:"title" json:"title"`
	Content    string             `bson:"content" json:"content"`
	Rating     float64            `bson:"rating" json:"rating"`
	Helpful    int                `bson:"helpful" json:"helpful"`
	NotHelpful int                `bson:"not_helpful" json:"not_helpful"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

// ReviewCollection returns the name of the MongoDB collection for reviews
func (Review) CollectionName() string {
	return "reviews"
}
