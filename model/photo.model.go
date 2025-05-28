package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Photo represents the structure for a photo in the database
type Photo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	URL         string             `bson:"url" json:"url" validate:"required,url"`
	Title       string             `bson:"title" json:"title" validate:"required,min=1,max=200"`
	Description string             `bson:"description" json:"description" validate:"max=1000"`
	Alt         string             `bson:"alt" json:"alt" validate:"required,max=200"`
	EntityID    primitive.ObjectID `bson:"entity_id" json:"entity_id" validate:"required"`
	EntityType  string             `bson:"entity_type" json:"entity_type" validate:"required,oneof=hotel restaurant cuisine"`
	IsPrimary   bool               `bson:"is_primary" json:"is_primary"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	Order       int                `bson:"order" json:"order" validate:"min=0"`
	Metadata    PhotoMetadata      `bson:"metadata" json:"metadata"`
	Tags        []string           `bson:"tags" json:"tags"`
	UploadedBy  primitive.ObjectID `bson:"uploaded_by" json:"uploaded_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type PhotoMetadata struct {
	Width       int    `bson:"width" json:"width" validate:"min=1"`
	Height      int    `bson:"height" json:"height" validate:"min=1"`
	Size        int64  `bson:"size" json:"size" validate:"min=1"` // Size in bytes
	Format      string `bson:"format" json:"format" validate:"required,oneof=jpg jpeg png webp"`
	ColorSpace  string `bson:"color_space" json:"color_space"`
	Orientation int    `bson:"orientation" json:"orientation"`
	Quality     int    `bson:"quality" json:"quality" validate:"min=1,max=100"`
}

// PhotoCollection returns the name of the MongoDB collection for photos
func (Photo) CollectionName() string {
	return "photos"
}

// Validate performs custom validation on the photo model
func (p *Photo) Validate() error {
	if strings.TrimSpace(p.URL) == "" {
		return errors.New("URL is required")
	}
	if strings.TrimSpace(p.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(p.EntityType) == "" {
		return errors.New("entity type is required")
	}
	if p.EntityID.IsZero() {
		return errors.New("entity ID is required")
	}

	// Validate entity type
	validEntityTypes := []string{"hotel", "restaurant", "cuisine"}
	if !contains(validEntityTypes, p.EntityType) {
		return fmt.Errorf("invalid entity type: %s", p.EntityType)
	}

	// Validate metadata if present
	if p.Metadata.Width > 0 || p.Metadata.Height > 0 {
		if p.Metadata.Width <= 0 || p.Metadata.Height <= 0 {
			return errors.New("both width and height must be positive if metadata is provided")
		}
	}

	return nil
}

// SetDefaults sets default values for the photo
func (p *Photo) SetDefaults() {
	if p.IsActive == false && p.CreatedAt.IsZero() {
		p.IsActive = true
	}
	if p.Order == 0 && p.CreatedAt.IsZero() {
		p.Order = 1
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	p.UpdatedAt = time.Now()

	// Set alt text if empty
	if p.Alt == "" {
		p.Alt = p.Title
	}
}

// GetThumbnailURL generates a thumbnail URL based on the original URL
func (p *Photo) GetThumbnailURL(size string) string {
	if size == "" {
		size = "medium"
	}

	// Simple thumbnail URL generation (would be more sophisticated in real implementation)
	parts := strings.Split(p.URL, ".")
	if len(parts) >= 2 {
		ext := parts[len(parts)-1]
		base := strings.Join(parts[:len(parts)-1], ".")
		return fmt.Sprintf("%s_thumb_%s.%s", base, size, ext)
	}
	return p.URL
}

// GetAspectRatio calculates the aspect ratio of the photo
func (p *Photo) GetAspectRatio() float64 {
	if p.Metadata.Height == 0 {
		return 0
	}
	return float64(p.Metadata.Width) / float64(p.Metadata.Height)
}

// IsLandscape returns true if the photo is in landscape orientation
func (p *Photo) IsLandscape() bool {
	return p.GetAspectRatio() > 1.0
}

// IsPortrait returns true if the photo is in portrait orientation
func (p *Photo) IsPortrait() bool {
	return p.GetAspectRatio() < 1.0
}

// IsSquare returns true if the photo is square
func (p *Photo) IsSquare() bool {
	ratio := p.GetAspectRatio()
	return ratio >= 0.95 && ratio <= 1.05 // Allow small tolerance
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
