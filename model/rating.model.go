package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Rating represents the structure for a rating in the database
type Rating struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	HotelID   primitive.ObjectID `bson:"hotel_id" json:"hotel_id"`
	Score     float64            `bson:"score" json:"score"`
	Comment   string             `bson:"comment" json:"comment"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// RatingCollection returns the name of the MongoDB collection for ratings
func (Rating) CollectionName() string {
	return "ratings"
}
