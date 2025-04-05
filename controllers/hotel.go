package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/janto-pee/Horizon-Travels.git/model"
	"github.com/janto-pee/Horizon-Travels.git/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Common success response function
func successResponse(data interface{}, pagination interface{}) gin.H {
	response := gin.H{"data": data}
	if pagination != nil {
		response["pagination"] = pagination
	}
	return response
}

// Common pagination struct
type Pagination struct {
	PageID   int64 `form:"page_id" binding:"required,min=1"`
	PageSize int64 `form:"page_size" binding:"required,min=5,max=100"`
}

// List Hotels with pagination
type ListHotelsRequest struct {
	Pagination
}

func ListHotels(c *gin.Context) {
	var req ListHotelsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	skip := (req.PageID - 1) * req.PageSize
	opts := options.Find().
		SetLimit(req.PageSize).
		SetSkip(skip).
		SetSort(bson.D{{"created_at", -1}})

	// Count total documents for better pagination
	total, err := util.Db.CountDocuments(ctx, bson.D{{}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cursor, err := util.Db.Find(ctx, bson.D{{}}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer cursor.Close(ctx)

	var hotels []model.Hotel
	if err = cursor.All(ctx, &hotels); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, successResponse(hotels, gin.H{
		"page_id":   req.PageID,
		"page_size": req.PageSize,
		"count":     len(hotels),
		"total":     total,
		"pages":     (total + req.PageSize - 1) / req.PageSize,
	}))
}

// Get Hotel by ID
func GetHotelByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid ID format")))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var hotel model.Hotel
	err = util.Db.FindOne(ctx, bson.D{{"_id", id}}).Decode(&hotel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, errorResponse(errors.New("hotel not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, hotel)
}

// Search Hotels
type SearchHotelsRequest struct {
	PageID   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=100"`
	Keyword  string `form:"keyword" binding:"required,min=1"`
}

func SearchHotels(c *gin.Context) {
	var req SearchHotelsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"$text", bson.D{{"$search", req.Keyword}}}}
	skip := (req.PageID - 1) * req.PageSize

	// Count total matching documents
	total, err := util.Db.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	opts := options.Find().
		SetLimit(int64(req.PageSize)).
		SetSkip(int64(skip)).
		// Add projection for text search score
		SetProjection(bson.M{"score": bson.M{"$meta": "textScore"}}).
		// Sort by text search score
		SetSort(bson.M{"score": bson.M{"$meta": "textScore"}})

	cursor, err := util.Db.Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer cursor.Close(ctx)

	var hotels []model.Hotel
	if err = cursor.All(ctx, &hotels); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, successResponse(hotels, gin.H{
		"page_id":   req.PageID,
		"page_size": req.PageSize,
		"count":     len(hotels),
		"total":     total,
		"pages":     (total + int64(req.PageSize) - 1) / int64(req.PageSize),
	}))
}

// Filter Hotels
type FilterHotelsRequest struct {
	PageID    int64   `form:"page_id" binding:"required,min=1"`
	PageSize  int64   `form:"page_size" binding:"required,min=5,max=100"`
	Provider  string  `form:"provider"`
	Location  string  `form:"location"`
	MinPrice  float64 `form:"min_price"`
	MaxPrice  float64 `form:"max_price"`
	MinRating float64 `form:"min_rating"`
}

func FilterHotels(c *gin.Context) {
	var req FilterHotelsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build filter
	filter := bson.D{}

	if req.Provider != "" {
		filter = append(filter, bson.E{"provider", bson.D{{"$regex", req.Provider}, {"$options", "i"}}})
	}

	if req.Location != "" {
		filter = append(filter, bson.E{"location", bson.D{{"$regex", req.Location}, {"$options", "i"}}})
	}

	if req.MinPrice > 0 || req.MaxPrice > 0 {
		priceFilter := bson.D{}
		if req.MinPrice > 0 {
			priceFilter = append(priceFilter, bson.E{"$gte", req.MinPrice})
		}
		if req.MaxPrice > 0 {
			priceFilter = append(priceFilter, bson.E{"$lte", req.MaxPrice})
		}
		filter = append(filter, bson.E{"price", priceFilter})
	}

	if req.MinRating > 0 {
		filter = append(filter, bson.E{"rating", bson.D{{"$gte", req.MinRating}}})
	}

	// Count total matching documents
	total, err := util.Db.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	skip := (req.PageID - 1) * req.PageSize
	opts := options.Find().
		SetLimit(req.PageSize).
		SetSkip(skip).
		SetSort(bson.D{{"rating", -1}})

	cursor, err := util.Db.Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer cursor.Close(ctx)

	var hotels []model.Hotel
	if err = cursor.All(ctx, &hotels); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, successResponse(hotels, gin.H{
		"page_id":   req.PageID,
		"page_size": req.PageSize,
		"count":     len(hotels),
		"total":     total,
		"pages":     (total + req.PageSize - 1) / req.PageSize,
	}))
}

