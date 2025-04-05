package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hotel struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title         string             `bson:"title" json:"title"`
	Content       string             `bson:"content" json:"content"`
	PrimaryInfo   string             `bson:"primary_info" json:"primary_info"`
	SecondaryInfo string             `bson:"secondary_info" json:"secondary_info"`
	AccentedLabel string             `bson:"accented_label" json:"accented_label"`
	Provider      string             `bson:"provider" json:"provider"`
	PriceDetails  string             `bson:"price_details" json:"price_details"`
	PriceSummary  string             `bson:"price_summary" json:"price_summary"`
	Location      string             `bson:"location" json:"location"`
	Price         float64            `bson:"price" json:"price"`
	Rating        float64            `bson:"rating" json:"rating"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}
