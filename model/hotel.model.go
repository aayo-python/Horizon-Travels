package model

import (
	"time"
)

type Hotel struct {
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	PrimaryInfo   string    `json:"primaryinfo"`
	SecondaryInfo string    `json:"secondaryinfo"`
	AccentedLabel string    `json:"accentedlabel"`
	Provider      string    `json:"provider"`
	PriceDetails  string    `json:"pricedetails"`
	PriceSummary  string    `json:"pricesummary"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