// Aggregate Hotels
func AggregateHotels(c *gin.Context) {
	var pipeline []bson.M
	if err := c.ShouldBindJSON(&pipeline); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cursor, err := util.Db.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, result)
}

// Create Hotel
type CreateHotelRequest struct {
	Title         string  `json:"title" binding:"required"`
	Content       string  `json:"content" binding:"required"`
	PrimaryInfo   string  `json:"primary_info"`
	SecondaryInfo string  `json:"secondary_info"`
	AccentedLabel string  `json:"accented_label"`
	Provider      string  `json:"provider" binding:"required"`
	PriceDetails  string  `json:"price_details"`
	PriceSummary  string  `json:"price_summary"`
	Price         float64 `json:"price" binding:"required"`
	Location      string  `json:"location" binding:"required"`
	Rating        float64 `json:"rating"`
}

func CreateHotel(c *gin.Context) {
	var req CreateHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	hotel := model.Hotel{
		ID:            primitive.NewObjectID(),
		Title:         req.Title,
		Content:       req.Content,
		PrimaryInfo:   req.PrimaryInfo,
		SecondaryInfo: req.SecondaryInfo,
		AccentedLabel: req.AccentedLabel,
		Provider:      req.Provider,
		PriceDetails:  req.PriceDetails,
		PriceSummary:  req.PriceSummary,
		Price:         req.Price,
		Location:      req.Location,
		Rating:        req.Rating,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	result, err := util.Db.InsertOne(ctx, hotel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	hotel.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, hotel)
}

// Update Hotel
type UpdateHotelRequest struct {
	Title         string  `json:"title"`
	Content       string  `json:"content"`
	PrimaryInfo   string  `json:"primary_info"`
	SecondaryInfo string  `json:"secondary_info"`
	AccentedLabel string  `json:"accented_label"`
	Provider      string  `json:"provider"`
	PriceDetails  string  `json:"price_details"`
	PriceSummary  string  `json:"price_summary"`
	Location      string  `json:"location"`
	Rating        float64 `json:"rating"`
}

func UpdateHotel(c *gin.Context) {
	idStr := c.Param("id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req UpdateHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check if hotel exists
	var existingHotel model.Hotel
	err = util.Db.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&existingHotel)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
		return
	}

	// Build update document
	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}

	if req.Title != "" {
		update["$set"].(bson.M)["title"] = req.Title
	}
	if req.Content != "" {
		update["$set"].(bson.M)["content"] = req.Content
	}
	if req.PrimaryInfo != "" {
		update["$set"].(bson.M)["primary_info"] = req.PrimaryInfo
	}
	if req.SecondaryInfo != "" {
		update["$set"].(bson.M)["secondary_info"] = req.SecondaryInfo
	}
	if req.AccentedLabel != "" {
		update["$set"].(bson.M)["accented_label"] = req.AccentedLabel
	}
	if req.Provider != "" {
		update["$set"].(bson.M)["provider"] = req.Provider
	}
	if req.PriceDetails != "" {
		update["$set"].(bson.M)["price_details"] = req.PriceDetails
	}
	if req.PriceSummary != "" {
		update["$set"].(bson.M)["price_summary"] = req.PriceSummary
	}
	if req.Location != "" {
		update["$set"].(bson.M)["location"] = req.Location
	}
	if req.Rating != 0 {
		update["$set"].(bson.M)["rating"] = req.Rating
	}

	result, err := util.Db.UpdateByID(context.TODO(), id, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No changes applied"})
		return
	}

	// Get updated hotel
	var updatedHotel model.Hotel
	err = util.Db.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&updatedHotel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, updatedHotel)
}

// Delete Hotel
func DeleteHotel(c *gin.Context) {
	idStr := c.Param("id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := util.Db.DeleteOne(c, bson.M{"id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete rating"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rating not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "hotel deleted"})
}
