package model

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hotel struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title         string             `bson:"title" json:"title" validate:"required,min=3,max=200"`
	Content       string             `bson:"content" json:"content" validate:"required,min=10"`
	PrimaryInfo   string             `bson:"primary_info" json:"primary_info"`
	SecondaryInfo string             `bson:"secondary_info" json:"secondary_info"`
	AccentedLabel string             `bson:"accented_label" json:"accented_label"`
	Provider      string             `bson:"provider" json:"provider" validate:"required"`
	PriceDetails  string             `bson:"price_details" json:"price_details"`
	PriceSummary  string             `bson:"price_summary" json:"price_summary"`
	Location      Location           `bson:"location" json:"location" validate:"required"`
	Address       string             `bson:"address" json:"address"`
	Price         float64            `bson:"price" json:"price" validate:"required,min=0"`
	Currency      string             `bson:"currency" json:"currency" validate:"required,len=3"`
	Rating        float64            `bson:"rating" json:"rating" validate:"min=0,max=5"`
	ReviewCount   int                `bson:"review_count" json:"review_count"`
	Amenities     []string           `bson:"amenities" json:"amenities"`
	RoomTypes     []RoomType         `bson:"room_types" json:"room_types"`
	ContactInfo   ContactInfo        `bson:"contact_info" json:"contact_info"`
	Policies      HotelPolicies      `bson:"policies" json:"policies"`
	Status        string             `bson:"status" json:"status" validate:"required,oneof=active inactive pending"`
	Tags          []string           `bson:"tags" json:"tags"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type Location struct {
	City        string    `bson:"city" json:"city" validate:"required"`
	State       string    `bson:"state" json:"state"`
	Country     string    `bson:"country" json:"country" validate:"required"`
	PostalCode  string    `bson:"postal_code" json:"postal_code"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"` // [longitude, latitude]
}

type RoomType struct {
	ID          string   `bson:"id" json:"id"`
	Name        string   `bson:"name" json:"name" validate:"required"`
	Description string   `bson:"description" json:"description"`
	Price       float64  `bson:"price" json:"price" validate:"required,min=0"`
	MaxGuests   int      `bson:"max_guests" json:"max_guests" validate:"required,min=1"`
	Amenities   []string `bson:"amenities" json:"amenities"`
	Available   bool     `bson:"available" json:"available"`
}

type ContactInfo struct {
	Phone   string `bson:"phone" json:"phone"`
	Email   string `bson:"email" json:"email" validate:"omitempty,email"`
	Website string `bson:"website" json:"website" validate:"omitempty,url"`
}

type HotelPolicies struct {
	CheckIn        string   `bson:"check_in" json:"check_in"`
	CheckOut       string   `bson:"check_out" json:"check_out"`
	Cancellation   string   `bson:"cancellation" json:"cancellation"`
	PetPolicy      string   `bson:"pet_policy" json:"pet_policy"`
	SmokingPolicy  string   `bson:"smoking_policy" json:"smoking_policy"`
	AgeRestriction int      `bson:"age_restriction" json:"age_restriction"`
	AcceptedCards  []string `bson:"accepted_cards" json:"accepted_cards"`
}

// HotelCollection returns the name of the MongoDB collection for hotels
func (Hotel) CollectionName() string {
	return "hotels"
}

// Validate performs custom validation on the hotel model
func (h *Hotel) Validate() error {
	if strings.TrimSpace(h.Title) == "" {
		return errors.New("title is required")
	}
	if h.Price < 0 {
		return errors.New("price cannot be negative")
	}
	if h.Rating < 0 || h.Rating > 5 {
		return errors.New("rating must be between 0 and 5")
	}
	if len(h.Currency) != 3 {
		return errors.New("currency must be a 3-letter code")
	}
	return nil
}

// SetDefaults sets default values for the hotel
func (h *Hotel) SetDefaults() {
	if h.Status == "" {
		h.Status = "pending"
	}
	if h.Currency == "" {
		h.Currency = "USD"
	}
	if h.CreatedAt.IsZero() {
		h.CreatedAt = time.Now()
	}
	h.UpdatedAt = time.Now()
}

// UpdateRating updates the hotel's average rating based on new review
func (h *Hotel) UpdateRating(newRating float64) {
	if h.ReviewCount == 0 {
		h.Rating = newRating
		h.ReviewCount = 1
	} else {
		totalRating := h.Rating * float64(h.ReviewCount)
		h.ReviewCount++
		h.Rating = (totalRating + newRating) / float64(h.ReviewCount)
	}
	h.UpdatedAt = time.Now()
}

// GetAveragePrice returns the average price across all room types
func (h *Hotel) GetAveragePrice() float64 {
	if len(h.RoomTypes) == 0 {
		return h.Price
	}

	total := 0.0
	for _, room := range h.RoomTypes {
		total += room.Price
	}
	return total / float64(len(h.RoomTypes))
}
